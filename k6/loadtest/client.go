package loadtest

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/go-logr/logr"
	"github.com/mysteryforge/gasper/k6/eth"
	"go.k6.io/k6/js/modules"
)

type Client interface {
	// Chain related
	ChainID(vu modules.VU, metrics *EthMetrics) (*big.Int, error)
	TxPoolStatus(vu modules.VU, metrics *EthMetrics) (*eth.PoolStatus, error)
	ReportBlockMetrics(vu modules.VU, metrics *EthMetrics) error

	// Transaction related
	SendTransaction(vu modules.VU, metrics *EthMetrics, options ...TransactionOption) (*common.Hash, error)
	SendERC20Transaction(vu modules.VU, metrics *EthMetrics, options ...TransactionOption) (*common.Hash, error)
	SendERC721Transaction(vu modules.VU, metrics *EthMetrics, options ...TransactionOption) (*common.Hash, error)

	// Contract related
	DeployContract(vu modules.VU, params *eth.DeployContractParams) (*common.Address, *common.Hash, *common.Address, error)
	TxContract(vu modules.VU, metrics *EthMetrics, params *eth.FnContractParams) (*common.Hash, error)
	CallContract(vu modules.VU, metrics *EthMetrics, params *eth.FnContractParams) ([]interface{}, error)
	GetTxInfo(vu modules.VU, metrics *EthMetrics, txHash common.Hash) (*eth.TransactionInfo, error)

	Call(vu modules.VU, metrics *EthMetrics, method string, args ...interface{}) (interface{}, error)

	RequestSharedWallet() (*eth.Wallet, error)
	ReleaseSharedWallet(address common.Address)

	// Cleanup
	Close()

	UID() string
}

type DefaultClient struct {
	ethClient         *eth.Client
	log               logr.Logger
	uid               string
	scenarioUID       string
	testers           *eth.WalletRegistry
	sponsors          *eth.WalletRegistry
	targetAddresses   *eth.TargetAddresses
	firstBlockNumber  uint64
	sendReportMux     *sync.Mutex
	erc20             *eth.ERC20
	erc721            *eth.ERC721
	batchFunder       *eth.BatchFunder
	deployedContracts map[common.Address]*bind.BoundContract
	reportBlock       *ReportBlock
	db                *eth.PebbleDb
	latestGasPrice    *eth.AtomicBigInt
	latestGasTip      *eth.AtomicBigInt
	txPool            *eth.PoolStatus
	txPoolRateLimiter *eth.TxPoolRateLimiter
	isLegacy          bool
	sharedWallets     map[string]*eth.Wallet
}

func NewClient(ctx context.Context, cfg *clientConfig, scenarioUID string, db *eth.PebbleDb, log logr.Logger) (*DefaultClient, error) {
	if err := validateClientConfig(cfg); err != nil {
		return nil, err
	}

	c := &DefaultClient{
		uid:               eth.SanitizeURL(cfg.HTTP),
		sendReportMux:     &sync.Mutex{},
		scenarioUID:       scenarioUID,
		deployedContracts: make(map[common.Address]*bind.BoundContract),
		db:                db,
		latestGasPrice:    eth.NewAtomicBigInt(nil),
		latestGasTip:      eth.NewAtomicBigInt(nil),
		txPool:            &eth.PoolStatus{},
		sharedWallets:     make(map[string]*eth.Wallet),
	}
	c.log = log.WithValues("uid", c.uid)

	var err error

	c.ethClient, err = eth.NewClient(ctx, cfg.HTTP)
	if err != nil {
		return nil, err
	}

	head, err := c.ethClient.Ec.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}
	if head.BaseFee == nil {
		c.isLegacy = true
	}

	c.targetAddresses, err = eth.NewTargetAddresses(cfg.TargetAddresses, cfg.NumTargetAddresses)
	if err != nil {
		return nil, err
	}

	c.firstBlockNumber, err = c.ethClient.Ec.BlockNumber(ctx)
	if err != nil {
		return nil, err
	}
	c.reportBlock = &ReportBlock{
		Number:        c.firstBlockNumber,
		TimestampMili: uint64(time.Now().UnixMilli()),
		Timestamp:     uint64(time.Now().Unix()),
	}
	latestBlock, err := c.ethClient.Ec.HeaderByNumber(ctx, new(big.Int).SetUint64(c.firstBlockNumber))
	if err != nil {
		return nil, err
	}
	timestampDelta := time.Since(time.Unix(int64(latestBlock.Time), 0)).Milliseconds()
	if timestampDelta > 10_000 {
		c.reportBlock.TimestampDelta = timestampDelta
	}

	if err := c.setupWalletsAndFund(
		ctx,
		cfg.PrivateKeys,
		cfg.Wallets,
		cfg.NumWallets,
		cfg.BatchFunderAddress,
		&cfg.FundAmount,
		cfg.MinGasPrice,
	); err != nil {
		return nil, err
	}

	if cfg.ERC20 {
		if err := c.setupERC20(ctx, cfg.ERC20Address, cfg.ERC20MintAmount, cfg.MinGasPrice); err != nil {
			return nil, err
		}
	}

	if cfg.ERC721 {
		if err := c.setupERC721(ctx, cfg.ERC721Address, cfg.ERC721Mint, cfg.MinGasPrice); err != nil {
			return nil, err
		}
	}

	if err := c.setupLatestGas(ctx, cfg.MinGasPrice); err != nil {
		return nil, err
	}

	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := c.setupLatestGas(ctx, cfg.MinGasPrice); err != nil {
					c.log.Error(err, "failed to get latest gas price")
				}
			}
		}
	}()

	if cfg.RateLimite != nil && *cfg.RateLimite > 0 {
		c.txPoolRateLimiter = eth.NewTxPoolRateLimiter(*cfg.RateLimite)
	}

	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := c.updateTxPoolStatus(ctx); err != nil {
					c.log.Error(err, "failed to get tx pool status")
				}
				if c.txPoolRateLimiter != nil && cfg.AdaptiveRateLimit {
					c.txPoolRateLimiter.UpdateTxPoolSize(uint64(c.txPool.Pending + c.txPool.Queued))
				}
			}
		}
	}()

	// listen for context cancellation and close the client
	go func() {
		<-ctx.Done()
		c.Close()
	}()

	return c, nil
}

