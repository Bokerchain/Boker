// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package boker_contract

import (
	"math/big"
	"strings"

	"github.com/Bokerchain/Boker/chain/accounts/abi"
	"github.com/Bokerchain/Boker/chain/accounts/abi/bind"
	"github.com/Bokerchain/Boker/chain/common"
	"github.com/Bokerchain/Boker/chain/core/types"
)

// BokerInterfaceABI is the input ABI used to generate the binding from.
const BokerInterfaceABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"getCandidates\",\"outputs\":[{\"name\":\"addresses\",\"type\":\"address[]\"},{\"name\":\"tickets\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"rotateVote\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"cName\",\"type\":\"string\"}],\"name\":\"contractAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"key\",\"type\":\"string\"}],\"name\":\"globalConfigString\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"checkAssignToken\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"version\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"createTime\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"bokerManager\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"tickVote\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"assignToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"addrManager\",\"type\":\"address\"}],\"name\":\"setManager\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"key\",\"type\":\"string\"}],\"name\":\"globalConfigInt\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getVoteRound\",\"outputs\":[{\"name\":\"round\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"addrManager\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

// BokerInterfaceBin is the compiled bytecode used for deploying new contracts.
const BokerInterfaceBin = `60806040526040805190810160405280600581526020017f312e302e3100000000000000000000000000000000000000000000000000000081525060019080519060200190620000519291906200024b565b50426002553480156200006357600080fd5b5060405160208062001f0b8339810180604052810190808051906020019092919050505080336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550620000e281620000ea640100000000026401000000009004565b5050620002fa565b60008190508073ffffffffffffffffffffffffffffffffffffffff1663519c28826040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401602060405180830381600087803b1580156200015457600080fd5b505af115801562000169573d6000803e3d6000fd5b505050506040513d60208110156200018057600080fd5b8101908080519060200190929190505050151562000206576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260118152602001807f6e6f7420426f6b65724d616e616765722100000000000000000000000000000081525060200191505060405180910390fd5b80600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200028e57805160ff1916838001178555620002bf565b82800160010185558215620002bf579182015b82811115620002be578251825591602001919060010190620002a1565b5b509050620002ce9190620002d2565b5090565b620002f791905b80821115620002f3576000816000905550600101620002d9565b5090565b90565b611c01806200030a6000396000f3006080604052600436106100db576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806306a49fce146101e257806311070fcd146102965780631eb726af146102ad578063378298bc14610356578063378c085c1461043857806354fd4d501461046757806361dcd7ab146104f757806366ebc1c6146105225780638765da94146105795780638da5cb5b146105a8578063a237213c146105ff578063d0ebdbe714610616578063d43c802114610659578063f0507573146106d6578063f2fde38b14610701575b6101196040805190810160405280600781526020017f46696e616e636500000000000000000000000000000000000000000000000000815250610744565b73ffffffffffffffffffffffffffffffffffffffff1663f9fdc7c034336002600781111561014357fe5b6040518463ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828152602001925050506000604051808303818588803b1580156101c757600080fd5b505af11580156101db573d6000803e3d6000fd5b5050505050005b3480156101ee57600080fd5b506101f76109f2565b604051808060200180602001838103835285818151815260200191508051906020019060200280838360005b8381101561023e578082015181840152602081019050610223565b50505050905001838103825284818151815260200191508051906020019060200280838360005b83811015610280578082015181840152602081019050610265565b5050505090500194505050505060405180910390f35b3480156102a257600080fd5b506102ab610b7e565b005b3480156102b957600080fd5b50610314600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610744565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561036257600080fd5b506103bd600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610ee6565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156103fd5780820151818401526020810190506103e2565b50505050905090810190601f16801561042a5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561044457600080fd5b5061044d611076565b604051808215151515815260200191505060405180910390f35b34801561047357600080fd5b5061047c611159565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156104bc5780820151818401526020810190506104a1565b50505050905090810190601f1680156104e95780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561050357600080fd5b5061050c6111f7565b6040518082815260200191505060405180910390f35b34801561052e57600080fd5b506105376111fd565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561058557600080fd5b5061058e611223565b604051808215151515815260200191505060405180910390f35b3480156105b457600080fd5b506105bd611306565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561060b57600080fd5b5061061461132b565b005b34801561062257600080fd5b50610657600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611693565b005b34801561066557600080fd5b506106c0600480360381019080803590602001908201803590602001908080601f01602080910402602001604051908101604052809392919081815260200183838082843782019150505050505091929192905050506116fa565b6040518082815260200191505060405180910390f35b3480156106e257600080fd5b506106eb611834565b6040518082815260200191505060405180910390f35b34801561070d57600080fd5b50610742600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611917565b005b600080600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663fca1f3c1846040518263ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004018080602001828103825283818151815260200191508051906020019080838360005b838110156107f15780820151818401526020810190506107d6565b50505050905090810190601f16801561081e5780820380516001836020036101000a031916815260200191505b5092505050600060405180830381600087803b15801561083d57600080fd5b505af1158015610851573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f8201168201806040525060a081101561087b57600080fd5b8101908080519060200190929190805164010000000081111561089d57600080fd5b828101905060208101848111156108b357600080fd5b81518560018202830111640100000000821117156108d057600080fd5b50509291906020018051906020019092919080516401000000008111156108f657600080fd5b8281019050602081018481111561090c57600080fd5b815185600182028301116401000000008211171561092957600080fd5b505092919060200180519060200190929190505050505092505050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141515156109e9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252600a8152602001807f616464722069732030210000000000000000000000000000000000000000000081525060200191505060405180910390fd5b80915050919050565b606080610a336040805190810160405280600481526020017f4e6f646500000000000000000000000000000000000000000000000000000000815250610744565b73ffffffffffffffffffffffffffffffffffffffff166306a49fce6040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401600060405180830381600087803b158015610a9657600080fd5b505af1158015610aaa573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f820116820180604052506040811015610ad457600080fd5b810190808051640100000000811115610aec57600080fd5b82810190506020810184811115610b0257600080fd5b8151856020820283011164010000000082111715610b1f57600080fd5b50509291906020018051640100000000811115610b3b57600080fd5b82810190506020810184811115610b5157600080fd5b8151856020820283011164010000000082111715610b6e57600080fd5b5050929190505050915091509091565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16635c975abb6040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401602060405180830381600087803b158015610c0457600080fd5b505af1158015610c18573d6000803e3d6000fd5b505050506040513d6020811015610c2e57600080fd5b8101908080519060200190929190505050151515610cb4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260078152602001807f706175736564210000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16634d4810836040805190810160405280600581526020017f61646d696e000000000000000000000000000000000000000000000000000000815250336040518363ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828103825284818151815260200191508051906020019080838360005b83811015610dc6578082015181840152602081019050610dab565b50505050905090810190601f168015610df35780820380516001836020036101000a031916815260200191505b509350505050600060405180830381600087803b158015610e1357600080fd5b505af1158015610e27573d6000803e3d6000fd5b50505050610e696040805190810160405280600481526020017f4e6f646500000000000000000000000000000000000000000000000000000000815250610744565b73ffffffffffffffffffffffffffffffffffffffff166311070fcd6040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401600060405180830381600087803b158015610ecc57600080fd5b505af1158015610ee0573d6000803e3d6000fd5b50505050565b6060600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16634a189f35836040518263ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004018080602001828103825283818151815260200191508051906020019080838360005b83811015610f92578082015181840152602081019050610f77565b50505050905090810190601f168015610fbf5780820380516001836020036101000a031916815260200191505b5092505050600060405180830381600087803b158015610fde57600080fd5b505af1158015610ff2573d6000803e3d6000fd5b505050506040513d6000823e3d601f19601f82011682018060405250602081101561101c57600080fd5b81019080805164010000000081111561103457600080fd5b8281019050602081018481111561104a57600080fd5b815185600182028301116401000000008211171561106757600080fd5b50509291905050509050919050565b60006110b66040805190810160405280600a81526020017f546f6b656e506f77657200000000000000000000000000000000000000000000815250610744565b73ffffffffffffffffffffffffffffffffffffffff1663378c085c6040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401602060405180830381600087803b15801561111957600080fd5b505af115801561112d573d6000803e3d6000fd5b505050506040513d602081101561114357600080fd5b8101908080519060200190929190505050905090565b60018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156111ef5780601f106111c4576101008083540402835291602001916111ef565b820191906000526020600020905b8154815290600101906020018083116111d257829003601f168201915b505050505081565b60025481565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60006112636040805190810160405280600481526020017f4e6f646500000000000000000000000000000000000000000000000000000000815250610744565b73ffffffffffffffffffffffffffffffffffffffff1663b0417e986040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401602060405180830381600087803b1580156112c657600080fd5b505af11580156112da573d6000803e3d6000fd5b505050506040513d60208110156112f057600080fd5b8101908080519060200190929190505050905090565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16635c975abb6040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401602060405180830381600087803b1580156113b157600080fd5b505af11580156113c5573d6000803e3d6000fd5b505050506040513d60208110156113db57600080fd5b8101908080519060200190929190505050151515611461576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260078152602001807f706175736564210000000000000000000000000000000000000000000000000081525060200191505060405180910390fd5b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16634d4810836040805190810160405280600581526020017f61646d696e000000000000000000000000000000000000000000000000000000815250336040518363ffffffff167c010000000000000000000000000000000000000000000000000000000002815260040180806020018373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001828103825284818151815260200191508051906020019080838360005b83811015611573578082015181840152602081019050611558565b50505050905090810190601f1680156115a05780820380516001836020036101000a031916815260200191505b509350505050600060405180830381600087803b1580156115c057600080fd5b505af11580156115d4573d6000803e3d6000fd5b505050506116166040805190810160405280600a81526020017f546f6b656e506f77657200000000000000000000000000000000000000000000815250610744565b73ffffffffffffffffffffffffffffffffffffffff1663a237213c6040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401600060405180830381600087803b15801561167957600080fd5b505af115801561168d573d6000803e3d6000fd5b50505050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415156116ee57600080fd5b6116f78161197e565b50565b6000600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16633230b078836040518263ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004018080602001828103825283818151815260200191508051906020019080838360005b838110156117a657808201518184015260208101905061178b565b50505050905090810190601f1680156117d35780820380516001836020036101000a031916815260200191505b5092505050602060405180830381600087803b1580156117f257600080fd5b505af1158015611806573d6000803e3d6000fd5b505050506040513d602081101561181c57600080fd5b81019080805190602001909291905050509050919050565b60006118746040805190810160405280600481526020017f4e6f646500000000000000000000000000000000000000000000000000000000815250610744565b73ffffffffffffffffffffffffffffffffffffffff1663f05075736040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401602060405180830381600087803b1580156118d757600080fd5b505af11580156118eb573d6000803e3d6000fd5b505050506040513d602081101561190157600080fd5b8101908080519060200190929190505050905090565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561197257600080fd5b61197b81611adb565b50565b60008190508073ffffffffffffffffffffffffffffffffffffffff1663519c28826040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401602060405180830381600087803b1580156119e757600080fd5b505af11580156119fb573d6000803e3d6000fd5b505050506040513d6020811015611a1157600080fd5b81019080805190602001909291905050501515611a96576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260118152602001807f6e6f7420426f6b65724d616e616765722100000000000000000000000000000081525060200191505060405180910390fd5b80600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614151515611b1757600080fd5b8073ffffffffffffffffffffffffffffffffffffffff166000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a3806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505600a165627a7a72305820f06ae67bf7dceb5b1fdebf0040d860a9a545bb8dcdbc9465ac70be6f14e8e9780029`

