package protocol

import (
	_ "bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"

	"github.com/Bokerchain/Boker/chain/accounts/abi"
	"github.com/Bokerchain/Boker/chain/common"
	"github.com/Bokerchain/Boker/chain/crypto/sha3"
	"github.com/Bokerchain/Boker/chain/log"
	"github.com/Bokerchain/Boker/chain/rlp"
)

const (
	ExtraVanity        = 32               //扩展字段的前缀字节数量
	ExtraSeal          = 65               //扩展字段的后缀字节数量
	InmemorySignatures = 4096             //保留在内存中的最近块签名的数量
	ProducerInterval   = int64(5)         //打包时间间隔（秒）
	TokenNoderInterval = int64(300)       //分配通证时间间隔(秒)
	EpochInterval      = int64(86400)     //一个周期的时间（86400秒 = 1天）
	MaxValidatorSize   = 1                //DPOS的验证者数量
	SafeSize           = 1                //安全的验证者数量
	ConsensusSize      = 1                //共识确认验证者数量
	BokerInterval      = time.Second      //分配通证时间间隔(秒)
	AssignTickInterval = time.Millisecond //分配通证时间间隔(秒)
	AssignInterval     = time.Minute      //分配通证时间间隔单位
	AssignTimer        = 5
)

//新增多个交易类型
type TxType uint8

const (
	Binary TxType = iota //原来的转账或者合约调用交易

	//设置验证者
	SetValidator

	/****设置和取消用户基础合约交易类型****/
	SetPersonalContract
	CancelPersonalContract

	/****设置和取消系统基础合约交易类型****/
	SetSystemContract
	CancelSystemContract

	/****用户基础合约交易类型****/
	RegisterCandidate //注册成为候选人(用户注册为候选人)
	VoteUser          //用户投票
	VoteCancel        //用户取消投票
	VoteEpoch         //产生当前的出块节点(在每次周期产生的时候触发)
	UserEvent         //用户数据上传

	/****系统基础合约交易类型****/
	AssignToken //分配通证(每次分配通证的时候触发)
	//AssignReward //出块节点的通证奖励(每次分配通证的时候触发)
)

//新增合约类型
type ContractType uint8

const (
	BinaryContract   ContractType = iota //普通合约类型
	SystemContract                       //系统基础合约
	PersonalContract                     //个人基础合约
)

var (
	BobbyUnit          *big.Int = big.NewInt(1e+17) //Bobby的单位
	BobbyMultiple      *big.Int = big.NewInt(110)   //倍数
	TransferUnit       *big.Int = big.NewInt(1e+17) //转账单位(这个数值仅用于每次给指定账号，方便指定账号给用户分配通证)
	TransferMultiple   *big.Int = big.NewInt(165)   //转账倍数
	SetValidatorVotes  *big.Int = big.NewInt(10000)
	MaxGasPrice        *big.Int = new(big.Int).SetUint64(0xffffffffffffffff) //最大的GasPrice
	MaxGasLimit        *big.Int = new(big.Int).SetUint64(0)                  //最大的GasLimit
	TimeOfFirstBlock            = int64(0)                                   //创世区块的时间偏移量
	ConfirmedBlockHead          = []byte("confirmed-block-head")
)

var (
	//用户触发的合约方法名（用户触发，但是不收取Gas费用）
	RegisterCandidateMethod = "registerCandidate" //申请候选节点
	VoteCandidateMethod     = "voteCandidate"     //投票候选节点
	CancelVoteMethod        = "cancelAllVotes"    //取消所有投票
	FireEventMethod         = "fireUserEvent"     //用户行为数据上报

	//基础链触发的基础合约
	AssignTokenMethod   = "assignToken"   //分配通证
	RotateVoteMethod    = "rotateVote"    //产生当前的出块节点(在每次周期产生的时候触发)
	TickCandidateMethod = "tickVote"      //投票时钟
	GetCandidateMethod  = "getCandidates" //获取候选人结果
)

var (
	EpochPrefix     = []byte("epoch-")     //存放周期信息
	ValidatorPrefix = []byte("validator-") //存放验证者投票信息
	BlockCntPrefix  = []byte("blockCnt-")  //存放投票数量
	ValidatorsKey   = []byte("validators") //存放所有的验证者列表
	SinglePrefix    = []byte("single-")    //
	Contracts       = []byte("contracts")  //
	AbiPrefix       = []byte("abi-")       //
)

