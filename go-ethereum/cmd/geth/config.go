package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"unicode"

	cli "gopkg.in/urfave/cli.v1"

	"github.com/boker/go-ethereum/cmd/utils"

	_ "github.com/boker/go-ethereum/boker"
	"github.com/boker/go-ethereum/contracts/release"
	"github.com/boker/go-ethereum/dashboard"
	"github.com/boker/go-ethereum/eth"
	"github.com/boker/go-ethereum/log"
	"github.com/boker/go-ethereum/node"
	"github.com/boker/go-ethereum/params"
	whisper "github.com/boker/go-ethereum/whisper/whisperv5"
	"github.com/naoina/toml"
)

var (
	dumpConfigCommand = cli.Command{
		Action:      utils.MigrateFlags(dumpConfig),
		Name:        "dumpconfig",
		Usage:       "Show configuration values",
		ArgsUsage:   "",
		Flags:       append(append(nodeFlags, rpcFlags...), whisperFlags...),
		Category:    "MISCELLANEOUS COMMANDS",
		Description: `The dumpconfig command shows configuration values.`,
	}

	configFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TOML configuration file",
	}
)

// These settings ensure that TOML keys use the same names as Go struct fields.
var tomlSettings = toml.Config{
	NormFieldName: func(rt reflect.Type, key string) string {
		return key
	},
	FieldToKey: func(rt reflect.Type, field string) string {
		return field
	},
	MissingField: func(rt reflect.Type, field string) error {
		link := ""
		if unicode.IsUpper(rune(rt.Name()[0])) && rt.PkgPath() != "main" {
			link = fmt.Sprintf(", see https://godoc.org/%s#%s for available fields", rt.PkgPath(), rt.Name())
		}
		return fmt.Errorf("field '%s' is not defined in %s%s", field, rt.String(), link)
	},
}

type ethstatsConfig struct {
	URL string `toml:",omitempty"`
}

type gethConfig struct {
	Eth       eth.Config     //Eth配置
	Shh       whisper.Config //
	Node      node.Config    //节点配置
	Ethstats  ethstatsConfig
	Dashboard dashboard.Config
}

func loadConfig(file string, cfg *gethConfig) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tomlSettings.NewDecoder(bufio.NewReader(f)).Decode(cfg)
	// Add file name to errors that have a line number.
	if _, ok := err.(*toml.LineError); ok {
		err = errors.New(file + ", " + err.Error())
	}
	return err
}

func defaultNodeConfig() node.Config {
	cfg := node.DefaultConfig
	cfg.Name = clientIdentifier
	cfg.Version = params.VersionWithCommit(gitCommit)
	cfg.HTTPModules = append(cfg.HTTPModules, "eth", "shh")
	cfg.WSModules = append(cfg.WSModules, "eth", "shh")
	cfg.IPCPath = "geth.ipc"
	return cfg
}

//根据配置信息产生一个节点和这个节点的配置
func makeConfigNode(ctx *cli.Context) (*node.Node, gethConfig) {

	log.Info("****makeConfigNode****")

	//加载默认配置
	cfg := gethConfig{
		Eth:       eth.DefaultConfig,
		Shh:       whisper.DefaultConfig,
		Node:      defaultNodeConfig(),
		Dashboard: dashboard.DefaultConfig,
	}
	log.Info("makeConfigNode gethConfig")

	//加载默认配置文件()
	if file := ctx.GlobalString(configFileFlag.Name); file != "" {
		if err := loadConfig(file, &cfg); err != nil {
			log.Error("%v", err)
		}
	}
	log.Info("makeConfigNode GlobalString")

	//应用标记
	utils.SetNodeConfig(ctx, &cfg.Node)
	log.Info("makeConfigNode SetNodeConfig")

	//根据配置信息生成一个节点
	stack, err := node.New(&cfg.Node)
	if err != nil {
		log.Error("Failed to create the protocol stack: %v", err)
	}
	log.Info("makeConfigNode node.New")

	//设置eth的配置
	utils.SetEthConfig(ctx, stack, &cfg.Eth)
	log.Info("makeConfigNode SetEthConfig")

	if ctx.GlobalIsSet(utils.EthStatsURLFlag.Name) {
		cfg.Ethstats.URL = ctx.GlobalString(utils.EthStatsURLFlag.Name)
	}
	log.Info("makeConfigNode GlobalIsSet")

	utils.SetShhConfig(ctx, stack, &cfg.Shh)
	log.Info("makeConfigNode SetShhConfig")

	utils.SetDashboardConfig(ctx, &cfg.Dashboard)
	log.Info("makeConfigNode SetDashboardConfig")

	return stack, cfg
}

// enableWhisper returns true in case one of the whisper flags is set.
func enableWhisper(ctx *cli.Context) bool {
	for _, flag := range whisperFlags {
		if ctx.GlobalIsSet(flag.GetName()) {
			return true
		}
	}
	return false
}

