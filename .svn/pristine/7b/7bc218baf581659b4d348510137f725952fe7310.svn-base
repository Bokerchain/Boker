pragma solidity ^0.4.8;

import "./BokerCommon.sol";

/**
 *  Contract Manager
 */
contract BokerContractManager is Ownable{

    mapping (string => address) private contracts;
    string[] names;

    function getContract(string cName) view public returns (address) {
        return contracts[cName];
    }

    function setContract(string cName, address addrManager) onlyOwner public {
        if(contracts[cName] == address(0)) {
            names.push(cName);
        }
        contracts[cName] = addrManager;
    }
    
    function getContractSize() view public returns (uint) {
        return names.length;
    }

    function getContractName(uint index) view public returns (string) {
        require(index < names.length, "index exceeds names length");
        return names[index];
    }
}

contract ContractManaged is Ownable{
    
    BokerContractManager public contractManager;
    string public contractName;

    constructor(string cName, address addrManager) public {
        require(addrManager != address(0));
        contractName = cName;
        contractManager = BokerContractManager(addrManager);
    }

    function getContract(string cName) view internal returns (address){
        return contractManager.getContract(cName);
    }

    function changeContractManager(address addrManager) onlyOwner external{
        contractManager = BokerContractManager(addrManager);
    }
}