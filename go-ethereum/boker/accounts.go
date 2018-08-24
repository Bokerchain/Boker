//播客链增加的特殊账号管理类
package boker

import (
	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/core/types"
	"github.com/boker/go-ethereum/log"
)

//播客链的账号管理
type BokerAccount struct {
	accounts map[common.Address]types.TxType
}

func NewAccount() *BokerAccount {

	bokerAccount := new(BokerAccount)
	bokerAccount.accounts = make(map[common.Address]types.TxType)
	bokerAccount.loadAccount()
	return bokerAccount
}

func (a *BokerAccount) loadAccount() {

	log.Info("load bokerchain account")

	//部署基础合约帐号
	a.accounts[common.HexToAddress("0x3347cc0f61122bcffb2de5089f6c9c5f968366b2")] = types.DeployVote
	//社区运营基金账户
	a.accounts[common.HexToAddress("0xd7fd311c8f97349670963d87f37a68794dfa80ff")] = types.DeployAssignToken
	//基金会账户
	a.accounts[common.HexToAddress("0xa0da98da40f8c4aba880ad8d219a5c82c8bc97c4")] = types.DeployAssignToken
	//团队账户
	a.accounts[common.HexToAddress("0x81b7fee82a6356351edbf1339a845b2480ad53c2")] = types.DeployAssignToken
}

func (a *BokerAccount) GetAccount(account common.Address) (types.TxType, error) {

	if len(a.accounts) > 0 {
		value, exist := a.accounts[account]
		if exist {
			return value, nil
		}
	}
	return types.Binary, ErrNotFoundAccount
}
