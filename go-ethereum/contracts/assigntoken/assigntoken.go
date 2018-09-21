package assigntoken

//go: 生成abi和bin文件 solc BokerAssignToken.sol BokerAssignTokenData.sol BokerAssignTokenDefine.sol BokerAssignTokenEventHandler.sol BokerAssignTokenImpl.sol ../BokerCommon.sol
//go: 生成go文件 abigen --abi BokerAssignToken.sol:BokerAssignToken.abi --bin BokerAssignToken.sol:BokerAssignToken.bin  --pkg assigntoken --out contract.go

import (
	"context"
	"time"

	"github.com/boker/go-ethereum/accounts/abi/bind"
	"github.com/boker/go-ethereum/boker/protocol"
	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/eth"
	"github.com/boker/go-ethereum/log"
)

//定期进行分配通证
type AssignTokenService struct {
	token    *Assigntoken    //分币session
	addr     common.Address  //合约地址
	ethereum *eth.Ethereum   //以太坊对象
	quit     chan chan error //退出chan
	start    bool            //是否已经启动
}

//创建一个新服务来定期执行
func NewAssignTokenService(ethereum *eth.Ethereum, address common.Address) (*AssignTokenService, error) {

	var assignToken *AssignTokenService = new(AssignTokenService)

	token, err := NewAssigntoken(address, eth.NewContractBackend(ethereum.ApiBackend))
	if err != nil {
		return nil, err
	}
	assignToken.token = token
	assignToken.addr = address
	assignToken.ethereum = ethereum
	assignToken.quit = make(chan chan error)
	assignToken.start = false

	return assignToken, nil
}

func (tokenService *AssignTokenService) createTransactOpts() *bind.TransactOpts {

	if coinbase, err := tokenService.ethereum.Coinbase(); err == nil {
		return bind.NewPasswordTransactor(tokenService.ethereum, coinbase)
	}
	return nil
}

func (tokenService *AssignTokenService) Start() {

	//启动tick函数
	tokenService.start = true
	go tokenService.tick()
}

func (tokenService *AssignTokenService) Stop() error {

	tokenService.start = false
	errc := make(chan error)
	tokenService.quit <- errc
	return <-errc
}

func (tokenService *AssignTokenService) tick() {

	timer := time.NewTimer(protocol.AssignTokenInterval * 1)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			tokenService.assignToken()
			timer.Reset(protocol.AssignTokenInterval * 1)
		case errc := <-tokenService.quit:
			errc <- nil
			return
		}
	}
}

//通证分配函数
func (tokenService *AssignTokenService) assignToken() {

	//判断出块节点是否是当前节点
	if err := tokenService.getCurrentTokenNoder(); err != nil {
		log.Error("Failed to Assign Token", "err", err)
		return
	}

	//调用时钟函数，判断是否周期发生改变
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	callOpts := &bind.CallOpts{Context: ctx}
	defer cancel()

	tokenBool, err := tokenService.token.CheckAssignToken(callOpts)
	if err != nil {
		if err == bind.ErrNoCode {
			log.Debug("Assign Token address not found", "Contract", tokenService.addr)
		} else {
			log.Error("Failed to retrieve current release", "err", err)
		}
		return
	} else {

		//开始分配通证
		if tokenBool {

			opts := tokenService.createTransactOpts()
			_, err := tokenService.token.AssignToken(opts)
			if err != nil {
				if err == bind.ErrNoCode {
					log.Debug("Release oracle not found", "Contract", tokenService.addr)
				} else {
					log.Error("Failed to retrieve current release", "err", err)
				}
				return
			}
		}
	}
}

func (tokenService *AssignTokenService) GetTokenAddr() common.Address {

	return tokenService.addr
}

func (tokenService *AssignTokenService) getCurrentTokenNoder() error {

	if tokenService.ethereum != nil {

		//得到当前的出块节点
		producer, err := tokenService.ethereum.BlockChain().CurrentBlock().DposCtx().GetCurrentProducer()
		if err != nil {
			return protocol.ErrInvalidTokenNoder
		}

		//得到当前挖矿节点
		coinbase, err := tokenService.ethereum.Coinbase()
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

func (tokenService *AssignTokenService) IsStart() bool { return tokenService.start }
