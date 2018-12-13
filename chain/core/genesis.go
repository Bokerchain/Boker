package core

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/Bokerchain/Boker/chain/boker/protocol"
	"github.com/Bokerchain/Boker/chain/common"
	"github.com/Bokerchain/Boker/chain/common/hexutil"
	"github.com/Bokerchain/Boker/chain/common/math"
	"github.com/Bokerchain/Boker/chain/core/state"
	"github.com/Bokerchain/Boker/chain/core/types"
	"github.com/Bokerchain/Boker/chain/ethdb"
	"github.com/Bokerchain/Boker/chain/log"
	"github.com/Bokerchain/Boker/chain/params"
	"github.com/Bokerchain/Boker/chain/rlp"
	"github.com/Bokerchain/Boker/chain/trie"
)

//go:generate gencodec -type Genesis -field-override genesisSpecMarshaling -out gen_genesis.go
//go:generate gencodec -type GenesisAccount -field-override genesisAccountMarshaling -out gen_genesis_account.go

var errGenesisNoConfig = errors.New("genesis has no chain configuration")

//创世区块配置定义
type Genesis struct {
	Config     *params.ChainConfig `json:"config"`
	Nonce      uint64              `json:"nonce"`
	Timestamp  uint64              `json:"timestamp"`
	ExtraData  []byte              `json:"extraData"`
	GasLimit   uint64              `json:"gasLimit"   gencodec:"required"`
	Difficulty *big.Int            `json:"difficulty" gencodec:"required"`
	Mixhash    common.Hash         `json:"mixHash"`
	Coinbase   common.Address      `json:"coinbase"`
	Alloc      GenesisAlloc        `json:"alloc"      gencodec:"required"`

	//这些字段用于一致性测试，请不要使用它们在实际的创世块中.
	Number     uint64      `json:"number"`
	GasUsed    uint64      `json:"gasUsed"`
	ParentHash common.Hash `json:"parentHash"`
}

//Json格式反序列化
func (ga *GenesisAlloc) UnmarshalJSON(data []byte) error {

	m := make(map[common.UnprefixedAddress]GenesisAccount)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	*ga = make(GenesisAlloc)
	for addr, a := range m {
		(*ga)[common.Address(addr)] = a
	}
	return nil
}

//定义创世区块中的账号信息
type GenesisAlloc map[common.Address]GenesisAccount
type GenesisAccount struct {
	Code       []byte                      `json:"code,omitempty"`
	Storage    map[common.Hash]common.Hash `json:"storage,omitempty"`
	Balance    *big.Int                    `json:"balance" gencodec:"required"` //Bobby数量
	Nonce      uint64                      `json:"nonce,omitempty"`             //用户的Nonce值
	PrivateKey []byte                      `json:"secretKey,omitempty"`         //用户私钥
}

// field type overrides for gencodec
type genesisSpecMarshaling struct {
	Nonce      math.HexOrDecimal64
	Timestamp  math.HexOrDecimal64
	ExtraData  hexutil.Bytes
	GasLimit   math.HexOrDecimal64
	GasUsed    math.HexOrDecimal64
	Number     math.HexOrDecimal64
	Difficulty *math.HexOrDecimal256
	Alloc      map[common.UnprefixedAddress]GenesisAccount
}

type genesisAccountMarshaling struct {
	Code       hexutil.Bytes
	Balance    *math.HexOrDecimal256
	Nonce      math.HexOrDecimal64
	Storage    map[storageJSON]storageJSON
	PrivateKey hexutil.Bytes
}

// storageJSON represents a 256 bit byte array, but allows less than 256 bits when
// unmarshaling from hex.
type storageJSON common.Hash

func (h *storageJSON) UnmarshalText(text []byte) error {

	//得到账号信息
	text = bytes.TrimPrefix(text, []byte("0x"))
	if len(text) > 64 {
		return fmt.Errorf("too many hex characters in storage key/value %q", text)
	}
	offset := len(h) - len(text)/2 // pad on the left
	if _, err := hex.Decode(h[offset:], text); err != nil {
		fmt.Println(err)
		return fmt.Errorf("invalid hex storage key/value %q", text)
	}
	return nil
}

