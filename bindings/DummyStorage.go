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

// DummyStorageMetaData contains all meta data concerning the DummyStorage contract.
var DummyStorageMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_initial\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"addValue\",\"inputs\":[{\"name\":\"_value\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"slot1\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"slot2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"slot3\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"}]",
	Bin: "0x6080604052348015600e575f5ffd5b506040516102eb3803806102eb8339818101604052810190602e91906079565b805f81905550806001819055508060028190555050609f565b5f5ffd5b5f819050919050565b605b81604b565b81146064575f5ffd5b50565b5f815190506073816054565b92915050565b5f60208284031215608b57608a6047565b5b5f6096848285016067565b91505092915050565b61023f806100ac5f395ff3fe608060405234801561000f575f5ffd5b506004361061004a575f3560e01c80631f457cb51461004e5780635b9af12b1461006c578063924fe31514610088578063d987e6b5146100a6575b5f5ffd5b6100566100c4565b6040516100639190610137565b60405180910390f35b6100866004803603810190610081919061017e565b6100c9565b005b610090610113565b60405161009d9190610137565b60405180910390f35b6100ae610119565b6040516100bb9190610137565b60405180910390f35b5f5481565b805f5f8282546100d991906101d6565b925050819055508060015f8282546100f191906101d6565b925050819055508060025f82825461010991906101d6565b9250508190555050565b60025481565b60015481565b5f819050919050565b6101318161011f565b82525050565b5f60208201905061014a5f830184610128565b92915050565b5f5ffd5b61015d8161011f565b8114610167575f5ffd5b50565b5f8135905061017881610154565b92915050565b5f6020828403121561019357610192610150565b5b5f6101a08482850161016a565b91505092915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f6101e08261011f565b91506101eb8361011f565b9250828201905080821115610203576102026101a9565b5b9291505056fea264697066735822122092d4ab40ca87ff5dcfd06e2e69eb236e353eae78e7ecfadd88e4e6bb241cdf5064736f6c634300081c0033",
}

// DummyStorageABI is the input ABI used to generate the binding from.
// Deprecated: Use DummyStorageMetaData.ABI instead.
var DummyStorageABI = DummyStorageMetaData.ABI

// DummyStorageBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DummyStorageMetaData.Bin instead.
var DummyStorageBin = DummyStorageMetaData.Bin

// DeployDummyStorage deploys a new Ethereum contract, binding an instance of DummyStorage to it.
func DeployDummyStorage(auth *bind.TransactOpts, backend bind.ContractBackend, _initial *big.Int) (common.Address, *types.Transaction, *DummyStorage, error) {
	parsed, err := DummyStorageMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DummyStorageBin), backend, _initial)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DummyStorage{DummyStorageCaller: DummyStorageCaller{contract: contract}, DummyStorageTransactor: DummyStorageTransactor{contract: contract}, DummyStorageFilterer: DummyStorageFilterer{contract: contract}}, nil
}

// DummyStorage is an auto generated Go binding around an Ethereum contract.
type DummyStorage struct {
	DummyStorageCaller     // Read-only binding to the contract
	DummyStorageTransactor // Write-only binding to the contract
	DummyStorageFilterer   // Log filterer for contract events
}

