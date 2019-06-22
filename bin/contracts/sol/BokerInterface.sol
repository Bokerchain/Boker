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
    function getInvitedFriendsCount(address addrUser) external view returns (uint256);
    function userCount() external view returns (uint256 total);
    // function transferTokenTo(address addrTo, uint256 reason) external payable;
    function giveTipTo(address addrTo) external payable;
    function getUserBindDapp(address addrUser) external view returns (address[] addrDapps, bool[] bindeds, uint256[] powerAdds);
    // function setUser(
    //     address addrUser, uint256 registerTime, uint256 invitedFriendsCount, address[] bindedDappAddresses, uint256[] bindedDappTimes, uint256 longtermPower, uint256 shorttermPower) 
    //     external;
}

interface IFile {
    function addFile(address uploader, address owner, uint256 fileId, string ipfsHash, string ipfsUrl, string aliDnaFileId) external;
    function addUserFile(address uploader, uint256 fileId) external;
    function fileInfoGet(uint256 fileId) external view returns (
        address uploader, address owner, string aliDnaFileId, string ipfsUrl, uint256 playCount, uint256 playTime, uint256 userCount, uint256 createTime); 
    function userFilesGet(address user, uint256 page, uint256 pageSize) external view returns (
        uint256[] fileIds, uint256[] playCounts, uint256[] playTimes, uint256[] userCounts, uint256[] createTimes);
    function fileStatisticsDailyGet() external view returns (uint256[] fileIds, uint256[] playCounts, uint256[] playTimes);
}

interface IFinance {
    function receiveTokenFrom(address addrFrom, uint256 reason) external payable;
}

interface IDapp {
    function dappSet(address dappAddr, uint256 dappType, uint256 userCount, uint256 monthlySales) external;
}

interface Ilog {
    // function userPowerLogAdd(
    //     address addrDapp, address addrUser, uint256 powerType, int256 powerChange, uint256 powerNew, uint256 reason, uint256 param1) external;

    function userPowerLogGet(address addrUser, uint256 page, uint256 pageSize) external view returns (
        address[] addrDapps, uint256[] powerTypes, int256[] powerChanges, uint256[] reasons, uint256[] param1s, uint256[] times);

    // function financeLogAdd(
    //     address addrFrom, address addrTo, uint256 tokensChange, uint256 tokensFrom, uint256 tokensTo, uint256 reason) external;

    function userFinanceLogGet(address addrUser, uint256 page, uint256 pageSize) external view returns (
        int256[] tokensChanges, uint256[] tokensAfters, uint256[] reasons, uint256[] times, uint256[] indexes);

    function userTipLogGet(address addrUser, uint256 page, uint256 pageSize) external view returns (
        address[] addrUsers, int256[] tokensChanges, uint256[] tokensAfters, uint256[] times);
}

