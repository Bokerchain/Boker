package core

import (
	"errors"
	"math/big"

	"github.com/boker/chain/boker/api"
	"github.com/boker/chain/boker/protocol"
	"github.com/boker/chain/common"
	"github.com/boker/chain/common/math"
	"github.com/boker/chain/core/vm"
	"github.com/boker/chain/log"
	"github.com/boker/chain/params"
)

var (
	Big0                         = big.NewInt(0)
	errInsufficientBalanceForGas = errors.New("insufficient balance to pay for gas")
)

/*
The State Transitioning Model
状态转换模型

A state transition is a change made when a transaction is applied to the current world state
状态转换 是指用当前的world state来执行交易，并改变当前的world state

The state transitioning model does all all the necessary work to work out a valid new state root.
状态转换做了所有所需的工作来产生一个新的有效的state root

1) Nonce handling	Nonce 处理
2) Pre pay gas		预先支付Gas
3) Create a new state object if the recipient is \0*32		如果接收人是空，那么创建一个新的state object
4) Value transfer	转账
== If contract creation ==
  4a) Attempt to run transaction data		尝试运行输入的数据
  4b) If valid, use result as code for the new state object	如果有效，那么用运行的结果作为新的state object的code
== end ==
5) Run Script section	运行脚本部分
6) Derive new state root	导出新的state root
*/
type StateTransition struct {
	gp         *GasPool //用来追踪区块内部的Gas的使用情况
	msg        Message
	gas        uint64
	gasPrice   *big.Int     // gas的价格
	initialGas *big.Int     // 最开始的gas
	value      *big.Int     // 转账的值
	data       []byte       // 输入数据
	extra      []byte       //扩展字段
	state      vm.StateDB   //StateDB对象
	evm        *vm.EVM      //虚拟机对象
	boker      bokerapi.Api //播客链的接口对象
}

// Message represents a message sent to a contract.
//发送一个合约的消息
type Message interface {
	From() common.Address
	To() *common.Address
	GasPrice() *big.Int
	Gas() *big.Int
	Value() *big.Int
	Nonce() uint64
	CheckNonce() bool
	Data() []byte
	Extra() []byte
	TxType() protocol.TxType
}

//计算Gas。
/*
	这一段代码用来计算使用的Gas数量，从以上的算法可以看出，Gas的算法为
	Gas = 创建合约费用（或者交易费用） + 占用字节费用（合约中实际有数据的长度 * 68） + 非占用字节费用（合约中实际没有数据的长度 * 4）
*/
func IntrinsicGas(data []byte, contractCreation, homestead bool) *big.Int {

	//初始化一个igas变量用来保存计算出来的Gas数量
	igas := new(big.Int)

	//判断是否需要进行合约部署费用
	if contractCreation && homestead {

		//部署合约 53000
		igas.SetUint64(params.TxGasContractCreation)
	} else {

		//交易合约 21000
		igas.SetUint64(params.TxGas)
	}

	if len(data) > 0 {

		//过滤掉合约中的空数据（bye = 0），得到数据长度
		var nz int64
		for _, byt := range data {
			if byt != 0 {
				nz++
			}
		}

		//合约中占用的字节数 * 68 = 占用字节费用
		m := big.NewInt(nz)
		m.Mul(m, new(big.Int).SetUint64(params.TxDataNonZeroGas))
		//将创建合约 + 数据长度所需要的费用设置为igas费用
		igas.Add(igas, m)
		//得到合约数据中为空（byt = 0）的数量
		m.SetInt64(int64(len(data)) - nz)
		//将这个数量 * 4 成为一个数值
		m.Mul(m, new(big.Int).SetUint64(params.TxDataZeroGas))
		//将这个数值和原来的数值进行相加
		igas.Add(igas, m)

	}
	return igas
}

// NewStateTransition initialises and returns a new state transition object.
//创建一个交易的状态对象
func NewStateTransition(evm *vm.EVM, msg Message, gp *GasPool) *StateTransition {
	return &StateTransition{
		gp:         gp,
		evm:        evm,
		msg:        msg,
		gasPrice:   msg.GasPrice(),
		initialGas: new(big.Int),
		value:      msg.Value(),
		data:       msg.Data(),
		extra:      msg.Extra(),
		state:      evm.StateDB,
	}
}

