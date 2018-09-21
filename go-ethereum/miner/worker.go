package miner

import (
	"bytes"
	"fmt"
	_ "io/ioutil"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/boker/go-ethereum/boker/protocol"
	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/consensus"
	"github.com/boker/go-ethereum/consensus/dpos"
	"github.com/boker/go-ethereum/consensus/misc"
	"github.com/boker/go-ethereum/core"
	"github.com/boker/go-ethereum/core/state"
	"github.com/boker/go-ethereum/core/types"
	"github.com/boker/go-ethereum/core/vm"
	"github.com/boker/go-ethereum/ethdb"
	"github.com/boker/go-ethereum/event"
	"github.com/boker/go-ethereum/log"
	"github.com/boker/go-ethereum/params"
	"github.com/boker/go-ethereum/trie"
	"gopkg.in/fatih/set.v0"
)

const (
	resultQueueSize  = 10
	miningLogAtDepth = 5

	// txChanSize is the size of channel listening to TxPreEvent.
	// The number is referenced from the size of tx pool.
	txChanSize = 4096
	// chainHeadChanSize is the size of channel listening to ChainHeadEvent.
	chainHeadChanSize = 10

	// chainSideChanSize is the size of channel listening to ChainSideEvent.
	chainSideChanSize = 10
)

// Work is the workers current environment and holds
// all of the current state information
type Work struct {
	config      *params.ChainConfig
	signer      types.Signer
	state       *state.StateDB // apply state changes here
	dposContext *types.DposContext
	ancestors   *set.Set     // ancestor set (used for checking uncle parent validity)
	family      *set.Set     // family set (used for checking uncle invalidity)
	uncles      *set.Set     // uncle set
	tcount      int          // tx count in cycle
	Block       *types.Block // the new block
	header      *types.Header
	txs         []*types.Transaction
	receipts    []*types.Receipt
	createdAt   time.Time
}

type Result struct {
	Work  *Work
	Block *types.Block
}

// worker is the main object which takes care of applying messages to the new state
type worker struct {
	config         *params.ChainConfig
	engine         consensus.Engine
	mu             sync.Mutex
	mux            *event.TypeMux
	txCh           chan core.TxPreEvent
	txSub          event.Subscription
	chainHeadCh    chan core.ChainHeadEvent
	chainHeadSub   event.Subscription
	wg             sync.WaitGroup
	recv           chan *Result
	eth            Backend
	chain          *core.BlockChain
	proc           core.Validator
	chainDb        ethdb.Database
	coinbase       common.Address
	extra          []byte
	currentMu      sync.Mutex
	current        *Work
	uncleMu        sync.Mutex
	possibleUncles map[common.Hash]*types.Block
	unconfirmed    *unconfirmedBlocks // set of locally mined blocks pending canonicalness confirmations
	mining         int32
	atWork         int32
	quitCh         chan struct{}
	stopper        chan struct{}
	isStart        bool
}

func newWorker(config *params.ChainConfig, engine consensus.Engine, coinbase common.Address, eth Backend, mux *event.TypeMux) *worker {

	//创建一个矿工
	worker := &worker{
		config:         config,
		engine:         engine,
		eth:            eth,
		mux:            mux,
		txCh:           make(chan core.TxPreEvent, txChanSize),
		chainHeadCh:    make(chan core.ChainHeadEvent, chainHeadChanSize),
		chainDb:        eth.ChainDb(),
		recv:           make(chan *Result, resultQueueSize),
		chain:          eth.BlockChain(),
		proc:           eth.BlockChain().Validator(),
		possibleUncles: make(map[common.Hash]*types.Block),
		coinbase:       coinbase,
		unconfirmed:    newUnconfirmedBlocks(eth.BlockChain(), miningLogAtDepth),
		quitCh:         make(chan struct{}, 1),
		stopper:        make(chan struct{}, 1),
		isStart:        false,
	}

	//订阅交易池的TxPreEvent事件
	worker.txSub = eth.TxPool().SubscribeTxPreEvent(worker.txCh)

	//订阅区块链的事件
	worker.chainHeadSub = eth.BlockChain().SubscribeChainHeadEvent(worker.chainHeadCh)

	go worker.update()
	go worker.wait()
	//worker.createNewWork()

	return worker
}

