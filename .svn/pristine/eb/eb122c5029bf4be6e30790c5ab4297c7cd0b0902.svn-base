pragma solidity ^0.4.8;

import "./BokerContractManager.sol";
import "./BokerCommon.sol";

contract BokerModule is Ownable, ContractManaged, Utility{

    string public moduleName;
    string public moduleImplName;
    string public moduleDataName;

    // modifier of only impl is set
    modifier onlyImplSet {
        require(implAddress() != address(0));
        _;
    }

    // modifier of only data is set
    modifier onlyDataSet {
        require(dataAddress() != address(0));
        _;
    }

    // modifier of only invoked by authorized user.
    modifier onlyAuthorized {
        require(
            msg.sender == owner ||
            msg.sender == moduleAddress() ||
            msg.sender == implAddress()
        );
        _;
    }

    // modifier of only invoked by impl.
    modifier onlyImpl {
        require(msg.sender == implAddress());
        _;
    }

    constructor(string mName, string cName, address addrManager) ContractManaged(cName, addrManager) public {
        moduleName = mName;
        moduleImplName = strConcat(moduleName, "Impl");
        moduleDataName = strConcat(moduleName, "Data");
    }

    function moduleAddress() view public returns(address) {
        return getContract(moduleName);
    }

    function implAddress() view public returns(address) {
        return getContract(moduleImplName);
    }

    function dataAddress() view public returns(address) {
        return getContract(moduleDataName);
    }
}