func BinaryMessage(evm *vm.EVM, msg Message, gp *GasPool, boker bokerapi.Api) ([]byte, []byte, *big.Int, bool, error) {

	st := NewStateTransition(evm, msg, gp)
	ret, _, gasUsed, failed, extra, err := st.TransitionDb(boker)
	return ret, extra, gasUsed, failed, err
}

//执行基本合约的消息
func baseMessage(evm *vm.EVM, msg Message, gp *GasPool, boker bokerapi.Api) ([]byte, []byte, *big.Int, bool, error) {

	st := NewStateTransition(evm, msg, gp)
	ret, _, _, failed, extra, err := st.BaseTransitionDb(boker)
	return ret, extra, new(big.Int).SetInt64(0), failed, err
}

//部署基础合约的消息
func contractMessage(evm *vm.EVM, msg Message, gp *GasPool, txType protocol.TxType, boker bokerapi.Api) ([]byte, []byte, *big.Int, bool, error) {

	st := NewStateTransition(evm, msg, gp)
	ret, _, _, failed, extra, err := st.ContractTransitionDb(txType, boker)
	return ret, extra, new(big.Int).SetInt64(0), failed, err
}

//通证分配合约的消息
func validatorMessage(evm *vm.EVM, msg Message, gp *GasPool, txType protocol.TxType, boker bokerapi.Api) ([]byte, []byte, *big.Int, bool, error) {

	st := NewStateTransition(evm, msg, gp)
	ret, _, _, failed, extra, err := st.ValidatorTransitionDb(txType, boker)
	return ret, extra, new(big.Int).SetInt64(0), failed, err
}

//获取交易的from信息
func (st *StateTransition) from() vm.AccountRef {

	f := st.msg.From()
	if !st.state.Exist(f) {
		st.state.CreateAccount(f)
	}

	return vm.AccountRef(f)
}

//获取交易的to信息
func (st *StateTransition) to() vm.AccountRef {
	if st.msg == nil {
		return vm.AccountRef{}
	}
	to := st.msg.To()
	if to == nil {
		return vm.AccountRef{} // contract creation
	}

	reference := vm.AccountRef(*to)
	if !st.state.Exist(*to) {
		st.state.CreateAccount(*to)
	}
	return reference
}

//获取交易中实际使用的Gas信息，从Gas中减去实际费用，得到剩余费用
func (st *StateTransition) useGas(amount uint64) error {
	if st.gas < amount {
		return vm.ErrOutOfGas
	}
	st.gas -= amount

	return nil
}

func (st *StateTransition) buyGas() error {

	//得到消息的Gas
	mgas := st.msg.Gas()
	if mgas.BitLen() > 64 {
		return vm.ErrOutOfGas
	}
	//计算Gas的价格合计 = Gas * GasPrice
	mgval := new(big.Int).Mul(mgas, st.gasPrice)

	var (
		state  = st.state
		sender = st.from()
	)
	//判断用户账户中有足够的费用
	if state.GetBalance(sender.Address()).Cmp(mgval) < 0 {
		return errInsufficientBalanceForGas
	}
	if err := st.gp.SubGas(mgas); err != nil {
		return err
	}
	st.gas += mgas.Uint64()

	st.initialGas.Set(mgas)
	state.SubBalance(sender.Address(), mgval)
	return nil
}

func (st *StateTransition) preCheck() error {
	msg := st.msg
	sender := st.from()

	// Make sure this transaction's nonce is correct
	if msg.CheckNonce() {
		nonce := st.state.GetNonce(sender.Address())
		if nonce < msg.Nonce() {
			return ErrNonceTooHigh
		} else if nonce > msg.Nonce() {
			return ErrNonceTooLow
		}
	}
	return st.buyGas()
}

