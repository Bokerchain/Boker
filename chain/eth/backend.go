package eth

import (
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/boker/chain/accounts"
	"github.com/boker/chain/boker/api"
	"github.com/boker/chain/boker/protocol"
	"github.com/boker/chain/common"
	"github.com/boker/chain/common/hexutil"
	"github.com/boker/chain/consensus"
	"github.com/boker/chain/consensus/dpos"
	"github.com/boker/chain/core"
	"github.com/boker/chain/core/bloombits"
	"github.com/boker/chain/core/types"
	"github.com/boker/chain/core/vm"
	"github.com/boker/chain/eth/downloader"
	"github.com/boker/chain/eth/filters"
	"github.com/boker/chain/eth/gasprice"
	"github.com/boker/chain/ethdb"
	"github.com/boker/chain/event"
	"github.com/boker/chain/internal/ethapi"
	"github.com/boker/chain/log"
	"github.com/boker/chain/miner"
	"github.com/boker/chain/node"
	"github.com/boker/chain/p2p"
	"github.com/boker/chain/params"
	"github.com/boker/chain/rlp"
	"github.com/boker/chain/rpc"
)

type LesServer interface {
	Start(srvr *p2p.Server)
	Stop()
	Protocols() []p2p.Protocol
	SetBloomBitsIndexer(bbIndexer *core.ChainIndexer)
}

//以太坊实现的全节点类
type Ethereum struct {
	config          *Config
	chainConfig     *params.ChainConfig            //配置信息
	shutdownChan    chan bool                      // Channel for shutting down the ethereum
	stopDbUpgrade   func() error                   // stop chain db sequential key upgrade
	txPool          *core.TxPool                   //交易池
	blockchain      *core.BlockChain               //链对象
	protocolManager *ProtocolManager               //网络协议管理
	lesServer       LesServer                      //轻量级客户端服务器
	chainDb         ethdb.Database                 //区块链数据库对象
	eventMux        *event.TypeMux                 //事件临界区
	engine          consensus.Engine               //共识引擎
	accountManager  *accounts.Manager              //账号管理
	bloomRequests   chan chan *bloombits.Retrieval //通道接收绽放数据检索请求
	bloomIndexer    *core.ChainIndexer             // Bloom indexer operating during block imports
	ApiBackend      *EthApiBackend                 //eth的API后台类
	miner           *miner.Miner                   //挖矿类
	gasPrice        *big.Int                       //Gas单价
	coinbase        common.Address                 //挖矿账号
	password        string                         //挖矿账号的密码
	selfValidator   common.Address                 //设置当前出块节点为挖矿节点
	networkId       uint64                         //网络ID
	netRPCService   *ethapi.PublicNetAPI           //网络Api接口
	lock            sync.RWMutex                   // Protects the variadic fields (e.g. gas price and coinbase)
	boker           bokerapi.Api                   //播客链新增加的接口
}

func (s *Ethereum) AddLesServer(ls LesServer) {
	s.lesServer = ls
	ls.SetBloomBitsIndexer(s.bloomIndexer)
}

//创建实例对象
func New(ctx *node.ServiceContext, config *Config) (*Ethereum, error) {

	if config.SyncMode == downloader.LightSync {
		return nil, errors.New("can't run eth.Ethereum in light sync mode, use les.LightEthereum")
	}
	if !config.SyncMode.IsValid() {
		return nil, fmt.Errorf("invalid sync mode %d", config.SyncMode)
	}
	chainDb, err := CreateDB(ctx, config, "chaindata")
	if err != nil {
		return nil, err
	}
	stopDbUpgrade := upgradeDeduplicateData(chainDb)

	//得到配置信息
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlock(chainDb, config.Genesis)
	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {
		return nil, genesisErr
	}

	eth := &Ethereum{
		config:         config,
		chainDb:        chainDb,
		chainConfig:    chainConfig,
		eventMux:       ctx.EventMux,
		accountManager: ctx.AccountManager,
		engine:         dpos.New(&params.DposConfig{}, chainDb),
		shutdownChan:   make(chan bool),
		stopDbUpgrade:  stopDbUpgrade,
		networkId:      config.NetworkId,
		gasPrice:       config.GasPrice,
		coinbase:       config.Coinbase,
		bloomRequests:  make(chan chan *bloombits.Retrieval),
		bloomIndexer:   NewBloomIndexer(chainDb, params.BloomBitsBlocks),
	}

	if !config.SkipBcVersionCheck {
		bcVersion := core.GetBlockChainVersion(chainDb)
		if bcVersion != core.BlockChainVersion && bcVersion != 0 {
			return nil, fmt.Errorf("Blockchain DB version mismatch (%d / %d). Run geth upgradedb.\n", bcVersion, core.BlockChainVersion)
		}
		core.WriteBlockChainVersion(chainDb, core.BlockChainVersion)
	}

	vmConfig := vm.Config{EnablePreimageRecording: config.EnablePreimageRecording}
	eth.blockchain, err = core.NewBlockChain(chainDb, eth.chainConfig, eth.engine, vmConfig)
	if err != nil {
		return nil, err
	}

	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		eth.blockchain.SetHead(compat.RewindTo)
		core.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}
	eth.bloomIndexer.Start(eth.blockchain)

	if config.TxPool.Journal != "" {
		config.TxPool.Journal = ctx.ResolvePath(config.TxPool.Journal)
	}

	//新建交易池
	eth.txPool = core.NewTxPool(config.TxPool, eth.chainConfig, eth.blockchain)

	//新建协议管理器
	if eth.protocolManager, err = NewProtocolManager(eth.chainConfig,
		config.SyncMode,
		config.NetworkId,
		eth.eventMux,
		eth.txPool,
		eth.engine,
		eth.blockchain, chainDb); err != nil {
		return nil, err
	}

	//新建矿工
	eth.miner = miner.New(eth, eth.chainConfig, eth.EventMux(), eth.engine)
	eth.miner.SetExtra(makeExtraData(config.ExtraData))

	//新建后台
	eth.ApiBackend = &EthApiBackend{eth, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.GasPrice
	}
	eth.ApiBackend.gpo = gasprice.NewOracle(eth.ApiBackend, gpoParams)

	return eth, nil
}

