pragma solidity ^0.4.8;

import "./BokerDefine.sol";
import "./BokerManager.sol";
import "./BokerDappData.sol";

contract BokerTokenPowerData  is BokerManaged {
    // the begin time of current assign cycle;
    uint256 public assignCycleBegin = 0;

    struct BindWalletPowerEntry {
        int256 upper;      // upper limit.
        uint256 power;     // gained power.
    }

    struct BindWalletConfig {
        uint256 dappType;
        BindWalletPowerEntry[] entries;
    }

    mapping (uint256=>BindWalletConfig) bindWalletConfig;

    constructor(address addrManager) BokerManaged(addrManager) public {
        assignCycleBegin = now;
        
        uint256 dappType = 1;
        BindWalletConfig storage config = bindWalletConfig[dappType];
        config.dappType = dappType;
        config.entries.push(BindWalletPowerEntry(10, 1));
        config.entries.push(BindWalletPowerEntry(50, 3));
        config.entries.push(BindWalletPowerEntry(100, 5));
        config.entries.push(BindWalletPowerEntry(500, 10));
        config.entries.push(BindWalletPowerEntry(1000, 15));
        config.entries.push(BindWalletPowerEntry(5000, 20));
        config.entries.push(BindWalletPowerEntry(10000, 30));
        config.entries.push(BindWalletPowerEntry(-1, 40));

        dappType = 2;
        config = bindWalletConfig[dappType];
        config.dappType = dappType;
        config.entries.push(BindWalletPowerEntry(10, 1));
        config.entries.push(BindWalletPowerEntry(50, 3));
        config.entries.push(BindWalletPowerEntry(100, 5));
        config.entries.push(BindWalletPowerEntry(200, 10));
        config.entries.push(BindWalletPowerEntry(500, 20));
        config.entries.push(BindWalletPowerEntry(1000, 30));
        config.entries.push(BindWalletPowerEntry(-1, 40));

        dappType = 3;
        config = bindWalletConfig[dappType];
        config.dappType = dappType;
        config.entries.push(BindWalletPowerEntry(10, 1));
        config.entries.push(BindWalletPowerEntry(50, 5));
        config.entries.push(BindWalletPowerEntry(100, 10));
        config.entries.push(BindWalletPowerEntry(200, 20));
        config.entries.push(BindWalletPowerEntry(500, 30));
        config.entries.push(BindWalletPowerEntry(-1, 40));
    }
    
    function bindWalletConfigGet(uint256 dappType) public view returns (int256[] memory uppers, uint256[] memory powers) {
        BindWalletConfig storage config = bindWalletConfig[dappType];
        if(config.dappType == DappTypeInvalid){
            return;
        }

        uint256 len = config.entries.length;
        uppers = new int256[](len);
        powers = new uint256[](len);
        for (uint256 index = 0; index < len; index++) {
            BindWalletPowerEntry storage entry = config.entries[index];
            uppers[index] = entry.upper;
            powers[index] = entry.power;
        }
    }

    function bindWalletConfigSet(uint256 dappType, uint256 index, int256 upper, uint256 power) onlyAdmin public {
        BindWalletConfig storage config = bindWalletConfig[dappType];
        if(config.dappType == DappTypeInvalid){
            config.dappType = dappType;
        }

        if(index >= config.entries.length) {
            config.entries.push(BindWalletPowerEntry(upper, power));
        }else {
            config.entries[index].upper = upper;
            config.entries[index].power = power;
        }
    }

    function bindWalletGetPower(address dappAddr) public view returns (uint256){
        (, uint256 dappType, uint256 dappUserCount, uint256 dappMonthlySales) = BokerDappData(contractAddress(ContractDappData)).dapps(dappAddr);
        if(dappType == DappTypeInvalid){
            return 0;
        }

        BindWalletConfig storage config = bindWalletConfig[dappType];
        if(config.dappType == DappTypeInvalid){
            return 0;
        }

        int256 lower = 0;
        int256 upper = 0;
        int256 value = 0;
        if(DappTypeHardware == dappType){
            value = int256(dappMonthlySales);
        }else{
            value = int256(dappUserCount);
        }

        for (uint256 index = 0; index < config.entries.length; index++) {
            BindWalletPowerEntry storage entry = config.entries[index];
            upper = entry.upper;

            if(-1 == upper) {
                if(value >= lower){
                    return entry.power;
                }
            }
            else {
                if((value >= lower) && (value < upper)){
                    return entry.power;
                }
            }

            lower = upper;
        }

        return 0;
    }

    /** @dev Set assign cycle begin time.
    * @param time Begin time of assign cycle.
    */
    function setAssignCycleBegin(uint256 time) external onlyContract {
        assignCycleBegin = time;
    }
}