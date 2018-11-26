pragma solidity ^0.4.8;

import "./BokerManager.sol";
import "./BokerDappData.sol";

contract BokerDapp is BokerManaged {

    constructor(address addrManager) BokerManaged(addrManager) public {
    }

    /** @dev Set dapp info.
    * @param dappAddr Address of from user.    
    * @param dappType Type of event.
    * @param userCount Address of to user.
    * @param monthlySales timestamps.
    */
    function dappSet(address dappAddr, uint256 dappType, uint256 userCount, uint256 monthlySales) external onlyContract {
        return BokerDappData(contractAddress(ContractDappData)).dappSet(dappAddr, dappType, userCount, monthlySales);
    }
}