pragma solidity ^0.4.8;


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

/**
 *  Boker define
 */
contract Define {
    uint public bobby = 1 ether;
}


contract Utility {
    using SafeMath for uint256;
    
    function abs(int256 value) internal pure returns (int256) {
        if(value < 0) {
            value = -value;
        }
        return value;
    }

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

    function suppressPureWarnings() view internal {
        now;
    }

    function strCompare(string _a, string _b) internal pure returns (int) {
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

    function strConcat(string _a, string _b, string _c, string _d, string _e) internal pure returns (string) {
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

    function strConcat(string _a, string _b, string _c, string _d) internal pure returns (string) {
        return strConcat(_a, _b, _c, _d, "");
    }

    function strConcat(string _a, string _b, string _c) internal pure returns (string) {
        return strConcat(_a, _b, _c, "", "");
    }

    function strConcat(string _a, string _b) internal pure returns (string) {
        return strConcat(_a, _b, "", "", "");
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

    function uint2str(uint i) internal pure returns (string){
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
}

// /**
//  *  Access Control by verifiers
//  */
// contract Verifiable {
//      //for test
//     function getVerifiers() pure public returns (address[]) {

//     }

//     // check if the address is one of verifiers
//     function isVerifier(address addr) pure public returns (bool) {
//         address[] memory verifiers = getVerifiers();
//         for(uint256 i=0; i<verifiers.length;i++) {
//             if (verifiers[i] == addr) {
//                 return true;
//             }
//         }
//         return false;
//     }

//     modifier onlyVerifier() {
//         require(isVerifier(msg.sender));
//         _;
//     }
// }

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

/**
 *  Access Control
 */
contract AccessControl {

    // The addresses of the accounts (or contracts) that can execute actions within each roles.
    address public ceoAddress;
    address public cfoAddress;
    address public cooAddress;

    // Pausable
    event Pause();
    event Unpause();

    bool public paused = false;

    constructor() public {
        ceoAddress = msg.sender;
    }

    // Access modifier for CEO-only functionality
    modifier onlyCEO() {
        require(msg.sender == ceoAddress);
        _;
    }

    // Access modifier for CFO-only functionality
    modifier onlyCFO() {
        require(msg.sender == cfoAddress);
        _;
    }

    // Access modifier for COO-only functionality
    modifier onlyCOO() {
        require(msg.sender == cooAddress);
        _;
    }

    modifier onlyCLevel() {
        require(
            msg.sender == cooAddress ||
            msg.sender == ceoAddress ||
            msg.sender == cfoAddress
        );
        _;
    }

    // Assigns a new address to act as the CEO. Only available to the current CEO.
    function setCEO(address _newCEO) external onlyCEO {
        require(_newCEO != address(0));

        ceoAddress = _newCEO;
    }

    // Assigns a new address to act as the CFO. Only available to the current CEO.
    function setCFO(address _newCFO) external onlyCEO {
        require(_newCFO != address(0));

        cfoAddress = _newCFO;
    }

    // Assigns a new address to act as the COO. Only available to the current CEO.
    function setCOO(address _newCOO) external onlyCEO {
        require(_newCOO != address(0));

        cooAddress = _newCOO;
    }

    // modifier to allow actions only when the contract IS paused
    modifier whenNotPaused() {
        require(!paused);
        _;
    }

    // modifier to allow actions only when the contract IS NOT paused
    modifier whenPaused {
        require(paused);
        _;
    }

    // called by ceo to pause, triggers stopped state
    function pause() public onlyCEO whenNotPaused returns (bool) {
        paused = true;
        emit Pause();
        return true;
    }

    // called by ceo to unpause, returns to normal state
    function unpause() public onlyCEO whenPaused returns (bool) {
        paused = false;
        emit Unpause();
        return true;
    }
}


contract Log is Ownable {

    struct LogEntry {
        uint256 time;
        uint8 level;
        string  key;
        address from;
        address to;
        uint256 v1;
        uint256 v2;
        uint256 v3;
        string remarks;
    }
    
    uint8 constant logLevelTrace = 0;
    uint8 constant logLevelDebug = 1;
    uint8 constant logLevelInfo = 2;
    uint8 constant logLevelWarning = 3;
    uint8 constant logLevelError = 4;
    uint8 constant logLevelClose = 5;

    LogEntry[] private _logs;

    uint8 public logEnabled = 1;
    uint8 public logLevel = logLevelTrace;
    string public logKeyDefault = "default";

    function enable(uint8 status) onlyOwner external {
        logEnabled = status;
    }

    function setLevel(uint8 level) onlyOwner external {
        logLevel = level;
    }

    function _log(uint8 level, string key, address from, address to, uint256 v1, uint256 v2, uint256 v3, string remarks) internal {
        if(0 == logEnabled) {
            return;
        }

        if(logLevel > level){
            return;
        }

        LogEntry memory entry;
        entry.level = level;
        entry.key = key;
        entry.time = now;
        entry.from = from;
        entry.to = to;
        entry.v1 = v1;
        entry.v2 = v2;
        entry.v3 = v3;
        entry.remarks = remarks;
        _logs.push(entry);
    }

    function _logTrace(string key, address from, address to, uint256 v1, uint256 v2, uint256 v3, string remarks) internal {
        _log(logLevelTrace, key, from, to, v1, v2, v3, remarks);
    }

    function _logDebug(string key, address from, address to, uint256 v1, uint256 v2, uint256 v3, string remarks) internal {
        _log(logLevelDebug, key, from, to, v1, v2, v3, remarks);
    }

    function _logInfo(string key, address from, address to, uint256 v1, uint256 v2, uint256 v3, string remarks) internal {
        _log(logLevelInfo, key, from, to, v1, v2, v3, remarks);
    }

    function _logWarning(string key, address from, address to, uint256 v1, uint256 v2, uint256 v3, string remarks) internal {
        _log(logLevelWarning, key, from, to, v1, v2, v3, remarks);
    }

    function _logError(string key, address from, address to, uint256 v1, uint256 v2, uint256 v3, string remarks) internal {
        _log(logLevelError, key, from, to, v1, v2, v3, remarks);
    }

    function clearLog() onlyOwner external {
        delete _logs;
    }

    function getLogSize() external view returns(uint256 size){
        size = _logs.length;
    }

    function getLog(uint256 _index) external view returns(uint8 level, uint time, string  key, address from, address to, uint256 v1, uint256 v2, uint256 v3, string remarks){
        require(_index < _logs.length);
        require(_index >= 0);
        LogEntry storage entry = _logs[_index];
        level = entry.level;
        key = entry.key;
        time = entry.time;
        from = entry.from;
        to = entry.to;
        v1 = entry.v1;
        v2 = entry.v2;
        v3 = entry.v3;
        remarks = entry.remarks;
    }
}

contract BokerDefine {

    uint256 constant internal bobby = 1 ether;
}