pragma solidity ^0.4.8;

import "./BokerManager.sol";

contract BokerDappData is BokerManaged {

    struct Dapp {
        address dappAddr;
        uint256 dappType;
        uint256 userCount;
        uint256 monthlySales;
    }

    mapping (address=>Dapp) public dapps;
    address[] public dappAddresses;

    constructor(address addrManager) BokerManaged(addrManager) public {
    }

    function dappCount() public view returns (uint256) {
        return dappAddresses.length;
    } 

    function dappSet(address dappAddr, uint256 dappType, uint256 userCount, uint256 monthlySales) public onlyContract {
        Dapp storage dapp = dapps[dappAddr];
        if(dapp.dappAddr == address(0)) {
            dapp.dappAddr = dappAddr;
            dappAddresses.push(dappAddr);
        }
        dapp.dappType = dappType;
        dapp.userCount = userCount;
        dapp.monthlySales = monthlySales;
    }

    function dappGetAdresses() public view returns (address[] memory addrDapps) {
        uint256 len = dappAddresses.length;
        addrDapps = new address[](len);
        for (uint256 index = 0; index < len; index++) {
            addrDapps[index] = dappAddresses[index];
        }
    }
}