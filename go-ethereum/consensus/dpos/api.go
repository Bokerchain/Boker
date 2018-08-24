package dpos

import (
	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/consensus"
	"github.com/boker/go-ethereum/core/types"
	"github.com/boker/go-ethereum/rpc"

	"math/big"
)

// API is a user facing RPC API to allow controlling the delegate and voting
// mechanisms of the delegated-proof-of-stake
type API struct {
	chain consensus.ChainReader
	dpos  *Dpos
}

// GetValidators retrieves the list of the validators at specified block
func (api *API) GetValidators(number *rpc.BlockNumber) ([]common.Address, error) {
	var header *types.Header
	if number == nil || *number == rpc.LatestBlockNumber {
		header = api.chain.CurrentHeader()
	} else {
		header = api.chain.GetHeaderByNumber(uint64(number.Int64()))
	}
	if header == nil {
		return nil, errUnknownBlock
	}

	epochTrie, err := types.NewEpochTrie(header.DposContext.EpochHash, api.dpos.db)
	if err != nil {
		return nil, err
	}
	dposContext := types.DposContext{}
	dposContext.SetEpoch(epochTrie)
	validators, err := dposContext.GetEpochTrie()
	if err != nil {
		return nil, err
	}
	return validators, nil
}

// GetConfirmedBlockNumber retrieves the latest irreversible block
func (api *API) GetConfirmedBlockNumber() (*big.Int, error) {
	var err error
	header := api.dpos.confirmedBlockHeader
	if header == nil {
		header, err = api.dpos.loadConfirmedBlockHeader(api.chain)
		if err != nil {
			return nil, err
		}
	}
	return header.Number, nil
}
