pragma solidity ^0.4.8;

import "./BokerVerifyVoteData.sol";
import "../BokerCommon.sol";

contract BokerVerifyVoteImpl is AccessControl, Log{
    using SafeMath for uint256;

    // Ensure that we are pointing to the right contract in our set call.
    bool public isBokerVerifyVoteImpl = true;    

    BokerVerifyVoteData public data;

    /** @dev
    * @param dataAddress  Address of data contract.
    */
    constructor(address dataAddress) public {
        _setData(dataAddress);
    }

    /** @dev Set address of data contract
    * @param dataAddress  Address of data contract.
    */
    function _setData(address dataAddress) private {
        BokerVerifyVoteData dataContract = BokerVerifyVoteData(dataAddress);

        // verify that a contract is what we expect
        require(dataContract.isBokerVerifyVoteData());

        data = dataContract;
    }

    function setData(address dataAddress) onlyCLevel external {
        _setData(dataAddress);
    }

    /** @dev Check if the address satisfiy the condition to be a candidate.
    * @param addrCandidate  Address of candidate.
    * @return can Bool value.
    */
    function _canBeCandidate(address addrCandidate) view private returns (bool can) {
        if(addrCandidate.balance >= data.candidateCondition()) {
            return true;
        }

        return false;
    }

    /** @dev Register to be candidate of verifier.
    * @param addrCandidate Address of candidate.
    */
    function registerCandidate(address addrCandidate) onlyCLevel external {
        _logDebug("registerCandidate", addrCandidate, address(0), 0, 0, 0, "");

        require(addrCandidate != address(0));
        require(_canBeCandidate(addrCandidate));

        data.addCandidate(addrCandidate);
    }

    /** @dev Get round of vote.
    * @return round Round of vote.
    */
    function getVoteRound() view external returns(uint256 round) {
        return data.getVoteRound();
    }
  
    /** @dev Get all candidates with tickets.
    * @return addresses The addresses of candidates.
    * @return tickets The tickets of candidates.
    */
    function getCandidates() external view returns(address[] memory addresses, uint256[] memory tickets) {
        return data.getCandidates();
    }

    /** @dev Vote for candidate.
    * @param addrVoter Address of voter.
    * @param addrCandidate Address of candidate.
    */
    function vote(address addrVoter, address addrCandidate) onlyCLevel external payable{
        _logInfo("vote", addrVoter, addrCandidate, msg.value, 0, 0, "");

        require(data.existCandidate(addrCandidate));
        require(msg.value > 0);

        data.increaseVoterDeposit.value(msg.value)(addrVoter);
        data.increaseVoterVote(addrVoter, addrCandidate, msg.value);
        data.updateVoterUnlockTime(addrVoter);
        data.increaseCandidateTicket(addrCandidate, msg.value);
    }

    /** @dev Cancel all votes of voter.
    * @param addrVoter Address of voter.
    * @param addresses Address of candidate.
    * @param tickets tickets of candidate.
    */
    function _cancelAllVotes(address addrVoter, address[] memory addresses, uint256[] memory tickets) private{
        uint256 len = addresses.length;
        for(uint256 index=0; index<len; index++){
            address addrCandidate = addresses[index];
            uint256 ticket = tickets[index];
            data.deleteVoterVote(addrVoter, addrCandidate);
            data.decreaseCandidateTicket(addrCandidate, ticket);
        }

        data.withdrawVoterDeposit(addrVoter);
    }

    /** @dev Cancel all votes of voter.
    * @param addrVoter Address of voter.
    */
    function cancelAllVotes(address addrVoter) onlyCLevel external {
        _logInfo("cancelAllVotes", addrVoter, address(0), 0, 0, 0, "");

        (address[] memory addresses, uint256[] memory tickets,) = data.getVoteInfo(addrVoter);
        _cancelAllVotes(addrVoter, addresses, tickets);       
    }

    function _isCurrentTickCycleEnd() view private returns (bool) {
        if(data.voteCycleBegin() + data.voteCyclePeriod() <= now) {
            return true;
        }
        return false;
    }

    /** @dev Check unclock deposit tokens.
    */
    function _checkUnlockDeposit() private {
        address[] memory addressesVoter = data.getVoters();

        uint256 len = addressesVoter.length;
        for(uint256 index=0; index<len; index++){
            address addrVoter = addressesVoter[index];

            (address[] memory addressesCandidate, uint256[] memory tickets, uint256 unlockTime) = data.getVoteInfo(addrVoter);
           
            //check if voter deposit can unlock;
            if(unlockTime.diff(now) <= data.unlockPrecision()){
                _cancelAllVotes(addrVoter, addressesCandidate, tickets);
            }
        }
    }

    /** @dev End current cycle, start next round
    * @return round Current round of vote.
    */
    function _rotateTickCycle() private returns (uint256 round) {
        data.setVoteCycleBegin(data.voteCycleBegin() + data.voteCyclePeriod());
        return data.increaseVoteCycleRound();
    }

    /** @dev Check if vote should be rotate.
    * @return needRotate If current vote round need rotate.
    */
    function tickVote() onlyCLevel view external returns (bool){
        if(!_isCurrentTickCycleEnd()) {
            return false;
        }

        return true;
    }

    /** @dev Rotate vote cycle.
    */
    function rotateVote() onlyCLevel external {
        _logDebug("rotateVote", msg.sender, address(0), 0, 0, 0, "");

        if(!_isCurrentTickCycleEnd()) {
            return;
        }

        //check unclock deposit tokens
        _checkUnlockDeposit();

        _rotateTickCycle();
    }
}