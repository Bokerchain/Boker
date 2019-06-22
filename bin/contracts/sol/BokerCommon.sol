pragma solidity ^0.4.8;

//https://github.com/Arachnid/solidity-stringutils

/**
 * @title SafeMath
 * @dev safe math operations.
 */
library SafeMath {
    function mul(uint256 a, uint256 b) internal pure returns (uint256) {
        if (a == 0) {
            return 0;
        }

        uint256 c = a * b;
        assert(c / a == b);
        return c;
    }

    function div(uint256 a, uint256 b) internal pure returns (uint256) {
        return a / b;
    }

    function sub(uint256 a, uint256 b) internal pure returns (uint256) {
        assert(b <= a);
        return a - b;
    }

    function add(uint256 a, uint256 b) internal pure returns (uint256) {
        uint256 c = a + b;
        assert(c >= a);
        return c;
    }
}

/**
* @title Ownable
* @dev The Ownable contract has an owner address, and provides basic authorization control
* functions, this simplifies the implementation of "user permissions".
*/
contract Ownable {
    address public owner;

    event OwnershipTransferred(
        address indexed previousOwner,
        address indexed newOwner
    );

    /**
    * @dev The Ownable constructor sets the original `owner` of the contract to the sender
    * account.
    */
    constructor() public {
        owner = msg.sender;
    }

    /**
    * @dev Throws if called by any account other than the owner.
    */
    modifier onlyOwner() {
        require(msg.sender == owner);
        _;
    }

    /**
    * @dev Allows the current owner to transfer control of the contract to a newOwner.
    * @param _newOwner The address to transfer ownership to.
    */
    function transferOwnership(address _newOwner) public onlyOwner {
        _transferOwnership(_newOwner);
    }

    /**
    * @dev Transfers control of the contract to a newOwner.
    * @param _newOwner The address to transfer ownership to.
    */
    function _transferOwnership(address _newOwner) internal {
        require(_newOwner != address(0));
        emit OwnershipTransferred(owner, _newOwner);
        owner = _newOwner;
    }
}

contract RoleAccessControl is Ownable{
    using StringUtil for string;

    struct RoleEntry {
        address addr;
        bool enable;
    }
    struct Role {
        string name;
        mapping (address => RoleEntry) userMap;
        address[] userArray;
    }
    mapping (string =>  Role) roles;
    string[] public roleNames;
    // mapping (string =>  mapping (address => bool)) roles;

    event RoleAdd(address indexed user, string name);
    event RoleRemove(address indexed user, string name);

    /**
    * @dev add a role to a user
    * @param _name the name of the role
    * @param _user address of user
    */
    function roleAdd(string _name, address _user) public onlyOwner {
        if (_name.empty() || (_user == address(0))) {
            return;
        }

        Role storage role = roles[_name];
        if (role.name.empty()) {
            role.name = _name;
            roleNames.push(_name);
        }
        
        RoleEntry storage entry = role.userMap[_user];
        if (entry.addr == address(0)) {
            entry.addr = _user;
            role.userArray.push(_user);
        }
        entry.enable = true;
        emit RoleAdd(_user, _name);
    }

    /**
    * @dev remove a role from a user
    * @param _name the name of the role
    * @param _user address of user
    */
    function roleRemove(string _name, address _user) public onlyOwner {
        roles[_name].userMap[_user].enable = false;
        emit RoleRemove(_user, _name);
    }

    /**
    * @dev check if user has this role
    * @param _name the name of the role
    * @param _user address of user
    * @return bool
    */
    function _roleHas(string _name, address _user) internal view returns (bool){
        return roles[_name].userMap[_user].enable;
    }

    /**
    * @dev reverts if addr does not have role
    * @param _name the name of the role
    * @param _user address of user
    * // reverts
    */
    function roleCheck(string _name, address _user) public view {
        require(_roleHas(_name, _user));   
    }

    /**
    * @dev determine if addr has role
    * @param _name the name of the role
    * @param _user address of user
    * @return bool
    */
    function roleHas(string _name, address _user) public view returns (bool) {
        return _roleHas(_name, _user);
    }

    function roleUsers(string _name) public view returns (address[] memory addrs) {
        Role storage role = roles[_name];
        if (role.userArray.length == 0) {
            return;
        }

        uint256 len = 0;
        for (uint256 j = 0; j < role.userArray.length; j++) {
            RoleEntry storage entry = role.userMap[role.userArray[j]];
            if (entry.enable) {
                len++;
            }
        }
        addrs = new address[](len);
        uint256 index = 0;
        for (j = 0; j < role.userArray.length; j++) {
            entry = role.userMap[role.userArray[j]];
            if (entry.enable) {
                addrs[index++] = entry.addr;
            }
        }
    }

    /**
    * @dev modifier to scope access to a single role (uses msg.sender as addr)
    * @param _role the name of the role
    * // reverts
    */
    modifier onlyRole(string _role) {
        roleCheck(_role, msg.sender);
        _;
    }
}

