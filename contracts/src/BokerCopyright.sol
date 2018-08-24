pragma solidity ^0.4.8;

import "./BokerCommon.sol";

contract BobbyCopyright {

    struct CopyrightInfo {
        uint256 id;
        string sign;
        address owner;
    }

    struct AuthorizeInfo {
        string[] channels;       //渠道d列表
    }

    mapping(string=>CopyrightInfo) private signToCopyright;
    string[] private copyrightSings;
    mapping (uint256 => address) private copyrightIdToOwner;


    // CopyrightInfo[] private copyrights;
    // mapping (uint256 => address) private copyrightIdToOwner
    // mapping (string => uint256) private copyrightSignToId


    //     AuthorizeInfo[] private authorizations;

    //事件
    event Register(address owner, string signature);   //确权
    event Authorize(address owner, string signature, string channel);   //授权

    function() public {

    }

    /**
     * query all copyrights of user
     */
    function copyrightsOfUser(address addr) public view returns(uint256[] copyrights) {
    }

    /**
     * query all authorizations of copyrights
     */
    function authorizationsOfCopyright(address addr) public view returns(uint256[] authorizations) {
    }

    /**
     * register
     */
    function register(address owner, string signature) public {

    }

    /**
     * register
     */
    function authorize(address owner, string signature, string channel) public {

    }
}
