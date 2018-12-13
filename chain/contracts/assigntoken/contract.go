// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package assigntoken

import (
	"math/big"
	"strings"

	"github.com/Bokerchain/Boker/chain/accounts/abi"
	"github.com/Bokerchain/Boker/chain/accounts/abi/bind"
	"github.com/Bokerchain/Boker/chain/common"
	"github.com/Bokerchain/Boker/chain/core/types"
)

// AssigntokenABI is the input ABI used to generate the binding from.
const AssigntokenABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"cfoAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ceoAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getLogSize\",\"outputs\":[{\"name\":\"size\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"logKeyDefault\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"assgineTokenPeriodDefault\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newCEO\",\"type\":\"address\"}],\"name\":\"setCEO\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"index\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newCOO\",\"type\":\"address\"}],\"name\":\"setCOO\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"enable\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getLog\",\"outputs\":[{\"name\":\"level\",\"type\":\"uint8\"},{\"name\":\"time\",\"type\":\"uint256\"},{\"name\":\"key\",\"type\":\"string\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"v1\",\"type\":\"uint256\"},{\"name\":\"v2\",\"type\":\"uint256\"},{\"name\":\"v3\",\"type\":\"uint256\"},{\"name\":\"remarks\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"eventType\",\"type\":\"uint256\"},{\"name\":\"addrFrom\",\"type\":\"address\"},{\"name\":\"addrTo\",\"type\":\"address\"},{\"name\":\"eventValue1\",\"type\":\"uint256\"},{\"name\":\"eventValue2\",\"type\":\"uint256\"}],\"name\":\"fireUserEvent\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"checkAssignToken\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"assginedTokensPerPeriodDefault\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newCFO\",\"type\":\"address\"}],\"name\":\"setCFO\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"clearLog\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"implAddress\",\"type\":\"address\"}],\"name\":\"setImpl\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"logLevel\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"impl\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"assignToken\",\"outputs\":[{\"name\":\"\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"cooAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"level\",\"type\":\"uint8\"}],\"name\":\"setLevel\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"logEnabled\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"}]"

