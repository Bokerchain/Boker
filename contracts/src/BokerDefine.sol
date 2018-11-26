pragma solidity ^0.4.8;

contract BokerDefine {
    uint256 constant internal bobby = 1 ether;

    //配置
    string constant CfgRegisterPowerAdd = "RegisterPowerAdd";
    string constant CfgInvitedPowerAdd = "InvitedPowerAdd";
    string constant CfgLoginDailyPowerAdd = "LoginDailyPowerAdd";
    string constant CfgCertificationPowerAdd = "CertificationPowerAdd";
    string constant CfgAssignPeriod = "AssignPeriod";
    string constant CfgUploadCountMax = "UploadCountMax";
    string constant CfgAssignTokenPerCycle = "AssignTokenPerCycle";
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
    string constant ContractDapp = "Dapp";
    string constant ContractDappData = "DappData";
    string constant ContractFile = "File";
    string constant ContractFileData = "FileData";
    string constant ContractFinance = "Finance";
    string constant ContractInterface = "Interface";
    string constant ContractInterfaceBase = "InterfaceBase";
    string constant ContractLog = "Log";
    string constant ContractLogData = "LogData";
    string constant ContractTokenPower = "TokenPower";
    string constant ContractTokenPowerData = "TokenPowerData";
    string constant ContractUser = "User";
    string constant ContractUserData = "UserData";
    string constant ContractNode = "Node";
    string constant ContractNodeData = "NodeData";
    
    enum Error {
        Ok,
        AddressIsContract,
        EventNotSupported,
        AlreadyRegistered,
        AlreadyDailyLogined,
        AlreadyBindWallet,
        AlreadyCertificated
    }   

    enum UserEventType {
        Register,
        LoginDaily,
        BindWallet,
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
        Register,
        Invited,
        Invitor,
        LoginDaily,
        BindWallet,
        Certification,
        Watch,
        VideoWatched,
        Upload,
        AssignTokenReset
    }

    enum FinanceReason {
        Transfer,
        Mine,
        AssignToken,
        Vote,
        VoteCancel,
        VoteUnlock,
        Withdraw
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