contract Config is Ownable{
    using Uint256Util for uint256;
    using StringUtil for string;
    
    struct ConfigEntry {
        uint256 index;
        string name;
        string valueString;
        uint256 valueInt;
    }
    mapping (string=>ConfigEntry) configEntries;
    string[] public configNames;

    function configSetString(string name, string value) public onlyOwner {
        ConfigEntry storage entry = configEntries[name];
        if (entry.index <= 0) {
            entry.index = configNames.length;
            entry.name = name;
            configNames.push(name);
        }
        entry.valueString = value;
        entry.valueInt = 0;
    }

    function configGetString(string name) public view returns (string) {
        return configEntries[name].valueString;
    }

    function configSetInt(string name, uint256 value) public onlyOwner {
        ConfigEntry storage entry = configEntries[name];
        if (entry.index <= 0) {
            entry.index = configNames.length;
            entry.name = name;
            configNames.push(name);
        }
        entry.valueInt = value;
        entry.valueString = "";
    }

    function configGetInt(string name) public view returns (uint256) {
        return configEntries[name].valueInt;
    }

    function configInfo() public view returns (string) {
        string memory str = "";
        for (uint256 index = 0; index < configNames.length; index++) {
            ConfigEntry storage entry = configEntries[configNames[index]];
            string memory value = entry.valueString;
            if (value.empty()) {
                value = entry.valueInt.toString();
            }
            str = str.concat(entry.name, "=", value, "|");
        }
        return str;
    }
}

library TimeUtil {
    using SafeMath for uint256;

    function isSameDay(uint256 time1, uint256 time2) internal pure returns (bool) {
         //GMT +8
        uint256 timeZone = uint256(8).mul(3600);
        uint256 secondsOneDay = uint256(24).mul(3600);
        uint256 day1 = time1.add(timeZone).div(secondsOneDay);
        uint256 day2 = time2.add(timeZone).div(secondsOneDay);
        if(day1 == day2){
            return true;
        }
        
        return false;
    }
}

library Int256Util {
     function abs(int256 value) internal pure returns (int256) {
        if(value < 0) {
            value = -value;
        }
        return value;
    }
}

library Uint256Util {
     function toString(uint i) internal pure returns (string){
        if (i == 0) return "0";
        uint j = i;
        uint len;
        while (j != 0){
            len++;
            j /= 10;
        }
        bytes memory bstr = new bytes(len);
        uint k = len - 1;
        while (i != 0){
            bstr[k--] = byte(48 + i % 10);
            i /= 10;
        }
        return string(bstr);
    }

    /** @dev Calculate the difference between a and b 
    * @param a Value a.
    * @param b Value b.
    * @return c always greater or equal 0.
    */
    function diff(uint256 a, uint256 b) internal pure returns (uint256) {
        if(a > b){
            return a - b;
        }
        else {
            return b - a;
        }
    }
}