func (c *DefaultClient) setupLatestGas(ctx context.Context, minGasPrice uint64) error {
	if !c.isLegacy {
		tipCap, err := c.ethClient.Ec.SuggestGasTipCap(ctx)
		if err != nil {
			return err
		}
		if tipCap.Cmp(big.NewInt(0)) == 0 {
			tipCap = big.NewInt(int64(minGasPrice)) // 1 gwei
		}
		c.latestGasTip.Store(tipCap)
	} else {
		gasPrice, err := c.ethClient.Ec.SuggestGasPrice(ctx)
		if err != nil {
			return err
		}
		if gasPrice.Cmp(big.NewInt(0)) == 0 {
			gasPrice = big.NewInt(int64(minGasPrice)) // 1 gwei
		}
		c.latestGasPrice.Store(gasPrice)
	}

	return nil

}

func (c *DefaultClient) setupWalletsAndFund(
	ctx context.Context,
	privateKeys []string,
	wallets []string,
	numNewWallets uint64,
	batchFunderAddress string,
	fundAmount *big.Int,
	minGasPrice uint64,
) error {
	ctx, cancel := context.WithTimeout(ctx, eth.WalletsTimeout)
	defer cancel()

	var err error
	if len(privateKeys) > 0 {
		c.sponsors, err = eth.NewWalletRegistryFromPrivateKeys(ctx, c.ethClient.Ec, privateKeys)
		if err != nil {
			return err
		}
	}

	if len(wallets) > 0 {
		c.testers, err = eth.NewWalletRegistryFromPrivateKeys(ctx, c.ethClient.Ec, wallets)
		if err != nil {
			return err
		}
	}

	if numNewWallets > 0 {
		if c.testers == nil {
			c.testers = eth.NewEmptyWalletRegistry()
		}
		if err := c.testers.GenerateAndStore(ctx, numNewWallets); err != nil {
			return err
		}

		if batchFunderAddress == "" {
			c.batchFunder, err = eth.DeployNewBatchFunder(ctx, c.ethClient, c.sponsors, minGasPrice)
			if err != nil {
				return err
			}
		} else {
			batchFunderAddress := common.HexToAddress(batchFunderAddress)
			c.batchFunder, err = eth.InitExistingBatchFunder(c.ethClient, batchFunderAddress)
			if err != nil {
				return err
			}
		}
		c.log.Info("Batch funder contract", "address", c.batchFunder.Address.Hex())
		if err := c.batchFunder.FundWallets(ctx, c.ethClient, c.sponsors, c.testers, fundAmount); err != nil {
			return err
		}
		c.log.Info("Funded wallets", "num_wallets", c.testers.Ln(), "fund_amount", fundAmount.String())
	}

	return nil
}

func (c *DefaultClient) setupERC20(ctx context.Context, address string, mintAmount big.Int, minGasPrice uint64) error {
	var err error
	if address == "" {
		c.erc20, err = eth.DeployNewERC20Contract(ctx, c.ethClient, c.sponsors, minGasPrice)
		if err != nil {
			return err
		}
	} else {
		erc20Address := common.HexToAddress(address)
		c.erc20, err = eth.InitExistingERC20Contract(c.ethClient.Ec, erc20Address)
		if err != nil {
			return err
		}
	}
	if mintAmount.Cmp(big.NewInt(0)) == 1 {
		if err := c.erc20.MintERC20Contract(ctx, c.ethClient, c.testers, mintAmount); err != nil {
			return err
		}
	}

	c.log.Info("ERC20 contract", "address", c.erc20.Address.Hex(), "mint_amount", mintAmount.String())
	return nil
}

func (c *DefaultClient) setupERC721(ctx context.Context, address string, mint bool, minGasPrice uint64) error {
	var err error
	if address == "" {
		c.erc721, err = eth.DeployNewERC721Contract(ctx, c.ethClient, c.sponsors, minGasPrice) // 1 gwei
		if err != nil {
			return err
		}
	} else {
		erc721Address := common.HexToAddress(address)
		c.erc721, err = eth.InitExistingERC721Contract(c.ethClient.Ec, erc721Address)
		if err != nil {
			return err
		}
	}
	if mint {
		if err := c.erc721.MintERC721Contract(ctx, c.ethClient, c.testers); err != nil {
			return err
		}
	}

	c.log.Info("ERC721 contract", "address", c.erc721.Address.Hex(), "mint", mint)
	return nil
}