//产生一个全节点
/*func makeFullNode(ctx *cli.Context) *node.Node {

	//产生一个节点的配置
	log.Info("makeConfigNode")
	stack, cfg := makeConfigNode(ctx)

	//注册构造函数
	log.Info("RegisterEthService")
	utils.RegisterEthService(stack, &cfg.Eth)

	log.Info("RegisterDashboardService")
	if ctx.GlobalBool(utils.DashboardEnabledFlag.Name) {
		utils.RegisterDashboardService(stack, &cfg.Dashboard)
	}

	log.Info("enableWhisper")
	// Whisper must be explicitly enabled by specifying at least 1 whisper flag
	shhEnabled := enableWhisper(ctx)
	if shhEnabled {
		if ctx.GlobalIsSet(utils.WhisperMaxMessageSizeFlag.Name) {
			cfg.Shh.MaxMessageSize = uint32(ctx.Int(utils.WhisperMaxMessageSizeFlag.Name))
		}
		if ctx.GlobalIsSet(utils.WhisperMinPOWFlag.Name) {
			cfg.Shh.MinimumAcceptedPOW = ctx.Float64(utils.WhisperMinPOWFlag.Name)
		}
		utils.RegisterShhService(stack, &cfg.Shh)
	}

	// Add the Ethereum Stats daemon if requested.
	log.Info("RegisterEthStatsService")
	if cfg.Ethstats.URL != "" {
		utils.RegisterEthStatsService(stack, cfg.Ethstats.URL)
	}

	// Add the release oracle service so it boots along with node.
	log.Info("NewReleaseService")
	if err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		config := release.Config{
			Oracle: relOracle,
			Major:  uint32(params.VersionMajor),
			Minor:  uint32(params.VersionMinor),
			Patch:  uint32(params.VersionPatch),
		}
		commit, _ := hex.DecodeString(gitCommit)
		copy(config.Commit[:], commit)
		return release.NewReleaseService(ctx, config)
	}); err != nil {
		utils.Fatalf("Failed to register the Geth release oracle service: %v", err)
	}

	log.Info("Set Ethereum Boker")
	if err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {

		//设置播客链指针
		var ethereum *eth.Ethereum
		if err := ctx.Service(&ethereum); err == nil {
			bokerChain := boker.New(ethereum)
			ethereum.SetBoker(bokerChain)
		}
		return nil, nil
	}); err != nil {
		utils.Fatalf("Failed to register the Geth release oracle service: %v", err)
	}

	return stack
}
*/

//产生一个全节点
func makeFullNode(ctx *cli.Context) *node.Node {

	//产生一个节点的配置
	stack, cfg := makeConfigNode(ctx)
	utils.RegisterEthService(stack, &cfg.Eth)

	if ctx.GlobalBool(utils.DashboardEnabledFlag.Name) {
		utils.RegisterDashboardService(stack, &cfg.Dashboard)
	}
	// Whisper must be explicitly enabled by specifying at least 1 whisper flag
	shhEnabled := enableWhisper(ctx)
	if shhEnabled {
		if ctx.GlobalIsSet(utils.WhisperMaxMessageSizeFlag.Name) {
			cfg.Shh.MaxMessageSize = uint32(ctx.Int(utils.WhisperMaxMessageSizeFlag.Name))
		}
		if ctx.GlobalIsSet(utils.WhisperMinPOWFlag.Name) {
			cfg.Shh.MinimumAcceptedPOW = ctx.Float64(utils.WhisperMinPOWFlag.Name)
		}
		utils.RegisterShhService(stack, &cfg.Shh)
	}

	// Add the Ethereum Stats daemon if requested.
	if cfg.Ethstats.URL != "" {
		utils.RegisterEthStatsService(stack, cfg.Ethstats.URL)
	}

	// Add the release oracle service so it boots along with node.
	if err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		config := release.Config{
			Oracle: relOracle,
			Major:  uint32(params.VersionMajor),
			Minor:  uint32(params.VersionMinor),
			Patch:  uint32(params.VersionPatch),
		}
		commit, _ := hex.DecodeString(gitCommit)
		copy(config.Commit[:], commit)
		return release.NewReleaseService(ctx, config)
	}); err != nil {
		utils.Fatalf("Failed to register the Geth release oracle service: %v", err)
	}

	//在这里启动调用基础合约的协程
	/*if err := stack.Register(func(ctx *node.ServiceContext) (node.Service, error) {

		log.Info("Create Base Contract Manager")
		//utils.Fatalf("Create Base Contract Manager")
		stack.Ctx = ctx
		return nil, contracts.NewBaseContractManager(ctx)
	}); err != nil {
		log.Error("Failed to register the Base Contract service: %v", err)
		utils.Fatalf("Failed to register the Base Contract service: %v", err)
	}*/
	return stack
}

// dumpConfig is the dumpconfig command.
func dumpConfig(ctx *cli.Context) error {
	_, cfg := makeConfigNode(ctx)
	comment := ""

	if cfg.Eth.Genesis != nil {
		cfg.Eth.Genesis = nil
		comment += "# Note: this config doesn't contain the genesis block.\n\n"
	}

	out, err := tomlSettings.Marshal(&cfg)
	if err != nil {
		return err
	}
	io.WriteString(os.Stdout, comment)
	os.Stdout.Write(out)
	return nil
}
