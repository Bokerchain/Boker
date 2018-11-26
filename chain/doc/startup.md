## 播客链超级节点运行流程

## 准备流程
# 同步时间：
	/usr/sbin/ntpdate te cn.pool.ntp.org > / > /dev/null 2>&1
	

## 运行流程

# 第一步：初始化创世文件
	geth --datadir "/projects/ethereum/geth/node" init genesis.json

# 第二步：启动geth
	nohup geth --nodiscover \
	--maxpeers 3 \
	--identity "bokerchain" \
	--rpc \
	--rpcaddr 0.0.0.0 \
	--rpccorsdomain "*" \
	--rpcvhosts '*' \
	--datadir "/projects/ethereum/geth/node" \
	--port 30303 \
	--rpcapi "db,eth,net,web3" \
	--networkid 96579 &

# 第三步：进入geth控制台
	geth attach ipc:/projects/ethereum/geth/node/geth.ipc

# 第四步：创建账号
	personal.newAccount()

# 第五步：设置帐号解锁（这里使用假定账号、密码）
	personal.unlockAccount("0x369bbedf102bd6f179e26ea0a7f434992ab9c0bf", "123456", 0)

# 第六步：设置自己为验证人
	miner.setLocalValidator()

# 第七步：设置验证人（这里使用假定账号、票数）
	eth.addValidator("0x369bbedf102bd6f179e26ea0a7f434992ab9c0bf", 10000)

# 第八步：启动挖矿
	miner.start()

# 第九步：终止挖矿
	miner.stop()
	
	
##	发布合约

# 第一步：通过solc或者vscode等工具编译得到abi文件和bin文件
	solc test.sol
	
# 第二步：新建一个js文件，比如Deploy.js，将一下内容拷贝至文件，其中abi替换成abi文件中的内容，bin替换成bin文件中的内容
	var testContract = web3.eth.contract(abi);
	var test = testContract.new(
	   {
	     from: web3.eth.accounts[0],
	     data: 'bin',
	     gas: '4700000'
	   }, function (e, contract){
	    console.log(e, contract);
	    if (typeof contract.address !== 'undefined') {
	         console.log('Contract mined! address: ' + contract.address + ' transactionHash: ' + contract.transactionHash);
	    }
	 })
	
# 第三步：进入geth控制台，输入命令，会返回交易hash，交易被打包后，可以看到输出信息Contract mined!
	loadScript("Deploy.js")
	
	
##	设置基础合约

	eth.setBaseContracts("0xff2e5867f89e7be22e8c4a3cd9fb59bfd31ce681", 2, "[{\"constant\":true,\"inputs\":[],\"name\":\"cfoAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ceoAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getLogSize\",\"outputs\":[{\"name\":\"size\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"logKeyDefault\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"assgineTokenPeriodDefault\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newCEO\",\"type\":\"address\"}],\"name\":\"setCEO\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"index\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newCOO\",\"type\":\"address\"}],\"name\":\"setCOO\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"enable\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getLog\",\"outputs\":[{\"name\":\"level\",\"type\":\"uint8\"},{\"name\":\"time\",\"type\":\"uint256\"},{\"name\":\"key\",\"type\":\"string\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"v1\",\"type\":\"uint256\"},{\"name\":\"v2\",\"type\":\"uint256\"},{\"name\":\"v3\",\"type\":\"uint256\"},{\"name\":\"remarks\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"eventType\",\"type\":\"uint256\"},{\"name\":\"addrFrom\",\"type\":\"address\"},{\"name\":\"addrTo\",\"type\":\"address\"},{\"name\":\"eventValue1\",\"type\":\"uint256\"},{\"name\":\"eventValue2\",\"type\":\"uint256\"}],\"name\":\"fireUserEvent\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"checkAssignToken\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"assginedTokensPerPeriodDefault\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newCFO\",\"type\":\"address\"}],\"name\":\"setCFO\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"clearLog\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"implAddress\",\"type\":\"address\"}],\"name\":\"setImpl\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"logLevel\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"impl\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"assignToken\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"cooAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"level\",\"type\":\"uint8\"}],\"name\":\"setLevel\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"logEnabled\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"}]")
	
	
## 查询信息

# 1：查询区块数量
	eth.blockNumber

# 2：查询区块内容
	eth.getBlock(1) 

# 3：查询交易内容（这里使用假定交易的Hash）
	eth.getTransaction("0xc51a2d7e722ef61bdb05be5a6725072f9181572073c9426cb9defe9aa91efdef")