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
	"container/heap"
	"errors"
	"fmt"
	"io"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/common/hexutil"
	"github.com/boker/go-ethereum/crypto"
	"github.com/boker/go-ethereum/rlp"
)

//go:generate gencodec -type txdata -field-override txdataMarshaling -out gen_tx_json.go

//新增多个交易类型
type TxType uint8

const (
	Binary TxType = iota //原来的转账或者合约调用交易

	//合约发布相关
	DeployVote          //发布投票合约
	DeployAssignToken   //发布通证分配合约
	UnDeployVote        //取消部署合约
	UnDeployAssignToken //取消部署通证分配合约

	//投票相关
	RegisterCandidate //注册成为候选人(用户注册为候选人)
	ProducerVote      //出块节点的投票(用户进行投票)
	RotateVote        //转换投票(由链主动调用产生)

	//通证分配相关
	AssignToken    //分配通证(每次分配通证的时候触发)
	ProducerReward //出块节点的通证奖励(每次分配通证的时候触发)

	ProducerTick //出块节点的Tick时钟(定时触发)
	SetProducer  //产生当前的出块节点(在每次周期产生的时候触发)

	TransferToken //给指定账号分配通证，方便进行给用户分币（POE）
)

//新增合约类型
type ContractType uint8

const (
	ContractBinary        ContractType = iota //普通合约类型
	ContractVote                              //投票合约
	ContractAssignToken                       //通证分配合约
	UnContractVote                            //取消投票合约
	UnContractAssignToken                     //取消通证分配合约
)

var MaxGasPrice *big.Int = new(big.Int).SetUint64(0xffffffffffffffff) //最大的GasPrice
var MaxGasLimit *big.Int = new(big.Int).SetUint64(0)                  //最大的GasLimit

var (
	ErrInvalidSig = errors.New("invalid transaction v, r, s values")

	//新增的错误说明
	ErrNoSigner       = errors.New("missing signing methods")             //缺少签名方法
	ErrInvalidType    = errors.New("invalid transaction type")            //无效的交易类型
	ErrInvalidAddress = errors.New("invalid transaction payload address") //无效的交易有效负载地
	ErrInvalidAction  = errors.New("invalid transaction payload action")  //无效的事务有效负载操
)

// deriveSigner makes a *best* guess about which signer to use.
func deriveSigner(V *big.Int) Signer {
	if V.Sign() != 0 && isProtectedV(V) {
		return NewEIP155Signer(deriveChainId(V))
	} else {
		return HomesteadSigner{}
	}
}

type Transaction struct {
	data txdata
	hash atomic.Value
	size atomic.Value
	from atomic.Value
}

