package core

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/Bokerchain/Boker/chain/boker/protocol"
	"github.com/Bokerchain/Boker/chain/common"
	"github.com/Bokerchain/Boker/chain/core/state"
	"github.com/Bokerchain/Boker/chain/core/types"
	"github.com/Bokerchain/Boker/chain/event"
	"github.com/Bokerchain/Boker/chain/log"
	"github.com/Bokerchain/Boker/chain/metrics"
	"github.com/Bokerchain/Boker/chain/params"
	"gopkg.in/karalabe/cookiejar.v2/collections/prque"
)

const (
	// chainHeadChanSize is the size of channel listening to ChainHeadEvent.
	chainHeadChanSize = 10
	// rmTxChanSize is the size of channel listening to RemovedTransactionEvent.
	rmTxChanSize = 10
)

var (
	ErrInvalidSender = errors.New("invalid sender")          //如果交易包含无效签名
	ErrNonceTooLow   = errors.New("nonce too low")           //Nonce太低
	ErrUnderpriced   = errors.New("transaction underpriced") //交易的Gas比交易池中配置的价格还低
	// ErrReplaceUnderpriced is returned if a transaction is attempted to be replaced
	// with a different one without the required price bump.
	ErrReplaceUnderpriced = errors.New("replacement transaction underpriced")
	ErrInsufficientFunds  = errors.New("insufficient funds for gas * price + value") //执行交易的总成本高于用户帐户的余额
	// ErrIntrinsicGas is returned if the transaction is specified to use less gas
	// than required to start the invocation.
	ErrIntrinsicGas = errors.New("intrinsic gas too low")
	ErrGasLimit     = errors.New("exceeds block gas limit") //交易要求的Gas超过限额，则返回当前块要求的最大允许Gas
	// ErrNegativeValue is a sanity error to ensure noone is able to specify a
	// transaction with a negative value.
	ErrNegativeValue = errors.New("negative value")
	// ErrOversizedData is returned if the input data of a transaction is greater
	// than some meaningful limit a user might use. This is not a consensus error
	// making the transaction invalid, rather a DOS protection.
	ErrOversizedData = errors.New("oversized data")           //超大数据
	ErrInvalidType   = errors.New("unknown transaction type") //未知交易类型
)

var (
	evictionInterval    = time.Minute     //检查可撤销交易的时间间隔
	statsReportInterval = 8 * time.Second //报告交易池统计信息的时间间隔
)

var (
	//待处理池的度量标准
	pendingDiscardCounter   = metrics.NewCounter("txpool/pending/discard")
	pendingReplaceCounter   = metrics.NewCounter("txpool/pending/replace")
	pendingRateLimitCounter = metrics.NewCounter("txpool/pending/ratelimit") // Dropped due to rate limiting
	pendingNofundsCounter   = metrics.NewCounter("txpool/pending/nofunds")   // Dropped due to out-of-funds

	//排队池的度量标准
	queuedDiscardCounter   = metrics.NewCounter("txpool/queued/discard")
	queuedReplaceCounter   = metrics.NewCounter("txpool/queued/replace")
	queuedRateLimitCounter = metrics.NewCounter("txpool/queued/ratelimit") // Dropped due to rate limiting
	queuedNofundsCounter   = metrics.NewCounter("txpool/queued/nofunds")   // Dropped due to out-of-funds

	// General tx metrics
	invalidTxCounter     = metrics.NewCounter("txpool/invalid")
	underpricedTxCounter = metrics.NewCounter("txpool/underpriced")
)

// TxStatus is the current status of a transaction as seen py the pool.
type TxStatus uint

const (
	TxStatusUnknown TxStatus = iota
	TxStatusQueued
	TxStatusPending
	TxStatusIncluded
)

// blockChain provides the state of blockchain and current gas limit to do
// some pre checks in tx pool and event subscribers.
type blockChain interface {
	CurrentBlock() *types.Block
	GetBlock(hash common.Hash, number uint64) *types.Block
	StateAt(root common.Hash) (*state.StateDB, error)
	SubscribeChainHeadEvent(ch chan<- ChainHeadEvent) event.Subscription
}

//miner是从pending中拿交易组装block的
type TxPoolConfig struct {
	NoLocals     bool          //Whether local transaction handling should be disabled
	Journal      string        //Journal of local transactions to survive node restarts
	Rejournal    time.Duration //重新生成本地交易日志的时间间隔
	PriceLimit   uint64        //最小的GasPrice Minimum gas price to enforce for acceptance into the pool
	PriceBump    uint64        //最小的Price Minimum price bump percentage to replace an already existing transaction (nonce)
	AccountSlots uint64        //Minimum number of executable transaction slots guaranteed per account
	GlobalSlots  uint64        //所有帐户的最大可执行交易槽数（即所有账户的最大交易数量）
	AccountQueue uint64        //Maximum number of non-executable transaction slots permitted per account
	GlobalQueue  uint64        //所有帐户的最大不可执行交易槽数（即所有账户的最大未执行交易数量）
	Lifetime     time.Duration //Maximum amount of time non-executable transaction are queued
}

//交易池的默认配置
var DefaultTxPoolConfig = TxPoolConfig{
	Journal:      "transactions.rlp",
	Rejournal:    time.Hour,
	PriceLimit:   1,
	PriceBump:    10,
	AccountSlots: 16,            //一个账户所能放的默认交易数量
	GlobalSlots:  4096,          //交易池默认最大数量
	AccountQueue: 64,            //账号队列长度
	GlobalQueue:  1024,          //总体队列长度
	Lifetime:     3 * time.Hour, //3小时
}

