package eth

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/time/rate"
)

type PoolStatus struct {
	BaseFee hexutil.Uint64 `json:"baseFee"`
	Pending hexutil.Uint64 `json:"pending"`
	Queued  hexutil.Uint64 `json:"queued"`
}

type TxPoolRateLimiter struct {
	steadyStateTxPoolSize      uint64
	adaptiveRateLimitIncrement uint64
	mu                         *sync.Mutex
	limiter                    *rate.Limiter
	backoffFactor              float64
}

func NewTxPoolRateLimiter(initialRate uint64) *TxPoolRateLimiter {
	limiter := rate.NewLimiter(rate.Limit(initialRate), 1)

	return &TxPoolRateLimiter{
		steadyStateTxPoolSize:      1000,
		adaptiveRateLimitIncrement: 50,
		limiter:                    limiter,
		backoffFactor:              2.0, // Default backoff factor
		mu:                         &sync.Mutex{},
	}
}

func (rl *TxPoolRateLimiter) UpdateTxPoolSize(size uint64) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	currentLimit := rl.limiter.Limit()

	if size < rl.steadyStateTxPoolSize {
		newLimit := currentLimit + rate.Limit(rl.adaptiveRateLimitIncrement)
		rl.limiter.SetLimit(newLimit)
	} else if size > rl.steadyStateTxPoolSize {
		newLimit := currentLimit / rate.Limit(rl.backoffFactor)
		if newLimit < 1 {
			newLimit = 1
		}
		rl.limiter.SetLimit(newLimit)
	}
}

func (rl *TxPoolRateLimiter) Wait(ctx context.Context) error {
	return rl.limiter.Wait(ctx)
}