func (c *DefaultClient) storeTransactionStartTime(hash common.Hash, t time.Time) {
	if c.db == nil {
		return
	}
	// TODO: potentially we will use worker pools here and sync.Pool if this becomes the bottleneck
	go func() {
		timeBytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(timeBytes, uint64(t.UnixMilli()))
		if err := c.db.Db().Set(c.db.GenKey("tx", hash.Hex()), timeBytes, pebble.NoSync); err != nil {
			c.log.Error(err, "failed to store tx hash", "hash", hash.Hex())
		}
	}()
}

func (c *DefaultClient) reportTransactionsLatency(vu modules.VU, metrics *EthMetrics, txs []string, t time.Time) {
	if c.db == nil {
		return
	}

	reportTx := func(hash string) {
		val, closer, err := c.db.Db().Get(c.db.GenKey("tx", hash))
		if err != nil {
			if err != pebble.ErrNotFound {
				c.log.Error(err, "failed to get tx hash", "hash", hash)
			}
			return
		}
		defer closer.Close() // nolint:errcheck
		startTime := time.UnixMilli(int64(binary.LittleEndian.Uint64(val)))
		ReportTimeToMineFromStats(vu, metrics, c.uid, t.Sub(startTime))
	}

	// TODO: potentially we will use worker pools here and sync.Pool if this becomes the bottleneck
	go func() {
		for _, tx := range txs {
			reportTx(tx)
		}
	}()
}

func (c *DefaultClient) Export() modules.Exports {
	return modules.Exports{}
}

func (c *DefaultClient) Close() {
	if c.db != nil {
		c.db.Close() // nolint:errcheck
	}
	c.ethClient.Close()
}

func (c *DefaultClient) UID() string {
	return c.uid
}

func (c *DefaultClient) RequestSharedWallet() (*eth.Wallet, error) {
	if c.testers == nil {
		return nil, fmt.Errorf("no available wallet")
	}
	w := c.testers.GetAvailableWallet()
	if w == nil {
		return nil, fmt.Errorf("no available wallet")
	}
	c.sharedWallets[w.Address.Hex()] = w
	return w, nil
}

func (c *DefaultClient) IsSharedWallet(address common.Address) bool {
	_, exists := c.sharedWallets[address.Hex()]
	return exists
}

func (c *DefaultClient) ReleaseSharedWallet(address common.Address) {
	if !c.IsSharedWallet(address) {
		return
	}
	delete(c.sharedWallets, address.Hex())
	c.testers.Unlock(address)
}

type ReportBlock struct {
	Number         uint64
	TimestampMili  uint64 // ms
	Timestamp      uint64 // seconds last block second
	TimestampDelta int64
}

func (c *DefaultClient) ReportBlockMetrics(vu modules.VU, metrics *EthMetrics) error {
	c.sendReportMux.Lock()
	defer c.sendReportMux.Unlock()

	ctx := vu.Context()

	blockNum, err := c.ethClient.Ec.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get block number: %w", err)
	}
	if blockNum <= c.reportBlock.Number {
		return nil
	}

	prevBlockTimestamp := c.reportBlock.Timestamp
	prevBlockTimestampMili := c.reportBlock.TimestampMili

	c.reportBlock.Number = blockNum
	c.reportBlock.Timestamp = uint64(time.Now().Unix())
	c.reportBlock.TimestampMili = uint64(time.Now().UnixMilli())

	block, err := c.ethClient.SlimBlockByNumber(ctx, new(big.Int).SetUint64(blockNum))
	if err != nil {
		return fmt.Errorf("failed to get block: %w", err)
	}
	if block == nil {
		return fmt.Errorf("block %d not found", blockNum)
	}

	var tps float64
	blockTimestampDiff := c.reportBlock.Timestamp - prevBlockTimestamp
	blockTimestampMiliDiff := c.reportBlock.TimestampMili - prevBlockTimestampMili
	txsLn := len(block.Transactions)
	if blockTimestampDiff >= 1 {
		tps = float64(txsLn) / float64(blockTimestampDiff)
	} else {
		tps = float64(txsLn)
	}

	mgas := float64(block.GasUsed) / (float64(blockTimestampMiliDiff) / 1000) / 1000000

	blockTimestampUnixMili := time.UnixMilli(int64(c.reportBlock.TimestampMili))
	block.Timestamp = block.Timestamp + uint64(c.reportBlock.TimestampDelta)
	ReportBlockMetrics(
		vu,
		metrics,
		c.uid,
		block,
		tps,
		mgas,
		blockTimestampMiliDiff,
		blockTimestampUnixMili,
	)

	c.reportTransactionsLatency(vu, metrics, block.Transactions, blockTimestampUnixMili)
	return nil
}

func (c *DefaultClient) Call(vu modules.VU, metrics *EthMetrics, method string, args ...interface{}) (interface{}, error) {
	var result interface{}
	err := c.ethClient.Rc.CallContext(vu.Context(), &result, method, args...)
	ReportReqDurationFromStats(vu, metrics, c.uid, method, time.Since(time.Now()))
	return result, err
}

