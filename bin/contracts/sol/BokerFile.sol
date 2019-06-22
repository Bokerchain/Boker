pragma solidity ^0.4.8;

import "./BokerCommon.sol";
import "./BokerDefine.sol";
import "./BokerFileData.sol";
import "./BokerManager.sol";


contract BokerFile is BokerManaged {

    constructor(address addrManager) BokerManaged(addrManager) public {
    }

    /** @dev Handle user event.
    * @param addrDapp Address of dapp.
    * @param addrFrom Address of from user.
    * @param eventType Type of event.
    * @param addrTo Address of to user.
    * @param timestamp Data3 of event.  
    * @param eventValue1 Data1 of event.
    * @param eventValue2 Data2 of event.
    */
    function onUserEvent(
        address addrDapp, address addrFrom, UserEventType eventType, address addrTo, uint256 timestamp, uint256 eventValue1, uint256 eventValue2) 
        external onlyContract returns (Error err) {
        addrDapp;
        addrTo;
        timestamp;
        BokerFileData fileData = BokerFileData(contractAddress(ContractFileData));
        if (eventType >= UserEventType.End) {
            err = Error.EventNotSupported;
        }else if (UserEventType.Watch == eventType) {
            fileData.updateWatch(addrFrom, eventValue1, eventValue2);
            fileData.updateStatistics(eventValue1, eventValue2);
        }
        if (err != Error.Ok) {
            return err;
        }
        
        return Error.Ok;
    }

    /** @dev Add file.
    * @param uploader Address of uploader.
    * @param owner Address of from owner.
    * @param fileId file id.
    * @param ipfsHash hash of ipfs file.
    * @param ipfsUrl url of ipfs file.  
    * @param aliDnaFileId fileid of ali dna.
    */
    function addFile(address uploader, address owner, uint256 fileId, string ipfsHash, string ipfsUrl, string aliDnaFileId) external onlyContract {
        BokerFileData(contractAddress(ContractFileData)).addFile(uploader, owner, fileId, ipfsHash, ipfsUrl, aliDnaFileId);
    }

    /** @dev Add user uploaded file.
    * @param uploader Address of uploader.
    * @param fileId file id.
    */
    function addUserFile(address uploader, uint256 fileId) external onlyContract {
        BokerFileData(contractAddress(ContractFileData)).addUserFile(uploader, fileId);
    }

    /** @dev Get file info.
    * @param fileId file id.
    * @return uploader address of uploader.
    * @return owner address of owner.
    * @return aliDnaFileId ali dna file id.
    * @return ipfsUrl ipfs url.
    * @return playCount play total count.
    * @return playTime play total time.
    * @return userCount play total user count.
    * @return createTime create time.
    */
    function fileInfoGet(uint256 fileId) external view returns (
        address uploader, address owner, string aliDnaFileId, string ipfsUrl, uint256 playCount, uint256 playTime, uint256 userCount, uint256 createTime) {
        (uploader,owner,,,ipfsUrl,aliDnaFileId,createTime,playCount,playTime,userCount) = BokerFileData(contractAddress(ContractFileData)).mapId2File(fileId);
    }

    /** @dev Get all files user uploaded.
    * @param user address of user.
    * @param page page number of result.
    * @param pageSize page size of result.
    * @return fileIds file id.
    * @return playCounts play total count.
    * @return playTimes play total time.
    * @return userCounts play total user count.
    * @return createTimes create time.
    */
    function userFilesGet(address user, uint256 page, uint256 pageSize) external view returns (
        uint256[] fileIds, uint256[] playCounts, uint256[] playTimes, uint256[] userCounts, uint256[] createTimes) {
        return BokerFileData(contractAddress(ContractFileData)).userFilesGet(user, page, pageSize);
    }

    /** @dev Get daily file statistics info
    * @return fileIds file id.
    * @return playCounts daily play total count.
    * @return playTimes daily play total time.
    */
    function fileStatisticsDailyGet() external view returns (uint256[] fileIds, uint256[] playCounts, uint256[] playTimes) {
        fileIds = BokerFileData(contractAddress(ContractFileData)).statisticsDailyFiles();
        uint256 len = fileIds.length;
        playCounts = new uint256[](len);
        playTimes = new uint256[](len);
        for (uint256 index = 0; index < len; index++) {
            uint256 fileId = fileIds[index];
            (uint256 playCountDaily, uint256 playTimeDaily) = BokerFileData(contractAddress(ContractFileData)).fileStatisticsDailyGet(fileId);
            playCounts[index] = playCountDaily;
            playTimes[index] = playTimeDaily;
        }
    }
}