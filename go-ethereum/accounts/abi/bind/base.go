package bind

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/boker/go-ethereum"
	"github.com/boker/go-ethereum/accounts/abi"
	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/core/types"
	"github.com/boker/go-ethereum/crypto"
	"github.com/boker/go-ethereum/eth"
	"github.com/boker/go-ethereum/include"
	"github.com/boker/go-ethereum/node"
)

// SignerFn is a signer function callback when a contract requires a method to
// sign the transaction before submission.
type SignerFn func(types.Signer, common.Address, *types.Transaction) (*types.Transaction, error)

// CallOpts is the collection of options to fine tune a contract call request.
type CallOpts struct {
	Pending bool            // Whether to operate on the pending state or the last known one
	From    common.Address  // Optional the sender address, otherwise the first account is used
	Context context.Context // Network context to support cancellation and timeouts (nil = no timeout)
}

//创建一个有效的以太坊交易
type TransactOpts struct {
	From     common.Address  // Ethereum account to send the transaction from
	Nonce    *big.Int        // Nonce to use for the transaction execution (nil = use pending state)
	Signer   SignerFn        // Method to use for signing the transaction (mandatory)
	Value    *big.Int        // Funds to transfer along along the transaction (nil = 0 = no funds)
	GasPrice *big.Int        // Gas price to use for the transaction execution (nil = gas price oracle)
	GasLimit *big.Int        // Gas limit to set for the transaction execution (nil = estimate + 10%)
	Context  context.Context // Network context to support cancellation and timeouts (nil = no timeout)
}

//BoundContract定义以太坊合约的基础包装器对象 它包含一组由方法使用的方法更高级别的合同绑定操作。
type BoundContract struct {
	address    common.Address     // Deployment address of the contract on the Ethereum blockchain
	abi        abi.ABI            // Reflect based ABI to access the correct Ethereum methods
	caller     ContractCaller     // Read interface to interact with the blockchain
	transactor ContractTransactor // Write interface to interact with the blockchain
}

var GethNode *node.Node

//NewBoundContract 创建一个通过其调用的低级合约接口并且交易可以通过。
func NewBoundContract(address common.Address,
	abi abi.ABI,
	caller ContractCaller,
	transactor ContractTransactor) *BoundContract {
	return &BoundContract{
		address:    address,
		abi:        abi,
		caller:     caller,
		transactor: transactor,
	}
}

func DeployContract(opts *TransactOpts, abi abi.ABI, bytecode []byte, backend ContractBackend, params ...interface{}) (common.Address, *types.Transaction, *BoundContract, error) {

	//赋值
	c := NewBoundContract(common.Address{}, abi, backend, backend)

	input, err := c.abi.Pack("", params...)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	tx, err := c.transact(opts, nil, append(bytecode, input...), types.Binary)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	c.address = crypto.CreateAddress(opts.From, tx.Nonce())
	return c.address, tx, c, nil
}

//调用合约方法，并将params作为输入值和将输出设置为result
func (c *BoundContract) Call(opts *CallOpts, result interface{}, method string, params ...interface{}) error {

	//判断opts是否为空
	if opts == nil {
		opts = new(CallOpts)
	}

	//打包输入，调用并解压缩结果
	input, err := c.abi.Pack(method, params...)
	if err != nil {
		return err
	}

	var (
		msg    = ethereum.CallMsg{From: opts.From, To: &c.address, Data: input}
		ctx    = ensureContext(opts.Context)
		code   []byte
		output []byte
	)

	if opts.Pending {

		pb, ok := c.caller.(PendingContractCaller)
		if !ok {
			return ErrNoPendingState
		}

		output, err = pb.PendingCallContract(ctx, msg)
		if err == nil && len(output) == 0 {
			// Make sure we have a contract to operate on, and bail out otherwise.
			if code, err = pb.PendingCodeAt(ctx, c.address); err != nil {
				return err
			} else if len(code) == 0 {
				return ErrNoCode
			}
		}

	} else {

		output, err = c.caller.CallContract(ctx, msg, nil)
		if err == nil && len(output) == 0 {
			// Make sure we have a contract to operate on, and bail out otherwise.
			if code, err = c.caller.CodeAt(ctx, c.address, nil); err != nil {
				return err
			} else if len(code) == 0 {
				return ErrNoCode
			}
		}

	}
	if err != nil {
		return err
	}
	return c.abi.Unpack(result, method, output)
}