// AssigntokenBin is the compiled bytecode used for deploying new contracts.
const AssigntokenBin = `60806040526000600260146101000a81548160ff0219169083151502179055506001600560006101000a81548160ff021916908360ff1602179055506000600560016101000a81548160ff021916908360ff1602179055506040805190810160405280600781526020017f64656661756c740000000000000000000000000000000000000000000000000081525060069080519060200190620000a49291906200013f565b506000600855348015620000b757600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555033600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550620001ee565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200018257805160ff1916838001178555620001b3565b82800160010185558215620001b3579182015b82811115620001b257825182559160200191906001019062000195565b5b509050620001c29190620001c6565b5090565b620001eb91905b80821115620001e7576000816000905550600101620001cd565b5090565b90565b611fd580620001fe6000396000f30060806040526004361061015f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630519ce79146102b45780630a0f81681461030b5780630c73a39214610362578063113feb6d1461038d57806315a59e351461041d57806327d7874c146104485780632986c0e51461048b5780632ba73c15146104b65780632f038fd5146104f95780633206b2c61461052957806333dfffc8146106ca578063378c085c1461074b5780633f4ba83a1461077a57806347fada6e146107a95780634e0a3379146107d45780635c50745e146108175780635c975abb1461082e578063691bd2ae1461085d5780637567772b146108a05780638456cb59146108d15780638abf6077146109005780638da5cb5b14610957578063a237213c146109ae578063b047fb50146109d9578063bd5546be14610a30578063e164f60c14610a60578063f2fde38b14610a91575b600073ffffffffffffffffffffffffffffffffffffffff16600760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16141515156101bd57600080fd5b6102136040805190810160405280600281526020017f2829000000000000000000000000000000000000000000000000000000000000815250336000346000806020604051908101604052806000815250610ad4565b600760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663dbd0353d346040518263ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004016000604051808303818588803b15801561029957600080fd5b505af11580156102ad573d6000803e3d6000fd5b5050505050005b3480156102c057600080fd5b506102c9610aee565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561031757600080fd5b50610320610b14565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561036e57600080fd5b50610377610b39565b6040518082815260200191505060405180910390f35b34801561039957600080fd5b506103a2610b46565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156103e25780820151818401526020810190506103c7565b50505050905090810190601f16801561040f5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561042957600080fd5b50610432610be4565b6040518082815260200191505060405180910390f35b34801561045457600080fd5b50610489600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610bea565b005b34801561049757600080fd5b506104a0610cc4565b6040518082815260200191505060405180910390f35b3480156104c257600080fd5b506104f7600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610cca565b005b34801561050557600080fd5b50610527600480360381019080803560ff169060200190929190505050610da5565b005b34801561053557600080fd5b5061055460048036038101908080359060200190929190505050610e1f565b604051808a60ff1660ff168152602001898152602001806020018873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018681526020018581526020018481526020018060200183810383528a818151815260200191508051906020019080838360005b83811015610620578082015181840152602081019050610605565b50505050905090810190601f16801561064d5780820380516001836020036101000a031916815260200191505b50838103825284818151815260200191508051906020019080838360005b8381101561068657808201518184015260208101905061066b565b50505050905090810190601f1680156106b35780820380516001836020036101000a031916815260200191505b509b50505050505050505050505060405180910390f35b3480156106d657600080fd5b5061074960048036038101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803590602001909291908035906020019092919050505061103c565b005b34801561075757600080fd5b50610760611379565b604051808215151515815260200191505060405180910390f35b34801561078657600080fd5b5061078f611441565b604051808215151515815260200191505060405180910390f35b3480156107b557600080fd5b506107be611507565b6040518082815260200191505060405180910390f35b3480156107e057600080fd5b50610815600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611517565b005b34801561082357600080fd5b5061082c6115f2565b005b34801561083a57600080fd5b5061084361165e565b604051808215151515815260200191505060405180910390f35b34801561086957600080fd5b5061089e600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611671565b005b3480156108ac57600080fd5b506108b56116d8565b604051808260ff1660ff16815260200191505060405180910390f35b3480156108dd57600080fd5b506108e66116eb565b604051808215151515815260200191505060405180910390f35b34801561090c57600080fd5b506109156117b2565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561096357600080fd5b5061096c6117d8565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156109ba57600080fd5b506109c36117fe565b6040518082815260200191505060405180910390f35b3480156109e557600080fd5b506109ee611870565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b348015610a3c57600080fd5b50610a5e600480360381019080803560ff169060200190929190505050611896565b005b348015610a6c57600080fd5b50610a75611910565b604051808260ff1660ff16815260200191505060405180910390f35b348015610a9d57600080fd5b50610ad2600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050611923565b005b610ae560028888888888888861198b565b50505050505050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000600480549050905090565b60068054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610bdc5780601f10610bb157610100808354040283529160200191610bdc565b820191906000526020600020905b815481529060010190602001808311610bbf57829003601f168201915b505050505081565b61012c81565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610c4557600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614151515610c8157600080fd5b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60085481565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610d2557600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614151515610d6157600080fd5b80600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610e0157600080fd5b80600560006101000a81548160ff021916908360ff16021790555050565b60008060606000806000806000606060006004805490508b101515610e4357600080fd5b60008b10151515610e5357600080fd5b60048b815481101515610e6257fe5b906000526020600020906009020190508060010160009054906101000a900460ff169950806002018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610f1e5780601f10610ef357610100808354040283529160200191610f1e565b820191906000526020600020905b815481529060010190602001808311610f0157829003601f168201915b50505050509750806000015498508060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1696508060040160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169550806005015494508060060154935080600701549250806008018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156110275780601f10610ffc57610100808354040283529160200191611027565b820191906000526020600020905b81548152906001019060200180831161100a57829003601f168201915b50505050509150509193959799909294969850565b600260149054906101000a900460ff1615151561105857600080fd5b600073ffffffffffffffffffffffffffffffffffffffff16600760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16141515156110b657600080fd5b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16148061115e57506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b806111b65750600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b15156111c157600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16141515156111fd57600080fd5b6112516040805190810160405280600d81526020017f66697265557365724576656e740000000000000000000000000000000000000081525085858886866020604051908101604052806000815250611be1565b600760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166333dfffc886868686866040518663ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808681526020018573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200182815260200195505050505050600060405180830381600087803b15801561135a57600080fd5b505af115801561136e573d6000803e3d6000fd5b505050505050505050565b6000600760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663378c085c6040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401602060405180830381600087803b15801561140157600080fd5b505af1158015611415573d6000803e3d6000fd5b505050506040513d602081101561142b57600080fd5b8101908080519060200190929190505050905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561149e57600080fd5b600260149054906101000a900460ff1615156114b957600080fd5b6000600260146101000a81548160ff0219169083151502179055507f7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b3360405160405180910390a16001905090565b670de0b6b3a76400006103de0281565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561157257600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16141515156115ae57600080fd5b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561164e57600080fd5b6004600061165c9190611d5b565b565b600260149054906101000a900460ff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415156116cc57600080fd5b6116d581611bfb565b50565b600560019054906101000a900460ff1681565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561174857600080fd5b600260149054906101000a900460ff1615151561176457600080fd5b6001600260146101000a81548160ff0219169083151502179055507f6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff62560405160405180910390a16001905090565b600760009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60006118566040805190810160405280600b81526020017f61737369676e546f6b656e000000000000000000000000000000000000000000815250336000806000806020604051908101604052806000815250611c45565b600860008154809291906001019190505550600854905090565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415156118f257600080fd5b80600560016101000a81548160ff021916908360ff16021790555050565b600560009054906101000a900460ff1681565b600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561197f57600080fd5b61198881611c5f565b50565b611993611d7f565b600560009054906101000a900460ff1660ff16600014156119b357611bd6565b8860ff16600560019054906101000a900460ff1660ff1611156119d557611bd6565b88816020019060ff16908160ff16815250508781604001819052504281600001818152505086816060019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505085816080019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050848160a0018181525050838160c0018181525050828160e001818152505081816101000181905250600481908060018154018082558091505090600182039060005260206000209060090201600090919290919091506000820151816000015560208201518160010160006101000a81548160ff021916908360ff1602179055506040820151816002019080519060200190611b07929190611dfb565b5060608201518160030160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060808201518160040160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060a0820151816005015560c0820151816006015560e08201518160070155610100820151816008019080519060200190611bd1929190611dfb565b505050505b505050505050505050565b611bf260008888888888888861198b565b50505050505050565b600081905080600760006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050565b611c5660018888888888888861198b565b50505050505050565b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614151515611c9b57600080fd5b8073ffffffffffffffffffffffffffffffffffffffff16600360009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a380600360006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b5080546000825560090290600052602060002090810190611d7c9190611e7b565b50565b6101206040519081016040528060008152602001600060ff16815260200160608152602001600073ffffffffffffffffffffffffffffffffffffffff168152602001600073ffffffffffffffffffffffffffffffffffffffff168152602001600081526020016000815260200160008152602001606081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10611e3c57805160ff1916838001178555611e6a565b82800160010185558215611e6a579182015b82811115611e69578251825591602001919060010190611e4e565b5b509050611e779190611f3c565b5090565b611f3991905b80821115611f35576000808201600090556001820160006101000a81549060ff0219169055600282016000611eb69190611f61565b6003820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690556004820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055600582016000905560068201600090556007820160009055600882016000611f2c9190611f61565b50600901611e81565b5090565b90565b611f5e91905b80821115611f5a576000816000905550600101611f42565b5090565b90565b50805460018160011615610100020316600290046000825580601f10611f875750611fa6565b601f016020900490600052602060002090810190611fa59190611f3c565b5b505600a165627a7a7230582042b02d1682ce6327c3ccf6582183f9eecf050ea00dc1f61cd72b7e7bb6b0090d0029`