//检查提供的用户配置,并更改任何不合理或不可行的配置
func (config *TxPoolConfig) sanitize() TxPoolConfig {

	conf := *config

	//重置无效的txpool日志时间
	if conf.Rejournal < time.Second {
		log.Warn("Sanitizing invalid txpool journal time", "provided", conf.Rejournal, "updated", time.Second)
		conf.Rejournal = time.Second
	}

	//重置无效的txpool价格限制
	if conf.PriceLimit < 1 {
		log.Warn("Sanitizing invalid txpool price limit", "provided", conf.PriceLimit, "updated", DefaultTxPoolConfig.PriceLimit)
		conf.PriceLimit = DefaultTxPoolConfig.PriceLimit
	}

	//重置无效的txpool价格暴跌
	if conf.PriceBump < 1 {
		log.Warn("Sanitizing invalid txpool price bump", "provided", conf.PriceBump, "updated", DefaultTxPoolConfig.PriceBump)
		conf.PriceBump = DefaultTxPoolConfig.PriceBump
	}
	return conf
}

// TxPool包含所有当前已知的交易。交易从网络收到或提交时进入池本地 当它们被包含在区块链中时，它们会退出交易池。
//交易池分隔可处理的交易（可以应用于当前状态）和未来的交易。 交易在这些之间移动随着时间的推移，它们会被接收和处理。
type TxPool struct {
	config        TxPoolConfig                       //交易池配置
	chainconfig   *params.ChainConfig                //链配置
	chain         blockChain                         //链
	gasPrice      *big.Int                           //最低的GasPrice限制
	txFeed        event.Feed                         //通过txFeed来订阅TxPool的消息
	scope         event.SubscriptionScope            //
	chainHeadCh   chan ChainHeadEvent                //订阅了区块头的消息，当有了新的区块头生成的时候会在这里收到通知
	chainHeadSub  event.Subscription                 //区块头消息的订阅器
	signer        types.Signer                       //封装了交易签名处理
	mu            sync.RWMutex                       //
	currentState  *state.StateDB                     //区块链头部当前状态
	pendingState  *state.ManagedState                //Pending state tracking virtual nonces
	currentMaxGas *big.Int                           //当前的交易Gas上限
	locals        *accountSet                        //Set of local transaction to exepmt from evicion rules
	journal       *txJournal                         //日志本地交易备份到磁盘
	pending       map[common.Address]*txList         //所有当前可处理的交易
	queue         map[common.Address]*txList         //不可处理的交易队列
	beats         map[common.Address]time.Time       //每个已知帐户的最后心跳
	all           map[common.Hash]*types.Transaction //允许查看的所有交易
	priced        *txPricedList                      //按价格排序的所有交易
	wg            sync.WaitGroup                     //for shutdown sync
	homestead     bool
}

//创建一个新的交易池，排序和过滤入站来自网络的交易
func NewTxPool(config TxPoolConfig, chainconfig *params.ChainConfig, chain blockChain) *TxPool {

	// Sanitize the input to ensure no vulnerable gas prices are set
	config = (&config).sanitize()

	//创建交易池并进行初始化
	pool := &TxPool{
		config:      config,
		chainconfig: chainconfig,
		chain:       chain,
		signer:      types.HomesteadSigner{},
		pending:     make(map[common.Address]*txList),
		queue:       make(map[common.Address]*txList),
		beats:       make(map[common.Address]time.Time),
		all:         make(map[common.Hash]*types.Transaction),
		chainHeadCh: make(chan ChainHeadEvent, chainHeadChanSize),
		gasPrice:    new(big.Int).SetUint64(config.PriceLimit),
	}

	pool.locals = newAccountSet(types.HomesteadSigner{})
	pool.priced = newTxPricedList(&pool.all)
	pool.reset(nil, chain.CurrentBlock().Header())

	//如果本地交易被允许,而且配置的Journal目录不为空,那么从指定的目录加载日志.
	//然后rotate交易日志. 因为老的交易可能已经失效了, 所以调用add方法之后再把被接收的交易写入日志.
	if !config.NoLocals && config.Journal != "" {
		pool.journal = newTxJournal(config.Journal)

		if err := pool.journal.load(pool.AddLocal); err != nil {
			log.Warn("Failed to load transaction journal", "err", err)
		}
		if err := pool.journal.rotate(pool.local()); err != nil {
			log.Warn("Failed to rotate transaction journal", "err", err)
		}
	}

	//从区块链订阅事件
	pool.chainHeadSub = pool.chain.SubscribeChainHeadEvent(pool.chainHeadCh)

	//启动交易池检测循环
	pool.wg.Add(1)
	go pool.loop()

	return pool
}

//启动交易池检测循环
func (pool *TxPool) loop() {
	defer pool.wg.Done()

	//启动定时器
	var prevPending, prevQueued, prevStales int

	//启动交易报告定时器
	report := time.NewTicker(statsReportInterval)
	defer report.Stop()

	//启动交易撤销定时器
	evict := time.NewTicker(evictionInterval)
	defer evict.Stop()

	//启动生成本地事务日志的定时器
	journal := time.NewTicker(pool.config.Rejournal)
	defer journal.Stop()

	// Track the previous head headers for transaction reorgs
	head := pool.chain.CurrentBlock()

	//等待并响应各种事件
	for {
		select {
		//监听到区块头的事件, 获取到新的区块头. 调用reset方法
		case ev := <-pool.chainHeadCh:
			if ev.Block != nil {
				pool.mu.Lock()
				//if pool.chainconfig.IsHomestead(ev.Block.Number()) {
				pool.homestead = true
				//}
				pool.reset(head.Header(), ev.Block.Header())
				head = ev.Block

				pool.mu.Unlock()
			}
		//由于系统停止而取消订阅
		case <-pool.chainHeadSub.Err():
			return

		//报告就是打印了一些日志
		case <-report.C:
			pool.mu.RLock()
			pending, queued := pool.stats()
			stales := pool.priced.stales
			pool.mu.RUnlock()

			if pending != prevPending || queued != prevQueued || stales != prevStales {
				log.Debug("Transaction pool status report", "executable", pending, "queued", queued, "stales", stales)
				prevPending, prevQueued, prevStales = pending, queued, stales
			}

		//处理超时的交易信息,
		case <-evict.C:
			pool.mu.Lock()
			for addr := range pool.queue {

				// Skip local transactions from the eviction mechanism
				if pool.locals.contains(addr) {
					continue
				}

				// Any non-locals old enough should be removed
				if time.Since(pool.beats[addr]) > pool.config.Lifetime {
					for _, tx := range pool.queue[addr].Flatten() {
						pool.removeTx(tx.Hash())
					}
				}
			}
			pool.mu.Unlock()

		//处理定时写交易日志的信息
		case <-journal.C:
			if pool.journal != nil {
				pool.mu.Lock()
				if err := pool.journal.rotate(pool.local()); err != nil {
					log.Warn("Failed to rotate local tx journal", "err", err)
				}
				pool.mu.Unlock()
			}
		}
	}
}

