pragma solidity ^0.4.8;

import "./BokerDefine.sol";
import "./BokerCommon.sol";
import "./BokerManager.sol";

contract BokerUserData  is BokerManaged {
    using SafeMath for uint256;
    using AddressUtil for address;

    struct User {
        address addr;                               // address of user
        uint256 index;                              // index in array of user
        uint256 registerTime;                       // time of register
        uint256 certificationTime;                  // time of real name certification
        uint256 bindWalletTime;                     // time of bind wallet
        uint256 lastLoginTime;                      
        uint256 invitedFriendsCount;

        mapping (uint256=>uint256) watchTime;       // watch time current assign period
        uint256[] watchFiles;                       // watch Files current assign period
        uint256 uploadCount;                        // upload count current assign period

        uint256 longtermPower;                      // long term power, always exists
        uint256 shorttermPower;                     // short term power, cleared at the end of every assgin cycle.
    }

    // mapping of users from address to User structure. 
    mapping(address => User) public users;
    address[] public userArray;

    
    constructor(address addrManager) BokerManaged(addrManager) public {
        //placeholder
        userArray.push(address(0));
    }

    function userCount() public view returns (uint256) {
        return userArray.length - 1;
    }

    function getUserAddresses() external view returns (address[] memory) {
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
    function getUserRegisterTime(address addrUser) public view returns (uint256){
        User storage user = users[addrUser];
        return user.registerTime;
    }

    /** @dev Set user register time
    * @param addrUser Address of user.
    * @param registerTime  Register time.
    */
    function setUserRegisterTime(address addrUser, uint256 registerTime) public onlyContract {
        //ignore contract address, can't process
        if (addrUser.isContract()) {
            return;
        }

        User storage user = _findAddUser(addrUser);
        user.registerTime = registerTime;
    }

    /** @dev Get user certification time.
    * @return certificationTime Time of real name certification.
    */
    function getUserCertificationTime(address addrUser) public view returns (uint256){
        User storage user = users[addrUser];
        return user.certificationTime;
    }

    /** @dev Set user certification time
    * @param addrUser Address of user.
    * @param certificationTime  Time of real name certification.
    */
    function setUserCertificationTime(address addrUser, uint256 certificationTime) public onlyContract {
        //ignore contract address, can't process
        if (addrUser.isContract()) {
            return;
        }

        User storage user = _findAddUser(addrUser);
        user.certificationTime = certificationTime;
    }


    /** @dev Get user bind wallet time.
    * @return bindWalletTime Bind wallet time.
    */
    function getUserBindWalletTime(address addrUser) public view returns (uint256){
        User storage user = users[addrUser];
        return user.bindWalletTime;
    }

    /** @dev Set user bind wallet time.
    * @param addrUser Address of user.
    * @param bindWalletTime  Bind wallet time.
    */
    function setUserBindWalletTime(address addrUser, uint256 bindWalletTime) public onlyContract {
        //ignore contract address, can't process
        if (addrUser.isContract()) {
            return;
        }

        User storage user = _findAddUser(addrUser);
        user.bindWalletTime = bindWalletTime;
    }

    /** @dev Get user last login time.
    * @param addrUser Address of user.
    * @return last login time.
    */
    function getLastLoginTime(address addrUser) public view returns (uint256){
        User storage user = users[addrUser];
        return user.lastLoginTime;
    }

    /** @dev Set user last login time
    * @param addrUser Address of user.
    * @param loginTime  login time.
    */
    function setUserLastLoginTime(address addrUser, uint256 loginTime) public onlyContract {
        //ignore contract address, can't process
        if (addrUser.isContract()) {
            return;
        }

        User storage user = _findAddUser(addrUser);
        user.lastLoginTime = loginTime;
    }

    /** @dev Add user watch time.
    * @param addrUser Address of user.
    * @param fileId fileId.
    * @param watchTime watchTime.
    */
    function addUserWatchTime(address addrUser, uint256 fileId, uint256 watchTime) public onlyContract {
        //ignore contract address, can't process
        if (addrUser.isContract()) {
            return;
        }

        User storage user = _findAddUser(addrUser);
        if (user.watchTime[fileId] <= 0) {
            user.watchFiles.push(fileId);
        }
        user.watchTime[fileId] = user.watchTime[fileId].add(watchTime);
    }

    /** @dev Get user watch time.
    * @param addrUser Address of user.
    * @return watchTime Total watch time current assign cycle.
    */
    function getUserWatchTime(address addrUser) external view returns (uint256 watchTotal, uint256[] memory fileIds, uint256[] memory watchTimes){
        User storage user = users[addrUser];
        uint256 len = user.watchFiles.length;
        fileIds = new uint256[](len);
        watchTimes = new uint256[](len);
        watchTotal = 0;
        for (uint256 index = 0; index < len; index++) {
            uint256 fileId = user.watchFiles[index];
            fileIds[index] = fileId;
            watchTimes[index] = user.watchTime[fileId];
            watchTotal = watchTotal.add(watchTimes[index]);
        }
    }
            // mapping (uint256=>uint256) watchTime;       // watch time current assign period
        // uint256 uploadCount;                        // upload count current assign period

    /** @dev Clear user watch time.
    * @param addrUser Address of user.
    */
    function clearUserWatchTime(address addrUser) external onlyContract {
        User storage user = _findAddUser(addrUser);
        for (uint256 index = 0; index < user.watchFiles.length; index++) {
            uint256 fileId = user.watchFiles[index];
            delete user.watchTime[fileId];
        }
        delete user.watchFiles;
    }

    /** @dev Set user upload count.
    * @param addrUser Address of user.
    */
    function clearUserUploadCount(address addrUser) external onlyContract {
        User storage user = _findAddUser(addrUser);
        user.uploadCount = 0;
    }

    /** @dev Add user upload count.
    * @param addrUser Address of user.
    * @param uploadCount Upload count.
    */
    function addUserUploadCount(address addrUser, uint256 uploadCount) public onlyContract {
        //ignore contract address, can't process
        if (addrUser.isContract()) {
            return;
        }

        User storage user = _findAddUser(addrUser);
        user.uploadCount = user.uploadCount.add(uploadCount);
    }

    /** @dev Get user upload count.
    * @param addrUser Address of user.
    * @return uploadCount Total upload count current assign cycle.
    */
    function getUserUploadCount(address addrUser) external view returns (uint256){
        User storage user = users[addrUser];
        return user.uploadCount;
    }

    /** @dev Get invited user count.
    * @param addrUser Address of user.
    * @return invitedCount Count of invited users.
    */
    function getInvitedFriendsCount(address addrUser) public view returns (uint256){
        User storage user = users[addrUser];
        return user.invitedFriendsCount;
    }

    /** @dev Add invited user.
    * @param addrUser Address of user.
    * @return last login time.
    */
    function addInvitedFriends(address addrUser) public onlyContract {
        //ignore contract address, can't process
        if (addrUser.isContract()) {
            return;
        }

        User storage user = _findAddUser(addrUser);
        user.invitedFriendsCount = user.invitedFriendsCount.add(1);
    }

    /** @dev Get user power.
    * @param addrUser Address of user.
    * @return longtermPower Longterm Power.
    * @return shorttermPower Shortterm Power.
    */
    function getUserPower(address addrUser) public view returns (uint256 longtermPower, uint256 shorttermPower){
        User storage user = users[addrUser];
        longtermPower = user.longtermPower;
        shorttermPower = user.shorttermPower;
    }

    /** @dev Set user short term power
    * @param addrUser Address of user.
    * @param value  Value of short term power.
    */
    function setUserShorttermPower(address addrUser, uint256 value) public onlyContract {
        //ignore contract address, can't process
        if (addrUser.isContract()) {
            return;
        }

        User storage user = _findAddUser(addrUser);
        user.shorttermPower = value;
    }

    /** @dev Set user long term power
    * @param addrUser Address of user.
    * @param value  Value of long term power.
    */
    function setUserLongtermPower(address addrUser, uint256 value) public onlyContract {
        //ignore contract address, can't process
        if (addrUser.isContract()) {
            return;
        }

        User storage user = _findAddUser(addrUser);
        user.longtermPower = value;
    }
 }