var (
	ErrNilBlockHeader             = errors.New("nil block header returned")                       //区块头为空
	ErrUnknownBlock               = errors.New("unknown block")                                   //未知区块
	ErrInvalidProducer            = errors.New("invalid current producer")                        //出块节点出错
	ErrInvalidTokenNoder          = errors.New("invalid current token noder")                     //当前分配通证节点出错
	ErrInvalidProducerTime        = errors.New("invalid time to mint the block")                  //不正确的出块时间
	ErrInvalidTokenTime           = errors.New("invalid time to assign token noder")              //错误的分币节点
	ErrInvalidCoinbase            = errors.New("invalid current mining coinbase")                 //当前挖矿账号错误
	ErrInvalidSystem              = errors.New("invalid current system")                          //当前系统的投票合约出错
	ErrMismatchSignerAndValidator = errors.New("mismatch block signer and validator")             //签名者和区块头中的验证者不是同一个
	ErrNoSigner                   = errors.New("missing signing methods")                         //缺少签名方法
	ErrInvalidType                = errors.New("invalid transaction type")                        //无效的交易类型
	ErrInvalidAddress             = errors.New("invalid transaction payload address")             //无效的交易有效负载地
	ErrInvalidAction              = errors.New("invalid transaction payload action")              //无效的事务有效负载操
	ErrLoadConfig                 = errors.New("load bokerchain config error")                    //加载配置信息出错
	ErrNotFoundAddress            = errors.New("not found bokerchain contract address")           //没有找到合约地址
	ErrNotFoundType               = errors.New("not found bokerchain contract type")              //没有找到合约类型
	ErrWriteJson                  = errors.New("write bokerchain json file error")                //写保存基础合约的Json格式出错
	ErrOpenFile                   = errors.New("open bokerchain json file error")                 //打开基础合约保存文件出错
	ErrWriteFile                  = errors.New("bokerchain write file error")                     //写基础合约保存文件出错
	ErrContractExist              = errors.New("bokerchain contract aleady exist")                //写基础合约保存文件出错
	ErrSystem                     = errors.New("system error")                                    //系统错误
	ErrNotFoundContract           = errors.New("not found bokerchain contract")                   //没有找到合约
	ErrNotFoundAccount            = errors.New("not found bokerchain account")                    //没有找到合约
	ErrNewContractService         = errors.New("create bokerchain base contract err")             //没有找到合约
	ErrSaveContractTrie           = errors.New("save contract trie err")                          //没有找到合约
	ErrLevel                      = errors.New("account level error")                             //没有找到合约
	ErrPointerIsNil               = errors.New("Trie Pointer is Nil")                             //Hash树指针是nil
	ErrTransactionType            = errors.New("Error Transaction Type")                          //交易类型错误
	ErrSpecialAccount             = errors.New("Current Account is`t BokerChain Special Account") //当前账号不是指定的特殊账号
	ErrValidatorsIsFull           = errors.New("Current Validators is Full")                      //当前验证者数量已满
	ErrExistsValidators           = errors.New("Current Validators Exists")                       //当前存在验证者
	ErrGenesisBlock               = errors.New("not genesis block")                               //区块需要不为0，即最近区块不是创世区块，证明已经工作
	ErrDecodeValidators           = errors.New("failed to decode validators")
	ErrEncodeValidators           = errors.New("failed to encode validators")
	ErrSetEpochTrieFail           = errors.New("failed set epoch trie")
	ErrEpochTrieNil               = errors.New("failed to producers length is zero")
	ErrToIsNil                    = errors.New("setValidator block header to is nil")
	ErrTxType                     = errors.New("failed to tx type")
)

//设置播客链配置
type BokerConfig struct {
	Address common.Address
}

type BokerBackendProto struct {
	SingleHash     common.Hash `json:"SingleRoot"        gencodec:"required"`
	ContractsHash  common.Hash `json:"ContractsRoot"    gencodec:"required"`
	ContracAbiHash common.Hash `json:"ContractABIRoot"    gencodec:"required"`
}

func (p *BokerBackendProto) Root() (h common.Hash) {

	hw := sha3.NewKeccak256()
	rlp.Encode(hw, p.SingleHash)
	rlp.Encode(hw, p.ContractsHash)
	rlp.Encode(hw, p.ContracAbiHash)
	hw.Sum(h[:0])
	return h
}

func ToBokerProto(singleHash common.Hash, contractsHash common.Hash, contractAbi common.Hash) *BokerBackendProto {

	return &BokerBackendProto{
		SingleHash:     singleHash,
		ContractsHash:  contractsHash,
		ContracAbiHash: contractAbi,
	}
}

