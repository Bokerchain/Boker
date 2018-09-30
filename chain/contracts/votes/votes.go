package votes

//go: 生成abi和bin文件 solc BokerVerifyVote.sol BokerVerifyVoteImpl.sol BokerVerifyVoteData.sol BokerCommon.sol
//go: 生成go文件 abigen --abi BokerVerifyVote.sol:BokerVerifyVote.abi --bin BokerVerifyVote.sol:BokerVerifyVote.bin  --pkg votes --out contract.go

import (
	"context"
	"math/big"
	"time"

	"github.com/boker/chain/accounts/abi/bind"
	"github.com/boker/chain/boker/protocol"
	"github.com/boker/chain/common"
	"github.com/boker/chain/eth"
	"github.com/boker/chain/log"
)

//投票服务
type VerifyVotesService struct {
	currentEpoch *big.Int        //当前周期序号
	votes        *Votes          //分币session
	addr         common.Address  //合约地址
	ethereum     *eth.Ethereum   //以太坊对象
	tickQuit     chan chan error //tick退出chan
	epochQuit    chan chan error //epoch退出chan
	start        bool            //是否已经启动
}

//创建一个新服务来定期执行
func NewVerifyVotesService(ethereum *eth.Ethereum, address common.Address) (*VerifyVotesService, error) {

	//创建投票对象
	var votesService *VerifyVotesService = new(VerifyVotesService)
	votes, err := NewVotes(address, eth.NewContractBackend(ethereum.ApiBackend))
	if err != nil {
		return nil, err
	}
	votesService.votes = votes
	votesService.addr = address
	votesService.ethereum = ethereum
	votesService.tickQuit = make(chan chan error)
	votesService.epochQuit = make(chan chan error)
	votesService.start = false

	return votesService, nil
}

func (votesService *VerifyVotesService) Start() {

	votesService.start = true
	go votesService.tick()
	go votesService.getEpoch()
}

func (votesService *VerifyVotesService) Stop() error {

	votesService.start = false

	errTick := make(chan error)
	votesService.tickQuit <- errTick

	errEpoch := make(chan error)
	votesService.epochQuit <- errEpoch

	return <-errEpoch
}

func (votesService *VerifyVotesService) createTransactOpts() *bind.TransactOpts {

	if coinbase, err := votesService.ethereum.Coinbase(); err == nil {
		return bind.NewPasswordTransactor(votesService.ethereum, coinbase)
	}
	return nil
}

//产生tick时钟
func (votesService *VerifyVotesService) tick() {

	timer := time.NewTimer(protocol.VotesInterval * 1)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			votesService.tickVotes()
			timer.Reset(protocol.VotesInterval * 1)
		case errc := <-votesService.tickQuit:
			errc <- nil
			return
		}
	}
}

func (votesService *VerifyVotesService) tickVotes() {

	//判断出块节点是否是当前节点
	if err := votesService.getCurrentProducer(); err != nil {
		log.Error("Failed to getCurrentProducer", "err", err)
		return
	}

	//调用时钟函数，判断是否周期发生改变
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	callOpts := &bind.CallOpts{Context: ctx}
	defer cancel()

	epochBool, err := votesService.votes.TickVote(callOpts)
	if err != nil {
		if err == bind.ErrNoCode {
			log.Error("tickVote method not found", "Contract", votesService.addr)
		} else {
			log.Error("Failed to tickVote", "err", err)
		}
		return
	} else {

		//调用转换票数函数
		if epochBool {

			opts := votesService.createTransactOpts()
			_, err := votesService.votes.RotateVote(opts)
			if err != nil {
				if err == bind.ErrNoCode {
					log.Error("rotateVote method not found", "Contract", votesService.addr)
				} else {
					log.Error("Failed to rotateVote", "err", err)
				}
				return
			}
		}
	}
}

//定期获取周期
func (votesService *VerifyVotesService) getEpoch() {

	timer := time.NewTimer(protocol.VotesInterval * 5)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			votesService.tickEpoch()
			timer.Reset(protocol.VotesInterval * 5)

		case errc := <-votesService.epochQuit:
			errc <- nil
			return
		}
	}
}

func (votesService *VerifyVotesService) tickEpoch() {

	//调用是否周期发生改变
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	callOpts := &bind.CallOpts{Context: ctx}
	defer cancel()

	epochIndex, err := votesService.votes.GetVoteRound(callOpts)
	if err != nil {
		if err == bind.ErrNoCode {
			log.Error("GetVoteRound method not found", "Contract", votesService.addr)
		} else {
			log.Error("Failed to GetVoteRound", "err", err)
		}
		return
	} else {

		//判断轮数是否和当前记录的是否一致，如果不一致，则重新获取数据
		if epochIndex != votesService.currentEpoch && epochIndex.Int64() == votesService.currentEpoch.Int64()+1 {

			//调用获取候选人列表数据
			candidateArray, err := votesService.votes.GetCandidates(callOpts)
			if err != nil {
				if err == bind.ErrNoCode {
					log.Error("GetCandidates method not found", "Contract", votesService.addr)
				} else {
					log.Error("Failed to GetCandidates", "err", err)
				}
				return
			} else {

				//判断候选人数组长度和候选人票数数组长度是否一致
				if len(candidateArray.Tickets) != len(candidateArray.Addresses) {
					log.Error("Failed Candidate Address Size Not equal to Candidate Vote Size", err)
					return
				}

				//检索以太坊服务依赖项以访问区块链
				if votesService.ethereum != nil {
					votesService.ethereum.BlockChain().CurrentBlock().DposCtx().SetValidatorVotes(candidateArray.Addresses, candidateArray.Tickets)
				}
			}

		}
	}
}

func (votesService *VerifyVotesService) GetVotesAddr() common.Address {

	return votesService.addr
}

func (votesService *VerifyVotesService) getCurrentProducer() error {

	if votesService.ethereum != nil {

		//得到当前的出块节点
		producer, err := votesService.ethereum.BlockChain().CurrentBlock().DposCtx().GetCurrentProducer()
		if err != nil {
			return protocol.ErrInvalidProducer
		}

		//得到当前挖矿节点
		coinbase, err := votesService.ethereum.Coinbase()
		if err != nil {
			return protocol.ErrInvalidCoinbase
		}

		//将当前出块节点和当前节点进行比较，如果是当前出块节点，则允许继续进行处理
		if producer == coinbase {
			return nil
		}
	}
	return protocol.ErrInvalidSystem
}

func (votesService *VerifyVotesService) IsStart() bool { return votesService.start }
