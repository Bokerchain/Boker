// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	_ "bytes"
	"errors"
	"math/big"

	"github.com/boker/chain/boker/api"
	"github.com/boker/chain/boker/protocol"
	"github.com/boker/chain/common"
	"github.com/boker/chain/consensus"
	"github.com/boker/chain/consensus/misc"
	"github.com/boker/chain/core/state"
	"github.com/boker/chain/core/types"
	"github.com/boker/chain/core/vm"
	"github.com/boker/chain/crypto"
	"github.com/boker/chain/log"
	"github.com/boker/chain/params"
)

//状态处理器，负责一个从一个节点到另一个节点
type StateProcessor struct {
	config *params.ChainConfig //链配置选项
	bc     *BlockChain         //规范块链
	engine consensus.Engine    //共识引擎
	boker  bokerapi.Api        //播客链的接口
}

//初始化一个新的状态处理器。
func NewStateProcessor(config *params.ChainConfig, bc *BlockChain, engine consensus.Engine) *StateProcessor {
	return &StateProcessor{
		config: config,
		bc:     bc,
		engine: engine,
	}
}

func (p *StateProcessor) Process(block *types.Block, statedb *state.StateDB, cfg vm.Config) (types.Receipts, []*types.Log, *big.Int, error) {

	var (
		receipts     types.Receipts
		totalUsedGas = big.NewInt(0)
		header       = block.Header()
		allLogs      []*types.Log
		gp           = new(GasPool).AddGas(block.GasLimit())
	)

	//根据任何硬叉规范改变块和状态
	if p.config.DAOForkSupport && p.config.DAOForkBlock != nil && p.config.DAOForkBlock.Cmp(block.Number()) == 0 {
		misc.ApplyDAOHardFork(statedb)
	}

	//得到区块中所有的交易，并将这些交易使用Dpos引擎进行执行。
	for i, tx := range block.Transactions() {

		//设置当前statedb状态,以便后面evm创建交易日志
		statedb.Prepare(tx.Hash(), block.Hash(), i)
		receipt, _, err := ApplyTransaction(p.config, block.DposCtx(), p.bc, nil, gp, statedb, header, tx, totalUsedGas, cfg, p.boker)
		if err != nil {
			return nil, nil, nil, err
		}

		//执行完毕的交易回执放入到回执数组中
		receipts = append(receipts, receipt)
		allLogs = append(allLogs, receipt.Logs...)
	}

	//执行完块中所有的交易，应用任何共识引擎特定的附加功能（例如块奖励）
	p.engine.Finalize(p.bc, header, statedb, block.Transactions(), block.Uncles(), receipts, block.DposCtx(), p.boker)

	//返回执行成功的回执数组/日志/以及总的使用Gas的数量
	return receipts, allLogs, totalUsedGas, nil
}

//处理普通交易
func binaryTransaction(config *params.ChainConfig,
	dposContext *types.DposContext,
	bc *BlockChain,
	author *common.Address,
	gp *GasPool,
	statedb *state.StateDB,
	header *types.Header,
	tx *types.Transaction,
	usedGas *big.Int,
	cfg vm.Config,
	msg types.Message,
	boker bokerapi.Api) (*types.Receipt, *big.Int, error) {

	msg, err := tx.AsMessage(types.MakeSigner(config, header.Number))
	if err != nil {
		return nil, nil, err
	}

	context := NewEVMContext(msg, header, bc, author)
	vmenv := vm.NewEVM(context, statedb, config, cfg)
	_, extra, gas, failed, err := BinaryMessage(vmenv, msg, gp, boker)
	if err != nil {
		return nil, nil, err
	}
	tx.SetExtra(extra)

	//用待处理的更改更新状态
	var root []byte
	if config.IsByzantium(header.Number) {
		statedb.Finalise(true)
	} else {
		root = statedb.IntermediateRoot(config.IsEIP158(header.Number)).Bytes()
	}
	usedGas.Add(usedGas, gas)

	//为交易创建一个新收据，存储tx使用的中间根和gas基于eip阶段，我们传递了根触发删除帐户。
	receipt := types.NewReceipt(root, failed, usedGas)
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = new(big.Int).Set(gas)

	//如果交易创建了合同，则将创建地址存储在收据中
	if msg.To() == nil {
		receipt.ContractAddress = crypto.CreateAddress(vmenv.Context.Origin, tx.Nonce())
		log.Info("binaryTransaction", "address", receipt.ContractAddress.String())
	}

	//设置收据日志并创建一个用于过滤的布尔值
	receipt.Logs = statedb.GetLogs(tx.Hash())
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})

	return receipt, gas, err
}