//设置扩展数据内容
func makeExtraData(extra []byte) []byte {

	//如果扩展数据长度为0，则使用默认扩展数据
	if len(extra) == 0 {
		// create default extradata
		extra, _ = rlp.EncodeToBytes([]interface{}{
			uint(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch),
			"geth",
			runtime.Version(),
			runtime.GOOS,
		})
	}
	//如果扩展数据长度大于最大的扩展数据长度（32），则设置扩展数据为nil
	if uint64(len(extra)) > params.MaximumExtraDataSize {
		log.Warn("Miner extra data exceed limit", "extra", hexutil.Bytes(extra), "limit", params.MaximumExtraDataSize)
		extra = nil
	}
	return extra
}

//创建链DB
func CreateDB(ctx *node.ServiceContext, config *Config, name string) (ethdb.Database, error) {
	db, err := ctx.OpenDatabase(name, config.DatabaseCache, config.DatabaseHandles)
	if err != nil {
		return nil, err
	}
	if db, ok := db.(*ethdb.LDBDatabase); ok {
		db.Meter("eth/db/chaindata/")
	}
	return db, nil
}

//返回以太坊提供的RPC调用
func (s *Ethereum) APIs() []rpc.API {

	//获取提供出去的API数组
	apis := ethapi.GetAPIs(s.ApiBackend, s.boker)

	//添加共识引擎支持的Api接口到apis中
	apis = append(apis, s.engine.APIs(s.BlockChain())...)

	//添加所有本地Api结构到apis中
	apis = append(apis, []rpc.API{
		{
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicEthereumAPI(s),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   NewPublicMinerAPI(s),
			Public:    true,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.protocolManager.downloader, s.eventMux, s.boker),
			Public:    true,
		}, {
			Namespace: "miner",
			Version:   "1.0",
			Service:   NewPrivateMinerAPI(s),
			Public:    false,
		}, {
			Namespace: "eth",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.ApiBackend, false, s.boker),
			Public:    true,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateAdminAPI(s),
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(s),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(s.chainConfig, s),
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		},
	}...)

	return apis
}

func (s *Ethereum) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

//得到当前的挖矿账号
func (s *Ethereum) Coinbase() (common.Address, error) {
	s.lock.RLock()
	coinbase := s.coinbase
	s.lock.RUnlock()

	if coinbase != (common.Address{}) {
		return coinbase, nil
	}

	//如果当前的挖矿账号为空，则得到当前节点钱包中的第一个账号作为挖矿账号
	if wallets := s.AccountManager().Wallets(); len(wallets) > 0 {
		if accounts := wallets[0].Accounts(); len(accounts) > 0 {
			return accounts[0].Address, nil
		}
	}
	return common.Address{}, fmt.Errorf("coinbase address must be explicitly specified")
}

//设置挖矿账号
func (self *Ethereum) SetCoinbase(coinbase common.Address) {
	self.lock.Lock()
	self.coinbase = coinbase
	self.lock.Unlock()

	self.miner.SetCoinbase(coinbase)
}

//设置当前本地的验证者
func (self *Ethereum) SetLocalValidator(validator common.Address) error {

	self.lock.Lock()

	engine, ok := self.engine.(*dpos.Dpos)
	if !ok {
		log.Error("Only the dpos engine was allowed")
		return protocol.ErrInvalidSystem
	}

	//获取当前出块节点的数量
	size, sizeErr := engine.GetProducerSize(self.BlockChain().CurrentBlock(), validator)
	if sizeErr != nil {
		return sizeErr
	}
	if size == 0 {
		//判断当前区块序号为0
		if self.BlockChain().CurrentBlock().Number().Int64() == 0 {

			//设置当前出块节点为当前节点
			self.selfValidator = validator

			/*producerErr := engine.SelfProducer(self.BlockChain().CurrentBlock(), validator)
			if producerErr != nil {
				return producerErr
			}*/
		} else {
			log.Error("current Block Number is`t 0", "number", self.BlockChain().CurrentBlock().Number())
			return protocol.ErrGenesisBlock
		}
	} else {
		log.Error("current Prodcuder Size is`t 0")
		return protocol.ErrExistsValidators
	}

	self.lock.Unlock()

	return nil
}

