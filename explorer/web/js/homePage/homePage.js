/*
 * 首页js
 */
var preventDef = function(e) {
//	e.preventDefault();	
	e.stopPropagation();
};
function getContract(){
	var xhr=new XMLHttpRequest();
	xhr.open("POST",$base+"/rest/contractAddress", false);
	xhr.onreadystatechange=function(){
		if (xhr.readyState == 4 && xhr.status == 200){
			//请求成功
			window.$contract = JSON.parse(xhr.responseText).interfaceAddress;
		}
	};
	xhr.send();
};
getContract();

var web3;
// 创建web3对象并连接到以太坊节点
if(typeof web3 !== 'undefined') {
	web3 = new Web3(web3.currentProvider);
} else {
	web3 = new Web3(new Web3.providers.HttpProvider($serverIp));
}
var infoContract = web3.eth.contract(ABI);
//参数配 合约地址
var info = infoContract.at($contract);

var vm = new Vue({
	el: '#homePage',
	mixins: [mixin],
	data: {
		isLogin: false, //是否登录
//		isFirstLogin: 0,//是否是首次登陆  0第一次登录
//		noticeList: [{
//			'content': 'test123'
//		}], //公告列表
		noticeList:[],
		activeIndex: '-1',//
		tabIndex: 0, //0表示‘私钥’tab   1表示‘二维码’tab
		userTotal:0,//居民人数
		tokenAssigned: 0, //昨日分出Bobby
		tokenToAssign: 0, //待分配Bobby
		tel: '',//手机号
		code: '',//验证码
		userId: 0,//用户ID
		password: '',//用户设置的 交易密码
		password1: '',
		password2:'',//用户修改 交易密码
		password3:'',
		inviteCode: '', //别人的邀请码
		nickname: '',//昵称
		address: '', //钱包地址
		myInviteCode: '', //自己的邀请码
		keystore: '', //私钥+交易密码 生成key-store
		balance: 0, //余额
		intNum: 0,
		floatNum: 0,
		longtermPower: 0,
		shorttermPower: 0, //临时算力
		privateKey: '', //私钥
		inputPassword: '', //用户导入私钥的时候 输入的交易密码
		bobbyDetailList: [], //bobby明细
		scoreDetailList: [], //算力明细
		inviteList:[],//当前用户的邀请人列表
		rewardRecordList:[],//打赏记录
		countDown:'-1',//倒计时用
		isClicked:false,
		isFromLast:-1,//如果点击提升算力的“邀请好友”跳到邀请奖励模态框 1
		transferToAddress:'',//打赏转账的对方的地址
		transferMoney:'',//打赏金额
		password4:''//用户打赏时输入的交易密码
	},
	created: function() {
		//把所有按钮置于非活性
		var _this = this;
		console.log(_this.tel.length)
		for(var i = 0; i < this.btnList.length - 1; i++) {
			this.btnList[i].active = false
		}
		_this.getInitData();
		_this.getGeneralInfo();
		var params = {
			keyBytes: 32,
			ivBytes: 16
		};
		console.log(_this.activeIndex);
		// synchronous
		//		var dk = keythereum.create(params);
		//		var aaa='123456';
		//		var bbb={"address":"3d77953e6a79710f062697b40c31501ff762bce7","crypto":{"cipher":"aes-128-ctr","ciphertext":"27d70880ed759ee28cd212ef0b0d099dc55226774535cbc537ab5f6a084f5a18","cipherparams":{"iv":"dd29a6894d566077a1915b06be76e3db"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"3b9f814aa699523cb35b77dd10aa8d7cb1a2c0dfa88dc8968c6a909204feeea2"},"mac":"07a9cb4669b347e1a75b880e5299f6c916fead1e113d37d851836b880d55365a"},"id":"54ab58c2-c853-48f9-97aa-1fed676e2c56","version":3}
		//		var privateKey = keythereum.recover(aaa, bbb);
		//		console.log(Bytes2Str(privateKey));
	},
	watch: {
		activeIndex: function(val) {
			if(val != '-1' && val != 2 && val != 3) {
				this.forbidScroll(this.activeIndex);
			} else {
				this.unforbidScroll(this.activeIndex);
				//重置thirdPatyList状态
				var list = this.thirdPatyList;
				if(list) {
					for(var thirdParty in list) {
						list[thirdParty].showDetail = false;
					}
				}
			}
		},
		thirdPatyList: {
			handler: function() {
				console.log(1)
				var _this = this;
				var list = this.thirdPatyList;
				var i=0;
				if(list) {
					for(var thirdParty in list) {
						i++;
						if(list[thirdParty].showDetail) {
							console.log(i);
							_this.isShowDetail = true;
							return;
						}
					}
					_this.isShowDetail = false;
				}
			},
			deep: true
		},
		tel: function(val){
			var ele = document.getElementById('getCodeBtn');
			if(val.length!=0){
				//可点击状态
				ele.style.backgroundImage="url(img/register/btn1@3x.png)";
			}else{
				//不可点击状态
				ele.style.backgroundImage="url(img/register/btn3@3x.png)";
			}
		},
		countDown:function(val){
			var ele = document.getElementById('getCodeBtn');
			if(val==-1){
				//可点击状态
				ele.style.backgroundImage="url(img/register/btn1@3x.png)";
			}else{
				//不可点击状态
				ele.style.backgroundImage="url(img/register/btn3@3x.png)";
			}
		}
	},
	methods: {
		getInitData: function() {
			var _this = this;
			//判断是否有缓存 
			if(localStorage.getItem('address')) {
				_this.address = localStorage.getItem('address');
				_this.getUser(_this.address);
			}
			if(localStorage.getItem('keystore')) {
				_this.nickname = localStorage.getItem('keystore');
			}
			if(localStorage.getItem('nickname')) {
				_this.nickname = localStorage.getItem('nickname');
			}
			if(localStorage.getItem('password')) {
				_this.password = localStorage.getItem('password');
			}
			if(localStorage.getItem('privateKey')) {
				_this.privateKey = localStorage.getItem('privateKey');
			}
			if(localStorage.getItem('tel')) {
				_this.tel = localStorage.getItem('tel');
			}
			if(localStorage.getItem('myInviteCode')) {
				_this.myInviteCode = localStorage.getItem('myInviteCode');
			}
			if(localStorage.getItem('userId')) {
				_this.userId = localStorage.getItem('userId');
			}
			//判断是否是第一次登录 判断cookie是否过期 过期上报登录信息
			console.log(getCookie('firstLoginTime'));
			if(getCookie('firstLoginTime')!=null){
				//未失效
				console.log("xxxxxxxxxxxxx");
			}else{
				console.log("yyyyyyyyyyyy");
				//通过缓存登录
				if(localStorage.getItem('address')){
					//失效
					var current_date = new Date().getTime();
					setCookie("firstLoginTime",current_date);
					console.log("首次登陆时间"+new Date());
					//行为数据上报
					//上报登录信息 eventType:1
					console.log("上报登录信息");
					this.$http.post($base + '/rest/fireUserEvent', {
						'addrFrom': _this.address,
						'eventType': 1,
					}).then(function(res) {
						console.log(res.body);
					}).catch(function(){
						_this.$toast({
						  message: '上报登录信息失败！',
						  position: 'middle',
						  duration: 2000
						});
					});
					_this.$toast({
						message: '今日登录，算力+1',
						position: 'middle',
						duration: 1000
					});
				}
			}
		},
		getGeneralInfo: function() {
			var _this = this;
			//获取账户余额等信息
			info.generalInfoGet(function(err, result) {
				if(result) {
					console.log(result);
					_this.userTotal  = result[0].toNumber();
					_this.tokenAssigned = result[1].toNumber()*Math.pow(10, -18); 
					_this.tokenAssigned=_this.tokenAssigned.toFixed(4);
					_this.tokenToAssign = result[2].toNumber()*Math.pow(10, -18);
					_this.tokenToAssign = _this.tokenToAssign.toFixed(4);
				} else {
					_this.$toast({
						message: '获取居民人数失败！',
						position: 'middle',
						duration: 1000
					});
				}
			});
		},
		//禁止滑动
		forbidScroll: function(ele) {
			document.body.style.overflow = 'hidden';
			var docEle='';
			switch (ele){
//				case '5':document.getElementById("modals").addEventListener("touchmove", preventDef, false);
//						document.getElementById("parentBox").removeEventListener("touchmove", preventDef, false);
				default:break;
			}
//			document.addEventListener("touchmove", preventDef, false); //禁止页面滑动
		},
		//解除禁止滑动
		unforbidScroll: function(ele) {
			document.body.style.overflow = ''; //出现滚动条
			switch (ele){
//				case '5':document.getElementById("modals").removeEventListener("touchmove", preventDef, false);console.log(1);
				default:break;
			}
//			document.removeEventListener("touchmove", preventDef, false);
		},
		clickBtn: function(index) {
			var _this = this;
			if(_this.isLogin == false) {
				_this.$toast({
					message: '请先登录！',
					position: 'middle',
					duration: 1000
				});
				return;
			}
			if(index==2){
				_this.getUserDappBindInfo();
				_this.btnList[2].active = true;
				for(var i = 0; i < _this.btnList.length - 1; i++) {
					if(i != index) {
						_this.btnList[i].active = false ;
					} 
				}
				return;
			}
			if(index==3){
				_this.getInviteInfo();
				_this.btnList[3].active = true;
				for(var i = 0; i < _this.btnList.length - 1; i++) {
					if(i != index) {
						_this.btnList[i].active = false ;
					} 
				}
				return;
			}
			if(index==4){
				document.getElementById('sideBar').style.transform='translate3d(0, 0, 0)';
				document.getElementById('sideBar').style.visibility='visible';
				document.getElementById('modal-overlay').addEventListener('touchstart',function(){
					_this.closeSideBar();
				},false);
			}
			if(index==1){
				document.getElementById('rightSideBar').style.transform='translate3d(0, 0, 0)';
				document.getElementById('rightSideBar').style.visibility='visible';
				document.getElementById('modal-overlay').addEventListener('touchstart',function(){
					_this.closeRightSideBar();
				},false);
			}
			for(var i = 0; i < _this.btnList.length ; i++) {
				if(i == index) {
					_this.activeIndex = i;
					_this.btnList[i].active = true ;
				} else {
					_this.btnList[i].active = false ;
				}
			}
		},
		clickLock: function() {
			this.activeIndex = '6-0';
		},
		getCode: function() {
			//获取短信码
			if(this.tel == '' || this.tel == null || this.tel == 'undefined') {
				this.$toast({
					message: '请输入手机号',
					position: 'middle',
					duration: 1000
				});
				return false;
			}
			if(this.tel.charAt(0) != "1" || this.tel.length != 11) {
				this.$toast({
					message: '手机号码格式错误',
					position: 'middle',
					duration: 1000
				});
				return false;
			}
			var _this = this;
			_this.isClicked=true;
			this.$http.post($base + '/rest/sendVerificationCode', {
				'nationCode': '86',
				'phone': _this.tel
			}).then(function(res) {
				if(res.body.result == 1) {
					this.$toast({
						message: '验证码已发送',
						position: 'middle',
						duration: 1000
					});
					//定时器
					_this.countDown=60;
					var t = setInterval(function(){
						_this.countDown--;
						if(_this.countDown==0){
							_this.countDown = -1;
							_this.isClicked = false;
							clearInterval(t);
						}
					},1e3);
				} else {
					this.$toast({
						message: res.body.msg,
						position: 'middle',
						duration: 1000
					});
				}
			}).catch(function(){
				_this.$indicator.close();
				_this.isClicked=false;
				_this.$toast({
				  message: '获取验证码失败，请刷新重试！',
				  position: 'middle',
				  duration: 2000
				});
			});
		},
		register: function() {
			if(this.tel == '' || this.tel == null || this.tel == 'undefined') {
				this.$toast({
					message: '请输入手机号',
					position: 'middle',
					duration: 1000
				});
				return false;
			}
			if(this.tel.charAt(0) != "1" || this.tel.length != 11) {
				this.$toast({
					message: '手机号码格式错误',
					position: 'middle',
					duration: 1000
				});
				return false;
			}
			if(this.code == '' || this.code == null || this.code == 'undefined') {
				this.$toast({
					message: '请输入验证码！',
					position: 'middle',
					duration: 1000
				});
				return false;
			}
			///^1[34578]\d{9}$/.test(phone)
			if(!(/^1[34578]\d{9}$/.test(this.tel))) {
				console.log('error');
				this.$toast({
					message: '手机号码格式错误',
					position: 'middle',
					duration: 1000
				});
				return false;
			}
			var _this = this;
			this.$http.post($base + '/rest/checkVerification', {
				'phone': _this.tel,
				'code': _this.code
			}).then(function(res) {
				if(res.body.result == 1) {
					if(res.body.isRegister == "0") {
						//未注册 去设置昵称
						_this.activeIndex = '6-1';
					} else {
						//已注册 去导入私钥
						_this.myInviteCode = res.body.inviteCode;
						_this.nickname = res.body.nickname;
						_this.activeIndex = '6-2';
					}
				} else {
					this.$toast({
						message: res.body.msg,
						position: 'middle',
						duration: 1000
					});
				}
			}).catch(function(){
				_this.$indicator.close();
				_this.$toast({
				  message: '验证失败，请重试！',
				  position: 'middle',
				  duration: 2000
				});
			});
		},
		setPassword: function() {
			var _this = this;
			if(_this.nickname == '' || _this.nickname == null || _this.nickname == 'undefined') {
				_this.$toast({
					message: '请填写昵称',
					position: 'middle',
					duration: 1000
				});
				return false;
			}
			console.log(_this.nickname);
			if(_this.nickname.indexOf(" ")!=-1){
				_this.$toast({
					message: '昵称中不能包含空格',
					position: 'middle',
					duration: 1000
				});
				return false;
			}
			if(_this.nickname.length > 10) {
				_this.$toast({
					message: '昵称长度应小于10个字符',
					position: 'middle',
					duration: 1000
				});
				return false;
			}
			if(_this.password.indexOf(" ")!=-1|| _this.password1.indexOf(" ")!=-1) {
				_this.$toast({
					message: '交易密码中不能包含空格',
					position: 'middle',
					duration: 1000
				});
				return false;
			}
			if(_this.password.length < 6 || _this.password1.length < 6) {
				_this.$toast({
					message: '请输入6位以上字符！',
					position: 'middle',
					duration: 1000
				});
				return false;
			}
			var reg = new RegExp("[\\u4E00-\\u9FFF]+","g");
		    if(reg.test(_this.password)||reg.test(_this.password1)){
		        _this.$toast({
					message: '交易密码不能包含中文！',
					position: 'middle',
					duration: 1000
				});
				return false;
		    }
			if(_this.password != _this.password1) {
				_this.$toast({
					message: '确认交易密码需与交易密码一致',
					position: 'middle',
					duration: 1000
				});
				return false;
			}
			
			_this.$indicator.open({
				text: '创建私钥中...'
			});
			this.$http.post($base + '/rest/register', {
				'phone': _this.tel,
				'nickname': _this.nickname,
				'password': _this.password,
				'inviteCode': _this.inviteCode
			}).then(function(res) {
				if(res.body.result == 1) {
					_this.address = '0x' + res.body.address;
					_this.myInviteCode = res.body.inviteCode;
					_this.keystore = res.body.keystore;
					_this.userId = res.body.userId;
					localStorage.setItem('tel', _this.tel);
					localStorage.setItem('keystore', _this.keystore); //手机号
					localStorage.setItem('password', _this.password); //交易密码
					localStorage.setItem('address', _this.address); //钱包地址  16进制
					localStorage.setItem('keystore', _this.keystore);
					localStorage.setItem('myInviteCode', _this.myInviteCode);
					localStorage.setItem('nickname', _this.nickname);
					localStorage.setItem('userId', _this.userId);
					var current_date = new Date().getTime();
					setCookie("firstLoginTime",current_date);
					console.log("首次登陆时间"+new Date());
					//行为数据上报
					if(_this.inviteCode!=""){
						//填写了邀请码
						_this.$http.post($base + '/rest/fireUserEvent', {
							'addrFrom': _this.address,
							'eventType': 0,
							'addrTo': _this.inviteCode
						}).then(function(res) {
							
						}).catch(function(){
							_this.$toast({
							  message: '上报注册信息失败！',
							  position: 'middle',
							  duration: 2000
							});
						});
					}else{
						//上报注册信息 eventType:0
						_this.$http.post($base + '/rest/fireUserEvent', {
							'addrFrom': _this.address,
							'eventType': 0,
							'addrTo': ''
						}).then(function(res) {
							
						}).catch(function(){
							_this.$toast({
							  message: '上报注册信息失败！',
							  position: 'middle',
							  duration: 2000
							});
						});
						console.log("上报登录信息");
					}
					//上报登录信息
					_this.$http.post($base + '/rest/fireUserEvent', {
						'addrFrom': _this.address,
						'eventType': 1
					}).then(function(res) {
						
					}).catch(function(){
						_this.$toast({
						  message: '上报登录信息失败！',
						  position: 'middle',
						  duration: 2000
						});
					});
					//根据keystore获取key
					setTimeout(function() {
						var params = {
							keyBytes: 32,
							ivBytes: 16
						};
						// synchronous
						var dk = keythereum.create(params);
						var privateKey = keythereum.recover(_this.password, eval('(' + _this.keystore + ')'));
						_this.privateKey = Bytes2Str(privateKey);
						localStorage.setItem('privateKey', Bytes2Str(privateKey));
						//获取账户余额等信息
						_this.getUser(_this.address);
						_this.activeIndex=-1;
						_this.isLogin = true;
						_this.$indicator.close();
						_this.$toast({
							message: "今日登录，算力+1",
							position: 'middle',
							duration: 1000
						});
					}, 500)
				} else {
					_this.$indicator.close();
					this.$toast({
						message: res.body.msg,
						position: 'middle',
						duration: 1000
					});
				}
			}).catch(function(){
				_this.$indicator.close();
				_this.$toast({
				  message: '注册失败，请重试！',
				  position: 'middle',
				  duration: 2000
				});
			});
		},
		importKey: function() {
			var _this = this;
			if(_this.privateKey == '' || _this.privateKey == null || _this.privateKey == 'undefined') {
				_this.$toast({
					message: '请输入私钥！',
					position: 'middle',
					duration: 1000
				});
				return false;
			}
			if(_this.privateKey.length!=64) {
				_this.$toast({
					message: '请输入正确的私钥！',
					position: 'middle',
					duration: 1000
				});
				return false;
			}
			if(_this.inputPassword == '' || _this.inputPassword == null || _this.inputPassword == 'undefined') {
				_this.$toast({
					message: '请输入交易密码！',
					position: 'middle',
					duration: 1000
				});
				return false;
			}
			if(_this.inputPassword.indexOf(" ")!=-1) {
				_this.$toast({
					message: '交易密码中不能包含空格',
					position: 'middle',
					duration: 1000
				});
				return false;
			}
			if(_this.inputPassword.length < 6) {
				_this.$toast({
					message: '交易密码请输入6位以上字符！',
					position: 'middle',
					duration: 1000
				});
				return false;
			}
			var reg = new RegExp("[\\u4E00-\\u9FFF]+","g");
		    if(reg.test(_this.inputPassword)){
		        _this.$toast({
					message: '交易密码不能包含中文！',
					position: 'middle',
					duration: 1000
				});
				return false;
		    }
			_this.$indicator.open({
				text: '导入私钥中...'
				//spinnerType: 'fading-circle'
			});
			var timeMark=false;
			setTimeout(function() {
				if(timeMark==false){
					_this.$toast({
						message: '私钥错误',
						position: 'middle',
						duration: 1000
					});
					_this.$indicator.close();
					return;
				}
			},10000)
			setTimeout(function() {
				var params = {
					keyBytes: 32,
					ivBytes: 16
				};
				var dk = keythereum.create(params);
				var options = {
					kdf: "scrypt",
					cipher: "aes-128-ctr",
					kdfparams: {
						c: 262144,
						dklen: 32,
						prf: "hmac-sha256"
					}
				};
				var keyObject = keythereum.dump(_this.inputPassword, _this.privateKey, dk.salt, dk.iv, options);
				timeMark=true;
				if(keyObject) {
					localStorage.setItem('keystore', JSON.stringify(keyObject));
					//var privateKey= keythereum.recover(_this.inputPassword, keyObject);
					//localStorage.setItem('privateKey',Bytes2Str(privateKey))
					localStorage.setItem('address', keyObject.address.startsWith("0x") ? keyObject.address : '0x' + keyObject.address);
					_this.address = localStorage.getItem('address');
					localStorage.setItem('privateKey', _this.privateKey);
					localStorage.setItem('password', _this.inputPassword);
					localStorage.setItem('tel', _this.tel);
					localStorage.setItem('nickname', _this.nickname);
					localStorage.setItem('myInviteCode', _this.myInviteCode);
					//导入		
				} else {
					_this.$toast({
						message: '导入失败',
						position: 'middle',
						duration: 1000
					});
					_this.$indicator.close();
					return;
				}
				//校验手机号 和 钱包地址的绑定关系
				_this.$http.post($base + '/rest/checkAddress', {
					'phone': _this.tel,
					'address': _this.address
				}).then(function(res) {
					_this.$indicator.close();
					if(res.body.result == 1) {
						if(res.body.isSelf == 1) {
							//存储userId
							_this.userId = res.body.userId;
							localStorage.setItem('userId', res.body.userId);
							//获取账户余额等信息
							info.getUser(localStorage.getItem('address'), function(err, result) {
								if(result) {
									isImportSuccess=true;
									_this.$toast({
										message: '私钥导入成功',
										position: 'middle',
										duration: 1000
									});
									var balance = result[0].toNumber()
									var longtermPower = result[1].toNumber()
									var shorttermPower = result[2].toNumber()
									if(balance != 0) {
										_this.balance = balance * Math.pow(10, -18);
										_this.intNum = Math.floor(_this.balance);
										_this.floatNum = _this.balance.toString().substring(_this.balance.toString().indexOf(".")+1,_this.balance.toString().indexOf(".")+4);
									}
									if(longtermPower != 0) {
										_this.longtermPower = longtermPower * Math.pow(10, -3);
									}
									if(shorttermPower != 0) {
										_this.shorttermPower = shorttermPower * Math.pow(10, -3);
									}
									var current_date = new Date().getTime();
									setCookie("firstLoginTime",current_date);
									console.log("用户首次登陆时间"+new Date());
									//上报登录信息	 eventType:1
									console.log("上报用户行为信息");
									//上报登录信息
									_this.$http.post($base + '/rest/fireUserEvent', {
										'addrFrom': _this.address,
										'eventType': 1,
									}).then(function(res) {
										
									}).catch(function(){
										_this.$toast({
										  message: '上报登录信息失败！',
										  position: 'middle',
										  duration: 2000
										});
									});
//									info.fireUserEvents(_this.address,1,0x00,parseInt(current_date/1000),0x00,0x00);
									_this.activeIndex = "-1";
									_this.isLogin = true;
									_this.password = _this.inputPassword;
								} else {
									_this.$toast({
										message: '导入失败',
										position: 'middle',
										duration: 1000
									});
									_this.removeLocalStorage();
									//_this.activeIndex = "-1";
								}
							});
						} else {
							_this.$toast({
								message: '私钥错误',
								position: 'middle',
								duration: 1000
							});
							_this.removeLocalStorage();
						}
					} else {
						_this.$toast({
							message: res.body.msg,
							position: 'middle',
							duration: 1000
						});
						_this.removeLocalStorage();
					}
				}).catch(function(){
					_this.$indicator.close();
					_this.$toast({
					  message: '获取数据失败，请刷新重试',
					  position: 'middle',
					  duration: 2000
					});
				});
			}, 500);

		},
		getUser: function(address) {
			//获取账户余额等信息
			var _this = this;
			info.getUser(address, function(err, result) {
				if(result) {
					console.log(result);
					_this.isLogin = true;
					var balance = result[0].toNumber();
					var longtermPower = result[1].toNumber();
					var shorttermPower = result[2].toNumber();
					console.log(balance+'-'+longtermPower+'-'+shorttermPower);
					if(balance != 0) {
						_this.balance = balance * Math.pow(10, -18);
						_this.intNum = Math.floor(_this.balance);
						_this.floatNum = _this.balance.toString().substring(_this.balance.toString().indexOf(".")+1,_this.balance.toString().indexOf(".")+4);
					}
					if(longtermPower != 0) {
						_this.longtermPower = longtermPower * Math.pow(10, -3);
					}
					if(shorttermPower != 0) {
						_this.shorttermPower = shorttermPower * Math.pow(10, -3);
					}
				} else {
					console.log(err);
					this.$toast({
						message: '获取账户信息失败',
						position: 'middle',
						duration: 1000
					});
				}
			});
			_this.activeIndex = "-1";
		},
		getBobbyDetail: function() {
			var _this = this;
			_this.bobbyDetailList=[];
			info.userFinanceLogGet(localStorage.getItem('address'), 1, 100, function(err, result) {
				if(result) {
					var tokensChanges = result[0];
					var reasons = result[2];
					var times = result[3];
					if(tokensChanges.length == 0) {
						_this.$toast({
							message: '暂无明细',
							position: 'middle',
							duration: 1000
						});
						_this.activeIndex = "7";
						return;
					} else {
						for(var i = 0; i < tokensChanges.length; i++) {
							var obj = {};
//							Transfer, //普通转账
//							Mine, //挖矿所得
//							AssignToken, //分币
//							Vote, //投票
//							VoteCancel, //取消投票
//							VoteUnlock, //投票解锁
//							FinanceWithdraw, //
//							Tip //打赏
							switch(reasons[i].toNumber()) {
								case 0:
									obj.reason = '普通转账';
									break;
								case 1:
									obj.reason = '挖矿所得';
									break;
								case 2:
									obj.reason = 'Bobby';
									break;
								case 3:
									obj.reason = '投票';
									break;
								case 4:
									obj.reason = '取消投票';
									break;
								case 5:
									obj.reason = '投票解锁';
									break;
								case 6:
									obj.reason = '/';
									break;
								case 7:
									obj.reason = '打赏';
									break;
								default:
									obj.reason = '/';
							}
							obj.tokensChange = tokensChanges[i].toNumber()*Math.pow(10, -18);
							obj.tokensChange=obj.tokensChange.toFixed(4);
							var longTime = new Date(times[i] * 1000).getTime();
							var todayZero = new Date(new Date(new Date().toLocaleDateString()).getTime()-1).getTime();
							var yesterdayZero = todayZero-86400000;
							if(longTime>todayZero){
								obj.time = new Date(times[i] * 1000).format("hh:mm");
							}else if(longTime<todayZero && longTime>yesterdayZero){
								obj.time = "昨天 "+new Date(times[i] * 1000).format("hh:mm");
							}else{
								obj.time = new Date(times[i] * 1000).format("yyyy-MM-dd hh:mm:ss");
							}
							_this.bobbyDetailList.push(obj);
						}
						_this.bobbyDetailList=_this.bobbyDetailList.reverse();
					}
					_this.activeIndex = "7";
				} else {
					_this.$toast({
						message: '获取Bobby明细失败！',
						position: 'middle',
						duration: 1000
					});
				}
			});
		},
		getScoreDetail: function() {
			var _this = this;
			_this.$indicator.open({
				text: '加载中...'
			});
			_this.$http.post($base + '/rest/getUserPowerLog', {
				'userId': _this.userId,
				'page': 1,
				'pageSize': 1000
			}).then(function(res) {
				if(res.body.result == 1) {
					if(res.body.data.length != 0) {
						_this.scoreDetailList = res.body.data;
						console.log(_this.scoreDetailList.length);
						for(var i=0;i<_this.scoreDetailList.length;i++){
//							var powerChange=_this.scoreDetailList[i].powerChange*0.001;
//							_this.scoreDetailList[i].powerChange=powerChange.toFixed(4);
							var longTime = new Date(_this.scoreDetailList[i].time*1000).getTime();
							var todayZero = new Date(new Date(new Date().toLocaleDateString()).getTime()-1).getTime();
							var yesterdayZero = todayZero-86400000;
							if(longTime>todayZero){
								_this.scoreDetailList[i].timeDetail = new Date(_this.scoreDetailList[i].time * 1000).format("hh:mm");
							}else if(longTime<todayZero && longTime>yesterdayZero){
								_this.scoreDetailList[i].timeDetail = "昨天 "+new Date(_this.scoreDetailList[i].time * 1000).format("hh:mm");
							}else{
								_this.scoreDetailList[i].timeDetail = new Date(_this.scoreDetailList[i].time * 1000).format("yyyy-MM-dd hh:mm:ss");
							}
							var a=_this.scoreDetailList[i].reason;
							if(a==7){
								_this.scoreDetailList[i].reasons=_this.scoreDetailList[i].remark;
							}
							switch(a) {
								//Register, //注册
								//Invited, //被邀请，填写邀请码
								//Invitor, //邀请别人
								//LoginDaily, //每日登录
								//BindDapp, //绑定渠道
								//Certification, //实名认证
								//Watch, //观看视频
								//VideoWatched, //视频被观看
								//Upload, //上传视频
								//ShorttermPowerReset //临时算力重置
								case 0:
									_this.scoreDetailList[i].reasonDetail = '注册Bobby钱包';
									_this.scoreDetailList[i].title = 'Bobby钱包';
									break;
								case 1:
									_this.scoreDetailList[i].reasonDetail = '通过邀请码注册';
									_this.scoreDetailList[i].title = 'Bobby钱包';
									break;
								case 2:
									_this.scoreDetailList[i].reasonDetail = '邀请好友';
									_this.scoreDetailList[i].title = 'Bobby钱包';
									break;
								case 3:
									_this.scoreDetailList[i].reasonDetail = '每日登录';
									_this.scoreDetailList[i].title = 'Bobby钱包';
									break;
								case 4:
									_this.scoreDetailList[i].reasonDetail = '绑定渠道';
									_this.scoreDetailList[i].title = '绑定渠道';
									break;
								case 5:
									_this.scoreDetailList[i].reasonDetail = '实名认证';
									_this.scoreDetailList[i].title = '';
									break;
								case 6:
									_this.scoreDetailList[i].reasonDetail = '观看视频';
									_this.scoreDetailList[i].title = '亦播客';
									break;
								case 7:
									_this.scoreDetailList[i].reasonDetail = '视频被观看';
									_this.scoreDetailList[i].title = '亦播客';
									break;
								case 8:
									_this.scoreDetailList[i].reasonDetail = '上传视频';
									_this.scoreDetailList[i].title = '亦播客';
									break;
								default:
									scoreDetailList[i].reasonDetail = '/';
							}
						}
						//时间倒序
						_this.scoreDetailList=_this.scoreDetailList.reverse();
					}
				} else {
					_this.$toast({
						message: '获取算力明细失败！',
						position: 'middle',
						duration: 1000
					});
				}
				_this.$indicator.close();
				_this.activeIndex = "8";
			}).catch(function(){
				_this.$indicator.close();
				_this.$toast({
				  message: '获取算力明细失败，请刷新重试！',
				  position: 'middle',
				  duration: 2000
				});
			});
		},
		transfer:function(){
			//打赏转账
			var _this = this;
			
			_this.$toast({
				message: '打赏失败',
				position: 'middle',
				duration: 1000
			});
		},
		getRewardDetail:function(){
			//打赏记录
			var _this = this;
			_this.rewardRecordList=[];
			info.userTipLogGet(localStorage.getItem('address'), 1, 100, function(err, result) {
				if(result) {
					console.log(result);
//					var tokensChanges = result[0];
//					var reasons = result[2];
//					var times = result[3];
//					if(tokensChanges.length == 0) {
//						_this.$toast({
//							message: '暂无明细',
//							position: 'middle',
//							duration: 1000
//						});
//						_this.activeIndex = "0-1";
//						return;
//					} else {
//						for(var i = 0; i < tokensChanges.length; i++) {
//							var obj = {};
//							obj.tokensChange = tokensChanges[i].toNumber()*Math.pow(10, -18);
//							obj.tokensChange=obj.tokensChange.toFixed(4);
//							var longTime = new Date(times[i] * 1000).getTime();
//							var todayZero = new Date(new Date(new Date().toLocaleDateString()).getTime()-1).getTime();
//							var yesterdayZero = todayZero-86400000;
//							if(longTime>todayZero){
//								obj.time = new Date(times[i] * 1000).format("hh:mm");
//							}else if(longTime<todayZero && longTime>yesterdayZero){
//								obj.time = "昨天 "+new Date(times[i] * 1000).format("hh:mm");
//							}else{
//								obj.time = new Date(times[i] * 1000).format("yyyy-MM-dd hh:mm:ss");
//							}
//							_this.bobbyDetailList.push(obj);
//						}
//						_this.bobbyDetailList=_this.bobbyDetailList.reverse();
//					}
					_this.activeIndex = "0-1";
				} else {
					_this.$toast({
						message: '获取打赏记录失败！',
						position: 'middle',
						duration: 1000
					});
				}
			});
		},
		revisePassword:function(){
			var _this=this;
			//重置密码
			if(_this.password2==''||_this.password3==''){
				_this.$toast({
					message: '请输入交易密码',
					position: 'middle',
					duration: 1000
				});
				return;
			}
			if(_this.password2!=_this.password3){
				_this.$toast({
					message: '两次输入不一致',
					position: 'middle',
					duration: 1000
				});
				return;
			}
			_this.$http.post($base + '/rest/restPassword', {
				'userId': _this.userId,
				'password':_this.password2
			}).then(function(res) {
				console.log(res);
				if(res.body.result == 1) {
					_this.password=_this.password2;
					localStorage.setItem('password', _this.password); //交易密码
					_this.$toast({
						message: '密码修改成功',
						position: 'middle',
						duration: 1000
					});
					//回到转账页面
					_this.activeIndex=0;
				} else {
					_this.$toast({
						message: '获取数据失败！',
						position: 'middle',
						duration: 1000
					});
				}
				_this.$indicator.close();
			}).catch(function(){
				_this.$indicator.close();
				_this.$toast({
				  message: '获取数据失败，请刷新重试！',
				  position: 'middle',
				  duration: 2000
				});
			});
		},
		copy: function() {
			var _this = this;
			var clipboard = new ClipboardJS('.bobbyInfo .copyBtn');
			clipboard.on('success', function(e) {
				_this.$toast({
					message: '钱包地址复制成功 ',
					position: 'middle',
					duration: 1000
				});
			});
			clipboard.on('error', function(e) {
				console.log(e);
			});
		},
		copyPrivateKey: function() {
			var _this = this;
			var clipboard = new ClipboardJS('.exportKeyModal .copyBtn');
			clipboard.on('success', function(e) {
				_this.$toast({
					message: '复制成功 ',
					position: 'middle',
					duration: 1000
				});
			});
			clipboard.on('error', function(e) {
				console.log(e);
			});
		},
		copyMyInviteCode: function() {
			var _this = this;
			var clipboard = new ClipboardJS('.inviteModal .copyBtn');
			clipboard.on('success', function(e) {
				_this.$toast({
					message: '复制成功 ',
					position: 'middle',
					duration: 1000
				});
			});
			clipboard.on('error', function(e) {
				console.log(e);
			});
		},
		getUserDappBindInfo:function(){
			var _this = this;
			_this.$indicator.open({
				text: '加载中...'
			});
			_this.$http.post($base + '/rest/getUserDappBindInfo', {
				'userId': _this.userId
			}).then(function(res) {
				console.log(res);
				if(res.body.result == 1) {
					/*if(res.body.bindInfoData.length != 0) {
						var list = res.body.bindInfoData;
						for(var i=0;i<list.length;i++){
							var obj={};
							obj.name='绑定“'+list[i].name+"”公众号";
							obj.describe='绑定“'+list[i].name+"”公众号";
							if(list[i].binded==1){
								obj.isBind=true;
							}else{
								obj.isBind=false;
							}
							obj.scores='+'+list[i].powerAdd*Math.pow(10, -3);
							obj.showDetail=false;
							_this.thirdPatyList.unshift(obj);
						}
						
					}*/
//				{
//					name: '绑定“高尔夫频道”公众号',
//					scores: '+5',
//					isBind: true,
//					showDetail: false,
//					describe: '微信搜索“高尔夫频道”公众号'
//				}
				} else {
					_this.$toast({
						message: '获取数据失败！',
						position: 'middle',
						duration: 1000
					});
				}
				_this.$indicator.close();
				_this.activeIndex = "2";
			}).catch(function(){
				_this.$indicator.close();
				_this.$toast({
				  message: '获取数据失败，请刷新重试！',
				  position: 'middle',
				  duration: 2000
				});
			});
		},
		getInviteInfo:function(){
			var _this=this;
			_this.$indicator.open({
				text: '加载中...'
			});
			_this.$http.post($base + '/rest/getInviteInfo', {
				'userId': _this.userId,
				'page':1,
				'pageSize':100
			}).then(function(res) {
				console.log(res);
				if(res.body.result == 1) {
					if(res.body.inviteList.length != 0) {
						_this.inviteList = res.body.inviteList;
					}
				} else {
					_this.$toast({
						message: '获取用户邀请信息失败！',
						position: 'middle',
						duration: 1000
					});
				}
				_this.$indicator.close();
				_this.activeIndex = "3";
			}).catch(function(){
				_this.$indicator.close();
				_this.$toast({
				  message: '获取用户邀请信息失败！',
				  position: 'middle',
				  duration: 2000
				});
			});
		},
		removeLocalStorage: function() {
			var itemArr=['keystore','address','privateKey','password','tel','nickname','myInviteCode','userId'];
			for(var i in itemArr){
				localStorage.removeItem(itemArr[i]);
			}
		},
		test1:function(){
//			console.log(1)
//			var testObj={
//				from: "0x24e13bea85ee49e12e005faf42974486071cb603", 
//				to: "1b9cfa3a084936946c2c966e59b6aab230c76680",
//				value :100
//			}
//			web3.eth.sendTransaction(testObj , null);
		  //合约地址
//		  var address = "0x881d33f3ef9d293c697470593e9319493692768a";
		  var account = "0x5fa7bc87479f3e4e72092e4a7ea59bd16d55a422";
		  var privateKey = "2edc22b49ca89fc3ca21528c1f432a5357be635a3fe17072ad73aa010fd6cb55";
		  web3.eth.getTransactionCount(account, function (err, nonce) {
//		    var data = info.increment.getData("foo", 32);
		    var tx = new ethereumjs.Tx({
		      type: 0,
		      nonce: 0,
		      gasPrice: web3.toHex(web3.toWei('5', 'gwei')),
		      gasLimit: 100000,
		      to: "0x1b9cfa3a084936946c2c966e59b6aab230c76680",//收钱地址
		      value: 5,
		      data: '123',
		      extra:''
		    });
		    tx.sign(ethereumjs.Buffer.Buffer.from(privateKey, 'hex'));
		    var raw = '0x' + tx.serialize().toString('hex');
		    console.log(tx);
		    web3.eth.sendRawTransaction(raw, function (err, transactionHash) {
				
		    });
		  });
		},
		checkShowDetail: function() {
			
		},
		closeSideBar : function(){
			var _this = this;
			document.getElementById('sideBar').style.transform='translate3d(-100%, 0, 0)';
			document.getElementById('sideBar').style.visibility='hidden';
			document.getElementById('modal-overlay').removeEventListener('touchstart',function(){},false);
			setTimeout(function(){
				_this.activeIndex=-1;
			},100);
		},
		closeRightSideBar : function(){
			var _this = this;
			document.getElementById('rightSideBar').style.transform='translate3d(100%, 0, 0)';
			document.getElementById('rightSideBar').style.visibility='hidden';
			document.getElementById('modal-overlay').removeEventListener('touchstart',function(){},false);
			setTimeout(function(){
				_this.activeIndex=-1;
			},100);
		},
		test : function(){
			//测试函数
			
		}
	}
})