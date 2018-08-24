//播客链主动产生交易
package boker

import (
	"errors"
	"math/big"

	"github.com/boker/go-ethereum/accounts/abi/bind"
	"github.com/boker/go-ethereum/accounts/abi/bind/backends"
	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/core"
	"github.com/boker/go-ethereum/core/types"
	"github.com/boker/go-ethereum/crypto"
	"github.com/boker/go-ethereum/eth"
)

//播客链的基础合约管理
type BokerTransaction struct {
	ethereum *eth.Ethereum
}

func NewTransaction(ethereum *eth.Ethereum) *BokerTransaction {

	bokerTransaction := new(BokerTransaction)
	bokerTransaction.ethereum = ethereum
	return bokerTransaction
}

//产生部署基础合约交易
func (t *BokerTransaction) DeployTransaction(txType types.TxType, address common.Address) error {

	if t.ethereum != nil {

		DeployKey, err := crypto.HexToECDSA(t.ethereum.BlockChain().Config().Producer.PrivateKey)
		if err != nil {
			return err
		}
		DeployAddr := crypto.PubkeyToAddress(DeployKey.PublicKey)
		DeployBalance := big.NewInt(0)
		DeployBalance.SetInt64(t.ethereum.BlockChain().Config().Producer.Balance)

		//构造backend和帐号
		backend := backends.NewSimulatedBackend(core.GenesisAlloc{DeployAddr: {Balance: DeployBalance}}, t.ethereum.Boker)
		opts := bind.NewKeyedTransactor(DeployKey)

		//得到Nonce
		nonce, err := backend.PendingNonceAt(opts.Context, t.ethereum.BlockChain().Config().Coinbase)

		//判断Value值是否为空
		value := opts.Value
		if value == nil {
			value = new(big.Int)
		}
		var input []byte
		rawTx := types.NewBaseTransaction(txType, nonce, address, value, input)

		//判断交易是否有签名者
		if opts.Signer == nil {
			return errors.New("no signer to authorize the transaction with")
		}

		if err := backend.SendTransaction(opts.Context, rawTx); err != nil {
			return err
		}
		return nil
	}
	return nil
}

//产生取消部署基础合约交易
func (t *BokerTransaction) UnDeployTransaction(txType types.TxType, address common.Address) error {

	if t.ethereum != nil {

		DeployKey, err := crypto.HexToECDSA(t.ethereum.BlockChain().Config().Producer.PrivateKey)
		if err != nil {
			return err
		}
		DeployAddr := crypto.PubkeyToAddress(DeployKey.PublicKey)
		DeployBalance := big.NewInt(0)
		DeployBalance.SetInt64(t.ethereum.BlockChain().Config().Producer.Balance)

		//构造backend和帐号
		backend := backends.NewSimulatedBackend(core.GenesisAlloc{DeployAddr: {Balance: DeployBalance}}, t.ethereum.Boker)
		opts := bind.NewKeyedTransactor(DeployKey)

		//得到Nonce
		nonce, err := backend.PendingNonceAt(opts.Context, t.ethereum.BlockChain().Config().Coinbase)

		//判断Value值是否为空
		value := opts.Value
		if value == nil {
			value = new(big.Int)
		}
		var input []byte
		rawTx := types.NewBaseTransaction(txType, nonce, address, value, input)

		//判断交易是否有签名者
		if opts.Signer == nil {
			return errors.New("no signer to authorize the transaction with")
		}

		if err := backend.SendTransaction(opts.Context, rawTx); err != nil {
			return err
		}
		return nil
	}
	return nil
}