func (c *DefaultClient) ChainID(vu modules.VU, metrics *EthMetrics) (*big.Int, error) {
	t := time.Now()
	r, err := c.ethClient.Ec.ChainID(vu.Context())
	ReportReqDurationFromStats(vu, metrics, c.uid, "chainId", time.Since(t))
	return r, err
}

func (c *DefaultClient) updateTxPoolStatus(ctx context.Context) error {
	var result eth.PoolStatus
	if err := c.ethClient.Rc.CallContext(ctx, &result, "txpool_status"); err != nil {
		return err
	}
	c.txPool = &result
	return nil
}

func (c *DefaultClient) TxPoolStatus(vu modules.VU, metrics *EthMetrics) (*eth.PoolStatus, error) {
	ReportTxPoolStatusFromStats(vu, metrics, c.uid, c.txPool)
	return c.txPool, nil
}

func (c *DefaultClient) GetTxInfo(vu modules.VU, metrics *EthMetrics, txHash common.Hash) (*eth.TransactionInfo, error) {
	t := time.Now()

	receipt, err := eth.WaitUntilMined(vu.Context(), c.ethClient.Ec, txHash, 20*time.Second, 500*time.Millisecond)
	if err != nil {
		return nil, err
	}

	ReportReqDurationFromStats(vu, metrics, c.uid, "getTxInfo", time.Since(t))

	return &eth.TransactionInfo{
		TxHash:  txHash.Hex(),
		Status:  receipt.Status,
		GasUsed: receipt.GasUsed,
	}, nil
}

type TransactionOptions struct {
	WaitForConfirmation bool
	ConfirmationDelay   time.Duration
	NoSend              bool
	OffsetNonce         uint64
	TxCount             uint64
	GasPriceMultiplier  uint64
	WalletAddresses     []*common.Address
}

type TransactionOption func(*TransactionOptions)

func WithTransactionConfirmationDelay(delay uint64) TransactionOption {
	return func(opts *TransactionOptions) {
		opts.WaitForConfirmation = true
		opts.ConfirmationDelay = time.Duration(delay) * time.Second
	}
}

func WithTransactionNoSend(send bool) TransactionOption {
	return func(opts *TransactionOptions) {
		opts.NoSend = send
	}
}

func WithTransactionOffsetNonce(offset uint64) TransactionOption {
	return func(opts *TransactionOptions) {
		opts.OffsetNonce = offset
	}
}

func WithTransactionTxCount(txCount uint64) TransactionOption {
	return func(opts *TransactionOptions) {
		opts.TxCount = txCount
	}
}

func WithTransactionGasPriceMultiplier(multiplier uint64) TransactionOption {
	return func(opts *TransactionOptions) {
		opts.GasPriceMultiplier = multiplier
	}
}

func WithTransactionWalletAddress(addresses []*common.Address) TransactionOption {
	return func(opts *TransactionOptions) {
		opts.WalletAddresses = addresses
	}
}

func DefaultTransactionOptions() *TransactionOptions {
	return &TransactionOptions{
		WaitForConfirmation: false,
		ConfirmationDelay:   eth.EthDefaultBlockTime,
		NoSend:              false,
		OffsetNonce:         0,
		TxCount:             1,
		GasPriceMultiplier:  1,
		WalletAddresses:     nil,
	}
}