func (h storageJSON) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}

// GenesisMismatchError is raised when trying to overwrite an existing
// genesis block with an incompatible one.
type GenesisMismatchError struct {
	Stored, New common.Hash
}

func (e *GenesisMismatchError) Error() string {
	return fmt.Sprintf("database already contains an incompatible genesis block (have %x, new %x)", e.Stored[:8], e.New[:8])
}

// SetupGenesisBlock writes or updates the genesis block in db.
// The block that will be used is:
//
//                          genesis == nil       genesis != nil
//                       +------------------------------------------
//     db has no genesis |  main-net default  |  genesis
//     db has genesis    |  from DB           |  genesis (if compatible)
//
// The stored chain configuration will be updated if it is compatible (i.e. does not
// specify a fork block below the local head block). In case of a conflict, the
// error is a *params.ConfigCompatError and the new, unwritten config is returned.
//
// The returned chain configuration is never nil.
func SetupGenesisBlock(db ethdb.Database, genesis *Genesis) (*params.ChainConfig, common.Hash, error) {

	log.Info("****SetupGenesisBlock****")

	if genesis != nil && genesis.Config == nil {
		return params.DposChainConfig, common.Hash{}, errGenesisNoConfig
	}

	//如果没有存储的genesis块，只需提交新块
	stored := GetCanonicalHash(db, 0)
	if (stored == common.Hash{}) {
		if genesis == nil {
			log.Info("Writing default main-net genesis block")
			genesis = DefaultGenesisBlock()
		} else {
			log.Info("Writing custom genesis block")
		}
		block, err := genesis.Commit(db)

		if err != nil {
			log.Error("Genesis Commit", "error", err)
		}
		log.Info("Genesis Commit", "Number", block.Number())

		return genesis.Config, block.Hash(), err
	}
	log.Info("GetCanonicalHash")

	// Check whether the genesis block is already written.
	if genesis != nil {
		block, _, _, _, _ := genesis.ToBlock()
		hash := block.Hash()
		if hash != stored {

			log.Info("Genesis ToBlock")
			return genesis.Config, block.Hash(), &GenesisMismatchError{stored, hash}
		}
	}
	log.Info("ToBlock")

	// Get the existing chain configuration.
	newcfg := genesis.configOrDefault(stored)
	storedcfg, err := GetChainConfig(db, stored)
	if err != nil {
		if err == ErrChainConfigNotFound {
			// This case happens if a genesis write was interrupted.
			log.Warn("Found genesis block without chain config")
			err = WriteChainConfig(db, stored, newcfg)
		}
		return newcfg, stored, err
	}
	log.Info("GetChainConfig")

	// Special case: don't change the existing config of a non-mainnet chain if no new
	// config is supplied. These chains would get AllProtocolChanges (and a compat error)
	// if we just continued here.
	if genesis == nil && stored != params.MainnetGenesisHash {
		return storedcfg, stored, nil
	}

	// Check config compatibility and write the config. Compatibility errors
	// are returned to the caller unless we're already at block zero.
	height := GetBlockNumber(db, GetHeadHeaderHash(db))
	if height == missingNumber {
		return newcfg, stored, fmt.Errorf("missing block number for head header hash")
	}
	log.Info("GetBlockNumber")

	compatErr := storedcfg.CheckCompatible(newcfg, height)
	if compatErr != nil && height != 0 && compatErr.RewindTo != 0 {
		return newcfg, stored, compatErr
	}
	log.Info("CheckCompatible")

	return newcfg, stored, WriteChainConfig(db, stored, newcfg)
}

func (g *Genesis) configOrDefault(ghash common.Hash) *params.ChainConfig {
	switch {
	case g != nil:
		return g.Config
	default:
		return params.DposChainConfig
	}
}