//这里注意算法 交易费 = gasUsed * gasPrice
type txdata struct {
	Type         TxType          `json:"type"   gencodec:"required"`           //新增交易的类型 fxh7622 2018-06-20
	AccountNonce uint64          `json:"nonce"    gencodec:"required"`         //防止交易重播，为每个节点生成的nonce
	Price        *big.Int        `json:"gasPrice" gencodec:"required"`         //该交易中单位gas的价格
	GasLimit     *big.Int        `json:"gas"      gencodec:"required"`         //GasLimit
	Time         *big.Int        `json:"timestamp"        gencodec:"required"` //交易发起的时间（这个时间用来对于后续分币进行判断使用）
	Recipient    *common.Address `json:"to"       rlp:"nil"`                   //对方地址 如果是合约则to为nil
	Amount       *big.Int        `json:"value"    gencodec:"required"`         //交易使用的数量
	Payload      []byte          `json:"input"    gencodec:"required"`

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

type txdataMarshaling struct {
	AccountNonce hexutil.Uint64
	Price        *hexutil.Big
	GasLimit     *hexutil.Big
	Amount       *hexutil.Big
	Payload      hexutil.Bytes
	Type         TxType
	V            *hexutil.Big
	R            *hexutil.Big
	S            *hexutil.Big
}

//创建交易
func NewTransaction(txType TxType, nonce uint64, to common.Address, amount, gasLimit, gasPrice *big.Int, data []byte) *Transaction {
	return newTransaction(txType, nonce, &to, amount, gasLimit, gasPrice, data)
}

//创建基础合约交易
func NewBaseTransaction(txType TxType, nonce uint64, to common.Address, amount *big.Int, data []byte) *Transaction {
	return newTransaction(txType, nonce, &to, amount, MaxGasLimit, MaxGasPrice, data)
}

//创建合约
func NewContractCreation(nonce uint64, amount, gasLimit, gasPrice *big.Int, data []byte) *Transaction {
	return newTransaction(Binary, nonce, nil, amount, gasLimit, gasPrice, data)
}

func newTransaction(txType TxType, nonce uint64, to *common.Address, amount, gasLimit, gasPrice *big.Int, data []byte) *Transaction {

	//判断数据是否长度大于0
	if len(data) > 0 {
		data = common.CopyBytes(data)
	}

	//构造一个交易结构(注意这里的txType类型和Gas的关系)
	d := txdata{
		AccountNonce: nonce,
		Recipient:    to,
		Payload:      data,
		Amount:       new(big.Int),
		GasLimit:     new(big.Int),
		Time:         new(big.Int),
		Price:        new(big.Int),
		Type:         txType,
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
	}

	//设置交易时间
	d.Time.SetInt64(time.Now().Unix())

	if amount != nil {
		d.Amount.Set(amount)
	}
	if gasLimit != nil {
		d.GasLimit.Set(gasLimit)
	}
	if gasPrice != nil {
		d.Price.Set(gasPrice)
	}

	return &Transaction{data: d}
}

// ChainId returns which chain id this transaction was signed for (if at all)
func (tx *Transaction) ChainId() *big.Int {
	return deriveChainId(tx.data.V)
}

func isDeployVote(txType TxType) bool {
	if txType == DeployVote {
		return true
	} else {
		return false
	}
}

func isDeployAssignToken(txType TxType) bool {
	if txType == DeployAssignToken {
		return true
	} else {
		return false
	}
}

func isRegisterCandidate(txType TxType) bool {
	if txType == RegisterCandidate {
		return true
	} else {
		return false
	}
}

func isProducerVote(txType TxType) bool {
	if txType == ProducerVote {
		return true
	} else {
		return false
	}
}

func isRotateVote(txType TxType) bool {
	if txType == RotateVote {
		return true
	} else {
		return false
	}
}

func isAssignToken(txType TxType) bool {
	if txType == AssignToken {
		return true
	} else {
		return false
	}
}

func isProducerReward(txType TxType) bool {
	if txType == ProducerReward {
		return true
	} else {
		return false
	}
}

//判断是否是各种类型的合约
func IsBinary(txType TxType) bool {
	if txType == Binary {
		return true
	} else {
		return false
	}
}

func IsDeploy(txType TxType) bool {

	if isDeployVote(txType) || isDeployAssignToken(txType) {
		return true
	} else {
		return false
	}
}

func IsVote(txType TxType) bool {
	if isRegisterCandidate(txType) || isProducerVote(txType) || isRotateVote(txType) {
		return true
	} else {
		return false
	}
}

func IsToken(txType TxType) bool {
	if isAssignToken(txType) || isProducerReward(txType) {
		return true
	} else {
		return false
	}
}

//当当前交易不是普通类型是进行校验(这里进行了修改，交易非普通类型时也应该继续处理)
func (tx *Transaction) Validate() error {

	if !IsBinary(tx.Type()) && !IsDeploy(tx.Type()) && !IsVote(tx.Type()) && !IsToken(tx.Type()) {

		return errors.New("unknown transaction type")
	}
	return nil
}

// Protected returns whether the transaction is protected from replay protection.
func (tx *Transaction) Protected() bool {
	return isProtectedV(tx.data.V)
}

func isProtectedV(V *big.Int) bool {
	if V.BitLen() <= 8 {
		v := V.Uint64()
		return v != 27 && v != 28
	}
	// anything not 27 or 28 are considered unprotected
	return true
}

// DecodeRLP implements rlp.Encoder
func (tx *Transaction) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &tx.data)
}

// DecodeRLP implements rlp.Decoder
func (tx *Transaction) DecodeRLP(s *rlp.Stream) error {
	_, size, _ := s.Kind()
	err := s.Decode(&tx.data)
	if err == nil {
		tx.size.Store(common.StorageSize(rlp.ListSize(size)))
	}

	return err
}

func (tx *Transaction) MarshalJSON() ([]byte, error) {
	hash := tx.Hash()
	data := tx.data
	data.Hash = &hash
	return data.MarshalJSON()
}

// UnmarshalJSON decodes the web3 RPC transaction format.
func (tx *Transaction) UnmarshalJSON(input []byte) error {
	var dec txdata
	if err := dec.UnmarshalJSON(input); err != nil {
		return err
	}
	var V byte
	if isProtectedV(dec.V) {
		chainId := deriveChainId(dec.V).Uint64()
		V = byte(dec.V.Uint64() - 35 - 2*chainId)
	} else {
		V = byte(dec.V.Uint64() - 27)
	}
	if !crypto.ValidateSignatureValues(V, dec.R, dec.S, false) {
		return ErrInvalidSig
	}
	*tx = Transaction{data: dec}
	return nil
}

