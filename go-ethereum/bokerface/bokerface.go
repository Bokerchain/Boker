package bokerface

import (
	"math/big"

	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/core/types"
)

//播客链定义的接口
type BokerInterface interface {
	GetAccount(account common.Address) (types.TxType, error)                   //得到账号级别
	GetContract(address common.Address) (types.ContractType, error)            //得到合约级别
	SetContract(address common.Address, contractType types.ContractType) error //设置合约级别
	SetDeployTransaction(txType types.TxType, address common.Address) error    //产生一个部署合约交易
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
