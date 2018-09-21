//播客链增加的特殊账号管理类
package boker

import (
	"strconv"

	"github.com/boker/go-ethereum/boker/protocol"
	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/contracts/assigntoken"
	"github.com/boker/go-ethereum/contracts/votes"
	"github.com/boker/go-ethereum/eth"
	"github.com/boker/go-ethereum/ethdb"
	"github.com/boker/go-ethereum/log"
	"github.com/boker/go-ethereum/rlp"
	"github.com/boker/go-ethereum/trie"
)

//基本合约配置信息
type ContractConfig struct {
	ContractType    protocol.ContractType `json:"contractType"`    //基础合约类型
	ContractAddress common.Address        `json:"contractAddress"` //合约地址
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
	baseTrie     *trie.Trie                               //基础合约Hash树，Key为基础合约的Address Value为合约的类型
	contractTrie *trie.Trie                               //所有基础合约整体保存树，这里保存所有的基础合约信息，但是不保存基础合约的类型信息
	ethereum     *eth.Ethereum                            //以太坊对象
	accounts     *BokerAccount                            //播客链账号对象
	transactions *BokerTransaction                        //播客链交易对象
	db           ethdb.Database                           //数据库
	contracts    map[common.Address]protocol.ContractType //基础合约的Map
	services     ContractService                          //合约服务类
}

func NewContract(db ethdb.Database,
	ethereum *eth.Ethereum,
	transactions *BokerTransaction,
	bokerProto *protocol.BokerBackendProto) (*BokerContracts, error) {

	log.Info("****NewContract****")

	log.Info("Create Boker Base Trie", "Hash", bokerProto.BaseHash.String())
	baseTrie, errTrie := NewBaseTrie(bokerProto.BaseHash, db)
	if errTrie != nil {
		return nil, errTrie
	}

	log.Info("Create Boker Contracts Trie", "Hash", bokerProto.ContractHash.String())
	contractTrie, errTrie := NewContractTrie(bokerProto.ContractHash, db)
	if errTrie != nil {
		return nil, errTrie
	}

	//创建对象
	base := new(BokerContracts)
	base.baseTrie = baseTrie
	base.contractTrie = contractTrie
	base.contracts = make(map[common.Address]protocol.ContractType)
	base.db = db
	base.ethereum = ethereum
	base.transactions = transactions

	//从树中加载合约信息
	log.Info("Load Boker Base Contract Config")
	var err error
	if err = base.loadTrieContract(); err != nil {
		log.Error("Load Boker Base Contract Trie", "error", err)
		return base, nil
	}

	//创建投票合约服务
	log.Info("Check Vote Contract Exists")
	var address common.Address
	address, err = base.getContractAddress(protocol.ContractVote)
	if err != nil {

		log.Debug("Vote Contract is`t Exists")
		return base, nil
	}

	log.Info("New Vote Contract Service")
	base.services.votesContract, err = votes.NewVerifyVotesService(base.ethereum, address)
	if err != nil {
		//return nil, include.ErrNewContractService
		return base, nil
	}
	log.Info("Start Vote Contract Service")
	base.services.votesContract.Start()

	//创建分配通证合约服务
	log.Info("Check Assign Token Contract Exists")
	address, err = base.getContractAddress(protocol.ContractAssignToken)
	if err != nil {
		return base, nil
	}

	log.Info("New Token Contract Service")
	base.services.tokenContract, err = assigntoken.NewAssignTokenService(base.ethereum, address)
	if err != nil {
		//return nil, include.ErrNewContractService
		return base, nil
	}
	log.Info("Start Token Contract Service")
	base.services.tokenContract.Start()

	return base, nil
}

//创建基础合约的Hash树
func NewBaseTrie(root common.Hash, db ethdb.Database) (*trie.Trie, error) {
	return trie.NewTrieWithPrefix(root, protocol.BasePrefix, db)
}

func NewContractTrie(root common.Hash, db ethdb.Database) (*trie.Trie, error) {
	return trie.NewTrieWithPrefix(root, protocol.ContractPrefix, db)
}

//设置合约到Hash树中
func (c *BokerContracts) SetContract(address common.Address, contractType protocol.ContractType) error {

	//设置基础合约
	log.Info("****SetContract****")

	exist, err := c.existContract(address)
	if err != nil {
		return err
	}
	if exist {
		return protocol.ErrContractExist
	}

	if err = c.baseTrie.TryUpdate(address.Bytes(), []byte(strconv.Itoa(int(contractType)))); err != nil {
		return err
	}
	c.contracts[address] = contractType

	//更新基础合约列表保存树
	if err = c.setContractsTrie(); err != nil {
		return protocol.ErrSaveContractTrie
	}

	//判断是否需要启动合约
	if contractType == protocol.ContractVote {

		c.services.votesContract, err = votes.NewVerifyVotesService(c.ethereum, address)
		if err != nil {
			return protocol.ErrNewContractService
		}
		c.services.votesContract.Start()
	} else if contractType == protocol.ContractAssignToken {

		c.services.tokenContract, err = assigntoken.NewAssignTokenService(c.ethereum, address)
		if err != nil {
			return protocol.ErrNewContractService
		}
		c.services.tokenContract.Start()
	}
	return nil
}