func (c *DefaultClient) SendTransaction(vu modules.VU, metrics *EthMetrics, options ...TransactionOption) (*common.Hash, error) {
	opts := DefaultTransactionOptions()
	for _, opt := range options {
		opt(opts)
	}

	if opts.WaitForConfirmation && opts.OffsetNonce > 0 {
		return nil, fmt.Errorf("cannot use offset nonce with confirmation")
	}

	var wallet *eth.Wallet
	for _, address := range opts.WalletAddresses {
		if c.IsSharedWallet(*address) {
			wallet = c.sharedWallets[address.Hex()]
			break
		}
	}
	if wallet == nil {
		wallet = c.testers.GetAvailableWallet()
		if wallet == nil {
			return nil, fmt.Errorf("no available wallet")
		}
		defer c.testers.Unlock(wallet.Address)
	}

	var hash common.Hash

	errCh := make(chan error, int(opts.TxCount))
	ctx := vu.Context()
	targetAddress := c.targetAddresses.Random()

	retryForNonce := false
	sendTx := func() {
		if c.txPoolRateLimiter != nil {
			if err := c.txPoolRateLimiter.Wait(ctx); err != nil {
				errCh <- fmt.Errorf("rate limit: %w", err)
			}
		}

		var nonce uint64
		wallet.Lock()
		if retryForNonce {
			retryForNonce = false
		} else {
			nonce = wallet.Nonce
			if opts.OffsetNonce > 0 {
				if wallet.OffsetNonce <= wallet.Nonce {
					wallet.OffsetNonce = wallet.Nonce + opts.OffsetNonce
				}
				nonce = wallet.OffsetNonce
			}
			if opts.OffsetNonce == 0 {
				wallet.Nonce++
			}
			wallet.OffsetNonce++
		}
		wallet.Unlock()

		var tx *types.Transaction
		if !c.isLegacy {
			head, err := c.ethClient.Ec.HeaderByNumber(vu.Context(), nil)
			if err != nil {
				errCh <- err
				return
			}

			// EIP-1559 transaction
			tipCap := new(big.Int).Mul(c.latestGasTip.Load(), big.NewInt(int64(opts.GasPriceMultiplier)))
			feeCap := new(big.Int).Add(
				tipCap,
				new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
			)
			tx = types.NewTx(&types.DynamicFeeTx{
				ChainID:   c.ethClient.ChainID,
				Nonce:     nonce,
				To:        targetAddress,
				Value:     big.NewInt(1),
				Gas:       21000,
				GasTipCap: tipCap,
				GasFeeCap: feeCap,
				Data:      nil,
			})
		} else {
			// Legacy transaction
			tx = types.NewTx(&types.LegacyTx{
				Nonce:    nonce,
				To:       targetAddress,
				Value:    big.NewInt(1),
				Gas:      21000,
				GasPrice: new(big.Int).Mul(c.latestGasPrice.Load(), big.NewInt(int64(opts.GasPriceMultiplier))),
				Data:     nil,
			})
		}

		signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(c.ethClient.ChainID), wallet.PrivateKey)
		if err != nil {
			errCh <- err
			return
		}

		t := time.Now()
		if opts.NoSend {
			cm := ethereum.CallMsg{
				From:       wallet.Address,
				To:         signedTx.To(),
				Gas:        signedTx.Gas(),
				Value:      signedTx.Value(),
				Data:       signedTx.Data(),
				AccessList: signedTx.AccessList(),
			}
			if !c.isLegacy {
				cm.GasFeeCap = signedTx.GasFeeCap()
				cm.GasTipCap = signedTx.GasTipCap()
			} else {
				cm.GasPrice = signedTx.GasPrice()
			}

			if _, err := c.ethClient.Ec.CallContract(vu.Context(), cm, nil); err != nil {
				if strings.Contains(err.Error(), "fee cap less than block base fee") {
					retryForNonce = true
				}
				errCh <- err
				return
			}
		} else {
			if err := c.ethClient.Ec.SendTransaction(vu.Context(), signedTx); err != nil {
				retryForNonce = true
				if strings.Contains(err.Error(), "replacement transaction underpriced") && retryForNonce {
					retryForNonce = false
				}
				if strings.Contains(err.Error(), "transaction underpriced") && retryForNonce {
					retryForNonce = false
				}
				if strings.Contains(err.Error(), "nonce too low") && retryForNonce {
					retryForNonce = false
				}
				if strings.Contains(err.Error(), "already known") && retryForNonce {
					retryForNonce = false
				}
				if strings.Contains(err.Error(), "could not replace existing") && retryForNonce {
					retryForNonce = false
				}
				if strings.Contains(err.Error(), "fee cap less than block base fee") {
					retryForNonce = true
				}

				errCh <- err
				return
			}
		}

		ReportEoaFromStats(vu, metrics, c.uid, 1, eth.TransactionTypeETH)
		ReportReqDurationFromStats(vu, metrics, c.uid, "sendTransaction", time.Since(t))

		hash = signedTx.Hash()
		c.storeTransactionStartTime(hash, t)

		if opts.WaitForConfirmation {
			_, err := eth.WaitUntilMined(vu.Context(), c.ethClient.Ec, hash, opts.ConfirmationDelay, 10*time.Millisecond)
			if err != nil {
				errCh <- err
				return
			}

			ReportReqDurationFromStats(vu, metrics, c.uid, "sendConfirmedTransaction", time.Since(t))
		}
	}

	for i := 0; i < int(opts.TxCount); i++ {
		if ctx.Err() != nil {
			break
		}

		sendTx()
	}

	close(errCh)
	for err := range errCh {
		if err != nil && !errors.Is(err, context.Canceled) {
			c.log.Error(err, "error in sending transactions")
		}
	}

	if hash == (common.Hash{}) {
		return nil, fmt.Errorf("failed to send any transaction")
	}

	return &hash, nil
}

