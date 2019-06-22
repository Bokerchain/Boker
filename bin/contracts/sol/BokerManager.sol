pragma solidity ^0.4.8;

import "./BokerDefine.sol";
import "./BokerCommon.sol";

/**
 *  Boker Manager
 */
contract BokerManager is RoleAccessControl, BokerDefine, Config {
    using StringUtil for string;
    using Uint256Util for uint256;

    string public version = "1.0.1";
    uint256 public createTime = now;
    
    struct ContractInfo {
        uint256 idx;                                // index of conftract
        string name;                                // name of contract
        address addr;                               // address of contract
        string version;                             // version of contract
        uint256 createTime;                         // create time of contract       
    }

    // Ensure that we are pointing to the right contract in our set call.
    bool public isBokerManager = true;

    //contract addresses
    mapping (string => ContractInfo) contracts;
    string[] public contractNames;

    // Pausable
    bool public paused = false;
    
    event Pause();
    event Unpause();

    constructor() public {
        //init role
        roleAdd(RoleAdmin, msg.sender);

        //init config
        configSetInt(CfgRegisterPowerAdd, 10 * POWER);
        configSetInt(CfgInviteCountMax, 10);
        configSetInt(CfgInvitedPowerAdd, 5 * POWER);
        configSetInt(CfgInvitorPowerAdd, 2 * POWER);
        configSetInt(CfgLoginDailyPowerAdd, 1 * POWER);
        configSetInt(CfgCertificationPowerAdd, 25 * POWER);
        configSetInt(CfgUploadCountMax, 50);
        configSetInt(CfgAssignPeriod, 5 minutes);
        configSetInt(CfgAssignTokenPerCycle, 990 * bobby);
        configSetInt(CfgAssignTokenTotal, 4000000000 * bobby);
        configSetInt(CfgPowerWatchOwnerRatio, 20);
        configSetInt(CfgAssignTokenLongtermRatio, 20);

        configSetInt(CfgVoteLockup,3 days);
        configSetInt(CfgVoteUnlockPrecision, 5 minutes);
        configSetInt(CfgVoteCyclePeriod, 1 days);
    }

    function setContract(string cName, address addrContract) public onlyOwner {
        require(addrContract != address(0), "addrContract is 0!");
        BokerManaged managedContract = BokerManaged(addrContract);
        string memory managedVersion = managedContract.version();
        require(bytes(managedVersion).length > 0, "invalid contract address!");

        ContractInfo storage info = contracts[cName];
        if(info.idx == 0) {
            info.idx = contractNames.length;
            contractNames.push(cName);
        }
        info.name = cName;
        info.addr = addrContract;
        info.version = managedVersion;
        info.createTime = managedContract.createTime();

        roleAdd(RoleContract, addrContract);
    }

    function findContract(string cName) public view returns (uint256, string, address, string, uint256) {
        ContractInfo storage info = contracts[cName];
        return (info.idx, info.name, info.addr, info.version, info.createTime);
    }

    function getContractSize() public view returns (uint) {
        return contractNames.length;
    }

    function contractAddresses() public view returns (address[] addrs) {
        uint256 len = contractNames.length;
        addrs = new address[](len);
        for (uint256 index = 0; index < contractNames.length; index++) {
            ContractInfo storage info = contracts[contractNames[index]];
            addrs[index] = info.addr;
        }
    }

    function contractInfo() public view returns (string) {
        string memory str = "";
        for (uint256 index = 0; index < contractNames.length; index++) {
            ContractInfo storage info = contracts[contractNames[index]];
            str = str.concat(info.name, "(", info.version, "_");
            str = str.concat(info.createTime.toString(), ")", "|");
        }
        return str;
    }

    // called by owner to pause, triggers stopped state
    function pause() public onlyOwner whenNotPaused{
        paused = true;
        emit Pause();
    }

    // called by owner to unpause, returns to normal state
    function unpause() public onlyOwner whenPaused{
        paused = false;
        emit Unpause();
    }
    
    // modifier to allow actions only when the contract IS paused
    modifier whenNotPaused() {
        require(!paused, "paused!");
        _;
    }

    // modifier to allow actions only when the contract IS NOT paused
    modifier whenPaused {
        require(paused, "not paused!");
        _;
    }
}

contract BokerManaged is Ownable, BokerDefine{
    string public version = "1.0.1";
    uint256 public createTime = now;

    BokerManager public bokerManager;    

    constructor(address addrManager) public {
        _setManager(addrManager);
    }

    function _setManager(address addrManager) private {
        BokerManager managerContract = BokerManager(addrManager);
        require(managerContract.isBokerManager(), "not BokerManager!");
        bokerManager = managerContract;
    }

    function setManager(address addrManager) public onlyOwner {
        _setManager(addrManager);
    }

    function contractAddress(string cName) public view returns(address) {
        (, , address addr, , ) = bokerManager.findContract(cName);
        require(addr != address(0), "addr is 0!");
        return addr;
    }

    function globalConfigString(string key) public view returns (string) {
        return bokerManager.configGetString(key);
    }

    function globalConfigInt(string key) public view returns (uint256) {
        return bokerManager.configGetInt(key);
    }

    modifier onlyRole(string _role) {
        bokerManager.roleCheck(_role, msg.sender);
        _;
    }

    modifier onlyAdmin() {
        bokerManager.roleCheck(RoleAdmin, msg.sender);
        _;
    }

    modifier onlyContract() {
        require(bokerManager.roleHas(RoleAdmin, msg.sender) || bokerManager.roleHas(RoleContract, msg.sender), "not authorized!");
        _;
    }

    modifier onlyDapp() {
        require(bokerManager.roleHas(RoleAdmin, msg.sender) || bokerManager.roleHas(RoleDapp, msg.sender), "not authorized!");
        _;
    }

    // modifier to allow actions only when the contract IS paused
    modifier whenNotPaused() {
        require(!bokerManager.paused(), "paused!");
        _;
    }

    // modifier to allow actions only when the contract IS NOT paused
    modifier whenPaused {
        require(bokerManager.paused(), "not paused!");
        _;
    }
}