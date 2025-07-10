package loadtest

import (
	"context"
	"math/big"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/go-logr/logr"
	"github.com/grafana/sobek"
	"github.com/mysteryforge/gasper/k6/eth"
	"github.com/mysteryforge/gasper/k6/mock"
	"github.com/stretchr/testify/require"
	k6common "go.k6.io/k6/js/common"
	"go.k6.io/k6/lib"
)

func TestNewClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	once := sync.Once{}
	var rpcUrl string
	var wallets []string
	var privateKeys []string
	startRPC := func() {
		once.Do(func() {
			url, pks := mock.StartTestAnvilContainer(t, ctx)
			rpcUrl = url
			wallets = pks[1:]
			privateKeys = pks[0:1]
		})
	}

	t.Run("successful client creation", func(t *testing.T) {
		startRPC()
		cfg := &clientConfig{
			HTTP:            rpcUrl,
			Wallets:         wallets,
			TargetAddresses: []string{"0xc78260046895c358dE4bE97210Efca3900544905"},
		}
		c, err := NewClient(ctx, cfg, "test", nil, logr.Discard())
		require.NoError(t, err)
		require.Equal(t, c.ethClient.ChainID.Int64(), int64(mock.TestChainID))
	})
	t.Run("invalid config", func(t *testing.T) {
		cfg := &clientConfig{
			HTTP: "invalid-url",
		}
		_, err := NewClient(ctx, cfg, "test", nil, logr.Discard())
		require.Error(t, err)
	})
	t.Run("funding new wallets", func(t *testing.T) {
		startRPC()
		cfg := &clientConfig{
			HTTP:            rpcUrl,
			NumWallets:      2,
			PrivateKeys:     privateKeys,
			FundAmount:      *big.NewInt(1000000000000000000),
			TargetAddresses: []string{"0xc78260046895c358dE4bE97210Efca3900544905"},
		}
		c, err := NewClient(ctx, cfg, "test", nil, logr.Discard())
		require.NoError(t, err)
		require.Equal(t, c.testers.Ln(), 2)
		testers := c.testers.All()
		for _, wallet := range testers {
			balance, err := c.ethClient.Ec.BalanceAt(ctx, wallet.Address, nil)
			require.NoError(t, err)
			require.Equal(t, balance.String(), "1000000000000000000")
			require.Equal(t, wallet.Nonce, uint64(0))
			require.Equal(t, wallet.OffsetNonce, uint64(0))
		}
	})
	t.Run("creating existing wallets", func(t *testing.T) {
		startRPC()
		cfg := &clientConfig{
			HTTP:            rpcUrl,
			Wallets:         wallets,
			TargetAddresses: []string{"0xc78260046895c358dE4bE97210Efca3900544905"},
		}
		c, err := NewClient(ctx, cfg, "test", nil, logr.Discard())
		require.NoError(t, err)
		require.Equal(t, c.testers.Ln(), len(wallets))
		testers := c.testers.All()
		for _, wallet := range testers {
			balance, err := c.ethClient.Ec.BalanceAt(ctx, wallet.Address, nil)
			require.NoError(t, err)
			require.Less(t, balance.Int64(), int64(1000000000000000000))
		}
	})
	t.Run("erc20 contract", func(t *testing.T) {
		startRPC()
		cfg := &clientConfig{
			HTTP:            rpcUrl,
			Wallets:         wallets,
			TargetAddresses: []string{"0xc78260046895c358dE4bE97210Efca3900544905"},
			PrivateKeys:     privateKeys,
			ERC20:           true,
			ERC20MintAmount: *big.NewInt(100000000),
		}
		c, err := NewClient(ctx, cfg, "test", nil, logr.Discard())
		require.NoError(t, err)
		require.NotNil(t, c.erc20)

		// create new client with existing contract
		cfg = &clientConfig{
			HTTP:            rpcUrl,
			Wallets:         wallets,
			TargetAddresses: []string{"0xc78260046895c358dE4bE97210Efca3900544905"},
			ERC20:           true,
			ERC20Address:    c.erc20.Address.Hex(),
			ERC20MintAmount: *big.NewInt(100000000),
		}
		c, err = NewClient(ctx, cfg, "test", nil, logr.Discard())
		require.NoError(t, err)
		require.NotNil(t, c.erc20)
	})
	t.Run("create erc721 contract", func(t *testing.T) {
		startRPC()
		cfg := &clientConfig{
			HTTP:            rpcUrl,
			Wallets:         wallets,
			TargetAddresses: []string{"0xc78260046895c358dE4bE97210Efca3900544905"},
			PrivateKeys:     privateKeys,
			ERC721:          true,
			ERC721Mint:      true,
		}
		c, err := NewClient(ctx, cfg, "test", nil, logr.Discard())
		require.NoError(t, err)
		require.NotNil(t, c.erc721)

		// create new client with existing contract
		cfg = &clientConfig{
			HTTP:            rpcUrl,
			Wallets:         wallets,
			TargetAddresses: []string{"0xc78260046895c358dE4bE97210Efca3900544905"},
			PrivateKeys:     privateKeys,
			ERC721:          true,
			ERC721Address:   c.erc721.Address.Hex(),
			ERC721Mint:      true,
		}
		c, err = NewClient(ctx, cfg, "test", nil, logr.Discard())
		require.NoError(t, err)
		require.NotNil(t, c.erc721)
	})
	t.Run("batch funder contract", func(t *testing.T) {
		startRPC()
		cfg := &clientConfig{
			HTTP:            rpcUrl,
			TargetAddresses: []string{"0xc78260046895c358dE4bE97210Efca3900544905"},
			PrivateKeys:     privateKeys,
			NumWallets:      5,
			FundAmount:      *big.NewInt(1000000000000000000),
		}
		c, err := NewClient(ctx, cfg, "test", nil, logr.Discard())
		require.NoError(t, err)
		require.NotNil(t, c.batchFunder)
		require.NotNil(t, c.batchFunder.Address)
		require.Equal(t, 5, c.testers.Ln())

		// Check all wallets were funded correctly
		testers := c.testers.All()
		for _, wallet := range testers {
			balance, err := c.ethClient.Ec.BalanceAt(ctx, wallet.Address, nil)
			require.NoError(t, err)
			require.Equal(t, "1000000000000000000", balance.String())
		}
	})
}

