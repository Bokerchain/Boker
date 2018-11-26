pragma solidity ^0.4.8;

import "./BokerCommon.sol";
import "./BokerManager.sol";

contract Test {
    struct User {
        address addr;                               // address of user
        uint256 tokens;
    }
    mapping (address=>User) users;

}

contract BokerTest {
    using SafeMath for uint256;
    using Uint256Util for uint256;
    using AddressUtil for address;
    
    int256 public value;
    
    uint256 public change = 0;
    uint256 public fromToken = 0;
    uint256 public toToken = 0;
    
    // constructor(address addrManager) BokerManaged(addrManager) public {
    // }
    
    function () public payable {
        change = msg.value;
        fromToken = msg.sender.balance;
        toToken = this.balance;
    }
    
    function testReceive() public payable {
        change = msg.value;
        fromToken = msg.sender.balance;
        toToken = this.balance;
    }

    function testGrant(address addrUser, uint256 mount) public {
        addrUser.transfer(mount);
        change = mount;
        fromToken = this.balance;
        toToken = addrUser.balance;
    }

    function isContract(address addr) public view returns (bool) {
        return addr.isContract();
    } 

    function test(uint256 page) public pure returns (uint256) {
        if(page <= 0){
            page = 1;
        }
        
        string memory v = page.toString();
        
        return page;
    }
    
    function pageRange(uint256 total, uint256 page, uint256 pageSize) public pure returns (uint256 start, uint256 end) {
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

    function userFilesGet(address user, uint256 page) public view returns (
        uint256[] playCounts, uint256[] playTimes, uint256[] userCounts, uint256[] createTimes) {
        if(page <= 0){
            page = 1;
        }
    }
    
    function getBalance(address addrUser) public view returns (uint256) {
        return addrUser.balance;
    }
}