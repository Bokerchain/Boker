package core

import (
	"fmt"
	"math/big"

	"github.com/boker/go-ethereum/boker/api"
	"github.com/boker/go-ethereum/boker/protocol"
	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/consensus/dpos"
	"github.com/boker/go-ethereum/consensus/ethash"
	"github.com/boker/go-ethereum/consensus/misc"
	"github.com/boker/go-ethereum/core/state"
	"github.com/boker/go-ethereum/core/types"
	"github.com/boker/go-ethereum/core/vm"
	"github.com/boker/go-ethereum/ethdb"
	"github.com/boker/go-ethereum/params"
)

// So we can deterministically seed different blockchains
var (
	canonicalSeed = 1
	forkSeed      = 2
)

// BlockGen creates blocks for testing.
// See GenerateChain for a detailed explanation.
type BlockGen struct {
	i       int
	parent  *types.Block
	chain   []*types.Block
	header  *types.Header
	statedb *state.StateDB

	gasPool  *GasPool
	txs      []*types.Transaction //交易数组
	receipts []*types.Receipt     //回执数组
	uncles   []*types.Header      //叔块数组

	config *params.ChainConfig
}

//设置生成的块的coinbase，只能调用一次。
func (b *BlockGen) SetCoinbase(addr common.Address) {
	if b.gasPool != nil {
		if len(b.txs) > 0 {
			panic("coinbase must be set before adding transactions")
		}
		panic("coinbase can only be set once")
	}
	b.header.Coinbase = addr
	b.gasPool = new(GasPool).AddGas(b.header.GasLimit)
}

//设置生成的块的额外数据字段
func (b *BlockGen) SetExtra(data []byte) {
	b.header.Extra = data
}

//添加一个交易到区块
func (b *BlockGen) AddTx(tx *types.Transaction, boker bokerapi.Api) {

	//判断gas池是否为nil
	if b.gasPool == nil {
		b.SetCoinbase(common.Address{})
	}

	//
	b.statedb.Prepare(tx.Hash(), common.Hash{}, len(b.txs))

	//应用交易，并返回回执
	receipt, _, err := ApplyTransaction(b.config,
		nil,
		nil,
		&b.header.Coinbase,
		b.gasPool,
		b.statedb,
		b.header,
		tx,
		b.header.GasUsed,
		vm.Config{},
		boker)

	//如果返回失败，则退出此协程
	if err != nil {
		panic(err)
	}

	//将交易加入到交易数组中
	b.txs = append(b.txs, tx)

	//将回执加入到回执数组中
	b.receipts = append(b.receipts, receipt)
}

//返回区块序号
func (b *BlockGen) Number() *big.Int {
	return new(big.Int).Set(b.header.Number)
}

func (b *BlockGen) AddUncheckedReceipt(receipt *types.Receipt) {
	b.receipts = append(b.receipts, receipt)
}

//根据账号地址返回下一个确认交易的Nonce，如果账号不存在则Panic
func (b *BlockGen) TxNonce(addr common.Address) uint64 {
	if !b.statedb.Exist(addr) {
		panic("account does not exist")
	}
	return b.statedb.GetNonce(addr)
}

//添加一个叔块
func (b *BlockGen) AddUncle(h *types.Header) {
	b.uncles = append(b.uncles, h)
}

//根据区块序号得到上一个区块（即父区块信息）
func (b *BlockGen) PrevBlock(index int) *types.Block {
	if index >= b.i {
		panic("block index out of range")
	}
	if index == -1 {
		return b.parent
	}
	return b.chain[index]
}

//修改区块时间，并计算区块难度（DPOS的难度始终是1）
func (b *BlockGen) OffsetTime(seconds int64) {
	b.header.Time.Add(b.header.Time, new(big.Int).SetInt64(seconds))
	if b.header.Time.Cmp(b.parent.Header().Time) <= 0 {
		panic("block time out of range")
	}
	b.header.Difficulty = ethash.CalcDifficulty(b.config, b.header.Time.Uint64(), b.parent.Header())
}

