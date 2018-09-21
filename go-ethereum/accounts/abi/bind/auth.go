// Copyright 2016 The go-ethereum Authors
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

package bind

import (
	"crypto/ecdsa"
	"errors"
	"io"
	"io/ioutil"
	"math/big"

	"github.com/boker/go-ethereum/accounts"
	"github.com/boker/go-ethereum/accounts/keystore"
	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/core/types"
	"github.com/boker/go-ethereum/crypto"
	"github.com/boker/go-ethereum/eth"
	"github.com/boker/go-ethereum/log"
)

// NewTransactor is a utility method to easily create a transaction signer from
// an encrypted json key stream and the associated passphrase.
func NewTransactor(keyin io.Reader, passphrase string) (*TransactOpts, error) {
	json, err := ioutil.ReadAll(keyin)
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(json, passphrase)
	if err != nil {
		return nil, err
	}
	return NewKeyedTransactor(key.PrivateKey), nil
}

// NewKeyedTransactor is a utility method to easily create a transaction signer
// from a single private key.
func NewKeyedTransactor(key *ecdsa.PrivateKey) *TransactOpts {
	keyAddr := crypto.PubkeyToAddress(key.PublicKey)
	return &TransactOpts{
		From: keyAddr,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != keyAddr {
				return nil, errors.New("not authorized to sign this account")
			}
			signature, err := crypto.Sign(signer.Hash(tx).Bytes(), key)
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}
}

//播客链中使用密码创建的Opts
func NewPasswordTransactor(ethereum *eth.Ethereum, addr common.Address) *TransactOpts {

	keyAddr := addr
	return &TransactOpts{
		From: keyAddr,

		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {

			if address != keyAddr {
				return nil, errors.New("not authorized to sign this account")
			}

			var chainID *big.Int
			if config := ethereum.ApiBackend.ChainConfig(); config.IsEIP155(ethereum.ApiBackend.CurrentBlock().Number()) {
				chainID = config.ChainId
			}

			account := accounts.Account{Address: address}
			wallet, err := ethereum.AccountManager().Find(account)
			if err != nil {
				log.Error("SubmitBokerTransaction AccountManager Find", "error", err)
				return nil, err
			}

			return wallet.SignTxWithPassphrase(account, ethereum.Password(), tx, chainID)
		},
	}
}
