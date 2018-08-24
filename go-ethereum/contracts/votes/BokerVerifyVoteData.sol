pragma solidity ^0.4.8;

import "../BokerCommon.sol";

contract BokerVerifyVoteData is AccessControl, Define, Log{
    using SafeMath for uint256;

    // Ensure that we are pointing to the right contract in our set call.
    bool public isBokerVerifyVoteData = true;

    // Define the minimum of tokens to be candidate.
    uint256 public candidateCondition = 1 ether;

    // Define deposit lock up period, unit: day.
    uint256 public voteLockup = 3 days;

    // Define precision of unlock time, unit: minute.
    uint256 public unlockPrecision = 5 minutes;

    // the begin time of current vote cycle;
    uint256 public voteCycleBegin = 0;

    // the period of vote cycle;
    uint256 public voteCyclePeriod = 1 days;

    // the round of vote cycle;
    uint256 public voteCycleRound = 0;

    // TODO event.

    struct Vote {
        uint256 index;                          // index of voter
        uint256 tickets;                        // ammount of tokens voted
    }

    struct Voter {
        address addr;                           // address of voter
        uint256 index;                          // index of voter
        uint256 deposit;                        // deposit tokens
        uint256 unlockTime;                     // time to unlock
        mapping (address=>Vote) votes;          // mapping of vote info, candidate address => Vote
        address[] candidateArray;               // array of candidate address .
    }

    struct Candidate {
        address addr;           // address of candidate
        uint256 index;          // index of candidate
        uint256 tickets;        // tickets
    }

    // mapping of voters from address to Voter structure.
    mapping(address => Voter) public voters;
    address[] public voterArray;

    // mapping of candidates from address to Candidate structure.
    mapping (address => Candidate) public candidates;
    // array of candidate address .
    address[] public candidateArray;
    
    constructor() public {
        //placeholder
        candidateArray.push(address(0));
        voterArray.push(address(0));

        voteCycleBegin = now;
    }

    /** @dev Set the condition to be a candidate.
    * @param condition  unit: ether.
    */
    function setCandidateCondition(uint256 condition) onlyCLevel external {
        candidateCondition = condition.mul(1 ether);
    }

    /** @dev Set the vote lock up period.
    * @param lockup  unit: days.
    */
    function setVoteLockup(uint256 lockup) onlyCLevel external {
        voteLockup = lockup.mul(1 days);
    }

    /** @dev Set the unlock precision.
    * @param precision  unit: minutes.
    */
    function setUnlockPrecision(uint256 precision) onlyCLevel external {
        unlockPrecision = precision.mul(1 minutes);
    }

    /** @dev Set vote cycle begin time.
    * @param time Begin time of vote cycle.
    */
    function setVoteCycleBegin(uint256 time) onlyCLevel external {
        voteCycleBegin = time;
    }

    /** @dev Set vote cycle period.
    * @param period Period of vote cycle.
    */
    function setVoteCyclePeriod(uint256 period) onlyCLevel external {
        voteCyclePeriod = period;
    }

    /** @dev Increase vote cycle round.
    * @return round New vote cycle round.
    */
    function increaseVoteCycleRound() onlyCLevel external returns (uint256) {
        voteCycleRound = voteCycleRound + 1;
        return voteCycleRound;
    }

    /** @dev Check if candidate exists
    * @param addrCandidate address of candidate.
    * @return exist Bool value.
    */
    function existCandidate(address addrCandidate) view external returns (bool exist) {
        if(0 != candidates[addrCandidate].index){
            return true;
        }
        return false;
    }

    /** @dev Add new candidate
    * @param addrCandidate address of candidate.
    */
    function addCandidate(address addrCandidate) onlyCLevel external {

        //Already added.
        if(0 != candidates[addrCandidate].index){
            return;
        }

        Candidate memory candidate;
        candidate.index = candidateArray.length;
        candidate.addr = addrCandidate;
        candidate.tickets = 0;
        candidates[addrCandidate] = candidate;
        candidateArray.push(addrCandidate);
    }

    /** @dev Get round of vote.
    * @return round Round of vote.
    */
    function getVoteRound() view external returns(uint256 round) {
        return voteCycleRound;
    }
  
    /** @dev Get all candidates with tickets.
    * @return addresses The addresses of candidates.
    * @return tickets The tickets of candidates.
    */
    function getCandidates() external view returns(address[] memory addresses, uint256[] memory tickets) {
        uint256 len = candidateArray.length;
        addresses = new address[](len - 1);
        tickets = new uint256[](len - 1);
        for(uint256 index=1; index<len; index++){
            address addr = candidateArray[index];
            addresses[index - 1] = addr;
            tickets[index - 1] = candidates[addr].tickets;
        }

        return (addresses, tickets);
    }

    /** @dev Add Voter.
    * @param addrVoter Address of voter.
    */
    function _addVoter(address addrVoter) private returns (Voter storage) {
        Voter memory voter;
        voter.index = voterArray.length;
        voter.addr = addrVoter;
        voter.deposit = 0;
        voter.unlockTime = voteCycleBegin + voteLockup;
        voters[addrVoter] = voter;
        voterArray.push(addrVoter);

        //placeholder
        voters[addrVoter].candidateArray.push(address(0));
        
        return voters[addrVoter];
    }

    /** @dev Get all voters.
    * @return addresses The addresses of voters.
    */
    function getVoters() external view returns(address[]) {
        uint256 len = 0;

        for(uint index=1; index<voterArray.length; index++) {
            address addrVoter = voterArray[index];
            if(voters[addrVoter].deposit > 0){
                len++;
            }
        }

        address[] memory addresses = new address[](len);
        uint256 pos = 0;
        for(index=1; index<voterArray.length; index++) {
            addrVoter = voterArray[index];
            if(voters[addrVoter].deposit > 0){
                addresses[pos] = addrVoter;
                pos++;
            }
        }

        return (addresses);
    }

    /** @dev Increase voter deposit.
    * @param addrVoter Address of candidate.
    */
    function increaseVoterDeposit(address addrVoter) onlyCLevel external payable{
        _logInfo("increaseVoterDeposit", addrVoter, address(0), msg.value, 0, 0, "");

        Voter storage voter = voters[addrVoter];
        if(0 == voter.index){
            voter = _addVoter(addrVoter);
        }

        voter.deposit = voter.deposit.add(msg.value);
    }

    /** @dev Withdraw voter deposit.
    * @param addrVoter Address of candidate.
    */
    function withdrawVoterDeposit(address addrVoter) onlyCLevel external{
        Voter storage voter = voters[addrVoter];
        require(voter.index > 0);        

        uint deposit = voter.deposit;
        require(address(this).balance >= deposit);

        voter.deposit = 0;
        addrVoter.transfer(deposit);
    }

    /** @dev Add Vote.
    * @param voter Pointer to voter.
    * @param addrCandidate Address of candidate.
    */
    function _addVote(Voter storage voter, address addrCandidate) private returns (Vote storage) {        
        Vote memory vote;
        vote.index = voter.candidateArray.length;
        vote.tickets = 0;
        voter.votes[addrCandidate] = vote;
        voter.candidateArray.push(addrCandidate);
        return voter.votes[addrCandidate];
    }

    /** @dev Update voter unlock time.
    * @param addrVoter Address of voter.
    */
    function updateVoterUnlockTime(address addrVoter) onlyCLevel external{
        Voter storage voter = voters[addrVoter];
        if(0 == voter.index){
            return;
        }
        voter.unlockTime = voteCycleBegin + voteLockup;
    }

    /** @dev Increase voter vote info.
    * @param addrVoter Address of voter.
    * @param addrCandidate Address of candidate.
    * @param amount Amount of tokens.
    */
    function increaseVoterVote(address addrVoter, address addrCandidate, uint256 amount) onlyCLevel external{
        Voter storage voter = voters[addrVoter];
        if(0 == voter.index){
            return;
        }

        Vote storage vote = voter.votes[addrCandidate];
        if(0 == vote.index) {
            vote = _addVote(voter, addrCandidate);
        }
        vote.tickets = vote.tickets.add(amount);
    }

    /** @dev Delete voter vote info.
    * @param addrVoter Address of voter.
    * @param addrCandidate Address of candidate.
    */
    function deleteVoterVote(address addrVoter, address addrCandidate) onlyCLevel external{
        Voter storage voter = voters[addrVoter];
        if(0 == voter.index){
            return;
        }

        Vote storage vote = voter.votes[addrCandidate];
        vote.tickets = 0;
    }

    /** @dev Get all vote info of voter.
    * @param addrVoter Address of voter.
    * @return addresses The addresses of candidates.
    * @return tickets The tickets of vote info.
    * @return unlockTime Unlock time.
    */
    function getVoteInfo(address addrVoter) view external returns(address[], uint256[], uint256) {
        Voter storage voter = voters[addrVoter];
        if(0 == voter.index){
            return;
        }

        uint256 len = voter.candidateArray.length;

        address[] memory addresses = new address[](len - 1);
        uint256[] memory tickets = new uint256[](len - 1);
        for(uint256 index=1; index<len; index++){
            address addrCandidate = voter.candidateArray[index];
            addresses[index-1] = addrCandidate;
            tickets[index-1] = voter.votes[addrCandidate].tickets;
        }

        return (addresses, tickets, voter.unlockTime);
    }

    /** @dev Increase candidate ticket.
    * @param addrCandidate Address of candidate.
    * @param amount Amount of tokens.
    */
    function increaseCandidateTicket(address addrCandidate, uint256 amount) onlyCLevel external{
        Candidate storage candidate = candidates[addrCandidate];
        if(0 == candidate.index){
            return;
        }

        candidate.tickets = candidate.tickets.add(amount);
    }

    /** @dev Decrease candidate ticket.
    * @param addrCandidate Address of candidate.
    * @param amount Ammount of tokens.
    */
    function decreaseCandidateTicket(address addrCandidate, uint256 amount) onlyCLevel external{
        Candidate storage candidate = candidates[addrCandidate];
        if(0 == candidate.index){
            return;
        }
        require(candidate.tickets >= amount);

        candidate.tickets = candidate.tickets.sub(amount);
    }
}