// lockedReset is a wrapper around reset to allow calling it in a thread safe
// manner. This method is only ever used in the tester!
func (pool *TxPool) lockedReset(oldHead, newHead *types.Header) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	pool.reset(oldHead, newHead)
}

//reset方法检索区块链的当前状态并且确保交易池的内容关于当前的区块链状态是有效的。主要功能包括：
//因为更换了区块头，所以原有的区块中有一些交易因为区块头的更换而作废，这部分交易需要重新加入到txPool里面等待插入新的区块
//生成新的currentState和pendingState
//因为状态的改变。将pending中的部分交易移到queue里面
//因为状态的改变，将queue里面的交易移入到pending里面。
func (pool *TxPool) reset(oldHead, newHead *types.Header) {

	//log.Info("(pool *TxPool) reset")

	//如果我们要重新安排旧状态，请重新注入所有丢弃的交易
	var reinject types.Transactions
	if oldHead != nil && oldHead.Hash() != newHead.ParentHash {

		// If the reorg is too deep, avoid doing it (will happen during fast sync)
		oldNum := oldHead.Number.Uint64()
		newNum := newHead.Number.Uint64()

		if depth := uint64(math.Abs(float64(oldNum) - float64(newNum))); depth > 64 {
			log.Warn("Skipping deep transaction reorg", "depth", depth)
		} else {
			// Reorg seems shallow enough to pull in all transactions into memory
			var discarded, included types.Transactions

			var (
				rem = pool.chain.GetBlock(oldHead.Hash(), oldHead.Number.Uint64())
				add = pool.chain.GetBlock(newHead.Hash(), newHead.Number.Uint64())
			)

			//如果老的高度大于新的.那么需要把多的全部删除.
			for rem.NumberU64() > add.NumberU64() {
				discarded = append(discarded, rem.Transactions()...)
				if rem = pool.chain.GetBlock(rem.ParentHash(), rem.NumberU64()-1); rem == nil {
					log.Error("Unrooted old chain seen by tx pool", "block", oldHead.Number, "hash", oldHead.Hash())
					return
				}
			}

			//如果新的高度大于老的, 那么需要增加.
			for add.NumberU64() > rem.NumberU64() {
				included = append(included, add.Transactions()...)
				if add = pool.chain.GetBlock(add.ParentHash(), add.NumberU64()-1); add == nil {
					log.Error("Unrooted new chain seen by tx pool", "block", newHead.Number, "hash", newHead.Hash())
					return
				}
			}

			//高度相同了.如果hash不同,那么需要往后找,一直找到他们相同hash根的节点.
			for rem.Hash() != add.Hash() {
				discarded = append(discarded, rem.Transactions()...)
				if rem = pool.chain.GetBlock(rem.ParentHash(), rem.NumberU64()-1); rem == nil {
					log.Error("Unrooted old chain seen by tx pool", "block", oldHead.Number, "hash", oldHead.Hash())
					return
				}
				included = append(included, add.Transactions()...)
				if add = pool.chain.GetBlock(add.ParentHash(), add.NumberU64()-1); add == nil {
					log.Error("Unrooted new chain seen by tx pool", "block", newHead.Number, "hash", newHead.Hash())
					return
				}
			}

			//找出所有存在discard里面,但是不在included里面的值.
			//需要等下把这些交易重新插入到pool里面。
			reinject = types.TxDifference(discarded, included)
		}
	}

	//将内部状态初始化为当前头部
	if newHead == nil {
		newHead = pool.chain.CurrentBlock().Header() // Special case during testing
	}
	//log.Info("(pool *TxPool) reset", "newHead", newHead.Number)

	statedb, err := pool.chain.StateAt(newHead.Root)
	if err != nil {
		log.Error("Failed to reset txpool state", "err", err)
		return
	}
	pool.currentState = statedb
	pool.pendingState = state.ManageState(statedb)
	pool.currentMaxGas = newHead.GasLimit
	pool.addTxsLocked(reinject, false)

	//验证pending transaction池里面的交易， 会移除所有已经存在区块链里面的交易，或者是因为其他交易导致不可用的交易(比如有一个更高的gasPrice)
	//demote 降级 将pending中的一些交易降级到queue里面。
	pool.demoteUnexecutables()

	// 根据pending队列的nonce更新所有账号的nonce
	for addr, list := range pool.pending {
		txs := list.Flatten() // Heavy but will be cached and is needed by the miner anyway
		pool.pendingState.SetNonce(addr, txs[len(txs)-1].Nonce()+1)

		//log.Info("(pool *TxPool) reset", "addr", addr, "nonce", txs[len(txs)-1].Nonce()+1)
	}

	//检查队列并尽可能地将事务移到pending，或删除那些已经失效的事务
	//promote 升级
	pool.promoteExecutables(nil)
}

// Stop terminates the transaction pool.
func (pool *TxPool) Stop() {
	// Unsubscribe all subscriptions registered from txpool
	pool.scope.Close()

	// Unsubscribe subscriptions registered from blockchain
	pool.chainHeadSub.Unsubscribe()
	pool.wg.Wait()

	if pool.journal != nil {
		pool.journal.close()
	}
	log.Info("Transaction pool stopped")
}

// SubscribeTxPreEvent registers a subscription of TxPreEvent and
// starts sending event to the given channel.
func (pool *TxPool) SubscribeTxPreEvent(ch chan<- TxPreEvent) event.Subscription {
	return pool.scope.Track(pool.txFeed.Subscribe(ch))
}

