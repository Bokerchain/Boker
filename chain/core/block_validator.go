package core

import (
	"fmt"
	"math/big"

	"github.com/Bokerchain/Boker/chain/common/math"
	"github.com/Bokerchain/Boker/chain/consensus"
	"github.com/Bokerchain/Boker/chain/core/state"
	"github.com/Bokerchain/Boker/chain/core/types"
	"github.com/Bokerchain/Boker/chain/params"
)

type BlockValidator struct {
	config *params.ChainConfig //链的配置信息
	bc     *BlockChain         //规范块链
	engine consensus.Engine    //验证采用的共识机制
}

var (
	curGaslimit = big.NewFloat(0)
	rtts        = big.NewFloat(0)
	rttd        = big.NewFloat(0)
)

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

func CalcGasLimit(parent *types.Block) *big.Int {

	contrib := new(big.Int).Mul(parent.GasUsed(), big.NewInt(3))
	contrib = contrib.Div(contrib, big.NewInt(2))
	contrib = contrib.Div(contrib, params.GasLimitBoundDivisor)

	//decay = parentGasLimit / 1024 -1
	decay := new(big.Int).Div(parent.GasLimit(), params.GasLimitBoundDivisor)
	decay.Sub(decay, big.NewInt(1))

	gl := new(big.Int).Sub(parent.GasLimit(), decay)
	gl = gl.Add(gl, contrib)

	gl.Set(math.BigMax(gl, params.MinGasLimit))

	if gl.Cmp(params.TargetGasLimit) < 0 {
		gl.Add(parent.GasLimit(), decay)
		gl.Set(math.BigMin(gl, params.TargetGasLimit))
	}
	return gl
}

//
func bound(lower *big.Float, middle *big.Float, upper *big.Float) *big.Float {

	if lower.Cmp(middle) == -1 {

		//lower < middle
		if middle.Cmp(upper) == -1 {

			//middle < upper
			return middle
		} else {

			//middle > upper
			return upper
		}
	} else if lower.Cmp(middle) == 0 {

		//lower == middle
		if middle.Cmp(upper) == -1 {

			//middle < upper
			return middle
		} else {

			//middle > upper
			return upper
		}

	} else {

		//lower > middle
		if lower.Cmp(upper) == -1 {

			//lower < upper
			return lower
		} else {

			//lower > upper
			return upper
		}
	}
}

//计算门限阀值均值
func calcSsthresh(gas *big.Int) *big.Int {

	//求均值
	gasFloat := new(big.Float).SetInt(gas)
	if curGaslimit.Cmp(big.NewFloat(0)) == 0 {

		rtts = gasFloat
		rttd = rttd.Quo(rtts, big.NewFloat(2.0))
	} else {

		rtts = rtts.Add(rtts.Mul(rtts, big.NewFloat(0.875)), rtts.Mul(gasFloat, big.NewFloat(0.125)))
		absFloat := new(big.Float).Mul(big.NewFloat(0.25), new(big.Float).Abs(new(big.Float).Sub(gasFloat, rtts)))
		rttd = rttd.Add(rttd.Mul(rttd, big.NewFloat(0.75)), absFloat)
	}
	curGaslimit = curGaslimit.Add(rtts, new(big.Float).Mul(rttd, big.NewFloat(4.0)))
	curGaslimit = bound(params.MinGasLimitFloat, curGaslimit, params.GasLimitSsthresh)

	//转换成big.Int类型
	ssthreshInt := new(big.Int)
	curGaslimit.Int(ssthreshInt)
	return ssthreshInt
}

//计算Gaslimit
/*func CalcGasLimit(parent *types.Block) *big.Int {

	return calcSsthresh(parent.GasUsed())
}*/