//设置部署基础交易
func contractSetTransaction(config *params.ChainConfig,
	dposContext *types.DposContext,
	bc *BlockChain,
	author *common.Address,
	gp *GasPool,
	statedb *state.StateDB,
	header *types.Header,
	tx *types.Transaction,
	usedGas *big.Int,
	cfg vm.Config,
	msg types.Message,
	boker bokerapi.Api) (*types.Receipt, *big.Int, error) {

	log.Info("****contractSetTransaction****")

	msg, err := tx.AsMessage(types.MakeSigner(config, header.Number))
	if err != nil {
		log.Error("contractSetTransaction tx.AsMessage", "msg", msg, "err", err)
		return nil, nil, err
	}

	context := NewEVMContext(msg, header, bc, author)
	vmenv := vm.NewEVM(context, statedb, config, cfg)
	_, extra, gas, failed, err := contractMessage(vmenv, msg, gp, msg.TxType(), boker)
	if err != nil {
		log.Error("contractSetTransaction contractMessage", "extra", extra, "gas", gas, "failed", failed, "err", err)
		return nil, nil, err
	}
	tx.SetExtra(extra)

	var root []byte
	if config.IsByzantium(header.Number) {
		statedb.Finalise(true)
	} else {
		root = statedb.IntermediateRoot(config.IsEIP158(header.Number)).Bytes()
	}
	usedGas.Add(usedGas, gas)

	receipt := types.NewReceipt(root, failed, usedGas)
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = new(big.Int).Set(gas)
	receipt.Logs = statedb.GetLogs(tx.Hash())
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	return receipt, gas, err
}

//用户投票合约
func baseTransaction(config *params.ChainConfig,
	dposContext *types.DposContext,
	bc *BlockChain,
	author *common.Address,
	gp *GasPool,
	statedb *state.StateDB,
	header *types.Header,
	tx *types.Transaction,
	usedGas *big.Int,
	cfg vm.Config,
	msg types.Message,
	boker bokerapi.Api) (*types.Receipt, *big.Int, error) {

	log.Info("****baseTransaction****")

	msg, err := tx.AsMessage(types.MakeSigner(config, header.Number))
	if err != nil {

		log.Error("baseTransaction AsMessage", "err", err)
		return nil, nil, err
	}

	//判断是否是分配通证合约
	if tx.Type() == protocol.AssignToken || tx.Type() == protocol.AssignReward {

		tokenNoder, err := dposContext.GetCurrentTokenNoder()
		if err != nil {

			log.Error("baseTransaction dposContext.GetCurrentTokenNoder", "err", err)
			return nil, nil, err
		}
		if tokenNoder != msg.From() {

			log.Error("baseTransaction failed tokenNoder != msg.From()", "tokenNoder", tokenNoder, "msg.From()", msg.From())
			return nil, nil, errors.New("from address not assign token producer")
		}
	}

	context := NewEVMContext(msg, header, bc, author)
	vmenv := vm.NewEVM(context, statedb, config, cfg)
	_, extra, gas, failed, err := baseMessage(vmenv, msg, gp, boker)
	if err != nil {
		log.Error("baseTransaction failed", "err", err)
		return nil, nil, err
	}
	tx.SetExtra(extra)

	var root []byte
	if config.IsByzantium(header.Number) {
		statedb.Finalise(true)
	} else {
		root = statedb.IntermediateRoot(config.IsEIP158(header.Number)).Bytes()
	}
	usedGas.Add(usedGas, gas)

	receipt := types.NewReceipt(root, failed, usedGas)
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = new(big.Int).Set(gas)
	receipt.Logs = statedb.GetLogs(tx.Hash())
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})

	log.Info("****baseTransaction End****", "gas", gas, "err", err)
	return receipt, gas, err
}