//GasPrice返回交易池强制执行的当前Gas价格
func (pool *TxPool) GasPrice() *big.Int {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return new(big.Int).Set(pool.gasPrice)
}

//更新交易池所需的最低价格，并删除低于此阈值的所有交易
func (pool *TxPool) SetGasPrice(price *big.Int) {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	pool.gasPrice = price
	for _, tx := range pool.priced.Cap(price, pool.locals) {
		pool.removeTx(tx.Hash())
	}
	//log.Info("Transaction pool price threshold updated", "price", price)
}

//返回交易池的虚拟托管状态
func (pool *TxPool) State() *state.ManagedState {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return pool.pendingState
}

// Stats retrieves the current pool stats, namely the number of pending and the
// number of queued (non-executable) transactions.
func (pool *TxPool) Stats() (int, int) {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return pool.stats()
}

//stats检索当前交易池的统计信息，即pending池和queue池的交易数量.
func (pool *TxPool) stats() (int, int) {

	//log.Info("(pool *TxPool) stats")

	pending := 0
	for _, list := range pool.pending {
		pending += list.Len()
	}
	queued := 0
	for _, list := range pool.queue {
		queued += list.Len()
	}

	//log.Info("(pool *TxPool) stats", "pending", pending, "queued", queued)

	return pending, queued
}

//检索交易池的数据内容，返回所有内容挂起和排队的交易，按帐户分组并按nonce排序
func (pool *TxPool) Content() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {

	log.Info("(pool *TxPool) Content")

	pool.mu.Lock()
	defer pool.mu.Unlock()

	pending := make(map[common.Address]types.Transactions)
	for addr, list := range pool.pending {
		pending[addr] = list.Flatten()
	}
	log.Info("(pool *TxPool) Content", "pending", len(pending))

	queued := make(map[common.Address]types.Transactions)
	for addr, list := range pool.queue {
		queued[addr] = list.Flatten()
	}
	log.Info("(pool *TxPool) Content", "queued", len(queued))

	return pending, queued
}

//待定检索按来源分组的所有当前可处理的交易帐户并按nonce排序。 返回的交易集是一个副本，可以是通过调用代码自由修改。
func (pool *TxPool) Pending() (map[common.Address]types.Transactions, error) {

	//log.Info("(pool *TxPool) Pending")

	pool.mu.Lock()
	defer pool.mu.Unlock()

	pending := make(map[common.Address]types.Transactions)
	for addr, list := range pool.pending {
		pending[addr] = list.Flatten()
	}
	log.Info("(pool *TxPool) Pending", "pending", len(pending))

	return pending, nil
}

//检索按来源分组的所有当前已知的本地交易帐户并按nonce排序。 返回的交易集是一个副本，可以是通过调用代码自由修改。
func (pool *TxPool) local() map[common.Address]types.Transactions {

	txs := make(map[common.Address]types.Transactions)
	for addr := range pool.locals.accounts {
		if pending := pool.pending[addr]; pending != nil {
			txs[addr] = append(txs[addr], pending.Flatten()...)
		}
		if queued := pool.queue[addr]; queued != nil {
			txs[addr] = append(txs[addr], queued.Flatten()...)
		}
	}
	return txs
}

//普通交易检验
func (pool *TxPool) normalValidateTx(tx *types.Transaction, local bool) error {

	log.Info("(pool *TxPool) normalValidateTx",
		"nonce", tx.Nonce(),
		"gas limit", tx.Gas(),
		"gas", tx.GasPrice(),
		"value", tx.Value(),
		"pool.currentMaxGas", pool.currentMaxGas)

	//如果当前的最大Gas数量小于交易所标记的Gas数量，则放回GasLimit错误(这里需要添加针对基础合约类型的判断，因为基础合约采用的Gas为最大值)
	if pool.currentMaxGas.Cmp(tx.Gas()) < 0 {
		return ErrGasLimit
	}

	//判断交易是否已经经过正确的签名
	from, err := types.Sender(types.HomesteadSigner{}, tx)
	if err != nil {
		return ErrInvalidSender
	}
	log.Info("(pool *TxPool) normalValidateTx", "from", from)

	// Drop non-local transactions under our own minimal accepted gas price
	local = local || pool.locals.contains(from) // account may be local even if the transaction arrived from the network

	//如果不是本地的交易,并且GasPrice低于我们的设置,那么也不会接收
	if !local && pool.gasPrice.Cmp(tx.GasPrice()) > 0 {
		return ErrUnderpriced
	}

	//判断交易中的Nonce是否大于发出交易用户的当前Nonce值
	if pool.currentState.GetNonce(from) > tx.Nonce() {
		return ErrNonceTooLow
	}

	//cost == Value + GasPrice * GasLimit
	//判断当前from用户的钱是否大于本次交易所花成本的最大值，如果小于则返回 ErrInsufficientFunds
	if pool.currentState.GetBalance(from).Cmp(tx.Cost()) < 0 {
		log.Error("(pool *TxPool) normalValidateTx", "balance", pool.currentState.GetBalance(from), "cost", tx.Cost())
		return ErrInsufficientFunds
	}

	//使用给定的数据，计算Gas。
	intrGas := IntrinsicGas(tx.Data(), tx.To() == nil, pool.homestead)

	//判断交易的Gas是否小于计算出来的Gas，如果小于则返回 ErrIntrinsicGas
	if tx.Gas().Cmp(intrGas) < 0 {
		return ErrIntrinsicGas
	}
	return nil
}

//普通交易检验
func (pool *TxPool) baseValidateTx(tx *types.Transaction, local bool) error {

	//判断交易是否已经经过正确的签名
	//from, err := types.Sender(pool.signer, tx)
	from, err := types.Sender(types.HomesteadSigner{}, tx)
	if err != nil {
		return ErrInvalidSender
	}

	//判断交易中的Nonce是否大于发出交易用户的当前Nonce值
	if pool.currentState.GetNonce(from) > tx.Nonce() {
		return ErrNonceTooLow
	}
	return nil
}