func (self *worker) CreateNewWork() {
	if !self.isStart {
		self.isStart = true
		self.createNewWork()
	}
}

func (self *worker) setCoinbase(addr common.Address) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.coinbase = addr
}

func (self *worker) setExtra(extra []byte) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.extra = extra
}

func (self *worker) pending() (*types.Block, *state.StateDB) {
	self.currentMu.Lock()
	defer self.currentMu.Unlock()

	if atomic.LoadInt32(&self.mining) == 0 {
		return types.NewBlock(
			self.current.header,
			self.current.txs,
			nil,
			self.current.receipts,
		), self.current.state.Copy()
	}
	return self.current.Block, self.current.state.Copy()
}

func (self *worker) pendingBlock() *types.Block {
	self.currentMu.Lock()
	defer self.currentMu.Unlock()

	if atomic.LoadInt32(&self.mining) == 0 {
		return types.NewBlock(
			self.current.header,
			self.current.txs,
			nil,
			self.current.receipts,
		)
	}
	return self.current.Block
}

//开始启动挖矿
func (self *worker) start() {
	self.mu.Lock()
	defer self.mu.Unlock()

	atomic.StoreInt32(&self.mining, 1)
	go self.mintLoop()
}

//矿工挖矿
func (self *worker) mintBlock(now int64) {

	//得到挖矿使用的共识引擎
	engine, ok := self.engine.(*dpos.Dpos)
	if !ok {
		log.Error("Only the dpos engine was allowed")
		return
	}

	//获取当前出块节点的数量
	size, sizeErr := engine.GetProducerSize(self.chain.CurrentBlock(), self.coinbase)
	if sizeErr != nil {
		return
	}
	if size == 0 {
		if self.chain.CurrentBlock().Number().Uint64() != 0 {
			log.Error("current block number is`t zero")
			return
		}
		if err := engine.CheckDeadline(self.chain.CurrentBlock(), now); err != nil {
			return
		}
		if self.chain.Boker().IsValidator(self.coinbase) {

			work, err := self.createNewWork()
			if err != nil {
				log.Error("Failed to create the new work", "err", err)
				return
			}

			//对区块进行封包处理
			result, err := self.engine.Seal(self.chain, work.Block, self.quitCh)
			if err != nil {
				log.Error("Failed to seal the block", "err", err)
				return
			}
			self.recv <- &Result{work, result}
		} else {
			log.Error("current coinbase is`t special account", "coinbase", self.coinbase)
		}

	} else {

		if err := engine.CheckDeadline(self.chain.CurrentBlock(), now); err != nil {
			return
		}

		err := engine.CheckProducer(self.chain.CurrentBlock(), now)
		if err != nil {
			switch err {
			case dpos.ErrWaitForPrevBlock,
				dpos.ErrMintFutureBlock,
				dpos.ErrInvalidBlockProducer,
				protocol.ErrInvalidMintBlockTime:
				{
					log.Error("Failed to mint the block, while ", "err", err)
					break
				}
			default:
				{
					log.Error("Failed to mint the block", "err", err)
					break
				}
			}
			return
		}

		//可以进行挖矿出块,创建一次挖矿矿工
		work, err := self.createNewWork()
		if err != nil {
			log.Error("Failed to create the new work", "err", err)
			return
		}
		//对区块进行封包处理
		result, err := self.engine.Seal(self.chain, work.Block, self.quitCh)
		if err != nil {
			log.Error("Failed to seal the block", "err", err)
			return
		}
		self.recv <- &Result{work, result}
	}
}

