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

// BatchFunderMetaData contains all meta data concerning the BatchFunder contract.
var BatchFunderMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"receive\",\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"batchSend\",\"inputs\":[{\"name\":\"recipients\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"rescueEth\",\"inputs\":[{\"name\":\"to\",\"type\":\"address\",\"internalType\":\"addresspayable\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"BatchFunded\",\"inputs\":[{\"name\":\"totalRecipients\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"totalAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Funded\",\"inputs\":[{\"name\":\"recipient\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"OwnableInvalidOwner\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"OwnableUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}]}]",
	Bin: "0x608060405234801561000f575f5ffd5b50335f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610081575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016100789190610196565b60405180910390fd5b6100908161009660201b60201c565b506101af565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f61018082610157565b9050919050565b61019081610176565b82525050565b5f6020820190506101a95f830184610187565b92915050565b610a4e806101bc5f395ff3fe60806040526004361061004d575f3560e01c80633b1ab44c14610058578063715018a6146100805780638da5cb5b14610096578063aa2f5220146100c0578063f2fde38b146100dc57610054565b3661005457005b5f5ffd5b348015610063575f5ffd5b5061007e60048036038101906100799190610669565b610104565b005b34801561008b575f5ffd5b506100946101c1565b005b3480156100a1575f5ffd5b506100aa6101d4565b6040516100b791906106b4565b60405180910390f35b6100da60048036038101906100d59190610761565b6101fb565b005b3480156100e7575f5ffd5b5061010260048036038101906100fd91906107e8565b610434565b005b61010c6104b8565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff160361017a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016101719061086d565b60405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff166108fc4790811502906040515f60405180830381858888f193505050501580156101bd573d5f5f3e3d5ffd5b5050565b6101c96104b8565b6101d25f61053f565b565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b5f5f90505b838390508110156103f2575f73ffffffffffffffffffffffffffffffffffffffff168484838181106102355761023461088b565b5b905060200201602081019061024a91906107e8565b73ffffffffffffffffffffffffffffffffffffffff16036102a0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161029790610902565b60405180910390fd5b5f8484838181106102b4576102b361088b565b5b90506020020160208101906102c991906107e8565b73ffffffffffffffffffffffffffffffffffffffff16836040516102ec9061094d565b5f6040518083038185875af1925050503d805f8114610326576040519150601f19603f3d011682016040523d82523d5f602084013e61032b565b606091505b505090508061036f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610366906109ab565b60405180910390fd5b8484838181106103825761038161088b565b5b905060200201602081019061039791906107e8565b73ffffffffffffffffffffffffffffffffffffffff167f5af8184bef8e4b45eb9f6ed7734d04da38ced226495548f46e0c8ff8d7d9a524846040516103dc91906109d8565b60405180910390a2508080600101915050610200565b507f6de1673d12d5d5bbbdff1e13ab8a981dd05a7d08eb9a306f66934b25f5c9044883839050346040516104279291906109f1565b60405180910390a1505050565b61043c6104b8565b5f73ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff16036104ac575f6040517f1e4fbdf70000000000000000000000000000000000000000000000000000000081526004016104a391906106b4565b60405180910390fd5b6104b58161053f565b50565b6104c0610600565b73ffffffffffffffffffffffffffffffffffffffff166104de6101d4565b73ffffffffffffffffffffffffffffffffffffffff161461053d57610501610600565b6040517f118cdaa700000000000000000000000000000000000000000000000000000000815260040161053491906106b4565b60405180910390fd5b565b5f5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050815f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b5f33905090565b5f5ffd5b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6106388261060f565b9050919050565b6106488161062e565b8114610652575f5ffd5b50565b5f813590506106638161063f565b92915050565b5f6020828403121561067e5761067d610607565b5b5f61068b84828501610655565b91505092915050565b5f61069e8261060f565b9050919050565b6106ae81610694565b82525050565b5f6020820190506106c75f8301846106a5565b92915050565b5f5ffd5b5f5ffd5b5f5ffd5b5f5f83601f8401126106ee576106ed6106cd565b5b8235905067ffffffffffffffff81111561070b5761070a6106d1565b5b602083019150836020820283011115610727576107266106d5565b5b9250929050565b5f819050919050565b6107408161072e565b811461074a575f5ffd5b50565b5f8135905061075b81610737565b92915050565b5f5f5f6040848603121561077857610777610607565b5b5f84013567ffffffffffffffff8111156107955761079461060b565b5b6107a1868287016106d9565b935093505060206107b48682870161074d565b9150509250925092565b6107c781610694565b81146107d1575f5ffd5b50565b5f813590506107e2816107be565b92915050565b5f602082840312156107fd576107fc610607565b5b5f61080a848285016107d4565b91505092915050565b5f82825260208201905092915050565b7f496e76616c6964206164647265737300000000000000000000000000000000005f82015250565b5f610857600f83610813565b915061086282610823565b602082019050919050565b5f6020820190508181035f8301526108848161084b565b9050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b7f496e76616c696420726563697069656e740000000000000000000000000000005f82015250565b5f6108ec601183610813565b91506108f7826108b8565b602082019050919050565b5f6020820190508181035f830152610919816108e0565b9050919050565b5f81905092915050565b50565b5f6109385f83610920565b91506109438261092a565b5f82019050919050565b5f6109578261092d565b9150819050919050565b7f455448207472616e73666572206661696c6564000000000000000000000000005f82015250565b5f610995601383610813565b91506109a082610961565b602082019050919050565b5f6020820190508181035f8301526109c281610989565b9050919050565b6109d28161072e565b82525050565b5f6020820190506109eb5f8301846109c9565b92915050565b5f604082019050610a045f8301856109c9565b610a1160208301846109c9565b939250505056fea26469706673582212206defedd2058f360dbccff396a42abff363680301b890d45a155b6e25e2269af564736f6c634300081c0033",
}

