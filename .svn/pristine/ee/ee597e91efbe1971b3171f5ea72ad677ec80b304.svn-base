pragma solidity ^0.4.8;

import "./BokerCommon.sol";
import "./BokerAddressManager.sol";

contract Tickable {
    bool public tickStatus = true;      //状态 true 启用 false 停用
    uint public interval;               //tick时间间隔，单位：秒
    
    
    constructor() public {
    }

    /**
    *  禁用tick
    */
    function disableTick() public {
        tickStatus = false;
    }

    /**
    *  设置tick时间间隔
    */
    function setInterval(uint intvl) public {
        interval = intvl;
    }

    function tick()  public;
}

/**
 *  TODO need access Control
 */
contract BokerTick is AccessControl, Verifiable {

    Tickable[] tickArray;
    
    function getTickSize() view public returns (uint){
        return tickArray.length;
    }

    function addTick(address addr) onlyCEO public {
        require(addr != address(0), "invalid address");

        tickArray.push(Tickable(addr));
    }

    function deleteTick(uint index) onlyCEO public {
        require(index < tickArray.length, "index exceeds tickArray length");

        uint len = tickArray.length;
        for (uint i = index; i < len-1; i++) {
            tickArray[i] = tickArray[i+1];
        }

        delete tickArray[len-1];
        tickArray.length--;
    }

    function getTick(uint index) view public returns (address addr, uint interval, bool tickStatus) {
        require(index < tickArray.length, "index exceeds tickArray length");

        return (address(tickArray[index]), tickArray[index].interval(), tickArray[index].tickStatus());
    }

    function disableTickAt(uint index) onlyCEO public {
        require(index < tickArray.length, "index exceeds tickArray length");

        tickArray[index].disableTick();
    }

    function setTickInterval(uint index, uint interval) public {
        require(index < tickArray.length, "index exceeds tickArray length");

        tickArray[index].setInterval(interval);
    }

    function tick() onlyVerifier public {
        for (uint i = 0; i < tickArray.length; i++) {
            if(tickArray[i].tickStatus()) {
                tickArray[i].tick();
            }
        }
    }
}