//矿工挖矿循环
func (self *worker) mintLoop() {

	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case now := <-ticker:
			self.mintBlock(now.Unix())
		case <-self.stopper:
			close(self.quitCh)
			self.quitCh = make(chan struct{}, 1)
			self.stopper = make(chan struct{}, 1)
			return
		}
	}

}

func (self *worker) stop() {
	if atomic.LoadInt32(&self.mining) == 0 {
		return
	}

	self.wg.Wait()

	self.mu.Lock()
	defer self.mu.Unlock()

	atomic.StoreInt32(&self.mining, 0)
	atomic.StoreInt32(&self.atWork, 0)
	close(self.stopper)
}

func (self *worker) update() {

	defer self.txSub.Unsubscribe()
	defer self.chainHeadSub.Unsubscribe()

	for {
		// A real event arrived, process interesting content
		select {
		// Handle ChainHeadEvent
		case <-self.chainHeadCh:
			close(self.quitCh)
			self.quitCh = make(chan struct{}, 1)

		// Handle TxPreEvent
		case ev := <-self.txCh:
			// Apply transaction to the pending state if we're not mining
			if atomic.LoadInt32(&self.mining) == 0 {
				self.currentMu.Lock()
				acc, _ := types.Sender(self.current.signer, ev.Tx)
				txs := map[common.Address]types.Transactions{acc: {ev.Tx}}
				txset := types.NewTransactionsByPriceAndNonce(self.current.signer, txs)

				self.current.commitTransactions(self.mux, txset, self.chain, self.coinbase)
				self.currentMu.Unlock()
			}
		// System stopped
		case <-self.txSub.Err():
			return
		case <-self.chainHeadSub.Err():
			return
		}
	}
}

func (self *worker) wait() {

	for {
		for result := range self.recv {
			atomic.AddInt32(&self.atWork, -1)

			if result == nil || result.Block == nil {
				continue
			}
			block := result.Block
			work := result.Work

			// Update the block hash in all logs since it is now available and not when the
			// receipt/log of individual transactions were created.
			for _, r := range work.receipts {
				for _, l := range r.Logs {
					l.BlockHash = block.Hash()
				}
			}
			for _, log := range work.state.Logs() {
				log.BlockHash = block.Hash()
			}

			//将区块和状态信息写入数据库
			stat, err := self.chain.WriteBlockAndState(block, work.receipts, work.state)
			if err != nil {
				log.Error("Failed writing block to chain", "err", err)
				continue
			}
			// check if canon block and write transactions
			if stat == core.CanonStatTy {
				// implicit by posting ChainHeadEvent
			}

			//广播块并宣布链插入事件(发送这个事件是为了把新挖出的区块广播给其他结点，事件处理代码位于eth/handler.go 中的 minedBroadcastLoop)
			self.mux.Post(core.NewMinedBlockEvent{Block: block})

			//发送ChainEvent事件
			var (
				events []interface{}
				logs   = work.state.Logs()
			)
			events = append(events, core.ChainEvent{Block: block, Hash: block.Hash(), Logs: logs})
			if stat == core.CanonStatTy {
				events = append(events, core.ChainHeadEvent{Block: block})
			}
			self.chain.PostChainEvents(events, logs)

			//将块插入待处理组中以等待确认
			self.unconfirmed.Insert(block.NumberU64(), block.Hash())
			//log.Info("Successfully sealed new block", "number", block.Number(), "hash", block.Hash())
		}
	}
}

func newBokerFromProto(db ethdb.Database, bokerProto *protocol.BokerBackendProto) (*trie.Trie, *trie.Trie, error) {

	//log.Info("****newBokerFromProto****")

	baseTrie, err := trie.NewTrieWithPrefix(bokerProto.BaseHash, protocol.BasePrefix, db)
	if err != nil {
		return nil, nil, err
	}

	contractTrie, err := trie.NewTrieWithPrefix(bokerProto.ContractHash, protocol.ContractPrefix, db)
	if err != nil {
		return nil, nil, err
	}
	return baseTrie, contractTrie, nil
}