func (c *DefaultClient) SendERC20Transaction(vu modules.VU, metrics *EthMetrics, options ...TransactionOption) (*common.Hash, error) {
	if c.erc20 == nil {
		return nil, fmt.Errorf("erc20 contract is not initialized")
	}

	opts := DefaultTransactionOptions()
	for _, opt := range options {
		opt(opts)
	}

	if opts.WaitForConfirmation && opts.OffsetNonce > 0 {
		return nil, fmt.Errorf("cannot use offset nonce with confirmation")
	}

	var wallet *eth.Wallet
	for _, address := range opts.WalletAddresses {
		if c.IsSharedWallet(*address) {
			wallet = c.sharedWallets[address.Hex()]
			break
		}
	}
	if wallet == nil {
		wallet = c.testers.GetAvailableWallet()
		if wallet == nil {
			return nil, fmt.Errorf("no available wallet")
		}
		defer c.testers.Unlock(wallet.Address)
	}

	var hash common.Hash

	errCh := make(chan error, int(opts.TxCount))
	ctx := vu.Context()
	targetAddress := c.targetAddresses.Random()

	retryForNonce := false
	sendTx := func() {
		if c.txPoolRateLimiter != nil {
			if err := c.txPoolRateLimiter.Wait(ctx); err != nil {
				errCh <- fmt.Errorf("rate limit: %w", err)
			}
		}

		var nonce uint64
		wallet.Lock()
		if retryForNonce {
			retryForNonce = false
		} else {
			nonce = wallet.Nonce
			if opts.OffsetNonce > 0 {
				if wallet.OffsetNonce <= wallet.Nonce {
					wallet.OffsetNonce = wallet.Nonce + opts.OffsetNonce
				}
				nonce = wallet.OffsetNonce
			}
			if opts.OffsetNonce == 0 {
				wallet.Nonce++
			}
			wallet.OffsetNonce++
		}
		wallet.Unlock()

		signFn := func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
			sig := types.LatestSignerForChainID(c.ethClient.ChainID)
			return types.SignTx(tx, sig, wallet.PrivateKey)
		}

		var (
			tx  *types.Transaction
			err error
			t   time.Time
		)
		tops := bind.TransactOpts{
			From:    wallet.Address,
			Signer:  signFn,
			Context: vu.Context(),
			Nonce:   big.NewInt(int64(nonce)),
		}

		if !c.isLegacy {
			head, err := c.ethClient.Ec.HeaderByNumber(vu.Context(), nil)
			if err != nil {
				errCh <- err
				return
			}

			tipCap := new(big.Int).Mul(c.latestGasTip.Load(), big.NewInt(int64(opts.GasPriceMultiplier)))
			feeCap := new(big.Int).Add(
				tipCap,
				new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
			)
			tops.GasFeeCap = feeCap
			tops.GasTipCap = tipCap
		} else {
			tops.GasPrice = new(big.Int).Mul(c.latestGasPrice.Load(), big.NewInt(int64(opts.GasPriceMultiplier)))
		}

		if opts.NoSend {
			t = time.Now()
			tops.NoSend = true
			tx, err = c.erc20.Contract.Transfer(&tops, *targetAddress, big.NewInt(1))
			if err != nil {
				errCh <- err
				return
			}
			cm := ethereum.CallMsg{
				From:       wallet.Address,
				To:         tx.To(),
				Gas:        tx.Gas(),
				Value:      tx.Value(),
				Data:       tx.Data(),
				AccessList: tx.AccessList(),
			}
			if !c.isLegacy {
				cm.GasFeeCap = tx.GasFeeCap()
				cm.GasTipCap = tx.GasTipCap()
			} else {
				cm.GasPrice = tx.GasPrice()
			}

			if _, err := c.ethClient.Ec.CallContract(vu.Context(), cm, nil); err != nil {
				if strings.Contains(err.Error(), "fee cap less than block base fee") {
					retryForNonce = true
				}
				errCh <- err
				return
			}
		} else {
			t = time.Now()
			tx, err = c.erc20.Contract.Transfer(&tops, *targetAddress, big.NewInt(1))
			if err != nil {
				retryForNonce = true
				if strings.Contains(err.Error(), "replacement transaction underpriced") && retryForNonce {
					retryForNonce = false
				}
				if strings.Contains(err.Error(), "transaction underpriced") && retryForNonce {
					retryForNonce = false
				}
				if strings.Contains(err.Error(), "nonce too low") && retryForNonce {
					retryForNonce = false
				}
				if strings.Contains(err.Error(), "already known") && retryForNonce {
					retryForNonce = false
				}
				if strings.Contains(err.Error(), "could not replace existing") && retryForNonce {
					retryForNonce = false
				}
				if strings.Contains(err.Error(), "fee cap less than block base fee") {
					retryForNonce = true
				}

				errCh <- err
				return
			}
		}

		ReportEoaFromStats(vu, metrics, c.uid, 1, eth.TransactionTypeERC20)
		ReportReqDurationFromStats(vu, metrics, c.uid, "sendERC20Transaction", time.Since(t))

		hash = tx.Hash()
		c.storeTransactionStartTime(hash, t)

		if opts.WaitForConfirmation {
			_, err := eth.WaitUntilMined(vu.Context(), c.ethClient.Ec, hash, opts.ConfirmationDelay, 10*time.Millisecond)
			if err != nil {
				errCh <- err
				return
			}

			ReportReqDurationFromStats(vu, metrics, c.uid, "sendConfirmedERC20Transaction", time.Since(t))
		}
	}
	for i := 0; i < int(opts.TxCount); i++ {
		if ctx.Err() != nil {
			break
		}

		sendTx()
	}

	close(errCh)
	for err := range errCh {
		if err != nil && !errors.Is(err, context.Canceled) {
			c.log.Error(err, "error in sending transactions")
		}
	}

	if hash == (common.Hash{}) {
		return nil, fmt.Errorf("failed to send any transaction")
	}

	return &hash, nil
}

