pragma solidity ^0.4.8;

contract BokerDefine {
    uint256 constant internal bobby = 1 ether;      //一个bobby币
    uint256 constant internal POWER = 1000;         //一个算力，

    //配置
    string constant CfgRegisterPowerAdd = "RegisterPowerAdd";
    string constant CfgInviteCountMax = "InviteCountMax";
    string constant CfgInvitedPowerAdd = "InvitedPowerAdd";
    string constant CfgInvitorPowerAdd = "InvitorPowerAdd";
    string constant CfgLoginDailyPowerAdd = "LoginDailyPowerAdd";
    string constant CfgCertificationPowerAdd = "CertificationPowerAdd";
    string constant CfgAssignPeriod = "AssignPeriod";    
    string constant CfgAssignTokenPerCycle = "AssignTokenPerCycle";
    string constant CfgAssignTokenTotal = "AssignTokenTotal";
    string constant CfgUploadCountMax = "UploadCountMax";
    string constant CfgPowerWatchOwnerRatio = "PowerWatchOwnerRatio";
    string constant CfgAssignTokenLongtermRatio = "AssignTokenLongtermRatio";

    string constant CfgVoteLockup = "VoteLockup";
    string constant CfgVoteUnlockPrecision = "VoteUnlockPrecision";
    string constant CfgVoteCyclePeriod = "VoteCyclePeriod";

    //角色
    string constant RoleContract = "contract";      //角色：合约
    string constant RoleAdmin = "admin";            //角色：管理
    string constant RoleDapp = "dapp";              //角色：渠道

    //合约
    string constant ContractManager = "BokerManager";
    string constant ContractDapp = "BokerDapp";
    string constant ContractDappData = "BokerDappData";
    string constant ContractFile = "BokerFile";
    string constant ContractFileData = "BokerFileData";
    string constant ContractFinance = "BokerFinance";
    string constant ContractLog = "BokerLog";
    string constant ContractLogData = "BokerLogData";
    string constant ContractTokenPower = "BokerTokenPower";
    string constant ContractTokenPowerData = "BokerTokenPowerData";
    string constant ContractUser = "BokerUser";
    string constant ContractUserData = "BokerUserData";
    string constant ContractNode = "BokerNode";
    string constant ContractNodeData = "BokerNodeData";
    string constant ContractDataTransfer = "BokerDataTransfer";
    string constant ContractInterface = "BokerInterface";
    string constant ContractInterfaceBase = "BokerInterfaceBase";
    
    enum Error {
        Ok,
        AddressIsContract,
        EventNotSupported,
        AlreadyRegistered,
        AlreadyDailyLogined,
        AlreadyBindDapp,
        AlreadyCertificated
    }   

    enum UserEventType {
        Register,
        LoginDaily,
        BindDapp,
        Watch,
        Upload,
        Certification,       

        End                      // end of event type, event type value should less than End
    }

    enum UserPowerType {
        Longterm,
        Shortterm
    }

    enum UserPowerReason {
        Register,               //注册
        Invited,                //被邀请，填写邀请码
        Invitor,                //邀请别人
        LoginDaily,             //每日登录
        BindDapp,               //绑定渠道
        Certification,          //实名认证
        Watch,                  //观看视频
        VideoWatched,           //视频被观看
        Upload,                 //上传视频
        ShorttermPowerReset     //临时算力重置
    }

    enum FinanceReason {
        Transfer,               //普通转账
        Mine,                   //挖矿所得
        AssignToken,            //分币
        Vote,                   //投票
        VoteCancel,             //取消投票
        VoteUnlock,             //投票解锁
        FinanceWithdraw,        //
        Tip                     //打赏
    }

    enum VoteType {
        Vote,
        Cancel,
        Unlock
    }

    uint256 constant DappTypeInvalid = 0;
    uint256 constant DappTypeHardware = 1;
    uint256 constant DappTypeDApp = 2;
    uint256 constant DappTypeH5 = 3;
}