package eth

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func ParsePrivateKey(privateKey string) (*Wallet, error) {
	pk, err := crypto.HexToECDSA(strings.TrimPrefix(privateKey, "0x"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	addr := crypto.PubkeyToAddress(pk.PublicKey)
	return &Wallet{Address: addr, PrivateKey: pk}, nil
}

func InitEthClient(ctx context.Context, uri string) (*rpc.Client, *ethclient.Client, *big.Int, error) {
	rc, err := rpc.DialContext(ctx, uri)
	if err != nil {
		return nil, nil, nil, err
	}
	ec := ethclient.NewClient(rc)

	sync, err := ec.SyncProgress(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to get sync progress: %w", err)
	}
	if sync != nil {
		return nil, nil, nil, fmt.Errorf("node is syncing")
	}

	chainID, err := ec.ChainID(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	return rc, ec, chainID, nil
}

func SanitizeURL(uri string) string {
	u, err := url.Parse(uri)
	if err != nil {
		return fmt.Sprintf("invalid URL: %s", err)
	}

	if u.User == nil {
		return uri
	}

	u.User = url.User("***")
	sanitized, err := url.QueryUnescape(u.String())
	if err != nil {
		return u.String()
	}
	return sanitized
}

func RepeatWithTimeout(ctx context.Context, timeout time.Duration, step time.Duration, fn func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	ticker := time.NewTicker(step)
	defer ticker.Stop()

	for ctx.Err() == nil {
		err := fn(ctx)
		if err == nil {
			return nil
		}

		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return err
		}
		<-ticker.C
	}
	return ctx.Err()
}

func WaitUntilMined(ctx context.Context, ec *ethclient.Client, hash common.Hash, timeout time.Duration, step time.Duration) (*types.Receipt, error) {
	var receipt *types.Receipt
	if err := RepeatWithTimeout(ctx, timeout, step, func(ctx context.Context) error {
		r, err := ec.TransactionReceipt(ctx, hash)
		if err != nil {
			return err
		}
		if r != nil && r.Status == types.ReceiptStatusSuccessful {
			receipt = r
			return nil
		}
		return fmt.Errorf("receipt status in unsuccessful")
	}); err != nil {
		return nil, fmt.Errorf("%w for hash %s", err, hash.Hex())
	}
	return receipt, nil
}

func DoIncreaseNonceWhenError(err error) bool {
	if strings.Contains(err.Error(), "replacement transaction underpriced") {
		return true
	}
	if strings.Contains(err.Error(), "transaction underpriced") {
		return true
	}
	if strings.Contains(err.Error(), "nonce too low") {
		return true
	}
	if strings.Contains(err.Error(), "already known") {
		return true
	}
	if strings.Contains(err.Error(), "could not replace existing") {
		return true
	}

	return false
}
