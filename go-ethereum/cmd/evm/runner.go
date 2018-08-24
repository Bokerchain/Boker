package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/pprof"
	"time"

	goruntime "runtime"

	"github.com/boker/go-ethereum/cmd/evm/internal/compiler"
	"github.com/boker/go-ethereum/cmd/utils"
	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/core"
	"github.com/boker/go-ethereum/core/state"
	"github.com/boker/go-ethereum/core/vm"
	"github.com/boker/go-ethereum/core/vm/runtime"
	"github.com/boker/go-ethereum/ethdb"
	"github.com/boker/go-ethereum/log"
	"github.com/boker/go-ethereum/params"
	cli "gopkg.in/urfave/cli.v1"
)

var runCommand = cli.Command{
	Action:      runCmd,
	Name:        "run",
	Usage:       "run arbitrary evm binary",
	ArgsUsage:   "<code>",
	Description: `The run command runs arbitrary EVM code.`,
}

//将读取给定的JSON格式genesis文件并返回初始化的Genesis结构
func readGenesis(genesisPath string) *core.Genesis {

	//确保我们有一个有效的创世纪JSON (genesisPath := ctx.Args().First())
	if len(genesisPath) == 0 {
		utils.Fatalf("Must supply path to genesis JSON file")
	}

	//打开创世配置文件
	file, err := os.Open(genesisPath)
	if err != nil {
		utils.Fatalf("Failed to read genesis file: %v", err)
	}
	defer file.Close()

	//得到创世配置的Json格式
	genesis := new(core.Genesis)
	if err := json.NewDecoder(file).Decode(genesis); err != nil {
		utils.Fatalf("invalid genesis file: %v", err)
	}
	return genesis
}

