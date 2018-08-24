// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package votes

import (
	"math/big"
	"strings"

	"github.com/boker/go-ethereum/accounts/abi"
	"github.com/boker/go-ethereum/accounts/abi/bind"
	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/core/types"
)

// VotesABI is the input ABI used to generate the binding from.
const VotesABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"cfoAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getCandidates\",\"outputs\":[{\"name\":\"addresses\",\"type\":\"address[]\"},{\"name\":\"tickets\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ceoAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getLogSize\",\"outputs\":[{\"name\":\"size\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"rotateVote\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"keyDefault\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"enabled\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addrCandidate\",\"type\":\"address\"}],\"name\":\"voteCandidate\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newCEO\",\"type\":\"address\"}],\"name\":\"setCEO\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newCOO\",\"type\":\"address\"}],\"name\":\"setCOO\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"enable\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getLog\",\"outputs\":[{\"name\":\"level\",\"type\":\"uint8\"},{\"name\":\"time\",\"type\":\"uint256\"},{\"name\":\"key\",\"type\":\"string\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"v1\",\"type\":\"uint256\"},{\"name\":\"v2\",\"type\":\"uint256\"},{\"name\":\"v3\",\"type\":\"uint256\"},{\"name\":\"remarks\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newCFO\",\"type\":\"address\"}],\"name\":\"setCFO\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"clearLog\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"implAddress\",\"type\":\"address\"}],\"name\":\"setImpl\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"logLevel\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tickVote\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"impl\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"cooAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"level\",\"type\":\"uint8\"}],\"name\":\"setLevel\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getVoteRound\",\"outputs\":[{\"name\":\"round\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"cancelAllVotes\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"registerCandidate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"}]"

