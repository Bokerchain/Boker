package main

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/boker/go-ethereum/accounts"
	"github.com/boker/go-ethereum/accounts/keystore"
	"github.com/boker/go-ethereum/cmd/utils"
	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/console"

	"github.com/boker/go-ethereum/accounts/abi/bind"
	"github.com/boker/go-ethereum/boker"
	//"github.com/boker/go-ethereum/contracts"
	"github.com/boker/go-ethereum/eth"
	"github.com/boker/go-ethereum/ethclient"
	"github.com/boker/go-ethereum/internal/debug"
	"github.com/boker/go-ethereum/log"
	"github.com/boker/go-ethereum/metrics"
	"github.com/boker/go-ethereum/node"
	"gopkg.in/urfave/cli.v1"
)

const (
	clientIdentifier = "geth" // Client identifier to advertise over the network
)

var (
	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit = ""

	//得到版本信息的合约地址.
	relOracle = common.HexToAddress("0xfa7b9770ca4cb04296cac84f37736d4041251cdf")

	// The app that holds all commands and flags.
	app = utils.NewApp(gitCommit, "the boker command line interface")

	//配置节点的标志
	nodeFlags = []cli.Flag{
		utils.IdentityFlag,
		utils.UnlockedAccountFlag,
		utils.PasswordFileFlag,
		utils.BootnodesFlag,
		utils.BootnodesV4Flag,
		utils.BootnodesV5Flag,
		utils.DataDirFlag,
		utils.KeyStoreDirFlag,
		utils.NoUSBFlag,
		utils.DashboardEnabledFlag,
		utils.DashboardAddrFlag,
		utils.DashboardPortFlag,
		utils.DashboardRefreshFlag,
		utils.TxPoolNoLocalsFlag,
		utils.TxPoolJournalFlag,
		utils.TxPoolRejournalFlag,
		utils.TxPoolPriceLimitFlag,
		utils.TxPoolPriceBumpFlag,
		utils.TxPoolAccountSlotsFlag,
		utils.TxPoolGlobalSlotsFlag,
		utils.TxPoolAccountQueueFlag,
		utils.TxPoolGlobalQueueFlag,
		utils.TxPoolLifetimeFlag,
		utils.FastSyncFlag,
		utils.LightModeFlag,
		utils.SyncModeFlag,
		utils.LightServFlag,
		utils.LightPeersFlag,
		utils.LightKDFFlag,
		utils.CacheFlag,
		utils.TrieCacheGenFlag,
		utils.ListenPortFlag,
		utils.MaxPeersFlag,
		utils.MaxPendingPeersFlag,
		utils.ValidatorFlag,
		utils.CoinbaseFlag,
		utils.GasPriceFlag,
		utils.MiningEnabledFlag,
		utils.TargetGasLimitFlag,
		utils.NATFlag,
		utils.NoDiscoverFlag,
		utils.DiscoveryV5Flag,
		utils.NetrestrictFlag,
		utils.NodeKeyFileFlag,
		utils.NodeKeyHexFlag,
		utils.VMEnableDebugFlag,
		utils.NetworkIdFlag,
		utils.RPCCORSDomainFlag,
		utils.EthStatsURLFlag,
		utils.MetricsEnabledFlag,
		utils.NoCompactionFlag,
		utils.GpoBlocksFlag,
		utils.GpoPercentileFlag,
		utils.ExtraDataFlag,
		configFileFlag,
	}

	rpcFlags = []cli.Flag{
		utils.RPCEnabledFlag,
		utils.RPCListenAddrFlag,
		utils.RPCPortFlag,
		utils.RPCApiFlag,
		utils.WSEnabledFlag,
		utils.WSListenAddrFlag,
		utils.WSPortFlag,
		utils.WSApiFlag,
		utils.WSAllowedOriginsFlag,
		utils.IPCDisabledFlag,
		utils.IPCPathFlag,
	}

	whisperFlags = []cli.Flag{
		utils.WhisperEnabledFlag,
		utils.WhisperMaxMessageSizeFlag,
		utils.WhisperMinPOWFlag,
	}
)

