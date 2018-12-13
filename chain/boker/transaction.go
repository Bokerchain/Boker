//播客链主动产生交易
package boker

import (
	"context"
	"math/big"

	"github.com/Bokerchain/Boker/chain/accounts"
	"github.com/Bokerchain/Boker/chain/boker/protocol"
	"github.com/Bokerchain/Boker/chain/common"
	"github.com/Bokerchain/Boker/chain/common/hexutil"
	"github.com/Bokerchain/Boker/chain/core/types"
	"github.com/Bokerchain/Boker/chain/eth"
	"github.com/Bokerchain/Boker/chain/internal/ethapi"
	"github.com/Bokerchain/Boker/chain/log"
)

//播客链的基础合约管理
type BokerTransaction struct {
	ethereum *eth.Ethereum
}

func NewTransaction(ethereum *eth.Ethereum) *BokerTransaction {

	return &BokerTransaction{
		ethereum: ethereum,
	}
}

func (t *BokerTransaction) SubmitBokerTransaction(ctx context.Context, txType protocol.TxType, to common.Address, extra string) error {

	//log.Info("****SubmitBokerTransaction****", "txType", txType, "to", to.String())
	if t.ethereum != nil {

		//得到From账号
		from, err := t.ethereum.ApiBackend.Coinbase()
		if err != nil {
			log.Error("bokerTransaction CoinBase", "error", err)
			return err
		}
		//log.Info("SubmitBokerTransaction CoinBase", "from", from.String())

		//设置参数（其中有些参数可以通过调用设置默认设置来进行获取）
		args := ethapi.SendTxArgs{
			From:     from,
			Type:     txType,
			Nonce:    nil,
			To:       &to,
			Gas:      nil,
			GasPrice: nil,
			Value:    nil,
			Data:     hexutil.Bytes([]byte(extra)),
			//Extra:    hexutil.Bytes([]byte(extra)),
		}

		//查找包含所请求签名者的钱包
		account := accounts.Account{Address: args.From}

		//根据帐号得到钱包信息
		wallet, err := t.ethereum.AccountManager().Find(account)
		if err != nil {
			log.Error("SubmitBokerTransaction AccountManager Find", "error", err)
			return err
		}

		//设置默认设置
		if err := args.SetDefaults(ctx, t.ethereum.ApiBackend); err != nil {
			log.Error("SubmitBokerTransaction SetDefaults", "error", err)
			return err
		}
		log.Info("SubmitBokerTransaction SetDefaults", "Nonce", args.Nonce.String(), "txType", args.Type)

		input := []byte("")
		tx := types.NewBaseTransaction(args.Type, (uint64)(*args.Nonce), (common.Address)(*args.To), (*big.Int)(args.Value), input)

		var chainID *big.Int
		if config := t.ethereum.ApiBackend.ChainConfig(); config.IsEIP155(t.ethereum.ApiBackend.CurrentBlock().Number()) {
			chainID = config.ChainId
		}

		//对该笔交易签名来确保该笔交易的真实有效性
		signed, err := wallet.SignTxWithPassphrase(account, t.ethereum.Password(), tx, chainID)
		if err != nil {
			log.Error("SubmitBokerTransaction SignTxWithPassphrase", "error", err)
			return err
		}
		//log.Info("SubmitBokerTransaction SetDefaults", "Nonce", args.Nonce.String(), "GasPrice", args.GasPrice)

		if _, err := ethapi.SubmitTransaction(ctx, t.ethereum.ApiBackend, signed); err != nil {
			log.Error("SubmitBokerTransaction SubmitTransaction", "error", err)
			return err
		}

		return nil
	}
	return protocol.ErrInvalidSystem
}
