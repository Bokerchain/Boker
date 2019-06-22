pragma solidity ^0.4.8;

import "./BokerManager.sol";
import "./BokerLog.sol";

contract BokerFinance is BokerManaged {
    using SafeMath for uint256;
    
    constructor(address addrManager) BokerManaged(addrManager) public {
    }

    function () public payable {
        //日志
        BokerLog(contractAddress(ContractLog)).financeLogAdd(
            msg.sender, this, msg.value, msg.sender.balance, address(this).balance, uint256(FinanceReason.Transfer));
    }

    function grantTokenTo(address addrTo, uint256 amount, uint256 reason) external onlyContract {
        addrTo.transfer(amount);

        //日志
        BokerLog(contractAddress(ContractLog)).financeLogAdd(
            this, addrTo, amount, address(this).balance, addrTo.balance, reason);
    }

    function receiveTokenFrom(address addrFrom, uint256 reason) external payable onlyContract {
        //日志
        BokerLog(contractAddress(ContractLog)).financeLogAdd(
            addrFrom, this, msg.value, addrFrom.balance, address(this).balance, reason);
    }

    function withdraw() public onlyOwner {
        uint256 amount = address(this).balance;
        owner.transfer(address(this).balance);
        
        //日志
        BokerLog(contractAddress(ContractLog)).financeLogAdd(
            this, owner, amount, address(this).balance, owner.balance, uint256(FinanceReason.FinanceWithdraw));
    }
}