func TestClientTransactions(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	once := sync.Once{}
	var rpcUrl string
	var wallets []string
	vu := mockVU{ctx: ctx}
	startRPC := func() {
		once.Do(func() {
			url, pks := mock.StartTestAnvilContainer(t, ctx)
			rpcUrl = url
			wallets = pks[1:]
		})
	}

	t.Run("send transaction", func(t *testing.T) {
		startRPC()
		cfg := &clientConfig{
			HTTP:            rpcUrl,
			Wallets:         wallets,
			TargetAddresses: []string{"0xc78260046895c358dE4bE97210Efca3900544905"},
		}
		c, err := NewClient(ctx, cfg, "test", nil, logr.Discard())
		require.NoError(t, err)
		hash, err := c.SendTransaction(&vu, nil)
		require.NoError(t, err)
		require.NotNil(t, hash)

		// Wait for transaction to be mined
		receipt, err := eth.WaitUntilMined(ctx, c.ethClient.Ec, *hash, eth.EthDefaultBlockTime, 10*time.Millisecond)
		require.NoError(t, err)
		require.Equal(t, receipt.Status, uint64(1))
	})

	t.Run("tx pool status", func(t *testing.T) {
		startRPC()
		cfg := &clientConfig{
			HTTP:            rpcUrl,
			Wallets:         wallets,
			TargetAddresses: []string{"0xc78260046895c358dE4bE97210Efca3900544905"},
		}
		c, err := NewClient(ctx, cfg, "test", nil, logr.Discard())
		require.NoError(t, err)

		status, err := c.TxPoolStatus(&vu, nil)
		require.NoError(t, err)
		require.NotNil(t, status)
		require.IsType(t, &eth.PoolStatus{
			BaseFee: 0,
			Pending: 0,
			Queued:  0,
		}, status)
	})

	t.Run("contract deployment and interaction", func(t *testing.T) {
		startRPC()
		cfg := &clientConfig{
			HTTP:            rpcUrl,
			Wallets:         wallets,
			TargetAddresses: []string{"0xc78260046895c358dE4bE97210Efca3900544905"},
		}
		c, err := NewClient(ctx, cfg, "test", nil, logr.Discard())
		require.NoError(t, err)

		// Create temporary ABI and BIN files for testing
		abiContent := `[{"inputs":[{"internalType":"uint256","name":"_value","type":"uint256"}],"name":"setValue","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"getValue","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"}]`
		binContent := "608060405234801561001057600080fd5b50610150806100206000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c8063209652551461003b5780635524107714610059575b600080fd5b610043610075565b60405161005091906100a1565b60405180910390f35b610073600480360381019061006e91906100ed565b61007e565b005b60008054905090565b8060008190555050565b6000819050919050565b61009b81610088565b82525050565b60006020820190506100b66000830184610092565b92915050565b600080fd5b6100ca81610088565b81146100d557600080fd5b50565b6000813590506100e7816100c1565b92915050565b600060208284031215610103576101026100bc565b5b6000610111848285016100d8565b9150509291505056fea2646970667358221220b20aa356f2c3ba45d8c19e7c2c92677c2c3e3d7845ae0ee6a9f8c3b9b7e8e98664736f6c63430008140033"

		abiFile := t.TempDir() + "/test.abi"
		binFile := t.TempDir() + "/test.bin"

		require.NoError(t, os.WriteFile(abiFile, []byte(abiContent), 0644))
		require.NoError(t, os.WriteFile(binFile, []byte(binContent), 0644))

		// Test contract deployment
		wallet, hash, addr, err := c.DeployContract(&vu, &eth.DeployContractParams{
			GasLimit: 300000,
			AbiPath:  abiFile,
			BinPath:  binFile,
			Args:     []interface{}{},
		})
		require.NoError(t, err)
		require.NotNil(t, wallet)
		require.NotNil(t, hash)
		require.NotNil(t, addr)

		// Test contract transaction
		txHash, err := c.TxContract(&vu, nil, &eth.FnContractParams{
			ContractAddress: *addr,
			Method:          "setValue",
			Args:            []interface{}{big.NewInt(100)},
		})
		require.NoError(t, err)
		require.NotNil(t, txHash)

		// Wait for transaction to be mined
		receipt, err := eth.WaitUntilMined(ctx, c.ethClient.Ec, *txHash, eth.EthDefaultBlockTime, 10*time.Millisecond)
		require.NoError(t, err)
		require.Equal(t, receipt.Status, uint64(1))

		// Test contract call
		result, err := c.CallContract(&vu, nil, &eth.FnContractParams{
			ContractAddress: *addr,
			Method:          "getValue",
			Args:            []interface{}{},
		})
		require.NoError(t, err)
		require.Equal(t, big.NewInt(100), result[0])
	})
}

