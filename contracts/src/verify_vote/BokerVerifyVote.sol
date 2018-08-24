pragma solidity ^0.4.8;

import "./BokerVerifyVoteImpl.sol";
import "../BokerCommon.sol";

contract BokerVerifyVote is AccessControl, Log {

    BokerVerifyVoteImpl public impl;

    // modifier of only implemented
    modifier onlyImplemented {
        require(address(impl) != address(0));
        _;
    }

    constructor() public {
    }

    function _setImpl(address implAddress) private {
        BokerVerifyVoteImpl implContract = BokerVerifyVoteImpl(implAddress);

        // verify that a contract is what we expect
        require(implContract.isBokerVerifyVoteImpl());

        impl = implContract;
    }

    function setImpl(address implAddress) onlyCLevel external {
        _setImpl(implAddress);
    }

    /** @dev Register to be candidate of verifier.
    */
    function registerCandidate() onlyImplemented external {
        _logDebug("registerCandidate", msg.sender, address(0), 0, 0, 0, "");

        impl.registerCandidate(msg.sender);
    } 

    /** @dev Get all candidates with tickets.
    * @return round Round of vote.
    * @return addresses The addresses of candidates in current round.
    * @return tickets The tickets of candidates in current round.
    */
    function getCandidates() onlyImplemented view external returns(uint256 round, address[] memory addresses, uint256[] memory tickets) {
        return impl.getCandidates();
    }

    /** @dev Vote for candidate.
    * @param addrCandidate Address of candidate.
    */
    function voteCandidate(address addrCandidate) onlyImplemented external payable{
        _logInfo("voteCandidate", msg.sender, addrCandidate, msg.value, 0, 0, "");

        impl.vote.value(msg.value)(msg.sender, addrCandidate);
    }

    /** @dev Cancel all votes of voter.
    */
    function cancelAllVotes() onlyImplemented external{
        _logInfo("cancelAllVotes", msg.sender, address(0), 0, 0, 0, "");

        impl.cancelAllVotes(msg.sender);
    }

    /** @dev Check if vote should be rotate.
    * @return needRotate If current vote round need rotate.
    */
    function tickVote() onlyImplemented view external returns (bool) {
        return impl.tickVote();
    }

    /** @dev Rotate vote cycle.
    */
    function rotateVote() onlyImplemented external{
        _logDebug("rotateVote", msg.sender, address(0), 0, 0, 0, "");

        return impl.rotateVote();
    }
}