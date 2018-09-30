![Image text](https://github.com/Bokerchain/Boker/blob/master/Boker.png)

## Bokerchain

Bokerchain is a public blockchain platform that serves the vertical area of audio & video. It is convenient for different intelligent terminal devices to access to Bokerchain. We can form a Video Application Union by providing SDK for various video APP, meeting the need for copyright protection, data sharing and benefit protection among all apps in Video Application Union. In this situation, we may make data among podcasts, advertisers and users more transparent, also maximizing the benefits.<br/>Our goal is to acheive video sharing, benefits sharing and user resources sharing, benefiting podcasts, advertisers and our users while providing entertainment.

## Contact us

* Whitepaper: [blkV2.1.pdf](http://yibokeclips.otvcloud.com/uploads/apks/blkV2.1.pdf)
* Website: [播客链](http://www.videoblkchain.com/)
* Telegram: [Bobbyglobal](https://t.me/Bobbyglobal)
* Twitter： [BokerChain](https://twitter.com/BokerBobby)
		
### Bokerchain Co-Founder
* WeChat: [区什么块什么链啊](Blockchain_fxh7622) 
* Twitter: [区什么块什么链啊](https://twitter.com/chain_fxh7622) 

### Smart Contract
* Twitter: [后青春期的诗](https://twitter.com/chain_stayreal)

## Building the source

Building Bokerchain requires both a Go (version 1.7.3 or later) and a C compiler. You can install them using your favourite package manager. Once the dependencies are installed, run

    make geth

or, to build the full suite of utilities:

    make all

## Executables

The Bokerchain project comes with several wrappers/executables found in the `cmd` directory.

| Command    | Description |
|:----------:|-------------|
| **`geth`** | Our main Ethereum CLI client. It is the entry point into the Ethereum network (main-, test- or private net), capable of running as a full node (default) archive node (retaining all historical state) or a light node (retrieving data live). It can be used by other processes as a gateway into the Ethereum network via JSON RPC endpoints exposed on top of HTTP, WebSocket and/or IPC transports. `geth --help` and the [CLI Wiki page](https://github.com/ethereum/chain/wiki/Command-Line-Options) for command line options. |
| `abigen` | Source code generator to convert Ethereum contract definitions into easy to use, compile-time type-safe Go packages. It operates on plain [Ethereum contract ABIs](https://github.com/ethereum/wiki/wiki/Ethereum-Contract-ABI) with expanded functionality if the contract bytecode is also available. However it also accepts Solidity source files, making development much more streamlined. Please see our [Native DApps](https://github.com/ethereum/chain/wiki/Native-DApps:-Go-bindings-to-Ethereum-contracts) wiki page for details. |

## Configuration

### genesis.json
```json
{
	"config": {
		"chainId": 0,
		"byzantiumBlock": 0,
		"eip155Block": 0,
		"eip158Block": 0
	},
	"alloc": {},
	"difficulty": "0x000001",
	"extraData": "",
	"gasLimit": "0xffffffff"
}

```
**Note: `"byzantiumBlock": 0` should be the config value, otherwise contract call may malfunction!**

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

## Starting up your Bokerchain nodes
### Prerequisit
Before start up your Bokerchain node, make sure the time is up to date. You can acheive this by adding crontab task as follows:
```
crontab -e
```
add commond:
```
*/10 * * * * /usr/sbin/ntpdate 1.cn.pool.ntp.org
```
It means synchronize time with 1.cn.pool.ntp.org every 10 minutes. You can adjust the commond according to your demand.

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

## Programatically interfacing Geth nodes

As a developer, sooner rather than later you'll want to start interacting with Bokerchain 
network via your own programs and not manually through the console. To aid this, Bokerchain has built in
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
via HTTP, WS or IPC to a Bokerchain node configured with the above flags and you'll need to speak [JSON-RPC](http://www.jsonrpc.org/specification)
on all transports. You can reuse the same connection for multiple requests!

**Note: Please understand the security implications of opening up an HTTP/WS based transport before
doing so! Hackers on the internet are actively trying to subvert Bokerchain nodes with exposed APIs!
Further, all browser tabs can access locally running webservers, so malicious webpages could try to
subvert locally available APIs!**

## Deploy Contracts

for example:
```solidity
//Test.sol
pragma solidity ^0.4.8;

contract Test {
    string public message;

    function Set(string m) public  {
        message = m;
    }

    function Show() public view returns (string){
        return message;
    }
}
```

### 1、Compile with solc
	solc Test.sol
you will get two files : `Test.sol:Test.abi` and `Test.sol:Test.bin`

### 2、Create and edit js file
create a js file like: depoly_test.js. Copy the following text:
```javascript
var myContract = web3.eth.contract([{"constant":false,"inputs":[{"name":"m","type":"string"}],"name":"Set","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"Show","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"message","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"}]);
var contract = implContract.new({
    from: web3.eth.accounts[0],
    data: '608060405234801561001057600080fd5b50610410806100206000396000f300608060405260043610610057576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806334541fcd1461005c578063db5e85a5146100c5578063e21f37ce14610155575b600080fd5b34801561006857600080fd5b506100c3600480360381019080803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506101e5565b005b3480156100d157600080fd5b506100da6101ff565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561011a5780820151818401526020810190506100ff565b50505050905090810190601f1680156101475780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561016157600080fd5b5061016a6102a1565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156101aa57808201518184015260208101905061018f565b50505050905090810190601f1680156101d75780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b80600090805190602001906101fb92919061033f565b5050565b606060008054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156102975780601f1061026c57610100808354040283529160200191610297565b820191906000526020600020905b81548152906001019060200180831161027a57829003601f168201915b5050505050905090565b60008054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156103375780601f1061030c57610100808354040283529160200191610337565b820191906000526020600020905b81548152906001019060200180831161031a57829003601f168201915b505050505081565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061038057805160ff19168380011785556103ae565b828001600101855582156103ae579182015b828111156103ad578251825591602001919060010190610392565b5b5090506103bb91906103bf565b5090565b6103e191905b808211156103dd5760008160009055506001016103c5565b5090565b905600a165627a7a723058207ba0ff699809e34d9ba764a2c026eaa14a08f74c77e650af328da93bdd81ea310029',
    gas: '4700000'
},
function(e, contract) {
    console.log(e, contract);
    if (typeof contract.address !== 'undefined') {
        console.log('Contract mined! address: ' + contract.address + ' transactionHash: ' + contract.transactionHash);
    }
})
```
  * the content between `web3.eth.contract(...)` in first line is same as that in `Test.sol:Test.abi`。You can replace it with your contract abi file.
  * the value for `data` in 4th line is same as that in `Test.sol:Test.bin`。You can replace it with your contract bin file.

### 3、Attach to geth.ipc
	geth attach ipc:geth.ipc
	
### 4、Load javascrip file
	loadScript("depoly_test.js")
wait for the contract to be mined. When it is done, you will get the message `"Contract mined!address:... transactionHash:..."`
You can get detail info by `eth.getTransaction` command

## Call contract function
```
contract.Set.sendTransaction("Hello Solidity",{from: eth.accounts[0]});
contract.Show()
```

## License

The Bokerchain binaries (i.e. all code inside of the `cmd` directory) is licensed under the
[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html), also included
in our repository in the `COPYING` file.
