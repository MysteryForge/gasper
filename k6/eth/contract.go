package eth

import (
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type FnContractParams struct {
	ContractAddress    common.Address
	Method             string
	GasPriceMultiplier uint64
	AccessList         types.AccessList
	Args               []interface{}
}

func ParseFnContractParams(params map[string]interface{}) (*FnContractParams, error) {
	contractAddr, ok := params["contract_address"].(string)
	if !ok {
		return nil, errors.New("fn contract_address must be a string")
	}

	method, ok := params["method"].(string)
	if !ok {
		return nil, errors.New("fn method must be a string")
	}

	rawArgs, ok := params["args"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("args must be an array")
	}

	gasPriceMultiplier, ok := params["gas_price_multiplier"].(int64)
	if !ok {
		gasPriceMultiplier = 1
	}

	args := make([]map[string]interface{}, 0, len(rawArgs))
	for _, arg := range rawArgs {
		m, ok := arg.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid arg format")
		}
		args = append(args, m)
	}

	parsedArgs, err := ParseContractArguments(args)
	if err != nil {
		return nil, err
	}

	var accessList types.AccessList
	if parsedAccessList, ok := params["access_list"].([]interface{}); ok {
		accessList, err = ParseAccessList(parsedAccessList)
		if err != nil {
			return nil, fmt.Errorf("failed to parse access list: %w", err)
		}
	}

	return &FnContractParams{
		ContractAddress:    common.HexToAddress(contractAddr),
		Method:             method,
		Args:               parsedArgs,
		GasPriceMultiplier: uint64(gasPriceMultiplier),
		AccessList:         accessList,
	}, nil
}

// ParseArguments parses the JS args array into Go-typed values
func ParseContractArguments(args []map[string]interface{}) ([]interface{}, error) {
	parsed := make([]interface{}, 0, len(args))

	for _, arg := range args {
		typ, ok := arg["type"].(string)
		if !ok {
			return nil, errors.New("argument missing 'type' field")
		}

		val, ok := arg["value"]
		if !ok {
			return nil, errors.New("argument missing 'value' field")
		}

		switch strings.ToLower(typ) {
		case "uint256", "uint":
			vStr := fmt.Sprintf("%v", val) // convert to string safely
			v := new(big.Int)
			_, ok := v.SetString(vStr, 10)
			if !ok {
				return nil, fmt.Errorf("invalid uint256 value: %v", val)
			}
			parsed = append(parsed, v)

		case "address":
			addressStr, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("address value must be a string: %v", val)
			}
			parsed = append(parsed, common.HexToAddress(addressStr))

		case "bool":
			boolVal, ok := val.(bool)
			if !ok {
				return nil, fmt.Errorf("bool value must be a boolean: %v", val)
			}
			parsed = append(parsed, boolVal)

		case "string":
			strVal, ok := val.(string)
			if !ok {
				return nil, fmt.Errorf("string value must be a string: %v", val)
			}
			parsed = append(parsed, strVal)

		default:
			return nil, fmt.Errorf("unsupported constructor arg type: %s", typ)
		}
	}

	return parsed, nil
}

func ParseAccessList(accessList []interface{}) (types.AccessList, error) {
	result := types.AccessList{}
	for _, item := range accessList {
		tup := types.AccessTuple{}
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("access list item must be an object: %v", item)
		}

		address, ok := itemMap["address"].(string)
		if !ok {
			return nil, fmt.Errorf("access list item must have an address field: %v", item)
		}
		tup.Address = common.HexToAddress(address)

		tup.StorageKeys = []common.Hash{}
		sk, ok := itemMap["storage_keys"]
		if ok {
			skList, ok := sk.([]interface{})
			if !ok {
				return nil, fmt.Errorf("storageKeys must be an array: %v", sk)
			}
			for _, skItem := range skList {
				skStr, ok := skItem.(string)
				if !ok {
					return nil, fmt.Errorf("storage key must be a string: %v", skItem)
				}
				tup.StorageKeys = append(tup.StorageKeys, common.HexToHash(skStr))
			}
		}

		result = append(result, tup)
	}

	return result, nil
}

type DeployContractParams struct {
	GasLimit           uint64
	AbiPath            string
	BinPath            string
	GasPriceMultiplier uint64
	Args               []interface{}
}

func ParseDeployContractParams(params map[string]interface{}) (*DeployContractParams, error) {
	gasLimit, ok := params["gas_limit"].(int64)
	if !ok {
		return nil, errors.New("gas_limit must be a number")
	}

	abiPath, ok := params["abi_path"].(string)
	if !ok {
		return nil, errors.New("abi_path must be a string")
	}

	binPath, ok := params["bin_path"].(string)
	if !ok {
		return nil, errors.New("bin_path must be a string")
	}

	rawArgs, ok := params["args"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("args must be an array")
	}

	gasPriceMultiplier, ok := params["gas_price_multiplier"].(int64)
	if !ok {
		gasPriceMultiplier = 1
	}

	args := make([]map[string]interface{}, 0, len(rawArgs))
	for _, arg := range rawArgs {
		m, ok := arg.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid arg format")
		}
		args = append(args, m)
	}

	parsedArgs, err := ParseContractArguments(args)
	if err != nil {
		return nil, err
	}

	return &DeployContractParams{
		GasLimit:           uint64(gasLimit),
		AbiPath:            abiPath,
		BinPath:            binPath,
		Args:               parsedArgs,
		GasPriceMultiplier: uint64(gasPriceMultiplier),
	}, nil
}