//设置设置验证人
func validatorTransaction(config *params.ChainConfig,
	dposContext *types.DposContext,
	bc *BlockChain,
	author *common.Address,
	gp *GasPool,
	statedb *state.StateDB,
	header *types.Header,
	tx *types.Transaction,
	usedGas *big.Int,
	cfg vm.Config,
	msg types.Message,
	boker bokerapi.Api) (*types.Receipt, *big.Int, error) {

	msg, err := tx.AsMessage(types.MakeSigner(config, header.Number))
	if err != nil {
		return nil, nil, err
	}

	context := NewEVMContext(msg, header, bc, author)
	vmenv := vm.NewEVM(context, statedb, config, cfg)

	producer, err := dposContext.GetCurrentProducer()
	if protocol.ErrEpochTrieNil == err {

		if bc.CurrentBlock().Number().Int64() == 0 {

			//log.Info("validatorTransaction validatorMessage", "txType", msg.TxType())
			_, extra, gas, failed, err := validatorMessage(vmenv, msg, gp, msg.TxType(), boker)
			if err != nil {
				return nil, nil, err
			}
			tx.SetExtra(extra)
			//log.Info("validatorTransaction", "tx.extra", tx.Extra(), "extra", extra)

			//设置验证者
			dposContext.Clean()
			dposContext.InsertValidator(*msg.To(), protocol.SetValidatorVotes)
			//log.Info("validatorTransaction InsertValidator", "root", dposContext.Root().String())

			root := statedb.IntermediateRoot(false).Bytes()
			usedGas.Add(usedGas, gas)
			//log.Info("validatorTransaction Add", "gas", gas, "usedGas", usedGas, "root", root)

			receipt := types.NewReceipt(root, failed, usedGas)
			receipt.TxHash = tx.Hash()
			receipt.GasUsed = new(big.Int).Set(gas)
			receipt.Logs = statedb.GetLogs(tx.Hash())
			receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
			//log.Info("validatorTransaction CreateBloom")

			return receipt, gas, err
		}
		return nil, nil, errors.New("current block number is`t zero")
	}

	if producer != msg.From() {
		return nil, nil, errors.New("from address not assign token producer")
	}

	_, extra, gas, failed, err := validatorMessage(vmenv, msg, gp, msg.TxType(), boker)
	if err != nil {
		return nil, nil, err
	}
	tx.SetExtra(extra)

	//设置验证者
	//log.Info("validatorTransaction validatorMessage", "txType", msg.TxType())
	dposContext.Clean()
	dposContext.InsertValidator(*msg.To(), protocol.SetValidatorVotes)

	var root []byte
	if config.IsByzantium(header.Number) {
		statedb.Finalise(true)
	} else {
		root = statedb.IntermediateRoot(config.IsEIP158(header.Number)).Bytes()
	}
	usedGas.Add(usedGas, gas)
	receipt := types.NewReceipt(root, failed, usedGas)
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = new(big.Int).Set(gas)
	receipt.Logs = statedb.GetLogs(tx.Hash())
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	return receipt, gas, err
}

//执行交易
func ApplyTransaction(config *params.ChainConfig,
	dposContext *types.DposContext,
	bc *BlockChain,
	author *common.Address,
	gp *GasPool,
	statedb *state.StateDB,
	header *types.Header,
	tx *types.Transaction,
	usedGas *big.Int,
	cfg vm.Config,
	boker bokerapi.Api) (*types.Receipt, *big.Int, error) {

	msg, err := tx.AsMessage(types.MakeSigner(config, header.Number))
	if err != nil {
		return nil, nil, err
	}
	log.Info("****ApplyTransaction****", "Number", header.Number.String(), "txType", msg.TxType(), "from", msg.From())

	if msg.TxType() == protocol.Binary {

		return binaryTransaction(config, dposContext, bc, author, gp, statedb, header, tx, usedGas, cfg, msg, boker)
	} else {

		if msg.To() == nil {
			return nil, nil, protocol.ErrToIsNil
		}

		//根据交易类型来区分
		switch msg.TxType() {

		case protocol.SetPersonalContract, protocol.CancelPersonalContract, protocol.SetSystemContract, protocol.CancelSystemContract:
			//设置合约(已经测试)
			return contractSetTransaction(config, dposContext, bc, author, gp, statedb, header, tx, usedGas, cfg, msg, boker)
		case protocol.VoteUser, protocol.VoteEpoch, protocol.AssignToken, protocol.AssignReward, protocol.RegisterCandidate: //基础交易(已经测试)

			return baseTransaction(config, dposContext, bc, author, gp, statedb, header, tx, usedGas, cfg, msg, boker)
		case protocol.SetValidator: //设置验证人(已经测试)

			return validatorTransaction(config, dposContext, bc, author, gp, statedb, header, tx, usedGas, cfg, msg, boker)
		default:

			return nil, nil, protocol.ErrInvalidType
		}
	}
}

func (p *StateProcessor) SetBoker(boker bokerapi.Api) { p.boker = boker }
