pragma solidity ^0.4.8;

import "./BokerModuleImpl.sol";
import "../BokerCommon.sol";

contract BokerModule is AccessControl{

    BokerModuleImpl public impl;

    constructor(address implAddress) public {
        _setImpl(implAddress);
    }

    function _setImpl(address implAddress) private {
        BokerModuleImpl implContract = BokerModuleImpl(implAddress);

        // verify that a contract is what we expect
        require(implContract.isBokerModuleImpl());

        impl = implContract;
    }

    function setImpl(address implAddress) onlyCEO external {
        _setImpl(implAddress);
    }
}