// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package assigntoken

import (
	"math/big"
	"strings"

	"github.com/boker/go-ethereum/accounts/abi"
	"github.com/boker/go-ethereum/accounts/abi/bind"
	"github.com/boker/go-ethereum/common"
	"github.com/boker/go-ethereum/core/types"
)

// AssigntokenABI is the input ABI used to generate the binding from.
const AssigntokenABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"cfoAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"ceoAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getLogSize\",\"outputs\":[{\"name\":\"size\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"assgineTokenPeriodDefault\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"keyDefault\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"enabled\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newCEO\",\"type\":\"address\"}],\"name\":\"setCEO\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newCOO\",\"type\":\"address\"}],\"name\":\"setCOO\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"enable\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"getLog\",\"outputs\":[{\"name\":\"level\",\"type\":\"uint8\"},{\"name\":\"time\",\"type\":\"uint256\"},{\"name\":\"key\",\"type\":\"string\"},{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"v1\",\"type\":\"uint256\"},{\"name\":\"v2\",\"type\":\"uint256\"},{\"name\":\"v3\",\"type\":\"uint256\"},{\"name\":\"remarks\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"eventType\",\"type\":\"uint256\"},{\"name\":\"addrFrom\",\"type\":\"address\"},{\"name\":\"addrTo\",\"type\":\"address\"},{\"name\":\"eventValue1\",\"type\":\"uint256\"},{\"name\":\"eventValue2\",\"type\":\"uint256\"}],\"name\":\"fireUserEvent\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"checkAssignToken\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"assginedTokensPerPeriodDefault\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newCFO\",\"type\":\"address\"}],\"name\":\"setCFO\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"clearLog\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"implAddress\",\"type\":\"address\"}],\"name\":\"setImpl\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"logLevel\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"impl\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"assignToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"cooAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"level\",\"type\":\"uint8\"}],\"name\":\"setLevel\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"}]"