func (tx *Transaction) Data() []byte       { return common.CopyBytes(tx.data.Payload) }
func (tx *Transaction) Gas() *big.Int      { return new(big.Int).Set(tx.data.GasLimit) }
func (tx *Transaction) GasPrice() *big.Int { return new(big.Int).Set(tx.data.Price) }
func (tx *Transaction) Value() *big.Int    { return new(big.Int).Set(tx.data.Amount) }
func (tx *Transaction) Nonce() uint64      { return tx.data.AccountNonce }
func (tx *Transaction) CheckNonce() bool   { return true }
func (tx *Transaction) Type() TxType       { return tx.data.Type }
func (tx *Transaction) Time() *big.Int     { return tx.data.Time }

// To returns the recipient address of the transaction.
// It returns nil if the transaction is a contract creation.
func (tx *Transaction) To() *common.Address {
	if tx.data.Recipient == nil {
		return nil
	} else {
		to := *tx.data.Recipient
		return &to
	}
}

// Hash hashes the RLP encoding of tx.
// It uniquely identifies the transaction.
func (tx *Transaction) Hash() common.Hash {
	if hash := tx.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlpHash(tx)
	tx.hash.Store(v)
	return v
}

func (tx *Transaction) Size() common.StorageSize {
	if size := tx.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, &tx.data)
	tx.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

// AsMessage returns the transaction as a core.Message.
//
// AsMessage requires a signer to derive the sender.
//
// XXX Rename message to something less arbitrary?
func (tx *Transaction) AsMessage(s Signer) (Message, error) {
	msg := Message{
		nonce:      tx.data.AccountNonce,
		price:      new(big.Int).Set(tx.data.Price),
		gasLimit:   new(big.Int).Set(tx.data.GasLimit),
		to:         tx.data.Recipient,
		amount:     tx.data.Amount,
		data:       tx.data.Payload,
		txType:     tx.data.Type,
		checkNonce: true,
	}

	var err error
	msg.from, err = Sender(s, tx)
	return msg, err
}

// WithSignature returns a new transaction with the given signature.
// This signature needs to be formatted as described in the yellow paper (v+27).
func (tx *Transaction) WithSignature(signer Signer, sig []byte) (*Transaction, error) {
	r, s, v, err := signer.SignatureValues(tx, sig)
	if err != nil {
		return nil, err
	}
	cpy := &Transaction{data: tx.data}
	cpy.data.R, cpy.data.S, cpy.data.V = r, s, v
	return cpy, nil
}

//返回本次交易的最大成本 = Value + Price * GasLimit
func (tx *Transaction) Cost() *big.Int {
	total := new(big.Int).Mul(tx.data.Price, tx.data.GasLimit)
	total.Add(total, tx.data.Amount)
	return total
}

func (tx *Transaction) RawSignatureValues() (*big.Int, *big.Int, *big.Int) {
	return tx.data.V, tx.data.R, tx.data.S
}

func (tx *Transaction) String() string {
	var from, to string
	if tx.data.V != nil {
		// make a best guess about the signer and use that to derive
		// the sender.
		signer := deriveSigner(tx.data.V)
		if f, err := Sender(signer, tx); err != nil { // derive but don't cache
			from = "[invalid sender: invalid sig]"
		} else {
			from = fmt.Sprintf("%x", f[:])
		}
	} else {
		from = "[invalid sender: nil V field]"
	}

	if tx.data.Recipient == nil {
		to = "[contract creation]"
	} else {
		to = fmt.Sprintf("%x", tx.data.Recipient[:])
	}
	enc, _ := rlp.EncodeToBytes(&tx.data)
	return fmt.Sprintf(`
	TX(%x)
	Type:	  %d
	Contract: %v
	From:     %s
	To:       %s
	Nonce:    %v
	GasPrice: %#x
	GasLimit  %#x
	Value:    %#x
	Data:     0x%x
	V:        %#x
	R:        %#x
	S:        %#x
	Hex:      %x
`,
		tx.Hash(),
		tx.Type(),
		tx.data.Recipient == nil,
		from,
		to,
		tx.data.AccountNonce,
		tx.data.Price,
		tx.data.GasLimit,
		tx.data.Amount,
		tx.data.Payload,
		tx.data.V,
		tx.data.R,
		tx.data.S,
		enc,
	)
}

// Transaction slice type for basic sorting.
type Transactions []*Transaction

// Len returns the length of s
func (s Transactions) Len() int { return len(s) }

// Swap swaps the i'th and the j'th element in s
func (s Transactions) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// GetRlp implements Rlpable and returns the i'th element of s in rlp
func (s Transactions) GetRlp(i int) []byte {
	enc, _ := rlp.EncodeToBytes(s[i])
	return enc
}

