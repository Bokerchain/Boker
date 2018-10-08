pragma solidity ^0.4.8;

import "./BokerModuleData.sol";
import "../BokerCommon.sol";

contract BokerModuleImpl is AccessControl{
    using SafeMath for uint256;

    // Ensure that we are pointing to the right contract in our set call.
    bool public isBokerModuleImpl = true;    

    BokerModuleData public data;

    /** @dev
    * @param dataAddress  Address of data contract.
    */
    constructor(address dataAddress) public {
        _setData(dataAddress);
    }

    function _setData(address dataAddress) private {
        BokerModuleData dataContract = BokerModuleData(dataAddress);

        // verify that a contract is what we expect
        require(dataContract.isBokerModuleData());

        data = dataContract;
    }

    /** @dev Set address of data contract
    * @param dataAddress  Address of data contract.
    */
    function setData(address dataAddress) onlyCLevel external {
        _setData(dataAddress);
    }
}