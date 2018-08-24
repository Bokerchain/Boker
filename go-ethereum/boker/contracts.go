//播客链增加的特殊账号管理类
package boker

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/contracts/assigntoken"
	"github.com/boker/go-ethereum/contracts/votes"
	"github.com/boker/go-ethereum/core/types"
	"github.com/boker/go-ethereum/eth"
	"github.com/boker/go-ethereum/ethdb"
	"github.com/boker/go-ethereum/include"
	"github.com/boker/go-ethereum/log"
	"github.com/boker/go-ethereum/rlp"
	"github.com/boker/go-ethereum/trie"
)

var (
	basePrefix     = []byte("base-")
	contractPrefix = []byte("contracts")
)

//基本合约配置信息
type ContractConfig struct {
	ContractType    types.ContractType `json:"contractType"`    //基础合约类型
	ContractAddress common.Address     `json:"contractAddress"` //合约地址
}

type ContractConfigs struct {
	baseContractArray []ContractConfig //基础合约
}

//合约Service类
type ContractService struct {
	votesContract *votes.VerifyVotesService       //投票合约
	tokenContract *assigntoken.AssignTokenService //通证分配合约
}

//播客链的基础合约管理
type BokerContracts struct {
	baseTrie     *trie.Trie                            //基础合约Hash树，Key为基础合约的Address Value为合约的类型
	contractTrie *trie.Trie                            //所有基础合约整体保存树，这里保存所有的基础合约信息，但是不保存基础合约的类型信息
	ethereum     *eth.Ethereum                         //以太坊对象
	accounts     *BokerAccount                         //播客链账号对象
	transactions *BokerTransaction                     //播客链交易对象
	db           ethdb.Database                        //数据库
	contracts    map[common.Address]types.ContractType //基础合约的Map
	services     ContractService                       //合约服务类
}

func NewContract(db ethdb.Database, ethereum *eth.Ethereum, transactions *BokerTransaction) (*BokerContracts, error) {

	log.Info("new base contract baseTrie")
	baseTrie, errTrie := NewBaseTrie(common.Hash{}, db)
	if errTrie != nil {
		return nil, errTrie
	}
	log.Info("new base contract contractTrie")
	contractTrie, errTrie := NewContractTrie(common.Hash{}, db)
	if errTrie != nil {
		return nil, errTrie
	}

	//创建对象
	base := new(BokerContracts)
	base.baseTrie = baseTrie
	base.contractTrie = contractTrie
	base.contracts = make(map[common.Address]types.ContractType)
	base.db = db
	base.ethereum = ethereum
	base.transactions = transactions

	//加载合约
	base.loadContract()
	base.readContractConfig()

	//创建投票合约服务
	var address common.Address
	var err error
	address, err = base.getContractAddress(types.ContractVote)
	if err != nil {
		return nil, err
	}
	base.services.votesContract, err = votes.NewVerifyVotesService(base.ethereum, include.BokerConfig{address})
	if err != nil {
		return nil, ErrNewContractService
	}
	base.services.votesContract.Start()

	//创建分配通证合约服务
	address, err = base.getContractAddress(types.ContractAssignToken)
	if err != nil {
		return nil, err
	}
	base.services.tokenContract, err = assigntoken.NewAssignTokenService(base.ethereum, include.BokerConfig{address})
	if err != nil {
		return nil, ErrNewContractService
	}
	base.services.tokenContract.Start()

	return base, nil
}

//创建基础合约的Hash树
func NewBaseTrie(root common.Hash, db ethdb.Database) (*trie.Trie, error) {
	return trie.NewTrieWithPrefix(root, basePrefix, db)
}

func NewContractTrie(root common.Hash, db ethdb.Database) (*trie.Trie, error) {
	return trie.NewTrieWithPrefix(root, contractPrefix, db)
}

//读取配置文件中关于合约配置的信息
func (c *BokerContracts) readContractConfig() (ContractConfigs, error) {

	//首先判断当前读取的用户是否是特殊账号用户
	address, err := c.ethereum.Coinbase()
	if err != nil {
		return ContractConfigs{}, err
	}
	txType, err := c.accounts.GetAccount(address)
	if err != nil {
		return ContractConfigs{}, err
	}
	if txType != types.DeployVote && txType != types.DeployAssignToken {
		return ContractConfigs{}, ErrLevel
	}

	//读取配置文件信息
	buffer, err := ioutil.ReadFile(JsonFileName)
	if err != nil {
		return ContractConfigs{}, err
	}
	config := ContractConfigs{}
	err = json.Unmarshal(buffer, &config)
	if err != nil {
		return ContractConfigs{}, err
	}

	//对写入的配置合约进行处理
	for _, v := range config.baseContractArray {

		contractType, err := c.getContractType(v.ContractAddress)
		if err == ErrNotFoundContract {

			//设置一个设置交易
			if contractType == types.ContractVote {
				c.transactions.DeployTransaction(types.DeployVote, v.ContractAddress)
			} else if contractType == types.ContractAssignToken {
				c.transactions.DeployTransaction(types.DeployAssignToken, v.ContractAddress)
			}
			continue
		}

		if contractType != v.ContractType {

			log.Error("Contract Type Atypism ", "config", v.ContractType, "trie", contractType)
		}
	}

	return config, nil
}