//对交易进行基本信息的验证
func (pool *TxPool) validateTx(tx *types.Transaction, local bool) error {

	//Dos攻击判断
	if tx.Size() > 32*1024 {
		return ErrOversizedData
	}

	//交易值是否进行签名判断
	if tx.Value().Sign() < 0 {
		return ErrNegativeValue
	}

	if types.IsBinary(tx.Type()) {

		//普通交易类型
		return pool.normalValidateTx(tx, local)
	} else if (tx.Type() >= protocol.SetValidator) && (tx.Type() <= protocol.AssignToken) {

		//基础合约交易类型
		return pool.baseValidateTx(tx, local)
	} else {

		//未知的交易类型
		return ErrInvalidType
	}
}

//验证交易并将其插入到future queue. 如果这个交易是替换了当前存在的某个交易,那么会返回之前的那个交易,这样外部就不用调用promote方法.
//如果某个新增加的交易被标记为local, 那么它的发送账户会进入白名单,这个账户的关联的交易将不会因为价格的限制或者其他的一些限制被删除.
func (pool *TxPool) add(tx *types.Transaction, local bool) (bool, error) {

	log.Info("(pool *TxPool) add",
		"to", tx.To(),
		"gas", tx.Gas(),
		"gasprice", tx.GasPrice(),
		"hash", tx.Hash().String(),
		"local", local)

	//检测此交易的Hash是否存在，如果存在，则表明是已经存在的交易，将其丢弃
	hash := tx.Hash()
	if pool.all[hash] != nil {
		log.Error("TxPool add Discarding already known transaction", "hash", hash)
		return false, fmt.Errorf("known transaction: %x", hash)
	}

	//对交易进行基本的验证，如果验证失败，则将其丢弃
	if err := pool.validateTx(tx, local); err != nil {
		log.Error("TxPool add Discarding invalid transaction", "hash", hash, "err", err)
		invalidTxCounter.Inc(1)
		return false, err
	}
	//log.Info("(pool *TxPool) add validateTx")

	//判断当前的交易池是否已经处于满状态（可执行最大槽数 + 不可执行最大槽数）
	if uint64(len(pool.all)) >= pool.config.GlobalSlots+pool.config.GlobalQueue {

		//如果新交易本身比目前跟踪的最低交易还低.那么不接收它
		if pool.priced.Underpriced(tx, pool.locals) {
			log.Error("TxPool add Discarding underpriced transaction", "hash", hash, "price", tx.GasPrice())
			underpricedTxCounter.Inc(1)
			return false, ErrUnderpriced
		}

		//如果交易大于交易池中的最小交易价格，则从交易池中删除最小价格交易，给本交易腾出空间
		drop := pool.priced.Discard(len(pool.all)-int(pool.config.GlobalSlots+pool.config.GlobalQueue-1), pool.locals)
		for _, tx := range drop {
			log.Error("TxPool add Discarding freshly underpriced transaction", "hash", tx.Hash(), "price", tx.GasPrice())
			underpricedTxCounter.Inc(1)
			pool.removeTx(tx.Hash())
		}
	}
	//log.Info("(pool *TxPool) add GlobalQueue")

	//根据交易签名获取本次交易的from用户
	from, _ := types.Sender(types.HomesteadSigner{}, tx)

	//判断本次交易的Nonce值是否已经存在于此用户的交易列表中
	if list := pool.pending[from]; list != nil && list.Overlaps(tx) {

		//本次交易的Nonce已经存在，检查是否满足所需的价格冲击
		inserted, old := list.Add(tx, pool.config.PriceBump)
		if !inserted {
			log.Error("TxPool add nonce Exists", "nonce", tx.Nonce())
			pendingDiscardCounter.Inc(1)
			return false, ErrReplaceUnderpriced
		}

		//新交易更好，取代旧交易
		if old != nil {

			//从交易池中根据原有交易的Hash删除原来的交易
			delete(pool.all, old.Hash())
			pool.priced.Removed()
			pendingReplaceCounter.Inc(1)
		}

		//在交易池中添加本次交易
		pool.all[tx.Hash()] = tx
		pool.priced.Put(tx)
		pool.journalTx(from, tx)

		log.Info("Pooled new executable transaction", "hash", hash, "from", from, "to", tx.To())

		//向所有订阅的频道发送，返回订阅者的数量
		go pool.txFeed.Send(TxPreEvent{tx})

		return old != nil, nil
	}
	//log.Info("(pool *TxPool) add Sender")

	//新交易不能替换pending里面的任意一个交易,那么把他push到futuren 队列里面.
	replace, err := pool.enqueueTx(hash, tx)
	if err != nil {
		return false, err
	}
	//log.Info("(pool *TxPool) add enqueueTx")

	//如果是本地的交易,会被记录进入journalTx
	if local {
		pool.locals.add(from)
	}
	pool.journalTx(from, tx)

	//log.Info("(pool *TxPool) add", "hash", hash, "from", from, "to", tx.To())

	return replace, nil
}

//将新交易插入到非可执行交易队列中,注意! 此方法假定池锁已被保留！
func (pool *TxPool) enqueueTx(hash common.Hash, tx *types.Transaction) (bool, error) {

	//log.Info("(pool *TxPool) enqueueTx", "hash", hash)

	//尝试将交易插入到将来的队列中
	from, _ := types.Sender(types.HomesteadSigner{}, tx)
	if pool.queue[from] == nil {
		pool.queue[from] = newTxList(false)
	}

	inserted, old := pool.queue[from].Add(tx, pool.config.PriceBump)
	if !inserted {
		// An older transaction was better, discard this
		queuedDiscardCounter.Inc(1)
		return false, ErrReplaceUnderpriced
	}

	// Discard any previous transaction and mark this
	if old != nil {
		delete(pool.all, old.Hash())
		pool.priced.Removed()
		queuedReplaceCounter.Inc(1)
	}
	pool.all[hash] = tx
	pool.priced.Put(tx)
	return old != nil, nil
}

