package test

//go: 生成abi和bin文件 solc BokerAssignToken.sol BokerAssignTokenData.sol BokerAssignTokenDefine.sol BokerAssignTokenEventHandler.sol BokerAssignTokenImpl.sol ../BokerCommon.sol
//go: 生成go文件 abigen --abi BokerAssignToken.sol:BokerAssignToken.abi --bin BokerAssignToken.sol:BokerAssignToken.bin  --pkg assigntoken --out contract.go

import (
	"context"
	"time"

	"github.com/Bokerchain/Boker/chain/accounts/abi/bind"
	"github.com/Bokerchain/Boker/chain/boker/protocol"
	"github.com/Bokerchain/Boker/chain/common"
	"github.com/Bokerchain/Boker/chain/eth"
	"github.com/Bokerchain/Boker/chain/log"
)

//定期进行分配通证
type TestService struct {
	testPtr  *Test           //分币session
	addr     common.Address  //合约地址
	ethereum *eth.Ethereum   //以太坊对象
	quit     chan chan error //退出chan
	start    bool            //是否已经启动
}

//创建一个新服务来定期执行
func NewTestService(ethereum *eth.Ethereum, address common.Address) (*TestService, error) {

	var testService *TestService = new(TestService)
	testPtr, err := NewTest(address, eth.NewContractBackend(ethereum.ApiBackend))
	if err != nil {
		return nil, err
	}
	testService.testPtr = testPtr
	testService.addr = address
	testService.ethereum = ethereum
	testService.quit = make(chan chan error)
	testService.start = false

	return testService, nil
}

func (testService *TestService) Start() {

	//启动tick函数
	testService.start = true
	go testService.tick()
}

func (testService *TestService) Stop() error {

	testService.start = false
	errc := make(chan error)
	testService.quit <- errc
	return <-errc
}

func (testService *TestService) tick() {

	timer := time.NewTimer(protocol.BokerInterval * 1)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			testService.assignToken()
			timer.Reset(protocol.BokerInterval * 1)
		case errc := <-testService.quit:
			errc <- nil
			return
		}
	}
}

//通证分配函数
func (testService *TestService) assignToken() {

	//调用时钟函数，判断是否周期发生改变
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	callOpts := &bind.CallOpts{Context: ctx}
	defer cancel()

	testInt, err := testService.testPtr.Test(callOpts)
	if err != nil {
		if err == bind.ErrNoCode {
			log.Info("Assign Token address not found", "Contract", testService.addr)
		} else {
			log.Error("ChechAssignToken Failed", "err", err)
		}
		return
	}
	//开始分配通证
	log.Info("Check AssignToken Success", "tokenInt", testInt)
}