//Abi函数参数信息
type ParamJson struct {
	Name  string `json:"name"`  //参数名称
	Type  string `json:"type"`  //参数类型
	Value string `json:"value"` //参数值
}
type MethodJson struct {
	Method string      `json:"method"`           //方法名称
	Params []ParamJson `json:"params,omitempty"` //方法参数列表
}

func getInterface(input abi.Argument) reflect.Value {

	//设置参数类型
	var param reflect.Value

	switch input.Type.T {
	case abi.SliceTy:

		log.Info("DecodeAbi abi.SliceTy")
		var paramType = []byte{}
		param = reflect.New(reflect.TypeOf(paramType))
	case abi.StringTy:

		log.Info("DecodeAbi abi.StringTy")
		var paramType = string("")
		param = reflect.New(reflect.TypeOf(paramType))
	case abi.IntTy:

		log.Info("DecodeAbi abi.IntTy", "input.Type.Type", input.Type.Type)
		switch input.Type.Type {

		case reflect.TypeOf(int8(0)):

			var paramType = int8(0)
			param = reflect.New(reflect.TypeOf(paramType))
		case reflect.TypeOf(int16(0)):

			var paramType = int16(0)
			param = reflect.New(reflect.TypeOf(paramType))
		case reflect.TypeOf(int32(0)):

			var paramType = int32(0)
			param = reflect.New(reflect.TypeOf(paramType))
		case reflect.TypeOf(int64(0)):

			var paramType = int64(0)
			param = reflect.New(reflect.TypeOf(paramType))
		case reflect.TypeOf(&big.Int{}):

			var paramType = big.NewInt(1)
			param = reflect.New(reflect.TypeOf(paramType))
		}
	case abi.UintTy:

		log.Info("DecodeAbi abi.UintTy", "input.Type.Type", input.Type.Type)
		switch input.Type.Type {

		case reflect.TypeOf(uint8(0)):

			var paramType = uint8(0)
			param = reflect.New(reflect.TypeOf(paramType))
		case reflect.TypeOf(uint16(0)):

			var paramType = uint16(0)
			param = reflect.New(reflect.TypeOf(paramType))
		case reflect.TypeOf(uint32(0)):

			var paramType = uint32(0)
			param = reflect.New(reflect.TypeOf(paramType))
		case reflect.TypeOf(uint64(0)):

			var paramType = uint64(0)
			param = reflect.New(reflect.TypeOf(paramType))
		case reflect.TypeOf(&big.Int{}):

			var paramType = big.NewInt(1)
			param = reflect.New(reflect.TypeOf(paramType))
		}
	case abi.BoolTy:

		log.Info("DecodeAbi abi.BoolTy")
		var paramType = bool(true)
		param = reflect.New(reflect.TypeOf(paramType))
	case abi.AddressTy:

		log.Info("DecodeAbi abi.AddressTy")
		var paramType = common.Address{}
		param = reflect.New(reflect.TypeOf(paramType))
	case abi.BytesTy:

		log.Info("DecodeAbi abi.BytesTy")
		var paramType = common.Hex2Bytes("0100000000000000000000000000000000000000000000000000000000000000")
		param = reflect.New(reflect.TypeOf(paramType))
	case abi.FunctionTy:

		log.Info("DecodeAbi abi.FunctionTy")
		var paramType = common.Hex2Bytes("0100000000000000000000000000000000000000000000000000000000000000")
		param = reflect.New(reflect.TypeOf(paramType))
	}
	return param
}

func GetParamCount(abiJson string, name string) int {

	abiDecoder, err := abi.JSON(strings.NewReader(abiJson))
	if err != nil {
		return 0
	}
	return len(abiDecoder.Methods[name].Inputs)
}