//将指定的交易添加到本地磁盘日志中（如果是）视为已从本地帐户发送.
func (pool *TxPool) journalTx(from common.Address, tx *types.Transaction) {

	//只有日志，如果它已启用且交易是本地的
	if pool.journal == nil || !pool.locals.contains(from) {
		return
	}
	if err := pool.journal.insert(tx); err != nil {
		log.Warn("Failed to journal local transaction", "err", err)
	}
}

//把某个交易加入到pending 队列. 这个方法假设已经获取到了锁
func (pool *TxPool) promoteTx(addr common.Address, hash common.Hash, tx *types.Transaction) {

	//尝试将交易插入挂起队列
	if pool.pending[addr] == nil {
		pool.pending[addr] = newTxList(true)
	}
	list := pool.pending[addr]

	inserted, old := list.Add(tx, pool.config.PriceBump)
	if !inserted {

		// 如果不能替换, 已经存在一个老的交易了. 删除.
		delete(pool.all, hash)
		pool.priced.Removed()

		pendingDiscardCounter.Inc(1)
		return
	}
	// Otherwise discard any previous transaction and mark this
	if old != nil {
		delete(pool.all, old.Hash())
		pool.priced.Removed()

		pendingReplaceCounter.Inc(1)
	}
	// Failsafe to work around direct pending inserts (tests)
	if pool.all[hash] == nil {
		pool.all[hash] = tx
		pool.priced.Put(tx)
	}

	// 把交易加入到队列,并发送消息告诉所有的订阅者, 这个订阅者在eth协议内部. 会接收这个消息并把这个消息通过网路广播出去.
	pool.beats[addr] = time.Now()
	pool.pendingState.SetNonce(addr, tx.Nonce()+1)

	go pool.txFeed.Send(TxPreEvent{tx})
}

//本地节点产生单条交易
func (pool *TxPool) AddLocal(tx *types.Transaction) error {

	//log.Info("(pool *TxPool) AddLocal", "Nonce", tx.Nonce())
	return pool.addTx(tx, !pool.config.NoLocals)
}

//网络中接收的单条交易
func (pool *TxPool) AddRemote(tx *types.Transaction) error {

	//log.Info("****AddRemote****", "Nonce", tx.Nonce())
	return pool.addTx(tx, false)
}

//本地节点产生一批交易
func (pool *TxPool) AddLocals(txs []*types.Transaction) []error {

	//log.Info("****AddLocals****", "len", len(txs))
	return pool.addTxs(txs, !pool.config.NoLocals)
}

//从网络中接收一批交易
func (pool *TxPool) AddRemotes(txs []*types.Transaction) []error {

	log.Info("(pool *TxPool) AddRemotes", "len", len(txs))
	return pool.addTxs(txs, false)
}

//将交易放入到交易池中
func (pool *TxPool) addTx(tx *types.Transaction, local bool) error {

	//log.Info("(pool *TxPool) addTx", "hash", tx.Hash())

	pool.mu.Lock()
	defer pool.mu.Unlock()

	//尝试注入交易并更新任何状态
	replace, err := pool.add(tx, local)
	if err != nil {
		return err
	}

	//如果我们添加了新的交易，请运行促销检查并返回
	if !replace {
		from, _ := types.Sender(types.HomesteadSigner{}, tx)
		pool.promoteExecutables([]common.Address{from})
	}
	return nil
}

// addTxs attempts to queue a batch of transactions if they are valid.
func (pool *TxPool) addTxs(txs []*types.Transaction, local bool) []error {
	pool.mu.Lock()
	defer pool.mu.Unlock()

	return pool.addTxsLocked(txs, local)
}

//尝试把有效的交易放入queue队列，调用这个函数的时候假设已经获取到锁
func (pool *TxPool) addTxsLocked(txs []*types.Transaction, local bool) []error {

	// Add the batch of transaction, tracking the accepted ones
	dirty := make(map[common.Address]struct{})
	errs := make([]error, len(txs))

	for i, tx := range txs {
		var replace bool
		if replace, errs[i] = pool.add(tx, local); errs[i] == nil {

			//replace 是替换的意思， 如果不是替换，那么就说明状态有更新，有可以下一步处理的可能。
			if !replace {
				//from, _ := types.Sender(pool.signer, tx) // already validated
				from, _ := types.Sender(types.HomesteadSigner{}, tx)
				dirty[from] = struct{}{}
			}
		}
	}

	// Only reprocess the internal state if something was actually added
	if len(dirty) > 0 {
		addrs := make([]common.Address, 0, len(dirty))
		for addr := range dirty {
			addrs = append(addrs, addr)
		}

		//传入了被修改的地址
		pool.promoteExecutables(addrs)
	}
	return errs
}

// Status returns the status (unknown/pending/queued) of a batch of transactions
// identified by their hashes.
func (pool *TxPool) Status(hashes []common.Hash) []TxStatus {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	status := make([]TxStatus, len(hashes))
	for i, hash := range hashes {
		if tx := pool.all[hash]; tx != nil {
			//from, _ := types.Sender(pool.signer, tx) // already validated
			from, _ := types.Sender(types.HomesteadSigner{}, tx)
			if pool.pending[from].txs.items[tx.Nonce()] != nil {
				status[i] = TxStatusPending
			} else {
				status[i] = TxStatusQueued
			}
		}
	}
	return status
}

//如果交易包含在池中，则返回返回交易，否则为空。
func (pool *TxPool) Get(hash common.Hash) *types.Transaction {
	pool.mu.RLock()
	defer pool.mu.RUnlock()

	return pool.all[hash]
}

