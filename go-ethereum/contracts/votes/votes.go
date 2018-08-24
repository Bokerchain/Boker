package votes

//go: 生成abi和bin文件 solc BokerVerifyVote.sol BokerVerifyVoteImpl.sol BokerVerifyVoteData.sol BokerCommon.sol
//go: 生成go文件 abigen --abi BokerVerifyVote.sol:BokerVerifyVote.abi --bin BokerVerifyVote.sol:BokerVerifyVote.bin  --pkg votes --out contract.go

import (
	"math/big"
	"time"

	"github.com/boker/go-ethereum/accounts/abi/bind"
	"github.com/boker/go-ethereum/accounts/abi/bind/backends"
	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/core"
	"github.com/boker/go-ethereum/crypto"
	"github.com/boker/go-ethereum/eth"
	"github.com/boker/go-ethereum/include"
	"github.com/boker/go-ethereum/log"
)

//投票服务
type VerifyVotesService struct {
	config       include.BokerConfig  //通证分配的配置
	votesSession VotesSession         //分币session
	addr         common.Address       //合约地址
	backend      bind.ContractBackend //后台对象
	ethereum     *eth.Ethereum        //以太坊对象
	currentEpoch *big.Int             //当前周期序号
	tickQuit     chan chan error      //tick退出chan
	epochQuit    chan chan error      //epoch退出chan
	start        bool                 //是否已经启动
}

//创建一个新服务来定期执行
func NewVerifyVotesService(ethereum *eth.Ethereum, config include.BokerConfig) (*VerifyVotesService, error) {

	//创建投票对象
	var votesService *VerifyVotesService = new(VerifyVotesService)
	transactOpts, backend, addr, contract, err := votesService.initContract()
	if err != nil {
		return nil, err
	}

	//定义一个分配币的session
	session := VotesSession{
		Contract:     contract,
		TransactOpts: *transactOpts,
	}
	votesService.votesSession = session
	votesService.addr = addr
	votesService.backend = backend
	votesService.config = config
	votesService.tickQuit = make(chan chan error)
	votesService.epochQuit = make(chan chan error)
	votesService.start = false

	return votesService, nil
}

func (votesService *VerifyVotesService) initContract() (*bind.TransactOpts, bind.ContractBackend, common.Address, *Votes, error) {

	//根据读取到的数据来进行处理
	DeployKey, err := crypto.HexToECDSA(votesService.ethereum.BlockChain().Config().Producer.PrivateKey)
	if err != nil {

	}
	DeployAddr := crypto.PubkeyToAddress(DeployKey.PublicKey)
	DeployBalance := big.NewInt(0)
	DeployBalance.SetInt64(votesService.ethereum.BlockChain().Config().Producer.Balance)

	//构造backend和帐号
	backend := backends.NewSimulatedBackend(core.GenesisAlloc{DeployAddr: {Balance: DeployBalance}}, votesService.ethereum.Boker)
	auth := bind.NewKeyedTransactor(DeployKey)

	//部署合约并得到合约地址
	addr, _, contract, err := DeployVotes(auth, backend)
	if err != nil {
		panic(err)
	}

	//提交合约
	backend.Commit()
	code, err := backend.CodeAt(nil, addr, nil)
	if err != nil {
		panic(err)
	}
	if len(code) == 0 {
		panic("empty code")
	}

	return auth, backend, addr, contract, nil
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

//产生tick时钟
func (votesService *VerifyVotesService) tick() {

	timer := time.NewTimer(include.VotesInterval * 1)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			votesService.tickVotes()
			timer.Reset(include.VotesInterval * 1)
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
	epochBool, err := votesService.votesSession.TickVote()
	if err != nil {
		if err == bind.ErrNoCode {
			log.Error("tickVote method not found", "Contract", votesService.config.Address)
		} else {
			log.Error("Failed to tickVote", "err", err)
		}
		return
	} else {

		//调用转换票数函数
		if epochBool {
			_, err := votesService.votesSession.RotateVote()
			if err != nil {
				if err == bind.ErrNoCode {
					log.Error("rotateVote method not found", "Contract", votesService.config.Address)
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

	timer := time.NewTimer(include.VotesInterval * 5)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			votesService.tickEpoch()
			timer.Reset(include.VotesInterval * 5)

		case errc := <-votesService.epochQuit:
			errc <- nil
			return
		}
	}
}

func (votesService *VerifyVotesService) tickEpoch() {

	//调用是否周期发生改变
	epochIndex, err := votesService.votesSession.GetVoteRound()
	if err != nil {
		if err == bind.ErrNoCode {
			log.Error("GetVoteRound method not found", "Contract", votesService.config.Address)
		} else {
			log.Error("Failed to GetVoteRound", "err", err)
		}
		return
	} else {

		//判断轮数是否和当前记录的是否一致，如果不一致，则重新获取数据
		if epochIndex != votesService.currentEpoch && epochIndex.Int64() == votesService.currentEpoch.Int64()+1 {

			//调用获取候选人列表数据
			candidateArray, err := votesService.votesSession.GetCandidates()
			if err != nil {
				if err == bind.ErrNoCode {
					log.Error("GetCandidates method not found", "Contract", votesService.config.Address)
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

	return votesService.config.Address
}

func (votesService *VerifyVotesService) getCurrentProducer() error {

	if votesService.ethereum != nil {

		//得到当前的出块节点
		producer, err := votesService.ethereum.BlockChain().CurrentBlock().DposCtx().GetCurrentProducer()
		if err != nil {
			return include.ErrInvalidProducer
		}

		//得到当前挖矿节点
		coinbase, err := votesService.ethereum.Coinbase()
		if err != nil {
			return include.ErrInvalidCoinbase
		}

		//将当前出块节点和当前节点进行比较，如果是当前出块节点，则允许继续进行处理
		if producer == coinbase {
			return nil
		}
	}
	return include.ErrInvalidSystem
}

func (votesService *VerifyVotesService) IsStart() bool { return votesService.start }
