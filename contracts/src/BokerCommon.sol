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

/**
 *  Boker utility
 */
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


contract Log is AccessControl {

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
    
    uint8 constant Trace = 0;
    uint8 constant Debug = 1;
    uint8 constant Info = 2;
    uint8 constant Warning = 3;
    uint8 constant Error = 4;
    uint8 constant Close = 5;

    LogEntry[] private _logs;

    uint8 public enabled = 1;
    uint8 public logLevel = Trace;
    string public keyDefault = "default";

    function enable(uint8 status) onlyCLevel external {
        enabled = status;
    }

    function setLevel(uint8 level) onlyCLevel external {
        logLevel = level;
    }

    function _log(uint8 level, string key, address from, address to, uint256 v1, uint256 v2, uint256 v3, string remarks) internal {
        if(0 == enabled) {
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
        _log(Trace, key, from, to, v1, v2, v3, remarks);
    }

    function _logDebug(string key, address from, address to, uint256 v1, uint256 v2, uint256 v3, string remarks) internal {
        _log(Debug, key, from, to, v1, v2, v3, remarks);
    }

    function _logInfo(string key, address from, address to, uint256 v1, uint256 v2, uint256 v3, string remarks) internal {
        _log(Info, key, from, to, v1, v2, v3, remarks);
    }

    function _logWarning(string key, address from, address to, uint256 v1, uint256 v2, uint256 v3, string remarks) internal {
        _log(Warning, key, from, to, v1, v2, v3, remarks);
    }

    function _logError(string key, address from, address to, uint256 v1, uint256 v2, uint256 v3, string remarks) internal {
        _log(Error, key, from, to, v1, v2, v3, remarks);
    }

    function clearLog() onlyCLevel external {
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