// BatchFunderABI is the input ABI used to generate the binding from.
// Deprecated: Use BatchFunderMetaData.ABI instead.
var BatchFunderABI = BatchFunderMetaData.ABI

// BatchFunderBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BatchFunderMetaData.Bin instead.
var BatchFunderBin = BatchFunderMetaData.Bin

// DeployBatchFunder deploys a new Ethereum contract, binding an instance of BatchFunder to it.
func DeployBatchFunder(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BatchFunder, error) {
	parsed, err := BatchFunderMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BatchFunderBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BatchFunder{BatchFunderCaller: BatchFunderCaller{contract: contract}, BatchFunderTransactor: BatchFunderTransactor{contract: contract}, BatchFunderFilterer: BatchFunderFilterer{contract: contract}}, nil
}

// BatchFunder is an auto generated Go binding around an Ethereum contract.
type BatchFunder struct {
	BatchFunderCaller     // Read-only binding to the contract
	BatchFunderTransactor // Write-only binding to the contract
	BatchFunderFilterer   // Log filterer for contract events
}

// BatchFunderCaller is an auto generated read-only Go binding around an Ethereum contract.
type BatchFunderCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchFunderTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BatchFunderTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchFunderFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BatchFunderFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchFunderSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BatchFunderSession struct {
	Contract     *BatchFunder      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// BatchFunderCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BatchFunderCallerSession struct {
	Contract *BatchFunderCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// BatchFunderTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BatchFunderTransactorSession struct {
	Contract     *BatchFunderTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// BatchFunderRaw is an auto generated low-level Go binding around an Ethereum contract.
type BatchFunderRaw struct {
	Contract *BatchFunder // Generic contract binding to access the raw methods on
}

// BatchFunderCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BatchFunderCallerRaw struct {
	Contract *BatchFunderCaller // Generic read-only contract binding to access the raw methods on
}

// BatchFunderTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BatchFunderTransactorRaw struct {
	Contract *BatchFunderTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBatchFunder creates a new instance of BatchFunder, bound to a specific deployed contract.