// DeployBokerInterface deploys a new Ethereum contract, binding an instance of BokerInterface to it.
func DeployBokerInterface(auth *bind.TransactOpts, backend bind.ContractBackend, addrManager common.Address) (common.Address, *types.Transaction, *BokerInterface, error) {
	parsed, err := abi.JSON(strings.NewReader(BokerInterfaceABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(BokerInterfaceBin), backend, addrManager)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BokerInterface{BokerInterfaceCaller: BokerInterfaceCaller{contract: contract}, BokerInterfaceTransactor: BokerInterfaceTransactor{contract: contract}}, nil
}

// BokerInterface is an auto generated Go binding around an Ethereum contract.
type BokerInterface struct {
	BokerInterfaceCaller     // Read-only binding to the contract
	BokerInterfaceTransactor // Write-only binding to the contract
}

// BokerInterfaceCaller is an auto generated read-only Go binding around an Ethereum contract.
type BokerInterfaceCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BokerInterfaceTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BokerInterfaceTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BokerInterfaceSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BokerInterfaceSession struct {
	Contract     *BokerInterface   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BokerInterfaceCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BokerInterfaceCallerSession struct {
	Contract *BokerInterfaceCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// BokerInterfaceTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BokerInterfaceTransactorSession struct {
	Contract     *BokerInterfaceTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// BokerInterfaceRaw is an auto generated low-level Go binding around an Ethereum contract.
type BokerInterfaceRaw struct {
	Contract *BokerInterface // Generic contract binding to access the raw methods on
}

// BokerInterfaceCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BokerInterfaceCallerRaw struct {
	Contract *BokerInterfaceCaller // Generic read-only contract binding to access the raw methods on
}

// BokerInterfaceTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BokerInterfaceTransactorRaw struct {
	Contract *BokerInterfaceTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBokerInterface creates a new instance of BokerInterface, bound to a specific deployed contract.
func NewBokerInterface(address common.Address, backend bind.ContractBackend) (*BokerInterface, error) {
	contract, err := bindBokerInterface(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BokerInterface{BokerInterfaceCaller: BokerInterfaceCaller{contract: contract}, BokerInterfaceTransactor: BokerInterfaceTransactor{contract: contract}}, nil
}

// NewBokerInterfaceCaller creates a new read-only instance of BokerInterface, bound to a specific deployed contract.
func NewBokerInterfaceCaller(address common.Address, caller bind.ContractCaller) (*BokerInterfaceCaller, error) {
	contract, err := bindBokerInterface(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &BokerInterfaceCaller{contract: contract}, nil
}

// NewBokerInterfaceTransactor creates a new write-only instance of BokerInterface, bound to a specific deployed contract.
func NewBokerInterfaceTransactor(address common.Address, transactor bind.ContractTransactor) (*BokerInterfaceTransactor, error) {
	contract, err := bindBokerInterface(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &BokerInterfaceTransactor{contract: contract}, nil
}

// bindBokerInterface binds a generic wrapper to an already deployed contract.
func bindBokerInterface(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(BokerInterfaceABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BokerInterface *BokerInterfaceRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BokerInterface.Contract.BokerInterfaceCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BokerInterface *BokerInterfaceRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BokerInterface.Contract.BokerInterfaceTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BokerInterface *BokerInterfaceRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BokerInterface.Contract.BokerInterfaceTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BokerInterface *BokerInterfaceCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _BokerInterface.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BokerInterface *BokerInterfaceTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BokerInterface.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BokerInterface *BokerInterfaceTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BokerInterface.Contract.contract.Transact(opts, method, params...)
}

// BokerManager is a free data retrieval call binding the contract method 0x66ebc1c6.
//
// Solidity: function bokerManager() constant returns(address)
func (_BokerInterface *BokerInterfaceCaller) BokerManager(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _BokerInterface.contract.Call(opts, out, "bokerManager")
	return *ret0, err
}

// BokerManager is a free data retrieval call binding the contract method 0x66ebc1c6.
//
// Solidity: function bokerManager() constant returns(address)
func (_BokerInterface *BokerInterfaceSession) BokerManager() (common.Address, error) {
	return _BokerInterface.Contract.BokerManager(&_BokerInterface.CallOpts)
}

// BokerManager is a free data retrieval call binding the contract method 0x66ebc1c6.
//
// Solidity: function bokerManager() constant returns(address)
func (_BokerInterface *BokerInterfaceCallerSession) BokerManager() (common.Address, error) {
	return _BokerInterface.Contract.BokerManager(&_BokerInterface.CallOpts)
}

// CheckAssignToken is a free data retrieval call binding the contract method 0x378c085c.
//
// Solidity: function checkAssignToken() constant returns(bool)
func (_BokerInterface *BokerInterfaceCaller) CheckAssignToken(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _BokerInterface.contract.Call(opts, out, "checkAssignToken")
	return *ret0, err
}

// CheckAssignToken is a free data retrieval call binding the contract method 0x378c085c.
//
// Solidity: function checkAssignToken() constant returns(bool)
func (_BokerInterface *BokerInterfaceSession) CheckAssignToken() (bool, error) {
	return _BokerInterface.Contract.CheckAssignToken(&_BokerInterface.CallOpts)
}

// CheckAssignToken is a free data retrieval call binding the contract method 0x378c085c.
//
// Solidity: function checkAssignToken() constant returns(bool)
func (_BokerInterface *BokerInterfaceCallerSession) CheckAssignToken() (bool, error) {
	return _BokerInterface.Contract.CheckAssignToken(&_BokerInterface.CallOpts)
}

// ContractAddress is a free data retrieval call binding the contract method 0x1eb726af.
//
// Solidity: function contractAddress(cName string) constant returns(address)
func (_BokerInterface *BokerInterfaceCaller) ContractAddress(opts *bind.CallOpts, cName string) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _BokerInterface.contract.Call(opts, out, "contractAddress", cName)
	return *ret0, err
}

// ContractAddress is a free data retrieval call binding the contract method 0x1eb726af.
//
// Solidity: function contractAddress(cName string) constant returns(address)
func (_BokerInterface *BokerInterfaceSession) ContractAddress(cName string) (common.Address, error) {
	return _BokerInterface.Contract.ContractAddress(&_BokerInterface.CallOpts, cName)
}

// ContractAddress is a free data retrieval call binding the contract method 0x1eb726af.
//
// Solidity: function contractAddress(cName string) constant returns(address)
func (_BokerInterface *BokerInterfaceCallerSession) ContractAddress(cName string) (common.Address, error) {
	return _BokerInterface.Contract.ContractAddress(&_BokerInterface.CallOpts, cName)
}

// CreateTime is a free data retrieval call binding the contract method 0x61dcd7ab.
//
// Solidity: function createTime() constant returns(uint256)
func (_BokerInterface *BokerInterfaceCaller) CreateTime(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BokerInterface.contract.Call(opts, out, "createTime")
	return *ret0, err
}

// CreateTime is a free data retrieval call binding the contract method 0x61dcd7ab.
//
// Solidity: function createTime() constant returns(uint256)
func (_BokerInterface *BokerInterfaceSession) CreateTime() (*big.Int, error) {
	return _BokerInterface.Contract.CreateTime(&_BokerInterface.CallOpts)
}

// CreateTime is a free data retrieval call binding the contract method 0x61dcd7ab.
//
// Solidity: function createTime() constant returns(uint256)
func (_BokerInterface *BokerInterfaceCallerSession) CreateTime() (*big.Int, error) {
	return _BokerInterface.Contract.CreateTime(&_BokerInterface.CallOpts)
}

// GetCandidates is a free data retrieval call binding the contract method 0x06a49fce.
//
// Solidity: function getCandidates() constant returns(addresses address[], tickets uint256[])
func (_BokerInterface *BokerInterfaceCaller) GetCandidates(opts *bind.CallOpts) (struct {
	Addresses []common.Address
	Tickets   []*big.Int
}, error) {
	ret := new(struct {
		Addresses []common.Address
		Tickets   []*big.Int
	})
	out := ret
	err := _BokerInterface.contract.Call(opts, out, "getCandidates")
	return *ret, err
}

// GetCandidates is a free data retrieval call binding the contract method 0x06a49fce.
//
// Solidity: function getCandidates() constant returns(addresses address[], tickets uint256[])
func (_BokerInterface *BokerInterfaceSession) GetCandidates() (struct {
	Addresses []common.Address
	Tickets   []*big.Int
}, error) {
	return _BokerInterface.Contract.GetCandidates(&_BokerInterface.CallOpts)
}

// GetCandidates is a free data retrieval call binding the contract method 0x06a49fce.
//
// Solidity: function getCandidates() constant returns(addresses address[], tickets uint256[])
func (_BokerInterface *BokerInterfaceCallerSession) GetCandidates() (struct {
	Addresses []common.Address
	Tickets   []*big.Int
}, error) {
	return _BokerInterface.Contract.GetCandidates(&_BokerInterface.CallOpts)
}

// GetVoteRound is a free data retrieval call binding the contract method 0xf0507573.
//
// Solidity: function getVoteRound() constant returns(round uint256)
func (_BokerInterface *BokerInterfaceCaller) GetVoteRound(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BokerInterface.contract.Call(opts, out, "getVoteRound")
	return *ret0, err
}

// GetVoteRound is a free data retrieval call binding the contract method 0xf0507573.
//
// Solidity: function getVoteRound() constant returns(round uint256)
func (_BokerInterface *BokerInterfaceSession) GetVoteRound() (*big.Int, error) {
	return _BokerInterface.Contract.GetVoteRound(&_BokerInterface.CallOpts)
}

// GetVoteRound is a free data retrieval call binding the contract method 0xf0507573.
//
// Solidity: function getVoteRound() constant returns(round uint256)
func (_BokerInterface *BokerInterfaceCallerSession) GetVoteRound() (*big.Int, error) {
	return _BokerInterface.Contract.GetVoteRound(&_BokerInterface.CallOpts)
}

// GlobalConfigInt is a free data retrieval call binding the contract method 0xd43c8021.
//
// Solidity: function globalConfigInt(key string) constant returns(uint256)
func (_BokerInterface *BokerInterfaceCaller) GlobalConfigInt(opts *bind.CallOpts, key string) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _BokerInterface.contract.Call(opts, out, "globalConfigInt", key)
	return *ret0, err
}

// GlobalConfigInt is a free data retrieval call binding the contract method 0xd43c8021.
//
// Solidity: function globalConfigInt(key string) constant returns(uint256)
func (_BokerInterface *BokerInterfaceSession) GlobalConfigInt(key string) (*big.Int, error) {
	return _BokerInterface.Contract.GlobalConfigInt(&_BokerInterface.CallOpts, key)
}

// GlobalConfigInt is a free data retrieval call binding the contract method 0xd43c8021.
//
// Solidity: function globalConfigInt(key string) constant returns(uint256)
func (_BokerInterface *BokerInterfaceCallerSession) GlobalConfigInt(key string) (*big.Int, error) {
	return _BokerInterface.Contract.GlobalConfigInt(&_BokerInterface.CallOpts, key)
}

// GlobalConfigString is a free data retrieval call binding the contract method 0x378298bc.
//
// Solidity: function globalConfigString(key string) constant returns(string)
func (_BokerInterface *BokerInterfaceCaller) GlobalConfigString(opts *bind.CallOpts, key string) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _BokerInterface.contract.Call(opts, out, "globalConfigString", key)
	return *ret0, err
}

// GlobalConfigString is a free data retrieval call binding the contract method 0x378298bc.
//
// Solidity: function globalConfigString(key string) constant returns(string)
func (_BokerInterface *BokerInterfaceSession) GlobalConfigString(key string) (string, error) {
	return _BokerInterface.Contract.GlobalConfigString(&_BokerInterface.CallOpts, key)
}

// GlobalConfigString is a free data retrieval call binding the contract method 0x378298bc.
//
// Solidity: function globalConfigString(key string) constant returns(string)
func (_BokerInterface *BokerInterfaceCallerSession) GlobalConfigString(key string) (string, error) {
	return _BokerInterface.Contract.GlobalConfigString(&_BokerInterface.CallOpts, key)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_BokerInterface *BokerInterfaceCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _BokerInterface.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_BokerInterface *BokerInterfaceSession) Owner() (common.Address, error) {
	return _BokerInterface.Contract.Owner(&_BokerInterface.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_BokerInterface *BokerInterfaceCallerSession) Owner() (common.Address, error) {
	return _BokerInterface.Contract.Owner(&_BokerInterface.CallOpts)
}

// TickVote is a free data retrieval call binding the contract method 0x8765da94.
//
// Solidity: function tickVote() constant returns(bool)
func (_BokerInterface *BokerInterfaceCaller) TickVote(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _BokerInterface.contract.Call(opts, out, "tickVote")
	return *ret0, err
}

// TickVote is a free data retrieval call binding the contract method 0x8765da94.
//
// Solidity: function tickVote() constant returns(bool)
func (_BokerInterface *BokerInterfaceSession) TickVote() (bool, error) {
	return _BokerInterface.Contract.TickVote(&_BokerInterface.CallOpts)
}

// TickVote is a free data retrieval call binding the contract method 0x8765da94.
//
// Solidity: function tickVote() constant returns(bool)
func (_BokerInterface *BokerInterfaceCallerSession) TickVote() (bool, error) {
	return _BokerInterface.Contract.TickVote(&_BokerInterface.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_BokerInterface *BokerInterfaceCaller) Version(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _BokerInterface.contract.Call(opts, out, "version")
	return *ret0, err
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_BokerInterface *BokerInterfaceSession) Version() (string, error) {
	return _BokerInterface.Contract.Version(&_BokerInterface.CallOpts)
}

// Version is a free data retrieval call binding the contract method 0x54fd4d50.
//
// Solidity: function version() constant returns(string)
func (_BokerInterface *BokerInterfaceCallerSession) Version() (string, error) {
	return _BokerInterface.Contract.Version(&_BokerInterface.CallOpts)
}

// AssignToken is a paid mutator transaction binding the contract method 0xa237213c.
//
// Solidity: function assignToken() returns()
func (_BokerInterface *BokerInterfaceTransactor) AssignToken(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BokerInterface.contract.Transact(opts, "assignToken")
}

// AssignToken is a paid mutator transaction binding the contract method 0xa237213c.
//
// Solidity: function assignToken() returns()
func (_BokerInterface *BokerInterfaceSession) AssignToken() (*types.Transaction, error) {
	return _BokerInterface.Contract.AssignToken(&_BokerInterface.TransactOpts)
}

// AssignToken is a paid mutator transaction binding the contract method 0xa237213c.
//
// Solidity: function assignToken() returns()
func (_BokerInterface *BokerInterfaceTransactorSession) AssignToken() (*types.Transaction, error) {
	return _BokerInterface.Contract.AssignToken(&_BokerInterface.TransactOpts)
}

// RotateVote is a paid mutator transaction binding the contract method 0x11070fcd.
//
// Solidity: function rotateVote() returns()
func (_BokerInterface *BokerInterfaceTransactor) RotateVote(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BokerInterface.contract.Transact(opts, "rotateVote")
}

// RotateVote is a paid mutator transaction binding the contract method 0x11070fcd.
//
// Solidity: function rotateVote() returns()
func (_BokerInterface *BokerInterfaceSession) RotateVote() (*types.Transaction, error) {
	return _BokerInterface.Contract.RotateVote(&_BokerInterface.TransactOpts)
}

// RotateVote is a paid mutator transaction binding the contract method 0x11070fcd.
//
// Solidity: function rotateVote() returns()
func (_BokerInterface *BokerInterfaceTransactorSession) RotateVote() (*types.Transaction, error) {
	return _BokerInterface.Contract.RotateVote(&_BokerInterface.TransactOpts)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(addrManager address) returns()
func (_BokerInterface *BokerInterfaceTransactor) SetManager(opts *bind.TransactOpts, addrManager common.Address) (*types.Transaction, error) {
	return _BokerInterface.contract.Transact(opts, "setManager", addrManager)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(addrManager address) returns()
func (_BokerInterface *BokerInterfaceSession) SetManager(addrManager common.Address) (*types.Transaction, error) {
	return _BokerInterface.Contract.SetManager(&_BokerInterface.TransactOpts, addrManager)
}

// SetManager is a paid mutator transaction binding the contract method 0xd0ebdbe7.
//
// Solidity: function setManager(addrManager address) returns()
func (_BokerInterface *BokerInterfaceTransactorSession) SetManager(addrManager common.Address) (*types.Transaction, error) {
	return _BokerInterface.Contract.SetManager(&_BokerInterface.TransactOpts, addrManager)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_BokerInterface *BokerInterfaceTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _BokerInterface.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_BokerInterface *BokerInterfaceSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _BokerInterface.Contract.TransferOwnership(&_BokerInterface.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_BokerInterface *BokerInterfaceTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _BokerInterface.Contract.TransferOwnership(&_BokerInterface.TransactOpts, _newOwner)
}