func (c *DefaultClient) SendERC721Transaction(vu modules.VU, metrics *EthMetrics, options ...TransactionOption) (*common.Hash, error) {
	if c.erc721 == nil {
		return nil, fmt.Errorf("erc721 contract is not initialized")
	}

	opts := DefaultTransactionOptions()
	for _, opt := range options {
		opt(opts)
	}

	if opts.WaitForConfirmation && opts.OffsetNonce > 0 {
		return nil, fmt.Errorf("cannot use offset nonce with confirmation")
	}

	var wallet *eth.Wallet
	for _, address := range opts.WalletAddresses {
		if c.IsSharedWallet(*address) {
			wallet = c.sharedWallets[address.Hex()]
			break
		}
	}
	if wallet == nil {
		wallet = c.testers.GetAvailableWallet()
		if wallet == nil {
			return nil, fmt.Errorf("no available wallet")
		}
		defer c.testers.Unlock(wallet.Address)
	}

	var hash common.Hash

	errCh := make(chan error, int(opts.TxCount))
	ctx := vu.Context()
	targetAddress := c.targetAddresses.Random()

	retryForNonce := false
	sendTx := func() {
		if c.txPoolRateLimiter != nil {
			if err := c.txPoolRateLimiter.Wait(ctx); err != nil {
				errCh <- fmt.Errorf("rate limit: %w", err)
			}
		}

		var nonce uint64
		wallet.Lock()
		if retryForNonce {
			retryForNonce = false
		} else {
			nonce = wallet.Nonce
			if opts.OffsetNonce > 0 {
				if wallet.OffsetNonce <= wallet.Nonce {
					wallet.OffsetNonce = wallet.Nonce + opts.OffsetNonce
				}
				nonce = wallet.OffsetNonce
			}
			if opts.OffsetNonce == 0 {
				wallet.Nonce++
			}
			wallet.OffsetNonce++
		}
		wallet.Unlock()

		signFn := func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
			sig := types.LatestSignerForChainID(c.ethClient.ChainID)
			return types.SignTx(tx, sig, wallet.PrivateKey)
		}

		var (
			tx  *types.Transaction
			err error
			t   time.Time
		)
		tops := bind.TransactOpts{
			From:    wallet.Address,
			Signer:  signFn,
			Context: vu.Context(),
			Nonce:   big.NewInt(int64(nonce)),
			Value:   big.NewInt(0),
		}

		if !c.isLegacy {
			head, err := c.ethClient.Ec.HeaderByNumber(vu.Context(), nil)
			if err != nil {
				errCh <- err
				return
			}

			tipCap := new(big.Int).Mul(c.latestGasTip.Load(), big.NewInt(int64(opts.GasPriceMultiplier)))
			feeCap := new(big.Int).Add(
				tipCap,
				new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
			)
			tops.GasFeeCap = feeCap
			tops.GasTipCap = tipCap
		} else {
			tops.GasPrice = new(big.Int).Mul(c.latestGasPrice.Load(), big.NewInt(int64(opts.GasPriceMultiplier)))
		}

		if opts.NoSend {
			t = time.Now()
			tops.NoSend = true
			tx, err = c.erc721.Contract.MintBatch(&tops, *targetAddress, big.NewInt(1))
			if err != nil {
				errCh <- err
				return
			}
			cm := ethereum.CallMsg{
				From:       wallet.Address,
				To:         tx.To(),
				Gas:        tx.Gas(),
				Value:      tx.Value(),
				Data:       tx.Data(),
				AccessList: tx.AccessList(),
			}
			if !c.isLegacy {
				cm.GasFeeCap = tx.GasFeeCap()
				cm.GasTipCap = tx.GasTipCap()
			} else {
				cm.GasPrice = tx.GasPrice()
			}
			if _, err := c.ethClient.Ec.CallContract(vu.Context(), cm, nil); err != nil {
				if strings.Contains(err.Error(), "fee cap less than block base fee") {
					retryForNonce = true
				}
				errCh <- err
				return
			}
		} else {
			t = time.Now()
			tx, err = c.erc721.Contract.MintBatch(&tops, *targetAddress, big.NewInt(1))
			if err != nil {
				retryForNonce = true
				if strings.Contains(err.Error(), "replacement transaction underpriced") && retryForNonce {
					retryForNonce = false
				}
				if strings.Contains(err.Error(), "transaction underpriced") && retryForNonce {
					retryForNonce = false
				}
				if strings.Contains(err.Error(), "nonce too low") && retryForNonce {
					retryForNonce = false
				}
				if strings.Contains(err.Error(), "already known") && retryForNonce {
					retryForNonce = false
				}
				if strings.Contains(err.Error(), "could not replace existing") && retryForNonce {
					retryForNonce = false
				}
				if strings.Contains(err.Error(), "fee cap less than block base fee") {
					retryForNonce = true
				}

				errCh <- err
				return
			}
		}

		ReportEoaFromStats(vu, metrics, c.uid, 1, eth.TransactionTypeERC721)
		ReportReqDurationFromStats(vu, metrics, c.uid, "sendERC721Transaction", time.Since(t))

		hash = tx.Hash()
		c.storeTransactionStartTime(hash, t)

		if opts.WaitForConfirmation {
			_, err := eth.WaitUntilMined(vu.Context(), c.ethClient.Ec, hash, opts.ConfirmationDelay, 10*time.Millisecond)
			if err != nil {
				errCh <- err
				return
			}

			ReportReqDurationFromStats(vu, metrics, c.uid, "sendConfirmedERC721Transaction", time.Since(t))
		}
	}

	for i := 0; i < int(opts.TxCount); i++ {
		if ctx.Err() != nil {
			break
		}

		sendTx()
	}

	close(errCh)
	for err := range errCh {
		if err != nil && !errors.Is(err, context.Canceled) {
			c.log.Error(err, "error in sending transactions")
		}
	}

	if hash == (common.Hash{}) {
		return nil, fmt.Errorf("failed to send any transaction")
	}

	return &hash, nil
}

