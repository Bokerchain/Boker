pragma solidity ^0.4.8;

import "./BokerDefine.sol";
import "./BokerManager.sol";
import "./BokerDappData.sol";

contract BokerTokenPowerData  is BokerManaged {
    // the begin time of current assign cycle;
    uint256 public assignCycleBegin = 0;
    uint256 public tokenAssignedTotal = 0;

    struct BindDappPowerEntry {
        int256 upper;      // upper limit.
        uint256 power;     // gained power.
    }

    struct BindDappConfig {
        uint256 dappType;
        BindDappPowerEntry[] entries;
    }

    mapping (uint256=>BindDappConfig) bindDappConfig;

    constructor(address addrManager) BokerManaged(addrManager) public {
        assignCycleBegin = now;
        
        uint256 dappType = 1;
        BindDappConfig storage config = bindDappConfig[dappType];
        config.dappType = dappType;
        config.entries.push(BindDappPowerEntry(10, 1 * POWER));
        config.entries.push(BindDappPowerEntry(50, 3 * POWER));
        config.entries.push(BindDappPowerEntry(100, 5 * POWER));
        config.entries.push(BindDappPowerEntry(500, 10 * POWER));
        config.entries.push(BindDappPowerEntry(1000, 15 * POWER));
        config.entries.push(BindDappPowerEntry(5000, 20 * POWER));
        config.entries.push(BindDappPowerEntry(10000, 30 * POWER));
        config.entries.push(BindDappPowerEntry(-1, 40 * POWER));

        dappType = 2;
        config = bindDappConfig[dappType];
        config.dappType = dappType;
        config.entries.push(BindDappPowerEntry(10, 1 * POWER));
        config.entries.push(BindDappPowerEntry(50, 3 * POWER));
        config.entries.push(BindDappPowerEntry(100, 5 * POWER));
        config.entries.push(BindDappPowerEntry(200, 10 * POWER));
        config.entries.push(BindDappPowerEntry(500, 20 * POWER));
        config.entries.push(BindDappPowerEntry(1000, 30 * POWER));
        config.entries.push(BindDappPowerEntry(-1, 40 * POWER));

        dappType = 3;
        config = bindDappConfig[dappType];
        config.dappType = dappType;
        config.entries.push(BindDappPowerEntry(10, 1 * POWER));
        config.entries.push(BindDappPowerEntry(50, 5 * POWER));
        config.entries.push(BindDappPowerEntry(100, 10 * POWER));
        config.entries.push(BindDappPowerEntry(200, 20 * POWER));
        config.entries.push(BindDappPowerEntry(500, 30 * POWER));
        config.entries.push(BindDappPowerEntry(-1, 40 * POWER));
    }
    
    function bindDappConfigGet(uint256 dappType) public view returns (int256[] memory uppers, uint256[] memory powers) {
        BindDappConfig storage config = bindDappConfig[dappType];
        if(config.dappType == DappTypeInvalid){
            return;
        }

        uint256 len = config.entries.length;
        uppers = new int256[](len);
        powers = new uint256[](len);
        for (uint256 index = 0; index < len; index++) {
            BindDappPowerEntry storage entry = config.entries[index];
            uppers[index] = entry.upper;
            powers[index] = entry.power;
        }
    }

    function bindDappConfigSet(uint256 dappType, uint256 index, int256 upper, uint256 power) onlyAdmin public {
        BindDappConfig storage config = bindDappConfig[dappType];
        if(config.dappType == DappTypeInvalid){
            config.dappType = dappType;
        }

        if(index >= config.entries.length) {
            config.entries.push(BindDappPowerEntry(upper, power));
        }else {
            config.entries[index].upper = upper;
            config.entries[index].power = power;
        }
    }

    function bindDappGetPower(address dappAddr) public view returns (uint256){
        (, uint256 dappType, uint256 dappUserCount, uint256 dappMonthlySales) = BokerDappData(contractAddress(ContractDappData)).dapps(dappAddr);
        if(dappType == DappTypeInvalid){
            return 0;
        }

        BindDappConfig storage config = bindDappConfig[dappType];
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
            BindDappPowerEntry storage entry = config.entries[index];
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

    function tokenAssignedTotalSet(uint256 val) external onlyContract {
        tokenAssignedTotal = val;
    }
}