//播客链增加的特殊账号管理类
package boker

import (
	"github.com/boker/chain/boker/protocol"
	"github.com/boker/chain/common"
	"github.com/boker/chain/log"
)

var (
	//发布投票合约账户
	VoteDeployAddress = common.HexToAddress("0xd7fd311c8f97349670963d87f37a68794dfa80ff")
	//发布分布合约账户
	TokenDeployAddress = common.HexToAddress("0xd7fd311c8f97349670963d87f37a68794dfa80ff")
	//设置验证人账户
	SetValidatorAddress = common.HexToAddress("0x97da0c2f933ff6aad55a0b9eb1933f5b0ae3cd9b")
	//社区账户
	CommunityAddress = common.HexToAddress("0xd7fd311c8f97349670963d87f37a68794dfa80ff")
	//基金账户
	FoundationAddress = common.HexToAddress("0xd7fd311c8f97349670963d87f37a68794dfa80ff")
	//团队账户
	TeamAddress = common.HexToAddress("0xd7fd311c8f97349670963d87f37a68794dfa80ff")
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

	//加载发布投票合约账户
	log.Info("loadVoteDeployAccount")
	bokerAccount.loadVoteDeployAccount()

	//加载发布分币合约账户
	log.Info("loadTokenDeployAccount")
	bokerAccount.loadTokenDeployAccount()

	//加载设置验证人账户
	log.Info("loadSetValidator")
	bokerAccount.loadSetValidator()

	//加载基金会账户
	log.Info("loadCommunityAccount")
	bokerAccount.loadCommunityAccount()

	return bokerAccount
}

func (a *BokerAccount) loadVoteDeployAccount() {

	accountLevel := AcccountLevel{
		level: make([]protocol.TxType, 0, 0),
	}
	accountLevel.level = append(accountLevel.level, protocol.SetVote)
	accountLevel.level = append(accountLevel.level, protocol.CancelVote)
	a.accounts[VoteDeployAddress] = accountLevel
}

func (a *BokerAccount) loadTokenDeployAccount() {

	accountLevel := AcccountLevel{
		level: make([]protocol.TxType, 0, 0),
	}
	accountLevel.level = append(accountLevel.level, protocol.SetAssignToken)
	accountLevel.level = append(accountLevel.level, protocol.CanclAssignToken)
	a.accounts[TokenDeployAddress] = accountLevel
}

func (a *BokerAccount) loadSetValidator() {

	//加载添加验证人账号
	accountLevel := AcccountLevel{
		level: make([]protocol.TxType, 0, 0),
	}
	accountLevel.level = append(accountLevel.level, protocol.SetValidator)
	a.accounts[SetValidatorAddress] = accountLevel
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

	//社区运营基金账户
	communityOperationsAccount := AcccountLevel{
		level: make([]protocol.TxType, 0, 0),
	}
	communityOperationsAccount.level = append(communityOperationsAccount.level, protocol.SetValidator)
	a.accounts[CommunityAddress] = communityOperationsAccount

	//基金会账户
	foundationAccount := AcccountLevel{
		level: make([]protocol.TxType, 0, 0),
	}
	foundationAccount.level = append(foundationAccount.level, protocol.SetValidator)
	a.accounts[FoundationAddress] = foundationAccount

	//团队账户
	teamAccount := AcccountLevel{
		level: make([]protocol.TxType, 0, 0),
	}
	teamAccount.level = append(teamAccount.level, protocol.SetValidator)
	a.accounts[TeamAddress] = teamAccount
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
