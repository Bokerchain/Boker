pragma solidity ^0.4.8;

import "../BokerCommon.sol";
import "./BokerAssignTokenData.sol";
import "./BokerAssignTokenDefine.sol";
import "./BokerAssignTokenEventHandler.sol";

contract BokerAssignTokenImpl is AccessControl, BokerAssignTokenDefine, BokerAssignTokenEventHandler, Log{
    using SafeMath for uint256;

    event AssignTokenBegin(uint256 tokensToAssgin, uint256 userCount, uint256 balance);
    event AssignToken(address indexed addrUser, uint256 amount);

    // Ensure that we are pointing to the right contract in our set call.
    bool public isBokerAssignTokenImpl = true;

    BokerAssignTokenData public data;

    /** @dev Constructor
    * @param dataAddress  Address of data contract.
    */
    constructor(address dataAddress) public {
        _setData(dataAddress);
    }

    function _setData(address dataAddress) private {
        BokerAssignTokenData dataContract = BokerAssignTokenData(dataAddress);

        // verify that a contract is what we expect
        require(dataContract.isBokerAssignTokenData());

        data = dataContract;
    }

    /** @dev Set address of data contract
    * @param dataAddress  Address of data contract.
    */
    function setData(address dataAddress) onlyCLevel external {
        _setData(dataAddress);
    }
    
    /** @dev Handover received tokens to data contract.
    */
    function handover() external payable {
        _logInfo("handover", msg.sender, address(0), msg.value, 0, 0, "");
        address(data).transfer(msg.value);
    }

    /** @dev Add user event.
    * @param ev Event data.
    */
    function _fireUserEvent(UserEvent memory ev) internal {
        UserEventHandler storage handler = eventHandlers[ev.eventType];
        if(!handler.exist){
            return;
        }
        
        if(!handler.check(data, ev)){
            return;
        }
    
        //calculate power;
        (uint256 powerFrom, uint256 powerTo) = handler.calculatePower(data, ev);
        _logTrace("_fireUserEvent->handler.calculatePower", ev.addrFrom, ev.addrTo, powerFrom, powerTo, 0, "");

        if(UserPowerType.Longterm == handler.powerType){
            data.increaseUserLongtermPower(ev.addrFrom, powerFrom);
            if(powerTo > 0){
                data.increaseUserLongtermPower(ev.addrTo, powerTo);
            }
        }
        else{
            data.increaseUserShorttermPower(ev.addrFrom, powerFrom);
            if(powerTo > 0){
                data.increaseUserShorttermPower(ev.addrTo, powerTo);
            }
        }

        //update user;
        handler.updateUser(data, ev);   
    }

    /** @dev Fire user event.
    * @param addrFrom Address of from user.
    * @param addrTo Address of to user.
    * @param eventType Type of event.
    * @param eventValue1 Data of event.
    * @param eventValue2 Data of event.
    */
    function fireUserEvent(
        uint256 eventType, address addrFrom, address addrTo, uint256 eventValue1, uint256 eventValue2) onlyCLevel external {
        if(addrFrom == address(0)){
            return;
        }
        _logTrace("fireUserEvent", addrFrom, addrTo, eventType, eventValue1, eventValue2, "");

        _fireUserEvent(UserEvent(eventType, addrFrom, addrTo, eventValue1, eventValue2));
    }

    /** @dev Assign token to users.
    */
    function assignToken() onlyCLevel external {
        _logDebug("assignToken", msg.sender, address(0), 0, 0, 0, "");

        uint256 tokensToAssgin = assginedTokensPerPeriodDefault.mul(data.assignCyclePeriod()).div(assgineTokenPeriodDefault);
        if(address(data).balance < tokensToAssgin) {
            tokensToAssgin = address(data).balance;
        }
        
        address[] memory addresses = data.getUsers();
        emit AssignTokenBegin(tokensToAssgin, addresses.length, address(data).balance);
        uint256 totalUserPower = data.getTotalUserPower();
        for(uint256 index = 0; index < addresses.length; index++){
            address addrUser = addresses[index];
            (uint256 longtermPower, uint256 shorttermPower) = data.getUserPower(addrUser);
            uint256 userPower = longtermPower + shorttermPower;
            uint256 tokensAssigned = userPower.mul(tokensToAssgin).div(totalUserPower);

            //clear user data and transger tokens to user.
            data.setUserShorttermPower(addrUser, 0);
            data.setUserWatchTime(addrUser, 0);
            data.setUserUploadCount(addrUser, 0);
            data.assignToken(addrUser, tokensAssigned);
        }
    }    
}