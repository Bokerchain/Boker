pragma solidity ^0.4.8;

import "./BokerCommon.sol";
import "./BokerManager.sol";

interface INode {
    function registerCandidate(address addrCandidate, string description, string team, string name) external;
    function vote(address addrVoter, address addrCandidate, uint256 tokens) external;
    function cancelAllVotes(address addrVoter) external;

    function getCandidates() external view  returns(address[] memory addresses, uint256[] memory tickets);
    function getCandidate(address addrCandidate) external view  returns(string description, string team, string name, uint256 tickets);
}

interface IUser {    
    function onUserEvent(
        address addrDapp, address addrFrom, uint256 eventType, address addrTo, uint256 timestamp, uint256 eventValue1, uint256 eventValue2) external;

    function getUser(address addrUser) external view  returns(uint256 balance, uint256 longtermPower, uint256 shorttermPower);
}

interface IFile {
    function addFile(address uploader, address owner, uint256 fileId, string ipfsHash, string ipfsUrl, string aliDnaFileId) external;
    function addUserFile(address uploader, uint256 fileId) external;
    function fileInfoGet(uint256 fileId) external view returns (
        address uploader, address owner, string aliDnaFileId, string ipfsUrl, uint256 playCount, uint256 playTime, uint256 userCount, uint256 createTime); 
    function userFilesGet(address user, uint256 page, uint256 pageSize) external view returns (
        uint256[] fileIds, uint256[] playCounts, uint256[] playTimes, uint256[] userCounts, uint256[] createTimes);
}

interface IFinance {
    function receiveTokenFrom(address addrFrom, uint256 reason) external payable;
}

interface IDapp {
    function dappSet(address dappAddr, uint256 dappType, uint256 userCount, uint256 monthlySales) external;
}

interface Ilog {
    function userPowerLogGet(address addrUser, uint256 page, uint256 pageSize) external view returns (
        uint256[] powerTypes, uint256[] powerOlds, uint256[] powerNews, uint256[] reasons, uint256[] times);
    
    function userFinanceLogGet(address addrUser, uint256 page, uint256 pageSize) external view returns (
        int256[] tokensChanges, uint256[] tokensAfters, uint256[] reasons, uint256[] times, uint256[] indexes);
}