//为当前周期创建一个新环境。
func (self *worker) makeCurrent(parent *types.Block, header *types.Header) error {

	//log.Info("****makeCurrent****")

	//根据父块状态创建一个新的StateDB实例
	state, err := self.chain.StateAt(parent.Root())
	if err != nil {
		return err
	}

	//创建一个共识实例
	dposContext, err := types.NewDposContextFromProto(self.chainDb, parent.Header().DposProto)
	if err != nil {
		return err
	}

	//加载播客链信息
	/*baseTrie, contractTrie, err := newBokerFromProto(self.chainDb, parent.Header().BokerProto)
	if err != nil {
		return err
	}*/
	/*log.Info("Read Boker",
	"BokerProto", parent.Header().BokerProto.Root().String(),
	"base", baseTrie.Hash().String(),
	"contract", contractTrie.Hash().String())*/

	//创建一个work实例
	work := &Work{
		config: self.config,
		//signer:      types.NewEIP155Signer(self.config.ChainId),
		signer:      types.HomesteadSigner{},
		state:       state,
		dposContext: dposContext,
		ancestors:   set.New(),
		family:      set.New(),
		uncles:      set.New(),
		header:      header,
		createdAt:   time.Now(),
	}

	// when 08 is processed ancestors contain 07 (quick block)
	for _, ancestor := range self.chain.GetBlocksFromHash(parent.Hash(), 7) {
		for _, uncle := range ancestor.Uncles() {
			work.family.Add(uncle.Hash())
		}
		work.family.Add(ancestor.Hash())
		work.ancestors.Add(ancestor.Hash())
	}

	// Keep track of transactions which return errors so they can be removed
	work.tcount = 0
	self.current = work
	return nil
}

