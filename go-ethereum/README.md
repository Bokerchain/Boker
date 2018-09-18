## Boker Chain

Official golang implementation of Boker Chain. 

[Website](http://www.videoblkchain.com/)<br/>
[Telegram](https://t.me/Bobbyglobal)<br/>
[Whitepaper](http://yibokeclips.otvcloud.com/uploads/apks/blkV2.1.pdf)<br/>
[Twitter](https://twitter.com/BokerBobby)	<br/>

Automated builds are available for stable releases and the unstable master branch.
Binary archives are published at https://geth.ethereum.org/downloads/.

## Building the source

Building Boker Chain requires both a Go (version 1.7 or later) and a C compiler. You can install them using your favourite package manager. Once the dependencies are installed, run

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
--datadir flag specify the data directory of your node

### 2、Run
	nohup geth --nodiscover --maxpeers 3 --identity "bokerchain" --rpc --rpcaddr 0.0.0.0 --rpccorsdomain "*" --rpcvhosts '*' --datadir "/projects/ethereum/geth/node" --port 30303 --rpcapi "db,eth,net,web3" --networkid 96579 &
--datadir		flag should be idential to that in first step  <br/>
--networkid`	flag specify your private net id

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

As an alternative to passing the numerous flags to the `geth` binary, you can also pass a configuration file via:

```
$ geth --config /path/to/your_config.toml
```

To get an idea how the file should look like you can use the `dumpconfig` subcommand to export your existing configuration:

```
$ geth --your-favourite-flags dumpconfig
```

*Note: This works only with geth v1.6.0 and above.*

#### Docker quick start

One of the quickest ways to get Ethereum up and running on your machine is by using Docker:

```
docker run -d --name ethereum-node -v /Users/alice/ethereum:/root \
           -p 8545:8545 -p 30303:30303 \
           ethereum/client-go --fast --cache=512
```

This will start geth in fast sync mode with a DB memory allowance of 512MB just as the above command does.  It will also create a persistent volume in your home directory for saving your blockchain as well as map the default ports. There is also an `alpine` tag available for a slim version of the image.

Do not forget `--rpcaddr 0.0.0.0`, if you want to access RPC from other containers and/or hosts. By default, `geth` binds to the local interface and RPC endpoints is not accessible from the outside.

### Programatically interfacing Geth nodes

As a developer, sooner rather than later you'll want to start interacting with Geth and the Ethereum
network via your own programs and not manually through the console. To aid this, Geth has built in
support for a JSON-RPC based APIs ([standard APIs](https://github.com/ethereum/wiki/wiki/JSON-RPC) and
[Geth specific APIs](https://github.com/ethereum/go-ethereum/wiki/Management-APIs)). These can be
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
via HTTP, WS or IPC to a Geth node configured with the above flags and you'll need to speak [JSON-RPC](http://www.jsonrpc.org/specification)
on all transports. You can reuse the same connection for multiple requests!

**Note: Please understand the security implications of opening up an HTTP/WS based transport before
doing so! Hackers on the internet are actively trying to subvert Ethereum nodes with exposed APIs!
Further, all browser tabs can access locally running webservers, so malicious webpages could try to
subvert locally available APIs!**

### Operating a private network

Maintaining your own private network is more involved as a lot of configurations taken for granted in
the official networks need to be manually set up.

#### Defining the private genesis state

First, you'll need to create the genesis state of your networks, which all nodes need to be aware of
and agree upon. This consists of a small JSON file (e.g. call it `genesis.json`):

```json
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
	

# boker.json, for Example

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

# genesis.json, for Example

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
```

The above fields should be fine for most purposes, although we'd recommend changing the `nonce` to
some random value so you prevent unknown remote nodes from being able to connect to you. If you'd
like to pre-fund some accounts for easier testing, you can populate the `alloc` field with account
configs:

```json
"alloc": {
  "0x0000000000000000000000000000000000000001": {"balance": "111111111"},
  "0x0000000000000000000000000000000000000002": {"balance": "222222222"}
}
```

With the genesis state defined in the above JSON file, you'll need to initialize **every** Geth node
with it prior to starting it up to ensure all blockchain parameters are correctly set:

```
$ geth init path/to/genesis.json
```

#### Creating the rendezvous point

With all nodes that you want to run initialized to the desired genesis state, you'll need to start a
bootstrap node that others can use to find each other in your network and/or over the internet. The
clean way is to configure and run a dedicated bootnode:

```
$ bootnode --genkey=boot.key
$ bootnode --nodekey=boot.key
```

With the bootnode online, it will display an [`enode` URL](https://github.com/ethereum/wiki/wiki/enode-url-format)
that other nodes can use to connect to it and exchange peer information. Make sure to replace the
displayed IP address information (most probably `[::]`) with your externally accessible IP to get the
actual `enode` URL.

*Note: You could also use a full fledged Geth node as a bootnode, but it's the less recommended way.*

#### Starting up your member nodes

With the bootnode operational and externally reachable (you can try `telnet <ip> <port>` to ensure
it's indeed reachable), start every subsequent Geth node pointed to the bootnode for peer discovery
via the `--bootnodes` flag. It will probably also be desirable to keep the data directory of your
private network separated, so do also specify a custom `--datadir` flag.

```
$ geth --datadir=path/to/custom/data/folder --bootnodes=<bootnode-enode-url-from-above>
```

*Note: Since your network will be completely cut off from the main and test networks, you'll also
need to configure a miner to process transactions and create new blocks for you.*

#### Running a private miner

Mining on the public Ethereum network is a complex task as it's only feasible using GPUs, requiring
an OpenCL or CUDA enabled `ethminer` instance. For information on such a setup, please consult the
[EtherMining subreddit](https://www.reddit.com/r/EtherMining/) and the [Genoil miner](https://github.com/Genoil/cpp-ethereum)
repository.

In a private network setting however, a single CPU miner instance is more than enough for practical
purposes as it can produce a stable stream of blocks at the correct intervals without needing heavy
resources (consider running on a single thread, no need for multiple ones either). To start a Geth
instance for mining, run it with all your usual flags, extended by:

```
$ geth <usual-flags> --mine --minerthreads=1 --coinbase=0x0000000000000000000000000000000000000000
```

Which will start mining blocks and transactions on a single CPU thread, crediting all proceedings to
the account specified by `--coinbase`. You can further tune the mining by changing the default gas
limit blocks converge to (`--targetgaslimit`) and the price transactions are accepted at (`--gasprice`).

## Contribution

Thank you for considering to help out with the source code! We welcome contributions from
anyone on the internet, and are grateful for even the smallest of fixes!

If you'd like to contribute to go-ethereum, please fork, fix, commit and send a pull request
for the maintainers to review and merge into the main code base. If you wish to submit more
complex changes though, please check up with the core devs first on [our gitter channel](https://gitter.im/ethereum/go-ethereum)
to ensure those changes are in line with the general philosophy of the project and/or get some
early feedback which can make both your efforts much lighter as well as our review and merge
procedures quick and simple.

Please make sure your contributions adhere to our coding guidelines:

 * Code must adhere to the official Go [formatting](https://golang.org/doc/effective_go.html#formatting) guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt/)).
 * Code must be documented adhering to the official Go [commentary](https://golang.org/doc/effective_go.html#commentary) guidelines.
 * Pull requests need to be based on and opened against the `master` branch.
 * Commit messages should be prefixed with the package(s) they modify.
   * E.g. "eth, rpc: make trace configs optional"

Please see the [Developers' Guide](https://github.com/ethereum/go-ethereum/wiki/Developers'-Guide)
for more details on configuring your environment, managing project dependencies and testing procedures.

## License

The go-ethereum library (i.e. all code outside of the `cmd` directory) is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html), also
included in our repository in the `COPYING.LESSER` file.

The go-ethereum binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included
in our repository in the `COPYING` file.