//删除某个交易， 并把所有后续的交易移动到future queue
func (pool *TxPool) removeTx(hash common.Hash) {

	log.Info("(pool *TxPool) removeTx", "hash", hash)

	// Fetch the transaction we wish to delete
	tx, ok := pool.all[hash]
	if !ok {
		return
	}
	addr, _ := types.Sender(types.HomesteadSigner{}, tx)

	// Remove it from the list of known transactions
	delete(pool.all, hash)
	pool.priced.Removed()

	// 把交易从pending删除， 并把因为这个交易的删除而变得无效的交易放到future queue 然后更新pendingState的状态
	if pending := pool.pending[addr]; pending != nil {
		if removed, invalids := pending.Remove(tx); removed {
			// If no more transactions are left, remove the list
			if pending.Empty() {
				delete(pool.pending, addr)
				delete(pool.beats, addr)
			} else {
				// Otherwise postpone any invalidated transactions
				for _, tx := range invalids {
					pool.enqueueTx(tx.Hash(), tx)
				}
			}
			// Update the account nonce if needed
			if nonce := tx.Nonce(); pool.pendingState.GetNonce(addr) > nonce {
				pool.pendingState.SetNonce(addr, nonce)
			}
			return
		}
	}

	//把交易从future queue删除.
	if future := pool.queue[addr]; future != nil {
		future.Remove(tx)
		if future.Empty() {
			delete(pool.queue, addr)
		}
	}
}

//把已经变得可以执行的交易从future queue 插入到pending queue. 在这个过程中，所有删除无效的交易（低随机数，低余额）。
func (pool *TxPool) promoteExecutables(accounts []common.Address) {

	//log.Info("(pool *TxPool) promoteExecutables")

	// accounts存储了所有潜在需要更新的账户。 如果账户传入为nil，代表所有已知的账户。
	if accounts == nil {
		accounts = make([]common.Address, 0, len(pool.queue))
		for addr := range pool.queue {
			accounts = append(accounts, addr)
		}
	}

	for _, addr := range accounts {

		//log.Info("(pool *TxPool) promoteExecutables", "addr", addr)

		list := pool.queue[addr]
		if list == nil {
			continue
		}

		//删除所有的nonce太低的交易
		//log.Info("(pool *TxPool) promoteExecutables Forward")
		for _, tx := range list.Forward(pool.currentState.GetNonce(addr)) {
			hash := tx.Hash()
			log.Info("(pool *TxPool) promoteExecutables Removed old queued transaction", "hash", hash)
			delete(pool.all, hash)
			pool.priced.Removed()
		}

		//删除所有余额不足的交易。
		//log.Info("(pool *TxPool) promoteExecutables GetBalance")
		drops, _ := list.Filter(pool.currentState.GetBalance(addr), pool.currentMaxGas)
		for _, tx := range drops {
			hash := tx.Hash()
			log.Info("(pool *TxPool) promoteExecutables Removed unpayable queued transaction", "hash", hash)
			delete(pool.all, hash)
			pool.priced.Removed()
			queuedNofundsCounter.Inc(1)
		}

		// 显示一下目前此用户所有的交易信息；
		log.Info("(pool *TxPool) Check User Tx", "addr", addr, "current nonce", pool.pendingState.GetNonce(addr))
		for _, tx := range list.Check(pool.pendingState.GetNonce(addr)) {
			log.Info("(pool *TxPool) User Tx", "hash", tx.Hash(), "nonce", tx.Nonce())
		}

		//得到所有的可以执行的交易，并promoteTx加入pending
		log.Info("(pool *TxPool) Ready", "addr", addr)
		for _, tx := range list.Ready(pool.pendingState.GetNonce(addr)) {
			hash := tx.Hash()
			pool.promoteTx(addr, hash, tx)
		}

		//删除所有超过限制的交易
		if !pool.locals.contains(addr) {
			for _, tx := range list.Cap(int(pool.config.AccountQueue)) {
				hash := tx.Hash()
				delete(pool.all, hash)
				pool.priced.Removed()
				queuedRateLimitCounter.Inc(1)
			}
		}
		// Delete the entire queue entry if it became empty.
		if list.Empty() {
			delete(pool.queue, addr)
		}
	}
	// If the pending limit is overflown, start equalizing allowances
	pending := uint64(0)
	for _, list := range pool.pending {
		pending += uint64(list.Len())
	}

	//如果pending的总数超过系统的配置。
	if pending > pool.config.GlobalSlots {
		pendingBeforeCap := pending
		// Assemble a spam order to penalize large transactors first
		spammers := prque.New()
		for addr, list := range pool.pending {

			// 首先把所有大于AccountSlots最小值的账户记录下来， 会从这些账户里面剔除一些交易。
			// 注意spammers是一个优先级队列，也就是说是按照交易的多少从大到小排序的。
			if !pool.locals.contains(addr) && uint64(list.Len()) > pool.config.AccountSlots {
				spammers.Push(addr, float32(list.Len()))
			}
		}
		// Gradually drop transactions from offenders
		offenders := []common.Address{}
		for pending > pool.config.GlobalSlots && !spammers.Empty() {

			/*
				模拟一下offenders队列的账户交易数量的变化情况。
				第一次循环   [10]    循环结束  [10]
				第二次循环   [10, 9] 循环结束  [9,9]
				第三次循环   [9, 9, 7] 循环结束 [7, 7, 7]
				第四次循环   [7, 7 , 7 ,2] 循环结束 [2, 2 ,2, 2]
			*/

			// Retrieve the next offender if not local address
			offender, _ := spammers.Pop()
			offenders = append(offenders, offender.(common.Address))

			// Equalize balances until all the same or below threshold
			if len(offenders) > 1 { // 第一次进入这个循环的时候， offenders队列里面有交易数量最大的两个账户

				//把最后加入的账户的交易数量当成本次的阈值
				threshold := pool.pending[offender.(common.Address)].Len()

				//遍历直到pending有效，或者是倒数第二个的交易数量等于最后一个的交易数量
				for pending > pool.config.GlobalSlots && pool.pending[offenders[len(offenders)-2]].Len() > threshold {

					//遍历除了最后一个账户以外的所有账户， 把他们的交易数量减去1.
					for i := 0; i < len(offenders)-1; i++ {
						list := pool.pending[offenders[i]]
						for _, tx := range list.Cap(list.Len() - 1) {

							// Drop the transaction from the global pools too
							hash := tx.Hash()
							delete(pool.all, hash)
							pool.priced.Removed()

							// Update the account nonce to the dropped transaction
							if nonce := tx.Nonce(); pool.pendingState.GetNonce(offenders[i]) > nonce {
								pool.pendingState.SetNonce(offenders[i], nonce)
							}
							log.Trace("Removed fairness-exceeding pending transaction", "hash", hash)
						}
						pending--
					}
				}
			}
		}

		//经过上面的循环，所有的超过AccountSlots的账户的交易数量都变成了之前的最小值。
		//如果还是超过阈值，那么在继续从offenders里面每次删除一个。
		if pending > pool.config.GlobalSlots && len(offenders) > 0 {
			for pending > pool.config.GlobalSlots && uint64(pool.pending[offenders[len(offenders)-1]].Len()) > pool.config.AccountSlots {
				for _, addr := range offenders {
					list := pool.pending[addr]
					for _, tx := range list.Cap(list.Len() - 1) {
						// Drop the transaction from the global pools too
						hash := tx.Hash()
						delete(pool.all, hash)
						pool.priced.Removed()

						// Update the account nonce to the dropped transaction
						if nonce := tx.Nonce(); pool.pendingState.GetNonce(addr) > nonce {
							pool.pendingState.SetNonce(addr, nonce)
						}
						log.Trace("Removed fairness-exceeding pending transaction", "hash", hash)
					}
					pending--
				}
			}
		}
		pendingRateLimitCounter.Inc(int64(pendingBeforeCap - pending))
	}

	//我们处理了pending的限制， 下面需要处理future queue的限制了。
	queued := uint64(0)
	for _, list := range pool.queue {
		queued += uint64(list.Len())
	}
	if queued > pool.config.GlobalQueue {

		//按心跳排序所有具有排队交易的帐户
		addresses := make(addresssByHeartbeat, 0, len(pool.queue))
		for addr := range pool.queue {
			if !pool.locals.contains(addr) { // don't drop locals
				addresses = append(addresses, addressByHeartbeat{addr, pool.beats[addr]})
			}
		}
		sort.Sort(addresses)

		// Drop transactions until the total is below the limit or only locals remain
		for drop := queued - pool.config.GlobalQueue; drop > 0 && len(addresses) > 0; {
			addr := addresses[len(addresses)-1]
			list := pool.queue[addr.address]

			addresses = addresses[:len(addresses)-1]

			// Drop all transactions if they are less than the overflow
			if size := uint64(list.Len()); size <= drop {
				for _, tx := range list.Flatten() {
					pool.removeTx(tx.Hash())
				}
				drop -= size
				queuedRateLimitCounter.Inc(int64(size))
				continue
			}
			// Otherwise drop only last few transactions
			txs := list.Flatten()
			for i := len(txs) - 1; i >= 0 && drop > 0; i-- {
				pool.removeTx(txs[i].Hash())
				drop--
				queuedRateLimitCounter.Inc(1)
			}
		}
	}
}