// AssigntokenBin is the compiled bytecode used for deploying new contracts.
const AssigntokenBin = `60806040526000600260146101000a81548160ff0219169083151502179055506001600460006101000a81548160ff021916908360ff1602179055506000600460016101000a81548160ff021916908360ff1602179055506040805190810160405280600781526020017f64656661756c740000000000000000000000000000000000000000000000000081525060059080519060200190620000a4929190620000f9565b50348015620000b257600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550620001a8565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200013c57805160ff19168380011785556200016d565b828001600101855582156200016d579182015b828111156200016c5782518255916020019190600101906200014f565b5b5090506200017c919062000180565b5090565b620001a591905b80821115620001a157600081600090555060010162000187565b5090565b90565b61210080620001b86000396000f30060806040526004361061013e576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680630519ce79146102935780630a0f8168146102ea5780630c73a3921461034157806315a59e351461036c5780631a08a37c14610397578063238dafe01461042757806327d7874c146104585780632ba73c151461049b5780632f038fd5146104de5780633206b2c61461050e57806333dfffc8146106af578063378c085c146107305780633f4ba83a1461075f57806347fada6e1461078e5780634e0a3379146107b95780635c50745e146107fc5780635c975abb14610813578063691bd2ae146108425780637567772b146108855780638456cb59146108b65780638abf6077146108e5578063a237213c1461093c578063b047fb5014610953578063bd5546be146109aa575b600073ffffffffffffffffffffffffffffffffffffffff16600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415151561019c57600080fd5b6101f26040805190810160405280600281526020017f28290000000000000000000000000000000000000000000000000000000000008152503360003460008060206040519081016040528060008152506109da565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663dbd0353d346040518263ffffffff167c01000000000000000000000000000000000000000000000000000000000281526004016000604051808303818588803b15801561027857600080fd5b505af115801561028c573d6000803e3d6000fd5b5050505050005b34801561029f57600080fd5b506102a86109f4565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156102f657600080fd5b506102ff610a1a565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561034d57600080fd5b50610356610a3f565b6040518082815260200191505060405180910390f35b34801561037857600080fd5b50610381610a4c565b6040518082815260200191505060405180910390f35b3480156103a357600080fd5b506103ac610a52565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156103ec5780820151818401526020810190506103d1565b50505050905090810190601f1680156104195780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561043357600080fd5b5061043c610af0565b604051808260ff1660ff16815260200191505060405180910390f35b34801561046457600080fd5b50610499600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610b03565b005b3480156104a757600080fd5b506104dc600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610bdd565b005b3480156104ea57600080fd5b5061050c600480360381019080803560ff169060200190929190505050610cb8565b005b34801561051a57600080fd5b5061053960048036038101908080359060200190929190505050610de1565b604051808a60ff1660ff168152602001898152602001806020018873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018681526020018581526020018481526020018060200183810383528a818151815260200191508051906020019080838360005b838110156106055780820151818401526020810190506105ea565b50505050905090810190601f1680156106325780820380516001836020036101000a031916815260200191505b50838103825284818151815260200191508051906020019080838360005b8381101561066b578082015181840152602081019050610650565b50505050905090810190601f1680156106985780820380516001836020036101000a031916815260200191505b509b50505050505050505050505060405180910390f35b3480156106bb57600080fd5b5061072e60048036038101908080359060200190929190803573ffffffffffffffffffffffffffffffffffffffff169060200190929190803573ffffffffffffffffffffffffffffffffffffffff1690602001909291908035906020019092919080359060200190929190505050610ffe565b005b34801561073c57600080fd5b5061074561133b565b604051808215151515815260200191505060405180910390f35b34801561076b57600080fd5b50610774611403565b604051808215151515815260200191505060405180910390f35b34801561079a57600080fd5b506107a36114c9565b6040518082815260200191505060405180910390f35b3480156107c557600080fd5b506107fa600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506114d9565b005b34801561080857600080fd5b506108116115b4565b005b34801561081f57600080fd5b506108286116cf565b604051808215151515815260200191505060405180910390f35b34801561084e57600080fd5b50610883600480360381019080803573ffffffffffffffffffffffffffffffffffffffff1690602001909291905050506116e2565b005b34801561089157600080fd5b5061089a611749565b604051808260ff1660ff16815260200191505060405180910390f35b3480156108c257600080fd5b506108cb61175c565b604051808215151515815260200191505060405180910390f35b3480156108f157600080fd5b506108fa611823565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561094857600080fd5b50610951611849565b005b34801561095f57600080fd5b506109686119b9565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156109b657600080fd5b506109d8600480360381019080803560ff1690602001909291905050506119df565b005b6109eb600288888888888888611b08565b50505050505050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000600380549050905090565b61012c81565b60058054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610ae85780601f10610abd57610100808354040283529160200191610ae8565b820191906000526020600020905b815481529060010190602001808311610acb57829003601f168201915b505050505081565b600460009054906101000a900460ff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610b5e57600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614151515610b9a57600080fd5b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141515610c3857600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614151515610c7457600080fd5b80600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161480610d6057506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b80610db85750600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b1515610dc357600080fd5b80600460006101000a81548160ff021916908360ff16021790555050565b60008060606000806000806000606060006003805490508b101515610e0557600080fd5b60008b10151515610e1557600080fd5b60038b815481101515610e2457fe5b906000526020600020906009020190508060010160009054906101000a900460ff169950806002018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610ee05780601f10610eb557610100808354040283529160200191610ee0565b820191906000526020600020905b815481529060010190602001808311610ec357829003601f168201915b50505050509750806000015498508060030160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1696508060040160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169550806005015494508060060154935080600701549250806008018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610fe95780601f10610fbe57610100808354040283529160200191610fe9565b820191906000526020600020905b815481529060010190602001808311610fcc57829003601f168201915b50505050509150509193959799909294969850565b600260149054906101000a900460ff1615151561101a57600080fd5b600073ffffffffffffffffffffffffffffffffffffffff16600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161415151561107857600080fd5b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16148061112057506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b806111785750600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b151561118357600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16141515156111bf57600080fd5b6112136040805190810160405280600d81526020017f66697265557365724576656e740000000000000000000000000000000000000081525085858886866020604051908101604052806000815250611d5e565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166333dfffc886868686866040518663ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401808681526020018573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200182815260200195505050505050600060405180830381600087803b15801561131c57600080fd5b505af1158015611330573d6000803e3d6000fd5b505050505050505050565b6000600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663378c085c6040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401602060405180830381600087803b1580156113c357600080fd5b505af11580156113d7573d6000803e3d6000fd5b505050506040513d60208110156113ed57600080fd5b8101908080519060200190929190505050905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561146057600080fd5b600260149054906101000a900460ff16151561147b57600080fd5b6000600260146101000a81548160ff0219169083151502179055507f7805862f689e2f13df9f062ff482ad3ad112aca9e0847911ed832e158c525b3360405160405180910390a16001905090565b670de0b6b3a76400006103de0281565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561153457600080fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415151561157057600080fd5b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16148061165c57506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b806116b45750600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b15156116bf57600080fd5b600360006116cd9190611e86565b565b600260149054906101000a900460ff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614151561173d57600080fd5b61174681611d78565b50565b600460019054906101000a900460ff1681565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415156117b957600080fd5b600260149054906101000a900460ff161515156117d557600080fd5b6001600260146101000a81548160ff0219169083151502179055507f6985a02210a168e66602d3235cb6db0e70f92b3ba4d376a33c0f3d9434bff62560405160405180910390a16001905090565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600260149054906101000a900460ff1615151561186557600080fd5b600073ffffffffffffffffffffffffffffffffffffffff16600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16141515156118c357600080fd5b6119196040805190810160405280600b81526020017f61737369676e546f6b656e000000000000000000000000000000000000000000815250336000806000806020604051908101604052806000815250611e6c565b600660009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663a237213c6040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401600060405180830381600087803b15801561199f57600080fd5b505af11580156119b3573d6000803e3d6000fd5b50505050565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600260009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161480611a8757506000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b80611adf5750600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16145b1515611aea57600080fd5b80600460016101000a81548160ff021916908360ff16021790555050565b611b10611eaa565b600460009054906101000a900460ff1660ff1660001415611b3057611d53565b8860ff16600460019054906101000a900460ff1660ff161115611b5257611d53565b88816020019060ff16908160ff16815250508781604001819052504281600001818152505086816060019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff168152505085816080019073ffffffffffffffffffffffffffffffffffffffff16908173ffffffffffffffffffffffffffffffffffffffff1681525050848160a0018181525050838160c0018181525050828160e001818152505081816101000181905250600381908060018154018082558091505090600182039060005260206000209060090201600090919290919091506000820151816000015560208201518160010160006101000a81548160ff021916908360ff1602179055506040820151816002019080519060200190611c84929190611f26565b5060608201518160030160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060808201518160040160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060a0820151816005015560c0820151816006015560e08201518160070155610100820151816008019080519060200190611d4e929190611f26565b505050505b505050505050505050565b611d6f600088888888888888611b08565b50505050505050565b60008190508073ffffffffffffffffffffffffffffffffffffffff1663d62983b66040518163ffffffff167c0100000000000000000000000000000000000000000000000000000000028152600401602060405180830381600087803b158015611de157600080fd5b505af1158015611df5573d6000803e3d6000fd5b505050506040513d6020811015611e0b57600080fd5b81019080805190602001909291905050501515611e2757600080fd5b80600660006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505050565b611e7d600188888888888888611b08565b50505050505050565b5080546000825560090290600052602060002090810190611ea79190611fa6565b50565b6101206040519081016040528060008152602001600060ff16815260200160608152602001600073ffffffffffffffffffffffffffffffffffffffff168152602001600073ffffffffffffffffffffffffffffffffffffffff168152602001600081526020016000815260200160008152602001606081525090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10611f6757805160ff1916838001178555611f95565b82800160010185558215611f95579182015b82811115611f94578251825591602001919060010190611f79565b5b509050611fa29190612067565b5090565b61206491905b80821115612060576000808201600090556001820160006101000a81549060ff0219169055600282016000611fe1919061208c565b6003820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff02191690556004820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055600582016000905560068201600090556007820160009055600882016000612057919061208c565b50600901611fac565b5090565b90565b61208991905b8082111561208557600081600090555060010161206d565b5090565b90565b50805460018160011615610100020316600290046000825580601f106120b257506120d1565b601f0160209004906000526020600020908101906120d09190612067565b5b505600a165627a7a72305820ecf661654fcaf56608ec2e66dac68c0b286e5a258b9f849cc8f6fa8fa38187840029`

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

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() constant returns(uint8)
func (_Assigntoken *AssigntokenCaller) Enabled(opts *bind.CallOpts) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _Assigntoken.contract.Call(opts, out, "enabled")
	return *ret0, err
}

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() constant returns(uint8)
func (_Assigntoken *AssigntokenSession) Enabled() (uint8, error) {
	return _Assigntoken.Contract.Enabled(&_Assigntoken.CallOpts)
}

