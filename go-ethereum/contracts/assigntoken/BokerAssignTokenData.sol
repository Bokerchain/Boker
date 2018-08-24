pragma solidity ^0.4.8;

import "../BokerCommon.sol";
import "./BokerAssignTokenDefine.sol";

contract BokerAssignTokenData  is AccessControl, BokerAssignTokenDefine, Log {
    using SafeMath for uint256;

    // Ensure that we are pointing to the right contract in our set call.
    bool public isBokerAssignTokenData = true;

    // mapping of users from address to User structure.
    mapping(address => User) public users;
    address[] public userArray;

    // the begin time of current assign cycle;
    uint256 public assignCycleBegin = 0;

    // the period of assign cycle;
    uint256 public assignCyclePeriod = 5 minutes;

    constructor() public {
        //placeholder
        userArray.push(address(0));
    }

    function () payable public {
        _logInfo("()", msg.sender, address(0), msg.value, 0, 0, "");
    }

    /** @dev Set assign cycle begin time.
    * @param time Begin time of assign cycle.
    */
    function setAssignCycleBegin(uint256 time) onlyCLevel external {
        assignCycleBegin = time;
    }

    /** @dev Set assign cycle period.
    * @param period Period of assign cycle.
    */
    function setAssignCyclePeriod(uint256 period) onlyCLevel external {
        assignCyclePeriod = period;
    }

    /** @dev Get user address array.
    * @return addresses  Addresses of users.
    */
    function getUsers() view external returns (address[] memory) {

        uint256 len = userArray.length;
        address[] memory addresses = new address[](len - 1);
        for(uint index = 1; index < len; index++) {
            addresses[index - 1] = userArray[index];
        }
        return addresses;
    }

    function _findAddUser(address addr) private returns (User storage user){
        user = users[addr];
        if(user.index <= 0) {
            user.addr = addr;
            user.index = userArray.length;
            userArray.push(addr);
        }
        return user;
    }

    /** @dev Get user register time.
    * @return registerTime Register time.
    */
    function getUserRegisterTime(address addrUser) view external returns (uint256){
        User storage user = users[addrUser];
        return user.registerTime;
    }

    /** @dev Set user register time
    * @param addrUser Address of user.
    * @param registerTime  Register time.
    */
    function setUserRegisterTime(address addrUser, uint256 registerTime) onlyCLevel external {
        User storage user = _findAddUser(addrUser);
        user.registerTime = registerTime;
    }

    /** @dev Get user certification time.
    * @return certificationTime Time of real name certification.
    */
    function getUserCertificationTime(address addrUser) view external returns (uint256){
        User storage user = users[addrUser];
        return user.certificationTime;
    }

    /** @dev Set user certification time
    * @param addrUser Address of user.
    * @param certificationTime  Time of real name certification.
    */
    function setUserCertificationTime(address addrUser, uint256 certificationTime) onlyCLevel external {
        User storage user = _findAddUser(addrUser);
        user.certificationTime = certificationTime;
    }


    /** @dev Get user bind wallet time.
    * @return bindWalletTime Bind wallet time.
    */
    function getUserBindWalletTime(address addrUser) view external returns (uint256){
        User storage user = users[addrUser];
        return user.bindWalletTime;
    }

    /** @dev Set user bind wallet time.
    * @param addrUser Address of user.
    * @param bindWalletTime  Bind wallet time.
    */
    function setUserBindWalletTime(address addrUser, uint256 bindWalletTime) onlyCLevel external {
        User storage user = _findAddUser(addrUser);
        user.bindWalletTime = bindWalletTime;
    }

    /** @dev Get user last login time.
    * @param addrUser Address of user.
    * @return last login time.
    */
    function getLastLoginTime(address addrUser) view external returns (uint256){
        User storage user = users[addrUser];
        return user.lastLoginTime;
    }

    /** @dev Set user last login time
    * @param addrUser Address of user.
    * @param loginTime  login time.
    */
    function setUserLastLoginTime(address addrUser, uint256 loginTime) onlyCLevel external {
        User storage user = _findAddUser(addrUser);
        user.lastLoginTime = loginTime;
    }

    /** @dev Get user watch time.
    * @param addrUser Address of user.
    * @return watchTime Total watch time current assign cycle.
    */
    function getUserWatchTime(address addrUser) view external returns (uint256){
        User storage user = users[addrUser];
        return user.watchTime;
    }

    /** @dev Set user watch time.
    * @param addrUser Address of user.
    * @param watchTime  watch time.
    */
    function setUserWatchTime(address addrUser, uint256 watchTime) onlyCLevel external {
        User storage user = _findAddUser(addrUser);
        user.watchTime = watchTime;
    }

    /** @dev Get user upload count.
    * @param addrUser Address of user.
    * @return uploadCount Total upload count current assign cycle.
    */
    function getUserUploadCount(address addrUser) view external returns (uint256){
        User storage user = users[addrUser];
        return user.uploadCount;
    }

    /** @dev Set user upload count.
    * @param addrUser Address of user.
    * @param uploadCount Upload count.
    */
    function setUserUploadCount(address addrUser, uint256 uploadCount) onlyCLevel external {
        User storage user = _findAddUser(addrUser);
        user.uploadCount = uploadCount;
    }

    /** @dev Get invited user count.
    * @param addrUser Address of user.
    * @return invitedCount Count of invited users.
    */
    function getInvitedFriendsCount(address addrUser) view external returns (uint256){
        User storage user = users[addrUser];
        return user.invitedFriendsCount;
    }

    /** @dev Add invited user.
    * @param addrUser Address of user.
    * @return last login time.
    */
    function addInvitedFriends(address addrUser) onlyCLevel external {
        User storage user = _findAddUser(addrUser);
        user.invitedFriendsCount = user.invitedFriendsCount.add(1);
    }

    /** @dev Get total user power.
    * @return totalPower Total user power.
    */
    function getTotalUserPower() view external returns (uint256){
        uint256 len = userArray.length;
        uint256 totalPower = 0;
        for(uint index = 1; index < len; index++) {
            address addrUser = userArray[index];
            totalPower = totalPower.add(users[addrUser].longtermPower).add(users[addrUser].shorttermPower);
        }
        return totalPower;
    }

    /** @dev Get user power.
    * @param addrUser Address of user.
    * @return longtermPower Longterm Power.
    * @return shorttermPower Shortterm Power.
    */
    function getUserPower(address addrUser) view external returns (uint256, uint256){
        User storage user = users[addrUser];
        return (user.longtermPower, user.shorttermPower);
    }

    /** @dev Set user short term power
    * @param addrUser Address of user.
    * @param value  Value of short term power.
    */
    function setUserShorttermPower(address addrUser, uint256 value) onlyCLevel external {
        User storage user = _findAddUser(addrUser);
        user.shorttermPower = value;
    }

    /** @dev Increase user long term power
    * @param addrUser Address of user.
    * @param value  Value of long term power.
    */
    function increaseUserLongtermPower(address addrUser, uint256 value) onlyCLevel external {
        User storage user = _findAddUser(addrUser);
        user.longtermPower = user.longtermPower.add(value);
    }

    /** @dev Increase user short term power
    * @param addrUser Address of user.
    * @param value  Value of short term power.
    */
    function increaseUserShorttermPower(address addrUser, uint256 value) onlyCLevel external {
        User storage user = _findAddUser(addrUser);
        user.shorttermPower = user.shorttermPower.add(value);
    }

    /** @dev Assign token to user.
    * @param addrUser Address of user.
    * @param amount  amount of tokens.
    */
    function assignToken(address addrUser, uint256 amount) onlyCLevel external {
        addrUser.transfer(amount);
    }
}