// demoteUnexecutables removes invalid and processed transactions from the pools
// executable/pending queue and any subsequent transactions that become unexecutable
// are moved back into the future queue.
func (pool *TxPool) demoteUnexecutables() {

	// Iterate over all accounts and demote any non-executable transactions
	for addr, list := range pool.pending {
		nonce := pool.currentState.GetNonce(addr)

		// 删除所有小于当前地址的nonce的交易，并从pool.all删除。
		for _, tx := range list.Forward(nonce) {
			hash := tx.Hash()
			log.Trace("Removed old pending transaction", "hash", hash)
			delete(pool.all, hash)
			pool.priced.Removed()
		}

		// 删除所有的太昂贵的交易。 用户的balance可能不够用。或者是out of gas
		drops, invalids := list.Filter(pool.currentState.GetBalance(addr), pool.currentMaxGas)
		for _, tx := range drops {
			hash := tx.Hash()
			log.Trace("Removed unpayable pending transaction", "hash", hash)
			delete(pool.all, hash)
			pool.priced.Removed()
			pendingNofundsCounter.Inc(1)
		}
		for _, tx := range invalids {
			hash := tx.Hash()
			log.Trace("Demoting pending transaction", "hash", hash)
			pool.enqueueTx(hash, tx)
		}

		// 如果存在一个空洞(nonce空洞)， 那么需要把所有的交易都放入future queue。
		// 这一步确实应该不可能发生，因为Filter已经把 invalids的都处理了。 应该不存在invalids的交易，也就是不存在空洞的。
		if list.Len() > 0 && list.txs.Get(nonce) == nil {
			for _, tx := range list.Cap(0) {
				hash := tx.Hash()
				log.Error("Demoting invalidated transaction", "hash", hash)
				pool.enqueueTx(hash, tx)
			}
		}
		// Delete the entire queue entry if it became empty.
		if list.Empty() {
			delete(pool.pending, addr)
			delete(pool.beats, addr)
		}
	}
}

// addressByHeartbeat is an account address tagged with its last activity timestamp.
type addressByHeartbeat struct {
	address   common.Address
	heartbeat time.Time
}

type addresssByHeartbeat []addressByHeartbeat

func (a addresssByHeartbeat) Len() int           { return len(a) }
func (a addresssByHeartbeat) Less(i, j int) bool { return a[i].heartbeat.Before(a[j].heartbeat) }
func (a addresssByHeartbeat) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// accountSet is simply a set of addresses to check for existence, and a signer
// capable of deriving addresses from transactions.
type accountSet struct {
	accounts map[common.Address]struct{}
	signer   types.Signer
}

// newAccountSet creates a new address set with an associated signer for sender
// derivations.
func newAccountSet(signer types.Signer) *accountSet {
	return &accountSet{
		accounts: make(map[common.Address]struct{}),
		signer:   signer,
	}
}

// contains checks if a given address is contained within the set.
func (as *accountSet) contains(addr common.Address) bool {
	_, exist := as.accounts[addr]
	return exist
}

// containsTx checks if the sender of a given tx is within the set. If the sender
// cannot be derived, this method returns false.
func (as *accountSet) containsTx(tx *types.Transaction) bool {
	if addr, err := types.Sender(as.signer, tx); err == nil {
		return as.contains(addr)
	}
	return false
}

// add inserts a new address into the set to track.
func (as *accountSet) add(addr common.Address) {
	as.accounts[addr] = struct{}{}
}