func DecodeAbi(abiJson string, name string, payload string) (MethodJson, error) {

	const definition = `[{"constant":false,"inputs":[{"name":"","type":"int256"},{"name":"str","type":"string"}],"name":"test","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"show","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"payable":false,"stateMutability":"nonpayable","type":"fallback"}]`
	payload = "069fd0a30000000000000000000000000000000000000000000000000000000000000005000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000046861686100000000000000000000000000000000000000000000000000000000"
	name = "test"

	//解析Abi格式成为Json格式
	abiDecoder, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		return MethodJson{}, err
	}
	log.Info("DecodeAbi", "Methods Inputs count", len(abiDecoder.Methods[name].Inputs))

	//剔除最前面的0x标记
	var decodeString string = ""
	hexFlag := strings.Index(payload, "0x")
	if hexFlag == -1 {
		decodeString = payload
	} else {
		decodeString = payload[2:]
	}

	//将字符串转换成[]Byte
	decodeBytes, err := hex.DecodeString(decodeString)
	if err != nil {
		return MethodJson{}, err
	}
	log.Info("DecodeAbi", "decodeBytes", decodeBytes)

	//根据函数的名称，设置函数的输入参数信息
	method, ok := abiDecoder.Methods[name]
	if !ok {
		return MethodJson{}, errors.New("")
	}

	//写入获取参数类型
	params := make([]interface{}, 0)
	for i := 0; i < len(method.Inputs); i++ {

		input := method.Inputs[i]

		//设置参数类型
		var param reflect.Value

		switch input.Type.T {
		case abi.SliceTy:

			log.Info("DecodeAbi abi.SliceTy")
			var paramType = []byte{}
			param = reflect.New(reflect.TypeOf(paramType))
		case abi.StringTy:

			log.Info("DecodeAbi abi.StringTy")
			var paramType = string("")
			param = reflect.New(reflect.TypeOf(paramType))
		case abi.IntTy:

			log.Info("DecodeAbi abi.IntTy", "input.Type.Type", input.Type.Type)
			switch input.Type.Type {

			case reflect.TypeOf(int8(0)):

				var paramType = int8(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(int16(0)):

				var paramType = int16(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(int32(0)):

				var paramType = int32(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(int64(0)):

				var paramType = int64(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(&big.Int{}):

				var paramType = big.NewInt(1)
				param = reflect.New(reflect.TypeOf(paramType))
			}
		case abi.UintTy:

			log.Info("DecodeAbi abi.UintTy", "input.Type.Type", input.Type.Type)
			switch input.Type.Type {

			case reflect.TypeOf(uint8(0)):

				var paramType = uint8(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(uint16(0)):

				var paramType = uint16(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(uint32(0)):

				var paramType = uint32(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(uint64(0)):

				var paramType = uint64(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(&big.Int{}):

				var paramType = big.NewInt(1)
				param = reflect.New(reflect.TypeOf(paramType))
			}
		case abi.BoolTy:

			log.Info("DecodeAbi abi.BoolTy")
			var paramType = bool(true)
			param = reflect.New(reflect.TypeOf(paramType))
		case abi.AddressTy:

			log.Info("DecodeAbi abi.AddressTy")
			var paramType = common.Address{}
			param = reflect.New(reflect.TypeOf(paramType))
		case abi.BytesTy:

			log.Info("DecodeAbi abi.BytesTy")
			var paramType = common.Hex2Bytes("0100000000000000000000000000000000000000000000000000000000000000")
			param = reflect.New(reflect.TypeOf(paramType))
		case abi.FunctionTy:

			log.Info("DecodeAbi abi.FunctionTy")
			var paramType = common.Hex2Bytes("0100000000000000000000000000000000000000000000000000000000000000")
			param = reflect.New(reflect.TypeOf(paramType))
		}
		params = append(params, param.Interface())
	}

	//解码
	if err := abiDecoder.InputUnpack(params, name, decodeBytes[4:]); err != nil {
		log.Error("DecodeAbi ", "err", err)
		return MethodJson{}, err
	}

	//将返回的信息放入到Json格式中
	json := MethodJson{}
	json.Method = name
	json.Params = make([]ParamJson, 0)

	for i := 0; i < len(params); i++ {

		valueOf := reflect.ValueOf(params[i])
		out := valueOf.Elem().Interface()
		s := fmt.Sprintf("%v", out)

		param := ParamJson{
			Name:  abiDecoder.Methods[name].Inputs[i].Name,
			Type:  abiDecoder.Methods[name].Inputs[i].Type.String(),
			Value: s,
		}
		json.Params = append(json.Params, param)
	}
	log.Info("DecodeAbi ", "Json", json)

	/*for i := 0; i < len(params); i++ {

		valueOf := reflect.ValueOf(params[i])
		out := valueOf.Elem().Interface()
		s := fmt.Sprintf("%v", out)
		outArray = append(outArray, s)
		log.Info("DecodeAbi ", "Param", s)
	}*/

	return json, nil
}

//这里提供了一个函数，来对合约的abi进行解析
/*func DecodeAbi(abiJson string, name string, payload string) (MethodJson, error) {

	const definition = `[{"constant":false,"inputs":[{"name":"","type":"int256"},{"name":"str","type":"string"}],"name":"test","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"show","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"payable":false,"stateMutability":"nonpayable","type":"fallback"}]`
	payload = "069fd0a30000000000000000000000000000000000000000000000000000000000000005000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000046861686100000000000000000000000000000000000000000000000000000000"
	name = "test"

	//解析Abi格式成为Json格式
	abiDecoder, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		return MethodJson{}, err
	}
	log.Info("DecodeAbi", "Methods Inputs count", len(abiDecoder.Methods[name].Inputs))

	//剔除最前面的0x标记
	var decodeString string = ""
	hexFlag := strings.Index(payload, "0x")
	if hexFlag == -1 {
		decodeString = payload
	} else {
		decodeString = payload[2:]
	}

	//将字符串转换成[]Byte
	decodeBytes, err := hex.DecodeString(decodeString)
	if err != nil {
		return MethodJson{}, err
	}
	log.Info("DecodeAbi", "decodeBytes", decodeBytes)

	//根据函数的名称，设置函数的输入参数信息
	method, ok := abiDecoder.Methods[name]
	if !ok {
		return MethodJson{}, errors.New("")
	}

	params := make([]interface{}, 0)
	for i := 0; i < len(method.Inputs); i++ {

		//设置参数类型
		var param reflect.Value

		paramType := getInterface(method.Inputs[i])
		param = reflect.New(reflect.TypeOf(paramType))

		switch method.Inputs[i].Type.T {
		case abi.SliceTy:

			log.Info("DecodeAbi abi.SliceTy")
			var paramType = []byte{}
			param = reflect.New(reflect.TypeOf(paramType))
		case abi.StringTy:

			log.Info("DecodeAbi abi.StringTy")
			var paramType = string("")
			param = reflect.New(reflect.TypeOf(paramType))
		case abi.IntTy:

			log.Info("DecodeAbi", "Methods Inputs index", i, "method.Inputs[i].Type.Type", method.Inputs[i].Type.Type)
			switch method.Inputs[i].Type.Type {

			case reflect.TypeOf(int8(0)):

				var paramType = int8(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(int16(0)):

				var paramType = int16(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(int32(0)):

				var paramType = int32(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(int64(0)):

				var paramType = int64(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(&big.Int{}):

				var paramType = big.NewInt(1)
				param = reflect.New(reflect.TypeOf(paramType))
			}
		case abi.UintTy:

			log.Info("DecodeAbi", "Methods Inputs index", i, "method.Inputs[i].Type.Type", method.Inputs[i].Type.Type)
			switch method.Inputs[i].Type.Type {

			case reflect.TypeOf(uint8(0)):

				var paramType = uint8(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(uint16(0)):

				var paramType = uint16(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(uint32(0)):

				var paramType = uint32(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(uint64(0)):

				var paramType = uint64(0)
				param = reflect.New(reflect.TypeOf(paramType))
			case reflect.TypeOf(&big.Int{}):

				var paramType = big.NewInt(1)
				param = reflect.New(reflect.TypeOf(paramType))
			}
		case abi.BoolTy:

			log.Info("DecodeAbi abi.BoolTy")
			var paramType = bool(true)
			param = reflect.New(reflect.TypeOf(paramType))
		case abi.AddressTy:

			log.Info("DecodeAbi abi.AddressTy")
			var paramType = common.Address{}
			param = reflect.New(reflect.TypeOf(paramType))
		case abi.BytesTy:

			log.Info("DecodeAbi abi.BytesTy")
			var paramType = common.Hex2Bytes("0100000000000000000000000000000000000000000000000000000000000000")
			param = reflect.New(reflect.TypeOf(paramType))
		case abi.FunctionTy:

			log.Info("DecodeAbi abi.FunctionTy")
			var paramType = common.Hex2Bytes("0100000000000000000000000000000000000000000000000000000000000000")
			param = reflect.New(reflect.TypeOf(paramType))
		}
		params = append(params, param.Interface())
	}
	if err := abiDecoder.InputUnpack(params, name, decodeBytes[4:]); err != nil {
		log.Error("DecodeAbi ", "err", err)
		return MethodJson{}, err
	}

	//将返回的信息放入到Json格式中
	json := MethodJson{}

	for i := 0; i < len(params); i++ {

		valueOf := reflect.ValueOf(params[i])
		out := valueOf.Elem().Interface()
		s := fmt.Sprintf("%v", out)
		outArray = append(outArray, s)
		log.Info("DecodeAbi ", "Param", s)
	}

	return json, nil
}*/