//得到当前分币帐号
func (c *BoundContract) getTokenNoder(opts *TransactOpts) (common.Address, error) {

	var ether *eth.Ethereum
	if err := GethNode.Service(&ether); err != nil {
		return common.Address{}, err
	}

	if ether.BlockChain().CurrentBlock() == nil {
		return common.Address{}, errors.New("failed to lookup token node")
	}

	return ether.BlockChain().CurrentBlock().DposCtx().GetCurrentTokenNoder()
}

//得到当前的验证者帐号
func (c *BoundContract) getProducer(opts *TransactOpts) (common.Address, error) {

	var ether *eth.Ethereum
	if err := GethNode.Service(&ether); err != nil {
		return common.Address{}, err
	}

	if ether.BlockChain().CurrentBlock() == nil {
		return common.Address{}, errors.New("failed to lookup token node")
	}

	return ether.BlockChain().CurrentBlock().DposCtx().GetCurrentProducer()
}

//使用输入的值作为参数调用合约方法
func (c *BoundContract) Transact(opts *TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {

	//打包合约参数
	input, err := c.abi.Pack(method, params...)
	if err != nil {
		return nil, err
	}

	var ether *eth.Ethereum
	if err := GethNode.Service(&ether); err != nil {
		return nil, err
	}

	contractType, err := ether.Boker.GetContract(c.address)
	if err != nil {
		return nil, err
	}

	if contractType == types.ContractVote { //当前调用的是投票基础合约

		if method == include.RegisterCandidateMethod {

			//候选人注册方法名
			return c.transact(opts, &c.address, input, types.RegisterCandidate)
		} else if method == include.VoteCandidateMethod {

			//投票方法名
			return c.transact(opts, &c.address, input, types.ProducerVote)
		} else if method == include.RotateVoteMethod {

			//转换投票方法名
			return c.transact(opts, &c.address, input, types.RotateVote)
		}
		return nil, errors.New("vote contract unknown method")

	} else if contractType == types.ContractAssignToken {

		//判断是否是通证分配协议
		if method == include.AssignTokenMethod {

			//得到当前的分币节点
			tokennoder, err := c.getTokenNoder(opts)
			if err != nil {
				return nil, errors.New("get assign token error")
			}

			//判断是否一致
			if tokennoder != opts.From {
				return nil, errors.New("current assign token not is from account")
			}
			return c.transact(opts, &c.address, input, types.AssignToken)
		}
		return nil, errors.New("assign token contract unknown method")

	}
	return c.transact(opts, &c.address, input, types.Binary)

}

func (c *BoundContract) Transfer(opts *TransactOpts) (*types.Transaction, error) {

	var ether *eth.Ethereum
	if err := GethNode.Service(&ether); err != nil {
		return nil, err
	}

	txType, err := ether.Boker.GetContract(c.address)
	if err != nil {
		return nil, err
	}

	return c.transact(opts, &c.address, nil, types.TxType(txType))
}

func (c *BoundContract) baseTransact(opts *TransactOpts, contract *common.Address, input []byte, transactTypes types.TxType) (*types.Transaction, error) {

	//判断Value值是否为空
	var err error
	value := opts.Value
	if value == nil {
		value = new(big.Int)
	}

	//判断Nonce值是否为空
	var nonce uint64
	if opts.Nonce == nil {

		//如果Nonce值为空，则初始化一个nonce值来进行初始化
		nonce, err = c.transactor.PendingNonceAt(ensureContext(opts.Context), opts.From)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
		}
	} else {
		nonce = opts.Nonce.Uint64()
	}

	/*这里不对Gas和GasLimit进行设置，因为在函数NewBaseTransaction里面已经针对基础业务进行了设置*/
	var rawTx *types.Transaction
	if contract == nil {

		//如果合约尚未创建，则创建合约
		//rawTx = types.NewBaseContractCreation(nonce, value, input)
		return nil, errors.New("not found base contract address")
	} else {

		//合约已经创建，则创建一个交易
		rawTx = types.NewBaseTransaction(transactTypes, nonce, c.address, value, input)
	}

	//判断交易是否有签名者
	if opts.Signer == nil {
		return nil, errors.New("no signer to authorize the transaction with")
	}

	//进行签名
	signedTx, err := opts.Signer(types.HomesteadSigner{}, opts.From, rawTx)
	if err != nil {
		return nil, err
	}

	//将交易注入pending池中
	if err := c.transactor.SendTransaction(ensureContext(opts.Context), signedTx); err != nil {
		return nil, err
	}
	return signedTx, nil
}

