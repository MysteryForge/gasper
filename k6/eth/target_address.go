package eth

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/sync/semaphore"
)

type TargetAddresses struct {
	list []*common.Address
	ln   int
	mu   *sync.RWMutex
}

func NewTargetAddresses(addresses []string, numAddresses uint64) (*TargetAddresses, error) {
	ln := len(addresses)
	targetAddresses := make([]*common.Address, 0, ln)
	for _, address := range addresses {
		targetAddress := common.HexToAddress(address)
		targetAddresses = append(targetAddresses, &targetAddress)
	}

	ta := &TargetAddresses{list: targetAddresses, ln: ln, mu: &sync.RWMutex{}}
	if err := ta.GenerateAndStore(context.Background(), numAddresses); err != nil {
		return nil, err
	}

	return ta, nil
}

func (ta *TargetAddresses) GenerateAndStore(ctx context.Context, numAddresses uint64) error {
	sem := semaphore.NewWeighted(MaxNumberOfCreatingWalletsAtOnce)
	errCh := make(chan error, int(numAddresses))

	for i := 0; i < int(numAddresses); i++ {
		if err := sem.Acquire(ctx, 1); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return err
		default:
		}

		go func() {
			defer sem.Release(1)

			wallet, err := NewWallet()
			if err != nil {
				errCh <- fmt.Errorf("failed to create wallet: %w", err)
				return
			}
			ta.Add(&wallet.Address)
		}()
	}

	if err := sem.Acquire(ctx, MaxNumberOfCreatingWalletsAtOnce); err != nil {
		return err
	}

	close(errCh)
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

func (ta *TargetAddresses) Add(address *common.Address) {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	ta.list = append(ta.list, address)
	ta.ln++
}

func (ta *TargetAddresses) Random() *common.Address {
	ta.mu.RLock()
	defer ta.mu.RUnlock()

	if len(ta.list) == 0 {
		return nil
	}
	randIndex := rand.Intn(ta.ln)
	return ta.list[randIndex]
}

func (ta *TargetAddresses) All() []*common.Address {
	ta.mu.RLock()
	defer ta.mu.RUnlock()

	return ta.list
}
