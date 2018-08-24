pragma solidity ^0.4.8;

contract Test{
    
    struct User {
        address addr;
        uint256 index;

        uint256 power;
    }

    mapping (address=>User) public users;
    address[] public userArray;

    constructor() public {
        userArray.push(address(0));
    }
    
    function() public payable {
        
    }

    function _findAddUser(address addr) private returns (User storage user){
        user = users[addr];
        if(user.index <= 0) {
            user.addr = addr;
            user.index = userArray.length;
            
            userArray.push(addr);
        }
        return user;
    }

    function addPower(address addr, uint256 power) public {
        User storage user = _findAddUser(addr);
        user.power = user.power + power * 4 / 5;
    }
    
    function balanceOf(address addr) view public returns (uint256) {
        return addr.balance;
    }
}