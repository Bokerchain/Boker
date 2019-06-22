pragma solidity ^0.4.8;

import "./BokerDefine.sol";
import "./BokerCommon.sol";
import "./BokerManager.sol";

contract BokerUserData  is BokerManaged {
    using SafeMath for uint256;
    using AddressUtil for address;

    struct WatchEntry {
        address addrDapp;
        mapping (uint256=>uint256) watchTimes;      // fileid to watch time
        uint256[] watchFiles;                       // all watched files
    }

    struct User {
        address addr;                               // address of user
        uint256 index;                              // index in array of user
        uint256 registerTime;                       // time of register
        uint256 certificationTime;                  // time of real name certification

        uint256 lastLoginTime;                      
        uint256 invitedFriendsCount;

        mapping (address=>uint256) bindedDappMap;       // binded dapps address => bindTime
        address[] bindedDappArray;                      //binded dapps addresses

        mapping (address=>WatchEntry) watchEntries;     // dapp address to watch entry current assign period
        address[] watchDapps;                           // watched dapps current assign period
        uint256 watchEntryTotal;

        uint256 uploadCount;                            // upload count current assign period

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

    function _findAddUser(address addr) private returns (User storage user) {
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
    function getUserRegisterTime(address addrUser) public view returns (uint256) {
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
    function getUserCertificationTime(address addrUser) public view returns (uint256) {
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


    /** @dev get dapp bind time.
    * @return bindTime dapp bind time.
    */
    function getDappBind(address addrUser, address addrDapp) public view returns (uint256 bindTime) {
        User storage user = users[addrUser];
        return user.bindedDappMap[addrDapp];
    }

    /** @dev Add binded dapp.
    * @param addrUser Address of user.
    * @param addrDapp Bind wallet time.
    */
    function addDappBind(address addrUser, address addrDapp) public onlyContract {
        //ignore contract address, can't process
        if (addrUser.isContract()) {
            return;
        }

        User storage user = _findAddUser(addrUser);
        user.bindedDappMap[addrDapp] = now;
    }

    /** @dev Get user last login time.
    * @param addrUser Address of user.
    * @return last login time.
    */
    function getLastLoginTime(address addrUser) public view returns (uint256) {
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
    function addUserWatchTime(address addrDapp, address addrUser, uint256 fileId, uint256 watchTime) public onlyContract {
        //ignore contract address, can't process
        if (addrUser.isContract()) {
            return;
        }

        User storage user = _findAddUser(addrUser);
        WatchEntry storage watchEntry = user.watchEntries[addrDapp];
        if (watchEntry.addrDapp == address(0)) {
            watchEntry.addrDapp = addrDapp;
            user.watchDapps.push(addrDapp);
        }

        if (watchEntry.watchTimes[fileId] <= 0) {
            watchEntry.watchFiles.push(fileId);
            user.watchEntryTotal = user.watchEntryTotal.add(1);
        }
        watchEntry.watchTimes[fileId] = watchEntry.watchTimes[fileId].add(watchTime);
    }

    /** @dev Get user watch time.
    * @param addrUser Address of user.
    * @return watchTime Total watch time current assign cycle.
    */
    function getUserWatchTime(address addrUser) external view returns (uint256 watchTotal, address[] addrDapps, uint256[] memory fileIds, uint256[] memory watchTimes) {
        User storage user = users[addrUser];
        uint256 len = user.watchEntryTotal;
        addrDapps = new address[](len);
        fileIds = new uint256[](len);
        watchTimes = new uint256[](len);
        watchTotal = 0;

        uint256 index = 0;
        for (uint256 i = 0; i < user.watchDapps.length; i++) {
            address addrDapp = user.watchDapps[i];
            WatchEntry storage watchEntry = user.watchEntries[addrDapp];
            for (uint256 j = 0; j < watchEntry.watchFiles.length; j++) {
                uint256 fileId = watchEntry.watchFiles[index];
                addrDapps[index] = addrDapp;
                fileIds[index] = fileId;
                watchTimes[index] = watchEntry.watchTimes[fileId];
                watchTotal = watchTotal.add(watchTimes[index]);
                index++;
            }
        }

    }
            // mapping (uint256=>uint256) watchTime;       // watch time current assign period
        // uint256 uploadCount;                        // upload count current assign period

    /** @dev Clear user watch time.
    * @param addrUser Address of user.
    */
    function clearUserWatchTime(address addrUser) external onlyContract {
        User storage user = _findAddUser(addrUser);
        for (uint256 i = 0; i < user.watchDapps.length; i++) {
            address addrDapp = user.watchDapps[i];
            WatchEntry storage watchEntry = user.watchEntries[addrDapp];
            for (uint256 j = 0; j < watchEntry.watchFiles.length; j++) {
                uint256 fileId = watchEntry.watchFiles[j];
                delete watchEntry.watchTimes[fileId];
            }
            delete watchEntry.watchFiles;
            delete user.watchEntries[addrDapp];
        }
        delete user.watchDapps;
        user.watchEntryTotal = 0;
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
        public onlyContract {
        
         //ignore contract address, can't process
        if (addrUser.isContract()) {
            return;
        }

        User storage user = _findAddUser(addrUser);
        user.registerTime = registerTime;
        user.invitedFriendsCount = invitedFriendsCount;
        user.longtermPower = longtermPower;
        user.shorttermPower = shorttermPower;
        for (uint256 i = 0; i < bindedDappAddresses.length; i++) {
            user.bindedDappMap[bindedDappAddresses[i]] = bindedDappTimes[i];
        }
    }
 }
