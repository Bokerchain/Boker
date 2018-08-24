// Package consensus implements different Ethereum consensus engines.
package consensus

import (
	"errors"

	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/core/state"
	"github.com/boker/go-ethereum/core/types"
	"github.com/boker/go-ethereum/params"
	"github.com/boker/go-ethereum/rpc"
)

var (
	// ErrUnknownAncestor is returned when validating a block requires an ancestor
	// that is unknown.
	ErrUnknownAncestor = errors.New("unknown ancestor")

	// ErrFutureBlock is returned when a block's timestamp is in the future according
	// to the current node.
	ErrFutureBlock = errors.New("block in the future")

	// ErrInvalidNumber is returned if a block's number doesn't equal it's parent's
	// plus one.
	ErrInvalidNumber = errors.New("invalid block number")
)

// ChainReader defines a small collection of methods needed to access the local
// blockchain during header and/or uncle verification.
type ChainReader interface {
	Config() *params.ChainConfig
	CurrentHeader() *types.Header
	GetHeader(hash common.Hash, number uint64) *types.Header
	GetHeaderByNumber(number uint64) *types.Header
	GetHeaderByHash(hash common.Hash) *types.Header
	GetBlock(hash common.Hash, number uint64) *types.Block
}

// Engine is an algorithm agnostic consensus engine.
type Engine interface {
	Author(header *types.Header) (common.Address, error)
	VerifyHeader(chain ChainReader, header *types.Header, seal bool) error
	VerifyHeaders(chain ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error)
	VerifyUncles(chain ChainReader, block *types.Block) error
	VerifySeal(chain ChainReader, header *types.Header) error
	Prepare(chain ChainReader, header *types.Header) error
	Finalize(chain ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt, dposContext *types.DposContext) (*types.Block, error)
	Seal(chain ChainReader, block *types.Block, stop <-chan struct{}) (*types.Block, error)
	APIs(chain ChainReader) []rpc.API
}

// PoW is a consensus engine based on proof-of-work.
type PoW interface {
	Engine

	// Hashrate returns the current mining hashrate of a PoW consensus engine.
	Hashrate() float64
}