func (self *worker) createNewWork() (*Work, error) {

	//初始化各种临界区
	self.mu.Lock()
	defer self.mu.Unlock()
	self.uncleMu.Lock()
	defer self.uncleMu.Unlock()
	self.currentMu.Lock()
	defer self.currentMu.Unlock()

	//起始时间以及父块
	tstart := time.Now()
	parent := self.chain.CurrentBlock()

	tstamp := tstart.Unix()
	if parent.Time().Cmp(new(big.Int).SetInt64(tstamp)) >= 0 {
		tstamp = parent.Time().Int64() + 1
	}

	if now := time.Now().Unix(); tstamp > now+1 {
		wait := time.Duration(tstamp-now) * time.Second
		log.Info("Mining too far in the future", "wait", common.PrettyDuration(wait))
		time.Sleep(wait)
	}

	//获取块头信息
	num := parent.Number()
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     num.Add(num, common.Big1),
		GasLimit:   core.CalcGasLimit(parent),
		GasUsed:    new(big.Int),
		Extra:      self.extra,
		Time:       big.NewInt(tstamp),
	}

	//如果我们正在挖掘，只设置coinbase（避免出现虚假区块奖励）
	if atomic.LoadInt32(&self.mining) == 1 {
		header.Coinbase = self.coinbase
	}

	//初始化共识引擎
	if err := self.engine.Prepare(self.chain, header); err != nil {
		return nil, fmt.Errorf("got error when preparing header, err: %s", err)
	}

	// If we are care about TheDAO hard-fork check whether to override the extra-data or not
	if daoBlock := self.config.DAOForkBlock; daoBlock != nil {
		// Check whether the block is among the fork extra-override range
		limit := new(big.Int).Add(daoBlock, params.DAOForkExtraRange)
		if header.Number.Cmp(daoBlock) >= 0 && header.Number.Cmp(limit) < 0 {
			// Depending whether we support or oppose the fork, override differently
			if self.config.DAOForkSupport {
				header.Extra = common.CopyBytes(params.DAOForkBlockExtra)
			} else if bytes.Equal(header.Extra, params.DAOForkBlockExtra) {
				header.Extra = []byte{} // If miner opposes, don't let it use the reserved extra-data
			}
		}
	}

	// Could potentially happen if starting to mine in an odd state.
	//log.Info("createNewWork makeCurrent", "number", header.Number)
	err := self.makeCurrent(parent, header)
	if err != nil {
		return nil, fmt.Errorf("got error when create mining context, err: %s", err)
	}

	//创建当前工作任务并检查所需的fork交易
	work := self.current
	if self.config.DAOForkSupport && self.config.DAOForkBlock != nil && self.config.DAOForkBlock.Cmp(header.Number) == 0 {
		misc.ApplyDAOHardFork(work.state)
	}

	//获取txpool的待处理交易列表的一个拷贝
	pending, err := self.eth.TxPool().Pending()
	if err != nil {
		return nil, fmt.Errorf("got error when fetch pending transactions, err: %s", err)
	}

	//将待处理交易列表,封装进一个TransactionsByPriceAndNonce类型的结构中。这个结构中包含一个heads字段，把交易按照gas price进行排序
	txs := types.NewTransactionsByPriceAndNonce(self.current.signer, pending)

	//调用commitTransactions把交易提交到EVM去执行
	work.commitTransactions(self.mux, txs, self.chain, self.coinbase)

	//遍历所有叔块
	var (
		uncles    []*types.Header
		badUncles []common.Hash
	)
	for hash, uncle := range self.possibleUncles {

		if len(uncles) == 2 {
			break
		}

		//调用commitUncle()把叔块header的hash添加进Work.uncles集合中(以太坊规定每个区块最多打包2个叔块的header)
		if err := self.commitUncle(work, uncle.Header()); err != nil {
			log.Trace("Bad uncle found and will be removed", "hash", hash)
			log.Trace(fmt.Sprint(uncle))

			badUncles = append(badUncles, hash)
		} else {
			log.Debug("Committing new uncle to block", "hash", hash)
			uncles = append(uncles, uncle.Header())
		}
	}
	for _, hash := range badUncles {
		delete(self.possibleUncles, hash)
	}

	//使用共识引擎打包新区块
	if self.eth == nil {
		log.Info("createNewWork check eth is nil")
	}
	if work.Block, err = self.engine.Finalize(self.chain, header, work.state, work.txs, uncles, work.receipts, work.dposContext, self.eth.Boker()); err != nil {
		return nil, fmt.Errorf("got error when finalize block for sealing, err: %s", err)
	}
	work.Block.DposContext = work.dposContext

	//更新新块的矿工数量,如果我们正在挖矿，那我们只打印日志
	if atomic.LoadInt32(&self.mining) == 1 {
		self.unconfirmed.Shift(work.Block.NumberU64() - 1)
	}

	//这里需要打印日志
	/*log.Info("****current statistics****")
	log.Info("Header", "lastNumber", header.Number, "dposProto", header.DposProto.Root().String(), "bokerProto", header.BokerProto.Root().String())
	log.Info("Miner", "coinbase", self.coinbase, "engine", "Dpos")
	log.Info("Validator", "validator", header.Validator)
	tokenNoder, tokenErr := work.Block.DposContext.GetTokenNoder(header.Time.Int64())
	if tokenErr == nil {
		log.Info("Token Noder", "assign", tokenNoder)
	}*/
	return work, nil
}

//叔块header的hash添加进Work.uncles集合中
func (self *worker) commitUncle(work *Work, uncle *types.Header) error {

	hash := uncle.Hash()
	if work.uncles.Has(hash) {
		return fmt.Errorf("uncle not unique")
	}
	if !work.ancestors.Has(uncle.ParentHash) {
		return fmt.Errorf("uncle's parent unknown (%x)", uncle.ParentHash[0:4])
	}
	if work.family.Has(hash) {
		return fmt.Errorf("uncle already in family (%x)", hash)
	}
	work.uncles.Add(uncle.Hash())
	return nil
}

