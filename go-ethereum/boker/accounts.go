//播客链增加的特殊账号管理类
package boker

import (
	"github.com/boker/go-ethereum/boker/protocol"
	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/log"
)

//播客链的账号管理
type AcccountLevel struct {
	level []protocol.TxType
}
type BokerAccount struct {
	accounts map[common.Address]AcccountLevel
}

func NewAccount() *BokerAccount {

	bokerAccount := new(BokerAccount)
	bokerAccount.accounts = make(map[common.Address]AcccountLevel)
	bokerAccount.loadAccount()
	return bokerAccount
}

func (a *BokerAccount) loadVoteDeployAccount() {

	log.Info("****loadVoteDeployAccount****")

	accountLevel := AcccountLevel{
		level: make([]protocol.TxType, 0, 0),
	}
	accountLevel.level = append(accountLevel.level, protocol.SetVote)
	accountLevel.level = append(accountLevel.level, protocol.CancelVote)
	a.accounts[common.HexToAddress("0xd7fd311c8f97349670963d87f37a68794dfa80ff")] = accountLevel
}

func (a *BokerAccount) loadTokenDeployAccount() {

	log.Info("****loadTokenDeployAccount****")

	accountLevel := AcccountLevel{
		level: make([]protocol.TxType, 0, 0),
	}
	accountLevel.level = append(accountLevel.level, protocol.SetAssignToken)
	accountLevel.level = append(accountLevel.level, protocol.CanclAssignToken)
	a.accounts[common.HexToAddress("0xd7fd311c8f97349670963d87f37a68794dfa80ff")] = accountLevel
}

func (a *BokerAccount) loadSetValidator() {

	log.Info("****loadSetValidator****")

	//加载添加验证人账号
	accountLevel := AcccountLevel{
		level: make([]protocol.TxType, 0, 0),
	}
	accountLevel.level = append(accountLevel.level, protocol.SetValidator)
	a.accounts[common.HexToAddress("0x97da0c2f933ff6aad55a0b9eb1933f5b0ae3cd9b")] = accountLevel
}

func (a *BokerAccount) IsValidator(address common.Address) bool {

	if len(a.accounts) <= 0 {
		return false
	}

	levels := a.accounts[address]
	if len(levels.level) <= 0 {
		return false
	}

	for _, v := range levels.level {
		if v == protocol.SetValidator {
			return true
		}
	}
	return false
}

func (a *BokerAccount) loadCommunityAccount() {

	log.Info("****loadCommunityAccount****")

	//社区运营基金账户
	communityOperationsAccount := AcccountLevel{
		level: make([]protocol.TxType, 0, 0),
	}
	communityOperationsAccount.level = append(communityOperationsAccount.level, protocol.SetValidator)
	a.accounts[common.HexToAddress("0xd7fd311c8f97349670963d87f37a68794dfa80ff")] = communityOperationsAccount

	//基金会账户
	foundationAccount := AcccountLevel{
		level: make([]protocol.TxType, 0, 0),
	}
	foundationAccount.level = append(foundationAccount.level, protocol.SetValidator)
	a.accounts[common.HexToAddress("0xd7fd311c8f97349670963d87f37a68794dfa80ff")] = foundationAccount

	//团队账户
	teamAccount := AcccountLevel{
		level: make([]protocol.TxType, 0, 0),
	}
	teamAccount.level = append(teamAccount.level, protocol.SetValidator)
	a.accounts[common.HexToAddress("0xd7fd311c8f97349670963d87f37a68794dfa80ff")] = teamAccount
}

func (a *BokerAccount) loadAccount() {

	log.Info("****loadAccount****")
	a.loadSetValidator()
}

func (a *BokerAccount) GetAccount(account common.Address) ([]protocol.TxType, error) {

	if len(a.accounts) > 0 {
		value, exist := a.accounts[account]
		if exist {
			return value.level, nil
		}
	}

	//测试使用
	return []protocol.TxType{protocol.SetValidator}, nil

	//return []protocol.TxType{protocol.Binary}, protocol.ErrSpecialAccount
}
