// Copyright 2014 The go-ethereum Authors
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

package types

import (
	"bytes"
	"fmt"
	"io"
	"math/big"

	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/common/hexutil"
	"github.com/boker/go-ethereum/rlp"
)

//go:generate gencodec -type Receipt -field-override receiptMarshaling -out gen_receipt_json.go

var (
	receiptStatusFailedRLP     = []byte{}
	receiptStatusSuccessfulRLP = []byte{0x01}
)

const (
	// ReceiptStatusFailed is the status code of a transaction if execution failed.
	ReceiptStatusFailed = uint(0)

	// ReceiptStatusSuccessful is the status code of a transaction if execution succeeded.
	ReceiptStatusSuccessful = uint(1)
)

//Receipt代表交易的结果
type Receipt struct {

	//共识字段
	PostState         []byte   `json:"root"` //保存了创建该Receipt对象时，整个Block内所有“帐户”的当时状态。Ethereum 里用stateObject来表示一个账户Account，这个账户可转帐(transfer value), 可执行tx, 它的唯一标示符是一个Address类型变量。 这个Receipt.PostState 就是当时所在Block里所有stateObject对象的RLP Hash值
	Status            uint     `json:"status"`
	CumulativeGasUsed *big.Int `json:"cumulativeGasUsed" gencodec:"required"`
	Bloom             Bloom    `json:"logsBloom"         gencodec:"required"` //Ethereum内部实现的一个256bit长Bloom Filter。 Bloom Filter概念定义可见wikipedia，它可用来快速验证一个新收到的对象是否处于一个已知的大量对象集合之中。这里Receipt的Bloom，被用以验证某个给定的Log是否处于Receipt已有的Log数组中
	Logs              []*Log   `json:"logs"              gencodec:"required"` //Log类型的数组，其中每一个Log对象记录了Tx中一小步的操作。所以，每一个tx的执行结果，由一个Receipt对象来表示；更详细的内容，由一组Log对象来记录。这个Log数组很重要，比如在不同Ethereum节点(Node)的相互同步过程中，待同步区块的Log数组有助于验证同步中收到的block是否正确和完整，所以会被单独同步(传输)。

	//实施领域（不要重新订购！）
	TxHash          common.Hash    `json:"transactionHash" gencodec:"required"`
	ContractAddress common.Address `json:"contractAddress"`
	GasUsed         *big.Int       `json:"gasUsed" gencodec:"required"`
}

type receiptMarshaling struct {
	PostState         hexutil.Bytes
	Status            hexutil.Uint
	CumulativeGasUsed *hexutil.Big
	GasUsed           *hexutil.Big
}

// receiptRLP is the consensus encoding of a receipt.
type receiptRLP struct {
	PostStateOrStatus []byte
	CumulativeGasUsed *big.Int
	Bloom             Bloom
	Logs              []*Log
}

type receiptStorageRLP struct {
	PostStateOrStatus []byte
	CumulativeGasUsed *big.Int
	Bloom             Bloom
	TxHash            common.Hash
	ContractAddress   common.Address
	Logs              []*LogForStorage
	GasUsed           *big.Int
}

// NewReceipt creates a barebone transaction receipt, copying the init fields.
func NewReceipt(root []byte, failed bool, cumulativeGasUsed *big.Int) *Receipt {
	r := &Receipt{PostState: common.CopyBytes(root), CumulativeGasUsed: new(big.Int).Set(cumulativeGasUsed)}
	if failed {
		r.Status = ReceiptStatusFailed
	} else {
		r.Status = ReceiptStatusSuccessful
	}
	return r
}

// EncodeRLP implements rlp.Encoder, and flattens the consensus fields of a receipt
// into an RLP stream. If no post state is present, byzantium fork is assumed.
func (r *Receipt) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &receiptRLP{r.statusEncoding(), r.CumulativeGasUsed, r.Bloom, r.Logs})
}

