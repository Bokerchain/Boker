//播客链增加的特殊账号管理类
package boker

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/boker/go-ethereum/boker/protocol"
	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/eth"
	"github.com/boker/go-ethereum/log"
	"github.com/boker/go-ethereum/params"
	"github.com/boker/go-ethereum/trie"
)

const JsonFileName = "boker.json" //播客链配置

//播客链中基础合约相关配置信息
type BaseContract struct {
	ContractType    uint64         `json:"contracttype"`    //基础合约类型
	DeployAddress   common.Address `json:"deployaddress"`   //部署账号
	ContractAddress common.Address `json:"contractaddress"` //合约地址
}

type BaseContractConfig struct {
	Bases []BaseContract `json:"bases,omitempty"`
}

//播客链用来分配通证的账号和私钥信息
type ProducerConfig struct {
	Coinbase common.Address `json:"coinbase"` //挖矿的Coinbase
	Password string         `json:"password"` //Coinbase的密码
}

type BokerConfig struct {
	Dpos      *params.DposConfig  `json:"dpos,omitempty"`      //Dpos的配置信息
	Contracts *BaseContractConfig `json:"contracts,omitempty"` //基础合约配置信息
	Producer  *ProducerConfig     `json:"producer,omitempty"`  //出块节点使用的信息
}

type BokerBackend struct {
	config       BokerConfig
	ethereum     *eth.Ethereum
	accounts     *BokerAccount
	contracts    *BokerContracts
	transactions *BokerTransaction
}

func New() *BokerBackend {

	log.Info("****New Boker****")
	boker := new(BokerBackend)
	boker.ethereum = nil
	boker.accounts = nil
	boker.contracts = nil
	boker.transactions = nil

	boker.loadConfig()
	return boker
}

func (boker *BokerBackend) Init(e *eth.Ethereum, bokerProto *protocol.BokerBackendProto) error {

	log.Info("****Init Boker****")

	//创建类
	boker.ethereum = e

	log.Info("Create Transaction Object")
	boker.transactions = NewTransaction(e)

	log.Info("Create Account Object")
	boker.accounts = NewAccount()

	var err error
	log.Info("Create Contract Object")
	boker.contracts, err = NewContract(e.ChainDb(), e, boker.transactions, bokerProto)
	if err != nil {
		return nil
	}
	return nil
}

//loadConfig 加载json格式的配置文件
func (boker *BokerBackend) loadConfig() error {

	log.Info("****loadConfig****")

	if boker.ethereum != nil {

		boker.config.Dpos = new(params.DposConfig)
		boker.config.Dpos.Validators = make([]common.Address, 0)

		boker.config.Contracts = new(BaseContractConfig)
		boker.config.Contracts.Bases = make([]BaseContract, 0)

		boker.config.Producer = new(ProducerConfig)

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
		boker.config.Dpos.Validators = make([]common.Address, 0)
		for _, v := range config.Dpos.Validators {
			boker.config.Dpos.Validators = append(boker.config.Dpos.Validators, v)
		}
		log.Info("Load BokerChain Dpos Validators", "Size", len(boker.config.Dpos.Validators))

		boker.config.Contracts.Bases = make([]BaseContract, 0)
		for _, v := range config.Contracts.Bases {
			boker.config.Contracts.Bases = append(boker.config.Contracts.Bases, v)
		}
		log.Info("Load BokerChain Base Contracts", "Size", len(boker.config.Contracts.Bases))

		boker.config.Producer.Coinbase = config.Producer.Coinbase
		boker.config.Producer.Password = config.Producer.Password
		log.Info("Load BokerChain Producer", "Coinbase", boker.config.Producer.Coinbase, "Password", boker.config.Producer.Password)

		return nil
	}
	return protocol.ErrSystem
}

//GetAccount 根据账号地址，得到账号等级
func (boker *BokerBackend) GetAccount(account common.Address) ([]protocol.TxType, error) {

	log.Info("****GetAccount****", "account", account.String())
	if boker.accounts == nil {

		log.Error("GetAccount function Boker accounts Objects is nil")
		return []protocol.TxType{protocol.Binary}, nil
	}

	return boker.accounts.GetAccount(account)
}

//GetContract 根据合约地址，得到合约等级
func (boker *BokerBackend) GetContract(address common.Address) (protocol.ContractType, error) {
	return boker.contracts.GetContract(address)
}

//SetContract 回写合约信息
func (boker *BokerBackend) SetContract(address common.Address, contractType protocol.ContractType) error {
	return boker.contracts.SetContract(address, contractType)
}

func (boker *BokerBackend) IsValidator(address common.Address) bool {

	//return boker.accounts.IsValidator(address)

	//测试使用
	return true
}

//SubmitBokerTransaction 设置一个播客链交易
func (boker *BokerBackend) SubmitBokerTransaction(ctx context.Context, txType protocol.TxType, to common.Address) error {
	return boker.transactions.SubmitBokerTransaction(ctx, txType, to)
}

func (boker *BokerBackend) CommitTrie() (*protocol.BokerBackendProto, error) {

	//提交基础合约交易
	contractRoot, err := boker.contracts.contractTrie.CommitTo(boker.ethereum.ChainDb())
	if err != nil {
		return nil, err
	}

	return &protocol.BokerBackendProto{
		ContractHash: contractRoot,
	}, nil
}

func (boker *BokerBackend) GetContractTrie() (*trie.Trie, *trie.Trie) {

	return boker.contracts.GetContractTrie()
}

func (boker *BokerBackend) GetMethodName(txType protocol.TxType) (string, string, error) {

	if txType < protocol.SetVote {
		return "", "", protocol.ErrTxType
	}

	switch txType {

	case protocol.SetValidator: //设置验证者
		return "", "", nil

	case protocol.RegisterCandidate: //注册成为候选人
		return "", protocol.RegisterCandidateMethod, nil

	case protocol.VoteUser: //用户投票
		return "", protocol.VoteCandidateMethod, nil

	case protocol.VoteEpoch: //产生当前的出块节点
		return "", protocol.RotateVoteMethod, nil

	case protocol.AssignToken: //分配通证
		return "", protocol.AssignTokenMethod, nil

	case protocol.AssignReward: //出块节点的通证奖励
		return "", "", nil

	default:
		return "", "", protocol.ErrTxType
	}
}
