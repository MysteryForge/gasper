package eth

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestWallet(t *testing.T) *Wallet {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	return &Wallet{
		Address:    address,
		PrivateKey: privateKey,
		Nonce:      0,
	}
}

func TestNewEmptyWalletRegistry(t *testing.T) {
	registry := NewEmptyWalletRegistry()
	assert.NotNil(t, registry.wallets)
	assert.NotNil(t, registry.locks)
	assert.Empty(t, registry.wallets)
	assert.Empty(t, registry.locks)
}

func TestWalletRegistry_Register(t *testing.T) {
	registry := NewEmptyWalletRegistry()
	wallet := createTestWallet(t)

	registry.Register(wallet)

	assert.Len(t, registry.wallets, 1)
	assert.Len(t, registry.locks, 1)
	assert.Equal(t, wallet, registry.wallets[wallet.Address])
	assert.False(t, registry.locks[wallet.Address])
}

func TestWalletRegistry_Unregister(t *testing.T) {
	registry := NewEmptyWalletRegistry()
	wallet := createTestWallet(t)

	registry.Register(wallet)
	assert.Len(t, registry.wallets, 1)

	registry.Unregister(wallet.Address)
	assert.Empty(t, registry.wallets)
	assert.Empty(t, registry.locks)
}

func TestWalletRegistry_Lock(t *testing.T) {
	registry := NewEmptyWalletRegistry()
	wallet := createTestWallet(t)
	registry.Register(wallet)

	// First lock should succeed
	success := registry.Lock(wallet.Address)
	assert.True(t, success)
	assert.True(t, registry.locks[wallet.Address])

	// Second lock should fail
	success = registry.Lock(wallet.Address)
	assert.False(t, success)
	assert.True(t, registry.locks[wallet.Address])

	// Lock for non-existent wallet should return false
	nonExistentAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	success = registry.Lock(nonExistentAddr)
	assert.False(t, success)
}

func TestWalletRegistry_Unlock(t *testing.T) {
	registry := NewEmptyWalletRegistry()
	wallet := createTestWallet(t)
	registry.Register(wallet)

	registry.Lock(wallet.Address)
	assert.True(t, registry.locks[wallet.Address])

	registry.Unlock(wallet.Address)
	assert.False(t, registry.locks[wallet.Address])

	// Unlocking an already unlocked wallet should not cause issues
	registry.Unlock(wallet.Address)
	assert.False(t, registry.locks[wallet.Address])

	// Unlocking non-existent wallet should not panic
	nonExistentAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	registry.Unlock(nonExistentAddr)
}

func TestWalletRegistry_GetAvailableWallet(t *testing.T) {
	registry := NewEmptyWalletRegistry()
	wallet1 := createTestWallet(t)
	wallet2 := createTestWallet(t)

	registry.Register(wallet1)
	registry.Register(wallet2)

	// First call should return an available wallet and lock it
	available := registry.GetAvailableWallet()
	assert.NotNil(t, available)
	assert.True(t, registry.locks[available.Address])

	// Second call should return the other wallet
	available2 := registry.GetAvailableWallet()
	assert.NotNil(t, available2)
	assert.NotEqual(t, available.Address, available2.Address)
	assert.True(t, registry.locks[available2.Address])

	// Third call should return nil as all wallets are locked
	available3 := registry.GetAvailableWallet()
	assert.Nil(t, available3)
}

func TestWalletRegistry_IsLocked(t *testing.T) {
	registry := NewEmptyWalletRegistry()
	wallet := createTestWallet(t)
	registry.Register(wallet)

	// Initially unlocked
	assert.False(t, registry.IsLocked(wallet.Address))

	// After locking
	registry.Lock(wallet.Address)
	assert.True(t, registry.IsLocked(wallet.Address))

	// After unlocking
	registry.Unlock(wallet.Address)
	assert.False(t, registry.IsLocked(wallet.Address))

	// Non-existent wallet
	nonExistentAddr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	assert.False(t, registry.IsLocked(nonExistentAddr))
}

func TestWalletRegistry_ConcurrentOperations(t *testing.T) {
	registry := NewEmptyWalletRegistry()
	wallet := createTestWallet(t)
	registry.Register(wallet)

	// Test concurrent access
	done := make(chan bool)
	go func() {
		for i := 0; i < 100; i++ {
			registry.Lock(wallet.Address)
			registry.Unlock(wallet.Address)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			registry.IsLocked(wallet.Address)
		}
		done <- true
	}()

	// Wait for both goroutines to finish
	<-done
	<-done
}
