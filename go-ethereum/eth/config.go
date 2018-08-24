// Copyright 2017 The go-ethereum Authors
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

package eth

import (
	"math/big"
	"os"
	"os/user"

	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/common/hexutil"
	"github.com/boker/go-ethereum/core"
	"github.com/boker/go-ethereum/eth/downloader"
	"github.com/boker/go-ethereum/eth/gasprice"
	"github.com/boker/go-ethereum/params"
)

// DefaultConfig contains default settings for use on the Ethereum main net.
var DefaultConfig = Config{
	SyncMode:      downloader.FullSync,
	NetworkId:     1357,
	LightPeers:    20,
	DatabaseCache: 128,
	GasPrice:      big.NewInt(18 * params.Shannon),
	TxPool:        core.DefaultTxPoolConfig,
	GPO: gasprice.Config{
		Blocks:     10,
		Percentile: 50,
	},
}

func init() {
	home := os.Getenv("HOME")
	if home == "" {
		if user, err := user.Current(); err == nil {
			home = user.HomeDir
		}
	}
}

//go:generate gencodec -type Config -field-override configMarshaling -formats toml -out gen_config.go

type Config struct {
	Genesis                 *core.Genesis       `toml:",omitempty"` //genesis块，如果数据库为空则插入。如果为nil，则使用以太坊主网块。
	NetworkId               uint64              //用于选择要连接的其它节点的网络ID
	SyncMode                downloader.SyncMode //是否同步模式
	LightServ               int                 `toml:",omitempty"` // Maximum percentage of time allowed for serving LES requests
	LightPeers              int                 `toml:",omitempty"` // Maximum number of LES client peers
	SkipBcVersionCheck      bool                `toml:"-"`
	DatabaseHandles         int                 `toml:"-"`
	DatabaseCache           int
	Coinbase                common.Address    `toml:",omitempty"` //矿工账号
	MinerThreads            int               `toml:",omitempty"` //挖矿线程数量
	ExtraData               []byte            `toml:",omitempty"` //扩展字段
	GasPrice                *big.Int          //交易价格
	TxPool                  core.TxPoolConfig //交易池配置
	GPO                     gasprice.Config   //Gas配置
	EnablePreimageRecording bool              //是否允许跟踪VM中的SHA3 preimages
	DocRoot                 string            `toml:"-"`
	PowFake                 bool              `toml:"-"`
	PowTest                 bool              `toml:"-"`
	PowShared               bool              `toml:"-"`
	Dpos                    bool              `toml:"-"`
}

type configMarshaling struct {
	ExtraData hexutil.Bytes
}