//创建一个特定的创世区块状态
func (g *Genesis) ToBlock() (*types.Block, *state.StateDB, *trie.Trie, *trie.Trie, *trie.Trie) {

	db, _ := ethdb.NewMemDatabase()
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(db))

	//创建导入的账号信息
	for addr, account := range g.Alloc {
		statedb.AddBalance(addr, account.Balance)
		statedb.SetCode(addr, account.Code)
		statedb.SetNonce(addr, account.Nonce)
		for key, value := range account.Storage {
			statedb.SetState(addr, key, value)
		}
	}
	root := statedb.IntermediateRoot(false)

	//添加Dpos配置
	dposContext := initGenesisDposContext(g, db)
	dposContextProto := dposContext.ToProto()
	log.Info("ToProto", "root", dposContextProto.Root().String())

	//添加播客链的设置
	singleTrie, contractsTrie, abiTrie, err := initBoker(db)
	if err != nil {
		fmt.Errorf("initGenesisBoker error")
		return nil, statedb, nil, nil, nil
	}
	bokerProto := protocol.ToBokerProto(singleTrie.Hash(), contractsTrie.Hash(), abiTrie.Hash())
	log.Info("ToBokerProto", "root", bokerProto.Root().String())

	head := &types.Header{
		Number:     new(big.Int).SetUint64(g.Number),
		Nonce:      types.EncodeNonce(g.Nonce),
		Time:       new(big.Int).SetUint64(g.Timestamp),
		ParentHash: g.ParentHash,
		Extra:      g.ExtraData,
		GasLimit:   new(big.Int).SetUint64(g.GasLimit),
		GasUsed:    new(big.Int).SetUint64(g.GasUsed),
		Difficulty: g.Difficulty,
		MixDigest:  g.Mixhash,
		Coinbase:   g.Coinbase,
		Root:       root,
		DposProto:  dposContextProto,
		BokerProto: bokerProto,
	}
	if g.GasLimit == 0 {
		head.GasLimit = params.GenesisGasLimit
	}
	if g.Difficulty == nil {
		head.Difficulty = params.GenesisDifficulty
	}

	block := types.NewBlock(head, nil, nil, nil)
	block.DposContext = dposContext

	return block, statedb, singleTrie, contractsTrie, abiTrie
}

// Commit writes the block and state of a genesis specification to the database.
// The block is committed as the canonical head block.
func (g *Genesis) Commit(db ethdb.Database) (*types.Block, error) {

	block, statedb, singleTrie, contractsTrie, abiTrie := g.ToBlock()

	// add dposcontext
	if _, err := block.DposContext.CommitTo(db); err != nil {
		return nil, err
	}
	//新增播客链数据保存
	if err := commitBoker(singleTrie, contractsTrie, abiTrie, db); err != nil {
		return nil, err
	}
	log.Info("Write Boker ", "singleTrieHash", singleTrie.Hash().String(), "ContractsHash", contractsTrie.Hash().String(), "abiTrieHash", abiTrie.Hash().String())

	if block.Number().Sign() != 0 {
		return nil, fmt.Errorf("can't commit genesis block with number > 0")
	}
	if _, err := statedb.CommitTo(db, false); err != nil {
		return nil, fmt.Errorf("cannot write state: %v", err)
	}
	if err := WriteTd(db, block.Hash(), block.NumberU64(), g.Difficulty); err != nil {
		return nil, err
	}
	if err := WriteBlock(db, block); err != nil {
		return nil, err
	}
	if err := WriteBlockReceipts(db, block.Hash(), block.NumberU64(), nil); err != nil {
		return nil, err
	}
	if err := WriteCanonicalHash(db, block.Hash(), block.NumberU64()); err != nil {
		return nil, err
	}
	if err := WriteHeadBlockHash(db, block.Hash()); err != nil {
		return nil, err
	}
	if err := WriteHeadHeaderHash(db, block.Hash()); err != nil {
		return nil, err
	}
	config := g.Config
	if config == nil {
		config = params.DposChainConfig
	}
	return block, WriteChainConfig(db, block.Hash(), config)
}

// GenesisBlockForTesting creates and writes a block in which addr has the given wei balance.
func GenesisBlockForTesting(db ethdb.Database, addr common.Address, balance *big.Int) *types.Block {
	g := Genesis{Alloc: GenesisAlloc{addr: {Balance: balance}}}
	return g.MustCommit(db)
}