func (c *BoundContract) normalTransact(opts *TransactOpts, contract *common.Address, input []byte, transactTypes types.TxType) (*types.Transaction, error) {

	//判断Value值是否为空
	var err error
	value := opts.Value
	if value == nil {
		value = new(big.Int)
	}

	//判断Nonce值是否为空
	var nonce uint64
	if opts.Nonce == nil {

		//如果Nonce值为空，则初始化一个nonce值来进行初始化
		nonce, err = c.transactor.PendingNonceAt(ensureContext(opts.Context), opts.From)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
		}
	} else {
		nonce = opts.Nonce.Uint64()
	}

	//如果GasPrice为空，则设置一个建议的GasPrice
	gasPrice := opts.GasPrice
	if gasPrice == nil {
		gasPrice, err = c.transactor.SuggestGasPrice(ensureContext(opts.Context)) //得到一个建议的GasPrice
		if err != nil {
			return nil, fmt.Errorf("failed to suggest gas price: %v", err)
		}
	}

	//如果GasLimit为空，则设置一个GasLimit
	gasLimit := opts.GasLimit
	if gasLimit == nil {

		//如果合约存在，则根据合约内容评估一个GasLimit
		if contract != nil {

			if code, err := c.transactor.PendingCodeAt(ensureContext(opts.Context), c.address); err != nil {
				return nil, err
			} else if len(code) == 0 {
				return nil, ErrNoCode
			}
		}

		//估算所需要的Gas
		msg := ethereum.CallMsg{From: opts.From, To: contract, Value: value, Data: input}
		gasLimit, err = c.transactor.EstimateGas(ensureContext(opts.Context), msg)
		if err != nil {
			return nil, fmt.Errorf("failed to estimate gas needed: %v", err) //估算所需gas失败
		}
	}

	//创建合约交易或者直接产生一个交易
	var rawTx *types.Transaction
	if contract == nil {
		//如果合约尚未创建，则创建合约
		rawTx = types.NewContractCreation(nonce, value, gasLimit, gasPrice, input)
	} else {
		//合约已经创建，则创建一个交易
		rawTx = types.NewTransaction(transactTypes, nonce, c.address, value, gasLimit, gasPrice, input)
	}

	//判断交易是否有签名者
	if opts.Signer == nil {
		return nil, errors.New("no signer to authorize the transaction with")
	}

	//进行签名
	signedTx, err := opts.Signer(types.HomesteadSigner{}, opts.From, rawTx)
	if err != nil {
		return nil, err
	}

	//将交易注入pending池中
	if err := c.transactor.SendTransaction(ensureContext(opts.Context), signedTx); err != nil {
		return nil, err
	}
	return signedTx, nil
}

func (c *BoundContract) transact(opts *TransactOpts, contract *common.Address, input []byte, transactTypes types.TxType) (*types.Transaction, error) {

	/*根据不同类型计算使用的Gas信息*/
	if transactTypes == types.Binary {

		//普通交易
		return c.normalTransact(opts, contract, input, transactTypes)
	} else if types.IsDeploy(transactTypes) || types.IsVote(transactTypes) || types.IsToken(transactTypes) {

		//基础合约交易
		return c.baseTransact(opts, contract, input, transactTypes)
	} else {

		//未知的类型
		return nil, errors.New("unknown transaction type")
	}
}

//
func ensureContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.TODO()
	}
	return ctx
}