// VotesBin is the compiled bytecode used for deploying new contracts.
const VotesBin = `60806040526000600260146101000a81548160ff0219169083151502179055506001600460006101000a81548160ff021916908360ff1602179055506000600460016101000a81548160ff021916908360ff1602179055506040805190810160405280600781526020017f64656661756c740000000000000000000000000000000000000000000000000081525060059080519060200190620000a4929190620000f9565b50348015620000b257600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550620001a8565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200013c57805160ff19168380011785556200016d565b828001600101855582156200016d579182015b828111156200016c5782518255916020019190600101906200014f565b5b5090506200017c919062000180565b5090565b620001a591905b80821115620001a157600081600090555060010162000187565b5090565b90565b61258080620001b86000396000f300608060405260043610610154576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630519ce791461015957806306a49fce146101b05780630a0f8168146102645780630c73a392146102bb57806311070fcd146102e65780631a08a37c146102fd578063238dafe01461038d57806326bb886d146103be57806327d7874c146103f45780632ba73c15146104375780632f038fd51461047a5780633206b2c6146104aa5780633f4ba83a1461064b5780634e0a33791461067a5780635c50745e146106bd5780635c975abb146106d4578063691bd2ae146107035780637567772b146107465780638456cb59146107775780638765da94146107a65780638abf6077146107d5578063b047fb501461082c578063bd5546be14610883578063f0507573146108b3578063f49f9df8146108de578063f7e0079e146108f5575b600080fd5b34801561016557600080fd5b5061016e61090c565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156101bc57600080fd5b506101c5610932565b604051808060200180602001838103835285818151815260200191508051906020019060200280838360005b8381101561020c5780820151818401526020810190506101f1565b50505050905001838103825284818151815260200191508051906020019060200280838360005b8381101561024e578082015181840152602081019050610233565b5050505090500194505050505060405180910390f35b34801561027057600080fd5b50610279610b01565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156102c757600080fd5b506102d0610b26565b6040518082815260200191505060405180910390f35b3480156102f257600080fd5b506102fb610b33565b005b34801561030957600080fd5b50610312610c87565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610352578082015181840152602081019050610337565b50505050905090810190601f16801561037f5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561039957600080fd5b506103a2610d25565b604051808260ff1660ff16815260200191505060405180910390f35b6103f2600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610d38565b005b34801561040057600080fd5b50610435600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610ef8565b005b34801561044357600080fd5b50610478600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610fd2565b005b34801561048657600080fd5b506104a8600480360381019080803560ff1690602001909291905050506110ad565b005b3480156104b657600080fd5b506104d5600480360381019080803590602001909291905050506111d6565b604051808a60ff1660ff168152602001898152602001806020018873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018681526020018581526020018481526020018060200183810383528a818151815260200191508051906020019080838360005b838110156105a1578082015181840152602081019050610586565b50505050905090810190601f1680156105ce5780820380516001836020036101000a031916815260200191505b50838103825284818151815260200191508051906020019080838360005b838110156106075780820151818401526020810190506105ec565b50505050905090810190601f1680156106345780820380516001836020036101000a031916815260200191505b509b50505050505050505050505060405180910390f35b34801561065757600080fd5b506106606113f3565b604051808215151515815260200191505060405180910390f35b34801561068657600080fd5b506106bb600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506114b9565b005b3480156106c957600080fd5b506106d2611594565b005b3480156106e057600080fd5b506106e96116af565b604051808215151515815260200191505060405180910390f35b34801561070f57600080fd5b50610744600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506116c2565b005b34801561075257600080fd5b5061075b6117d9565b604051808260ff1660ff16815260200191505060405180910390f35b34801561078357600080fd5b5061078c6117ec565b604051808215151515815260200191505060405180910390f35b3480156107b257600080fd5b506107bb6118b3565b604051808215151515815260200191505060405180910390f35b3480156107e157600080fd5b506107ea6119d8565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561083857600080fd5b506108416119fe565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561088f57600080fd5b506108b1600480360381019080803560ff169060200190929190505050611a24565b005b3480156108bf57600080fd5b506108c8611b4d565b6040518082815260200191505060405180910390f35b3480156108ea57600080fd5b506108f3611c72565b005b34801561090157600080fd5b5061090a611dfd565b005b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b606080600073ffffffffffffffffffffffffffffffffffffffff16600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415151561099357600080fd5b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166306a49fce6040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401600060405180830381600087803b158015610a1957600080fd5b505af1158015610a2d573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f820116820180604052506040811015610a5757600080fd5b810190808051640100000000811115610a6f57600080fd5b82810190506020810184811115610a8557600080fd5b8151856020820283011164010000000082111715610aa257600080fd5b50509291906020018051640100000000811115610abe57600080fd5b82810190506020810184811115610ad457600080fd5b8151856020820283011164010000000082111715610af157600080fd5b5050929190505050915091509091565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000600380549050905090565b600073ffffffffffffffffffffffffffffffffffffffff16600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614151515610b9157600080fd5b610be76040805190810160405280600a81526020017f726f74617465566f746500000000000000000000000000000000000000000000815250336000806000806020604051908101604052806000815250611f88565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166311070fcd6040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401600060405180830381600087803b158015610c6d57600080fd5b505af1158015610c81573d6000803e3d6000fd5b50505050565b60058054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610d1d5780601f10610cf257610100808354040283529160200191610d1d565b820191906000526020600020905b815481529060010190602001808311610d0057829003601f168201915b505050505081565b600460009054906101000a900460ff1681565b600073ffffffffffffffffffffffffffffffffffffffff16600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614151515610d9657600080fd5b610deb6040805190810160405280600d81526020017f766f746543616e646964617465000000000000000000000000000000000000008152503383346000806020604051908101604052806000815250611fa2565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16631afa74e53433846040518463ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001925050506000604051808303818588803b158015610edc57600080fd5b505af1158015610ef0573d6000803e3d6000fd5b505050505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610f5357600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614151515610f8f57600080fd5b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561102d57600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415151561106957600080fd5b80600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16148061115557506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b806111ad5750600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b15156111b857600080fd5b80600460006101000a81548160ff021916908360ff16021790555050565b60008060606000806000806000606060006003805490508b1015156111fa57600080fd5b60008b1015151561120a57600080fd5b60038b81548110151561121957fe5b906000526020600020906009020190508060010160009054906101000a900460ff169950806002018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156112d55780601f106112aa576101008083540402835291602001916112d5565b820191906000526020600020905b8154815290600101906020018083116112b857829003601f168201915b50505050509750806000015498508060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1696508060040160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169550806005015494508060060154935080600701549250806008018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156113de5780601f106113b3576101008083540402835291602001916113de565b820191906000526020600020905b8154815290600101906020018083116113c157829003601f168201915b50505050509150509193959799909294969850565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561145057600080fd5b600260149054906101000a900460ff16151561146b57600080fd5b6000600260146101000a81548160ff0219169083151502179055507f7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b3360405160405180910390a16001905090565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561151457600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415151561155057600080fd5b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16148061163c57506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b806116945750600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b151561169f57600080fd5b600360006116ad9190612306565b565b600260149054906101000a900460ff1681565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16148061176a57506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b806117c25750600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b15156117cd57600080fd5b6117d681611fbc565b50565b600460019054906101000a900460ff1681565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561184957600080fd5b600260149054906101000a900460ff1615151561186557600080fd5b6001600260146101000a81548160ff0219169083151502179055507f6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff62560405160405180910390a16001905090565b60008073ffffffffffffffffffffffffffffffffffffffff16600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415151561191257600080fd5b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16638765da946040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401602060405180830381600087803b15801561199857600080fd5b505af11580156119ac573d6000803e3d6000fd5b505050506040513d60208110156119c257600080fd5b8101908080519060200190929190505050905090565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161480611acc57506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b80611b245750600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b1515611b2f57600080fd5b80600460016101000a81548160ff021916908360ff16021790555050565b60008073ffffffffffffffffffffffffffffffffffffffff16600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614151515611bac57600080fd5b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663f05075736040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401602060405180830381600087803b158015611c3257600080fd5b505af1158015611c46573d6000803e3d6000fd5b505050506040513d6020811015611c5c57600080fd5b8101908080519060200190929190505050905090565b600073ffffffffffffffffffffffffffffffffffffffff16600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614151515611cd057600080fd5b611d266040805190810160405280600e81526020017f63616e63656c416c6c566f746573000000000000000000000000000000000000815250336000806000806020604051908101604052806000815250611fa2565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663281c231b336040518263ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001915050600060405180830381600087803b158015611de357600080fd5b505af1158015611df7573d6000803e3d6000fd5b50505050565b600073ffffffffffffffffffffffffffffffffffffffff16600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614151515611e5b57600080fd5b611eb16040805190810160405280601181526020017f726567697374657243616e646964617465000000000000000000000000000000815250336000806000806020604051908101604052806000815250611f88565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166384aa3b93336040518263ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001915050600060405180830381600087803b158015611f6e57600080fd5b505af1158015611f82573d6000803e3d6000fd5b50505050565b611f996001888888888888886120b0565b50505050505050565b611fb36002888888888888886120b0565b50505050505050565b60008190508073ffffffffffffffffffffffffffffffffffffffff16632723ec486040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401602060405180830381600087803b15801561202557600080fd5b505af1158015612039573d6000803e3d6000fd5b505050506040513d602081101561204f57600080fd5b8101908080519060200190929190505050151561206b57600080fd5b80600660006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050565b6120b861232a565b600460009054906101000a900460ff1660ff16600014156120d8576122fb565b8860ff16600460019054906101000a900460ff1660ff1611156120fa576122fb565b88816020019060ff16908160ff16815250508781604001819052504281600001818152505086816060019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505085816080019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050848160a0018181525050838160c0018181525050828160e001818152505081816101000181905250600381908060018154018082558091505090600182039060005260206000209060090201600090919290919091506000820151816000015560208201518160010160006101000a81548160ff021916908360ff160217905550604082015181600201908051906020019061222c9291906123a6565b5060608201518160030160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060808201518160040160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060a0820151816005015560c0820151816006015560e082015181600701556101008201518160080190805190602001906122f69291906123a6565b505050505b505050505050505050565b50805460008255600902906000526020600020908101906123279190612426565b50565b6101206040519081016040528060008152602001600060ff16815260200160608152602001600073ffffffffffffffffffffffffffffffffffffffff168152602001600073ffffffffffffffffffffffffffffffffffffffff168152602001600081526020016000815260200160008152602001606081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106123e757805160ff1916838001178555612415565b82800160010185558215612415579182015b828111156124145782518255916020019190600101906123f9565b5b50905061242291906124e7565b5090565b6124e491905b808211156124e0576000808201600090556001820160006101000a81549060ff0219169055600282016000612461919061250c565b6003820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690556004820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690556005820160009055600682016000905560078201600090556008820160006124d7919061250c565b5060090161242c565b5090565b90565b61250991905b808211156125055760008160009055506001016124ed565b5090565b90565b50805460018160011615610100020316600290046000825580601f106125325750612551565b601f01602090049060005260206000209081019061255091906124e7565b5b505600a165627a7a72305820b9b29f156f9418ecd1f9921c50d6efd9c79c89daed5fb4a124bb72a4c50a33e50029`

