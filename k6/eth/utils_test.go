package eth

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/mysteryforge/gasper/k6/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePrivateKey(t *testing.T) {
	t.Run("valid private key", func(t *testing.T) {
		privateKey := "52fb3ff54731f7609d97b6b0195aa1fac56b95141c4b71eaa4f08af23558c63b"
		wallet, err := ParsePrivateKey(privateKey)
		require.NoError(t, err)
		assert.NotNil(t, wallet.Address)
		assert.NotNil(t, wallet.PrivateKey)
	})

	t.Run("valid private key with 0x prefix", func(t *testing.T) {
		privateKey := "0x52fb3ff54731f7609d97b6b0195aa1fac56b95141c4b71eaa4f08af23558c63b"
		wallet, err := ParsePrivateKey(privateKey)
		require.NoError(t, err)
		assert.NotNil(t, wallet.Address)
		assert.NotNil(t, wallet.PrivateKey)
	})

	t.Run("invalid private key", func(t *testing.T) {
		privateKey := "invalid"
		wallet, err := ParsePrivateKey(privateKey)
		assert.Error(t, err)
		assert.Nil(t, wallet)
		assert.Contains(t, err.Error(), "failed to parse private key")
	})
}

func TestInitEthClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rpcURL, _ := mock.StartTestAnvilContainer(t, ctx)
	rc, ec, chainID, err := InitEthClient(ctx, rpcURL)
	require.NoError(t, err)
	defer rc.Close()
	defer ec.Close()

	assert.NotNil(t, rc)
	assert.NotNil(t, ec)
	assert.NotNil(t, chainID)
}

func TestSanitizeURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic url without credentials",
			input:    "http://localhost:8545",
			expected: "http://localhost:8545",
		},
		{
			name:     "url with username and password",
			input:    "http://user:pass@localhost:8545",
			expected: "http://***@localhost:8545",
		},
		{
			name:     "url with only username",
			input:    "http://user@localhost:8545",
			expected: "http://***@localhost:8545",
		},
		{
			name:     "https url with credentials",
			input:    "https://admin:secret@example.com",
			expected: "https://***@example.com",
		},
		{
			name:     "url with path and query",
			input:    "http://user:pass@example.com/path?query=value",
			expected: "http://***@example.com/path?query=value",
		},
		{
			name:     "websocket url with credentials",
			input:    "ws://user:pass@localhost:8546",
			expected: "ws://***@localhost:8546",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeURL(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestRepeatWithTimeout(t *testing.T) {
	t.Run("succeeds immediately", func(t *testing.T) {
		ctx := context.Background()
		calls := 0
		err := RepeatWithTimeout(ctx, time.Second, time.Millisecond*100, func(ctx context.Context) error {
			calls++
			return nil
		})
		require.NoError(t, err)
		require.Equal(t, 1, calls, "function should be called exactly once")
	})

	t.Run("succeeds after retries", func(t *testing.T) {
		ctx := context.Background()
		calls := 0
		err := RepeatWithTimeout(ctx, time.Second, time.Millisecond*10, func(ctx context.Context) error {
			calls++
			if calls < 3 {
				return fmt.Errorf("temporary error")
			}
			return nil
		})
		require.NoError(t, err)
		require.Equal(t, 3, calls, "function should be called exactly three times")
	})

	t.Run("fails on timeout", func(t *testing.T) {
		ctx := context.Background()
		calls := 0
		err := RepeatWithTimeout(ctx, time.Millisecond*50, time.Millisecond*10, func(ctx context.Context) error {
			calls++
			return fmt.Errorf("persistent error")
		})
		require.Error(t, err)
		require.True(t, errors.Is(err, context.DeadlineExceeded))
		require.Greater(t, calls, 1, "function should be called multiple times")
	})

	t.Run("respects parent context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		calls := 0

		go func() {
			time.Sleep(time.Millisecond * 50)
			cancel()
		}()

		err := RepeatWithTimeout(ctx, time.Second, time.Millisecond*10, func(ctx context.Context) error {
			calls++
			return fmt.Errorf("error")
		})
		require.Error(t, err)
		require.True(t, errors.Is(err, context.Canceled))
		require.Greater(t, calls, 1, "function should be called multiple times")
	})

	t.Run("handles function returning context errors", func(t *testing.T) {
		ctx := context.Background()
		calls := 0
		err := RepeatWithTimeout(ctx, time.Second, time.Millisecond*10, func(ctx context.Context) error {
			calls++
			return context.DeadlineExceeded
		})
		require.Error(t, err)
		require.True(t, errors.Is(err, context.DeadlineExceeded))
		require.Equal(t, 1, calls, "function should be called exactly once")
	})
}
