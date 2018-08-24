pragma solidity ^0.4.8;

import "./BokerCommon.sol";

/**
 *  Address Manager
 */
contract BokerAddressManager is AccessControl{

    mapping (string => address) private addresses;
    string[] names;

    function getAddress(string name) view public returns (address) {
        return addresses[name];
    }

    function setAddress(string name, address addr) onlyCLevel public {
        if(addresses[name] == address(0)) {
            names.push(name);
        }
        addresses[name] = addr;
    }
    
    function getAddressLen() view public returns (uint) {
        return names.length;
    }

    function getName(uint index) view public returns (string) {
        require(index < names.length, "index exceeds names length");
        return names[index];
    }
}

/**
 *  Address manageable
 */
contract AddressManageable is AccessControl{
    
    BokerAddressManager public addressManager;
    string public contractName;

    constructor(string name, address manager) public {
        contractName = name;
        addressManager = BokerAddressManager(manager);
    }

    function getAddress(string name) view public returns (address){
        return addressManager.getAddress(name);
    }

    function changeManager(string name) onlyCLevel view public returns (address){
        return addressManager.getAddress(name);
    }
}