// DeployVotes deploys a new Ethereum contract, binding an instance of Votes to it.
func DeployVotes(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Votes, error) {
	parsed, err := abi.JSON(strings.NewReader(VotesABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(VotesBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Votes{VotesCaller: VotesCaller{contract: contract}, VotesTransactor: VotesTransactor{contract: contract}}, nil
}

// Votes is an auto generated Go binding around an Ethereum contract.
type Votes struct {
	VotesCaller     // Read-only binding to the contract
	VotesTransactor // Write-only binding to the contract
}

// VotesCaller is an auto generated read-only Go binding around an Ethereum contract.
type VotesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VotesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VotesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VotesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VotesSession struct {
	Contract     *Votes            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VotesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VotesCallerSession struct {
	Contract *VotesCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// VotesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VotesTransactorSession struct {
	Contract     *VotesTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VotesRaw is an auto generated low-level Go binding around an Ethereum contract.
type VotesRaw struct {
	Contract *Votes // Generic contract binding to access the raw methods on
}

// VotesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VotesCallerRaw struct {
	Contract *VotesCaller // Generic read-only contract binding to access the raw methods on
}

// VotesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VotesTransactorRaw struct {
	Contract *VotesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVotes creates a new instance of Votes, bound to a specific deployed contract.
func NewVotes(address common.Address, backend bind.ContractBackend) (*Votes, error) {
	contract, err := bindVotes(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Votes{VotesCaller: VotesCaller{contract: contract}, VotesTransactor: VotesTransactor{contract: contract}}, nil
}

// NewVotesCaller creates a new read-only instance of Votes, bound to a specific deployed contract.
func NewVotesCaller(address common.Address, caller bind.ContractCaller) (*VotesCaller, error) {
	contract, err := bindVotes(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &VotesCaller{contract: contract}, nil
}

// NewVotesTransactor creates a new write-only instance of Votes, bound to a specific deployed contract.
func NewVotesTransactor(address common.Address, transactor bind.ContractTransactor) (*VotesTransactor, error) {
	contract, err := bindVotes(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &VotesTransactor{contract: contract}, nil
}

// bindVotes binds a generic wrapper to an already deployed contract.
func bindVotes(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VotesABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Votes *VotesRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Votes.Contract.VotesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Votes *VotesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Votes.Contract.VotesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Votes *VotesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Votes.Contract.VotesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Votes *VotesCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Votes.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Votes *VotesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Votes.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Votes *VotesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Votes.Contract.contract.Transact(opts, method, params...)
}

// CeoAddress is a free data retrieval call binding the contract method 0x0a0f8168.
//
// Solidity: function ceoAddress() constant returns(address)
func (_Votes *VotesCaller) CeoAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Votes.contract.Call(opts, out, "ceoAddress")
	return *ret0, err
}

// CeoAddress is a free data retrieval call binding the contract method 0x0a0f8168.
//
// Solidity: function ceoAddress() constant returns(address)
func (_Votes *VotesSession) CeoAddress() (common.Address, error) {
	return _Votes.Contract.CeoAddress(&_Votes.CallOpts)
}

// CeoAddress is a free data retrieval call binding the contract method 0x0a0f8168.
//
// Solidity: function ceoAddress() constant returns(address)
func (_Votes *VotesCallerSession) CeoAddress() (common.Address, error) {
	return _Votes.Contract.CeoAddress(&_Votes.CallOpts)
}

// CfoAddress is a free data retrieval call binding the contract method 0x0519ce79.
//
// Solidity: function cfoAddress() constant returns(address)
func (_Votes *VotesCaller) CfoAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Votes.contract.Call(opts, out, "cfoAddress")
	return *ret0, err
}

// CfoAddress is a free data retrieval call binding the contract method 0x0519ce79.
//
// Solidity: function cfoAddress() constant returns(address)
func (_Votes *VotesSession) CfoAddress() (common.Address, error) {
	return _Votes.Contract.CfoAddress(&_Votes.CallOpts)
}

// CfoAddress is a free data retrieval call binding the contract method 0x0519ce79.
//
// Solidity: function cfoAddress() constant returns(address)
func (_Votes *VotesCallerSession) CfoAddress() (common.Address, error) {
	return _Votes.Contract.CfoAddress(&_Votes.CallOpts)
}

// CooAddress is a free data retrieval call binding the contract method 0xb047fb50.
//
// Solidity: function cooAddress() constant returns(address)
func (_Votes *VotesCaller) CooAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Votes.contract.Call(opts, out, "cooAddress")
	return *ret0, err
}

// CooAddress is a free data retrieval call binding the contract method 0xb047fb50.
//
// Solidity: function cooAddress() constant returns(address)
func (_Votes *VotesSession) CooAddress() (common.Address, error) {
	return _Votes.Contract.CooAddress(&_Votes.CallOpts)
}

// CooAddress is a free data retrieval call binding the contract method 0xb047fb50.
//
// Solidity: function cooAddress() constant returns(address)
func (_Votes *VotesCallerSession) CooAddress() (common.Address, error) {
	return _Votes.Contract.CooAddress(&_Votes.CallOpts)
}

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() constant returns(uint8)
func (_Votes *VotesCaller) Enabled(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Votes.contract.Call(opts, out, "enabled")
	return *ret0, err
}

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() constant returns(uint8)
func (_Votes *VotesSession) Enabled() (uint8, error) {
	return _Votes.Contract.Enabled(&_Votes.CallOpts)
}

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() constant returns(uint8)
func (_Votes *VotesCallerSession) Enabled() (uint8, error) {
	return _Votes.Contract.Enabled(&_Votes.CallOpts)
}

// GetCandidates is a free data retrieval call binding the contract method 0x06a49fce.
//
// Solidity: function getCandidates() constant returns(addresses address[], tickets uint256[])
func (_Votes *VotesCaller) GetCandidates(opts *bind.CallOpts) (struct {
	Addresses []common.Address
	Tickets   []*big.Int
}, error) {
	ret := new(struct {
		Addresses []common.Address
		Tickets   []*big.Int
	})
	out := ret
	err := _Votes.contract.Call(opts, out, "getCandidates")
	return *ret, err
}

// GetCandidates is a free data retrieval call binding the contract method 0x06a49fce.
//
// Solidity: function getCandidates() constant returns(addresses address[], tickets uint256[])
func (_Votes *VotesSession) GetCandidates() (struct {
	Addresses []common.Address
	Tickets   []*big.Int
}, error) {
	return _Votes.Contract.GetCandidates(&_Votes.CallOpts)
}

// GetCandidates is a free data retrieval call binding the contract method 0x06a49fce.
//
// Solidity: function getCandidates() constant returns(addresses address[], tickets uint256[])
func (_Votes *VotesCallerSession) GetCandidates() (struct {
	Addresses []common.Address
	Tickets   []*big.Int
}, error) {
	return _Votes.Contract.GetCandidates(&_Votes.CallOpts)
}

// GetLog is a free data retrieval call binding the contract method 0x3206b2c6.
//
// Solidity: function getLog(_index uint256) constant returns(level uint8, time uint256, key string, from address, to address, v1 uint256, v2 uint256, v3 uint256, remarks string)
func (_Votes *VotesCaller) GetLog(opts *bind.CallOpts, _index *big.Int) (struct {
	Level   uint8
	Time    *big.Int
	Key     string
	From    common.Address
	To      common.Address
	V1      *big.Int
	V2      *big.Int
	V3      *big.Int
	Remarks string
}, error) {
	ret := new(struct {
		Level   uint8
		Time    *big.Int
		Key     string
		From    common.Address
		To      common.Address
		V1      *big.Int
		V2      *big.Int
		V3      *big.Int
		Remarks string
	})
	out := ret
	err := _Votes.contract.Call(opts, out, "getLog", _index)
	return *ret, err
}

// GetLog is a free data retrieval call binding the contract method 0x3206b2c6.
//
// Solidity: function getLog(_index uint256) constant returns(level uint8, time uint256, key string, from address, to address, v1 uint256, v2 uint256, v3 uint256, remarks string)
func (_Votes *VotesSession) GetLog(_index *big.Int) (struct {
	Level   uint8
	Time    *big.Int
	Key     string
	From    common.Address
	To      common.Address
	V1      *big.Int
	V2      *big.Int
	V3      *big.Int
	Remarks string
}, error) {
	return _Votes.Contract.GetLog(&_Votes.CallOpts, _index)
}

// GetLog is a free data retrieval call binding the contract method 0x3206b2c6.
//
// Solidity: function getLog(_index uint256) constant returns(level uint8, time uint256, key string, from address, to address, v1 uint256, v2 uint256, v3 uint256, remarks string)
func (_Votes *VotesCallerSession) GetLog(_index *big.Int) (struct {
	Level   uint8
	Time    *big.Int
	Key     string
	From    common.Address
	To      common.Address
	V1      *big.Int
	V2      *big.Int
	V3      *big.Int
	Remarks string
}, error) {
	return _Votes.Contract.GetLog(&_Votes.CallOpts, _index)
}

// GetLogSize is a free data retrieval call binding the contract method 0x0c73a392.
//
// Solidity: function getLogSize() constant returns(size uint256)
func (_Votes *VotesCaller) GetLogSize(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Votes.contract.Call(opts, out, "getLogSize")
	return *ret0, err
}

// GetLogSize is a free data retrieval call binding the contract method 0x0c73a392.
//
// Solidity: function getLogSize() constant returns(size uint256)
func (_Votes *VotesSession) GetLogSize() (*big.Int, error) {
	return _Votes.Contract.GetLogSize(&_Votes.CallOpts)
}

// GetLogSize is a free data retrieval call binding the contract method 0x0c73a392.
//
// Solidity: function getLogSize() constant returns(size uint256)
func (_Votes *VotesCallerSession) GetLogSize() (*big.Int, error) {
	return _Votes.Contract.GetLogSize(&_Votes.CallOpts)
}

// GetVoteRound is a free data retrieval call binding the contract method 0xf0507573.
//
// Solidity: function getVoteRound() constant returns(round uint256)
func (_Votes *VotesCaller) GetVoteRound(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Votes.contract.Call(opts, out, "getVoteRound")
	return *ret0, err
}

// GetVoteRound is a free data retrieval call binding the contract method 0xf0507573.
//
// Solidity: function getVoteRound() constant returns(round uint256)
func (_Votes *VotesSession) GetVoteRound() (*big.Int, error) {
	return _Votes.Contract.GetVoteRound(&_Votes.CallOpts)
}

// GetVoteRound is a free data retrieval call binding the contract method 0xf0507573.
//
// Solidity: function getVoteRound() constant returns(round uint256)
func (_Votes *VotesCallerSession) GetVoteRound() (*big.Int, error) {
	return _Votes.Contract.GetVoteRound(&_Votes.CallOpts)
}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() constant returns(address)
func (_Votes *VotesCaller) Impl(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Votes.contract.Call(opts, out, "impl")
	return *ret0, err
}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() constant returns(address)
func (_Votes *VotesSession) Impl() (common.Address, error) {
	return _Votes.Contract.Impl(&_Votes.CallOpts)
}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() constant returns(address)
func (_Votes *VotesCallerSession) Impl() (common.Address, error) {
	return _Votes.Contract.Impl(&_Votes.CallOpts)
}

// KeyDefault is a free data retrieval call binding the contract method 0x1a08a37c.
//
// Solidity: function keyDefault() constant returns(string)
func (_Votes *VotesCaller) KeyDefault(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Votes.contract.Call(opts, out, "keyDefault")
	return *ret0, err
}

// KeyDefault is a free data retrieval call binding the contract method 0x1a08a37c.
//
// Solidity: function keyDefault() constant returns(string)
func (_Votes *VotesSession) KeyDefault() (string, error) {
	return _Votes.Contract.KeyDefault(&_Votes.CallOpts)
}

// KeyDefault is a free data retrieval call binding the contract method 0x1a08a37c.
//
// Solidity: function keyDefault() constant returns(string)
func (_Votes *VotesCallerSession) KeyDefault() (string, error) {
	return _Votes.Contract.KeyDefault(&_Votes.CallOpts)
}

// LogLevel is a free data retrieval call binding the contract method 0x7567772b.
//
// Solidity: function logLevel() constant returns(uint8)
func (_Votes *VotesCaller) LogLevel(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Votes.contract.Call(opts, out, "logLevel")
	return *ret0, err
}

// LogLevel is a free data retrieval call binding the contract method 0x7567772b.
//
// Solidity: function logLevel() constant returns(uint8)
func (_Votes *VotesSession) LogLevel() (uint8, error) {
	return _Votes.Contract.LogLevel(&_Votes.CallOpts)
}

// LogLevel is a free data retrieval call binding the contract method 0x7567772b.
//
// Solidity: function logLevel() constant returns(uint8)
func (_Votes *VotesCallerSession) LogLevel() (uint8, error) {
	return _Votes.Contract.LogLevel(&_Votes.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_Votes *VotesCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Votes.contract.Call(opts, out, "paused")
	return *ret0, err
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_Votes *VotesSession) Paused() (bool, error) {
	return _Votes.Contract.Paused(&_Votes.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_Votes *VotesCallerSession) Paused() (bool, error) {
	return _Votes.Contract.Paused(&_Votes.CallOpts)
}

// TickVote is a free data retrieval call binding the contract method 0x8765da94.
//
// Solidity: function tickVote() constant returns(bool)
func (_Votes *VotesCaller) TickVote(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Votes.contract.Call(opts, out, "tickVote")
	return *ret0, err
}

// TickVote is a free data retrieval call binding the contract method 0x8765da94.
//
// Solidity: function tickVote() constant returns(bool)
func (_Votes *VotesSession) TickVote() (bool, error) {
	return _Votes.Contract.TickVote(&_Votes.CallOpts)
}

// TickVote is a free data retrieval call binding the contract method 0x8765da94.
//
// Solidity: function tickVote() constant returns(bool)
func (_Votes *VotesCallerSession) TickVote() (bool, error) {
	return _Votes.Contract.TickVote(&_Votes.CallOpts)
}

// CancelAllVotes is a paid mutator transaction binding the contract method 0xf49f9df8.
//
// Solidity: function cancelAllVotes() returns()
func (_Votes *VotesTransactor) CancelAllVotes(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Votes.contract.Transact(opts, "cancelAllVotes")
}

// CancelAllVotes is a paid mutator transaction binding the contract method 0xf49f9df8.
//
// Solidity: function cancelAllVotes() returns()
func (_Votes *VotesSession) CancelAllVotes() (*types.Transaction, error) {
	return _Votes.Contract.CancelAllVotes(&_Votes.TransactOpts)
}

// CancelAllVotes is a paid mutator transaction binding the contract method 0xf49f9df8.
//
// Solidity: function cancelAllVotes() returns()
func (_Votes *VotesTransactorSession) CancelAllVotes() (*types.Transaction, error) {
	return _Votes.Contract.CancelAllVotes(&_Votes.TransactOpts)
}

// ClearLog is a paid mutator transaction binding the contract method 0x5c50745e.
//
// Solidity: function clearLog() returns()
func (_Votes *VotesTransactor) ClearLog(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Votes.contract.Transact(opts, "clearLog")
}

// ClearLog is a paid mutator transaction binding the contract method 0x5c50745e.
//
// Solidity: function clearLog() returns()
func (_Votes *VotesSession) ClearLog() (*types.Transaction, error) {
	return _Votes.Contract.ClearLog(&_Votes.TransactOpts)
}

// ClearLog is a paid mutator transaction binding the contract method 0x5c50745e.
//
// Solidity: function clearLog() returns()
func (_Votes *VotesTransactorSession) ClearLog() (*types.Transaction, error) {
	return _Votes.Contract.ClearLog(&_Votes.TransactOpts)
}

// Enable is a paid mutator transaction binding the contract method 0x2f038fd5.
//
// Solidity: function enable(status uint8) returns()
func (_Votes *VotesTransactor) Enable(opts *bind.TransactOpts, status uint8) (*types.Transaction, error) {
	return _Votes.contract.Transact(opts, "enable", status)
}

// Enable is a paid mutator transaction binding the contract method 0x2f038fd5.
//
// Solidity: function enable(status uint8) returns()
func (_Votes *VotesSession) Enable(status uint8) (*types.Transaction, error) {
	return _Votes.Contract.Enable(&_Votes.TransactOpts, status)
}

// Enable is a paid mutator transaction binding the contract method 0x2f038fd5.
//
// Solidity: function enable(status uint8) returns()
func (_Votes *VotesTransactorSession) Enable(status uint8) (*types.Transaction, error) {
	return _Votes.Contract.Enable(&_Votes.TransactOpts, status)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns(bool)
func (_Votes *VotesTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Votes.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns(bool)
func (_Votes *VotesSession) Pause() (*types.Transaction, error) {
	return _Votes.Contract.Pause(&_Votes.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns(bool)
func (_Votes *VotesTransactorSession) Pause() (*types.Transaction, error) {
	return _Votes.Contract.Pause(&_Votes.TransactOpts)
}

// RegisterCandidate is a paid mutator transaction binding the contract method 0xf7e0079e.
//
// Solidity: function registerCandidate() returns()
func (_Votes *VotesTransactor) RegisterCandidate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Votes.contract.Transact(opts, "registerCandidate")
}

// RegisterCandidate is a paid mutator transaction binding the contract method 0xf7e0079e.
//
// Solidity: function registerCandidate() returns()
func (_Votes *VotesSession) RegisterCandidate() (*types.Transaction, error) {
	return _Votes.Contract.RegisterCandidate(&_Votes.TransactOpts)
}

// RegisterCandidate is a paid mutator transaction binding the contract method 0xf7e0079e.
//
// Solidity: function registerCandidate() returns()
func (_Votes *VotesTransactorSession) RegisterCandidate() (*types.Transaction, error) {
	return _Votes.Contract.RegisterCandidate(&_Votes.TransactOpts)
}

// RotateVote is a paid mutator transaction binding the contract method 0x11070fcd.
//
// Solidity: function rotateVote() returns()
func (_Votes *VotesTransactor) RotateVote(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Votes.contract.Transact(opts, "rotateVote")
}

// RotateVote is a paid mutator transaction binding the contract method 0x11070fcd.
//
// Solidity: function rotateVote() returns()
func (_Votes *VotesSession) RotateVote() (*types.Transaction, error) {
	return _Votes.Contract.RotateVote(&_Votes.TransactOpts)
}

// RotateVote is a paid mutator transaction binding the contract method 0x11070fcd.
//
// Solidity: function rotateVote() returns()
func (_Votes *VotesTransactorSession) RotateVote() (*types.Transaction, error) {
	return _Votes.Contract.RotateVote(&_Votes.TransactOpts)
}

// SetCEO is a paid mutator transaction binding the contract method 0x27d7874c.
//
// Solidity: function setCEO(_newCEO address) returns()
func (_Votes *VotesTransactor) SetCEO(opts *bind.TransactOpts, _newCEO common.Address) (*types.Transaction, error) {
	return _Votes.contract.Transact(opts, "setCEO", _newCEO)
}

// SetCEO is a paid mutator transaction binding the contract method 0x27d7874c.
//
// Solidity: function setCEO(_newCEO address) returns()
func (_Votes *VotesSession) SetCEO(_newCEO common.Address) (*types.Transaction, error) {
	return _Votes.Contract.SetCEO(&_Votes.TransactOpts, _newCEO)
}

// SetCEO is a paid mutator transaction binding the contract method 0x27d7874c.
//
// Solidity: function setCEO(_newCEO address) returns()
func (_Votes *VotesTransactorSession) SetCEO(_newCEO common.Address) (*types.Transaction, error) {
	return _Votes.Contract.SetCEO(&_Votes.TransactOpts, _newCEO)
}

// SetCFO is a paid mutator transaction binding the contract method 0x4e0a3379.
//
// Solidity: function setCFO(_newCFO address) returns()
func (_Votes *VotesTransactor) SetCFO(opts *bind.TransactOpts, _newCFO common.Address) (*types.Transaction, error) {
	return _Votes.contract.Transact(opts, "setCFO", _newCFO)
}

// SetCFO is a paid mutator transaction binding the contract method 0x4e0a3379.
//
// Solidity: function setCFO(_newCFO address) returns()
func (_Votes *VotesSession) SetCFO(_newCFO common.Address) (*types.Transaction, error) {
	return _Votes.Contract.SetCFO(&_Votes.TransactOpts, _newCFO)
}

// SetCFO is a paid mutator transaction binding the contract method 0x4e0a3379.
//
// Solidity: function setCFO(_newCFO address) returns()
func (_Votes *VotesTransactorSession) SetCFO(_newCFO common.Address) (*types.Transaction, error) {
	return _Votes.Contract.SetCFO(&_Votes.TransactOpts, _newCFO)
}

// SetCOO is a paid mutator transaction binding the contract method 0x2ba73c15.
//
// Solidity: function setCOO(_newCOO address) returns()
func (_Votes *VotesTransactor) SetCOO(opts *bind.TransactOpts, _newCOO common.Address) (*types.Transaction, error) {
	return _Votes.contract.Transact(opts, "setCOO", _newCOO)
}

// SetCOO is a paid mutator transaction binding the contract method 0x2ba73c15.
//
// Solidity: function setCOO(_newCOO address) returns()
func (_Votes *VotesSession) SetCOO(_newCOO common.Address) (*types.Transaction, error) {
	return _Votes.Contract.SetCOO(&_Votes.TransactOpts, _newCOO)
}

// SetCOO is a paid mutator transaction binding the contract method 0x2ba73c15.
//
// Solidity: function setCOO(_newCOO address) returns()
func (_Votes *VotesTransactorSession) SetCOO(_newCOO common.Address) (*types.Transaction, error) {
	return _Votes.Contract.SetCOO(&_Votes.TransactOpts, _newCOO)
}

// SetImpl is a paid mutator transaction binding the contract method 0x691bd2ae.
//
// Solidity: function setImpl(implAddress address) returns()
func (_Votes *VotesTransactor) SetImpl(opts *bind.TransactOpts, implAddress common.Address) (*types.Transaction, error) {
	return _Votes.contract.Transact(opts, "setImpl", implAddress)
}

// SetImpl is a paid mutator transaction binding the contract method 0x691bd2ae.
//
// Solidity: function setImpl(implAddress address) returns()
func (_Votes *VotesSession) SetImpl(implAddress common.Address) (*types.Transaction, error) {
	return _Votes.Contract.SetImpl(&_Votes.TransactOpts, implAddress)
}

// SetImpl is a paid mutator transaction binding the contract method 0x691bd2ae.
//
// Solidity: function setImpl(implAddress address) returns()
func (_Votes *VotesTransactorSession) SetImpl(implAddress common.Address) (*types.Transaction, error) {
	return _Votes.Contract.SetImpl(&_Votes.TransactOpts, implAddress)
}

// SetLevel is a paid mutator transaction binding the contract method 0xbd5546be.
//
// Solidity: function setLevel(level uint8) returns()
func (_Votes *VotesTransactor) SetLevel(opts *bind.TransactOpts, level uint8) (*types.Transaction, error) {
	return _Votes.contract.Transact(opts, "setLevel", level)
}

// SetLevel is a paid mutator transaction binding the contract method 0xbd5546be.
//
// Solidity: function setLevel(level uint8) returns()
func (_Votes *VotesSession) SetLevel(level uint8) (*types.Transaction, error) {
	return _Votes.Contract.SetLevel(&_Votes.TransactOpts, level)
}

// SetLevel is a paid mutator transaction binding the contract method 0xbd5546be.
//
// Solidity: function setLevel(level uint8) returns()
func (_Votes *VotesTransactorSession) SetLevel(level uint8) (*types.Transaction, error) {
	return _Votes.Contract.SetLevel(&_Votes.TransactOpts, level)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns(bool)
func (_Votes *VotesTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Votes.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns(bool)
func (_Votes *VotesSession) Unpause() (*types.Transaction, error) {
	return _Votes.Contract.Unpause(&_Votes.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns(bool)
func (_Votes *VotesTransactorSession) Unpause() (*types.Transaction, error) {
	return _Votes.Contract.Unpause(&_Votes.TransactOpts)
}

// VoteCandidate is a paid mutator transaction binding the contract method 0x26bb886d.
//
// Solidity: function voteCandidate(addrCandidate address) returns()
func (_Votes *VotesTransactor) VoteCandidate(opts *bind.TransactOpts, addrCandidate common.Address) (*types.Transaction, error) {
	return _Votes.contract.Transact(opts, "voteCandidate", addrCandidate)
}

// VoteCandidate is a paid mutator transaction binding the contract method 0x26bb886d.
//
// Solidity: function voteCandidate(addrCandidate address) returns()
func (_Votes *VotesSession) VoteCandidate(addrCandidate common.Address) (*types.Transaction, error) {
	return _Votes.Contract.VoteCandidate(&_Votes.TransactOpts, addrCandidate)
}

// VoteCandidate is a paid mutator transaction binding the contract method 0x26bb886d.
//
// Solidity: function voteCandidate(addrCandidate address) returns()
func (_Votes *VotesTransactorSession) VoteCandidate(addrCandidate common.Address) (*types.Transaction, error) {
	return _Votes.Contract.VoteCandidate(&_Votes.TransactOpts, addrCandidate)
}
