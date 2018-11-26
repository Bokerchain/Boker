package bokerapi

import (
	"context"
	"math/big"

	"github.com/boker/chain/boker/protocol"
	"github.com/boker/chain/common"
	"github.com/boker/chain/trie"
)

//播客链定义的接口
type Api interface {
	GetAccount(account common.Address) ([]protocol.TxType, error)                                                //得到账号级别
	GetContract(address common.Address) (protocol.ContractType, error)                                           //得到合约级别
	SetContract(address common.Address, contractType protocol.ContractType, isCancel bool, abiJson string) error //设置合约级别
	GetContractAddr(protocol.ContractType) (common.Address, error)                                               //得到合约帐号
	SubmitBokerTransaction(ctx context.Context, txType protocol.TxType, to common.Address, extra string) error   //产生一个设置验证者交易
	IsValidator(address common.Address) bool                                                                     //必须是特殊账号
	GetContractTrie() (*trie.Trie, *trie.Trie, *trie.Trie)                                                       //得到合约树
	GetMethodName(txType protocol.TxType) (string, string, error)                                                //根据交易类型得到方法名称（只适用于基础合约）
}

func ExistsTxType(txType protocol.TxType, txTypes []protocol.TxType) bool {

	if len(txTypes) <= 0 {
		return false
	}

	for _, v := range txTypes {
		if v == txType {
			return true
		}
	}
	return false
}

type SortableAddress struct {
	Address common.Address
	Weight  *big.Int
}
type SortableAddresses []*SortableAddress

func (p SortableAddresses) Swap(i, j int) { p[i], p[j] = p[j], p[i] }
func (p SortableAddresses) Len() int      { return len(p) }
func (p SortableAddresses) Less(i, j int) bool {
	if p[i].Weight.Cmp(p[j].Weight) < 0 {
		return false
	} else if p[i].Weight.Cmp(p[j].Weight) > 0 {
		return true
	} else {
		return p[i].Address.String() < p[j].Address.String()
	}
}
