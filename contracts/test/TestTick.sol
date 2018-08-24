pragma solidity ^0.4.8;

import "../src/BokerTick.sol";

contract TestTick is Tickable{
    string public message;

    function tick()  public {
        message = "TestTick.tick()";
    }
    
    function messageClear() public returns (int) {
        int a = -1;
        return a;
    }
}