func (st *StateTransition) getExtra(boker bokerapi.Api) string {

	var contractType protocol.ContractType
	var err error

	//获取合约等级
	if boker == nil {
		log.Info("boker is nil")
		return ""
	}
	if st.msg.To() == nil {
		log.Info("st.msg.To() is nil")
		return ""
	}

	contractType, err = boker.GetContract(*st.msg.To())
	if err != nil {
		return ""
	}

	//判断合约是否是普通合约
	if contractType <= protocol.BinaryContract {
		return ""
	}

	//根据交易类型得到合约的abiJson格式和方法名称
	var name, abiJson string
	abiJson, name, err = boker.GetMethodName(st.msg.TxType())
	if err != nil {
		return ""
	}

	//判断输入参数是否大于0
	if protocol.GetParamCount(abiJson, name) <= 0 {
		return ""
	}

	//解码abi
	var methodJson protocol.MethodJson
	methodJson, err = protocol.DecodeAbi(abiJson, name, string(st.data))
	if err == nil {
		return ""
	}

	//得到最后一个输入参数内容
	return methodJson.Params[len(methodJson.Params)-1].Value
}

//通过应用当前消息并返回结果来转换状态包括操作所需的气体以及用过的气体。 如果它返回错误失败了，表示存在共识问题。
func (st *StateTransition) TransitionDb(boker bokerapi.Api) (ret []byte, requiredGas, usedGas *big.Int, failed bool, extra []byte, err error) {

	//log.Info("****TransitionDb****")
	if err = st.preCheck(); err != nil {
		return
	}

	msg := st.msg
	sender := st.from()
	homestead := true
	contractCreation := msg.To() == nil

	//计算Gas数量
	intrinsicGas := IntrinsicGas(st.data, contractCreation, homestead)
	//log.Info("TransitionDb ", "intrinsicGas", intrinsicGas)

	//判断Gas长度是否超长（长度大于64位）
	if intrinsicGas.BitLen() > 64 {
		return nil, nil, nil, false, []byte(""), vm.ErrOutOfGas
	}

	if err = st.useGas(intrinsicGas.Uint64()); err != nil {
		return nil, nil, nil, false, []byte(""), err
	}

	var (
		evm   = st.evm
		vmerr error
	)

	//判断合约是否存在
	if contractCreation {
		ret, _, st.gas, vmerr = evm.Create(sender, st.data, st.gas, st.value)
	} else {

		//得到扩展字段
		extra = []byte(st.getExtra(boker))
		st.state.SetNonce(sender.Address(), st.state.GetNonce(sender.Address())+1)
		ret, st.gas, vmerr = evm.Call(sender, st.to().Address(), st.data, st.gas, st.value)

		log.Info("evm Call", "ret", ret, "gas", st.gas)
	}
	if vmerr != nil {

		if vmerr == vm.ErrInsufficientBalance {
			return nil, nil, nil, false, []byte(""), vmerr
		}
	}

	//交易所需要的Gas费用
	requiredGas = new(big.Int).Set(st.gasUsed())

	//退还Gas
	st.refundGas()
	st.state.AddBalance(st.evm.Coinbase, new(big.Int).Mul(st.gasUsed(), st.gasPrice))

	return ret, requiredGas, st.gasUsed(), vmerr != nil, extra, err
}

//基本交易执行
func (st *StateTransition) BaseTransitionDb(boker bokerapi.Api) (ret []byte, requiredGas, usedGas *big.Int, failed bool, extra []byte, err error) {

	if err = st.preCheck(); err != nil {
		return
	}

	//判断Gas长度是否超长（长度大于64位）
	intrinsicGas := IntrinsicGas(st.data, false, true)
	if intrinsicGas.BitLen() > 64 {

		return nil, nil, nil, false, []byte(""), vm.ErrOutOfGas
	}

	var (
		evm   = st.evm
		vmerr error
	)

	extra = []byte(st.getExtra(boker))
	st.state.SetNonce(st.from().Address(), st.state.GetNonce(st.from().Address())+1)

	//由于是基础业务，因此这里设置gas为最大值
	st.gas = protocol.MaxGasPrice.Uint64()
	ret, st.gas, vmerr = evm.Call(st.from(), st.to().Address(), st.data, st.gas, st.value)

	if vmerr != nil {
		log.Debug("VM returned with error", "err", vmerr)
		if vmerr == vm.ErrInsufficientBalance {
			return nil, nil, nil, false, []byte(""), vmerr
		}
	}

	return ret, new(big.Int).SetInt64(0), new(big.Int).SetInt64(0), vmerr != nil, extra, err
}

