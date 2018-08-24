pragma solidity ^0.4.8;

import "./BokerAssignTokenDefine.sol";
import "./BokerAssignTokenData.sol";
import "../BokerCommon.sol";

contract BokerAssignTokenEventHandler is AccessControl, BokerAssignTokenDefine, Utility {
    using SafeMath for uint256;

    struct UserEventHandler {
        bool exist;
        UserPowerType powerType;

        /** @dev Check conditions before event is added.
        * @param data Address of BokerAssignTokenData contract.
        * @param ev Event data.
        * @param ev Event data.
        * @return true if can be added.
        */
        function(BokerAssignTokenData, UserEvent memory) view internal returns (bool) check;

        /** @dev Calculate user power
        * @param data Address of BokerAssignTokenData contract.
        * @param ev Event data.
        * @return powerFrom Power assigned to from user.
        * @return powerTo Power assigned to to user.
         */
        function(BokerAssignTokenData, UserEvent memory) view internal returns (uint256, uint256) calculatePower;

        /** @dev Update user info
        * @param data Address of BokerAssignTokenData contract.
        * @param ev Event data.
         */
        function(BokerAssignTokenData, UserEvent memory) internal updateUser;
    }

    // user event handlers
    mapping (uint256=>UserEventHandler) internal eventHandlers;

    constructor() public {
        _registerUserEventHandler(
            UserEventType.Register, UserPowerType.Longterm, _eventRegisterCheck, _eventRegisterCalculatePower, _eventRegisterUpdateUser);
        
        _registerUserEventHandler(
            UserEventType.LoginDaily, UserPowerType.Shortterm, _eventLoginDailyCheck, _eventLoginDailyCalculatePower, _eventLoginDailyUpdateUser);

        _registerUserEventHandler(
            UserEventType.Certification, UserPowerType.Longterm, _eventCertificationCheck, _eventCertificationCalculatePower, _eventCertificationUpdateUser);

        _registerUserEventHandler(
            UserEventType.BindWallet, UserPowerType.Longterm, _eventBindWalletCheck, _eventBindWalletCalculatePower, _eventBindWalletUpdateUser);

        _registerUserEventHandler(
            UserEventType.Watch, UserPowerType.Shortterm, _eventWatchCheck, _eventWatchCalculatePower, _eventWatchUpdateUser);

        _registerUserEventHandler(
            UserEventType.Upload, UserPowerType.Shortterm, _eventUploadCheck, _eventUploadCalculatePower, _eventUploadUpdateUser);

        uint256 channelType = 0;

        channelType = 1;
        BindWalletChannelConfig storage config = bindWalletChannelConfig[channelType];
        config.channelType = channelType;
        config.exist = true;
        config.entries.push(BindWalletChannelPowerEntry(10, 1));
        config.entries.push(BindWalletChannelPowerEntry(50, 3));
        config.entries.push(BindWalletChannelPowerEntry(100, 5));
        config.entries.push(BindWalletChannelPowerEntry(500, 10));
        config.entries.push(BindWalletChannelPowerEntry(1000, 15));
        config.entries.push(BindWalletChannelPowerEntry(5000, 20));
        config.entries.push(BindWalletChannelPowerEntry(10000, 30));
        config.entries.push(BindWalletChannelPowerEntry(-1, 40));

        channelType = 2;
        config = bindWalletChannelConfig[channelType];
        config.channelType = channelType;
        config.exist = true;
        config.entries.push(BindWalletChannelPowerEntry(10, 1));
        config.entries.push(BindWalletChannelPowerEntry(50, 3));
        config.entries.push(BindWalletChannelPowerEntry(100, 5));
        config.entries.push(BindWalletChannelPowerEntry(200, 10));
        config.entries.push(BindWalletChannelPowerEntry(500, 20));
        config.entries.push(BindWalletChannelPowerEntry(1000, 30));
        config.entries.push(BindWalletChannelPowerEntry(-1, 40));

        channelType = 3;
        config = bindWalletChannelConfig[channelType];
        config.channelType = channelType;
        config.exist = true;
        config.entries.push(BindWalletChannelPowerEntry(10, 1));
        config.entries.push(BindWalletChannelPowerEntry(50, 5));
        config.entries.push(BindWalletChannelPowerEntry(100, 10));
        config.entries.push(BindWalletChannelPowerEntry(200, 20));
        config.entries.push(BindWalletChannelPowerEntry(500, 30));
        config.entries.push(BindWalletChannelPowerEntry(-1, 40));
    }

    function _registerUserEventHandler(
        UserEventType eventType,
        UserPowerType powerType,
        function(BokerAssignTokenData, UserEvent memory) view internal returns (bool) canAdd,
        function(BokerAssignTokenData, UserEvent memory) view internal returns (uint256, uint256) calculatePower,
        function(BokerAssignTokenData, UserEvent memory) internal updateUser
    ) private {
        eventHandlers[uint256(eventType)] = UserEventHandler(true, powerType, canAdd, calculatePower, updateUser);
    }

    /***************************** Register ***********************************/
    uint256 public eventRegisterPowerAddDefault = 10;
    uint256 public eventInviteCodePowerAddDefault = 5;

    function _eventRegisterCheck(BokerAssignTokenData data, UserEvent memory ev) view internal returns (bool) {
        if(data.getUserRegisterTime(ev.addrFrom) > 0) {
            return false;
        }

        if(ev.addrFrom == ev.addrTo) {
            return false;
        }

        return true;
    }

    function _eventRegisterCalculatePower(BokerAssignTokenData data, UserEvent memory ev) view internal returns (uint256, uint256) {
        //addrFrom register, addrTo invitor if exist.
        uint256 powerFrom = eventRegisterPowerAddDefault;
        uint256 powerTo = 0;
        if(ev.addrTo != address(0)) {
            //register.
            powerFrom = powerFrom.add(eventInviteCodePowerAddDefault);

            //invitor.
            powerTo = powerTo.add(data.getInvitedFriendsCount(ev.addrTo) + 1).mul(2);
        }

        return (powerFrom, powerTo);
    }

    function _eventRegisterUpdateUser(BokerAssignTokenData data, UserEvent memory ev) internal {
        data.setUserRegisterTime(ev.addrFrom, now);
        data.addInvitedFriends(ev.addrTo);
    }

    /***************************** Register ***********************************/

    /***************************** Login Daily ***********************************/
    uint256 public eventLoginDailyPowerAddDefault = 1;

    function _eventLoginDailyCheck(BokerAssignTokenData data, UserEvent memory ev) view internal returns (bool) {
        if(!isSameDay(now, data.getLastLoginTime(ev.addrFrom))) {
            return true;
        }

        return false;
    }

    function _eventLoginDailyCalculatePower(BokerAssignTokenData data, UserEvent memory ev) view internal returns (uint256, uint256) {
        //in case warning.
        data;
        ev;

        return (eventLoginDailyPowerAddDefault, 0);
    }

    function _eventLoginDailyUpdateUser(BokerAssignTokenData data, UserEvent memory ev) internal {
        data.setUserLastLoginTime(ev.addrFrom, now);
    }

    /***************************** Login Daily ***********************************/

    /***************************** Certification ***********************************/
    uint256 public eventCertificationPowerAddDefault = 25;

    function _eventCertificationCheck(BokerAssignTokenData data, UserEvent memory ev) view internal returns (bool) {

        if(data.getUserCertificationTime(ev.addrFrom) > 0) {
            return false;
        }
        return true;
    }

    function _eventCertificationCalculatePower(BokerAssignTokenData data, UserEvent memory ev) view internal returns (uint256, uint256) {
        //in case warning.
        data;
        ev;

        return (eventCertificationPowerAddDefault, 0);
    }

    function _eventCertificationUpdateUser(BokerAssignTokenData data, UserEvent memory ev) internal {
        data.setUserCertificationTime(ev.addrFrom, now);
    }

    /***************************** Certification ***********************************/

    /***************************** Bind Wallet ***********************************/
    struct BindWalletChannelPowerEntry {
        int256 upper;      // upper limit.
        uint256 power;     // gained power.
    }

    struct BindWalletChannelConfig {
        uint256 channelType;
        bool exist;
        BindWalletChannelPowerEntry[] entries;
    }

    struct BindWalletChannel {
        uint256 channelId;
        uint256 channelType;
        int256 value;
    }

    mapping (uint256=>BindWalletChannelConfig) bindWalletChannelConfig;
    mapping (uint256=>BindWalletChannel) bindWalletChannels;

    function eventBindWalletGetConfig(uint256 channelType) view public returns (int256[] memory uppers, uint256[] memory powers) {
        BindWalletChannelConfig storage config = bindWalletChannelConfig[channelType];
        if(!config.exist){
            return;
        }

        uint256 len = config.entries.length;
        uppers = new int256[](len);
        powers = new uint256[](len);
        for (uint256 index = 0; index < len; index++) {
            BindWalletChannelPowerEntry storage entry = config.entries[index];
            uppers[index] = entry.upper;
            powers[index] = entry.power;
        }
    }

    function eventBindWalletSetConfig(uint256 channelType, uint256 index, int256 upper, uint256 power) onlyCEO public {
        BindWalletChannelConfig storage config = bindWalletChannelConfig[channelType];
        if(!config.exist){
            config.channelType = channelType;
            config.exist = true;
        }

        if(index >= config.entries.length) {
            config.entries.push(BindWalletChannelPowerEntry(upper, power));
        }
        else
        {
            config.entries[index].upper = upper;
            config.entries[index].power = power;
        }
    }

    function eventBindWalletGetChannel(uint256 channelId) view public returns (uint256 channelType, int256 value) {
        BindWalletChannel storage channel = bindWalletChannels[channelId];
        if(channel.channelId <= 0){
            return;
        }

        channelType = channel.channelType;
        value = channel.value;
    }

    function eventBindWalletSetChannel(uint256 channelId, uint256 channelType, int256 value) onlyCEO public {
        BindWalletChannel storage channel = bindWalletChannels[channelId];
        if(channel.channelId <= 0){
            channel.channelId = channelId;
            channel.channelType = channelType;
            channel.value = value;            
        }
        else {
            channel.channelType = channelType;
            channel.value = value;
        }
    }

    function eventBindWalletGetPower(uint256 channelId) view public returns (uint256){
        BindWalletChannel storage channel = bindWalletChannels[channelId];
        if(channel.channelId <= 0){
            return 0;
        }

        BindWalletChannelConfig storage config = bindWalletChannelConfig[channel.channelType];
        if(!config.exist){
            return 0;
        }

        int256 lower = 0;
        int256 upper = 0;
        for (uint256 index = 0; index < config.entries.length; index++) {
            BindWalletChannelPowerEntry storage entry = config.entries[index];
            upper = entry.upper;

            if(-1 == upper) {
                if(channel.value >= lower){
                    return entry.power;
                }
            }
            else {
                if((channel.value >= lower) && (channel.value < upper)){
                    return entry.power;
                }
            }

            lower = upper;
        }

        return 0;
    }

    function _eventBindWalletCheck(BokerAssignTokenData data, UserEvent memory ev) view internal returns (bool) {

        if(data.getUserBindWalletTime(ev.addrFrom) > 0) {
            return false;
        }
        return true;
    }

    function _eventBindWalletCalculatePower(BokerAssignTokenData data, UserEvent memory ev) view internal returns (uint256, uint256) {
        //in case warning.
        data;

        return (eventBindWalletGetPower(ev.eventValue1), 0);
    }

    function _eventBindWalletUpdateUser(BokerAssignTokenData data, UserEvent memory ev) internal {
        data.setUserBindWalletTime(ev.addrFrom, now);
    }

    /***************************** Bind Wallet ***********************************/

    /***************************** Watch ***********************************/
    function watchMax(BokerAssignTokenData data) view public returns (uint256) {
        return data.assignCyclePeriod();
    }

    //C=A(1-e^(-xt))，A=50.11859563，t=0.020154417
    function watchGetPower(uint256 value) pure public returns (uint256 power) {
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

    function _eventWatchCheck(BokerAssignTokenData data, UserEvent memory ev) view internal returns (bool) {
        if(data.getUserWatchTime(ev.addrFrom) >= watchMax(data)) {
            return false;
        }
        return true;
    }

    function _eventWatchCalculatePower(BokerAssignTokenData data, UserEvent memory ev) view internal returns (uint256, uint256) {
        //in case warning.
        data;
        suppressPureWarnings();

        // check if exceed max.
        uint256 max = watchMax(data);
        uint256 watchOrigin = data.getUserWatchTime(ev.addrFrom);
        if(watchOrigin.add(ev.eventValue1) > max){
            ev.eventValue1 = max.sub(watchOrigin);
        }

        // 80% to watcher, 20% to uploader.
        uint256 power = watchGetPower(ev.eventValue1);
        uint256 powerWatcher = power.mul(4).div(5);
        uint256 powerUploader = power.div(5);
        return (powerWatcher, powerUploader);
    }

    function _eventWatchUpdateUser(BokerAssignTokenData data, UserEvent memory ev) internal {
        uint256 max = watchMax(data);
        uint256 watchTime = data.getUserWatchTime(ev.addrFrom).add(ev.eventValue1);
        if(watchTime > max){
            watchTime = max;
        }
        data.setUserWatchTime(ev.addrFrom, watchTime);
    }
    /***************************** Watch ***********************************/

    /***************************** Upload ***********************************/
    uint256 public eventUploadCountMaxPer5Minutes = 50;

    function uploadMax(BokerAssignTokenData data) view public returns (uint256) {
        return (eventUploadCountMaxPer5Minutes.mul(data.assignCyclePeriod()).div(assgineTokenPeriodDefault));
    }

    //C=A(1-e^(-xt))，A=100.5336244，t=0.104770921
    function uploadGetPower(uint256 value) pure public returns (uint256 power) {
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

    function _eventUploadCheck(BokerAssignTokenData data, UserEvent memory ev) view internal returns (bool) {
        if(data.getUserUploadCount(ev.addrFrom) >= uploadMax(data)) {
            return false;
        }
        return true;
    }

    function _eventUploadCalculatePower(BokerAssignTokenData data, UserEvent memory ev) view internal returns (uint256, uint256) {
        //in case warning.
        data;
        suppressPureWarnings();

        // check if exceed max.
        uint256 max = uploadMax(data);
        uint256 uploadOrigin = data.getUserUploadCount(ev.addrFrom);
        if(uploadOrigin.add(ev.eventValue1) > max){
            ev.eventValue1 = max.sub(uploadOrigin);
        }

        return (uploadGetPower(ev.eventValue1), 0);
    }

    function _eventUploadUpdateUser(BokerAssignTokenData data, UserEvent memory ev) internal {
        uint256 max = uploadMax(data);
        uint256 uploadCount = data.getUserUploadCount(ev.addrFrom).add(ev.eventValue1);
        if(uploadCount > max){
            uploadCount = max;
        }
        data.setUserUploadCount(ev.addrFrom, uploadCount);
    }
    /***************************** Upload ***********************************/
}