interface ITokenPower {
    function tokenInfoGet() external view returns (uint256 tokenAssigned, uint256 tokenToAssign);
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
    function fireUserEvents(
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
    
    // function fireUserEvent(
    //     address addrFrom, uint256 eventType, address addrTo, uint256 timestamp, uint256 eventValue1, uint256 eventValue2) 
    //     external whenNotPaused onlyDapp {
    //     require(addrFrom != address(0), "addrFrom is 0!");

    //     IUser(contractAddress(ContractUser)).onUserEvent(msg.sender, addrFrom, eventType, addrTo, timestamp, eventValue1, eventValue2);
    // }

    /** @dev Get user.
    * @param addrUser Address of user.
    * @return balance balance of user
    * @return longtermPower long term power
    * @return shorttermPower short term power
    */
    function getUser(address addrUser) external view  returns(uint256 balance, uint256 longtermPower, uint256 shorttermPower) {
        return IUser(contractAddress(ContractUser)).getUser(addrUser);
    }

    /** @dev Get invited user count.
    * @param addrUser Address of user.
    * @return invitedCount Count of invited users.
    */
    function getInvitedFriendsCount(address addrUser) external view returns (uint256) {
        return IUser(contractAddress(ContractUser)).getInvitedFriendsCount(addrUser);
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

    /** @dev Get daily file statistics info
    * @return fileIds file id.
    * @return playCounts daily play total count.
    * @return playTimes daily play total time.
    */
    function fileStatisticsDailyGet() external view returns (uint256[] fileIds, uint256[] playCounts, uint256[] playTimes) {
        return IFile(contractAddress(ContractFile)).fileStatisticsDailyGet();
    }

    /** @dev Set dapp info.
    * @param dappAddr Address of from user.    
    * @param dappType Type of event.
    * @param userCount Address of to user.
    * @param monthlySales timestamps.
    */
    function dappSet(address dappAddr, uint256 dappType, uint256 userCount, uint256 monthlySales) external onlyAdmin {
        IDapp(contractAddress(ContractDapp)).dappSet(dappAddr, dappType, userCount, monthlySales);
    }

    // /** @dev Set user.
    // * @param addrUser address of user.
    // * @param registerTime register time of user.
    // * @param invitedFriendsCount total invited friends count.
    // * @param bindedDappAddresses addresses of binded dapps.
    // * @param bindedDappTimes binded times of binded dapps.
    // * @param longtermPower longterm power of user.
    // * @param shorttermPower shortterm power of user.
    // */
    // function setUser(
    //     address addrUser, uint256 registerTime, uint256 invitedFriendsCount, address[] bindedDappAddresses, uint256[] bindedDappTimes, uint256 longtermPower, uint256 shorttermPower) 
    //     external onlyAdmin {
    //     IUser(contractAddress(ContractUser)).setUser(addrUser, registerTime, invitedFriendsCount, bindedDappAddresses, bindedDappTimes, longtermPower, shorttermPower);
    // }

    // function userPowerLogAdd(
    //     address addrDapp, address addrUser, uint256 powerType, int256 powerChange, uint256 powerNew, uint256 reason, uint256 param1) external {
    //     Ilog(contractAddress(ContractLog)).userPowerLogAdd(addrDapp, addrUser, powerType, powerChange, powerNew, reason, param1);
    // }

    // function financeLogAdd(
    //     address addrFrom, address addrTo, uint256 tokensChange, uint256 tokensFrom, uint256 tokensTo, uint256 reason) external {
    //     Ilog(contractAddress(ContractLog)).financeLogAdd(addrFrom, addrTo, tokensChange, tokensFrom, tokensTo, reason);
    // }

    /** @dev Get user power logs
    * @param addrUser Address of user.    
    * @param page page number of result.
    * @param pageSize page size of result.
    * @return addrDapps power type
    * @return powerTypes power type
    * @return powerChanges power change
    * @return reasons change reason
    * @return param1s change param1
    * @return times change time
    */
    function userPowerLogGet(address addrUser, uint256 page, uint256 pageSize) external view returns (
            address[] addrDapps, uint256[] powerTypes, int256[] powerChanges, uint256[] reasons, uint256[] param1s, uint256[] times) {
        return Ilog(contractAddress(ContractLog)).userPowerLogGet(addrUser, page, pageSize);
    }

    /** @dev Get user finance logs
    * @param addrUser Address of user.    
    * @param page page number of result.
    * @param pageSize page size of result.
    * @return tokensChanges tokens changed
    * @return tokensAfters tokens of from user after change
    * @return reasons change reason
    * @return times change time
    * @return indexes index in origin logs
    */
    function userFinanceLogGet(address addrUser, uint256 page, uint256 pageSize) external view returns (
        int256[] tokensChanges, uint256[] tokensAfters, uint256[] reasons, uint256[] times, uint256[] indexes) {
        return Ilog(contractAddress(ContractLog)).userFinanceLogGet(addrUser, page, pageSize);
    }

    /** @dev Get user tip logs
    * @param addrUser Address of user.    
    * @param page page number of result.
    * @param pageSize page size of result.
    * @return addrUsers address of other user.
    * @return tokensChanges tokens changed
    * @return tokensAfters tokens of from user after change
    * @return times change time
    */
    function userTipLogGet(address addrUser, uint256 page, uint256 pageSize) public view returns (
        address[] addrUsers, int256[] tokensChanges, uint256[] tokensAfters, uint256[] times) {
        return Ilog(contractAddress(ContractLog)).userTipLogGet(addrUser, page, pageSize);
    }

    /** @dev Get general info
    * @return userTotal total user number registered.    
    * @return tokenAssigned tokens assigned last day.
    * @return tokenToAssign tokens left to assign.
    */
    function generalInfoGet() public view returns (uint256 userTotal, uint256 tokenAssigned, uint256 tokenToAssign) {
        userTotal = IUser(contractAddress(ContractUser)).userCount();
        (tokenAssigned, tokenToAssign) = ITokenPower(contractAddress(ContractTokenPower)).tokenInfoGet(); 
    }

    // /** @dev transfer bobby to user, need gas.
    // * @param addrTo Address of to user.
    // * @param reason reason of transfer.
    // */
    // function transferTokenTo(address addrTo, uint256 reason) external payable {
    //     IUser(contractAddress(ContractUser)).transferTokenTo(addrTo, reason);
    // }

    /** @dev give tip to user
    * @param addrTo Address of to user.
    */
    function giveTipTo(address addrTo) external payable {
        IUser(contractAddress(ContractUser)).giveTipTo(addrTo);
    }

    /** @dev get user binded dapp info
    * @param addrUser Address of user.
    * @return addrDapps addresses of dapps.
    * @return bindeds if has bind.
    * @return powerAdds power to add if bind.
    */
    function getUserBindDapp(address addrUser) external view returns (address[] addrDapps, bool[] bindeds, uint256[] powerAdds) {
        return IUser(contractAddress(ContractUser)).getUserBindDapp(addrUser);
    }
}