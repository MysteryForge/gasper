package eth

import (
	"math/big"
	"sync/atomic"
)

// AtomicBigInt provides atomic operations for big.Int values
type AtomicBigInt struct {
	value atomic.Value
}

// NewAtomicBigInt creates a new AtomicBigInt with an initial value
func NewAtomicBigInt(initial *big.Int) *AtomicBigInt {
	a := &AtomicBigInt{}
	if initial != nil {
		a.Store(new(big.Int).Set(initial))
	} else {
		a.Store(new(big.Int))
	}
	return a
}

// Load atomically loads the big.Int value
func (a *AtomicBigInt) Load() *big.Int {
	if v := a.value.Load(); v != nil {
		return v.(*big.Int)
	}
	return new(big.Int)
}

// Store atomically stores the big.Int value
func (a *AtomicBigInt) Store(val *big.Int) {
	// Always store a copy to prevent external modification
	if val != nil {
		a.value.Store(new(big.Int).Set(val))
	} else {
		a.value.Store(new(big.Int))
	}
}

// CompareAndSwap atomically swaps the value if the current value equals the expected value
func (a *AtomicBigInt) CompareAndSwap(expected, newValue *big.Int) bool {
	for {
		current := a.Load()
		if current.Cmp(expected) != 0 {
			return false
		}
		newCopy := new(big.Int).Set(newValue)
		if a.value.CompareAndSwap(current, newCopy) {
			return true
		}
	}
}