library PageUtil {
    using SafeMath for uint256;
    
    function pageRange(uint256 total, uint256 page, uint256 pageSize) internal pure returns (uint256 start, uint256 end) {
        if(total <= 0) {
            return (0, 0);
        }

        if(page <= 0) {
            page = 1;
        }

        if(pageSize == 0) {
            start = 0;
            end = total - 1;
            return (start, end);
        }

	    start = (page - 1).mul(pageSize);
        if(start >= total) {
            page = total/pageSize + 1;
            if(0 == total%pageSize) {
                page = page - 1;
            }
            start = (page - 1).mul(pageSize);
        }

	    end = start + pageSize - 1;
        if(end >= total) {
            end = total - 1;
        }

	    return (start, end);
    }
}

library AddressUtil {
    function codeSize(address _addr) internal view returns(uint256 _size) {
        assembly {  _size := extcodesize(_addr) }
    }

    function isContract(address _addr) internal view returns (bool) {
        uint256 size;
        assembly { size := extcodesize(_addr) }
        return size > 0;
    }
}

library StringUtil {

    function empty(string _a) internal pure returns (bool) {
        return (bytes(_a).length == 0);
    }

    function compare(string _a, string _b) internal pure returns (int) {
        bytes memory a = bytes(_a);
        bytes memory b = bytes(_b);
        uint minLength = a.length;
        if (b.length < minLength) minLength = b.length;
        for (uint i = 0; i < minLength; i ++)
            if (a[i] < b[i])
                return -1;
            else if (a[i] > b[i])
                return 1;
        if (a.length < b.length)
            return -1;
        else if (a.length > b.length)
            return 1;
        else
            return 0;
    }

    function indexOf(string _haystack, string _needle) internal pure returns (int) {
        bytes memory h = bytes(_haystack);
        bytes memory n = bytes(_needle);
        if(h.length < 1 || n.length < 1 || (n.length > h.length))
            return -1;
        else if(h.length > (2**128 -1))
            return -1;
        else
        {
            uint subindex = 0;
            for (uint i = 0; i < h.length; i ++)
            {
                if (h[i] == n[0])
                {
                    subindex = 1;
                    while(subindex < n.length && (i + subindex) < h.length && h[i + subindex] == n[subindex])
                    {
                        subindex++;
                    }
                    if(subindex == n.length)
                        return int(i);
                }
            }
            return -1;
        }
    }

    function concat(string _a, string _b, string _c, string _d, string _e) internal pure returns (string) {
        bytes memory _ba = bytes(_a);
        bytes memory _bb = bytes(_b);
        bytes memory _bc = bytes(_c);
        bytes memory _bd = bytes(_d);
        bytes memory _be = bytes(_e);
        string memory abcde = new string(_ba.length + _bb.length + _bc.length + _bd.length + _be.length);
        bytes memory babcde = bytes(abcde);
        uint k = 0;
        for (uint i = 0; i < _ba.length; i++) babcde[k++] = _ba[i];
        for (i = 0; i < _bb.length; i++) babcde[k++] = _bb[i];
        for (i = 0; i < _bc.length; i++) babcde[k++] = _bc[i];
        for (i = 0; i < _bd.length; i++) babcde[k++] = _bd[i];
        for (i = 0; i < _be.length; i++) babcde[k++] = _be[i];
        return string(babcde);
    }

    function concat(string _a, string _b, string _c, string _d) internal pure returns (string) {
        return concat(_a, _b, _c, _d, "");
    }

    function concat(string _a, string _b, string _c) internal pure returns (string) {
        return concat(_a, _b, _c, "", "");
    }

    function concat(string _a, string _b) internal pure returns (string) {
        return concat(_a, _b, "", "", "");
    }

    // parseInt
    function parseInt(string _a) internal pure returns (uint) {
        return parseInt(_a, 0);
    }

    // parseInt(parseFloat*10^_b)
    function parseInt(string _a, uint _b) internal pure returns (uint) {
        bytes memory bresult = bytes(_a);
        uint mint = 0;
        bool decimals = false;
        for (uint i=0; i<bresult.length; i++){
            if ((bresult[i] >= 48)&&(bresult[i] <= 57)){
                if (decimals){
                   if (_b == 0) break;
                    else _b--;
                }
                mint *= 10;
                mint += uint(bresult[i]) - 48;
            } else if (bresult[i] == 46) decimals = true;
        }
        if (_b > 0) mint *= 10**_b;
        return mint;
    }
}

