
# 类说明(eth/api.go)

	* PublicEthereumAPI	提供了访问以太网完全节点相关的API信息
	* PublicMinerAPI		提供用来控制矿工的API，它仅提供对数据进行操作的方法，这些方法在可公开访问且不会带来安全风险。
	* PrivateMinerAPI	提供用来控制矿工的私有RPC方法，由于这些方法可能被外部用户所滥用，所以必须被认为是不安全的，供不信任的用户使用。
	* PrivateAdminAPI	以太坊全节点相关API的集合，通过私有管理端点公开。
	* PublicDebugAPI		公开的以太坊全节点API，通过公共调试端点
	* PrivateDebugAPI	公开的以太坊全节点API，私有调试端点
	
## 这些类中均有 *Ethereum 类