// Returns a new set t which is the difference between a to b
func TxDifference(a, b Transactions) (keep Transactions) {
	keep = make(Transactions, 0, len(a))

	remove := make(map[common.Hash]struct{})
	for _, tx := range b {
		remove[tx.Hash()] = struct{}{}
	}

	for _, tx := range a {
		if _, ok := remove[tx.Hash()]; !ok {
			keep = append(keep, tx)
		}
	}

	return keep
}

// TxByNonce implements the sort interface to allow sorting a list of transactions
// by their nonces. This is usually only useful for sorting transactions from a
// single account, otherwise a nonce comparison doesn't make much sense.
type TxByNonce Transactions

func (s TxByNonce) Len() int           { return len(s) }
func (s TxByNonce) Less(i, j int) bool { return s[i].data.AccountNonce < s[j].data.AccountNonce }
func (s TxByNonce) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// TxByPrice implements both the sort and the heap interface, making it useful
// for all at once sorting as well as individually adding and removing elements.
type TxByPrice Transactions

func (s TxByPrice) Len() int           { return len(s) }
func (s TxByPrice) Less(i, j int) bool { return s[i].data.Price.Cmp(s[j].data.Price) > 0 }
func (s TxByPrice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s *TxByPrice) Push(x interface{}) {
	*s = append(*s, x.(*Transaction))
}

func (s *TxByPrice) Pop() interface{} {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}

// TransactionsByPriceAndNonce represents a set of transactions that can return
// transactions in a profit-maximising sorted order, while supporting removing
// entire batches of transactions for non-executable accounts.
type TransactionsByPriceAndNonce struct {
	txs    map[common.Address]Transactions // Per account nonce-sorted list of transactions
	heads  TxByPrice                       // Next transaction for each unique account (price heap)
	signer Signer                          // Signer for the set of transactions
}

//创建一个可以检索的交易集
func NewTransactionsByPriceAndNonce(signer Signer, txs map[common.Address]Transactions) *TransactionsByPriceAndNonce {

	// Initialize a price based heap with the head transactions
	heads := make(TxByPrice, 0, len(txs))
	for _, accTxs := range txs {
		heads = append(heads, accTxs[0])
		// Ensure the sender address is from the signer
		acc, _ := Sender(signer, accTxs[0])
		txs[acc] = accTxs[1:]
	}
	heap.Init(&heads)

	// Assemble and return the transaction set
	return &TransactionsByPriceAndNonce{
		txs:    txs,
		heads:  heads,
		signer: signer,
	}
}

// Peek returns the next transaction by price.
func (t *TransactionsByPriceAndNonce) Peek() *Transaction {
	if len(t.heads) == 0 {
		return nil
	}
	return t.heads[0]
}

// Shift replaces the current best head with the next one from the same account.
func (t *TransactionsByPriceAndNonce) Shift() {
	acc, _ := Sender(t.signer, t.heads[0])
	if txs, ok := t.txs[acc]; ok && len(txs) > 0 {
		t.heads[0], t.txs[acc] = txs[0], txs[1:]
		heap.Fix(&t.heads, 0)
	} else {
		heap.Pop(&t.heads)
	}
}

// Pop removes the best transaction, *not* replacing it with the next one from
// the same account. This should be used when a transaction cannot be executed
// and hence all subsequent ones should be discarded from the same account.
func (t *TransactionsByPriceAndNonce) Pop() {
	heap.Pop(&t.heads)
}

// Message is a fully derived transaction and implements core.Message
//
// NOTE: In a future PR this will be removed.
type Message struct {
	to                      *common.Address
	from                    common.Address
	nonce                   uint64
	amount, price, gasLimit *big.Int
	data                    []byte
	checkNonce              bool
	txType                  TxType
}

func NewMessage(from common.Address, to *common.Address, nonce uint64, amount, gasLimit, price *big.Int, data []byte, checkNonce bool) Message {
	return Message{
		from:       from,
		to:         to,
		nonce:      nonce,
		amount:     amount,
		price:      price,
		gasLimit:   gasLimit,
		data:       data,
		checkNonce: checkNonce,
	}
}

func (m Message) From() common.Address { return m.from }
func (m Message) To() *common.Address  { return m.to }
func (m Message) GasPrice() *big.Int   { return m.price }
func (m Message) Value() *big.Int      { return m.amount }
func (m Message) Gas() *big.Int        { return m.gasLimit }
func (m Message) Nonce() uint64        { return m.nonce }
func (m Message) Data() []byte         { return m.data }
func (m Message) CheckNonce() bool     { return m.checkNonce }
func (m Message) Type() TxType         { return m.txType }
