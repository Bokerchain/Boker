## Boker Chain

Official golang implementation of Boker Chain.

## Contact us

* Whitepaper: [blkV2.1.pdf](http://yibokeclips.otvcloud.com/uploads/apks/blkV2.1.pdf)
* Website: [播客链](http://www.videoblkchain.com/)
* Telegram: [Bobbyglobal](https://t.me/Bobbyglobal)
* Twitter： [BokerChain](https://twitter.com/BokerBobby)
		
### Boker Chain CTO
* WeChat: [区什么块什么链啊](Blockchain_fxh7622) 
* Twitter: [区什么块什么链啊](https://twitter.com/chain_fxh7622) 

### Smart Contract
* Twitter: [后青春期的诗](https://twitter.com/chain_stayreal)

## Building the source

Building Boker Chain requires both a Go (version 1.7.3 or later) and a C compiler. You can install them using your favourite package manager. Once the dependencies are installed, run

    make geth

or, to build the full suite of utilities:

    make all

## Executables

The Boker Chain project comes with several wrappers/executables found in the `cmd` directory.

| Command    | Description |
|:----------:|-------------|
| **`geth`** | Our main Ethereum CLI client. It is the entry point into the Ethereum network (main-, test- or private net), capable of running as a full node (default) archive node (retaining all historical state) or a light node (retrieving data live). It can be used by other processes as a gateway into the Ethereum network via JSON RPC endpoints exposed on top of HTTP, WebSocket and/or IPC transports. `geth --help` and the [CLI Wiki page](https://github.com/ethereum/go-ethereum/wiki/Command-Line-Options) for command line options. |
| `abigen` | Source code generator to convert Ethereum contract definitions into easy to use, compile-time type-safe Go packages. It operates on plain [Ethereum contract ABIs](https://github.com/ethereum/wiki/wiki/Ethereum-Contract-ABI) with expanded functionality if the contract bytecode is also available. However it also accepts Solidity source files, making development much more streamlined. Please see our [Native DApps](https://github.com/ethereum/go-ethereum/wiki/Native-DApps:-Go-bindings-to-Ethereum-contracts) wiki page for details. |


## Starting up your Boker Chain nodes

### 1、Initialize with genesis.json
	geth --datadir "/projects/ethereum/geth/node" init genesis.json
* `--datadir` flag specify the data directory of your node

### 2、Run
	nohup geth --nodiscover --maxpeers 3 --identity "bokerchain" --rpc --rpcaddr 0.0.0.0 --rpccorsdomain "*" --rpcvhosts '*' --datadir "/projects/ethereum/geth/node" --port 30303 --rpcapi "db,eth,net,web3" --networkid 96579 &
	
* `--rpc` Enable the HTTP-RPC server
* `--rpcaddr` HTTP-RPC server listening interface (default: "localhost"). If you want to access RPC from other containers and/or hosts, it should be set to `--rpcaddr 0.0.0.0`.
* `--rpcapi` API's offered over the HTTP-RPC interface (default: "eth,net,web3")
* `--rpccorsdomain` Comma separated list of domains from which to accept cross origin requests (browser enforced)
* `--datadir`		flag should be idential to that in first step
* `--networkid`	flag specify your private net id

### 3、Enter console mode
	geth attach ipc:/projects/ethereum/geth/node/geth.ipc


### 4、Create new account
	personal.newAccount()

### 5、Unlock account
	personal.unlockAccount("0x1d4443a3eff8a5df88ecd9d91a58037585289be0", "123456", 0)

### 6、Set current account as first validator to process first transaction
	miner.setLocalValidator()

### 7、Set new validator
	eth.addValidator("0x1d4443a3eff8a5df88ecd9d91a58037585289be0", 10000)

### 8、Start mining
	miner.start()


## Configuration

### genesis.json
```json
{
	"config": {
		"chainId": 0,
		"homesteadBlock": 0,
		"eip155Block": 0,
		"eip158Block": 0
	},
	"alloc": {},
	"difficulty": "0x000001",
	"extraData": "",
	"gasLimit": "0xffffffff"
}

```

### boker.json
```json
{
    "dpos":{
        "validators":[
            "0xdd165ba267593d2acc837fc507c2e94e802817d9"
        ]
    },
    "contracts":{
        "bases":[
            {
                "contracttype":2,
                "deployaddress":"0xdd165ba267593d2acc837fc507c2e94e802817d9",
                "contractaddress":"0xd7fd311c8f97349670963d87f37a68794dfa80ff"
            }
        ]
    },
    "producer":{
        "coinbase":"0x1aa228dde26f02e1cc5551cb4f1d74d0e998d24a",
        "password":"123456"
    }
}
```

### Programatically interfacing Geth nodes

As a developer, sooner rather than later you'll want to start interacting with Boker Chain 
network via your own programs and not manually through the console. To aid this, Boker Chain has built in
support for a JSON-RPC based APIs. These can be
exposed via HTTP, WebSockets and IPC (unix sockets on unix based platforms, and named pipes on Windows).

The IPC interface is enabled by default and exposes all the APIs supported by Geth, whereas the HTTP
and WS interfaces need to manually be enabled and only expose a subset of APIs due to security reasons.
These can be turned on/off and configured as you'd expect.

HTTP based JSON-RPC API options:

  * `--rpc` Enable the HTTP-RPC server
  * `--rpcaddr` HTTP-RPC server listening interface (default: "localhost")
  * `--rpcport` HTTP-RPC server listening port (default: 8545)
  * `--rpcapi` API's offered over the HTTP-RPC interface (default: "eth,net,web3")
  * `--rpccorsdomain` Comma separated list of domains from which to accept cross origin requests (browser enforced)
  * `--ws` Enable the WS-RPC server
  * `--wsaddr` WS-RPC server listening interface (default: "localhost")
  * `--wsport` WS-RPC server listening port (default: 8546)
  * `--wsapi` API's offered over the WS-RPC interface (default: "eth,net,web3")
  * `--wsorigins` Origins from which to accept websockets requests
  * `--ipcdisable` Disable the IPC-RPC server
  * `--ipcapi` API's offered over the IPC-RPC interface (default: "admin,debug,eth,miner,net,personal,shh,txpool,web3")
  * `--ipcpath` Filename for IPC socket/pipe within the datadir (explicit paths escape it)

You'll need to use your own programming environments' capabilities (libraries, tools, etc) to connect
via HTTP, WS or IPC to a Boker Chain node configured with the above flags and you'll need to speak [JSON-RPC](http://www.jsonrpc.org/specification)
on all transports. You can reuse the same connection for multiple requests!

**Note: Please understand the security implications of opening up an HTTP/WS based transport before
doing so! Hackers on the internet are actively trying to subvert Boker Chain nodes with exposed APIs!
Further, all browser tabs can access locally running webservers, so malicious webpages could try to
subvert locally available APIs!**


## License

The Boker Chain binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included
in our repository in the `COPYING` file.