//初始化节点
func init() {

	//初始化CLI应用程序并启动Geth
	app.Action = geth

	//隐藏版本，通过命令可以查看版本信息
	app.HideVersion = true

	//作者信息
	app.Copyright = "Copyright 2017-2018 The Bokerchain Authors"

	//定义一个指令数组
	app.Commands = []cli.Command{

		//注册链Cmd指令，可以查看chaincmd.go
		initCommand,   //初始化指令
		importCommand, //从一个文件导入链
		exportCommand, //导出链到指定文件
		copydbCommand,
		removedbCommand,
		dumpCommand,

		//注册监控CMD指令，可以查看monitorcmd.go
		monitorCommand,

		//注册账号指令，可以查看accountcmd.go
		accountCommand,
		walletCommand,

		//注册控制台CMD指令，可以查看consolecmd.go
		consoleCommand,
		attachCommand,
		javascriptCommand,

		//注册版本以及BUG指令，可以查看misccmd.go
		versionCommand,
		bugCommand,
		licenseCommand,

		//注册调试config指令，可以查看config.go
		dumpConfigCommand,
	}
	sort.Sort(cli.CommandsByName(app.Commands))

	//
	app.Flags = append(app.Flags, nodeFlags...)
	app.Flags = append(app.Flags, rpcFlags...)
	app.Flags = append(app.Flags, consoleFlags...)
	app.Flags = append(app.Flags, debug.Flags...)
	app.Flags = append(app.Flags, whisperFlags...)

	app.Before = func(ctx *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		if err := debug.Setup(ctx); err != nil {
			return err
		}
		// Start system runtime metrics collection
		go metrics.CollectProcessMetrics(3 * time.Second)

		utils.SetupNetwork(ctx)
		return nil
	}

	app.After = func(ctx *cli.Context) error {
		debug.Exit()
		console.Stdin.Close() // Resets terminal mode.
		return nil
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

//如果没有运行特殊的子命令，geth是进入系统的主要入口点,它根据命令行参数创建一个默认节点并运行它阻塞模式，等待它关闭。
func geth(ctx *cli.Context) error {

	//生成一个*node.Node对象stack
	log.Info("Start geth With Context", "ctx", ctx)

	//这里的gethNode是一个全局变量
	bind.GethNode = makeFullNode(ctx)
	log.Info("Full Node Create:", "Node", bind.GethNode)

	//启动这个节点
	log.Info("Start Node Begin")
	startNode(ctx, bind.GethNode)
	log.Info("Start Node End")

	//

	//节点进入等待
	bind.GethNode.Wait()
	return nil
}

//启动系统节点和所有已注册的协议，之后它解锁任何请求的帐户，并启动RPC / IPC接口和矿工
/*func startNode(ctx *cli.Context, stack *node.Node) {

	//启动节点
	log.Info("startNode")
	utils.StartNode(stack)

	//解锁一些特定需求的账户
	log.Info("stack.AccountManager().Backends")
	ks := stack.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)

	//加载相应的密码信息
	passwords := utils.MakePasswordList(ctx)
	unlocks := strings.Split(ctx.GlobalString(utils.UnlockedAccountFlag.Name), ",")
	for i, account := range unlocks {
		if trimmed := strings.TrimSpace(account); trimmed != "" {
			unlockAccount(ctx, ks, trimmed, i, passwords)
		}
	}
	//注册钱包事件处理程序以打开和自动派生钱包
	events := make(chan accounts.WalletEvent, 16)
	stack.AccountManager().Subscribe(events)

	go func() {

		//创建一个rpcclient
		log.Info("stack.Attach()")
		rpcClient, err := stack.Attach()
		if err != nil {
			utils.Fatalf("Failed to attach to self: %v", err)
		}
		stateReader := ethclient.NewClient(rpcClient)

		//打开已经附上的钱包
		log.Info("range stack.AccountManager().Wallets()")
		for _, wallet := range stack.AccountManager().Wallets() {
			if err := wallet.Open(""); err != nil {
				log.Warn("Failed to open wallet", "url", wallet.URL(), "err", err)
			}
		}

		//监听钱包活动直到终止
		log.Info("range events {")
		for event := range events {
			switch event.Kind {
			case accounts.WalletArrived:
				if err := event.Wallet.Open(""); err != nil {
					log.Warn("New wallet appeared, failed to open", "url", event.Wallet.URL(), "err", err)
				}
			case accounts.WalletOpened:
				status, _ := event.Wallet.Status()
				log.Info("New wallet appeared", "url", event.Wallet.URL(), "status", status)

				if event.Wallet.URL().Scheme == "ledger" {
					event.Wallet.SelfDerive(accounts.DefaultLedgerBaseDerivationPath, stateReader)
				} else {
					event.Wallet.SelfDerive(accounts.DefaultBaseDerivationPath, stateReader)
				}

			case accounts.WalletDropped:
				log.Info("Old wallet dropped", "url", event.Wallet.URL())
				event.Wallet.Close()
			}
		}
	}()

	//这里查看一下读取的配置信息
	//chainConfig := ethereum.BlockChain().Config()
	//log.Info("Dpos Validators ", "Size", len(chainConfig.Dpos.Validators))
	//log.Info("Contract Bases ", "Size", len(chainConfig.Contracts.Bases))

	//如果设置为可用，则启动辅助Services
	if ctx.GlobalBool(utils.MiningEnabledFlag.Name) {

		//挖矿只有在运行完整的以太坊节点时才是有意义的
		var ethereum *eth.Ethereum
		if err := stack.Service(&ethereum); err != nil {
			utils.Fatalf("ethereum service not running: %v", err)
		}

		//从CLI和开始挖矿中设置GasPrice的限制
		ethereum.TxPool().SetGasPrice(utils.GlobalBig(ctx, utils.GasPriceFlag.Name))

		if err := ethereum.StartMining(true); err != nil {

			utils.Fatalf("Failed to start mining: %v", err)
		}
	}
}
*/

//启动系统节点和所有已注册的协议，之后它解锁任何请求的帐户，并启动RPC / IPC接口和矿工
func startNode(ctx *cli.Context, stack *node.Node) {

	//启动节点
	utils.StartNode(stack)

	//解锁一些特定需求的账户
	ks := stack.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)

	//加载相应的密码信息
	passwords := utils.MakePasswordList(ctx)
	unlocks := strings.Split(ctx.GlobalString(utils.UnlockedAccountFlag.Name), ",")
	for i, account := range unlocks {
		if trimmed := strings.TrimSpace(account); trimmed != "" {
			unlockAccount(ctx, ks, trimmed, i, passwords)
		}
	}
	//注册钱包事件处理程序以打开和自动派生钱包
	events := make(chan accounts.WalletEvent, 16)
	stack.AccountManager().Subscribe(events)

	go func() {
		//创建一个rpcclient
		rpcClient, err := stack.Attach()
		if err != nil {
			utils.Fatalf("Failed to attach to self: %v", err)
		}
		stateReader := ethclient.NewClient(rpcClient)

		//打开已经附上的钱包
		for _, wallet := range stack.AccountManager().Wallets() {
			if err := wallet.Open(""); err != nil {
				log.Warn("Failed to open wallet", "url", wallet.URL(), "err", err)
			}
		}
		//监听钱包活动直到终止
		for event := range events {
			switch event.Kind {
			case accounts.WalletArrived:
				if err := event.Wallet.Open(""); err != nil {
					log.Warn("New wallet appeared, failed to open", "url", event.Wallet.URL(), "err", err)
				}
			case accounts.WalletOpened:
				status, _ := event.Wallet.Status()
				log.Info("New wallet appeared", "url", event.Wallet.URL(), "status", status)

				if event.Wallet.URL().Scheme == "ledger" {
					event.Wallet.SelfDerive(accounts.DefaultLedgerBaseDerivationPath, stateReader)
				} else {
					event.Wallet.SelfDerive(accounts.DefaultBaseDerivationPath, stateReader)
				}

			case accounts.WalletDropped:
				log.Info("Old wallet dropped", "url", event.Wallet.URL())
				event.Wallet.Close()
			}
		}
	}()

	//如果设置为可用，则启动辅助Services
	if ctx.GlobalBool(utils.MiningEnabledFlag.Name) {

		//挖矿只有在运行完整的以太坊节点时才是有意义的
		var ethereum *eth.Ethereum
		if err := stack.Service(&ethereum); err != nil {
			utils.Fatalf("ethereum service not running: %v", err)
		}
		bokerChain := boker.New(ethereum)
		ethereum.SetBoker(bokerChain)

		//从CLI和开始挖矿中设置GasPrice的限制
		ethereum.TxPool().SetGasPrice(utils.GlobalBig(ctx, utils.GasPriceFlag.Name))
		if err := ethereum.StartMining(true); err != nil {
			utils.Fatalf("Failed to start mining: %v", err)
		}
	}
}