//启动挖矿
func (s *Ethereum) StartMining(local bool) error {

	//得到当前的coinbase，并检测当前coinbase是否为nil
	coinbase, err := s.Coinbase()
	if err != nil {
		log.Error("Cannot start mining without coinbase", "err", err)
		return fmt.Errorf("coinbase missing: %v", err)
	}

	//根据当前的挖矿账号得到Dpos使用的签名函数
	if dpos, ok := s.engine.(*dpos.Dpos); ok {
		wallet, err := s.accountManager.Find(accounts.Account{Address: coinbase})
		if wallet == nil || err != nil {
			log.Error("Coinbase account unavailable locally", "err", err)
			return fmt.Errorf("signer missing: %v", err)
		}
		dpos.Authorize(coinbase, wallet.SignHash)
	}

	if local {
		// If local (CPU) mining is started, we can disable the transaction rejection
		// mechanism introduced to speed sync times. CPU mining on mainnet is ludicrous
		// so noone will ever hit this path, whereas marking sync done on CPU mining
		// will ensure that private networks work in single miner mode too.
		atomic.StoreUint32(&s.protocolManager.acceptTxs, 1)
	}
	go s.miner.Start(coinbase)
	return nil
}

func (s *Ethereum) Boker() bokerapi.Api {

	return s.boker
}

func (s *Ethereum) SetBoker(boker bokerapi.Api) {

	s.boker = boker
	s.blockchain.SetBoker(boker)
}

func (s *Ethereum) Password() string {

	return s.password
}

func (s *Ethereum) SetPassword(password string) {

	s.lock.Lock()
	s.password = password
	s.lock.Unlock()
}

//解码
func (s *Ethereum) DecodeParams(code []byte) ([]byte, error) {

	//
	return nil, nil
}

func (s *Ethereum) StopMining()                        { s.miner.Stop() }
func (s *Ethereum) IsMining() bool                     { return s.miner.Mining() }
func (s *Ethereum) Miner() *miner.Miner                { return s.miner }
func (s *Ethereum) AccountManager() *accounts.Manager  { return s.accountManager }
func (s *Ethereum) BlockChain() *core.BlockChain       { return s.blockchain }
func (s *Ethereum) TxPool() *core.TxPool               { return s.txPool }
func (s *Ethereum) EventMux() *event.TypeMux           { return s.eventMux }
func (s *Ethereum) Engine() consensus.Engine           { return s.engine }
func (s *Ethereum) ChainDb() ethdb.Database            { return s.chainDb }
func (s *Ethereum) IsListening() bool                  { return true } //总是处于监听
func (s *Ethereum) EthVersion() int                    { return int(s.protocolManager.SubProtocols[0].Version) }
func (s *Ethereum) NetVersion() uint64                 { return s.networkId }
func (s *Ethereum) Downloader() *downloader.Downloader { return s.protocolManager.downloader }

//返回所有当前配置的网络协议
func (s *Ethereum) Protocols() []p2p.Protocol {
	if s.lesServer == nil {
		return s.protocolManager.SubProtocols
	}
	return append(s.protocolManager.SubProtocols, s.lesServer.Protocols()...)
}

// Start implements node.Service, starting all internal goroutines needed by the
// Ethereum protocol implementation.
func (s *Ethereum) Start(srvr *p2p.Server) error {

	// Start the bloom bits servicing goroutines
	s.startBloomHandlers()

	//启动RPC服务
	s.netRPCService = ethapi.NewPublicNetAPI(srvr, s.NetVersion())

	// Figure out a max peers count based on the server limits
	maxPeers := srvr.MaxPeers
	if s.config.LightServ > 0 {
		maxPeers -= s.config.LightPeers
		if maxPeers < srvr.MaxPeers/2 {
			maxPeers = srvr.MaxPeers / 2
		}
	}

	//启动P2P网络
	s.protocolManager.Start(maxPeers)
	if s.lesServer != nil {
		s.lesServer.Start(srvr)
	}
	return nil
}

func (s *Ethereum) Stop() error {

	if s.stopDbUpgrade != nil {
		s.stopDbUpgrade()
	}

	s.bloomIndexer.Close()
	s.blockchain.Stop()
	s.protocolManager.Stop()

	if s.lesServer != nil {
		s.lesServer.Stop()
	}

	s.txPool.Stop()
	s.miner.Stop()
	s.eventMux.Stop()
	s.chainDb.Close()
	close(s.shutdownChan)

	return nil
}
