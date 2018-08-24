# 播客链官方

	* [播客链官网](http://www.videoblkchain.com/)
	* [播客链电报群](https://t.me/Bobbyglobal) 
	* [播客链白皮书下载地址](http://yibokeclips.otvcloud.com/uploads/apks/blkV2.1.pdf) 
	* [播客链推特](https://twitter.com/BokerBobby)		
		
# 播客链技术负责人	
	* 微信公众号 [区什么块什么链啊](Blockchain_fxh7622) 
	* 推特 [区什么块什么链啊](https://twitter.com/chain_fxh7622) 
	
# 智能合约负责人
	* 推特 [后青春期的诗](https://twitter.com/chain_stayreal)
	

# 播客链配置文件(boker.json)示例配置

{
	"boker": {
		"dpos": {
			"validators": ["0x1aa228dde26f02e1cc5551cb4f1d74d0e998d24a"]
		},
		"contracts": [{
				"contractType": "1",
				"deployAddress": "0x1aa228dde26f02e1cc5551cb4f1d74d0e998d24a",
				"contractAddress": "0x1aa228dde26f02e1cc5551cb4f1d74d0e998d24a"
			},
			{
				"contractType": "2",
				"deployAddress": "0x1aa228dde26f02e1cc5551cb4f1d74d0e998d24a",
				"contractAddress": "0x1aa228dde26f02e1cc5551cb4f1d74d0e998d24a"
			}
		]
	}
}

# 在创世配置文件(genesis.json)中也可以进行设置,示例配置

{
	"config": {
		"chainId": 0,
		"homesteadBlock": 0,
		"eip155Block": 0,
		"eip158Block": 0,
		"boker": {
			"dpos": {
				"validators": ["0x1aa228dde26f02e1cc5551cb4f1d74d0e998d24a"]
			},
			"contracts": [{
					"contractType": "1",
					"deployAddress": "0x1aa228dde26f02e1cc5551cb4f1d74d0e998d24a",
					"contractAddress": "0x1aa228dde26f02e1cc5551cb4f1d74d0e998d24a"
				},
				{
					"contractType": "2",
					"deployAddress": "0x1aa228dde26f02e1cc5551cb4f1d74d0e998d24a",
					"contractAddress": "0x1aa228dde26f02e1cc5551cb4f1d74d0e998d24a"
				}
			]
		}
	},
	"alloc": {
		"0x1aa228dde26f02e1cc5551cb4f1d74d0e998d24a": {
			"balance": "10000"
		}
	},
	"difficulty": "0x000001",
	"extraData": "",
	"gasLimit": "0xffffffff"
}