// Enabled is a free data retrieval call binding the contract method 0x238dafe0.
//
// Solidity: function enabled() constant returns(uint8)
func (_Assigntoken *AssigntokenCallerSession) Enabled() (uint8, error) {
	return _Assigntoken.Contract.Enabled(&_Assigntoken.CallOpts)
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

// KeyDefault is a free data retrieval call binding the contract method 0x1a08a37c.
//
// Solidity: function keyDefault() constant returns(string)
func (_Assigntoken *AssigntokenCaller) KeyDefault(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Assigntoken.contract.Call(opts, out, "keyDefault")
	return *ret0, err
}

// KeyDefault is a free data retrieval call binding the contract method 0x1a08a37c.
//
// Solidity: function keyDefault() constant returns(string)
func (_Assigntoken *AssigntokenSession) KeyDefault() (string, error) {
	return _Assigntoken.Contract.KeyDefault(&_Assigntoken.CallOpts)
}

// KeyDefault is a free data retrieval call binding the contract method 0x1a08a37c.
//
// Solidity: function keyDefault() constant returns(string)
func (_Assigntoken *AssigntokenCallerSession) KeyDefault() (string, error) {
	return _Assigntoken.Contract.KeyDefault(&_Assigntoken.CallOpts)
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
// Solidity: function assignToken() returns()
func (_Assigntoken *AssigntokenTransactor) AssignToken(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Assigntoken.contract.Transact(opts, "assignToken")
}

// AssignToken is a paid mutator transaction binding the contract method 0xa237213c.
//
// Solidity: function assignToken() returns()
func (_Assigntoken *AssigntokenSession) AssignToken() (*types.Transaction, error) {
	return _Assigntoken.Contract.AssignToken(&_Assigntoken.TransactOpts)
}

// AssignToken is a paid mutator transaction binding the contract method 0xa237213c.
//
// Solidity: function assignToken() returns()
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
