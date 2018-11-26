pragma solidity ^0.4.8;

import "./BokerCommon.sol";
import "./BokerDefine.sol";
import "./BokerManager.sol";
import "./BokerFile.sol";
import "./BokerUserData.sol";
import "./BokerTokenPower.sol";
import "./BokerLog.sol";

contract BokerUser is BokerManaged {
    using SafeMath for uint256;
    using AddressUtil for address;

    constructor(address addrManager) BokerManaged(addrManager) public {
    }

    /** @dev Handle user event.
    * @param addrFrom Address of from user.
    * @param addrTo Address of to user.
    * @param eventType Type of event.
    * @param eventValue1 Data1 of event.
    * @param eventValue2 Data2 of event.
    */
    function onUserEvent(
        address addrDapp, address addrFrom, uint256 eventType, address addrTo, uint256 timestamp, uint256 eventValue1, uint256 eventValue2) 
        external onlyContract{

        if (eventType >= uint256(UserEventType.End)) {
            return;
        }
        
        //ignore contract address, can't process
        if (addrFrom.isContract() || addrTo.isContract()) {
            //TODO error log
            return;
        }

        BokerTokenPower(contractAddress(ContractTokenPower)).onUserEvent(addrDapp, addrFrom, UserEventType(eventType), addrTo, timestamp, eventValue1, eventValue2);
        BokerFile(contractAddress(ContractFile)).onUserEvent(addrDapp, addrFrom, UserEventType(eventType), addrTo, timestamp, eventValue1, eventValue2);
    }

    /** @dev increase user short term power
    * @param addrUser Address of user.
    * @param value  Value of short term power.
    */
    function increaseUserShorttermPower(address addrUser, uint256 value, UserPowerReason reason) external onlyContract {
        //ignore contract address, can't process
        if (addrUser.isContract()) {
            return;
        }

        if (value == 0) {
            return;
        }

        BokerUserData userData = BokerUserData(contractAddress(ContractUserData));
        (, uint256 shorttemPowerOld) = userData.getUserPower(addrUser);
        uint256 shorttemPowerNew = shorttemPowerOld.add(value);
        userData.setUserShorttermPower(addrUser, shorttemPowerNew);
        
        //记录日志
        BokerLog(contractAddress(ContractLog)).userPowerLogAdd(
            addrUser, UserPowerType.Shortterm, shorttemPowerOld, shorttemPowerNew, reason);
    }

    /** @dev increase user long term power
    * @param addrUser Address of user.
    * @param value  Value of long term power.
    */
    function increaseUserLongtermPower(address addrUser, uint256 value, UserPowerReason reason) external onlyContract {
        //ignore contract address, can't process
        if (addrUser.isContract()) {
            return;
        }

        if (value == 0) {
            return;
        }

        BokerUserData userData = BokerUserData(contractAddress(ContractUserData));
        (uint256 longtermPowerOld,) = userData.getUserPower(addrUser);
        uint256 longtermPowerNew = longtermPowerOld.add(value);
        userData.setUserLongtermPower(addrUser, longtermPowerNew);
        
        //记录日志
        BokerLog(contractAddress(ContractLog)).userPowerLogAdd(
            addrUser, UserPowerType.Longterm, longtermPowerOld, longtermPowerNew, reason);
    }

    function clearUserShorttermPower(address addrUser, UserPowerReason reason) external onlyContract {
        //ignore contract address, can't process
        if (addrUser.isContract()) {
            return;
        }

        BokerUserData userData = BokerUserData(contractAddress(ContractUserData));
        (, uint256 shorttemPowerOld) = userData.getUserPower(addrUser);
        if (shorttemPowerOld == 0) {
            return;
        }
        userData.setUserShorttermPower(addrUser, 0);

        //记录日志
        BokerLog(contractAddress(ContractLog)).userPowerLogAdd(
            addrUser, UserPowerType.Shortterm, shorttemPowerOld, 0, reason);
    }

    /** @dev Get user.
    * @param addrUser Address of user.
    * @return balance balance of user
    * @return longtermPower long term power
    * @return shorttermPower short term power
    */
    function getUser(address addrUser) external view  returns(uint256 balance, uint256 longtermPower, uint256 shorttermPower) {
        (,,,,,,, longtermPower, shorttermPower) = BokerUserData(contractAddress(ContractUserData)).users(addrUser);
        return (addrUser.balance, longtermPower, shorttermPower);
    }
}