func TestClientERC20Operations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	once := sync.Once{}
	var rpcUrl string
	var wallets []string
	var privateKeys []string
	vu := mockVU{ctx: ctx}
	startRPC := func() {
		once.Do(func() {
			url, pks := mock.StartTestAnvilContainer(t, ctx)
			rpcUrl = url
			wallets = pks[1:]
			privateKeys = pks[0:1]
		})
	}

	t.Run("deploy and interact with ERC20", func(t *testing.T) {
		startRPC()
		cfg := &clientConfig{
			HTTP:            rpcUrl,
			Wallets:         wallets,
			PrivateKeys:     privateKeys,
			ERC20:           true,
			ERC20MintAmount: *big.NewInt(100000000000000),
			TargetAddresses: []string{"0xc78260046895c358dE4bE97210Efca3900544905"},
		}
		c, err := NewClient(ctx, cfg, "test", nil, logr.Discard())
		require.NoError(t, err)
		require.NotNil(t, c.erc20)
		require.NotNil(t, c.erc20.Address)

		// Test ERC20 transaction
		hash, err := c.SendERC20Transaction(&vu, nil)
		require.NoError(t, err)
		require.NotNil(t, hash)
	})
}

func TestClientERC721Operations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	once := sync.Once{}
	var rpcUrl string
	var wallets []string
	var privateKeys []string
	vu := mockVU{ctx: ctx}
	startRPC := func() {
		once.Do(func() {
			url, pks := mock.StartTestAnvilContainer(t, ctx)
			rpcUrl = url
			wallets = pks[1:]
			privateKeys = pks[0:1]
		})
	}
	t.Run("deploy and interact with ERC721", func(t *testing.T) {
		startRPC()
		cfg := &clientConfig{
			HTTP:            rpcUrl,
			Wallets:         wallets,
			PrivateKeys:     privateKeys,
			ERC721:          true,
			ERC721Mint:      true,
			TargetAddresses: []string{"0xc78260046895c358dE4bE97210Efca3900544905"},
		}
		c, err := NewClient(ctx, cfg, "test", nil, logr.Discard())
		require.NoError(t, err)
		require.NotNil(t, c.erc721)
		require.NotNil(t, c.erc721.Address)

		// Test ERC721 transaction
		hash, err := c.SendERC721Transaction(&vu, nil)
		require.NoError(t, err)
		require.NotNil(t, hash)
	})
}

