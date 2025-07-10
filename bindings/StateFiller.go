// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// StateFillerMetaData contains all meta data concerning the StateFiller contract.
var StateFillerMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"count\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"deleteRandomState\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getItems\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"size\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"}]",
	Bin: "0x608060405234801561000f575f5ffd5b50604051610681380380610681833981810160405281019061003191906100ee565b3360015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505f5f90505b818110156100b0575f81908060018154018082558091505060019003905f5260205f20015f90919091909150558080600101915050610076565b5050610119565b5f5ffd5b5f819050919050565b6100cd816100bb565b81146100d7575f5ffd5b50565b5f815190506100e8816100c4565b92915050565b5f60208284031215610103576101026100b7565b5b5f610110848285016100da565b91505092915050565b61055b806101265f395ff3fe608060405234801561000f575f5ffd5b506004361061004a575f3560e01c8063796c5e941461004e5780637b3118e11461007e5780638da5cb5b14610088578063949d225d146100a6575b5f5ffd5b6100686004803603810190610063919061024a565b6100c4565b6040516100759190610284565b60405180910390f35b6100866100e8565b005b6100906101e3565b60405161009d91906102dc565b60405180910390f35b6100ae610208565b6040516100bb9190610284565b60405180910390f35b5f5f82815481106100d8576100d76102f5565b5b905f5260205f2001549050919050565b5f5f805490501161012e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101259061037c565b60405180910390fd5b5f5f805490505f80549050423360405160200161014d939291906103ff565b604051602081830303815290604052805190602001205f1c61016f9190610468565b90505f60015f8054905061018391906104c5565b81548110610194576101936102f5565b5b905f5260205f2001545f82815481106101b0576101af6102f5565b5b905f5260205f2001819055505f8054806101cd576101cc6104f8565b5b600190038181905f5260205f20015f9055905550565b60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b5f5f80549050905090565b5f5ffd5b5f819050919050565b61022981610217565b8114610233575f5ffd5b50565b5f8135905061024481610220565b92915050565b5f6020828403121561025f5761025e610213565b5b5f61026c84828501610236565b91505092915050565b61027e81610217565b82525050565b5f6020820190506102975f830184610275565b92915050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6102c68261029d565b9050919050565b6102d6816102bc565b82525050565b5f6020820190506102ef5f8301846102cd565b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f82825260208201905092915050565b7f4e6f206d6f726520737461746520746f2064656c6574650000000000000000005f82015250565b5f610366601783610322565b915061037182610332565b602082019050919050565b5f6020820190508181035f8301526103938161035a565b9050919050565b5f819050919050565b6103b46103af82610217565b61039a565b82525050565b5f8160601b9050919050565b5f6103d0826103ba565b9050919050565b5f6103e1826103c6565b9050919050565b6103f96103f4826102bc565b6103d7565b82525050565b5f61040a82866103a3565b60208201915061041a82856103a3565b60208201915061042a82846103e8565b601482019150819050949350505050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f61047282610217565b915061047d83610217565b92508261048d5761048c61043b565b5b828206905092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6104cf82610217565b91506104da83610217565b92508282039050818111156104f2576104f1610498565b5b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603160045260245ffdfea2646970667358221220958f9180b7d3acdacf2a86340842fe7fe5589bd0f5de2c7faf23de17058a1c2a64736f6c634300081c0033",
}

// StateFillerABI is the input ABI used to generate the binding from.
// Deprecated: Use StateFillerMetaData.ABI instead.
var StateFillerABI = StateFillerMetaData.ABI

// StateFillerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use StateFillerMetaData.Bin instead.
var StateFillerBin = StateFillerMetaData.Bin

// DeployStateFiller deploys a new Ethereum contract, binding an instance of StateFiller to it.
func DeployStateFiller(auth *bind.TransactOpts, backend bind.ContractBackend, count *big.Int) (common.Address, *types.Transaction, *StateFiller, error) {
	parsed, err := StateFillerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(StateFillerBin), backend, count)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &StateFiller{StateFillerCaller: StateFillerCaller{contract: contract}, StateFillerTransactor: StateFillerTransactor{contract: contract}, StateFillerFilterer: StateFillerFilterer{contract: contract}}, nil
}

// StateFiller is an auto generated Go binding around an Ethereum contract.
type StateFiller struct {
	StateFillerCaller     // Read-only binding to the contract
	StateFillerTransactor // Write-only binding to the contract
	StateFillerFilterer   // Log filterer for contract events
}

