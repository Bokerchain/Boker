pragma solidity ^0.4.8;


import "./BokerCommon.sol";
import "./BokerManager.sol";

interface INode {
    function getVoteRound() external view  returns(uint256 round);
    function getCandidates() external view  returns(address[] memory addresses, uint256[] memory tickets);
    function checkVote() external view returns (bool);
    function rotateVote() external;    
}

interface ITokenPower {
    function checkAssignToken() external view returns (bool);
    function assignToken() external;
}

interface IFinance {    
    function receiveTokenFrom(address addrFrom, uint256 reason) external payable;
}

contract BokerInterfaceBase is BokerManaged{
    constructor(address addrManager) BokerManaged(addrManager) public {
    }

    function () public payable {
    }

    /** @dev check if time to assgin token.
    * @return If true need invode assgin token.
    */
    function checkAssignToken() external view returns (bool) {
        return ITokenPower(contractAddress(ContractTokenPower)).checkAssignToken();
    }

    /** @dev Assign token called periodically by chain.
    */
    function assignToken() external whenNotPaused onlyAdmin {
        IFinance(contractAddress(ContractFinance)).receiveTokenFrom.value(address(this).balance)(this, uint256(FinanceReason.Mine));
        ITokenPower(contractAddress(ContractTokenPower)).assignToken();
    }

    /** @dev Get round of vote.
    * @return round Round of vote.
    */
    function getVoteRound() external view  returns(uint256 round) {
        return INode(contractAddress(ContractNode)).getVoteRound();
    }

    /** @dev Get all candidates with tickets.
    * @return addresses The addresses of candidates in current round.
    * @return tickets The tickets of candidates in current round.
    */
    function getCandidates() external view returns(address[] memory addresses, uint256[] memory tickets) {
        return INode(contractAddress(ContractNode)).getCandidates();
    }

    /** @dev Check if vote should be rotate.
    * @return needRotate If current vote round need rotate.
    */
    function tickVote() external view returns (bool) {
        return INode(contractAddress(ContractNode)).checkVote();
    }

    /** @dev Rotate vote cycle.
    */
    function rotateVote() external whenNotPaused onlyAdmin {
        INode(contractAddress(ContractNode)).rotateVote();
    }
}