func (env *Work) commitTransactions(mux *event.TypeMux, txs *types.TransactionsByPriceAndNonce, bc *core.BlockChain, coinbase common.Address) {

	//给GasPrice一个初值GasLimit
	gp := new(core.GasPool).AddGas(env.header.GasLimit)

	var coalescedLogs []*types.Log
	for {
		//获取待处理交易池中一个交易,如果为空则退出
		tx := txs.Peek()
		if tx == nil {
			break
		}

		from, _ := types.Sender(env.signer, tx)
		if tx.Protected() && !env.config.IsEIP155(env.header.Number) {
			log.Trace("Ignoring reply protected transaction", "hash", tx.Hash(), "eip155", env.config.EIP155Block)
			txs.Pop()
			continue
		}

		//开始执行交易
		env.state.Prepare(tx.Hash(), common.Hash{}, env.tcount)
		err, logs := env.commitTransaction(tx, bc, coinbase, gp)

		//根据返回错误进行处理
		switch err {
		case core.ErrGasLimitReached: //当前区块超过了Gas限制
			log.Trace("Gas limit exceeded for current block", "sender", from)
			txs.Pop() //所有后续交易都将被跳过

		case core.ErrNonceTooLow: //Nonce太低
			log.Trace("Skipping transaction with low nonce", "sender", from, "nonce", tx.Nonce())
			txs.Shift() //移动到下一笔交易

		case core.ErrNonceTooHigh: //Nonce太高
			log.Trace("Skipping account with hight nonce", "sender", from, "nonce", tx.Nonce())
			txs.Pop()

		case nil:
			coalescedLogs = append(coalescedLogs, logs...)
			env.tcount++
			txs.Shift()

		default:
			// Strange error, discard the transaction and get the next in line (note, the
			// nonce-too-high clause will prevent us from executing in vain).
			log.Debug("Transaction failed, account skipped", "hash", tx.Hash(), "err", err)
			txs.Shift()
		}
	}

	if len(coalescedLogs) > 0 || env.tcount > 0 {
		// make a copy, the state caches the logs and these logs get "upgraded" from pending to mined
		// logs by filling in the block hash when the block was mined by the local miner. This can
		// cause a race condition if a log was "upgraded" before the PendingLogsEvent is processed.
		cpy := make([]*types.Log, len(coalescedLogs))
		for i, l := range coalescedLogs {
			cpy[i] = new(types.Log)
			*cpy[i] = *l
		}
		go func(logs []*types.Log, tcount int) {
			if len(logs) > 0 {
				mux.Post(core.PendingLogsEvent{Logs: logs})
			}
			if tcount > 0 {
				mux.Post(core.PendingStateEvent{})
			}
		}(cpy, env.tcount)
	}
}

func (env *Work) commitTransaction(tx *types.Transaction, bc *core.BlockChain, coinbase common.Address, gp *core.GasPool) (error, []*types.Log) {

	//首先获取当前状态的快照
	snap := env.state.Snapshot()
	dposSnap := env.dposContext.Snapshot()

	//执行交易
	receipt, _, err := core.ApplyTransaction(env.config,
		env.dposContext,
		bc,
		&coinbase,
		gp,
		env.state,
		env.header,
		tx,
		env.header.GasUsed,
		vm.Config{},
		bc.Boker())

	if err != nil {

		//交易执行失败，则回滚到之前的快照状态并返回错误，该账户的所有后续交易都将被跳过
		env.state.RevertToSnapshot(snap)
		env.dposContext.RevertToSnapShot(dposSnap)
		return err, nil
	}

	//交易执行成功，则记录该交易以及交易执行的回执（receipt）并返回，然后移动到下一笔交易
	env.txs = append(env.txs, tx)
	env.receipts = append(env.receipts, receipt)

	return nil, receipt.Logs
}