// StateFillerCaller is an auto generated read-only Go binding around an Ethereum contract.
type StateFillerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateFillerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StateFillerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateFillerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StateFillerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StateFillerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StateFillerSession struct {
	Contract     *StateFiller      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StateFillerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StateFillerCallerSession struct {
	Contract *StateFillerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// StateFillerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StateFillerTransactorSession struct {
	Contract     *StateFillerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// StateFillerRaw is an auto generated low-level Go binding around an Ethereum contract.
type StateFillerRaw struct {
	Contract *StateFiller // Generic contract binding to access the raw methods on
}

// StateFillerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StateFillerCallerRaw struct {
	Contract *StateFillerCaller // Generic read-only contract binding to access the raw methods on
}

// StateFillerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StateFillerTransactorRaw struct {
	Contract *StateFillerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStateFiller creates a new instance of StateFiller, bound to a specific deployed contract.
func NewStateFiller(address common.Address, backend bind.ContractBackend) (*StateFiller, error) {
	contract, err := bindStateFiller(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &StateFiller{StateFillerCaller: StateFillerCaller{contract: contract}, StateFillerTransactor: StateFillerTransactor{contract: contract}, StateFillerFilterer: StateFillerFilterer{contract: contract}}, nil
}

// NewStateFillerCaller creates a new read-only instance of StateFiller, bound to a specific deployed contract.
func NewStateFillerCaller(address common.Address, caller bind.ContractCaller) (*StateFillerCaller, error) {
	contract, err := bindStateFiller(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StateFillerCaller{contract: contract}, nil
}

// NewStateFillerTransactor creates a new write-only instance of StateFiller, bound to a specific deployed contract.
func NewStateFillerTransactor(address common.Address, transactor bind.ContractTransactor) (*StateFillerTransactor, error) {
	contract, err := bindStateFiller(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StateFillerTransactor{contract: contract}, nil
}

// NewStateFillerFilterer creates a new log filterer instance of StateFiller, bound to a specific deployed contract.
func NewStateFillerFilterer(address common.Address, filterer bind.ContractFilterer) (*StateFillerFilterer, error) {
	contract, err := bindStateFiller(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StateFillerFilterer{contract: contract}, nil
}

// bindStateFiller binds a generic wrapper to an already deployed contract.
func bindStateFiller(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := StateFillerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StateFiller *StateFillerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StateFiller.Contract.StateFillerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StateFiller *StateFillerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StateFiller.Contract.StateFillerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StateFiller *StateFillerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StateFiller.Contract.StateFillerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_StateFiller *StateFillerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _StateFiller.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_StateFiller *StateFillerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StateFiller.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_StateFiller *StateFillerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _StateFiller.Contract.contract.Transact(opts, method, params...)
}

// GetItems is a free data retrieval call binding the contract method 0x796c5e94.
//
// Solidity: function getItems(uint256 index) view returns(uint256)
func (_StateFiller *StateFillerCaller) GetItems(opts *bind.CallOpts, index *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _StateFiller.contract.Call(opts, &out, "getItems", index)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetItems is a free data retrieval call binding the contract method 0x796c5e94.
//
// Solidity: function getItems(uint256 index) view returns(uint256)
func (_StateFiller *StateFillerSession) GetItems(index *big.Int) (*big.Int, error) {
	return _StateFiller.Contract.GetItems(&_StateFiller.CallOpts, index)
}

// GetItems is a free data retrieval call binding the contract method 0x796c5e94.
//
// Solidity: function getItems(uint256 index) view returns(uint256)
func (_StateFiller *StateFillerCallerSession) GetItems(index *big.Int) (*big.Int, error) {
	return _StateFiller.Contract.GetItems(&_StateFiller.CallOpts, index)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StateFiller *StateFillerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _StateFiller.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StateFiller *StateFillerSession) Owner() (common.Address, error) {
	return _StateFiller.Contract.Owner(&_StateFiller.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_StateFiller *StateFillerCallerSession) Owner() (common.Address, error) {
	return _StateFiller.Contract.Owner(&_StateFiller.CallOpts)
}

// Size is a free data retrieval call binding the contract method 0x949d225d.
//
// Solidity: function size() view returns(uint256)
func (_StateFiller *StateFillerCaller) Size(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _StateFiller.contract.Call(opts, &out, "size")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Size is a free data retrieval call binding the contract method 0x949d225d.
//
// Solidity: function size() view returns(uint256)
func (_StateFiller *StateFillerSession) Size() (*big.Int, error) {
	return _StateFiller.Contract.Size(&_StateFiller.CallOpts)
}

// Size is a free data retrieval call binding the contract method 0x949d225d.
//
// Solidity: function size() view returns(uint256)
func (_StateFiller *StateFillerCallerSession) Size() (*big.Int, error) {
	return _StateFiller.Contract.Size(&_StateFiller.CallOpts)
}

// DeleteRandomState is a paid mutator transaction binding the contract method 0x7b3118e1.
//
// Solidity: function deleteRandomState() returns()
func (_StateFiller *StateFillerTransactor) DeleteRandomState(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _StateFiller.contract.Transact(opts, "deleteRandomState")
}

// DeleteRandomState is a paid mutator transaction binding the contract method 0x7b3118e1.
//
// Solidity: function deleteRandomState() returns()
func (_StateFiller *StateFillerSession) DeleteRandomState() (*types.Transaction, error) {
	return _StateFiller.Contract.DeleteRandomState(&_StateFiller.TransactOpts)
}

// DeleteRandomState is a paid mutator transaction binding the contract method 0x7b3118e1.
//
// Solidity: function deleteRandomState() returns()
func (_StateFiller *StateFillerTransactorSession) DeleteRandomState() (*types.Transaction, error) {
	return _StateFiller.Contract.DeleteRandomState(&_StateFiller.TransactOpts)
}
