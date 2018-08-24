package miner

import (
	"fmt"
	"sync/atomic"

	"github.com/boker/go-ethereum/accounts"
	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/consensus"
	"github.com/boker/go-ethereum/core"
	"github.com/boker/go-ethereum/core/state"
	"github.com/boker/go-ethereum/core/types"
	"github.com/boker/go-ethereum/eth/downloader"
	"github.com/boker/go-ethereum/ethdb"
	"github.com/boker/go-ethereum/event"
	"github.com/boker/go-ethereum/log"
	"github.com/boker/go-ethereum/params"
)

// Backend wraps all methods required for mining.
type Backend interface {
	AccountManager() *accounts.Manager
	BlockChain() *core.BlockChain
	TxPool() *core.TxPool
	ChainDb() ethdb.Database
}

//创建区块并搜索POW值
type Miner struct {
	mux         *event.TypeMux
	worker      *worker
	coinbase    common.Address
	mining      int32
	eth         Backend
	engine      consensus.Engine
	canStart    int32 // can start indicates whether we can start the mining operation
	shouldStart int32 // should start indicates whether we should start after sync
}

func New(eth Backend, config *params.ChainConfig, mux *event.TypeMux, engine consensus.Engine) *Miner {

	//创建一个矿工
	miner := &Miner{
		eth:      eth,
		mux:      mux,
		engine:   engine,
		worker:   newWorker(config, engine, common.Address{}, eth, mux),
		canStart: 1,
	}

	//启动一个协程，等待事件到来
	go miner.update()

	return miner
}

// update keeps track of the downloader events. Please be aware that this is a one shot type of update loop.
// It's entered once and as soon as `Done` or `Failed` has been broadcasted the events are unregistered and
// the loop is exited. This to prevent a major security vuln where external parties can DOS you with blocks
// and halt your mining operation for as long as the DOS continues.
func (self *Miner) update() {
	events := self.mux.Subscribe(downloader.StartEvent{}, downloader.DoneEvent{}, downloader.FailedEvent{})
out:
	for ev := range events.Chan() {
		switch ev.Data.(type) {
		case downloader.StartEvent:
			atomic.StoreInt32(&self.canStart, 0)
			if self.Mining() {
				self.Stop()
				atomic.StoreInt32(&self.shouldStart, 1)
				log.Info("Mining aborted due to sync")
			}
		case downloader.DoneEvent, downloader.FailedEvent:
			shouldStart := atomic.LoadInt32(&self.shouldStart) == 1

			atomic.StoreInt32(&self.canStart, 1)
			atomic.StoreInt32(&self.shouldStart, 0)
			if shouldStart {
				self.Start(self.coinbase)
			}
			// unsubscribe. we're only interested in this event once
			events.Unsubscribe()
			// stop immediately and ignore all further pending events
			break out
		}
	}
}

func (self *Miner) Start(coinbase common.Address) {

	atomic.StoreInt32(&self.shouldStart, 1)

	//设置当前挖矿的账号
	log.Info("Miner Start Set Coinbase Address", "Info", coinbase)
	self.worker.setCoinbase(coinbase)
	self.coinbase = coinbase
	
	//
	self.worker.config.Coinbase = coinbase
	
	if atomic.LoadInt32(&self.canStart) == 0 {
		log.Info("Network syncing, will start miner afterwards")
		return
	}
	atomic.StoreInt32(&self.mining, 1)

	//启动挖矿
	log.Info("Starting mining operation")
	self.worker.start()
}

func (self *Miner) Stop() {

	//停止挖矿
	self.worker.stop()
	atomic.StoreInt32(&self.mining, 0)
	atomic.StoreInt32(&self.shouldStart, 0)
}

func (self *Miner) Mining() bool {
	return atomic.LoadInt32(&self.mining) > 0
}

func (self *Miner) HashRate() int64 {
	return 0
}

func (self *Miner) SetExtra(extra []byte) error {
	if uint64(len(extra)) > params.MaximumExtraDataSize {
		return fmt.Errorf("Extra exceeds max length. %d > %v", len(extra), params.MaximumExtraDataSize)
	}
	self.worker.setExtra(extra)
	return nil
}

// Pending returns the currently pending block and associated state.
func (self *Miner) Pending() (*types.Block, *state.StateDB) {
	return self.worker.pending()
}

// PendingBlock returns the currently pending block.
//
// Note, to access both the pending block and the pending state
// simultaneously, please use Pending(), as the pending state can
// change between multiple method calls
func (self *Miner) PendingBlock() *types.Block {
	return self.worker.pendingBlock()
}

func (self *Miner) SetCoinbase(addr common.Address) {
	self.coinbase = addr
	self.worker.setCoinbase(addr)
}
