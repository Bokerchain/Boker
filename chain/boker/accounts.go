//播客链增加的特殊账号管理类
package boker

import (
	"github.com/Bokerchain/Boker/chain/boker/protocol"
	"github.com/Bokerchain/Boker/chain/common"
	"github.com/Bokerchain/Boker/chain/log"
)

//帐号常量
const DeloyAddress = "0xd7fd311c8f97349670963d87f37a68794dfa80ff"
const SystemAddress = "0xd7fd311c8f97349670963d87f37a68794dfa80ff"

var (
	DeployPersonalContractAddress = common.HexToAddress(DeloyAddress)  //部署用户基础合约帐号
	DeploySystemContractAddress   = common.HexToAddress(DeloyAddress)  //部署系统基础合约帐号
	SetValidatorAddress           = common.HexToAddress(SystemAddress) //设置验证人账户
	CommunityAddress              = common.HexToAddress(SystemAddress) //社区账户
	FoundationAddress             = common.HexToAddress(SystemAddress) //基金账户
	TeamAddress                   = common.HexToAddress(SystemAddress) //团队账户
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
	log.Info("deployPersonalContractAccount")
	bokerAccount.deployPersonalContractAccount()

	//加载发布分币合约账户
	log.Info("deploySystemContractAccount")
	bokerAccount.deploySystemContractAccount()

	//加载设置验证人账户
	log.Info("loadSetValidator")
	bokerAccount.loadSetValidator()

	//加载基金会账户
	log.Info("loadCommunityAccount")
	bokerAccount.loadCommunityAccount()

	return bokerAccount
}

func (a *BokerAccount) deployPersonalContractAccount() {

	accountLevel := AcccountLevel{
		level: make([]protocol.TxType, 0, 0),
	}
	accountLevel.level = append(accountLevel.level, protocol.SetPersonalContract)
	accountLevel.level = append(accountLevel.level, protocol.CancelPersonalContract)
	a.accounts[DeployPersonalContractAddress] = accountLevel
}

func (a *BokerAccount) deploySystemContractAccount() {

	accountLevel := AcccountLevel{
		level: make([]protocol.TxType, 0, 0),
	}
	accountLevel.level = append(accountLevel.level, protocol.SetSystemContract)
	accountLevel.level = append(accountLevel.level, protocol.CancelSystemContract)
	a.accounts[DeploySystemContractAddress] = accountLevel
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
