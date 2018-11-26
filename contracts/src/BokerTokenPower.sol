pragma solidity ^0.4.8;

import "./BokerDefine.sol";
import "./BokerCommon.sol";
import "./BokerManager.sol";
import "./BokerUser.sol";
import "./BokerUserData.sol";
import "./BokerTokenPowerData.sol";
import "./BokerFileData.sol";
import "./BokerFinance.sol";

contract BokerTokenPower  is BokerManaged {
    using SafeMath for uint256;
    using TimeUtil for uint256;
    using AddressUtil for address;

    constructor(address addrManager) BokerManaged(addrManager) public {
    }    

    function onUserEvent(
        address addrDapp, address addrFrom, UserEventType eventType, address addrTo, uint256 timestamp, uint256 eventValue1, uint256 eventValue2) 
        public onlyContract returns (Error err) {

        //ignore contract address, can't process
        if (addrFrom.isContract() || addrTo.isContract()) {
            return Error.AddressIsContract;
        }
        
        if (eventType >= UserEventType.End) {
            err = Error.EventNotSupported;
        }else if (UserEventType.Register == eventType) {
            err = _onUserEventRegister(addrFrom, addrTo);
        }else if (UserEventType.LoginDaily == eventType) {
            err = _onUserEventLoginDaily(addrFrom);
        }else if (UserEventType.BindWallet == eventType) {
            err = _onUserEventBindWallet(addrDapp, addrFrom);
        }else if (UserEventType.Watch == eventType) {
            err = _onUserEventWatch(addrFrom, eventValue1, eventValue2);
        }else if (UserEventType.Upload == eventType) {
            err = _onUserEventUpload(addrFrom);
        }else if (UserEventType.Certification == eventType) {
            err = _onUserEventCertification(addrFrom);
        }
        if (err != Error.Ok) {
            return err;
        }
        
        return Error.Ok;
    }
    
    function _onUserEventRegister(address addrUser, address addrInvitor) private returns (Error err) {       
        BokerUserData userData = BokerUserData(contractAddress(ContractUserData));
        if(userData.getUserRegisterTime(addrUser) > 0) {
            return Error.AlreadyRegistered;
        }

        BokerUser user = BokerUser(contractAddress(ContractUser));
        user.increaseUserLongtermPower(addrUser, globalConfigInt(CfgRegisterPowerAdd), UserPowerReason.Register);
        uint256 invitorPowerAdd = 0;
        if(addrInvitor != address(0) && addrInvitor != addrUser) {
            user.increaseUserLongtermPower(addrUser, globalConfigInt(CfgInvitedPowerAdd), UserPowerReason.Invited);
            invitorPowerAdd = (userData.getInvitedFriendsCount(addrInvitor) + 1).mul(2);
            user.increaseUserLongtermPower(addrInvitor, invitorPowerAdd, UserPowerReason.Invitor);
            userData.addInvitedFriends(addrInvitor);
        }
        userData.setUserRegisterTime(addrUser, now);
        return Error.Ok;
    }

    function _onUserEventLoginDaily(address addrUser) private returns (Error err) {
        BokerUserData userData = BokerUserData(contractAddress(ContractUserData));
        uint256 lastLoginTime = userData.getLastLoginTime(addrUser);
        if(now.isSameDay(lastLoginTime)) {
            return Error.AlreadyDailyLogined;
        }

        BokerUser user = BokerUser(contractAddress(ContractUser));
        user.increaseUserLongtermPower(addrUser, globalConfigInt(CfgLoginDailyPowerAdd), UserPowerReason.LoginDaily);
        userData.setUserLastLoginTime(addrUser, now);
        return Error.Ok;
    }

    function _onUserEventBindWallet(address addrDapp, address addrUser) private returns (Error err) {
        BokerUserData userData = BokerUserData(contractAddress(ContractUserData));
        if(userData.getUserBindWalletTime(addrUser) > 0) {
            return Error.AlreadyBindWallet;
        }

        BokerTokenPowerData tokenPowerData = BokerTokenPowerData(contractAddress(ContractTokenPowerData));
        BokerUser user = BokerUser(contractAddress(ContractUser));
        user.increaseUserLongtermPower(addrUser, tokenPowerData.bindWalletGetPower(addrDapp), UserPowerReason.BindWallet);
        userData.setUserBindWalletTime(addrUser, now);
        return Error.Ok;
    }

    //C=A(1-e^(-xt))，A=50.11859563，t=0.020154417
    function watchGetPower(uint256 value) public view returns (uint256 power) {
        uint256 watchTimeMax = globalConfigInt(CfgAssignPeriod);
        if(value > watchTimeMax){
            value = watchTimeMax;
        }

        if((value>=1) && (value<=1)) {return 1;}
        else if((value>=2) && (value<=2)) {return 2;}
        else if((value>=3) && (value<=3)) {return 3;}
        else if((value>=4) && (value<=4)) {return 4;}
        else if((value>=5) && (value<=5)) {return 5;}
        else if((value>=6) && (value<=6)) {return 6;}
        else if((value>=7) && (value<=8)) {return 7;}
        else if((value>=9) && (value<=9)) {return 8;}
        else if((value>=10) && (value<=10)) {return 9;}
        else if((value>=11) && (value<=11)) {return 10;}
        else if((value>=12) && (value<=12)) {return 11;}
        else if((value>=13) && (value<=14)) {return 12;}
        else if((value>=15) && (value<=15)) {return 13;}
        else if((value>=16) && (value<=16)) {return 14;}
        else if((value>=17) && (value<=18)) {return 15;}
        else if((value>=19) && (value<=19)) {return 16;}
        else if((value>=20) && (value<=21)) {return 17;}
        else if((value>=22) && (value<=22)) {return 18;}
        else if((value>=23) && (value<=24)) {return 19;}
        else if((value>=25) && (value<=26)) {return 20;}
        else if((value>=27) && (value<=27)) {return 21;}
        else if((value>=28) && (value<=29)) {return 22;}
        else if((value>=30) && (value<=31)) {return 23;}
        else if((value>=32) && (value<=33)) {return 24;}
        else if((value>=34) && (value<=35)) {return 25;}
        else if((value>=36) && (value<=37)) {return 26;}
        else if((value>=38) && (value<=39)) {return 27;}
        else if((value>=40) && (value<=41)) {return 28;}
        else if((value>=42) && (value<=44)) {return 29;}
        else if((value>=45) && (value<=46)) {return 30;}
        else if((value>=47) && (value<=49)) {return 31;}
        else if((value>=50) && (value<=51)) {return 32;}
        else if((value>=52) && (value<=54)) {return 33;}
        else if((value>=55) && (value<=57)) {return 34;}
        else if((value>=58) && (value<=61)) {return 35;}
        else if((value>=62) && (value<=64)) {return 36;}
        else if((value>=65) && (value<=68)) {return 37;}
        else if((value>=69) && (value<=72)) {return 38;}
        else if((value>=73) && (value<=76)) {return 39;}
        else if((value>=77) && (value<=81)) {return 40;}
        else if((value>=82) && (value<=87)) {return 41;}
        else if((value>=88) && (value<=93)) {return 42;}
        else if((value>=94) && (value<=100)) {return 43;}
        else if((value>=101) && (value<=108)) {return 44;}
        else if((value>=109) && (value<=118)) {return 45;}
        else if((value>=119) && (value<=130)) {return 46;}
        else if((value>=131) && (value<=146)) {return 47;}
        else if((value>=147) && (value<=170)) {return 48;}
        else if((value>=171) && (value<=218)) {return 49;}
        else if((value>=219) && (value<=300)) {return 50;}
        else {return 0;}
    }

    function _onUserEventWatch(address addrUser, uint256 fileId, uint256 watchTime) private returns (Error err) {
        BokerUserData userData = BokerUserData(contractAddress(ContractUserData));
        userData.addUserWatchTime(addrUser, fileId, watchTime);
        return Error.Ok;
    }

    //C=A(1-e^(-xt))，A=100.5336244，t=0.104770921
    function uploadGetPower(uint256 value) public view returns (uint256 power) {
        uint256 uploadCountMax = globalConfigInt(CfgUploadCountMax);
        if(value > uploadCountMax){
            value = uploadCountMax;
        }

        if((value>=1) && (value<=1)) {return 10;}
        else if((value>=2) && (value<=2)) {return 19;}
        else if((value>=3) && (value<=3)) {return 27;}
        else if((value>=4) && (value<=4)) {return 34;}
        else if((value>=5) && (value<=5)) {return 41;}
        else if((value>=6) && (value<=6)) {return 47;}
        else if((value>=7) && (value<=7)) {return 52;}
        else if((value>=8) && (value<=8)) {return 57;}
        else if((value>=9) && (value<=9)) {return 61;}
        else if((value>=10) && (value<=10)) {return 65;}
        else if((value>=11) && (value<=11)) {return 69;}
        else if((value>=12) && (value<=12)) {return 72;}
        else if((value>=13) && (value<=13)) {return 75;}
        else if((value>=14) && (value<=14)) {return 77;}
        else if((value>=15) && (value<=15)) {return 80;}
        else if((value>=16) && (value<=16)) {return 82;}
        else if((value>=17) && (value<=17)) {return 84;}
        else if((value>=18) && (value<=18)) {return 85;}
        else if((value>=19) && (value<=19)) {return 87;}
        else if((value>=20) && (value<=20)) {return 88;}
        else if((value>=21) && (value<=21)) {return 89;}
        else if((value>=22) && (value<=22)) {return 91;}
        else if((value>=23) && (value<=24)) {return 92;}
        else if((value>=25) && (value<=25)) {return 93;}
        else if((value>=26) && (value<=26)) {return 94;}
        else if((value>=27) && (value<=28)) {return 95;}
        else if((value>=29) && (value<=30)) {return 96;}
        else if((value>=31) && (value<=33)) {return 97;}
        else if((value>=34) && (value<=37)) {return 98;}
        else if((value>=38) && (value<=43)) {return 99;}
        else if((value>=44) && (value<=50)) {return 100;}
        else {return 0;}
    }

    function _onUserEventUpload(address addrUser) private returns (Error err) {
        BokerUserData userData = BokerUserData(contractAddress(ContractUserData));
        userData.addUserUploadCount(addrUser, 1);
        return Error.Ok;
    }

    function _onUserEventCertification(address addrUser) private returns (Error err) {       
        BokerUserData userData = BokerUserData(contractAddress(ContractUserData));

        if(userData.getUserCertificationTime(addrUser) > 0) {
            return Error.AlreadyCertificated;
        }

        BokerUser user = BokerUser(contractAddress(ContractUser));
        user.increaseUserLongtermPower(addrUser, globalConfigInt(CfgCertificationPowerAdd), UserPowerReason.Certification);
        userData.setUserCertificationTime(addrUser, now);
        return Error.Ok;
    }

    function _assignCycleIsEnd() private view returns (bool) {
        if(BokerTokenPowerData(contractAddress(ContractTokenPowerData)).assignCycleBegin().add(globalConfigInt(CfgAssignPeriod)) <= now) {
            return true;
        }
        return false;
    }

    /** @dev check if time to assgin token.
    * @return If true need invode assgin token.
    */
    function checkAssignToken() external view  returns (bool) {
        if(_assignCycleIsEnd()) {
            return true;
        }

        return false;
    }

    function _calculateWatchPower(address addrUser) private {
        BokerUserData userData = BokerUserData(contractAddress(ContractUserData));
        (uint256 watchTotal, uint256[] memory fileIds, uint256[] memory watchTimes) = userData.getUserWatchTime(addrUser);

        BokerFileData fileData = BokerFileData(contractAddress(ContractFileData));
        BokerUser user = BokerUser(contractAddress(ContractUser));
        uint256 powerWatch = watchGetPower(watchTotal);
        uint256 powerWatcher = 0;
        for (uint256 i = 0; i < fileIds.length; i++) {
            uint256 powerFile = watchTimes[i].mul(powerWatch).div(watchTotal);
            // 80% to watcher, 20% to owner.       
            uint256 powerOwner = powerFile.mul(globalConfigInt(CfgPowerWatchOwnerRatio)).div(100);
            powerWatcher = powerWatcher.add(powerFile.sub(powerOwner));
            address addrOwner = fileData.fileOwnerGet(fileIds[i]);
            if(powerOwner > 0){
                user.increaseUserShorttermPower(addrOwner, powerOwner, UserPowerReason.VideoWatched);
            }
        }
        if(powerWatcher > 0){
            user.increaseUserShorttermPower(addrUser, powerWatcher, UserPowerReason.Watch);
        }
    }
    
    function _calculateUploadPower(address addrUser) private {
        BokerUserData userData = BokerUserData(contractAddress(ContractUserData));
        BokerUser user = BokerUser(contractAddress(ContractUser));
        uint256 uploadCount = userData.getUserUploadCount(addrUser);
        uint256 powerUpload = uploadGetPower(uploadCount);
        if(powerUpload > 0){
            user.increaseUserShorttermPower(addrUser, powerUpload, UserPowerReason.Upload);
        }
    }

    function _calculateUserPower() private {
        address[] memory addresses = BokerUserData(contractAddress(ContractUserData)).getUserAddresses();
        for(uint256 index = 0; index < addresses.length; index++){
            address addrUser = addresses[index];

            //calculate watch power
            _calculateWatchPower(addrUser);
            
            //calculate upload
            _calculateUploadPower(addrUser);
        } 
    }
    
    function _assignTokensGet() private view returns (uint256 longtermTokens, uint256 shorttermTokens) {
        uint256 tokensToAsign = globalConfigInt(CfgAssignTokenPerCycle);
        if(contractAddress(ContractFinance).balance < tokensToAsign) {
            tokensToAsign = contractAddress(ContractFinance).balance;
        }
        longtermTokens = tokensToAsign.mul(globalConfigInt(CfgAssignTokenLongtermRatio)).div(100);
        shorttermTokens = tokensToAsign.sub(longtermTokens);
    }

    function _userPowerGet() private view returns (uint256 longtermPowerTotal, uint256 shorttermPowerTotal) {
        BokerUserData userData = BokerUserData(contractAddress(ContractUserData));
        longtermPowerTotal = 0;
        shorttermPowerTotal = 0;
        address[] memory addresses = userData.getUserAddresses();
        for(uint256 index = 0; index < addresses.length; index++){
            (uint256 longtermPower, uint256 shorttermPower) = userData.getUserPower(addresses[index]);
            longtermPowerTotal = longtermPowerTotal.add(longtermPower);
            shorttermPowerTotal = shorttermPowerTotal.add(shorttermPower);
        }
    }

    function _assignToken(uint256 longtermTokens, uint256 shorttermTokens, uint256 longtermPowerTotal, uint256 shorttermPowerTotal) private {
        BokerUserData userData = BokerUserData(contractAddress(ContractUserData));
        BokerUser user = BokerUser(contractAddress(ContractUser));

        address[] memory addresses = userData.getUserAddresses();
        for(uint256 index = 0; index < addresses.length; index++){
            (uint256 longtermPower, uint256 shorttermPower) = userData.getUserPower(addresses[index]);

            uint256 tokensAssigned = 0;
            if (longtermTokens > 0 && longtermPowerTotal > 0) {
                tokensAssigned = tokensAssigned.add(longtermTokens.mul(longtermPower).div(longtermPowerTotal));
            }
            if (shorttermTokens > 0 && shorttermPowerTotal > 0) {
                tokensAssigned = tokensAssigned.add(shorttermTokens.mul(shorttermPower).div(shorttermPowerTotal));
            }

            //clear user data and transger tokens to user.
            user.clearUserShorttermPower(addresses[index], UserPowerReason.AssignTokenReset);
            userData.clearUserWatchTime(addresses[index]);
            userData.clearUserUploadCount(addresses[index]);

            if(tokensAssigned > 0 && !addresses[index].isContract()) {
                BokerFinance(contractAddress(ContractFinance)).grantTokenTo(addresses[index], tokensAssigned, uint256(FinanceReason.AssignToken));
            }
        }
    }

    /** @dev Assign token to users.
    */
    function assignToken() external onlyContract {      
        if(!_assignCycleIsEnd()) {
            return;
        }

        _calculateUserPower();
        (uint256 longtermTokens, uint256 shorttermTokens) = _assignTokensGet();
        (uint256 longtermPowerTotal, uint256 shorttermPowerTotal) = _userPowerGet();
        _assignToken(longtermTokens, shorttermTokens, longtermPowerTotal, shorttermPowerTotal);
        
        BokerTokenPowerData tokenPowerData = BokerTokenPowerData(contractAddress(ContractTokenPowerData));
        uint256 assignCycleBeginOld = tokenPowerData.assignCycleBegin();
        uint256 assignPeriod = globalConfigInt(CfgAssignPeriod);
        uint256 roundAdd = now.sub(assignCycleBeginOld).div(assignPeriod);
        tokenPowerData.setAssignCycleBegin(assignCycleBeginOld.add(roundAdd.mul(assignPeriod)));

        //add log
        BokerLog(contractAddress(ContractLog)).assignTokenLogAdd(longtermTokens, shorttermTokens, longtermPowerTotal, shorttermPowerTotal);
    }
}