//window.$base = 'http://10.10.50.245:8081';
//window.$base = 'http://192.168.1.198:8081';
//window.$base = 'http://10.10.50.249:8081';
window.$base = 'http://bobby-test.otvcloud.com:20000';
//window.$chain='http://101.132.132.167:8081';
//window.$serverIp = 'http://172.200.1.172:9001';//节点地址
window.$serverIp = 'http://172.200.2.195:9001';
//window.$serverIp='http://10.10.50.245:8545';
//window.$contract = '0x16f2f065B07EEE97Bf38d5761D10Be782fE38EBD';//合约地址
//window.$contract = '0x3Ccb288c1791786cF7037FdC95BbF108BD5D0E85';
//window.$contract = '0x3Ccb288c1791786cF7037FdC95BbF108BD5D0E85' ;
//window.$contract = '0x29bE144591B76a2cf8eEb25C2d742CCc1398155d';

function Bytes2Str(arr) {
	var str = "";
	for(var i = 0; i < arr.length; i++)    {
		var tmp = arr[i].toString(16);
		if(tmp.length == 1)       {
			tmp = "0" + tmp;
		}
		str += tmp;
	}
	return str;
}

function setCookie(name, value) {
	//	var Days = 30;
	//	var exp = new Date(new Date(new Date().toLocaleDateString()).getTime()+24*60*60*1000-1);
	//	console.log(exp.toGMTString());
	//	document.cookie = name + "="+ escape (value) + ";expires=" + exp.toGMTString();
	var curDate = new Date();

	//当前时间戳
	var curTamp = curDate.getTime();

	//当日凌晨的时间戳,减去一毫秒是为了防止后续得到的时间不会达到00:00:00的状态
	var curWeeHours = new Date(curDate.toLocaleDateString()).getTime() - 1;

	//当日已经过去的时间（毫秒）
	var passedTamp = curTamp - curWeeHours;

	//当日剩余时间
	var leftTamp = 24 * 60 * 60 * 1000 - passedTamp;
	var leftTime = new Date();
	leftTime.setTime(leftTamp + curTamp);
	//创建cookie
	document.cookie = name + "=" + escape(value) + ";expires=" + leftTime.toGMTString();
	console.log(leftTime.toGMTString());
}

function getCookie(name) {
	var arr, reg = new RegExp("(^| )" + name + "=([^;]*)(;|$)");
	if(arr = document.cookie.match(reg))
		return unescape(arr[2]);
	else
		return null;
}

Date.prototype.format = function(fmt)   
{ //author: meizz   
  var o = {   
    "M+" : this.getMonth()+1,                 //月份   
    "d+" : this.getDate(),                    //日   
    "h+" : this.getHours(),                   //小时   
    "m+" : this.getMinutes(),                 //分   
    "s+" : this.getSeconds(),                 //秒   
    "q+" : Math.floor((this.getMonth()+3)/3), //季度   
    "S"  : this.getMilliseconds()             //毫秒   
  }; 
  if(/(y+)/.test(fmt))   
    fmt=fmt.replace(RegExp.$1, (this.getFullYear()+"").substr(4 - RegExp.$1.length));   
  for(var k in o)   
    if(new RegExp("("+ k +")").test(fmt))   
  fmt = fmt.replace(RegExp.$1, (RegExp.$1.length==1) ? (o[k]) : (("00"+ o[k]).substr((""+ o[k]).length)));   
  return fmt;   
}  

