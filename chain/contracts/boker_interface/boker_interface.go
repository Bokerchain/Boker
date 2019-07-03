package boker_contract

//go: 生成abi和bin文件 solc BokerAssignToken.sol BokerAssignTokenData.sol BokerAssignTokenDefine.sol BokerAssignTokenEventHandler.sol BokerAssignTokenImpl.sol ../BokerCommon.sol
//go: 生成go文件 abigen --abi BokerAssignToken.sol:BokerAssignToken.abi --bin BokerAssignToken.sol:BokerAssignToken.bin  --pkg assigntoken --out contract.go

import (
	"context"
	"math/big"
	"time"

	"github.com/Bokerchain/Boker/chain/accounts/abi/bind"
	"github.com/Bokerchain/Boker/chain/boker/protocol"
	"github.com/Bokerchain/Boker/chain/common"
	_ "github.com/Bokerchain/Boker/chain/common/hexutil"
	"github.com/Bokerchain/Boker/chain/eth"
	"github.com/Bokerchain/Boker/chain/log"
)

//播客链合约接口
type BokerInterfaceService struct {
	bokerInterface *BokerInterface //播客链合约接口
	currentEpoch   *big.Int        //当前周期序号
	addr           common.Address  //合约地址
	ethereum       *eth.Ethereum   //以太坊对象
	tickQuit       chan chan error //tick退出chan
	epochQuit      chan chan error //epoch退出chan
	quit           chan chan error //退出chan
	start          bool            //是否已经启动
	startTimer     int64           //启动时间
	useAssignTimer int64           //
}

//创建接口服务
func NewBokerInterfaceService(ethereum *eth.Ethereum, address common.Address) (*BokerInterfaceService, error) {

	var s *BokerInterfaceService = new(BokerInterfaceService)

	bokerInterface, err := NewBokerInterface(address, eth.NewContractBackend(ethereum.ApiBackend))
	if err != nil {
		return nil, err
	}
	s.bokerInterface = bokerInterface
	s.addr = address
	s.ethereum = ethereum
	s.tickQuit = make(chan chan error)
	s.epochQuit = make(chan chan error)
	s.quit = make(chan chan error)
	s.start = false
	s.currentEpoch = big.NewInt(0)
	s.startTimer = time.Now().Unix()
	s.useAssignTimer = time.Now().Unix()

	return s, nil
}

func (s *BokerInterfaceService) createTransactOpts() *bind.TransactOpts {

	if coinbase, err := s.ethereum.Coinbase(); err == nil {

		opts := bind.NewPasswordTransactor(s.ethereum, coinbase)
		return opts
	}
	return nil
}

func (s *BokerInterfaceService) Start() {

	s.start = true
	go s.tick()
	go s.getEpoch()
	go s.assignToken()
}

func (s *BokerInterfaceService) Stop() error {

	s.start = false

	errTick := make(chan error)
	s.tickQuit <- errTick

	errEpoch := make(chan error)
	s.epochQuit <- errEpoch

	errc := make(chan error)
	s.quit <- errc

	return <-errc
}

func (s *BokerInterfaceService) tick() {

	timer := time.NewTimer(protocol.BokerInterval * 1)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			s.business()
			timer.Reset(protocol.BokerInterval * 1)
		case errc := <-s.quit:
			errc <- nil
			return
		}
	}
}

func (s *BokerInterfaceService) tickVotes() {

	//调用时钟函数，判断是否周期发生改变
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	callOpts := &bind.CallOpts{Context: ctx}
	defer cancel()

	log.Info("(s *BokerInterfaceService) tickVotes")
	epochBool, err := s.bokerInterface.BokerInterfaceCaller.TickVote(callOpts)
	if err != nil {

		if err == bind.ErrNoCode {
			log.Error("tickVote method not found", "Contract", s.addr)
		} else {
			log.Error("Failed to tickVote", "err", err)
		}
		return
	} else {

		//调用转换票数函数
		if epochBool {

			log.Info("(s *BokerInterfaceService) tickVotes")

			opts := s.createTransactOpts()
			/*if opts.Nonce == nil {
				nonce := s.ethereum.TxPool().State().GetNonce(opts.From)
				opts.Nonce = new(big.Int).SetUint64(nonce)
			}*/
			log.Info("Create TickVote TransactOpts", "From", opts.From.String(), "Nonce", opts.Nonce)

			_, err := s.bokerInterface.BokerInterfaceTransactor.RotateVote(opts)
			if err != nil {

				if err == bind.ErrNoCode {
					log.Error("rotateVote method not found", "Contract", s.addr)
				} else {
					log.Error("Failed to rotateVote", "err", err)
				}
				return
			}
		}
	}
}

//定期获取周期
func (s *BokerInterfaceService) getEpoch() {

	timer := time.NewTimer(protocol.BokerInterval * 5)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			s.vote()
			timer.Reset(protocol.BokerInterval * 5)

		case errc := <-s.epochQuit:
			errc <- nil
			return
		}
	}
}

