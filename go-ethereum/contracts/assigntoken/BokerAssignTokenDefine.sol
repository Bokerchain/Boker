pragma solidity ^0.4.8;

import "../BokerCommon.sol";

contract BokerAssignTokenDefine is BokerDefine {

    uint256 constant public assgineTokenPeriodDefault = 5 minutes;
    uint256 constant public assginedTokensPerPeriodDefault = 990 * bobby;

    struct User {
        address addr;                               // address of user
        uint256 index;                              // index in array of user
        uint256 registerTime;                       // time of register
        uint256 certificationTime;                  // time of real name certification
        uint256 bindWalletTime;                     // time of bind wallet
        uint256 lastLoginTime;                      
        uint256 invitedFriendsCount;
        uint256 watchTime;                          // total watch time current assign cycle
        uint256 uploadCount;                        // total upload count current assign cycle

        uint256 longtermPower;                      // long term power, always exists
        uint256 shorttermPower;                     // short term power, cleared at the end of every assgin cycle.
    }

    struct UserEvent {
        uint256 eventType;
        address addrFrom;
        address addrTo;        
        uint256 eventValue1;
        uint256 eventValue2;
    }

    enum UserEventType {
        Register,
        LoginDaily,
        Certification,          // real name  certification.
        BindWallet,
        Watch,
        Upload
    }

    enum UserPowerType {
        Longterm,
        Shortterm
    }
}