// DecodeRLP implements rlp.Decoder, and loads the consensus fields of a receipt
// from an RLP stream.
func (r *Receipt) DecodeRLP(s *rlp.Stream) error {
	var dec receiptRLP
	if err := s.Decode(&dec); err != nil {
		return err
	}
	if err := r.setStatus(dec.PostStateOrStatus); err != nil {
		return err
	}
	r.CumulativeGasUsed, r.Bloom, r.Logs = dec.CumulativeGasUsed, dec.Bloom, dec.Logs
	return nil
}

func (r *Receipt) setStatus(postStateOrStatus []byte) error {
	switch {
	case bytes.Equal(postStateOrStatus, receiptStatusSuccessfulRLP):
		r.Status = ReceiptStatusSuccessful
	case bytes.Equal(postStateOrStatus, receiptStatusFailedRLP):
		r.Status = ReceiptStatusFailed
	case len(postStateOrStatus) == len(common.Hash{}):
		r.PostState = postStateOrStatus
	default:
		return fmt.Errorf("invalid receipt status %x", postStateOrStatus)
	}
	return nil
}

func (r *Receipt) statusEncoding() []byte {
	if len(r.PostState) == 0 {
		if r.Status == ReceiptStatusFailed {
			return receiptStatusFailedRLP
		}
		return receiptStatusSuccessfulRLP
	}
	return r.PostState
}

// String implements the Stringer interface.
func (r *Receipt) String() string {
	if len(r.PostState) == 0 {
		return fmt.Sprintf("receipt{status=%d cgas=%v bloom=%x logs=%v}", r.Status, r.CumulativeGasUsed, r.Bloom, r.Logs)
	}
	return fmt.Sprintf("receipt{med=%x cgas=%v bloom=%x logs=%v}", r.PostState, r.CumulativeGasUsed, r.Bloom, r.Logs)
}

// ReceiptForStorage is a wrapper around a Receipt that flattens and parses the
// entire content of a receipt, as opposed to only the consensus fields originally.
type ReceiptForStorage Receipt

// EncodeRLP implements rlp.Encoder, and flattens all content fields of a receipt
// into an RLP stream.
func (r *ReceiptForStorage) EncodeRLP(w io.Writer) error {
	enc := &receiptStorageRLP{
		PostStateOrStatus: (*Receipt)(r).statusEncoding(),
		CumulativeGasUsed: r.CumulativeGasUsed,
		Bloom:             r.Bloom,
		TxHash:            r.TxHash,
		ContractAddress:   r.ContractAddress,
		Logs:              make([]*LogForStorage, len(r.Logs)),
		GasUsed:           r.GasUsed,
	}
	for i, log := range r.Logs {
		enc.Logs[i] = (*LogForStorage)(log)
	}
	return rlp.Encode(w, enc)
}

// DecodeRLP implements rlp.Decoder, and loads both consensus and implementation
// fields of a receipt from an RLP stream.
func (r *ReceiptForStorage) DecodeRLP(s *rlp.Stream) error {
	var dec receiptStorageRLP
	if err := s.Decode(&dec); err != nil {
		return err
	}
	if err := (*Receipt)(r).setStatus(dec.PostStateOrStatus); err != nil {
		return err
	}
	// Assign the consensus fields
	r.CumulativeGasUsed, r.Bloom = dec.CumulativeGasUsed, dec.Bloom
	r.Logs = make([]*Log, len(dec.Logs))
	for i, log := range dec.Logs {
		r.Logs[i] = (*Log)(log)
	}
	// Assign the implementation fields
	r.TxHash, r.ContractAddress, r.GasUsed = dec.TxHash, dec.ContractAddress, dec.GasUsed
	return nil
}

// Receipts is a wrapper around a Receipt array to implement DerivableList.
type Receipts []*Receipt

// Len returns the number of receipts in this list.
func (r Receipts) Len() int { return len(r) }

// GetRlp returns the RLP encoding of one receipt from the list.
func (r Receipts) GetRlp(i int) []byte {
	bytes, err := rlp.EncodeToBytes(r[i])
	if err != nil {
		panic(err)
	}
	return bytes
}