func (s *BokerInterfaceService) vote() {

	//log.Info("Check Bokerchain Vote")

	//判断投票;
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	callOpts := &bind.CallOpts{Context: ctx}
	defer cancel()

	//获取周期
	log.Info("(s *BokerInterfaceService) vote GetVoteRound")
	epochIndex, err := s.bokerInterface.BokerInterfaceCaller.GetVoteRound(callOpts)
	if err != nil {
		if err == bind.ErrNoCode {
			log.Error("GetVoteRound method not found", "Contract", s.addr)
		} else {
			log.Error("Failed to GetVoteRound", "err", err)
		}
		return
	} else {

		log.Info("GetVoteRound method Reuslt", "epochIndex", epochIndex)

		//判断轮数是否和当前记录的是否一致，如果不一致，则重新获取数据
		if epochIndex != s.currentEpoch && epochIndex.Int64() == s.currentEpoch.Int64()+1 {

			log.Info("(s *BokerInterfaceService) vote epochIndex != s.currentEpoch && epochIndex.Int64() == s.currentEpoch.Int64()+1")

			//调用获取候选人列表数据
			candidateArray, err := s.bokerInterface.BokerInterfaceCaller.GetCandidates(callOpts)
			if err != nil {
				if err == bind.ErrNoCode {
					log.Error("GetCandidates method not found", "Contract", s.addr)
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
				if s.ethereum != nil {
					s.ethereum.BlockChain().CurrentBlock().DposCtx().SetValidatorVotes(candidateArray.Addresses, candidateArray.Tickets)
				}
			}

			s.currentEpoch.Add(s.currentEpoch, big.NewInt(1))
		}
	}

}

func (s *BokerInterfaceService) assignToken() {

	var lastTxTime int64 = 0

	for {

		time.Sleep(time.Duration(500) * protocol.AssignTickInterval)

		if lastTxTime != 0 && lastTxTime == time.Now().Unix() {
			continue
		}

		//得到第一个区块
		blocks := s.ethereum.BlockChain().GetBlockByNumber(0)
		if blocks == nil {
			continue
		}
		//得到第一个区块的时间
		firstTimer := blocks.Time().Int64()
		now := time.Now().Unix()
		offset := now - firstTimer

		if offset%protocol.TokenNoderInterval == 0 {

			log.Info("Bokerchain Assign Token Start", "Now", now, "firstTimer", firstTimer)

			//通证分配节点是否是当前节点
			if err := s.getCurrentTokenNoder(); err != nil {
				log.Error("Failed to Assign Token", "err", err)
				return
			}
			log.Info("Bokerchain Assign Token Noder Check Success")

			opts := s.createTransactOpts()
			tx, err := s.bokerInterface.BokerInterfaceTransactor.AssignToken(opts, now)
			if err != nil {
				if err == bind.ErrNoCode {
					log.Info("Bokerchain Assign Token Address Not Found", "Contract", s.addr)
				} else {
					log.Error("Bokerchain Assign Token Failed", "err", err)
				}
				return
			} else {

				if tx != nil {

					lastTxTime = tx.Time().Int64()
				}

				log.Info("Bokerchain Assign Token End")
			}

		}
	}

}

//通证分配函数
func (s *BokerInterfaceService) business() {

	//判断出块节点是否是当前节点
	if err := s.getCurrentProducer(); err != nil {
		return
	}

	//执行投票
	log.Info("(s *BokerInterfaceService) business")
	s.tickVotes()
}

func (s *BokerInterfaceService) GetTokenAddr() common.Address {

	return s.addr
}

func (s *BokerInterfaceService) getCurrentProducer() error {

	if s.ethereum != nil {

		//得到第一个区块
		blocks := s.ethereum.BlockChain().GetBlockByNumber(0)
		if blocks == nil {
			return protocol.ErrInvalidSystem
		}
		//得到第一个区块的时间
		firstTimer := blocks.Time().Int64()

		//得到当前的出块节点
		producer, err := s.ethereum.BlockChain().CurrentBlock().DposCtx().GetCurrentProducer(firstTimer)
		if err != nil {
			return protocol.ErrInvalidProducer
		}

		//得到当前挖矿节点
		coinbase, err := s.ethereum.Coinbase()
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

func (s *BokerInterfaceService) getCurrentTokenNoder() error {

	if s.ethereum != nil {

		//得到第一个区块
		blocks := s.ethereum.BlockChain().GetBlockByNumber(0)
		if blocks == nil {
			return protocol.ErrInvalidSystem
		}
		//得到第一个区块的时间
		firstTimer := blocks.Time().Int64()

		//得到当前的出块节点
		tokenNoder, err := s.ethereum.BlockChain().CurrentBlock().DposCtx().GetCurrentTokenNoder(firstTimer)
		if err != nil {
			return protocol.ErrInvalidTokenNoder
		}

		//得到当前挖矿节点
		coinbase, err := s.ethereum.Coinbase()
		if err != nil {
			return protocol.ErrInvalidCoinbase
		}

		//将当前出块节点和当前节点进行比较，如果是当前出块节点，则允许继续进行处理
		if tokenNoder == coinbase {
			return nil
		}
	}
	return protocol.ErrInvalidSystem
}

func (s *BokerInterfaceService) getNowTokenNoder(now int64) error {

	if s.ethereum != nil {

		//得到第一个区块
		blocks := s.ethereum.BlockChain().GetBlockByNumber(0)
		if blocks == nil {
			return protocol.ErrInvalidSystem
		}
		//得到第一个区块的时间
		firstTimer := blocks.Time().Int64()

		//得到当前的出块节点
		tokenNoder, err := s.ethereum.BlockChain().CurrentBlock().DposCtx().GetNowTokenNoder(firstTimer, now)
		if err != nil {
			return protocol.ErrInvalidTokenNoder
		}

		//得到当前挖矿节点
		coinbase, err := s.ethereum.Coinbase()
		if err != nil {
			return protocol.ErrInvalidCoinbase
		}

		//将当前出块节点和当前节点进行比较，如果是当前出块节点，则允许继续进行处理
		if tokenNoder == coinbase {
			return nil
		}
	}
	return protocol.ErrInvalidSystem
}

func (s *BokerInterfaceService) IsStart() bool { return s.start }