func (c *DefaultClient) DeployContract(vu modules.VU, params *eth.DeployContractParams) (*common.Address, *common.Hash, *common.Address, error) {
	// Read ABI and binary files
	abiBytes, err := os.ReadFile(params.AbiPath)
	if err != nil {
		return nil, nil, nil, err
	}

	binBytes, err := os.ReadFile(params.BinPath)
	if err != nil {
		return nil, nil, nil, err
	}

	wallet := c.testers.GetAvailableWallet()
	if wallet == nil {
		return nil, nil, nil, err
	}
	defer c.testers.Unlock(wallet.Address)

	tops, err := bind.NewKeyedTransactorWithChainID(wallet.PrivateKey, c.ethClient.ChainID)
	if err != nil {
		return nil, nil, nil, err
	}
	tops.Nonce = big.NewInt(int64(wallet.Nonce))

	if !c.isLegacy {
		head, err := c.ethClient.Ec.HeaderByNumber(vu.Context(), nil)
		if err != nil {
			return nil, nil, nil, err
		}

		tipCap := new(big.Int).Mul(c.latestGasTip.Load(), big.NewInt(int64(params.GasPriceMultiplier)))
		feeCap := new(big.Int).Add(
			tipCap,
			new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
		)
		tops.GasFeeCap = feeCap
		tops.GasTipCap = tipCap
	} else {
		tops.GasPrice = new(big.Int).Mul(c.latestGasPrice.Load(), big.NewInt(int64(params.GasPriceMultiplier)))
	}

	parsedAbi, err := abi.JSON(strings.NewReader(string(abiBytes)))
	if err != nil {
		return nil, nil, nil, err
	}

	addr, tx, contract, err := bind.DeployContract(
		tops,
		parsedAbi,
		common.FromHex(string(binBytes)),
		c.ethClient.Ec,
		params.Args...,
	)
	if err != nil {
		if eth.DoIncreaseNonceWhenError(err) {
			wallet.Nonce++
			wallet.OffsetNonce++
		}
		return nil, nil, nil, err
	}

	hash := tx.Hash()
	wallet.Nonce++
	wallet.OffsetNonce++

	if _, err := eth.WaitUntilMined(vu.Context(), c.ethClient.Ec, hash, eth.EthDefaultBlockTime, 10*time.Millisecond); err != nil {
		return nil, nil, nil, err
	}

	c.deployedContracts[addr] = contract

	return &wallet.Address, &hash, &addr, nil
}

func (c *DefaultClient) TxContract(vu modules.VU, metrics *EthMetrics, params *eth.FnContractParams) (*common.Hash, error) {
	contract, ok := c.deployedContracts[params.ContractAddress]
	if !ok {
		return nil, fmt.Errorf("contract not found")
	}

	wallet := c.testers.GetAvailableWallet()
	if wallet == nil {
		return nil, fmt.Errorf("no available wallet")
	}
	defer c.testers.Unlock(wallet.Address)

	signFn := func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
		sig := types.LatestSignerForChainID(c.ethClient.ChainID)
		return types.SignTx(tx, sig, wallet.PrivateKey)
	}
	tops := &bind.TransactOpts{
		From:       wallet.Address,
		Signer:     signFn,
		Context:    vu.Context(),
		Nonce:      big.NewInt(int64(wallet.Nonce)),
		AccessList: params.AccessList,
	}

	if !c.isLegacy {
		head, err := c.ethClient.Ec.HeaderByNumber(vu.Context(), nil)
		if err != nil {
			return nil, err
		}

		tipCap := new(big.Int).Mul(c.latestGasTip.Load(), big.NewInt(int64(params.GasPriceMultiplier)))
		feeCap := new(big.Int).Add(
			tipCap,
			new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
		)
		tops.GasFeeCap = feeCap
		tops.GasTipCap = tipCap
	} else {
		tops.GasPrice = new(big.Int).Mul(c.latestGasPrice.Load(), big.NewInt(int64(params.GasPriceMultiplier)))
	}

	t := time.Now()
	tx, err := contract.Transact(tops, params.Method, params.Args...)
	if err != nil {
		if eth.DoIncreaseNonceWhenError(err) {
			wallet.Nonce++
			wallet.OffsetNonce++
		}
		return nil, err
	}
	wallet.Nonce++
	wallet.OffsetNonce++

	ReportReqDurationFromStats(vu, metrics, c.uid, "callContract", time.Since(t))

	hash := tx.Hash()
	return &hash, nil
}

func (c *DefaultClient) CallContract(vu modules.VU, metrics *EthMetrics, params *eth.FnContractParams) ([]interface{}, error) {
	contract, ok := c.deployedContracts[params.ContractAddress]
	if !ok {
		return nil, fmt.Errorf("contract not found")
	}

	var out []interface{}
	t := time.Now()
	if err := contract.Call(&bind.CallOpts{
		From:    params.ContractAddress,
		Context: vu.Context(),
	}, &out, params.Method, params.Args...); err != nil {
		return nil, err
	}
	ReportReqDurationFromStats(vu, metrics, c.uid, "callContract", time.Since(t))
	return out, nil
}