func NewBatchFunder(address common.Address, backend bind.ContractBackend) (*BatchFunder, error) {
	contract, err := bindBatchFunder(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BatchFunder{BatchFunderCaller: BatchFunderCaller{contract: contract}, BatchFunderTransactor: BatchFunderTransactor{contract: contract}, BatchFunderFilterer: BatchFunderFilterer{contract: contract}}, nil
}

// NewBatchFunderCaller creates a new read-only instance of BatchFunder, bound to a specific deployed contract.
func NewBatchFunderCaller(address common.Address, caller bind.ContractCaller) (*BatchFunderCaller, error) {
	contract, err := bindBatchFunder(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BatchFunderCaller{contract: contract}, nil
}

// NewBatchFunderTransactor creates a new write-only instance of BatchFunder, bound to a specific deployed contract.
func NewBatchFunderTransactor(address common.Address, transactor bind.ContractTransactor) (*BatchFunderTransactor, error) {
	contract, err := bindBatchFunder(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BatchFunderTransactor{contract: contract}, nil
}

// NewBatchFunderFilterer creates a new log filterer instance of BatchFunder, bound to a specific deployed contract.
func NewBatchFunderFilterer(address common.Address, filterer bind.ContractFilterer) (*BatchFunderFilterer, error) {
	contract, err := bindBatchFunder(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BatchFunderFilterer{contract: contract}, nil
}

// bindBatchFunder binds a generic wrapper to an already deployed contract.
func bindBatchFunder(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BatchFunderMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BatchFunder *BatchFunderRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BatchFunder.Contract.BatchFunderCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BatchFunder *BatchFunderRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchFunder.Contract.BatchFunderTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BatchFunder *BatchFunderRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BatchFunder.Contract.BatchFunderTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BatchFunder *BatchFunderCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BatchFunder.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BatchFunder *BatchFunderTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchFunder.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BatchFunder *BatchFunderTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BatchFunder.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BatchFunder *BatchFunderCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BatchFunder.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BatchFunder *BatchFunderSession) Owner() (common.Address, error) {
	return _BatchFunder.Contract.Owner(&_BatchFunder.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BatchFunder *BatchFunderCallerSession) Owner() (common.Address, error) {
	return _BatchFunder.Contract.Owner(&_BatchFunder.CallOpts)
}

// BatchSend is a paid mutator transaction binding the contract method 0xaa2f5220.
//
// Solidity: function batchSend(address[] recipients, uint256 amount) payable returns()
func (_BatchFunder *BatchFunderTransactor) BatchSend(opts *bind.TransactOpts, recipients []common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BatchFunder.contract.Transact(opts, "batchSend", recipients, amount)
}

// BatchSend is a paid mutator transaction binding the contract method 0xaa2f5220.
//
// Solidity: function batchSend(address[] recipients, uint256 amount) payable returns()
func (_BatchFunder *BatchFunderSession) BatchSend(recipients []common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BatchFunder.Contract.BatchSend(&_BatchFunder.TransactOpts, recipients, amount)
}

// BatchSend is a paid mutator transaction binding the contract method 0xaa2f5220.
//
// Solidity: function batchSend(address[] recipients, uint256 amount) payable returns()
func (_BatchFunder *BatchFunderTransactorSession) BatchSend(recipients []common.Address, amount *big.Int) (*types.Transaction, error) {
	return _BatchFunder.Contract.BatchSend(&_BatchFunder.TransactOpts, recipients, amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BatchFunder *BatchFunderTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchFunder.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BatchFunder *BatchFunderSession) RenounceOwnership() (*types.Transaction, error) {
	return _BatchFunder.Contract.RenounceOwnership(&_BatchFunder.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_BatchFunder *BatchFunderTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _BatchFunder.Contract.RenounceOwnership(&_BatchFunder.TransactOpts)
}

// RescueEth is a paid mutator transaction binding the contract method 0x3b1ab44c.
//
// Solidity: function rescueEth(address to) returns()
func (_BatchFunder *BatchFunderTransactor) RescueEth(opts *bind.TransactOpts, to common.Address) (*types.Transaction, error) {
	return _BatchFunder.contract.Transact(opts, "rescueEth", to)
}

// RescueEth is a paid mutator transaction binding the contract method 0x3b1ab44c.
//
// Solidity: function rescueEth(address to) returns()
func (_BatchFunder *BatchFunderSession) RescueEth(to common.Address) (*types.Transaction, error) {
	return _BatchFunder.Contract.RescueEth(&_BatchFunder.TransactOpts, to)
}

// RescueEth is a paid mutator transaction binding the contract method 0x3b1ab44c.
//
// Solidity: function rescueEth(address to) returns()
func (_BatchFunder *BatchFunderTransactorSession) RescueEth(to common.Address) (*types.Transaction, error) {
	return _BatchFunder.Contract.RescueEth(&_BatchFunder.TransactOpts, to)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BatchFunder *BatchFunderTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _BatchFunder.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BatchFunder *BatchFunderSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BatchFunder.Contract.TransferOwnership(&_BatchFunder.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_BatchFunder *BatchFunderTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _BatchFunder.Contract.TransferOwnership(&_BatchFunder.TransactOpts, newOwner)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BatchFunder *BatchFunderTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchFunder.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BatchFunder *BatchFunderSession) Receive() (*types.Transaction, error) {
	return _BatchFunder.Contract.Receive(&_BatchFunder.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_BatchFunder *BatchFunderTransactorSession) Receive() (*types.Transaction, error) {
	return _BatchFunder.Contract.Receive(&_BatchFunder.TransactOpts)
}

// BatchFunderBatchFundedIterator is returned from FilterBatchFunded and is used to iterate over the raw logs and unpacked data for BatchFunded events raised by the BatchFunder contract.
type BatchFunderBatchFundedIterator struct {
	Event *BatchFunderBatchFunded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BatchFunderBatchFundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BatchFunderBatchFunded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BatchFunderBatchFunded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BatchFunderBatchFundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BatchFunderBatchFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BatchFunderBatchFunded represents a BatchFunded event raised by the BatchFunder contract.
type BatchFunderBatchFunded struct {
	TotalRecipients *big.Int
	TotalAmount     *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterBatchFunded is a free log retrieval operation binding the contract event 0x6de1673d12d5d5bbbdff1e13ab8a981dd05a7d08eb9a306f66934b25f5c90448.
//
// Solidity: event BatchFunded(uint256 totalRecipients, uint256 totalAmount)
func (_BatchFunder *BatchFunderFilterer) FilterBatchFunded(opts *bind.FilterOpts) (*BatchFunderBatchFundedIterator, error) {

	logs, sub, err := _BatchFunder.contract.FilterLogs(opts, "BatchFunded")
	if err != nil {
		return nil, err
	}
	return &BatchFunderBatchFundedIterator{contract: _BatchFunder.contract, event: "BatchFunded", logs: logs, sub: sub}, nil
}

// WatchBatchFunded is a free log subscription operation binding the contract event 0x6de1673d12d5d5bbbdff1e13ab8a981dd05a7d08eb9a306f66934b25f5c90448.
//
// Solidity: event BatchFunded(uint256 totalRecipients, uint256 totalAmount)
func (_BatchFunder *BatchFunderFilterer) WatchBatchFunded(opts *bind.WatchOpts, sink chan<- *BatchFunderBatchFunded) (event.Subscription, error) {

	logs, sub, err := _BatchFunder.contract.WatchLogs(opts, "BatchFunded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BatchFunderBatchFunded)
				if err := _BatchFunder.contract.UnpackLog(event, "BatchFunded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBatchFunded is a log parse operation binding the contract event 0x6de1673d12d5d5bbbdff1e13ab8a981dd05a7d08eb9a306f66934b25f5c90448.
//
// Solidity: event BatchFunded(uint256 totalRecipients, uint256 totalAmount)
func (_BatchFunder *BatchFunderFilterer) ParseBatchFunded(log types.Log) (*BatchFunderBatchFunded, error) {
	event := new(BatchFunderBatchFunded)
	if err := _BatchFunder.contract.UnpackLog(event, "BatchFunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BatchFunderFundedIterator is returned from FilterFunded and is used to iterate over the raw logs and unpacked data for Funded events raised by the BatchFunder contract.
type BatchFunderFundedIterator struct {
	Event *BatchFunderFunded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BatchFunderFundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BatchFunderFunded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BatchFunderFunded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BatchFunderFundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BatchFunderFundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BatchFunderFunded represents a Funded event raised by the BatchFunder contract.
type BatchFunderFunded struct {
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterFunded is a free log retrieval operation binding the contract event 0x5af8184bef8e4b45eb9f6ed7734d04da38ced226495548f46e0c8ff8d7d9a524.
//
// Solidity: event Funded(address indexed recipient, uint256 amount)
func (_BatchFunder *BatchFunderFilterer) FilterFunded(opts *bind.FilterOpts, recipient []common.Address) (*BatchFunderFundedIterator, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BatchFunder.contract.FilterLogs(opts, "Funded", recipientRule)
	if err != nil {
		return nil, err
	}
	return &BatchFunderFundedIterator{contract: _BatchFunder.contract, event: "Funded", logs: logs, sub: sub}, nil
}

// WatchFunded is a free log subscription operation binding the contract event 0x5af8184bef8e4b45eb9f6ed7734d04da38ced226495548f46e0c8ff8d7d9a524.
//
// Solidity: event Funded(address indexed recipient, uint256 amount)
func (_BatchFunder *BatchFunderFilterer) WatchFunded(opts *bind.WatchOpts, sink chan<- *BatchFunderFunded, recipient []common.Address) (event.Subscription, error) {

	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _BatchFunder.contract.WatchLogs(opts, "Funded", recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BatchFunderFunded)
				if err := _BatchFunder.contract.UnpackLog(event, "Funded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseFunded is a log parse operation binding the contract event 0x5af8184bef8e4b45eb9f6ed7734d04da38ced226495548f46e0c8ff8d7d9a524.
//
// Solidity: event Funded(address indexed recipient, uint256 amount)
func (_BatchFunder *BatchFunderFilterer) ParseFunded(log types.Log) (*BatchFunderFunded, error) {
	event := new(BatchFunderFunded)
	if err := _BatchFunder.contract.UnpackLog(event, "Funded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BatchFunderOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the BatchFunder contract.
type BatchFunderOwnershipTransferredIterator struct {
	Event *BatchFunderOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *BatchFunderOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BatchFunderOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(BatchFunderOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *BatchFunderOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BatchFunderOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BatchFunderOwnershipTransferred represents a OwnershipTransferred event raised by the BatchFunder contract.
type BatchFunderOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BatchFunder *BatchFunderFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*BatchFunderOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BatchFunder.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &BatchFunderOwnershipTransferredIterator{contract: _BatchFunder.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BatchFunder *BatchFunderFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *BatchFunderOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _BatchFunder.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BatchFunderOwnershipTransferred)
				if err := _BatchFunder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_BatchFunder *BatchFunderFilterer) ParseOwnershipTransferred(log types.Log) (*BatchFunderOwnershipTransferred, error) {
	event := new(BatchFunderOwnershipTransferred)
	if err := _BatchFunder.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
