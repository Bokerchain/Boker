## 播客链超级节点运行流程

# 运行流程

# 第一步：初始化创世文件
	geth --datadir "/projects/ethereum/geth/node" init genesis.json

# 第二步：启动geth
	nohup geth --nodiscover  \
	--maxpeers 3 \
	--identity "bokerchain" \
	--rpc \
	--rpccorsdomain "*" \
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
	miner.selfValidator()

# 第七步：设置验证人（这里使用假定账号、票数）
	eth.addValidator("0x369bbedf102bd6f179e26ea0a7f434992ab9c0bf", 10000)

# 第八步：启动挖矿
	miner.start()

# 第九步：终止挖矿
	miner.stop()
	
	
## 查询信息

# 1：查询区块数量
	eth.blockNumber

# 2：查询区块内容
	eth.getBlock(1) 

# 3：查询交易内容（这里使用假定交易的Hash）
	eth.getTransaction("0xc51a2d7e722ef61bdb05be5a6725072f9181572073c9426cb9defe9aa91efdef")