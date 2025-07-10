package loadtest

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mysteryforge/gasper/k6/eth"
	"github.com/stretchr/testify/assert"
	"go.k6.io/k6/js/modules"
)

type mockClient struct {
	uid                 string
	chainID             *big.Int
	ChainIDFunc         func(vu modules.VU, metrics *EthMetrics) (*big.Int, error)
	TxPoolStatusFunc    func(vu modules.VU, metrics *EthMetrics) (*eth.PoolStatus, error)
	SendTransactionFunc func(vu modules.VU, metrics *EthMetrics, options ...TransactionOption) (*common.Hash, error)
	CallContractFunc    func(vu modules.VU, metrics *EthMetrics, params *eth.FnContractParams) ([]interface{}, error)
	CallFunc            func(vu modules.VU, metrics *EthMetrics, method string, params []interface{}) (interface{}, error)
}

func (m *mockClient) UID() string {
	return m.uid
}

func (m *mockClient) ChainID(vu modules.VU, metrics *EthMetrics) (*big.Int, error) {
	if m.ChainIDFunc != nil {
		return m.ChainIDFunc(vu, metrics)
	}
	return m.chainID, nil
}

func (m *mockClient) TxPoolStatus(vu modules.VU, metrics *EthMetrics) (*eth.PoolStatus, error) {
	if m.TxPoolStatusFunc != nil {
		return m.TxPoolStatusFunc(vu, metrics)
	}
	return &eth.PoolStatus{Pending: 0, Queued: 0}, nil
}

func (m *mockClient) SendTransaction(vu modules.VU, metrics *EthMetrics, options ...TransactionOption) (*common.Hash, error) {
	if m.SendTransactionFunc != nil {
		return m.SendTransactionFunc(vu, metrics, options...)
	}
	hash := common.HexToHash("0x3ad06070b524694608f48556b19fab89d0a5b7b558ce8e753c8648c3e0ca6b38")
	return &hash, nil
}

func (m *mockClient) CallContract(vu modules.VU, metrics *EthMetrics, params *eth.FnContractParams) ([]interface{}, error) {
	if m.CallContractFunc != nil {
		return m.CallContractFunc(vu, metrics, params)
	}
	return []interface{}{}, nil
}

func (m *mockClient) SendERC20Transaction(vu modules.VU, metrics *EthMetrics, options ...TransactionOption) (*common.Hash, error) {
	hash := common.HexToHash("0x3ad06070b524694608f48556b19fab89d0a5b7b558ce8e753c8648c3e0ca6b38")
	return &hash, nil
}

func (m *mockClient) SendERC721Transaction(vu modules.VU, metrics *EthMetrics, options ...TransactionOption) (*common.Hash, error) {
	hash := common.HexToHash("0x3ad06070b524694608f48556b19fab89d0a5b7b558ce8e753c8648c3e0ca6b38")
	return &hash, nil
}

func (m *mockClient) DeployContract(vu modules.VU, params *eth.DeployContractParams) (*common.Address, *common.Hash, *common.Address, error) {
	addr := common.HexToAddress("0x2A71e39B76B99645FDaFDfa9d38c0a51815d0941")
	hash := common.HexToHash("0x8bdde6587bb8f486bb71a605939bbdd19084e64ed8569bc21e92d4d279cb16c9")
	wallet := common.HexToAddress("0x653DB51224aBa51949534A895f522e50687f3C13")
	return &wallet, &hash, &addr, nil
}

func (m *mockClient) TxContract(vu modules.VU, metrics *EthMetrics, params *eth.FnContractParams) (*common.Hash, error) {
	hash := common.HexToHash("0x3ad06070b524694608f48556b19fab89d0a5b7b558ce8e753c8648c3e0ca6b38")
	return &hash, nil
}

func (m *mockClient) GetTxInfo(vu modules.VU, metrics *EthMetrics, txHash common.Hash) (*eth.TransactionInfo, error) {
	return &eth.TransactionInfo{
		TxHash:  "0x3ad06070b524694608f48556b19fab89d0a5b7b558ce8e753c8648c3e0ca6b38",
		GasUsed: 21000,
	}, nil

}

func (m *mockClient) ReportBlockMetrics(vu modules.VU, metrics *EthMetrics) error {
	return nil
}

func (m *mockClient) RequestSharedWallet() (*eth.Wallet, error) {
	return nil, nil
}

func (m *mockClient) ReleaseSharedWallet(address common.Address) {}

func (m *mockClient) Call(vu modules.VU, metrics *EthMetrics, method string, args ...interface{}) (interface{}, error) {
	return nil, nil
}

func (m *mockClient) Close() {}

