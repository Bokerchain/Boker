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
	"math/big"

	"github.com/boker/go-ethereum/bokerface"
	"github.com/boker/go-ethereum/core/state"
	"github.com/boker/go-ethereum/core/types"
	"github.com/boker/go-ethereum/core/vm"
)

//区块验证者接口
type Validator interface {
	ValidateBody(block *types.Block) error
	ValidateState(block, parent *types.Block, state *state.StateDB, receipts types.Receipts, usedGas *big.Int) error
	ValidateDposState(block *types.Block) error
}

//区块处理器接口
type Processor interface {
	Process(block *types.Block, statedb *state.StateDB, cfg vm.Config) (types.Receipts, []*types.Log, *big.Int, error)
	SetBoker(boker bokerface.BokerInterface)
}
