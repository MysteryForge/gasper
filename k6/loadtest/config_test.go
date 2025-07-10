package loadtest

import (
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadConfigYML(t *testing.T) {
	t.Run("valid single config with all fields", func(t *testing.T) {
		tmpFile := createTempConfigFile(t, `
http: http://localhost:8123
num_wallets: 5
fund_amount: 90000000000000000000
private_keys: [0x52fb3ff54731f7609d97b6b0195aa1fac56b95141c4b71eaa4f08af23558c63b]
target_addresses: [0xc78260046895c358dE4bE97210Efca3900544905]
erc20: true
erc20_address: 0x11730920CC1DFa1dbBC2436B4E360662c50db10d
erc20_mint_amount: 100000000000000000000
erc721: true
erc721_address: 0x22730920CC1DFa1dbBC2436B4E360662c50db10d
erc721_mint: true
`)
		defer os.Remove(tmpFile) // nolint: errcheck

		configs, err := ReadConfigYML(tmpFile)
		require.NoError(t, err)
		require.Len(t, configs, 1)

		cfg := configs[0]
		assert.Equal(t, "http://localhost:8123", cfg.HTTP)
		assert.Equal(t, uint64(5), cfg.NumWallets)
		assert.Equal(t, []string{"0x52fb3ff54731f7609d97b6b0195aa1fac56b95141c4b71eaa4f08af23558c63b"}, cfg.PrivateKeys)
		assert.Equal(t, []string{"0xc78260046895c358dE4bE97210Efca3900544905"}, cfg.TargetAddresses)

		expectedAmount := new(big.Int)
		expectedAmount.SetString("90000000000000000000", 10)
		assert.Equal(t, expectedAmount, &cfg.FundAmount)

		assert.True(t, cfg.ERC20)
		assert.Equal(t, "0x11730920CC1DFa1dbBC2436B4E360662c50db10d", cfg.ERC20Address)
		expectedERC20Amount := new(big.Int)
		expectedERC20Amount.SetString("100000000000000000000", 10)
		assert.Equal(t, expectedERC20Amount, &cfg.ERC20MintAmount)

		assert.True(t, cfg.ERC721)
		assert.Equal(t, "0x22730920CC1DFa1dbBC2436B4E360662c50db10d", cfg.ERC721Address)
		assert.True(t, cfg.ERC721Mint)
	})

	t.Run("valid multiple configs", func(t *testing.T) {
		tmpFile := createTempConfigFile(t, `
- http: http://localhost:8123
  num_wallets: 5
  fund_amount: 90000000000000000000
  private_keys: [0x52fb3ff54731f7609d97b6b0195aa1fac56b95141c4b71eaa4f08af23558c63b]
  target_addresses: [0xc78260046895c358dE4bE97210Efca3900544905]
  erc20: true
  erc20_address: 0x11730920CC1DFa1dbBC2436B4E360662c50db10d
  erc20_mint_amount: 100000000000000000000
- http: http://localhost:8124
  wallets: ["wallet1", "wallet2"]
  target_addresses: [0xd78260046895c358dE4bE97210Efca3900544906]
  erc721: true
  erc721_address: 0x22730920CC1DFa1dbBC2436B4E360662c50db10d

`)
		defer os.Remove(tmpFile) // nolint: errcheck

		configs, err := ReadConfigYML(tmpFile)
		require.NoError(t, err)
		require.Len(t, configs, 2)

		assert.Equal(t, "http://localhost:8123", configs[0].HTTP)
		assert.Equal(t, "http://localhost:8124", configs[1].HTTP)
		assert.Equal(t, []string{"wallet1", "wallet2"}, configs[1].Wallets)
		assert.Equal(t, ".scratch/load_test.db", configs[0].DBPath)

		assert.True(t, configs[0].ERC20)
		assert.True(t, configs[1].ERC721)
	})

	t.Run("invalid file path", func(t *testing.T) {
		_, err := ReadConfigYML("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no config file provided")

		_, err = ReadConfigYML("nonexistent.yml")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read config file")
	})

	t.Run("invalid file extension", func(t *testing.T) {
		tmpFile := createTempConfigFile(t, "invalid")
		defer os.Remove(tmpFile) // nolint: errcheck

		newPath := tmpFile[:len(tmpFile)-4] + ".txt"
		err := os.Rename(tmpFile, newPath)
		require.NoError(t, err)
		defer os.Remove(newPath) // nolint: errcheck

		_, err = ReadConfigYML(newPath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid config file extension")
	})

	t.Run("invalid yaml syntax", func(t *testing.T) {
		tmpFile := createTempConfigFile(t, `
http: http://localhost:8123
  invalid:
    yaml: syntax
`)
		defer os.Remove(tmpFile) // nolint: errcheck

		_, err := ReadConfigYML(tmpFile)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse config file")
	})
}

func TestValidateClientConfig(t *testing.T) {
	t.Run("valid config with all fields", func(t *testing.T) {
		cfg := &clientConfig{
			HTTP:            "http://localhost:8123",
			PrivateKeys:     []string{"0x52fb3ff54731f7609d97b6b0195aa1fac56b95141c4b71eaa4f08af23558c63b"},
			NumWallets:      5,
			FundAmount:      *big.NewInt(1000000),
			TargetAddresses: []string{"0xc78260046895c358dE4bE97210Efca3900544905"},
			ERC20:           true,
			ERC20Address:    "0x11730920CC1DFa1dbBC2436B4E360662c50db10d",
			ERC20MintAmount: *big.NewInt(1000000),
			ERC721:          true,
			ERC721Address:   "0x22730920CC1DFa1dbBC2436B4E360662c50db10d",
			ERC721Mint:      true,
		}
		err := validateClientConfig(cfg)
		assert.NoError(t, err)
	})

	t.Run("invalid ERC20 config", func(t *testing.T) {
		cfg := &clientConfig{
			HTTP:            "http://localhost:8123",
			ERC20:           true,
			Wallets:         []string{"wallet1"},
			TargetAddresses: []string{"0xc78260046895c358dE4bE97210Efca3900544905"},
		}
		err := validateClientConfig(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erc20_address is required")
	})

	t.Run("invalid ERC721 config", func(t *testing.T) {
		cfg := &clientConfig{
			HTTP:            "http://localhost:8123",
			ERC721:          true,
			Wallets:         []string{"wallet1"},
			TargetAddresses: []string{"0xc78260046895c358dE4bE97210Efca3900544905"},
		}
		err := validateClientConfig(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erc721_address is required")
	})

	t.Run("invalid ERC20 mint amount", func(t *testing.T) {
		cfg := &clientConfig{
			HTTP:            "http://localhost:8123",
			ERC20:           true,
			ERC20Address:    "0x11730920CC1DFa1dbBC2436B4E360662c50db10d",
			ERC20MintAmount: *big.NewInt(100),
		}
		err := validateClientConfig(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot transfer to erc20_contracts")
	})

	t.Run("invalid HTTP URL", func(t *testing.T) {
		cfg := &clientConfig{
			HTTP: "invalid-url",
		}
		err := validateClientConfig(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid Http URL scheme")
	})

	t.Run("invalid HTTP scheme", func(t *testing.T) {
		cfg := &clientConfig{
			HTTP: "ws://localhost:8123",
		}
		err := validateClientConfig(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid Http URL scheme")
	})

	t.Run("missing private key with num_wallets", func(t *testing.T) {
		cfg := &clientConfig{
			HTTP:       "http://localhost:8123",
			NumWallets: 5,
		}
		err := validateClientConfig(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "private keys is required when num_wallets > 0")
	})

	t.Run("invalid fund amount", func(t *testing.T) {
		cfg := &clientConfig{
			HTTP:        "http://localhost:8123",
			PrivateKeys: []string{"0x52fb3ff54731f7609d97b6b0195aa1fac56b95141c4b71eaa4f08af23558c63b"},
			NumWallets:  5,
			FundAmount:  *big.NewInt(0),
		}
		err := validateClientConfig(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "fund amount should be greater then 0")
	})

	t.Run("missing target address", func(t *testing.T) {
		cfg := &clientConfig{
			HTTP:        "http://localhost:8123",
			NumWallets:  5,
			PrivateKeys: []string{"0x52fb3ff54731f7609d97b6b0195aa1fac56b95141c4b71eaa4f08af23558c63b"},
			FundAmount:  *big.NewInt(1000000),
		}
		err := validateClientConfig(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "target address is required")
	})
}

// Helper function to create temporary config file
func createTempConfigFile(t *testing.T, content string) string {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yml")
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	require.NoError(t, err)
	return tmpFile
}