func TestChainID(t *testing.T) {
	t.Run("successful chain ID retrieval", func(t *testing.T) {
		clients := &Clients{
			list: []Client{
				&mockClient{uid: "0", chainID: big.NewInt(1)},
				&mockClient{uid: "1", chainID: big.NewInt(2)},
			},
		}
		results := clients.ChainID(nil, nil)
		assert.Equal(t, 2, len(results))
		assert.Equal(t, big.NewInt(1), results["0"].Data)
		assert.Equal(t, big.NewInt(2), results["1"].Data)
	})
	t.Run("error handling", func(t *testing.T) {
		clients := &Clients{
			list: []Client{
				&mockClient{uid: "0", chainID: big.NewInt(1)},
				&mockClient{uid: "1", chainID: big.NewInt(2), ChainIDFunc: func(vu modules.VU, metrics *EthMetrics) (*big.Int, error) {
					return nil, assert.AnError
				}},
			},
		}
		results := clients.ChainID(nil, nil)
		assert.Equal(t, 2, len(results))
		assert.Equal(t, big.NewInt(1), results["0"].Data)
		assert.Equal(t, assert.AnError, results["1"].Err)
	})
}

func TestTxPoolStatus(t *testing.T) {
	t.Run("successful pool status retrieval", func(t *testing.T) {
		clients := &Clients{
			list: []Client{
				&mockClient{uid: "0"},
				&mockClient{uid: "1"},
			},
		}
		results := clients.TxPoolStatus(nil, nil)
		assert.Equal(t, 2, len(results))
		assert.Equal(t, &eth.PoolStatus{Pending: 0, Queued: 0}, results["0"].Data)
		assert.Equal(t, &eth.PoolStatus{Pending: 0, Queued: 0}, results["1"].Data)
	})

	t.Run("error handling", func(t *testing.T) {
		clients := &Clients{
			list: []Client{
				&mockClient{uid: "0"},
				&mockClient{uid: "1", TxPoolStatusFunc: func(vu modules.VU, metrics *EthMetrics) (*eth.PoolStatus, error) {
					return nil, assert.AnError
				}},
			},
		}
		results := clients.TxPoolStatus(nil, nil)
		assert.Equal(t, 2, len(results))
		assert.Equal(t, &eth.PoolStatus{Pending: 0, Queued: 0}, results["0"].Data)
		assert.Equal(t, assert.AnError, results["1"].Err)
	})
}

func TestSendTransaction(t *testing.T) {
	t.Run("successful transaction sending", func(t *testing.T) {
		clients := &Clients{
			list: []Client{
				&mockClient{uid: "0"},
				&mockClient{uid: "1"},
			},
		}
		results := clients.SendTransaction(nil, nil, map[string]interface{}{"tx_count": 1})
		assert.Equal(t, 2, len(results))
		assert.Equal(t, "0x3ad06070b524694608f48556b19fab89d0a5b7b558ce8e753c8648c3e0ca6b38", results["0"].Data)
		assert.Equal(t, "0x3ad06070b524694608f48556b19fab89d0a5b7b558ce8e753c8648c3e0ca6b38", results["1"].Data)
	})

	t.Run("error handling", func(t *testing.T) {
		clients := &Clients{
			list: []Client{
				&mockClient{uid: "0"},
				&mockClient{uid: "1", SendTransactionFunc: func(vu modules.VU, metrics *EthMetrics, options ...TransactionOption) (*common.Hash, error) {
					return nil, assert.AnError
				}},
			},
		}
		results := clients.SendTransaction(nil, nil, map[string]interface{}{"tx_count": 1})
		assert.Equal(t, 2, len(results))
		assert.Equal(t, "0x3ad06070b524694608f48556b19fab89d0a5b7b558ce8e753c8648c3e0ca6b38", results["0"].Data)
		assert.Equal(t, assert.AnError, results["1"].Err)
	})
}

func TestSendERC20Transaction(t *testing.T) {
	t.Run("successful ERC20 transaction", func(t *testing.T) {
		clients := &Clients{
			list: []Client{
				&mockClient{uid: "0"},
				&mockClient{uid: "1"},
			},
		}
		results := clients.SendERC20Transaction(nil, nil, map[string]interface{}{"tx_count": 1})
		assert.Equal(t, 2, len(results))
		assert.Equal(t, "0x3ad06070b524694608f48556b19fab89d0a5b7b558ce8e753c8648c3e0ca6b38", results["0"].Data)
		assert.Equal(t, "0x3ad06070b524694608f48556b19fab89d0a5b7b558ce8e753c8648c3e0ca6b38", results["1"].Data)
	})
}

