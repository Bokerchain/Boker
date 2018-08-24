pragma solidity ^0.4.8;

import "../BokerCommon.sol";
import "./BokerAssignTokenDefine.sol";
import "./BokerAssignTokenImpl.sol";

contract BokerAssignToken is AccessControl, BokerAssignTokenDefine, Log{

    BokerAssignTokenImpl public impl;

    // modifier of only implemented
    modifier onlyImplemented {
        require(address(impl) != address(0));
        _;
    }

    constructor() public {
    }

    function _setImpl(address implAddress) private {
        BokerAssignTokenImpl implContract = BokerAssignTokenImpl(implAddress);

        // verify that a contract is what we expect
        require(implContract.isBokerAssignTokenImpl());

        impl = implContract;
    }

    function setImpl(address implAddress) onlyCEO external {
        _setImpl(implAddress);
    }

    function () onlyImplemented payable public {
        _logInfo("()", msg.sender, address(0), msg.value, 0, 0, "");
        impl.handover.value(msg.value)();
    }

    /** @dev Fire user event.
    * @param addrFrom Address of from user.
    * @param addrTo Address of to user.
    * @param eventType Type of event.
    * @param eventValue1 Data1 of event.
    * @param eventValue2 Data2 of event.
    */
    function fireUserEvent(
        uint256 eventType, address addrFrom, address addrTo, uint256 eventValue1, uint256 eventValue2) whenNotPaused onlyImplemented onlyCLevel external {
        require(addrFrom != address(0));
        
        _logTrace("fireUserEvent", addrFrom, addrTo, eventType, eventValue1, eventValue2, "");
        
        impl.fireUserEvent(eventType, addrFrom, addrTo, eventValue1, eventValue2);
    }

    /** @dev check if time to assgin token.
    * @return If true need invode assgin token.
    */
    function checkAssignToken() external view returns (bool) {
        return impl.checkAssignToken();
    }

    /** @dev Assign token called periodically by chain.
    */
    function assignToken() whenNotPaused onlyImplemented external {
        _logDebug("assignToken", msg.sender, address(0), 0, 0, 0, "");

        impl.assignToken();
    }
}