func GenerateChain(config *params.ChainConfig, parent *types.Block, db ethdb.Database, n int, gen func(int, *BlockGen)) ([]*types.Block, []types.Receipts) {

	//如果config为nil则加载Dpos配置信息
	if config == nil {
		config = params.DposChainConfig
	}

	//创建区块以及回执数组
	blocks, receipts := make(types.Blocks, n), make([]types.Receipts, n)

	//
	genblock := func(i int, h *types.Header, statedb *state.StateDB) (*types.Block, types.Receipts) {
		b := &BlockGen{parent: parent, i: i, chain: blocks, header: h, statedb: statedb, config: config}
		// Mutate the state and block according to any hard-fork specs
		if daoBlock := config.DAOForkBlock; daoBlock != nil {
			limit := new(big.Int).Add(daoBlock, params.DAOForkExtraRange)
			if h.Number.Cmp(daoBlock) >= 0 && h.Number.Cmp(limit) < 0 {
				if config.DAOForkSupport {
					h.Extra = common.CopyBytes(params.DAOForkBlockExtra)
				}
			}
		}
		if config.DAOForkSupport && config.DAOForkBlock != nil && config.DAOForkBlock.Cmp(h.Number) == 0 {
			misc.ApplyDAOHardFork(statedb)
		}
		// Execute any user modifications to the block and finalize it
		if gen != nil {
			gen(i, b)
		}

		//累计奖励
		dpos.AccumulateRewards(config, statedb, h, b.uncles)
		root, err := statedb.CommitTo(db, config.IsEIP158(h.Number))
		if err != nil {
			panic(fmt.Sprintf("state write error: %v", err))
		}
		h.Root = root
		h.DposProto = parent.Header().DposProto
		h.BokerProto = parent.Header().BokerProto
		return types.NewBlock(h, b.txs, b.uncles, b.receipts), b.receipts
	}

	//
	for i := 0; i < n; i++ {
		statedb, err := state.New(parent.Root(), state.NewDatabase(db))
		if err != nil {
			panic(err)
		}
		header := makeHeader(config, parent, statedb)
		block, receipt := genblock(i, header, statedb)
		blocks[i] = block
		receipts[i] = receipt
		parent = block
	}
	return blocks, receipts
}

//创建一个区块头
func makeHeader(config *params.ChainConfig, parent *types.Block, state *state.StateDB) *types.Header {

	//计算时间（这里注意以太坊允许有10秒中的时间差）
	var time *big.Int
	if parent.Time() == nil {
		time = big.NewInt(10)
	} else {
		time = new(big.Int).Add(parent.Time(), big.NewInt(10)) // block time is fixed at 10 seconds
	}

	//创建区块头
	return &types.Header{
		Root:       state.IntermediateRoot(config.IsEIP158(parent.Number())),
		ParentHash: parent.Hash(),
		Coinbase:   parent.Coinbase(),
		Difficulty: parent.Difficulty(),
		DposProto:  &types.DposContextProto{},
		BokerProto: &protocol.BokerBackendProto{},
		GasLimit:   CalcGasLimit(parent),
		GasUsed:    new(big.Int),
		Number:     new(big.Int).Add(parent.Number(), common.Big1),
		Time:       time,
	}
}

// newCanonical creates a chain database, and injects a deterministic canonical
// chain. Depending on the full flag, if creates either a full block chain or a
// header only chain.
/*func newCanonical(n int, full bool) (ethdb.Database, *BlockChain, error) {
	// Initialize a fresh chain with only a genesis block
	gspec := new(Genesis)
	db, _ := ethdb.NewMemDatabase()
	genesis := gspec.MustCommit(db)

	blockchain, _ := NewBlockChain(db, params.AllEthashProtocolChanges, ethash.NewFaker(), vm.Config{})
	// Create and inject the requested chain
	if n == 0 {
		return db, blockchain, nil
	}
	if full {
		// Full block-chain requested
		blocks := makeBlockChain(genesis, n, db, canonicalSeed)
		_, err := blockchain.InsertChain(blocks)
		return db, blockchain, err
	}
	// Header-only chain requested
	headers := makeHeaderChain(genesis.Header(), n, db, canonicalSeed)
	_, err := blockchain.InsertHeaderChain(headers, 1)
	return db, blockchain, err
}*/

//在父区块根节点创建一个确定性的链
func makeHeaderChain(parent *types.Header, n int, db ethdb.Database, seed int) []*types.Header {
	blocks := makeBlockChain(types.NewBlockWithHeader(parent), n, db, seed)
	headers := make([]*types.Header, len(blocks))
	for i, block := range blocks {
		headers[i] = block.Header()
	}
	return headers
}

// makeBlockChain creates a deterministic chain of blocks rooted at parent.
func makeBlockChain(parent *types.Block, n int, db ethdb.Database, seed int) []*types.Block {
	blocks, _ := GenerateChain(params.DposChainConfig, parent, db, n, func(i int, b *BlockGen) {
		b.SetCoinbase(common.Address{0: byte(seed), 19: byte(i)})
	})
	return blocks
}