contract BokerInterface is BokerManaged{
    // Ensure that we are pointing to the right contract in our set call.
    bool public isBokerInterface = true;

    constructor(address addrManager) BokerManaged(addrManager) public {
    }

    function () public payable {
        IFinance(contractAddress(ContractFinance)).receiveTokenFrom.value(msg.value)(msg.sender, uint256(FinanceReason.Transfer));
    }

    /** @dev Fire user event.
    * @param addrFrom Address of from user.    
    * @param eventTypes Type of event.
    * @param addrTos Address of to user.
    * @param timestamps timestamps.
    * @param eventValue1s eventValue1s.
    * @param eventValue2s eventValue2s.
    */
    function fireUserEvent(
        address addrFrom, uint256[] eventTypes, address[] addrTos, uint256[] timestamps, uint256[] eventValue1s, uint256[] eventValue2s) 
        external whenNotPaused onlyDapp {
        require(addrFrom != address(0), "addrFrom is 0!");

        uint256 len = eventTypes.length;
        require(len == addrTos.length, "addrTos length invalid!");
        require(len == timestamps.length, "timestamps length invalid!");
        require(len == eventValue1s.length, "eventValue1s length invalid!");
        require(len == eventValue2s.length, "eventValue2s length invalid!");

        for (uint256 index = 0; index < len; index++) {
            IUser(contractAddress(ContractUser)).onUserEvent(
                msg.sender, addrFrom, eventTypes[index], addrTos[index], timestamps[index], eventValue1s[index], eventValue2s[index]);
        }
    }
    
    function fireUserEvent(
        address addrFrom, uint256 eventType, address addrTo, uint256 timestamp, uint256 eventValue1, uint256 eventValue2) 
        external whenNotPaused onlyDapp {
        require(addrFrom != address(0), "addrFrom is 0!");

        IUser(contractAddress(ContractUser)).onUserEvent(msg.sender, addrFrom, eventType, addrTo, timestamp, eventValue1, eventValue2);
    }

    /** @dev Get user.
    * @param addrUser Address of user.
    * @return balance balance of user
    * @return longtermPower long term power
    * @return shorttermPower short term power
    */
    function getUser(address addrUser) external view  returns(uint256 balance, uint256 longtermPower, uint256 shorttermPower) {
        return IUser(contractAddress(ContractUser)).getUser(addrUser);
    }

    /** @dev Register to be candidate of verifier.
    * @param description description of node
    * @param team description of team
    * @param name name of node
    */
    function registerCandidate(string description, string team, string name) external whenNotPaused {
        INode(contractAddress(ContractNode)).registerCandidate(msg.sender, description, team, name);
    }

    /** @dev Vote for candidate.
    * @param addrCandidate Address of candidate.
    */
    function voteCandidate(address addrCandidate) external payable whenNotPaused {
        IFinance(contractAddress(ContractFinance)).receiveTokenFrom.value(msg.value)(msg.sender, uint256(FinanceReason.Vote));
        INode(contractAddress(ContractNode)).vote(msg.sender, addrCandidate, msg.value);
    }

    /** @dev Cancel all votes of voter.
    */
    function cancelAllVotes() external whenNotPaused {
        INode(contractAddress(ContractNode)).cancelAllVotes(msg.sender);
    }

    /** @dev Get all candidates with tickets.
    * @return addresses The addresses of candidates.
    * @return tickets The tickets of candidates.
    */
    function getCandidates() external view returns(address[] memory addresses, uint256[] memory tickets) {
        return INode(contractAddress(ContractNode)).getCandidates();
    }

    /** @dev Get candidate.
    * @param addrCandidate Address of candidate.
    * @return description description of node
    * @return team description of team
    * @return name name of node
    * @return tickets tickets of node
    */
    function getCandidate(address addrCandidate) external view  returns(string description, string team, string name, uint256 tickets) {
        return INode(contractAddress(ContractNode)).getCandidate(addrCandidate);
    }

    /** @dev Add file.
    * @param uploader Address of uploader.
    * @param owner Address of from owner.
    * @param fileId file id.
    * @param ipfsHash hash of ipfs file.
    * @param ipfsUrl url of ipfs file.  
    * @param aliDnaFileId fileid of ali dna.
    */
    function addFile(address uploader, address owner, uint256 fileId, string ipfsHash, string ipfsUrl, string aliDnaFileId) public onlyAdmin {
        IFile(contractAddress(ContractFile)).addFile(uploader, owner, fileId, ipfsHash, ipfsUrl, aliDnaFileId);
    }

    /** @dev Add user uploaded file.
    * @param uploader Address of uploader.
    * @param fileId file id.
    */
    function addUserFile(address uploader, uint256 fileId) public onlyAdmin {
        IFile(contractAddress(ContractFile)).addUserFile(uploader, fileId);
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
        return IFile(contractAddress(ContractFile)).fileInfoGet(fileId);
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
        return IFile(contractAddress(ContractFile)).userFilesGet(user, page, pageSize);
    }

    /** @dev Set dapp info.
    * @param dappAddr Address of from user.    
    * @param dappType Type of event.
    * @param userCount Address of to user.
    * @param monthlySales timestamps.
    */
    function dappSet(address dappAddr, uint256 dappType, uint256 userCount, uint256 monthlySales) external onlyAdmin {
        return IDapp(contractAddress(ContractDapp)).dappSet(dappAddr, dappType, userCount, monthlySales);
    }

    /** @dev Get user power logs
    * @param addrUser Address of user.    
    * @param page page number of result.
    * @param pageSize page size of result.
    * @return powerTypes power type
    * @return powerOlds power old
    * @return powerNews power new
    * @return reasons change reason
    * @return times change time
    */
    function userPowerLogGet(address addrUser, uint256 page, uint256 pageSize) external view returns (
        uint256[] powerTypes, uint256[] powerOlds, uint256[] powerNews, uint256[] reasons, uint256[] times) {
        return Ilog(contractAddress(ContractLog)).userPowerLogGet(addrUser, page, pageSize);
    }

    /** @dev Get user finance logs
    * @param addrUser Address of user.    
    * @param page page number of result.
    * @param pageSize page size of result.
    * @return addrFroms address of from user
    * @return addrTos address of to user
    * @return tokensChanges tokens changed
    * @return tokensFroms tokens of from user after change
    * @return tokensTos tokens of to user after change
    * @return reasons change reason
    * @return times change time
    * @return indexes index in origin logs
    */
    function userFinanceLogGet(address addrUser, uint256 page, uint256 pageSize) public view returns (
        int256[] tokensChanges, uint256[] tokensAfters, uint256[] reasons, uint256[] times, uint256[] indexes) {
        return Ilog(contractAddress(ContractLog)).userFinanceLogGet(addrUser, page, pageSize);
    }
}