pragma solidity ^0.4.8;

import "./BokerManager.sol";
import "./BokerNodeData.sol";
import "./BokerFinance.sol";
import "./BokerLog.sol";

contract BokerNode is BokerManaged {
    using SafeMath for uint256;
    using Uint256Util for uint256;

    constructor(address addrManager) BokerManaged(addrManager) public {
    }

    /** @dev Register to be candidate of verifier.
    * @param addrCandidate Address of candidate.
    * @param description description of node
    * @param team description of team
    * @param name name of node
    */
    function registerCandidate(address addrCandidate, string description, string team, string name) external onlyContract {
        require(addrCandidate != address(0), "addrCandidate is 0");
        BokerNodeData(contractAddress(ContractNodeData)).addCandidate(addrCandidate, description, team, name);
    }

    /** @dev Vote for candidate.
    * @param addrVoter Address of voter.
    * @param addrCandidate Address of candidate.
    * @param tokens tokens.
    */
    function vote(address addrVoter, address addrCandidate, uint256 tokens) external onlyContract {
        BokerNodeData nodeData = BokerNodeData(contractAddress(ContractNodeData));
        require(nodeData.existCandidate(addrCandidate), "addrCandidate not found!");
        require(tokens > 0, "tokens <= 0!");

        nodeData.increaseVoterDeposit(addrVoter, tokens);
        nodeData.increaseVoterVote(addrVoter, addrCandidate, tokens);
        nodeData.updateVoterUnlockTime(addrVoter);
        nodeData.increaseCandidateTicket(addrCandidate, tokens);

        //日志
        BokerLog(contractAddress(ContractLog)).voteLogAdd(addrVoter, addrCandidate, uint256(VoteType.Vote), tokens);
    }
    
    /** @dev Cancel all votes of voter.
    * @param addrVoter Address of voter.
    * @param addresses Address of candidate.
    * @param tickets tickets of candidate.
    */
    function _cancelAllVotes(address addrVoter, address[] memory addresses, uint256[] memory tickets) private {
        BokerNodeData nodeData = BokerNodeData(contractAddress(ContractNodeData));
        for(uint256 index = 0; index < addresses.length; index++){
            address addrCandidate = addresses[index];
            uint256 ticket = tickets[index];
            nodeData.clearVoterVote(addrVoter, addrCandidate);
            nodeData.decreaseCandidateTicket(addrCandidate, ticket);
        }
    }

    /** @dev Cancel all votes of voter.
    * @param addrVoter Address of voter.
    */
    function cancelAllVotes(address addrVoter) external onlyContract {
        BokerNodeData nodeData = BokerNodeData(contractAddress(ContractNodeData));
        (address[] memory addresses, uint256[] memory tickets, uint256 unlockTime, uint256 deposit) = nodeData.getVoteInfo(addrVoter);
        unlockTime;
        _cancelAllVotes(addrVoter, addresses, tickets);

        //return deposit
        nodeData.setVoterDeposit(addrVoter, 0);
        BokerFinance(contractAddress(ContractFinance)).grantTokenTo(addrVoter, deposit, uint256(FinanceReason.VoteCancel));

        //日志
        BokerLog(contractAddress(ContractLog)).voteLogAdd(addrVoter, address(0), uint256(VoteType.Cancel), deposit);
    }

    /** @dev Get round of vote.
    * @return round Round of vote.
    */
    function getVoteRound() external view returns(uint256 round) {
        return BokerNodeData(contractAddress(ContractNodeData)).voteCycleRound();
    }

    /** @dev Get all candidates with tickets.
    * @return addresses The addresses of candidates.
    * @return tickets The tickets of candidates.
    */
    function getCandidates() external view returns(address[] memory addresses, uint256[] memory tickets) {
        return BokerNodeData(contractAddress(ContractNodeData)).getCandidates();
    }

    /** @dev Get candidate.
    * @param addrCandidate Address of candidate.
    * @return description description of node
    * @return team description of team
    * @return name name of node
    * @return tickets tickets of node
    */
    function getCandidate(address addrCandidate) external view  returns(string description, string team, string name, uint256 tickets) {
        return BokerNodeData(contractAddress(ContractNodeData)).getCandidate(addrCandidate);
    }

    function _voteCycleIsEnd() private view returns (bool) {
        BokerNodeData nodeData = BokerNodeData(contractAddress(ContractNodeData));
        if(nodeData.voteCycleBegin() + globalConfigInt(CfgVoteCyclePeriod) <= now) { 
            return true;
        }
        return false;
    }

    /** @dev Check if vote should be rotate.
    * @return needRotate If current vote round need rotate.
    */
    function checkVote() external view returns (bool){
        if(!_voteCycleIsEnd()) {
            return false;
        }
        return true;
    }

    function _checkUnlockDeposit() private {
        BokerNodeData nodeData = BokerNodeData(contractAddress(ContractNodeData));
        address[] memory addressesVoter = nodeData.getVoters();

        uint256 len = addressesVoter.length;
        for(uint256 index = 0; index < len; index++){
            address addrVoter = addressesVoter[index];
            (address[] memory addresses, uint256[] memory tickets, uint256 unlockTime, uint256 deposit) = nodeData.getVoteInfo(addrVoter);
           
            //check if voter deposit can unlock;
            if(unlockTime.diff(now) <= globalConfigInt(CfgVoteUnlockPrecision)){
                _cancelAllVotes(addrVoter, addresses, tickets);
                //return deposit
                nodeData.setVoterDeposit(addrVoter, 0);
                BokerFinance(contractAddress(ContractFinance)).grantTokenTo(addrVoter, deposit, uint256(FinanceReason.VoteUnlock));

                //日志
                BokerLog(contractAddress(ContractLog)).voteLogAdd(addrVoter, address(0), uint256(VoteType.Unlock), deposit);    
            }
        }
    }

    function _rotateTickCycle() private returns (uint256 round) {
        BokerNodeData nodeData = BokerNodeData(contractAddress(ContractNodeData));
        uint256 voteCycleBeginOld = nodeData.voteCycleBegin();
        uint256 voteCyclePeriod = globalConfigInt(CfgVoteCyclePeriod);
        uint256 roundAdd = now.sub(voteCycleBeginOld).div(voteCyclePeriod);
        nodeData.setVoteCycleBegin(voteCycleBeginOld.add(roundAdd.mul(voteCyclePeriod)));
        return nodeData.increaseVoteCycleRound(roundAdd);
    }

    /** @dev Rotate vote cycle.
    */
    function rotateVote() external onlyContract {
        if(!_voteCycleIsEnd()) {
            return;
        }

        _checkUnlockDeposit();

        BokerNodeData nodeData = BokerNodeData(contractAddress(ContractNodeData));
        uint256 roundOld = nodeData.voteCycleRound();
        _rotateTickCycle();

        //日志
        BokerLog(contractAddress(ContractLog)).voteRotateLogAdd(roundOld); 
    }
}