//设置合约到Hash树中
func (c *BokerContracts) SetContract(address common.Address, contractType types.ContractType) error {

	//设置基础合约
	log.Info("SetContract find base contract exist")
	exist, err := c.existContract(address)
	if err != nil {
		return err
	}
	if exist {
		return ErrContractExist
	}

	log.Info("SetContract contract is`t exist TryUpdate")
	if err = c.baseTrie.TryUpdate(address.Bytes(), []byte(strconv.Itoa(int(contractType)))); err != nil {
		return err
	}
	c.contracts[address] = contractType

	//更新基础合约列表保存树
	if err = c.setContractsTrie(); err != nil {
		return ErrSaveContractTrie
	}

	//判断是否需要启动合约
	if contractType == types.ContractVote {

		c.services.votesContract, err = votes.NewVerifyVotesService(c.ethereum, include.BokerConfig{address})
		if err != nil {
			return ErrNewContractService
		}
		c.services.votesContract.Start()
	} else if contractType == types.ContractAssignToken {

		c.services.tokenContract, err = assigntoken.NewAssignTokenService(c.ethereum, include.BokerConfig{address})
		if err != nil {
			return ErrNewContractService
		}
		c.services.tokenContract.Start()
	}
	return nil
}

//设置合约到Hash树中
func (c *BokerContracts) CancelContract(address common.Address) error {

	//检测合约是否存在
	log.Info("CancelContract find base contract exist")
	contractType, err := c.GetContract(address)
	if err != nil {
		return err
	}

	//从Hash树中删除
	log.Info("CancelContract contract is`t exist TryUpdate")
	if err := c.baseTrie.TryDelete(address.Bytes()); err != nil {
		return err
	}
	delete(c.contracts, address)

	//更新基础合约列表保存树
	if err = c.setContractsTrie(); err != nil {
		return ErrSaveContractTrie
	}

	//终止合约运行
	if contractType == types.ContractVote {

		if (c.services.votesContract != nil) && c.services.votesContract.IsStart() {

		}
	} else if contractType == types.ContractAssignToken {

		if (c.services.tokenContract != nil) && c.services.tokenContract.IsStart() {

		}
	}
	return nil
}

//将所有的合约地址进行设置
func (c *BokerContracts) setContractsTrie() error {

	//转换成切片数组
	log.Info("rotate contract map to array")
	var contracts []common.Address = make([]common.Address, 0)
	for k, _ := range c.contracts {
		contracts = append(contracts, k)
	}

	log.Info("setContracts contract context encode")
	key := []byte("contracts")
	contractsRLP, err := rlp.EncodeToBytes(contracts)
	if err != nil {

		log.Error("failed to encode contracts to rlp", "error", err)
		return err
	}
	log.Info("setContracts contract hash trie Update")
	c.contractTrie.Update(key, contractsRLP)
	return nil
}

//得到所有的合约地址
func (c *BokerContracts) getContractsTrie() ([]common.Address, error) {

	log.Info("getContracts get contracts")
	var contracts []common.Address
	key := []byte("contracts")
	contractsRLP := c.contractTrie.Get(key)
	if err := rlp.DecodeBytes(contractsRLP, &contracts); err != nil {

		log.Error("failed to decode contracts", "error", err)
		return nil, err
	}
	return contracts, nil
}

//从Hash树中获取合约类型
func (c *BokerContracts) readContractType(address common.Address) (types.ContractType, error) {

	//根据地址获取数据
	key := address.Bytes()
	v, err := c.baseTrie.TryGet(key)
	if err != nil {
		return types.ContractBinary, err
	}

	//转换成交易类型
	contractType, err := strconv.Atoi(string(v[:]))
	return types.ContractType(contractType), nil
}

//根据合约类型得到合约地址
func (c *BokerContracts) getContractAddress(contractType types.ContractType) (common.Address, error) {

	if len(c.contracts) <= 0 {
		return common.Address{}, ErrNotFoundContract
	}

	for k, v := range c.contracts {

		if v == contractType {
			return k, nil
		}
	}
	return common.Address{}, ErrNotFoundContract
}

//根据合约地址得到合约类型
func (c *BokerContracts) getContractType(address common.Address) (types.ContractType, error) {

	contractType, exist := c.contracts[address]
	if exist {
		return contractType, nil
	} else {
		return types.ContractBinary, ErrNotFoundContract
	}
}

//加载基础合约信息
func (c *BokerContracts) loadContract() error {

	//获取所有合约
	contracts, err := c.getContractsTrie()
	if err != nil {
		return err
	}

	//根据所有合约得到合约类型
	for _, address := range contracts {
		contractType, err := c.readContractType(address)
		if err != nil {
			continue
		}
		c.contracts[address] = contractType
	}
	return nil
}

//查找合约账户
func (c *BokerContracts) GetContract(address common.Address) (types.ContractType, error) {

	if len(c.contracts) > 0 {
		value, exist := c.contracts[address]
		if exist {
			return value, nil
		}
	}
	return types.ContractBinary, ErrNotFoundContract
}

//判断此合约是否已经存在
func (c *BokerContracts) existContract(address common.Address) (bool, error) {

	if len(c.contracts) <= 0 {
		return false, ErrNotFoundContract
	}
	_, exist := c.contracts[address]
	if exist {
		return true, nil
	} else {
		return false, ErrNotFoundContract
	}
}
