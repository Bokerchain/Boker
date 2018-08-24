//播客链增加的特殊账号管理类
package boker

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/core/types"
	"github.com/boker/go-ethereum/eth"
	"github.com/boker/go-ethereum/log"
	"github.com/boker/go-ethereum/params"
)

const JsonFileName = "boker.json" //播客链配置

var (
	ErrLoadConfig         = errors.New("load bokerchain config error")          //加载配置信息出错
	ErrNotFoundAddress    = errors.New("not found bokerchain contract address") //没有找到合约地址
	ErrNotFoundType       = errors.New("not found bokerchain contract type")    //没有找到合约类型
	ErrWriteJson          = errors.New("write bokerchain json file error")      //写保存基础合约的Json格式出错
	ErrOpenFile           = errors.New("open bokerchain json file error")       //打开基础合约保存文件出错
	ErrWriteFile          = errors.New("bokerchain write file error")           //写基础合约保存文件出错
	ErrContractExist      = errors.New("bokerchain contract aleady exist")      //写基础合约保存文件出错
	ErrSystem             = errors.New("system error")                          //系统错误
	ErrNotFoundContract   = errors.New("not found bokerchain contract")         //没有找到合约
	ErrNotFoundAccount    = errors.New("not found bokerchain account")          //没有找到合约
	ErrNewContractService = errors.New("create bokerchain base contract err")   //没有找到合约
	ErrSaveContractTrie   = errors.New("save contract trie err")                //没有找到合约
	ErrLevel              = errors.New("account level error")                   //没有找到合约
)

type BokerConfig struct {
	Dpos      *params.DposConfig         `json:"dpos,omitempty"`      //Dpos的配置信息
	Contracts *params.BaseContractConfig `json:"contracts,omitempty"` //基础合约配置信息
	Producer  *params.ProducerConfig     `json:"producer,omitempty"`  //出块节点使用的信息
}

type BokerBackend struct {
	ethereum     *eth.Ethereum
	accounts     *BokerAccount
	contracts    *BokerContracts
	transactions *BokerTransaction
}

func New(ether *eth.Ethereum) *BokerBackend {

	bokerBackend := new(BokerBackend)

	//创建类
	bokerBackend.ethereum = ether
	bokerBackend.transactions = NewTransaction(bokerBackend.ethereum)
	bokerBackend.accounts = NewAccount()

	var err error
	bokerBackend.contracts, err = NewContract(ether.ChainDb(), bokerBackend.ethereum, bokerBackend.transactions)
	if err != nil {
		return nil
	}

	//加载配置信息
	bokerBackend.loadConfig()
	return bokerBackend
}

//loadConfig 加载json格式的配置文件
func (boker *BokerBackend) loadConfig() error {

	if boker.ethereum != nil {

		//获取到配置信息
		chainConfig := boker.ethereum.BlockChain().Config()
		if chainConfig == nil {
			log.Error("Load Ethereum Config File Error", "Error", ErrLoadConfig)
			return ErrLoadConfig
		}

		//这里添加一次判断，判断配置信息中是否为nil，为了防止崩溃
		if chainConfig.Dpos == nil {
			chainConfig.Dpos = new(params.DposConfig)
			chainConfig.Dpos.Validators = make([]common.Address, 0)
		}
		if chainConfig.Contracts == nil {
			chainConfig.Contracts = new(params.BaseContractConfig)
			chainConfig.Contracts.Bases = make([]params.BaseContract, 0)
		}
		if chainConfig.Producer == nil {
			chainConfig.Producer = new(params.ProducerConfig)
		}

		//读取文件
		buffer, err := ioutil.ReadFile(JsonFileName)
		if err != nil {
			log.Error("Read boker.json File Error", "Error", err)
			return err
		}

		config := BokerConfig{}
		err = json.Unmarshal(buffer, &config)
		if err != nil {
			log.Error("Unmarshal boker.json File Error", "Error", err)
			return err
		}

		//清空原来的数据
		chainConfig.Dpos.Validators = make([]common.Address, 0)
		for _, v := range config.Dpos.Validators {
			chainConfig.Dpos.Validators = append(chainConfig.Dpos.Validators, v)
		}
		chainConfig.Contracts.Bases = make([]params.BaseContract, 0)
		for _, v := range config.Contracts.Bases {
			chainConfig.Contracts.Bases = append(chainConfig.Contracts.Bases, v)
		}
		chainConfig.Producer.Account = config.Producer.Account
		chainConfig.Producer.PrivateKey = config.Producer.PrivateKey

		return nil
	}
	return ErrSystem
}

//GetAccount 根据账号地址，得到账号等级
func (boker *BokerBackend) GetAccount(account common.Address) (types.TxType, error) {
	return boker.accounts.GetAccount(account)
}

//GetContract 根据合约地址，得到合约等级
func (boker *BokerBackend) GetContract(address common.Address) (types.ContractType, error) {
	return boker.contracts.GetContract(address)
}

//SetContract 回写合约信息
func (boker *BokerBackend) SetContract(address common.Address, contractType types.ContractType) error {
	return boker.contracts.SetContract(address, contractType)
}

//SetDeployTransaction 设置一个部署播客链基础合约交易
func (boker *BokerBackend) SetDeployTransaction(txType types.TxType, address common.Address) error {

	if txType == types.DeployAssignToken || txType == types.DeployVote {

		boker.transactions.DeployTransaction(txType, address)
	} else if txType == types.UnDeployAssignToken || txType == types.UnDeployVote {

		boker.transactions.UnDeployTransaction(txType, address)
	}
	return nil
}

//CancelDeployTransaction 取消一个播客链基础合约交易
func (boker *BokerBackend) CancelDeployTransaction(txType types.TxType, address common.Address) error {

	if txType == types.DeployAssignToken || txType == types.DeployVote {

		boker.transactions.DeployTransaction(txType, address)
	} else if txType == types.UnDeployAssignToken || txType == types.UnDeployVote {

		boker.transactions.UnDeployTransaction(txType, address)
	}
	return nil
}