// contract Log is Ownable {

//     struct LogEntry {
//         uint256 time;
//         uint8 level;
//         string  key;
//         address from;
//         address to;
//         uint256 v1;
//         uint256 v2;
//         uint256 v3;
//         string remarks;
//     }
    
//     uint8 constant logLevelTrace = 0;
//     uint8 constant logLevelDebug = 1;
//     uint8 constant logLevelInfo = 2;
//     uint8 constant logLevelWarning = 3;
//     uint8 constant logLevelError = 4;
//     uint8 constant logLevelClose = 5;

//     LogEntry[] private _logs;

//     uint8 public logEnabled = 1;
//     uint8 public logLevel = logLevelTrace;
//     string public logKeyDefault = "default";

//     function enable(uint8 status) onlyOwner external {
//         logEnabled = status;
//     }

//     function setLevel(uint8 level) onlyOwner external {
//         logLevel = level;
//     }

//     function _log(uint8 level, string key, address from, address to, uint256 v1, uint256 v2, uint256 v3, string remarks) internal {
//         if(0 == logEnabled) {
//             return;
//         }

//         if(logLevel > level){
//             return;
//         }

//         LogEntry memory entry;
//         entry.level = level;
//         entry.key = key;
//         entry.time = now;
//         entry.from = from;
//         entry.to = to;
//         entry.v1 = v1;
//         entry.v2 = v2;
//         entry.v3 = v3;
//         entry.remarks = remarks;
//         _logs.push(entry);
//     }

//     function _logTrace(string key, address from, address to, uint256 v1, uint256 v2, uint256 v3, string remarks) internal {
//         _log(logLevelTrace, key, from, to, v1, v2, v3, remarks);
//     }

//     function _logDebug(string key, address from, address to, uint256 v1, uint256 v2, uint256 v3, string remarks) internal {
//         _log(logLevelDebug, key, from, to, v1, v2, v3, remarks);
//     }

//     function _logInfo(string key, address from, address to, uint256 v1, uint256 v2, uint256 v3, string remarks) internal {
//         _log(logLevelInfo, key, from, to, v1, v2, v3, remarks);
//     }

//     function _logWarning(string key, address from, address to, uint256 v1, uint256 v2, uint256 v3, string remarks) internal {
//         _log(logLevelWarning, key, from, to, v1, v2, v3, remarks);
//     }

//     function _logError(string key, address from, address to, uint256 v1, uint256 v2, uint256 v3, string remarks) internal {
//         _log(logLevelError, key, from, to, v1, v2, v3, remarks);
//     }

//     function clearLog() onlyOwner external {
//         delete _logs;
//     }

//     function getLogSize() external view returns(uint256 size){
//         size = _logs.length;
//     }

//     function getLog(uint256 _index) external view returns(uint8 level, uint time, string  key, address from, address to, uint256 v1, uint256 v2, uint256 v3, string remarks){
//         require(_index < _logs.length);
//         require(_index >= 0);
//         LogEntry storage entry = _logs[_index];
//         level = entry.level;
//         key = entry.key;
//         time = entry.time;
//         from = entry.from;
//         to = entry.to;
//         v1 = entry.v1;
//         v2 = entry.v2;
//         v3 = entry.v3;
//         remarks = entry.remarks;
//     }
// }