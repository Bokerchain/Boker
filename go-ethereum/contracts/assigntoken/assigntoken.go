package assigntoken

//go: 生成abi和bin文件 solc BokerAssignToken.sol BokerAssignTokenData.sol BokerAssignTokenDefine.sol BokerAssignTokenEventHandler.sol BokerAssignTokenImpl.sol ../BokerCommon.sol
//go: 生成go文件 abigen --abi BokerAssignToken.sol:BokerAssignToken.abi --bin BokerAssignToken.sol:BokerAssignToken.bin  --pkg assigntoken --out contract.go

import (
	_ "errors"
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

//定义部署合约的用户信息
/*var (
	DeployKey, _  = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	DeployAddr    = crypto.PubkeyToAddress(DeployKey.PublicKey)
	DeployBalance = big.NewInt(1000000)
)*/

//定期进行分配通证
type AssignTokenService struct {
	config       include.BokerConfig  //通证分配的配置
	tokenSession AssigntokenSession   //分币session
	addr         common.Address       //合约地址
	backend      bind.ContractBackend //后台对象
	ethereum     *eth.Ethereum        //以太坊对象
	quit         chan chan error      //退出chan
	start        bool                 //是否已经启动
}

//创建一个新服务来定期执行
func NewAssignTokenService(ethereum *eth.Ethereum, config include.BokerConfig) (*AssignTokenService, error) {

	var assignToken *AssignTokenService = new(AssignTokenService)
	transactOpts, backend, addr, contract, err := assignToken.initContract()
	if err != nil {
		return nil, err
	}

	//定义一个分配币的session
	session := AssigntokenSession{
		Contract:     contract,
		TransactOpts: *transactOpts,
	}
	assignToken.tokenSession = session
	assignToken.addr = addr
	assignToken.backend = backend
	assignToken.config = config
	assignToken.ethereum = ethereum
	assignToken.quit = make(chan chan error)
	assignToken.start = false

	return assignToken, nil
}

func (tokenService *AssignTokenService) initContract() (*bind.TransactOpts, bind.ContractBackend, common.Address, *Assigntoken, error) {

	//根据读取到的数据来进行处理
	DeployKey, err := crypto.HexToECDSA(tokenService.ethereum.BlockChain().Config().Producer.PrivateKey)
	if err != nil {

	}
	DeployAddr := crypto.PubkeyToAddress(DeployKey.PublicKey)
	DeployBalance := big.NewInt(0)
	DeployBalance.SetInt64(tokenService.ethereum.BlockChain().Config().Producer.Balance)

	//构造backend和帐号
	backend := backends.NewSimulatedBackend(core.GenesisAlloc{DeployAddr: {Balance: DeployBalance}}, tokenService.ethereum.Boker)
	auth := bind.NewKeyedTransactor(DeployKey)

	//部署合约并得到合约地址
	addr, _, contract, err := DeployAssigntoken(auth, backend)
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

	timer := time.NewTimer(include.AssignTokenInterval * 1)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			tokenService.assignToken()
			timer.Reset(include.AssignTokenInterval * 1)
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
	tokenBool, err := tokenService.tokenSession.CheckAssignToken()
	if err != nil {
		if err == bind.ErrNoCode {
			log.Debug("Assign Token address not found", "Contract", tokenService.config.Address)
		} else {
			log.Error("Failed to retrieve current release", "err", err)
		}
		return
	} else {

		//开始分配通证
		if tokenBool {

			_, err := tokenService.tokenSession.AssignToken()
			if err != nil {
				if err == bind.ErrNoCode {
					log.Debug("Release oracle not found", "Contract", tokenService.config.Address)
				} else {
					log.Error("Failed to retrieve current release", "err", err)
				}
				return
			}
		}
	}
}

func (tokenService *AssignTokenService) GetTokenAddr() common.Address {

	return tokenService.config.Address
}

func (tokenService *AssignTokenService) getCurrentTokenNoder() error {

	if tokenService.ethereum != nil {

		//得到当前的出块节点
		producer, err := tokenService.ethereum.BlockChain().CurrentBlock().DposCtx().GetCurrentProducer()
		if err != nil {
			return include.ErrInvalidTokenNoder
		}

		//得到当前挖矿节点
		coinbase, err := tokenService.ethereum.Coinbase()
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

func (tokenService *AssignTokenService) IsStart() bool { return tokenService.start }