func TestSendERC721Transaction(t *testing.T) {
	t.Run("successful ERC721 transaction", func(t *testing.T) {
		clients := &Clients{
			list: []Client{
				&mockClient{uid: "0"},
				&mockClient{uid: "1"},
			},
		}
		results := clients.SendERC721Transaction(nil, nil, map[string]interface{}{"tx_count": 1})
		assert.Equal(t, 2, len(results))
		assert.Equal(t, "0x3ad06070b524694608f48556b19fab89d0a5b7b558ce8e753c8648c3e0ca6b38", results["0"].Data)
		assert.Equal(t, "0x3ad06070b524694608f48556b19fab89d0a5b7b558ce8e753c8648c3e0ca6b38", results["1"].Data)
	})
}

func TestDeployContract(t *testing.T) {
	t.Run("successful contract deployment", func(t *testing.T) {
		clients := &Clients{
			list: []Client{
				&mockClient{uid: "0"},
				&mockClient{uid: "1"},
			},
		}
		results := clients.DeployContract(nil, map[string]interface{}{"abi_path": "test.abi", "bin_path": "test.bin", "gas_limit": int64(300000), "args": []interface{}{}})
		assert.Equal(t, 2, len(results))

		expectedResponse := &ContractDeploymentResponse{
			TransactionHash: "0x8bdde6587bb8f486bb71a605939bbdd19084e64ed8569bc21e92d4d279cb16c9",
			ContractAddress: "0x2A71e39B76B99645FDaFDfa9d38c0a51815d0941",
			OwnerWallet:     "0x653DB51224aBa51949534A895f522e50687f3C13",
		}
		assert.Equal(t, expectedResponse, results["0"].Data)
		assert.Equal(t, expectedResponse, results["1"].Data)
	})
}

func TestCallContract(t *testing.T) {
	t.Run("successful contract call", func(t *testing.T) {
		clients := &Clients{
			list: []Client{
				&mockClient{uid: "0"},
				&mockClient{uid: "1"},
			},
		}
		results := clients.CallContract(nil, nil, map[string]interface{}{"contract_address": "0x2A71e39B76B99645FDaFDfa9d38c0a51815d0941", "method": "test", "args": []interface{}{}})
		assert.Equal(t, 2, len(results))
		assert.Equal(t, []interface{}{}, results["0"].Data)
		assert.Equal(t, []interface{}{}, results["1"].Data)
	})

	t.Run("error handling", func(t *testing.T) {
		clients := &Clients{
			list: []Client{
				&mockClient{uid: "0"},
				&mockClient{uid: "1", CallContractFunc: func(vu modules.VU, metrics *EthMetrics, params *eth.FnContractParams) ([]interface{}, error) {
					return nil, assert.AnError
				}},
			},
		}
		results := clients.CallContract(nil, nil, map[string]interface{}{"contract_address": "0x2A71e39B76B99645FDaFDfa9d38c0a51815d0941", "method": "test", "args": []interface{}{}})
		assert.Equal(t, 2, len(results))
		assert.Equal(t, []interface{}{}, results["0"].Data)
		assert.Equal(t, assert.AnError, results["1"].Err)
	})
}

func TestTxContract(t *testing.T) {
	t.Run("successful contract transaction", func(t *testing.T) {
		clients := &Clients{
			list: []Client{
				&mockClient{uid: "0"},
				&mockClient{uid: "1"},
			},
		}
		results := clients.TxContract(nil, nil, map[string]interface{}{"contract_address": "0x2A71e39B76B99645FDaFDfa9d38c0a51815d0941", "method": "test", "args": []interface{}{}})
		assert.Equal(t, 2, len(results))
		assert.Equal(t, "0x3ad06070b524694608f48556b19fab89d0a5b7b558ce8e753c8648c3e0ca6b38", results["0"].Data)
		assert.Equal(t, "0x3ad06070b524694608f48556b19fab89d0a5b7b558ce8e753c8648c3e0ca6b38", results["1"].Data)
	})
}

func TestReportBlockMetrics(t *testing.T) {
	t.Run("successful metrics reporting", func(t *testing.T) {
		clients := &Clients{
			list: []Client{
				&mockClient{uid: "0"},
				&mockClient{uid: "1"},
			},
		}
		results := clients.ReportBlockMetrics(nil, nil)
		assert.Equal(t, 2, len(results))
		assert.Nil(t, results["0"].Data)
		assert.Nil(t, results["0"].Err)
		assert.Nil(t, results["1"].Data)
		assert.Nil(t, results["1"].Err)
	})
}

func TestCall(t *testing.T) {
	t.Run("successful eth_chainId call", func(t *testing.T) {
		clients := &Clients{
			list: []Client{
				&mockClient{uid: "0"},
				&mockClient{uid: "1"},
			},
		}
		results := clients.Call(nil, nil, "eth_chaiId", []interface{}{})
		assert.Equal(t, 2, len(results))
		fmt.Println(results)
		assert.Nil(t, results["0"].Data)
		assert.Nil(t, results["0"].Err)
		assert.Nil(t, results["1"].Data)
		assert.Nil(t, results["1"].Err)
	})
}
