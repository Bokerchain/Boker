package core

import (
	"fmt"
	"math/big"

	"github.com/boker/chain/common/math"
	"github.com/boker/chain/consensus"
	"github.com/boker/chain/core/state"
	"github.com/boker/chain/core/types"
	"github.com/boker/chain/params"
)

type BlockValidator struct {
	config *params.ChainConfig //链的配置信息
	bc     *BlockChain         //规范块链
	engine consensus.Engine    //验证采用的共识机制
}

func NewBlockValidator(config *params.ChainConfig,
	blockchain *BlockChain,
	engine consensus.Engine) *BlockValidator {

	return &BlockValidator{
		config: config,
		engine: engine,
		bc:     blockchain,
	}
}

//验证给定的区块叔块以及验证区块交易头
func (v *BlockValidator) ValidateBody(block *types.Block) error {

	//检查区块是否已知，如果不知道，它是否可链接
	if v.bc.HasBlockAndState(block.Hash()) {
		return ErrKnownBlock
	}
	if !v.bc.HasBlockAndState(block.ParentHash()) {
		return consensus.ErrUnknownAncestor
	}

	//此时已知区块头有效性，检查叔块和交易
	header := block.Header()
	if err := v.engine.VerifyUncles(v.bc, block); err != nil {
		return err
	}
	if hash := types.CalcUncleHash(block.Uncles()); hash != header.UncleHash {

		return fmt.Errorf("uncle root hash mismatch: have %x, want %x", hash, header.UncleHash)
	}
	if hash := types.DeriveSha(block.Transactions()); hash != header.TxHash {

		return fmt.Errorf("transaction root hash mismatch: have %x, want %x", hash, header.TxHash)
	}
	return nil
}

//验证者状态后发生的各种更改转换，例如用过的Gas，收据根和状态根本身 如果验证者成功，ValidateState将返回数据库批处理否则为nil并返回错误。
func (v *BlockValidator) ValidateState(block, parent *types.Block,
	statedb *state.StateDB,
	receipts types.Receipts,
	usedGas *big.Int) error {

	header := block.Header()
	if block.GasUsed().Cmp(usedGas) != 0 {
		return fmt.Errorf("invalid gas used (remote: %v local: %v)", block.GasUsed(), usedGas)
	}

	// Validate the received block's bloom with the one derived from the generated receipts.
	// For valid blocks this should always validate to true.
	rbloom := types.CreateBloom(receipts)
	if rbloom != header.Bloom {
		return fmt.Errorf("invalid bloom (remote: %x  local: %x)", header.Bloom, rbloom)
	}

	// Tre receipt Trie's root (R = (Tr [[H1, R1], ... [Hn, R1]]))
	receiptSha := types.DeriveSha(receipts)
	if receiptSha != header.ReceiptHash {
		return fmt.Errorf("invalid receipt root hash (remote: %x local: %x)", header.ReceiptHash, receiptSha)
	}

	// Validate the state root against the received state root and throw
	// an error if they don't match.
	if root := statedb.IntermediateRoot(v.config.IsEIP158(header.Number)); header.Root != root {
		return fmt.Errorf("invalid merkle root (remote: %x local: %x)", header.Root, root)
	}
	return nil
}

func (v *BlockValidator) ValidateDposState(block *types.Block) error {

	//验证区块中记录的Dpos根和区块头中记录的Dpos根是否一致
	header := block.Header()
	localRoot := block.DposCtx().Root()
	remoteRoot := header.DposProto.Root()
	if remoteRoot != localRoot {
		return fmt.Errorf("invalid dpos root (remote: %x local: %x)", remoteRoot, localRoot)
	}
	return nil
}

//计算Gaslimit(使用父节点) 此处可以采用TCP的慢门限阀值方式来处理
func CalcGasLimit(parent *types.Block) *big.Int {

	/*
		parent.GasUsed() 父交易消耗的总gas数量
			以下代码实现功能为：
			contrib = (parentGasUsed * 3 / 2) / 1024
			contrib :=((parent.GasUsed() * 3) / 2) / params.GasLimitBoundDivisor

			原码为：contrib := (parent.GasUsed() + parent.GasUsed()/2) / params.GasLimitBoundDivisor

	*/
	contrib := new(big.Int).Mul(parent.GasUsed(), big.NewInt(3))
	contrib = contrib.Div(contrib, big.NewInt(2))
	contrib = contrib.Div(contrib, params.GasLimitBoundDivisor)

	//decay = parentGasLimit / 1024 -1
	decay := new(big.Int).Div(parent.GasLimit(), params.GasLimitBoundDivisor)
	decay.Sub(decay, big.NewInt(1))

	/*
		当父节点使用的Gas大于父节点的GasLimit的2/3的时候，则加大GasLimit，否则减小GasLimit
	*/

	//limit := parent.GasLimit() - decay + contrib
	gl := new(big.Int).Sub(parent.GasLimit(), decay)
	gl = gl.Add(gl, contrib)

	/*
		if limit < params.MinGasLimit {
			limit = params.MinGasLimit
		}
	*/
	gl.Set(math.BigMax(gl, params.MinGasLimit))

	/*
		however, if we're now below the target (TargetGasLimit) we increase the
		 limit as much as we can (parentGasLimit / 1024 -1)
		但是，如果我们现在低于目标（TargetGasLimit），我们会增加尽可能多地限制（parentGasLimit / 1024 -1）

		if limit < params.TargetGasLimit {
			limit = parent.GasLimit() + decay
			if limit > params.TargetGasLimit {
				limit = params.TargetGasLimit
			}
		}
	*/

	if gl.Cmp(params.TargetGasLimit) < 0 {
		gl.Add(parent.GasLimit(), decay)
		gl.Set(math.BigMin(gl, params.TargetGasLimit))
	}
	return gl
}