//部署基础合约执行
func (st *StateTransition) ContractTransitionDb(txType protocol.TxType, boker bokerapi.Api) (ret []byte, requiredGas, usedGas *big.Int, failed bool, extra []byte, err error) {

	if err = st.preCheck(); err != nil {
		return
	}

	//为了测试方便，将这一行进行屏蔽
	/*txLevel, err := boker.GetAccount(st.msg.From())
	if err != nil {
		return nil, nil, nil, false, []byte(""), protocol.ErrLevel
	}

	if !bokerapi.ExistsTxType(protocol.SetVote, txLevel) &&
		!bokerapi.ExistsTxType(protocol.SetAssignToken, txLevel) &&
		!bokerapi.ExistsTxType(protocol.CancelVote, txLevel) &&
		!bokerapi.ExistsTxType(protocol.CanclAssignToken, txLevel) {
		return nil, nil, nil, false, []byte(""), protocol.ErrLevel
	}*/

	if txType == protocol.SetPersonalContract {
		boker.SetContract(*st.msg.To(), protocol.PersonalContract, false, string(st.extra))
	} else if txType == protocol.SetSystemContract {
		boker.SetContract(*st.msg.To(), protocol.SystemContract, false, string(st.extra))
	} else if txType == protocol.CancelPersonalContract {
		boker.SetContract(*st.msg.To(), protocol.PersonalContract, true, "")
	} else if txType == protocol.CancelSystemContract {
		boker.SetContract(*st.msg.To(), protocol.SystemContract, true, "")
	}

	extra = st.extra
	st.state.SetNonce(st.from().Address(), st.state.GetNonce(st.from().Address())+1)
	return ret, new(big.Int).SetInt64(0), new(big.Int).SetInt64(0), false, extra, err
}

//投票交易操作
func (st *StateTransition) VoteTransitionDb(txType protocol.TxType, boker bokerapi.Api) (ret []byte, requiredGas, usedGas *big.Int, failed bool, extra []byte, err error) {

	if err = st.preCheck(); err != nil {
		return
	}
	extra = []byte(st.getExtra(boker))
	st.state.SetNonce(st.from().Address(), st.state.GetNonce(st.from().Address())+1)
	return ret, new(big.Int).SetInt64(0), new(big.Int).SetInt64(0), false, extra, err
}

//投票交易操作
func (st *StateTransition) ValidatorTransitionDb(txType protocol.TxType, boker bokerapi.Api) (ret []byte, requiredGas, usedGas *big.Int, failed bool, extra []byte, err error) {

	if err = st.preCheck(); err != nil {
		return
	}
	extra = []byte(st.getExtra(boker))
	st.state.SetNonce(st.from().Address(), st.state.GetNonce(st.from().Address())+1)
	return ret, new(big.Int).SetInt64(0), new(big.Int).SetInt64(0), false, extra, err
}

//退还Gas
func (st *StateTransition) refundGas() {
	// Return eth for remaining gas to the sender account,
	// exchanged at the original rate.
	sender := st.from() // err already checked
	remaining := new(big.Int).Mul(new(big.Int).SetUint64(st.gas), st.gasPrice)
	st.state.AddBalance(sender.Address(), remaining)

	// Apply refund counter, capped to half of the used gas.
	uhalf := remaining.Div(st.gasUsed(), common.Big2)
	refund := math.BigMin(uhalf, st.state.GetRefund())
	st.gas += refund.Uint64()

	st.state.AddBalance(sender.Address(), refund.Mul(refund, st.gasPrice))

	// Also return remaining gas to the block gas counter so it is
	// available for the next transaction.
	st.gp.AddGas(new(big.Int).SetUint64(st.gas))
}

func (st *StateTransition) gasUsed() *big.Int {
	return new(big.Int).Sub(st.initialGas, new(big.Int).SetUint64(st.gas))
}
