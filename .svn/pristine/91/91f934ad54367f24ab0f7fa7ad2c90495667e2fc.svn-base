pragma solidity ^0.4.8;

import "../../github.com/oraclize/ethereum-api/oraclizeAPI_0.5.sol";

contract VoteImpl is usingOraclize{
    string public message;

    function clearMessage() public {
        delete message;
    }

    function deposit(address addrVoter) external payable {
        // message = strConcat("deposit () called msg.value=", "");
        message = strConcat("deposit () called msg.value=", uint2str(msg.value), " balance=", uint2str(this.balance));
        addrVoter;
    }
}

contract Vote is usingOraclize{

    VoteImpl public impl;
    string public message;

    constructor(address implAddress) public {
        impl = VoteImpl(implAddress);
    }

    function deposit() external payable{
        message = strConcat("deposit () called msg.value=", uint2str(msg.value), " balance=", uint2str(this.balance));
        impl.deposit.value(msg.value)(msg.sender);
    }
}