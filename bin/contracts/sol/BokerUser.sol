pragma solidity ^0.4.8;

import "./BokerCommon.sol";
import "./BokerDefine.sol";
import "./BokerManager.sol";
import "./BokerFile.sol";
import "./BokerUserData.sol";
import "./BokerTokenPower.sol";
import "./BokerTokenPowerData.sol";
import "./BokerLog.sol";
import "./BokerDapp.sol";
import "./BokerFinance.sol";

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
    function increaseUserShorttermPower(address addrDapp, address addrUser, uint256 value, UserPowerReason reason, uint256 param1) external onlyContract {
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
            addrDapp, addrUser, uint256(UserPowerType.Shortterm), int256(value), shorttemPowerNew, uint256(reason), param1);
    }

    /** @dev increase user long term power
    * @param addrUser Address of user.
    * @param value  Value of long term power.
    */
    function increaseUserLongtermPower(address addrDapp, address addrUser, uint256 value, UserPowerReason reason, uint256 param1) external onlyContract {
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
            addrDapp, addrUser, uint256(UserPowerType.Longterm), int256(value), longtermPowerNew, uint256(reason), param1);
    }

    function clearUserShorttermPower(address addrDapp, address addrUser, UserPowerReason reason, uint256 param1) external onlyContract {
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
            addrDapp, addrUser, uint256(UserPowerType.Shortterm), -1*int256(shorttemPowerOld), 0, uint256(reason), param1);
    }

    /** @dev Get user.
    * @param addrUser Address of user.
    * @return balance balance of user
    * @return longtermPower long term power
    * @return shorttermPower short term power
    */
    function getUser(address addrUser) external view  returns(uint256 balance, uint256 longtermPower, uint256 shorttermPower) {
        (,,,,,,,, longtermPower, shorttermPower) = BokerUserData(contractAddress(ContractUserData)).users(addrUser);
        return (addrUser.balance, longtermPower, shorttermPower);
    }

    /** @dev Get user count.
    * @return total count of user
    */
    function userCount() external view returns (uint256 total) {
        return BokerUserData(contractAddress(ContractUserData)).userCount();
    }

    /** @dev Get invited user count.
    * @param addrUser Address of user.
    * @return invitedCount Count of invited users.
    */
    function getInvitedFriendsCount(address addrUser) external view returns (uint256) {
        return BokerUserData(contractAddress(ContractUserData)).getInvitedFriendsCount(addrUser);
    }

    /** @dev transfer bobby to user, need gas.
    * @param addrTo Address of to user.
    * @param reason reason of transfer.
    */
    function transferTokenTo(address addrTo, uint256 reason) external payable {
        if (msg.value <= 0) {
            return;
        }
        addrTo.transfer(msg.value);

        //记录日志
        BokerLog(contractAddress(ContractLog)).financeLogAdd(
            msg.sender, addrTo, msg.value, msg.sender.balance, addrTo.balance, reason);
    }

    /** @dev give tip to user
    * @param addrTo Address of to user.
    */
    function giveTipTo(address addrTo) external payable {
        this.transferTokenTo(addrTo, uint256(FinanceReason.Tip));
    }


    /** @dev get user binded dapp info
    * @param addrUser Address of user.
    * @return addrDapps addresses of dapps.
    * @return bindeds if has bind.
    * @return powerAdds power to add if bind.
    */
    function getUserBindDapp(address addrUser) external view returns (address[] addrDapps, bool[] bindeds, uint256[] powerAdds) {
        addrDapps = BokerDapp(contractAddress(ContractDapp)).dappGetAdresses();
        uint256 len = addrDapps.length;
        bindeds = new bool[](len);
        powerAdds = new uint256[](len);
        for (uint256 index = 0; index < len; index++) {
            address addrDapp = addrDapps[index];
            bool binded = false;
            if (BokerUserData(contractAddress(ContractUserData)).getDappBind(addrUser, addrDapp) > 0) {
                binded = true;
            }
            bindeds[index] = binded;
            powerAdds[index] = BokerTokenPowerData(contractAddress(ContractTokenPowerData)).bindDappGetPower(addrDapp);
        }
    }

    /** @dev Set user.
    * @param addrUser address of user.
    * @param registerTime register time of user.
    * @param invitedFriendsCount total invited friends count.
    * @param bindedDappAddresses addresses of binded dapps.
    * @param bindedDappTimes binded times of binded dapps.
    * @param longtermPower longterm power of user.
    * @param shorttermPower shortterm power of user.
    */
    function setUser(
        address addrUser, uint256 registerTime, uint256 invitedFriendsCount, address[] bindedDappAddresses, uint256[] bindedDappTimes, uint256 longtermPower, uint256 shorttermPower) 
        external onlyContract {
        BokerUserData(contractAddress(ContractUserData)).setUser(addrUser, registerTime, invitedFriendsCount, bindedDappAddresses, bindedDappTimes, longtermPower, shorttermPower);
    }

}