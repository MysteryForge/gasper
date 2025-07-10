package eth

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"golang.org/x/sync/semaphore"
)

type Wallet struct {
	Address     common.Address
	PrivateKey  *ecdsa.PrivateKey
	Nonce       uint64
	OffsetNonce uint64
	mu          sync.RWMutex
}

func NewWallet() (*Wallet, error) {
	pk, err := crypto.GenerateKey()
	if err != nil {
		return nil, err
	}
	addr := crypto.PubkeyToAddress(pk.PublicKey)

	return &Wallet{Address: addr, PrivateKey: pk}, nil
}

func (w *Wallet) Lock() {
	w.mu.Lock()
}

func (w *Wallet) Unlock() {
	w.mu.Unlock()
}

func (w *Wallet) IncNonce() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Nonce++
}

func (w *Wallet) IncOffsetNonce() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.OffsetNonce++
}

func (w *Wallet) IncreaseBothNonces() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.Nonce++
	w.OffsetNonce++
}

func (w *Wallet) RefreshNonce(ec *ethclient.Client) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	nonce, err := ec.PendingNonceAt(context.Background(), w.Address)
	if err != nil {
		return fmt.Errorf("failed to get nonce for wallet %s: %w", w.Address, err)
	}
	w.Nonce = nonce
	w.OffsetNonce = nonce
	return nil
}

type WalletRegistry struct {
	wallets map[common.Address]*Wallet
	locks   map[common.Address]bool
	ln      int
	mu      sync.RWMutex
}

func NewEmptyWalletRegistry() *WalletRegistry {
	return &WalletRegistry{
		wallets: make(map[common.Address]*Wallet),
		locks:   make(map[common.Address]bool),
	}
}

func NewWalletRegistryFromPrivateKeys(ctx context.Context, ec *ethclient.Client, privateKeys []string) (*WalletRegistry, error) {
	if len(privateKeys) == 0 {
		return nil, fmt.Errorf("private keys is empty")
	}

	sem := semaphore.NewWeighted(MaxNumberOfCreatingWalletsAtOnce)
	errCh := make(chan error, len(privateKeys))

	wr := NewEmptyWalletRegistry()

	for _, pk := range privateKeys {
		if err := sem.Acquire(ctx, 1); err != nil {
			return nil, err
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case err := <-errCh:
			return nil, err
		default:
		}

		go func() {
			defer sem.Release(1)

			wallet, err := ParsePrivateKey(pk)
			if err != nil {
				errCh <- err
				return
			}
			if err := wallet.RefreshNonce(ec); err != nil {
				errCh <- fmt.Errorf("failed to refresh nonce for wallet %s: %w", wallet.Address, err)
				return
			}
			wr.Register(wallet)
		}()
	}

	if err := sem.Acquire(ctx, MaxNumberOfCreatingWalletsAtOnce); err != nil {
		return nil, err
	}

	close(errCh)
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	return wr, nil
}

func (wr *WalletRegistry) GenerateAndStore(ctx context.Context, numWallets uint64) error {
	sem := semaphore.NewWeighted(MaxNumberOfCreatingWalletsAtOnce)
	errCh := make(chan error, int(numWallets))

	for i := 0; i < int(numWallets); i++ {
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
				// return fmt.Errorf("failed to create wallet: %w", err)
				errCh <- fmt.Errorf("failed to create wallet: %w", err)
				return
			}
			wr.Register(wallet)
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

func (wr *WalletRegistry) All() map[common.Address]*Wallet {
	wr.mu.RLock()
	defer wr.mu.RUnlock()
	return wr.wallets
}

func (wr *WalletRegistry) Register(wallet *Wallet) {
	wr.mu.Lock()
	defer wr.mu.Unlock()
	wr.ln++
	wr.wallets[wallet.Address] = wallet
	wr.locks[wallet.Address] = false
}

func (wr *WalletRegistry) Unregister(addr common.Address) {
	wr.mu.Lock()
	defer wr.mu.Unlock()
	wr.ln--
	delete(wr.wallets, addr)
	delete(wr.locks, addr)
}

func (wr *WalletRegistry) Lock(addr common.Address) bool {
	wr.mu.Lock()
	defer wr.mu.Unlock()
	_, exists := wr.locks[addr]
	if !exists {
		return false
	}
	if !wr.locks[addr] {
		wr.locks[addr] = true
		return true
	}
	return false
}

func (wr *WalletRegistry) Unlock(addr common.Address) bool {
	wr.mu.Lock()
	defer wr.mu.Unlock()
	_, exists := wr.locks[addr]
	if !exists {
		return false
	}
	wr.locks[addr] = false
	return true
}

func (wr *WalletRegistry) GetAvailableWallet() *Wallet {
	wr.mu.Lock()
	defer wr.mu.Unlock()

	for addr, wallet := range wr.wallets {
		if !wr.locks[addr] {
			wr.locks[addr] = true
			return wallet
		}
	}
	return nil
}

func (wr *WalletRegistry) IsLocked(addr common.Address) bool {
	wr.mu.RLock()
	defer wr.mu.RUnlock()
	return wr.locks[addr]
}

func (wr *WalletRegistry) Ln() int {
	wr.mu.RLock()
	defer wr.mu.RUnlock()
	return wr.ln
}