// MustCommit writes the genesis block and state to db, panicking on error.
// The block is committed as the canonical head block.
func (g *Genesis) MustCommit(db ethdb.Database) *types.Block {
	block, err := g.Commit(db)
	if err != nil {
		panic(err)
	}
	return block
}

func DefaultGenesisBlock() *Genesis {
	return &Genesis{
		Config:     params.DposChainConfig,
		Nonce:      66,
		Timestamp:  1522052340,
		ExtraData:  hexutil.MustDecode("0x11bbe8db4e347b4e8c937c1c8370e4b5ed33adb3db69cbdb7a38e1e50b1b82fa"),
		GasLimit:   4712388,
		Difficulty: big.NewInt(17179869184),
		Alloc:      decodePrealloc(mainnetAllocData),
	}
}

func decodePrealloc(data string) GenesisAlloc {
	var p []struct{ Addr, Balance *big.Int }
	if err := rlp.NewStream(strings.NewReader(data), 0).Decode(&p); err != nil {
		panic(err)
	}
	ga := make(GenesisAlloc, len(p))
	for _, account := range p {
		ga[common.BigToAddress(account.Addr)] = GenesisAccount{Balance: account.Balance}
	}
	return ga
}

//DPOS的初始化设置
func initGenesisDposContext(g *Genesis, db ethdb.Database) *types.DposContext {

	dc, err := types.NewDposContextFromProto(db, &types.DposContextProto{})
	if err != nil {
		return nil
	}

	//由于第一次创建，因此需要提交一次周期树
	var validators []common.Address = make([]common.Address, 0)
	dc.SetEpochTrie(validators)

	var producers []common.Address
	validatorsRLP := dc.EpochTrie().Get(protocol.ValidatorsKey)
	if err := rlp.DecodeBytes(validatorsRLP, &producers); err != nil {

		log.Info("failed to decode validators", "error", err)
		return nil
	}

	return dc
}

//****创建播客链相关Hash树信息****//
func initBoker(db ethdb.Database) (*trie.Trie, *trie.Trie, *trie.Trie, error) {

	log.Info("****initBoker****")

	var singleTrie, contractsTrie, abiTrie *trie.Trie
	var err error
	var root common.Hash

	//创建基础合约树
	if singleTrie, err = trie.NewTrieWithPrefix(root, protocol.SinglePrefix, db); err != nil {
		return nil, nil, nil, err
	}
	var bases []common.Address = make([]common.Address, 0)
	basesRLP, err := rlp.EncodeToBytes(bases)
	if err != nil {
		log.Error("failed to encode contract to rlp", "error", err)
		return nil, nil, nil, err
	}
	singleTrie.Update(protocol.SinglePrefix, basesRLP)

	//创建合约总和树
	if contractsTrie, err = trie.NewTrieWithPrefix(root, protocol.Contracts, db); err != nil {
		return nil, nil, nil, err
	}
	var contracts []common.Address = make([]common.Address, 0)
	contractsRLP, err := rlp.EncodeToBytes(contracts)
	if err != nil {
		log.Error("failed to encode contracts to rlp", "error", err)
		return nil, nil, nil, err
	}
	contractsTrie.Update(protocol.Contracts, contractsRLP)

	//创建合约abi树
	if abiTrie, err = trie.NewTrieWithPrefix(root, protocol.AbiPrefix, db); err != nil {
		return nil, nil, nil, err
	}
	var abi []common.Address = make([]common.Address, 0)
	abiRLP, err := rlp.EncodeToBytes(abi)
	if err != nil {
		log.Error("failed to encode contracts to rlp", "error", err)
		return nil, nil, nil, err
	}
	abiTrie.Update(protocol.AbiPrefix, abiRLP)

	return singleTrie, contractsTrie, abiTrie, nil
}

//****创建播客链相关Hash树信息****//
func commitBoker(singleTrie *trie.Trie, contractsTrie *trie.Trie, abiTrie *trie.Trie, db ethdb.Database) error {

	log.Info("****commitBoker****")

	if _, err := singleTrie.CommitTo(db); err != nil {
		return err
	}
	if _, err := contractsTrie.CommitTo(db); err != nil {
		return err
	}
	if _, err := abiTrie.CommitTo(db); err != nil {
		return err
	}
	return nil
}
