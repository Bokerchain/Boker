pragma solidity ^0.4.8;

import "./BokerManager.sol";

contract BokerNodeData is BokerManaged {
    using SafeMath for uint256;

    // the begin time of current vote cycle;
    uint256 public voteCycleBegin = 0;

    // the round of vote cycle;
    uint256 public voteCycleRound = 0;

    constructor(address addrManager) BokerManaged(addrManager) public {
        //placeholder
        candidateArray.push(address(0));
        voterArray.push(address(0));

        voteCycleBegin = now;
    }    

    //candidates
    struct Candidate {
        address addr;           // address of candidate
        uint256 index;          // index of candidate
        uint256 tickets;        // tickets

        string description;     // description of node
        string team;            // description of team
        string name;            // name of node
    }
    mapping (address => Candidate) public candidates;
    address[] public candidateArray;

    /** @dev Add new candidate
    * @param addrCandidate Address of candidate.
    * @param description description of node
    * @param team description of team
    * @param name name of node
    */
    function addCandidate(address addrCandidate, string description, string team, string name) external onlyContract {
        Candidate storage candidate = candidates[addrCandidate];
        if(candidate.index > 0) {
            return;
        }
        candidate.addr = addrCandidate;
        candidate.index = candidateArray.length;
        candidate.tickets = 0;
        candidate.description = description;
        candidate.team = team;
        candidate.name = name;
        candidateArray.push(addrCandidate);
    }

        /** @dev Get all candidates with tickets.
    * @return addresses The addresses of candidates.
    * @return tickets The tickets of candidates.
    */
    function getCandidates() external view returns(address[] memory addresses, uint256[] memory tickets) {
        uint256 len = candidateArray.length;
        addresses = new address[](len - 1);
        tickets = new uint256[](len - 1);
        for(uint256 index = 1; index < len; index++){
            address addr = candidateArray[index];
            addresses[index - 1] = addr;
            tickets[index - 1] = candidates[addr].tickets;
        }
    }

    /** @dev Get candidate.
    * @param addrCandidate Address of candidate.
    * @return description description of node
    * @return team description of team
    * @return name name of node
    * @return tickets tickets of node
    */
    function getCandidate(address addrCandidate) external view  returns(string description, string team, string name, uint256 tickets) {
        Candidate storage candidate = candidates[addrCandidate];
        if(0 == candidate.index){
            return;
        }
        description = candidate.description;
        team = candidate.team;
        name = candidate.name;
        tickets = candidate.tickets;
        return;
    }

    /** @dev Check if candidate exists
    * @param addrCandidate address of candidate.
    * @return exist Bool value.
    */
    function existCandidate(address addrCandidate) external view returns (bool exist) {
        if(0 != candidates[addrCandidate].index){
            return true;
        }
        return false;
    }

    /** @dev Increase candidate ticket.
    * @param addrCandidate Address of candidate.
    * @param tokens Amount of tokens.
    */
    function increaseCandidateTicket(address addrCandidate, uint256 tokens) external onlyContract {
        Candidate storage candidate = candidates[addrCandidate];
        if(0 == candidate.index){
            return;
        }
        candidate.tickets = candidate.tickets.add(tokens);
    }

    /** @dev Decrease candidate ticket.
    * @param addrCandidate Address of candidate.
    * @param amount Ammount of tokens.
    */
    function decreaseCandidateTicket(address addrCandidate, uint256 amount) onlyContract external{
        Candidate storage candidate = candidates[addrCandidate];
        if(0 == candidate.index){
            return;
        }
        require(candidate.tickets >= amount, "candidate.tickets < amount!");

        candidate.tickets = candidate.tickets.sub(amount);
    }

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
    
    mapping(address => Voter) public voters;
    address[] public voterArray;

    function _findAddVote(Voter storage voter, address addrCandidate) private returns (Vote storage) {
        Vote storage vote = voter.votes[addrCandidate];
        vote.tickets = 0;
        vote.index = voter.candidateArray.length;
        voter.candidateArray.push(addrCandidate);
        return vote;
    }

    /** @dev Increase voter deposit.
    * @param addrVoter Address of candidate.
    * @param tokens tokens.
    */
    function increaseVoterDeposit(address addrVoter, uint256 tokens) external onlyContract {
        Voter storage voter = voters[addrVoter];
        if(0 == voter.index){            
            voter.addr = addrVoter;
            voter.deposit = 0;
            voter.unlockTime = voteCycleBegin + globalConfigInt(CfgVoteLockup);
            voter.index = voterArray.length;
            voterArray.push(addrVoter);

            //placeholder
            voters[addrVoter].candidateArray.push(address(0));
        }

        voter.deposit = voter.deposit.add(tokens);
    }

    /** @dev Increase voter vote info.
    * @param addrVoter Address of voter.
    * @param addrCandidate Address of candidate.
    * @param tokens Amount of tokens.
    */
    function increaseVoterVote(address addrVoter, address addrCandidate, uint256 tokens) external onlyContract {
        Voter storage voter = voters[addrVoter];
        if(0 == voter.index){
            return;
        }

        Vote storage vote = _findAddVote(voter, addrCandidate);
        vote.tickets = vote.tickets.add(tokens);
    }

    /** @dev Clear voter vote info.
    * @param addrVoter Address of voter.
    * @param addrCandidate Address of candidate.
    */
    function clearVoterVote(address addrVoter, address addrCandidate) external onlyContract {
        Voter storage voter = voters[addrVoter];
        if(0 == voter.index){
            return;
        }

        Vote storage vote = voter.votes[addrCandidate];
        vote.tickets = 0;
    }

    /** @dev Update voter unlock time.
    * @param addrVoter Address of voter.
    */
    function updateVoterUnlockTime(address addrVoter) external onlyContract {
        Voter storage voter = voters[addrVoter];
        if(0 == voter.index){
            return;
        }
        voter.unlockTime = voteCycleBegin + globalConfigInt(CfgVoteLockup);
    }

    /** @dev Get all vote info of voter.
    * @param addrVoter Address of voter.
    * @return addresses The addresses of candidates.
    * @return tickets The tickets of vote info.
    * @return unlockTime Unlock time.
    */
    function getVoteInfo(address addrVoter) external view returns(address[] addresses, uint256[] tickets, uint256 unlockTime, uint256 deposit) {
        Voter storage voter = voters[addrVoter];
        if(0 == voter.index){
            return;
        }

        unlockTime = voter.unlockTime;
        deposit = voter.deposit;
        uint256 len = voter.candidateArray.length;
        addresses = new address[](len - 1);
        tickets = new uint256[](len - 1);
        for(uint256 index = 1; index < len; index++){
            address addrCandidate = voter.candidateArray[index];
            addresses[index-1] = addrCandidate;
            tickets[index-1] = voter.votes[addrCandidate].tickets;
        }
    }

    /** @dev Get all voters.
    * @return addresses The addresses of voters.
    */
    function getVoters() external view returns(address[] addresses) {
        uint256 len = 0;

        for(uint index = 1; index < voterArray.length; index++) {
            address addrVoter = voterArray[index];
            if(voters[addrVoter].deposit > 0){
                len++;
            }
        }

        addresses = new address[](len);
        uint256 pos = 0;
        for(index = 1; index < voterArray.length; index++) {
            addrVoter = voterArray[index];
            if(voters[addrVoter].deposit > 0){
                addresses[pos] = addrVoter;
                pos++;
            }
        }
    }

    /** @dev clear Voter Deposit
    * @param addrVoter Address of voter.
    */
    function setVoterDeposit(address addrVoter, uint256 amount) external onlyContract {
        Voter storage voter = voters[addrVoter];
        if(0 == voter.index){
            return;
        }
        voter.deposit = amount;
    }

    /** @dev Set vote cycle begin time.
    * @param time Begin time of vote cycle.
    */
    function setVoteCycleBegin(uint256 time) external onlyContract {
        voteCycleBegin = time;
    }

    /** @dev Increase vote cycle round.
    * @return round New vote cycle round.
    */
    function increaseVoteCycleRound(uint256 roundAdd) external onlyContract returns (uint256) {
        voteCycleRound = voteCycleRound.add(roundAdd);
        return voteCycleRound;
    }
}