var mixin = {
	data: {
		isShowDetail: false, //提升算例页 有没有展开信息介绍
		isShowVideoDetail: false, //是否展开"观看视频"的详细信息
		inviteNums: 3, //当前用户邀请人数
		/*0-5 主按钮、6-0:注册、6-1：昵称、6-2：导入私钥*/
		btnList: [{ //6个button 从右上开始 0-5
			className: 'reward',
			imgSrc: 'img/reward/reward@3x.png',
			imgSrcActive: 'img/reward/rewardActive@3x.png',
			btnName: '打赏转账',
			active: false
		}, {
			className: 'activity',
			imgSrc: 'img/activity/activity@3x.png',
			imgSrcActive: 'img/activity/activityActive@3x.png',
			btnName: '活动',
			active: false
		}, {
			className: 'score',
			imgSrc: 'img/score/score@3x.png',
			imgSrcActive: 'img/score/scoreActive@3x.png',
			btnName: '提升算力',
			active: false
		}, {
			className: 'invite',
			imgSrc: 'img/invite/invite@3x.png',
			imgSrcActive: 'img/invite/inviteActive@3x.png',
			btnName: '邀请奖励',
			active: false
		}, {
			className: 'user',
			imgSrc: 'img/user/headImg@3x.png',
			imgSrcActive: 'img/user/headImg@3x.png',
			btnName: '',
			active: false
		}, {
			className: 'community',
			imgSrc: 'img/community/community@3x.png',
			imgSrcActive: 'img/community/community@3x.png',
			btnName: '社区攻略',
			active: false
		}],
		//绑定的第三方平台的信息
		thirdPatyList: [
//		{
//			name: '绑定“"bobby wallet”公众号',
//			scores: '+10',
//			isBind: true,
//			showDetail: false,
//			describe: '微信搜索“高尔夫频道”公众号'
//		}, 
		{
			name: '绑定“高尔夫频道”公众号',
			scores: '+5',
			isBind: false,
			showDetail: false,
			describe: '微信搜索“高尔夫频道”公众号'
		}, {
			name: '绑定“掌上威视”公众号',
			scores: '+5',
			isBind: false,
			showDetail: false,
			describe: '微信搜索“掌上威视”公众号'
		}, {
			name: '绑定“亦播客”APP',
			scores: '+10',
			isBind: false,
			showDetail: false,
			describe: '下载打开“亦播客”APP'
		}, {
			name: '绑定“微投屏”APP',
			scores: '+10',
			isBind: false,
			showDetail: false,
			describe: '下载打开“微投屏”APP'
		}],
		//可兑换的商品列表
		goodsList: [{
			name: '京东E卡',
			price: 11880,
			surplus: 0,
			imgSrc: 'img/homePage/jd@3x.png',
			imgSrcNone: 'img/homePage/jdnone@3x.png'
		}, {
			name: '优酷季卡',
			price: 21880,
			surplus: 0,
			imgSrc: 'img/homePage/yk@3x.png',
			imgSrcNone: 'img/homePage/yknone@3x.png'
		}, {
			name: '芒果TV年卡',
			price: 31880,
			surplus: 0,
			imgSrc: 'img/homePage/mg@3x.png',
			imgSrcNone: 'img/homePage/mgnone@3x.png'
		}]
	}
}
var ABI = [{
	"constant": true,
	"inputs": [],
	"name": "getCandidates",
	"outputs": [{
		"name": "addresses",
		"type": "address[]"
	}, {
		"name": "tickets",
		"type": "uint256[]"
	}],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
}, {
	"constant": true,
	"inputs": [],
	"name": "generalInfoGet",
	"outputs": [{
		"name": "userTotal",
		"type": "uint256"
	}, {
		"name": "tokenAssigned",
		"type": "uint256"
	}, {
		"name": "tokenToAssign",
		"type": "uint256"
	}],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
}, {
	"constant": false,
	"inputs": [{
		"name": "uploader",
		"type": "address"
	}, {
		"name": "owner",
		"type": "address"
	}, {
		"name": "fileId",
		"type": "uint256"
	}, {
		"name": "ipfsHash",
		"type": "string"
	}, {
		"name": "ipfsUrl",
		"type": "string"
	}, {
		"name": "aliDnaFileId",
		"type": "string"
	}],
	"name": "addFile",
	"outputs": [],
	"payable": false,
	"stateMutability": "nonpayable",
	"type": "function"
}, {
	"constant": true,
	"inputs": [{
		"name": "cName",
		"type": "string"
	}],
	"name": "contractAddress",
	"outputs": [{
		"name": "",
		"type": "address"
	}],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
}, {
	"constant": false,
	"inputs": [{
		"name": "addrTo",
		"type": "address"
	}],
	"name": "giveTipTo",
	"outputs": [],
	"payable": true,
	"stateMutability": "payable",
	"type": "function"
}, {
	"constant": false,
	"inputs": [{
		"name": "addrCandidate",
		"type": "address"
	}],
	"name": "voteCandidate",
	"outputs": [],
	"payable": true,
	"stateMutability": "payable",
	"type": "function"
}, {
	"constant": false,
	"inputs": [{
		"name": "description",
		"type": "string"
	}, {
		"name": "team",
		"type": "string"
	}, {
		"name": "name",
		"type": "string"
	}],
	"name": "registerCandidate",
	"outputs": [],
	"payable": false,
	"stateMutability": "nonpayable",
	"type": "function"
}, {
	"constant": true,
	"inputs": [{
		"name": "key",
		"type": "string"
	}],
	"name": "globalConfigString",
	"outputs": [{
		"name": "",
		"type": "string"
	}],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
}, {
	"constant": true,
	"inputs": [],
	"name": "version",
	"outputs": [{
		"name": "",
		"type": "string"
	}],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
}, {
	"constant": false,
	"inputs": [{
		"name": "addrFrom",
		"type": "address"
	}, {
		"name": "eventTypes",
		"type": "uint256[]"
	}, {
		"name": "addrTos",
		"type": "address[]"
	}, {
		"name": "timestamps",
		"type": "uint256[]"
	}, {
		"name": "eventValue1s",
		"type": "uint256[]"
	}, {
		"name": "eventValue2s",
		"type": "uint256[]"
	}],
	"name": "fireUserEvents",
	"outputs": [],
	"payable": false,
	"stateMutability": "nonpayable",
	"type": "function"
}, {
	"constant": true,
	"inputs": [],
	"name": "createTime",
	"outputs": [{
		"name": "",
		"type": "uint256"
	}],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
}, {
	"constant": true,
	"inputs": [],
	"name": "bokerManager",
	"outputs": [{
		"name": "",
		"type": "address"
	}],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
}, {
	"constant": true,
	"inputs": [{
		"name": "addrUser",
		"type": "address"
	}],
	"name": "getUser",
	"outputs": [{
		"name": "balance",
		"type": "uint256"
	}, {
		"name": "longtermPower",
		"type": "uint256"
	}, {
		"name": "shorttermPower",
		"type": "uint256"
	}],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
}, {
	"constant": true,
	"inputs": [{
		"name": "user",
		"type": "address"
	}, {
		"name": "page",
		"type": "uint256"
	}, {
		"name": "pageSize",
		"type": "uint256"
	}],
	"name": "userFilesGet",
	"outputs": [{
		"name": "fileIds",
		"type": "uint256[]"
	}, {
		"name": "playCounts",
		"type": "uint256[]"
	}, {
		"name": "playTimes",
		"type": "uint256[]"
	}, {
		"name": "userCounts",
		"type": "uint256[]"
	}, {
		"name": "createTimes",
		"type": "uint256[]"
	}],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
}, {
	"constant": false,
	"inputs": [{
		"name": "uploader",
		"type": "address"
	}, {
		"name": "fileId",
		"type": "uint256"
	}],
	"name": "addUserFile",
	"outputs": [],
	"payable": false,
	"stateMutability": "nonpayable",
	"type": "function"
}, {
	"constant": false,
	"inputs": [{
		"name": "addrTo",
		"type": "address"
	}, {
		"name": "reason",
		"type": "uint256"
	}],
	"name": "transferTokenTo",
	"outputs": [],
	"payable": true,
	"stateMutability": "payable",
	"type": "function"
}, {
	"constant": false,
	"inputs": [{
		"name": "dappAddr",
		"type": "address"
	}, {
		"name": "dappType",
		"type": "uint256"
	}, {
		"name": "userCount",
		"type": "uint256"
	}, {
		"name": "monthlySales",
		"type": "uint256"
	}],
	"name": "dappSet",
	"outputs": [],
	"payable": false,
	"stateMutability": "nonpayable",
	"type": "function"
}, {
	"constant": true,
	"inputs": [],
	"name": "owner",
	"outputs": [{
		"name": "",
		"type": "address"
	}],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
}, {
	"constant": true,
	"inputs": [{
		"name": "addrUser",
		"type": "address"
	}, {
		"name": "page",
		"type": "uint256"
	}, {
		"name": "pageSize",
		"type": "uint256"
	}],
	"name": "userTipLogGet",
	"outputs": [{
		"name": "addrUsers",
		"type": "address[]"
	}, {
		"name": "tokensChanges",
		"type": "int256[]"
	}, {
		"name": "tokensAfters",
		"type": "uint256[]"
	}, {
		"name": "times",
		"type": "uint256[]"
	}],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
}, {
	"constant": true,
	"inputs": [{
		"name": "addrCandidate",
		"type": "address"
	}],
	"name": "getCandidate",
	"outputs": [{
		"name": "description",
		"type": "string"
	}, {
		"name": "team",
		"type": "string"
	}, {
		"name": "name",
		"type": "string"
	}, {
		"name": "tickets",
		"type": "uint256"
	}],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
}, {
	"constant": true,
	"inputs": [{
		"name": "addrUser",
		"type": "address"
	}],
	"name": "getUserBindDapp",
	"outputs": [{
		"name": "addrDapps",
		"type": "address[]"
	}, {
		"name": "bindeds",
		"type": "bool[]"
	}, {
		"name": "powerAdds",
		"type": "uint256[]"
	}],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
}, {
	"constant": true,
	"inputs": [{
		"name": "addrUser",
		"type": "address"
	}, {
		"name": "page",
		"type": "uint256"
	}, {
		"name": "pageSize",
		"type": "uint256"
	}],
	"name": "userFinanceLogGet",
	"outputs": [{
		"name": "tokensChanges",
		"type": "int256[]"
	}, {
		"name": "tokensAfters",
		"type": "uint256[]"
	}, {
		"name": "reasons",
		"type": "uint256[]"
	}, {
		"name": "times",
		"type": "uint256[]"
	}, {
		"name": "indexes",
		"type": "uint256[]"
	}],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
}, {
	"constant": true,
	"inputs": [{
		"name": "addrUser",
		"type": "address"
	}],
	"name": "getInvitedFriendsCount",
	"outputs": [{
		"name": "",
		"type": "uint256"
	}],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
}, {
	"constant": false,
	"inputs": [{
		"name": "addrManager",
		"type": "address"
	}],
	"name": "setManager",
	"outputs": [],
	"payable": false,
	"stateMutability": "nonpayable",
	"type": "function"
}, {
	"constant": true,
	"inputs": [{
		"name": "fileId",
		"type": "uint256"
	}],
	"name": "fileInfoGet",
	"outputs": [{
		"name": "uploader",
		"type": "address"
	}, {
		"name": "owner",
		"type": "address"
	}, {
		"name": "aliDnaFileId",
		"type": "string"
	}, {
		"name": "ipfsUrl",
		"type": "string"
	}, {
		"name": "playCount",
		"type": "uint256"
	}, {
		"name": "playTime",
		"type": "uint256"
	}, {
		"name": "userCount",
		"type": "uint256"
	}, {
		"name": "createTime",
		"type": "uint256"
	}],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
}, {
	"constant": true,
	"inputs": [{
		"name": "addrUser",
		"type": "address"
	}, {
		"name": "page",
		"type": "uint256"
	}, {
		"name": "pageSize",
		"type": "uint256"
	}],
	"name": "userPowerLogGet",
	"outputs": [{
		"name": "addrDapps",
		"type": "uint256[]"
	}, {
		"name": "powerTypes",
		"type": "uint256[]"
	}, {
		"name": "powerChanges",
		"type": "int256[]"
	}, {
		"name": "reasons",
		"type": "uint256[]"
	}, {
		"name": "param1s",
		"type": "uint256[]"
	}, {
		"name": "times",
		"type": "uint256[]"
	}],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
}, {
	"constant": true,
	"inputs": [],
	"name": "isBokerInterface",
	"outputs": [{
		"name": "",
		"type": "bool"
	}],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
}, {
	"constant": true,
	"inputs": [{
		"name": "key",
		"type": "string"
	}],
	"name": "globalConfigInt",
	"outputs": [{
		"name": "",
		"type": "uint256"
	}],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
}, {
	"constant": false,
	"inputs": [{
		"name": "_newOwner",
		"type": "address"
	}],
	"name": "transferOwnership",
	"outputs": [],
	"payable": false,
	"stateMutability": "nonpayable",
	"type": "function"
}, {
	"constant": false,
	"inputs": [],
	"name": "cancelAllVotes",
	"outputs": [],
	"payable": false,
	"stateMutability": "nonpayable",
	"type": "function"
}, {
	"constant": false,
	"inputs": [{
		"name": "addrFrom",
		"type": "address"
	}, {
		"name": "eventType",
		"type": "uint256"
	}, {
		"name": "addrTo",
		"type": "address"
	}, {
		"name": "timestamp",
		"type": "uint256"
	}, {
		"name": "eventValue1",
		"type": "uint256"
	}, {
		"name": "eventValue2",
		"type": "uint256"
	}],
	"name": "fireUserEvent",
	"outputs": [],
	"payable": false,
	"stateMutability": "nonpayable",
	"type": "function"
}, {
	"inputs": [{
		"name": "addrManager",
		"type": "address"
	}],
	"payable": false,
	"stateMutability": "nonpayable",
	"type": "constructor"
}, {
	"payable": true,
	"stateMutability": "payable",
	"type": "fallback"
}, {
	"anonymous": false,
	"inputs": [{
		"indexed": true,
		"name": "previousOwner",
		"type": "address"
	}, {
		"indexed": true,
		"name": "newOwner",
		"type": "address"
	}],
	"name": "OwnershipTransferred",
	"type": "event"
}]