// DummyStorageCaller is an auto generated read-only Go binding around an Ethereum contract.
type DummyStorageCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DummyStorageTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DummyStorageTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DummyStorageFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DummyStorageFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DummyStorageSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DummyStorageSession struct {
	Contract     *DummyStorage     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DummyStorageCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DummyStorageCallerSession struct {
	Contract *DummyStorageCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// DummyStorageTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DummyStorageTransactorSession struct {
	Contract     *DummyStorageTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// DummyStorageRaw is an auto generated low-level Go binding around an Ethereum contract.
type DummyStorageRaw struct {
	Contract *DummyStorage // Generic contract binding to access the raw methods on
}

// DummyStorageCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DummyStorageCallerRaw struct {
	Contract *DummyStorageCaller // Generic read-only contract binding to access the raw methods on
}

// DummyStorageTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DummyStorageTransactorRaw struct {
	Contract *DummyStorageTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDummyStorage creates a new instance of DummyStorage, bound to a specific deployed contract.
func NewDummyStorage(address common.Address, backend bind.ContractBackend) (*DummyStorage, error) {
	contract, err := bindDummyStorage(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DummyStorage{DummyStorageCaller: DummyStorageCaller{contract: contract}, DummyStorageTransactor: DummyStorageTransactor{contract: contract}, DummyStorageFilterer: DummyStorageFilterer{contract: contract}}, nil
}

// NewDummyStorageCaller creates a new read-only instance of DummyStorage, bound to a specific deployed contract.
func NewDummyStorageCaller(address common.Address, caller bind.ContractCaller) (*DummyStorageCaller, error) {
	contract, err := bindDummyStorage(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DummyStorageCaller{contract: contract}, nil
}

// NewDummyStorageTransactor creates a new write-only instance of DummyStorage, bound to a specific deployed contract.
func NewDummyStorageTransactor(address common.Address, transactor bind.ContractTransactor) (*DummyStorageTransactor, error) {
	contract, err := bindDummyStorage(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DummyStorageTransactor{contract: contract}, nil
}

// NewDummyStorageFilterer creates a new log filterer instance of DummyStorage, bound to a specific deployed contract.
func NewDummyStorageFilterer(address common.Address, filterer bind.ContractFilterer) (*DummyStorageFilterer, error) {
	contract, err := bindDummyStorage(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DummyStorageFilterer{contract: contract}, nil
}

// bindDummyStorage binds a generic wrapper to an already deployed contract.
func bindDummyStorage(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DummyStorageMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DummyStorage *DummyStorageRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DummyStorage.Contract.DummyStorageCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DummyStorage *DummyStorageRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DummyStorage.Contract.DummyStorageTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DummyStorage *DummyStorageRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DummyStorage.Contract.DummyStorageTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DummyStorage *DummyStorageCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DummyStorage.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DummyStorage *DummyStorageTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DummyStorage.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DummyStorage *DummyStorageTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DummyStorage.Contract.contract.Transact(opts, method, params...)
}

// Slot1 is a free data retrieval call binding the contract method 0x1f457cb5.
//
// Solidity: function slot1() view returns(uint256)
func (_DummyStorage *DummyStorageCaller) Slot1(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DummyStorage.contract.Call(opts, &out, "slot1")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Slot1 is a free data retrieval call binding the contract method 0x1f457cb5.
//
// Solidity: function slot1() view returns(uint256)
func (_DummyStorage *DummyStorageSession) Slot1() (*big.Int, error) {
	return _DummyStorage.Contract.Slot1(&_DummyStorage.CallOpts)
}

// Slot1 is a free data retrieval call binding the contract method 0x1f457cb5.
//
// Solidity: function slot1() view returns(uint256)
func (_DummyStorage *DummyStorageCallerSession) Slot1() (*big.Int, error) {
	return _DummyStorage.Contract.Slot1(&_DummyStorage.CallOpts)
}

// Slot2 is a free data retrieval call binding the contract method 0xd987e6b5.
//
// Solidity: function slot2() view returns(uint256)
func (_DummyStorage *DummyStorageCaller) Slot2(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DummyStorage.contract.Call(opts, &out, "slot2")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Slot2 is a free data retrieval call binding the contract method 0xd987e6b5.
//
// Solidity: function slot2() view returns(uint256)
func (_DummyStorage *DummyStorageSession) Slot2() (*big.Int, error) {
	return _DummyStorage.Contract.Slot2(&_DummyStorage.CallOpts)
}

// Slot2 is a free data retrieval call binding the contract method 0xd987e6b5.
//
// Solidity: function slot2() view returns(uint256)
func (_DummyStorage *DummyStorageCallerSession) Slot2() (*big.Int, error) {
	return _DummyStorage.Contract.Slot2(&_DummyStorage.CallOpts)
}

// Slot3 is a free data retrieval call binding the contract method 0x924fe315.
//
// Solidity: function slot3() view returns(uint256)
func (_DummyStorage *DummyStorageCaller) Slot3(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DummyStorage.contract.Call(opts, &out, "slot3")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Slot3 is a free data retrieval call binding the contract method 0x924fe315.
//
// Solidity: function slot3() view returns(uint256)
func (_DummyStorage *DummyStorageSession) Slot3() (*big.Int, error) {
	return _DummyStorage.Contract.Slot3(&_DummyStorage.CallOpts)
}

// Slot3 is a free data retrieval call binding the contract method 0x924fe315.
//
// Solidity: function slot3() view returns(uint256)
func (_DummyStorage *DummyStorageCallerSession) Slot3() (*big.Int, error) {
	return _DummyStorage.Contract.Slot3(&_DummyStorage.CallOpts)
}

// AddValue is a paid mutator transaction binding the contract method 0x5b9af12b.
//
// Solidity: function addValue(uint256 _value) returns()
func (_DummyStorage *DummyStorageTransactor) AddValue(opts *bind.TransactOpts, _value *big.Int) (*types.Transaction, error) {
	return _DummyStorage.contract.Transact(opts, "addValue", _value)
}

// AddValue is a paid mutator transaction binding the contract method 0x5b9af12b.
//
// Solidity: function addValue(uint256 _value) returns()
func (_DummyStorage *DummyStorageSession) AddValue(_value *big.Int) (*types.Transaction, error) {
	return _DummyStorage.Contract.AddValue(&_DummyStorage.TransactOpts, _value)
}

// AddValue is a paid mutator transaction binding the contract method 0x5b9af12b.
//
// Solidity: function addValue(uint256 _value) returns()
func (_DummyStorage *DummyStorageTransactorSession) AddValue(_value *big.Int) (*types.Transaction, error) {
	return _DummyStorage.Contract.AddValue(&_DummyStorage.TransactOpts, _value)
}
