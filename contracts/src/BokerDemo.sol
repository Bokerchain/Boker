pragma solidity ^0.4.8;

import "./BokerManager.sol";


contract BokerDemo is BokerManaged {

    constructor(address addrManager) BokerManaged(addrManager) public {
    }
}