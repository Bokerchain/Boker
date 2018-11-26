pragma solidity ^0.4.8;

import "./BokerManager.sol";


contract BokerLogData is BokerManaged {
    using PageUtil for uint256;

    constructor(address addrManager) BokerManaged(addrManager) public {
    }

    // 算力变化日志
    struct UserPowerLogEntry {
        uint256 powerType;
        uint256 powerOld;
        uint256 powerNew;
        uint256 reason;
        uint256 time;
    }
    struct UserPowerLog {
        address addrUser;
        UserPowerLogEntry[] entries;
    }
    mapping (address=>UserPowerLog) userPowerLogs;  

    function userPowerLogAdd(
        address addrUser, UserPowerType powerType, uint256 powerOld, uint256 powerNew, UserPowerReason reason) public onlyContract {
        UserPowerLog storage userPowerLog = userPowerLogs[addrUser];
        userPowerLog.addrUser = addrUser;
        userPowerLog.entries.push(UserPowerLogEntry(uint256(powerType), powerOld, powerNew, uint256(reason), now));
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
    function userPowerLogGet(address addrUser, uint256 page, uint256 pageSize) public view returns (
            uint256[] powerTypes, uint256[] powerOlds, uint256[] powerNews, uint256[] reasons, uint256[] times) {
        UserPowerLog storage userPowerLog = userPowerLogs[addrUser];

        if(userPowerLog.addrUser == address(0)) {
            return;
        }

        (uint256 start, uint256 end) = userPowerLog.entries.length.pageRange(page, pageSize);
        uint256 len = end - start + 1;
        powerTypes = new uint256[](len);
        powerOlds = new uint256[](len);
        powerNews = new uint256[](len);
        reasons = new uint256[](len);
        times = new uint256[](len);
        for (uint256 index = start; index <= end; index++) {
            UserPowerLogEntry storage entry = userPowerLog.entries[index];
            powerTypes[index-start] = entry.powerType;
            powerOlds[index-start] = entry.powerOld;
            powerNews[index-start] = entry.powerNew;
            reasons[index-start] = entry.reason;
            times[index-start] = entry.time;
        }
    }

    //分币日志
    struct AssignTokenLog {
        uint256 longtermTokens;
        uint256 shorttermTokens;
        uint256 longtermPowerTotal;
        uint256 shorttermPowerTotal;
        uint256 time;
    }
    AssignTokenLog[] assignTokenLogs;

    function assignTokenLogAdd(
        uint256 longtermTokens, uint256 shorttermTokens, uint256 longtermPowerTotal, uint256 shorttermPowerTotal) public onlyContract {
        assignTokenLogs.push(AssignTokenLog(longtermTokens, shorttermTokens, longtermPowerTotal, shorttermPowerTotal, now));
    }

    function assignTokenLogGet(uint256 index) public view returns (
        uint256 longtermTokens, uint256 shorttermTokens, uint256 longtermPowerTotal, uint256 shorttermPowerTotal, uint256 time) {
        require(index < assignTokenLogs.length, "index exceeds assignTokenLogs length");
        AssignTokenLog storage assignTokenLog = assignTokenLogs[index];
        longtermTokens = assignTokenLog.longtermTokens;
        shorttermTokens = assignTokenLog.shorttermTokens;
        longtermPowerTotal = assignTokenLog.longtermPowerTotal;
        shorttermPowerTotal = assignTokenLog.shorttermPowerTotal;
        time = assignTokenLog.time;
    }

    function assignTokenLogGet(uint256 page, uint256 pageSize) public view returns (
        uint256[] longtermTokenses, uint256[] shorttermTokenses, uint256[] longtermPowerTotals, uint256[] shorttermPowerTotals, uint256[] times) {
        (uint256 start, uint256 end) = assignTokenLogs.length.pageRange(page, pageSize);
        uint256 len = end - start + 1;
        longtermTokenses = new uint256[](len);
        shorttermTokenses = new uint256[](len);
        longtermPowerTotals = new uint256[](len);
        shorttermPowerTotals = new uint256[](len);
        times = new uint256[](len);
        for (uint256 index = start; index <= end; index++) {
            AssignTokenLog storage log = assignTokenLogs[index];
            longtermTokenses[index-start] = log.longtermTokens;
            shorttermTokenses[index-start] = log.shorttermTokens;
            longtermPowerTotals[index-start] = log.longtermPowerTotal;
            shorttermPowerTotals[index-start] = log.shorttermPowerTotal;
            times[index-start] = log.time;
        }
    }

    //投票日志
    struct VoteLogEntry {
        uint256   voteType;
        uint256 tokens;
        address addrCandidate;
        uint256 time;
    }
    struct VoteLog {
        address addrUser;
        VoteLogEntry[] entries;
    }
    mapping (address=>VoteLog) voteLogs;

    function voteLogAdd(address addrUser, address addrCandidate, uint256 voteType, uint256 tokens) public onlyContract {
        VoteLog storage voteLog = voteLogs[addrUser];
        voteLog.addrUser = addrUser;
        voteLog.entries.push(VoteLogEntry(uint256(voteType), tokens, addrCandidate, now));
    }

    function voteLogGet(address addrUser, uint256 page, uint256 pageSize) public view returns (
        uint256[] voteTypes, uint256[] tokenses, address[] addrCandidates, uint256[] times) {
        VoteLog storage voteLog = voteLogs[addrUser];

        if(voteLog.addrUser == address(0)) {
            return;
        }

        (uint256 start, uint256 end) = voteLog.entries.length.pageRange(page, pageSize);
        uint256 len = end - start + 1;
        voteTypes = new uint256[](len);
        tokenses = new uint256[](len);
        addrCandidates = new address[](len);
        times = new uint256[](len);
        for (uint256 index = start; index <= end; index++) {
            VoteLogEntry storage entry = voteLog.entries[index];
            voteTypes[index-start] = entry.voteType;
            tokenses[index-start] = entry.tokens;
            addrCandidates[index-start] = entry.addrCandidate;
            times[index-start] = entry.time;
        }
    }

    //vote rotate 日志
    struct VoteRotateLog {
        uint256 round;
        uint256 time;
    }
    VoteRotateLog[] voteRotateLogs;

    function voteRotateLogAdd(uint256 round) public onlyContract {
        voteRotateLogs.push(VoteRotateLog(round, now));
    }

    function voteRotateLogGet(uint256 index) public view returns (uint256 round, uint256 time) {
        require(index < voteRotateLogs.length, "index exceeds voteRotateLogs length");
        VoteRotateLog storage voteRotateLog = voteRotateLogs[index];
        round = voteRotateLog.round;
        time = voteRotateLog.time;
    }

    function voteRotateLogGet(uint256 page, uint256 pageSize) public view returns (uint256[] rounds, uint256[] times) {
        (uint256 start, uint256 end) = voteRotateLogs.length.pageRange(page, pageSize);
        uint256 len = end - start + 1;
        rounds = new uint256[](len);
        times = new uint256[](len);
        for (uint256 index = start; index <= end; index++) {
            VoteRotateLog storage voteRotateLog = voteRotateLogs[index];
            rounds[index-start] = voteRotateLog.round;
            times[index-start] = voteRotateLog.time;
        }
    }

    //finance日志
    struct FinanceLogEntry {
        address addrFrom;
        address addrTo;
        uint256 tokensChange;        // tokens changed
        uint256 tokensFrom;          // tokens of from user after change
        uint256 tokensTo;            // tokens of to user after change
        uint256 reason;
        uint256 time;
    }
    FinanceLogEntry[] public financeLogs;

    struct UserFinanceLog {
        address addrUser;
        uint256[] indexes;
    }
    mapping (address=>UserFinanceLog) userFinanceLogs;

    function financeLogAdd(
        address addrFrom, address addrTo, uint256 tokensChange, uint256 tokensFrom, uint256 tokensTo, uint256 reason) public onlyContract {
        financeLogs.push(FinanceLogEntry(addrFrom, addrTo, tokensChange, tokensFrom, tokensTo, uint256(reason), now));

        //from user
        UserFinanceLog storage fromLog = userFinanceLogs[addrFrom];
        if (fromLog.addrUser == address(0)) {
            fromLog.addrUser = addrFrom;
        }
        fromLog.indexes.push(financeLogs.length - 1);

        //to user
        UserFinanceLog storage toLog = userFinanceLogs[addrTo];
        if (toLog.addrUser == address(0)) {
            toLog.addrUser = addrTo;
        }
        toLog.indexes.push(financeLogs.length - 1);
    }

    function _userFinanceLogArrayInit(uint256 len) private pure 
        returns (int256[] tokensChanges, uint256[] tokensAfters, uint256[] reasons, uint256[] times, uint256[] indexes) {
        tokensChanges = new int256[](len);
        tokensAfters = new uint256[](len);
        reasons = new uint256[](len);
        times = new uint256[](len);
        indexes = new uint256[](len);
    }

    function _userFinanceLogGet(UserFinanceLog storage userFinanceLog, uint256 start, uint256 end) private view returns (
        int256[] tokensChanges, uint256[] tokensAfters, uint256[] reasons, uint256[] times, uint256[] indexes) {
        (tokensChanges, tokensAfters, reasons, times, indexes) = _userFinanceLogArrayInit(end - start + 1);
        for (uint256 index = start; index <= end; index++) {
            indexes[index-start] = userFinanceLog.indexes[index];
            FinanceLogEntry storage entry = financeLogs[userFinanceLog.indexes[index]];
            if (userFinanceLog.addrUser == entry.addrFrom) {
                tokensChanges[index-start] = -int256(entry.tokensChange);
                tokensAfters[index-start] = entry.tokensFrom;
            } else {
                tokensChanges[index-start] = int256(entry.tokensChange);
                tokensAfters[index-start] = entry.tokensTo;
            }
            reasons[index-start] = entry.reason;
            times[index-start] = entry.time;
        }
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
        UserFinanceLog storage userFinanceLog = userFinanceLogs[addrUser];
        if(userFinanceLog.addrUser == address(0)) {
            return;
        }
        (uint256 start, uint256 end) = userFinanceLog.indexes.length.pageRange(page, pageSize);
        (tokensChanges, tokensAfters, reasons, times, indexes) = _userFinanceLogGet(userFinanceLog, start, end);
    }
}