// DeployAssigntoken deploys a new Ethereum contract, binding an instance of Assigntoken to it.
func DeployAssigntoken(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Assigntoken, error) {
	parsed, err := abi.JSON(strings.NewReader(AssigntokenABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AssigntokenBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Assigntoken{AssigntokenCaller: AssigntokenCaller{contract: contract}, AssigntokenTransactor: AssigntokenTransactor{contract: contract}}, nil
}

// Assigntoken is an auto generated Go binding around an Ethereum contract.
type Assigntoken struct {
	AssigntokenCaller     // Read-only binding to the contract
	AssigntokenTransactor // Write-only binding to the contract
}

// AssigntokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type AssigntokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AssigntokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AssigntokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AssigntokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AssigntokenSession struct {
	Contract     *Assigntoken      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AssigntokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AssigntokenCallerSession struct {
	Contract *AssigntokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// AssigntokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AssigntokenTransactorSession struct {
	Contract     *AssigntokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// AssigntokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type AssigntokenRaw struct {
	Contract *Assigntoken // Generic contract binding to access the raw methods on
}

// AssigntokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AssigntokenCallerRaw struct {
	Contract *AssigntokenCaller // Generic read-only contract binding to access the raw methods on
}

// AssigntokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AssigntokenTransactorRaw struct {
	Contract *AssigntokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAssigntoken creates a new instance of Assigntoken, bound to a specific deployed contract.
func NewAssigntoken(address common.Address, backend bind.ContractBackend) (*Assigntoken, error) {
	contract, err := bindAssigntoken(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Assigntoken{AssigntokenCaller: AssigntokenCaller{contract: contract}, AssigntokenTransactor: AssigntokenTransactor{contract: contract}}, nil
}

// NewAssigntokenCaller creates a new read-only instance of Assigntoken, bound to a specific deployed contract.
func NewAssigntokenCaller(address common.Address, caller bind.ContractCaller) (*AssigntokenCaller, error) {
	contract, err := bindAssigntoken(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &AssigntokenCaller{contract: contract}, nil
}

// NewAssigntokenTransactor creates a new write-only instance of Assigntoken, bound to a specific deployed contract.
func NewAssigntokenTransactor(address common.Address, transactor bind.ContractTransactor) (*AssigntokenTransactor, error) {
	contract, err := bindAssigntoken(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &AssigntokenTransactor{contract: contract}, nil
}

// bindAssigntoken binds a generic wrapper to an already deployed contract.
func bindAssigntoken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AssigntokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Assigntoken *AssigntokenRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Assigntoken.Contract.AssigntokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Assigntoken *AssigntokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Assigntoken.Contract.AssigntokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Assigntoken *AssigntokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Assigntoken.Contract.AssigntokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Assigntoken *AssigntokenCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Assigntoken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Assigntoken *AssigntokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Assigntoken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Assigntoken *AssigntokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Assigntoken.Contract.contract.Transact(opts, method, params...)
}

// AssgineTokenPeriodDefault is a free data retrieval call binding the contract method 0x15a59e35.
//
// Solidity: function assgineTokenPeriodDefault() constant returns(uint256)
func (_Assigntoken *AssigntokenCaller) AssgineTokenPeriodDefault(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assigntoken.contract.Call(opts, out, "assgineTokenPeriodDefault")
	return *ret0, err
}

// AssgineTokenPeriodDefault is a free data retrieval call binding the contract method 0x15a59e35.
//
// Solidity: function assgineTokenPeriodDefault() constant returns(uint256)
func (_Assigntoken *AssigntokenSession) AssgineTokenPeriodDefault() (*big.Int, error) {
	return _Assigntoken.Contract.AssgineTokenPeriodDefault(&_Assigntoken.CallOpts)
}

// AssgineTokenPeriodDefault is a free data retrieval call binding the contract method 0x15a59e35.
//
// Solidity: function assgineTokenPeriodDefault() constant returns(uint256)
func (_Assigntoken *AssigntokenCallerSession) AssgineTokenPeriodDefault() (*big.Int, error) {
	return _Assigntoken.Contract.AssgineTokenPeriodDefault(&_Assigntoken.CallOpts)
}

// AssginedTokensPerPeriodDefault is a free data retrieval call binding the contract method 0x47fada6e.
//
// Solidity: function assginedTokensPerPeriodDefault() constant returns(uint256)
func (_Assigntoken *AssigntokenCaller) AssginedTokensPerPeriodDefault(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assigntoken.contract.Call(opts, out, "assginedTokensPerPeriodDefault")
	return *ret0, err
}

// AssginedTokensPerPeriodDefault is a free data retrieval call binding the contract method 0x47fada6e.
//
// Solidity: function assginedTokensPerPeriodDefault() constant returns(uint256)
func (_Assigntoken *AssigntokenSession) AssginedTokensPerPeriodDefault() (*big.Int, error) {
	return _Assigntoken.Contract.AssginedTokensPerPeriodDefault(&_Assigntoken.CallOpts)
}

// AssginedTokensPerPeriodDefault is a free data retrieval call binding the contract method 0x47fada6e.
//
// Solidity: function assginedTokensPerPeriodDefault() constant returns(uint256)
func (_Assigntoken *AssigntokenCallerSession) AssginedTokensPerPeriodDefault() (*big.Int, error) {
	return _Assigntoken.Contract.AssginedTokensPerPeriodDefault(&_Assigntoken.CallOpts)
}

// CeoAddress is a free data retrieval call binding the contract method 0x0a0f8168.
//
// Solidity: function ceoAddress() constant returns(address)
func (_Assigntoken *AssigntokenCaller) CeoAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Assigntoken.contract.Call(opts, out, "ceoAddress")
	return *ret0, err
}

// CeoAddress is a free data retrieval call binding the contract method 0x0a0f8168.
//
// Solidity: function ceoAddress() constant returns(address)
func (_Assigntoken *AssigntokenSession) CeoAddress() (common.Address, error) {
	return _Assigntoken.Contract.CeoAddress(&_Assigntoken.CallOpts)
}

// CeoAddress is a free data retrieval call binding the contract method 0x0a0f8168.
//
// Solidity: function ceoAddress() constant returns(address)
func (_Assigntoken *AssigntokenCallerSession) CeoAddress() (common.Address, error) {
	return _Assigntoken.Contract.CeoAddress(&_Assigntoken.CallOpts)
}

// CfoAddress is a free data retrieval call binding the contract method 0x0519ce79.
//
// Solidity: function cfoAddress() constant returns(address)
func (_Assigntoken *AssigntokenCaller) CfoAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Assigntoken.contract.Call(opts, out, "cfoAddress")
	return *ret0, err
}

// CfoAddress is a free data retrieval call binding the contract method 0x0519ce79.
//
// Solidity: function cfoAddress() constant returns(address)
func (_Assigntoken *AssigntokenSession) CfoAddress() (common.Address, error) {
	return _Assigntoken.Contract.CfoAddress(&_Assigntoken.CallOpts)
}

// CfoAddress is a free data retrieval call binding the contract method 0x0519ce79.
//
// Solidity: function cfoAddress() constant returns(address)
func (_Assigntoken *AssigntokenCallerSession) CfoAddress() (common.Address, error) {
	return _Assigntoken.Contract.CfoAddress(&_Assigntoken.CallOpts)
}

// CheckAssignToken is a free data retrieval call binding the contract method 0x378c085c.
//
// Solidity: function checkAssignToken() constant returns(bool)
func (_Assigntoken *AssigntokenCaller) CheckAssignToken(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Assigntoken.contract.Call(opts, out, "checkAssignToken")
	return *ret0, err
}

// CheckAssignToken is a free data retrieval call binding the contract method 0x378c085c.
//
// Solidity: function checkAssignToken() constant returns(bool)
func (_Assigntoken *AssigntokenSession) CheckAssignToken() (bool, error) {
	return _Assigntoken.Contract.CheckAssignToken(&_Assigntoken.CallOpts)
}

// CheckAssignToken is a free data retrieval call binding the contract method 0x378c085c.
//
// Solidity: function checkAssignToken() constant returns(bool)
func (_Assigntoken *AssigntokenCallerSession) CheckAssignToken() (bool, error) {
	return _Assigntoken.Contract.CheckAssignToken(&_Assigntoken.CallOpts)
}

// CooAddress is a free data retrieval call binding the contract method 0xb047fb50.
//
// Solidity: function cooAddress() constant returns(address)
func (_Assigntoken *AssigntokenCaller) CooAddress(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Assigntoken.contract.Call(opts, out, "cooAddress")
	return *ret0, err
}

// CooAddress is a free data retrieval call binding the contract method 0xb047fb50.
//
// Solidity: function cooAddress() constant returns(address)
func (_Assigntoken *AssigntokenSession) CooAddress() (common.Address, error) {
	return _Assigntoken.Contract.CooAddress(&_Assigntoken.CallOpts)
}

// CooAddress is a free data retrieval call binding the contract method 0xb047fb50.
//
// Solidity: function cooAddress() constant returns(address)
func (_Assigntoken *AssigntokenCallerSession) CooAddress() (common.Address, error) {
	return _Assigntoken.Contract.CooAddress(&_Assigntoken.CallOpts)
}

// GetLog is a free data retrieval call binding the contract method 0x3206b2c6.
//
// Solidity: function getLog(_index uint256) constant returns(level uint8, time uint256, key string, from address, to address, v1 uint256, v2 uint256, v3 uint256, remarks string)
func (_Assigntoken *AssigntokenCaller) GetLog(opts *bind.CallOpts, _index *big.Int) (struct {
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
	err := _Assigntoken.contract.Call(opts, out, "getLog", _index)
	return *ret, err
}

// GetLog is a free data retrieval call binding the contract method 0x3206b2c6.
//
// Solidity: function getLog(_index uint256) constant returns(level uint8, time uint256, key string, from address, to address, v1 uint256, v2 uint256, v3 uint256, remarks string)
func (_Assigntoken *AssigntokenSession) GetLog(_index *big.Int) (struct {
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
	return _Assigntoken.Contract.GetLog(&_Assigntoken.CallOpts, _index)
}

// GetLog is a free data retrieval call binding the contract method 0x3206b2c6.
//
// Solidity: function getLog(_index uint256) constant returns(level uint8, time uint256, key string, from address, to address, v1 uint256, v2 uint256, v3 uint256, remarks string)
func (_Assigntoken *AssigntokenCallerSession) GetLog(_index *big.Int) (struct {
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
	return _Assigntoken.Contract.GetLog(&_Assigntoken.CallOpts, _index)
}

// GetLogSize is a free data retrieval call binding the contract method 0x0c73a392.
//
// Solidity: function getLogSize() constant returns(size uint256)
func (_Assigntoken *AssigntokenCaller) GetLogSize(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assigntoken.contract.Call(opts, out, "getLogSize")
	return *ret0, err
}

// GetLogSize is a free data retrieval call binding the contract method 0x0c73a392.
//
// Solidity: function getLogSize() constant returns(size uint256)
func (_Assigntoken *AssigntokenSession) GetLogSize() (*big.Int, error) {
	return _Assigntoken.Contract.GetLogSize(&_Assigntoken.CallOpts)
}

// GetLogSize is a free data retrieval call binding the contract method 0x0c73a392.
//
// Solidity: function getLogSize() constant returns(size uint256)
func (_Assigntoken *AssigntokenCallerSession) GetLogSize() (*big.Int, error) {
	return _Assigntoken.Contract.GetLogSize(&_Assigntoken.CallOpts)
}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() constant returns(address)
func (_Assigntoken *AssigntokenCaller) Impl(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Assigntoken.contract.Call(opts, out, "impl")
	return *ret0, err
}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() constant returns(address)
func (_Assigntoken *AssigntokenSession) Impl() (common.Address, error) {
	return _Assigntoken.Contract.Impl(&_Assigntoken.CallOpts)
}

// Impl is a free data retrieval call binding the contract method 0x8abf6077.
//
// Solidity: function impl() constant returns(address)
func (_Assigntoken *AssigntokenCallerSession) Impl() (common.Address, error) {
	return _Assigntoken.Contract.Impl(&_Assigntoken.CallOpts)
}

// Index is a free data retrieval call binding the contract method 0x2986c0e5.
//
// Solidity: function index() constant returns(int256)
func (_Assigntoken *AssigntokenCaller) Index(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Assigntoken.contract.Call(opts, out, "index")
	return *ret0, err
}

// Index is a free data retrieval call binding the contract method 0x2986c0e5.
//
// Solidity: function index() constant returns(int256)
func (_Assigntoken *AssigntokenSession) Index() (*big.Int, error) {
	return _Assigntoken.Contract.Index(&_Assigntoken.CallOpts)
}

// Index is a free data retrieval call binding the contract method 0x2986c0e5.
//
// Solidity: function index() constant returns(int256)
func (_Assigntoken *AssigntokenCallerSession) Index() (*big.Int, error) {
	return _Assigntoken.Contract.Index(&_Assigntoken.CallOpts)
}

// LogEnabled is a free data retrieval call binding the contract method 0xe164f60c.
//
// Solidity: function logEnabled() constant returns(uint8)
func (_Assigntoken *AssigntokenCaller) LogEnabled(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assigntoken.contract.Call(opts, out, "logEnabled")
	return *ret0, err
}

// LogEnabled is a free data retrieval call binding the contract method 0xe164f60c.
//
// Solidity: function logEnabled() constant returns(uint8)
func (_Assigntoken *AssigntokenSession) LogEnabled() (uint8, error) {
	return _Assigntoken.Contract.LogEnabled(&_Assigntoken.CallOpts)
}

// LogEnabled is a free data retrieval call binding the contract method 0xe164f60c.
//
// Solidity: function logEnabled() constant returns(uint8)
func (_Assigntoken *AssigntokenCallerSession) LogEnabled() (uint8, error) {
	return _Assigntoken.Contract.LogEnabled(&_Assigntoken.CallOpts)
}

// LogKeyDefault is a free data retrieval call binding the contract method 0x113feb6d.
//
// Solidity: function logKeyDefault() constant returns(string)
func (_Assigntoken *AssigntokenCaller) LogKeyDefault(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Assigntoken.contract.Call(opts, out, "logKeyDefault")
	return *ret0, err
}

// LogKeyDefault is a free data retrieval call binding the contract method 0x113feb6d.
//
// Solidity: function logKeyDefault() constant returns(string)
func (_Assigntoken *AssigntokenSession) LogKeyDefault() (string, error) {
	return _Assigntoken.Contract.LogKeyDefault(&_Assigntoken.CallOpts)
}

// LogKeyDefault is a free data retrieval call binding the contract method 0x113feb6d.
//
// Solidity: function logKeyDefault() constant returns(string)
func (_Assigntoken *AssigntokenCallerSession) LogKeyDefault() (string, error) {
	return _Assigntoken.Contract.LogKeyDefault(&_Assigntoken.CallOpts)
}

// LogLevel is a free data retrieval call binding the contract method 0x7567772b.
//
// Solidity: function logLevel() constant returns(uint8)
func (_Assigntoken *AssigntokenCaller) LogLevel(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assigntoken.contract.Call(opts, out, "logLevel")
	return *ret0, err
}

// LogLevel is a free data retrieval call binding the contract method 0x7567772b.
//
// Solidity: function logLevel() constant returns(uint8)
func (_Assigntoken *AssigntokenSession) LogLevel() (uint8, error) {
	return _Assigntoken.Contract.LogLevel(&_Assigntoken.CallOpts)
}

// LogLevel is a free data retrieval call binding the contract method 0x7567772b.
//
// Solidity: function logLevel() constant returns(uint8)
func (_Assigntoken *AssigntokenCallerSession) LogLevel() (uint8, error) {
	return _Assigntoken.Contract.LogLevel(&_Assigntoken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Assigntoken *AssigntokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Assigntoken.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Assigntoken *AssigntokenSession) Owner() (common.Address, error) {
	return _Assigntoken.Contract.Owner(&_Assigntoken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Assigntoken *AssigntokenCallerSession) Owner() (common.Address, error) {
	return _Assigntoken.Contract.Owner(&_Assigntoken.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_Assigntoken *AssigntokenCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Assigntoken.contract.Call(opts, out, "paused")
	return *ret0, err
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_Assigntoken *AssigntokenSession) Paused() (bool, error) {
	return _Assigntoken.Contract.Paused(&_Assigntoken.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_Assigntoken *AssigntokenCallerSession) Paused() (bool, error) {
	return _Assigntoken.Contract.Paused(&_Assigntoken.CallOpts)
}

// AssignToken is a paid mutator transaction binding the contract method 0xa237213c.
//
// Solidity: function assignToken() returns(int256)
func (_Assigntoken *AssigntokenTransactor) AssignToken(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Assigntoken.contract.Transact(opts, "assignToken")
}

// AssignToken is a paid mutator transaction binding the contract method 0xa237213c.
//
// Solidity: function assignToken() returns(int256)
func (_Assigntoken *AssigntokenSession) AssignToken() (*types.Transaction, error) {
	return _Assigntoken.Contract.AssignToken(&_Assigntoken.TransactOpts)
}

// AssignToken is a paid mutator transaction binding the contract method 0xa237213c.
//
// Solidity: function assignToken() returns(int256)
func (_Assigntoken *AssigntokenTransactorSession) AssignToken() (*types.Transaction, error) {
	return _Assigntoken.Contract.AssignToken(&_Assigntoken.TransactOpts)
}

// ClearLog is a paid mutator transaction binding the contract method 0x5c50745e.
//
// Solidity: function clearLog() returns()
func (_Assigntoken *AssigntokenTransactor) ClearLog(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Assigntoken.contract.Transact(opts, "clearLog")
}

// ClearLog is a paid mutator transaction binding the contract method 0x5c50745e.
//
// Solidity: function clearLog() returns()
func (_Assigntoken *AssigntokenSession) ClearLog() (*types.Transaction, error) {
	return _Assigntoken.Contract.ClearLog(&_Assigntoken.TransactOpts)
}

// ClearLog is a paid mutator transaction binding the contract method 0x5c50745e.
//
// Solidity: function clearLog() returns()
func (_Assigntoken *AssigntokenTransactorSession) ClearLog() (*types.Transaction, error) {
	return _Assigntoken.Contract.ClearLog(&_Assigntoken.TransactOpts)
}

// Enable is a paid mutator transaction binding the contract method 0x2f038fd5.
//
// Solidity: function enable(status uint8) returns()
func (_Assigntoken *AssigntokenTransactor) Enable(opts *bind.TransactOpts, status uint8) (*types.Transaction, error) {
	return _Assigntoken.contract.Transact(opts, "enable", status)
}

// Enable is a paid mutator transaction binding the contract method 0x2f038fd5.
//
// Solidity: function enable(status uint8) returns()
func (_Assigntoken *AssigntokenSession) Enable(status uint8) (*types.Transaction, error) {
	return _Assigntoken.Contract.Enable(&_Assigntoken.TransactOpts, status)
}

// Enable is a paid mutator transaction binding the contract method 0x2f038fd5.
//
// Solidity: function enable(status uint8) returns()
func (_Assigntoken *AssigntokenTransactorSession) Enable(status uint8) (*types.Transaction, error) {
	return _Assigntoken.Contract.Enable(&_Assigntoken.TransactOpts, status)
}

// FireUserEvent is a paid mutator transaction binding the contract method 0x33dfffc8.
//
// Solidity: function fireUserEvent(eventType uint256, addrFrom address, addrTo address, eventValue1 uint256, eventValue2 uint256) returns()
func (_Assigntoken *AssigntokenTransactor) FireUserEvent(opts *bind.TransactOpts, eventType *big.Int, addrFrom common.Address, addrTo common.Address, eventValue1 *big.Int, eventValue2 *big.Int) (*types.Transaction, error) {
	return _Assigntoken.contract.Transact(opts, "fireUserEvent", eventType, addrFrom, addrTo, eventValue1, eventValue2)
}

// FireUserEvent is a paid mutator transaction binding the contract method 0x33dfffc8.
//
// Solidity: function fireUserEvent(eventType uint256, addrFrom address, addrTo address, eventValue1 uint256, eventValue2 uint256) returns()
func (_Assigntoken *AssigntokenSession) FireUserEvent(eventType *big.Int, addrFrom common.Address, addrTo common.Address, eventValue1 *big.Int, eventValue2 *big.Int) (*types.Transaction, error) {
	return _Assigntoken.Contract.FireUserEvent(&_Assigntoken.TransactOpts, eventType, addrFrom, addrTo, eventValue1, eventValue2)
}

// FireUserEvent is a paid mutator transaction binding the contract method 0x33dfffc8.
//
// Solidity: function fireUserEvent(eventType uint256, addrFrom address, addrTo address, eventValue1 uint256, eventValue2 uint256) returns()
func (_Assigntoken *AssigntokenTransactorSession) FireUserEvent(eventType *big.Int, addrFrom common.Address, addrTo common.Address, eventValue1 *big.Int, eventValue2 *big.Int) (*types.Transaction, error) {
	return _Assigntoken.Contract.FireUserEvent(&_Assigntoken.TransactOpts, eventType, addrFrom, addrTo, eventValue1, eventValue2)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns(bool)
func (_Assigntoken *AssigntokenTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Assigntoken.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns(bool)
func (_Assigntoken *AssigntokenSession) Pause() (*types.Transaction, error) {
	return _Assigntoken.Contract.Pause(&_Assigntoken.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns(bool)
func (_Assigntoken *AssigntokenTransactorSession) Pause() (*types.Transaction, error) {
	return _Assigntoken.Contract.Pause(&_Assigntoken.TransactOpts)
}

// SetCEO is a paid mutator transaction binding the contract method 0x27d7874c.
//
// Solidity: function setCEO(_newCEO address) returns()
func (_Assigntoken *AssigntokenTransactor) SetCEO(opts *bind.TransactOpts, _newCEO common.Address) (*types.Transaction, error) {
	return _Assigntoken.contract.Transact(opts, "setCEO", _newCEO)
}

// SetCEO is a paid mutator transaction binding the contract method 0x27d7874c.
//
// Solidity: function setCEO(_newCEO address) returns()
func (_Assigntoken *AssigntokenSession) SetCEO(_newCEO common.Address) (*types.Transaction, error) {
	return _Assigntoken.Contract.SetCEO(&_Assigntoken.TransactOpts, _newCEO)
}

// SetCEO is a paid mutator transaction binding the contract method 0x27d7874c.
//
// Solidity: function setCEO(_newCEO address) returns()
func (_Assigntoken *AssigntokenTransactorSession) SetCEO(_newCEO common.Address) (*types.Transaction, error) {
	return _Assigntoken.Contract.SetCEO(&_Assigntoken.TransactOpts, _newCEO)
}

// SetCFO is a paid mutator transaction binding the contract method 0x4e0a3379.
//
// Solidity: function setCFO(_newCFO address) returns()
func (_Assigntoken *AssigntokenTransactor) SetCFO(opts *bind.TransactOpts, _newCFO common.Address) (*types.Transaction, error) {
	return _Assigntoken.contract.Transact(opts, "setCFO", _newCFO)
}

// SetCFO is a paid mutator transaction binding the contract method 0x4e0a3379.
//
// Solidity: function setCFO(_newCFO address) returns()
func (_Assigntoken *AssigntokenSession) SetCFO(_newCFO common.Address) (*types.Transaction, error) {
	return _Assigntoken.Contract.SetCFO(&_Assigntoken.TransactOpts, _newCFO)
}

// SetCFO is a paid mutator transaction binding the contract method 0x4e0a3379.
//
// Solidity: function setCFO(_newCFO address) returns()
func (_Assigntoken *AssigntokenTransactorSession) SetCFO(_newCFO common.Address) (*types.Transaction, error) {
	return _Assigntoken.Contract.SetCFO(&_Assigntoken.TransactOpts, _newCFO)
}

// SetCOO is a paid mutator transaction binding the contract method 0x2ba73c15.
//
// Solidity: function setCOO(_newCOO address) returns()
func (_Assigntoken *AssigntokenTransactor) SetCOO(opts *bind.TransactOpts, _newCOO common.Address) (*types.Transaction, error) {
	return _Assigntoken.contract.Transact(opts, "setCOO", _newCOO)
}

// SetCOO is a paid mutator transaction binding the contract method 0x2ba73c15.
//
// Solidity: function setCOO(_newCOO address) returns()
func (_Assigntoken *AssigntokenSession) SetCOO(_newCOO common.Address) (*types.Transaction, error) {
	return _Assigntoken.Contract.SetCOO(&_Assigntoken.TransactOpts, _newCOO)
}

// SetCOO is a paid mutator transaction binding the contract method 0x2ba73c15.
//
// Solidity: function setCOO(_newCOO address) returns()
func (_Assigntoken *AssigntokenTransactorSession) SetCOO(_newCOO common.Address) (*types.Transaction, error) {
	return _Assigntoken.Contract.SetCOO(&_Assigntoken.TransactOpts, _newCOO)
}

// SetImpl is a paid mutator transaction binding the contract method 0x691bd2ae.
//
// Solidity: function setImpl(implAddress address) returns()
func (_Assigntoken *AssigntokenTransactor) SetImpl(opts *bind.TransactOpts, implAddress common.Address) (*types.Transaction, error) {
	return _Assigntoken.contract.Transact(opts, "setImpl", implAddress)
}

// SetImpl is a paid mutator transaction binding the contract method 0x691bd2ae.
//
// Solidity: function setImpl(implAddress address) returns()
func (_Assigntoken *AssigntokenSession) SetImpl(implAddress common.Address) (*types.Transaction, error) {
	return _Assigntoken.Contract.SetImpl(&_Assigntoken.TransactOpts, implAddress)
}

// SetImpl is a paid mutator transaction binding the contract method 0x691bd2ae.
//
// Solidity: function setImpl(implAddress address) returns()
func (_Assigntoken *AssigntokenTransactorSession) SetImpl(implAddress common.Address) (*types.Transaction, error) {
	return _Assigntoken.Contract.SetImpl(&_Assigntoken.TransactOpts, implAddress)
}

// SetLevel is a paid mutator transaction binding the contract method 0xbd5546be.
//
// Solidity: function setLevel(level uint8) returns()
func (_Assigntoken *AssigntokenTransactor) SetLevel(opts *bind.TransactOpts, level uint8) (*types.Transaction, error) {
	return _Assigntoken.contract.Transact(opts, "setLevel", level)
}

// SetLevel is a paid mutator transaction binding the contract method 0xbd5546be.
//
// Solidity: function setLevel(level uint8) returns()
func (_Assigntoken *AssigntokenSession) SetLevel(level uint8) (*types.Transaction, error) {
	return _Assigntoken.Contract.SetLevel(&_Assigntoken.TransactOpts, level)
}

// SetLevel is a paid mutator transaction binding the contract method 0xbd5546be.
//
// Solidity: function setLevel(level uint8) returns()
func (_Assigntoken *AssigntokenTransactorSession) SetLevel(level uint8) (*types.Transaction, error) {
	return _Assigntoken.Contract.SetLevel(&_Assigntoken.TransactOpts, level)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Assigntoken *AssigntokenTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _Assigntoken.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Assigntoken *AssigntokenSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _Assigntoken.Contract.TransferOwnership(&_Assigntoken.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Assigntoken *AssigntokenTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _Assigntoken.Contract.TransferOwnership(&_Assigntoken.TransactOpts, _newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns(bool)
func (_Assigntoken *AssigntokenTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Assigntoken.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns(bool)
func (_Assigntoken *AssigntokenSession) Unpause() (*types.Transaction, error) {
	return _Assigntoken.Contract.Unpause(&_Assigntoken.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns(bool)
func (_Assigntoken *AssigntokenTransactorSession) Unpause() (*types.Transaction, error) {
	return _Assigntoken.Contract.Unpause(&_Assigntoken.TransactOpts)
}