//设置合约到Hash树中
func (c *BokerContracts) CancelContract(address common.Address) error {

	//检测合约是否存在
	log.Info("****CancelContract****")
	contractType, err := c.GetContract(address)
	if err != nil {
		return err
	}

	//从Hash树中删除
	if err := c.baseTrie.TryDelete(address.Bytes()); err != nil {
		return err
	}
	delete(c.contracts, address)

	//更新基础合约列表保存树
	if err = c.setContractsTrie(); err != nil {
		return protocol.ErrSaveContractTrie
	}

	//终止合约运行
	if contractType == protocol.ContractVote {

		if (c.services.votesContract != nil) && c.services.votesContract.IsStart() {

		}
	} else if contractType == protocol.ContractAssignToken {

		if (c.services.tokenContract != nil) && c.services.tokenContract.IsStart() {

		}
	}
	return nil
}

//将所有的合约地址进行设置
func (c *BokerContracts) setContractsTrie() error {

	//转换成切片数组
	log.Info("****setContractsTrie****")

	var contracts []common.Address = make([]common.Address, 0)
	for k, _ := range c.contracts {
		contracts = append(contracts, k)
	}
	contractsRLP, err := rlp.EncodeToBytes(contracts)
	if err != nil {

		log.Error("failed to encode contracts to rlp", "error", err)
		return err
	}
	c.contractTrie.Update(protocol.ContractPrefix, contractsRLP)
	return nil
}

//得到所有的合约地址
func (c *BokerContracts) getContractsTrie() ([]common.Address, error) {

	log.Info("****getContractsTrie****")

	var contracts []common.Address
	if c.contractTrie == nil {
		log.Error("contract Trie is nil")
		return []common.Address{}, protocol.ErrPointerIsNil
	}

	contractsRLP, err := c.contractTrie.TryGet(protocol.ContractPrefix)
	if err != nil {
		return []common.Address{}, err
	}

	if err := rlp.DecodeBytes(contractsRLP, &contracts); err != nil {
		log.Error("failed to decode contracts", "error", err)
		return []common.Address{}, err
	}

	return contracts, nil
}

//从Hash树中获取合约类型
func (c *BokerContracts) readContractType(address common.Address) (protocol.ContractType, error) {

	log.Info("****readContractType****")

	//根据地址获取数据
	key := address.Bytes()
	v, err := c.baseTrie.TryGet(key)
	if err != nil {
		return protocol.ContractBinary, err
	}

	//转换成交易类型
	contractType, err := strconv.Atoi(string(v[:]))
	return protocol.ContractType(contractType), nil
}

//根据合约类型得到合约地址
func (c *BokerContracts) getContractAddress(contractType protocol.ContractType) (common.Address, error) {

	if len(c.contracts) <= 0 {
		return common.Address{}, protocol.ErrNotFoundContract
	}

	for k, v := range c.contracts {

		if v == contractType {
			return k, nil
		}
	}
	return common.Address{}, protocol.ErrNotFoundContract
}

//根据合约地址得到合约类型
func (c *BokerContracts) getContractType(address common.Address) (protocol.ContractType, error) {

	contractType, exist := c.contracts[address]
	if exist {
		return contractType, nil
	} else {
		return protocol.ContractBinary, protocol.ErrNotFoundContract
	}
}

//加载基础合约信息
func (c *BokerContracts) loadTrieContract() error {

	log.Info("****loadContract****")

	//获取所有合约
	contracts, err := c.getContractsTrie()
	if err != nil {
		return err
	}
	log.Info("Load Boker Base Contract", "Size", len(contracts))
	if len(contracts) <= 0 {
		return nil
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
func (c *BokerContracts) GetContract(address common.Address) (protocol.ContractType, error) {

	if len(c.contracts) > 0 {
		value, exist := c.contracts[address]
		if exist {
			return value, nil
		}
	}
	return protocol.ContractBinary, protocol.ErrNotFoundContract
}

//判断此合约是否已经存在
func (c *BokerContracts) existContract(address common.Address) (bool, error) {

	if len(c.contracts) <= 0 {
		return false, protocol.ErrNotFoundContract
	}
	_, exist := c.contracts[address]
	if exist {
		return true, nil
	} else {
		return false, protocol.ErrNotFoundContract
	}
}

func (c *BokerContracts) GetContractTrie() (*trie.Trie, *trie.Trie) {

	return c.baseTrie, c.contractTrie
}