//
func runCmd(ctx *cli.Context) error {

	//设置日志相关
	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
	glogger.Verbosity(log.Lvl(ctx.GlobalInt(VerbosityFlag.Name)))
	log.Root().SetHandler(glogger)
	logconfig := &vm.LogConfig{
		DisableMemory: ctx.GlobalBool(DisableMemoryFlag.Name),
		DisableStack:  ctx.GlobalBool(DisableStackFlag.Name),
	}

	var (
		tracer      vm.Tracer
		debugLogger *vm.StructLogger
		statedb     *state.StateDB
		chainConfig *params.ChainConfig
		sender      = common.StringToAddress("sender")
		receiver    = common.StringToAddress("receiver")
	)

	//创建调试日志
	if ctx.GlobalBool(MachineFlag.Name) {
		tracer = NewJSONLogger(logconfig, os.Stdout)
	} else if ctx.GlobalBool(DebugFlag.Name) {
		debugLogger = vm.NewStructLogger(logconfig)
		tracer = debugLogger
	} else {
		debugLogger = vm.NewStructLogger(logconfig)
	}

	//如果有创世文件，则加载创世配置
	if ctx.GlobalString(GenesisFlag.Name) != "" {
		gen := readGenesis(ctx.GlobalString(GenesisFlag.Name))
		_, statedb = gen.ToBlock()
		chainConfig = gen.Config
	} else {
		db, _ := ethdb.NewMemDatabase()
		statedb, _ = state.New(common.Hash{}, state.NewDatabase(db))
	}

	//
	if ctx.GlobalString(SenderFlag.Name) != "" {
		sender = common.HexToAddress(ctx.GlobalString(SenderFlag.Name))
	}
	statedb.CreateAccount(sender)

	if ctx.GlobalString(ReceiverFlag.Name) != "" {
		receiver = common.HexToAddress(ctx.GlobalString(ReceiverFlag.Name))
	}

	var (
		code []byte
		ret  []byte
		err  error
	)
	// The '--code' or '--codefile' flag overrides code in state
	//带有'--code' 或者 '--codefile'参数的执行
	if ctx.GlobalString(CodeFileFlag.Name) != "" {
		var hexcode []byte
		var err error
		// If - is specified, it means that code comes from stdin
		if ctx.GlobalString(CodeFileFlag.Name) == "-" {
			//Try reading from stdin
			if hexcode, err = ioutil.ReadAll(os.Stdin); err != nil {
				fmt.Printf("Could not load code from stdin: %v\n", err)
				os.Exit(1)
			}
		} else {
			// Codefile with hex assembly
			if hexcode, err = ioutil.ReadFile(ctx.GlobalString(CodeFileFlag.Name)); err != nil {
				fmt.Printf("Could not load code from file: %v\n", err)
				os.Exit(1)
			}
		}
		code = common.Hex2Bytes(string(bytes.TrimRight(hexcode, "\n")))

	} else if ctx.GlobalString(CodeFlag.Name) != "" {
		code = common.Hex2Bytes(ctx.GlobalString(CodeFlag.Name))
	} else if fn := ctx.Args().First(); len(fn) > 0 {
		// EASM-file to compile
		src, err := ioutil.ReadFile(fn)
		if err != nil {
			return err
		}
		bin, err := compiler.Compile(fn, src, false)
		if err != nil {
			return err
		}
		code = common.Hex2Bytes(bin)
	}

	//设置初始的Gas
	initialGas := ctx.GlobalUint64(GasFlag.Name)
	runtimeConfig := runtime.Config{
		Origin:   sender,
		State:    statedb,
		GasLimit: initialGas,
		GasPrice: utils.GlobalBig(ctx, PriceFlag.Name),
		Value:    utils.GlobalBig(ctx, ValueFlag.Name),
		EVMConfig: vm.Config{
			Tracer:             tracer,
			Debug:              ctx.GlobalBool(DebugFlag.Name) || ctx.GlobalBool(MachineFlag.Name),
			DisableGasMetering: ctx.GlobalBool(DisableGasMeteringFlag.Name),
		},
	}

	//设置CPU的配置信息
	if cpuProfilePath := ctx.GlobalString(CPUProfileFlag.Name); cpuProfilePath != "" {
		f, err := os.Create(cpuProfilePath)
		if err != nil {
			fmt.Println("could not create CPU profile: ", err)
			os.Exit(1)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			fmt.Println("could not start CPU profile: ", err)
			os.Exit(1)
		}
		defer pprof.StopCPUProfile()
	}

	//设置运行时的链配置信息
	if chainConfig != nil {
		runtimeConfig.ChainConfig = chainConfig
	}
	tstart := time.Now()
	var leftOverGas uint64
	if ctx.GlobalBool(CreateFlag.Name) {
		input := append(code, common.Hex2Bytes(ctx.GlobalString(InputFlag.Name))...)
		ret, _, leftOverGas, err = runtime.Create(input, &runtimeConfig)
	} else {
		if len(code) > 0 {
			statedb.SetCode(receiver, code)
		}
		ret, leftOverGas, err = runtime.Call(receiver, common.Hex2Bytes(ctx.GlobalString(InputFlag.Name)), &runtimeConfig)
	}
	execTime := time.Since(tstart)

	if ctx.GlobalBool(DumpFlag.Name) {
		statedb.IntermediateRoot(true)
		fmt.Println(string(statedb.Dump()))
	}

	if memProfilePath := ctx.GlobalString(MemProfileFlag.Name); memProfilePath != "" {
		f, err := os.Create(memProfilePath)
		if err != nil {
			fmt.Println("could not create memory profile: ", err)
			os.Exit(1)
		}
		if err := pprof.WriteHeapProfile(f); err != nil {
			fmt.Println("could not write memory profile: ", err)
			os.Exit(1)
		}
		f.Close()
	}

	if ctx.GlobalBool(DebugFlag.Name) {
		if debugLogger != nil {
			fmt.Fprintln(os.Stderr, "#### TRACE ####")
			vm.WriteTrace(os.Stderr, debugLogger.StructLogs())
		}
		fmt.Fprintln(os.Stderr, "#### LOGS ####")
		vm.WriteLogs(os.Stderr, statedb.Logs())
	}

	if ctx.GlobalBool(StatDumpFlag.Name) {
		var mem goruntime.MemStats
		goruntime.ReadMemStats(&mem)
		fmt.Fprintf(os.Stderr, `evm execution time: %v
heap objects:       %d
allocations:        %d
total allocations:  %d
GC calls:           %d
Gas used:           %d

`, execTime, mem.HeapObjects, mem.Alloc, mem.TotalAlloc, mem.NumGC, initialGas-leftOverGas)
	}
	if tracer != nil {
		tracer.CaptureEnd(ret, initialGas-leftOverGas, execTime, err)
	} else {
		fmt.Printf("0x%x\n", ret)
		if err != nil {
			fmt.Printf(" error: %v\n", err)
		}
	}

	return nil
}
