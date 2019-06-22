pragma solidity ^0.4.8;

import "./BokerManager.sol";
import "./BokerLogData.sol";


contract BokerLog is BokerManaged {

    constructor(address addrManager) BokerManaged(addrManager) public {
    }

    function commonLogAdd(string data, uint256 time) external onlyContract {
        BokerLogData(contractAddress(ContractLogData)).commonLogAdd(data, time);
    }

    function userPowerLogAdd(
        address addrDapp, address addrUser, uint256 powerType, int256 powerChange, uint256 powerNew, uint256 reason, uint256 param1) external onlyContract {
        BokerLogData(contractAddress(ContractLogData)).userPowerLogAdd(addrDapp, addrUser, powerType, powerChange, powerNew, reason, param1);
    }

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
        return BokerLogData(contractAddress(ContractLogData)).userPowerLogGet(addrUser, page, pageSize);
    }

    function assignTokenLogAdd(
        uint256 longtermTokens, uint256 shorttermTokens, uint256 longtermPowerTotal, uint256 shorttermPowerTotal) public onlyContract {
        BokerLogData(contractAddress(ContractLogData)).assignTokenLogAdd(longtermTokens, shorttermTokens, longtermPowerTotal, shorttermPowerTotal);
    }

    function assignTokenLogGetIndex(uint256 index) public view returns (
        uint256 longtermTokens, uint256 shorttermTokens, uint256 longtermPowerTotal, uint256 shorttermPowerTotal, uint256 time) {
        return BokerLogData(contractAddress(ContractLogData)).assignTokenLogGetIndex(index);
    }

    function assignTokenLogGet(uint256 page, uint256 pageSize) public view returns (
        uint256[] longtermTokenses, uint256[] shorttermTokenses, uint256[] longtermPowerTotals, uint256[] shorttermPowerTotals, uint256[] times) {
        return BokerLogData(contractAddress(ContractLogData)).assignTokenLogGet(page, pageSize);
    }

    function voteLogAdd(address addrUser, address addrCandidate, uint256 voteType, uint256 tokens) public onlyContract {
        BokerLogData(contractAddress(ContractLogData)).voteLogAdd(addrUser, addrCandidate, voteType, tokens);
    }

    function voteLogGet(address addrUser, uint256 page, uint256 pageSize) public view returns (
        uint256[] voteTypes, uint256[] tokenses, address[] addrCandidates, uint256[] times) {
        return BokerLogData(contractAddress(ContractLogData)).voteLogGet(addrUser, page, pageSize);
    }

    function voteRotateLogAdd(uint256 round) public onlyContract {
        BokerLogData(contractAddress(ContractLogData)).voteRotateLogAdd(round);
    }

    function voteRotateLogGetIndex(uint256 index) public view returns (uint256 round, uint256 time) {
        return BokerLogData(contractAddress(ContractLogData)).voteRotateLogGetIndex(index);
    }

    function voteRotateLogGet(uint256 page, uint256 pageSize) public view returns (uint256[] rounds, uint256[] times) {
        return BokerLogData(contractAddress(ContractLogData)).voteRotateLogGet(page, pageSize);
    }

    function financeLogAdd(
        address addrFrom, address addrTo, uint256 tokensChange, uint256 tokensFrom, uint256 tokensTo, uint256 reason) external onlyContract {
        BokerLogData(contractAddress(ContractLogData)).financeLogAdd(addrFrom, addrTo, tokensChange, tokensFrom, tokensTo, reason);
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
        return BokerLogData(contractAddress(ContractLogData)).userFinanceLogGet(addrUser, page, pageSize);
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
    function userTipLogGet(address addrUser, uint256 page, uint256 pageSize) external view returns (
        address[] addrUsers, int256[] tokensChanges, uint256[] tokensAfters, uint256[] times) {
        return BokerLogData(contractAddress(ContractLogData)).userTipLogGet(addrUser, page, pageSize);
    }
}