func TestTxInfoOperation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	once := sync.Once{}
	var rpcUrl string
	var wallets []string
	var privateKeys []string
	vu := mockVU{ctx: ctx}
	startRPC := func() {
		once.Do(func() {
			url, pks := mock.StartTestAnvilContainer(t, ctx)
			rpcUrl = url
			wallets = pks[1:]
			privateKeys = pks[0:1]
		})
	}
	t.Run("send eoa and check gas utilisation", func(t *testing.T) {
		startRPC()
		cfg := &clientConfig{
			HTTP:            rpcUrl,
			Wallets:         wallets,
			PrivateKeys:     privateKeys,
			ERC721:          true,
			ERC721Mint:      true,
			TargetAddresses: []string{"0xc78260046895c358dE4bE97210Efca3900544905"},
		}
		c, err := NewClient(ctx, cfg, "test", nil, logr.Discard())
		require.NoError(t, err)
		require.NotNil(t, c.erc721)
		require.NotNil(t, c.erc721.Address)

		// send a simple eoa transaction
		hash, err := c.SendTransaction(&vu, nil)
		require.NoError(t, err)
		require.NotNil(t, hash)

		// now get the tx info
		txInfo, err := c.GetTxInfo(&vu, nil, *hash)
		require.NoError(t, err)
		require.NotNil(t, txInfo)
		require.Equal(t, txInfo.GasUsed, uint64(21000))
	})
}

type mockVU struct {
	ctx context.Context
}

func (m *mockVU) Context() context.Context {
	return m.ctx
}

func (m *mockVU) InitEnv() *k6common.InitEnvironment {
	return nil
}

func (m *mockVU) State() *lib.State {
	return nil
}

func (m *mockVU) Runtime() *sobek.Runtime {
	return nil
}

func (m *mockVU) Events() k6common.Events {
	return k6common.Events{}
}

func (m *mockVU) RegisterCallback() (enqueueCallback func(func() error)) {
	return nil
}

// func newMockEthMetrics() *EthMetrics {
// 	r := metrics.NewRegistry()
// 	return &EthMetrics{
// 		RequestDuration:   r.MustNewMetric("mock_gasper_req_duration", metrics.Trend, metrics.Time),
// 		TimeToMine:        r.MustNewMetric("mock_gasper_time_to_mine", metrics.Trend, metrics.Time),
// 		Block:             r.MustNewMetric("mock_gasper_block", metrics.Counter, metrics.Default),
// 		GasUsed:           r.MustNewMetric("mock_gasper_gas_used", metrics.Trend, metrics.Default),
// 		TPS:               r.MustNewMetric("mock_gasper_tps", metrics.Trend, metrics.Default),
// 		EOA:               r.MustNewMetric("mock_gasper_eoa", metrics.Counter, metrics.Default),
// 		BlockTime:         r.MustNewMetric("mock_gasper_block_time", metrics.Trend, metrics.Time),
// 		BlockPerSec:       r.MustNewMetric("mock_gasper_block_per_sec", metrics.Trend, metrics.Default),
// 		PoolStatusPending: r.MustNewMetric("mock_gasper_pool_status_pending", metrics.Trend, metrics.Default),
// 		PoolStatusQueued:  r.MustNewMetric("mock_gasper_pool_status_queued", metrics.Trend, metrics.Default),
// 	}
// }
