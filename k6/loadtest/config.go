package loadtest

import (
	"fmt"
	"math/big"
	"net/url"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	MaxWalletsNumContract = 1000
)

type clientConfig struct {
	HTTP string `yaml:"http" js:"http"`

	PrivateKeys        []string `yaml:"private_keys,omitempty" js:"privateKeys,omitempty"`                // private key of the account which will be used to fund new wallets
	BatchFunderAddress string   `yaml:"batch_funder_address,omitempty" js:"batchFunderAddress,omitempty"` // address of the batch funder contract

	NumWallets uint64   `yaml:"num_wallets,omitempty" js:"numWallets,omitempty"` // number of new wallets to create and fund
	FundAmount big.Int  `yaml:"fund_amount,omitempty" js:"fundAmount,omitempty"` // amount to fund new wallets with
	Wallets    []string `yaml:"wallets,omitempty" js:"wallets,omitempty"`        // prefunded wallets

	TargetAddresses    []string `yaml:"target_addresses,omitempty" js:"targetAddresses,omitempty"`        // target address for transactions
	NumTargetAddresses uint64   `yaml:"num_target_addresses,omitempty" js:"numTargetAddresses,omitempty"` // number of new target addresses to use`

	ERC20           bool    `yaml:"erc20,omitempty" js:"erc20,omitempty"`                       // whether to test ERC20 token contract
	ERC20Address    string  `yaml:"erc20_address,omitempty" js:"erc20Address,omitempty"`        // address of the ERC20 token contract
	ERC20MintAmount big.Int `yaml:"erc20_mint_amount,omitempty" js:"erc20MintAmount,omitempty"` // amount to transfer in ERC20 token contract per wallet

	ERC721        bool   `yaml:"erc721,omitempty" js:"erc721,omitempty"`                // whether to test ERC721 token contract
	ERC721Address string `yaml:"erc721_address,omitempty" js:"erc721Address,omitempty"` // address of the ERC721 token contract
	ERC721Mint    bool   `yaml:"erc721_mint,omitempty" js:"erc721Mint,omitempty"`       // whether to mint ERC721 token contract on startup

	RateLimite        *uint64 `yaml:"rate_limit,omitempty" js:"rateLimit,omitempty"`                  // rate limit for transaction sending in tx/s
	AdaptiveRateLimit bool    `yaml:"adaptive_rate_limit,omitempty" js:"adaptiveRateLimit,omitempty"` // whether to use adaptive rate limiting, using tx pool size

	MinGasPrice uint64 `yaml:"min_gas_price,omitempty" js:"minGasPrice,omitempty"` // minimum gas price to use for transactions

	DBPath string `yaml:"db_path,omitempty" js:"dbPath,omitempty"` // path to the database where we store transaction hashes
}

func ReadConfigYML(pth string) ([]*clientConfig, error) {
	if pth == "" {
		return nil, fmt.Errorf("no config file provided")
	}
	if filepath.Ext(pth) != ".yml" && filepath.Ext(pth) != ".yaml" {
		return nil, fmt.Errorf("invalid config file extension: %s", filepath.Ext(pth))
	}
	pth = filepath.Clean(pth)
	d, err := os.ReadFile(pth)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfgs := make([]*clientConfig, 0)
	if err := yaml.Unmarshal(d, &cfgs); err != nil {
		var cfg clientConfig
		if err := yaml.Unmarshal(d, &cfg); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}

		cfgs = append(cfgs, &cfg)
		return cfgs, nil
	}

	for i, cfg := range cfgs {
		if err := validateClientConfig(cfg); err != nil {
			return nil, fmt.Errorf("failed to validate config %d: %w", i, err)
		}
		if cfg.DBPath == "" {
			cfg.DBPath = ".scratch/load_test.db"
		}
		if cfg.RateLimite == nil {
			cfg.RateLimite = new(uint64)
			*cfg.RateLimite = 4
		}
	}

	return cfgs, nil
}

func validateClientConfig(cfg *clientConfig) error {
	uri, err := url.Parse(cfg.HTTP)
	if err != nil {
		return fmt.Errorf("failed to parse Http URL: %w", err)
	}

	if uri.Scheme != "http" && uri.Scheme != "https" {
		return fmt.Errorf("invalid Http URL scheme: %s", uri.Scheme)
	}

	if cfg.NumWallets > 0 && len(cfg.PrivateKeys) == 0 {
		return fmt.Errorf("private keys is required when num_wallets > 0")
	}

	if cfg.NumWallets > 0 && (cfg.FundAmount.Cmp(big.NewInt(0)) == 0 || cfg.FundAmount.Cmp(big.NewInt(0)) < 0) {
		return fmt.Errorf("fund amount should be greater then 0: %v when num_wallets > 0", cfg.FundAmount)
	}

	numWallets := len(cfg.Wallets) + int(cfg.NumWallets)
	if numWallets > 0 && len(cfg.TargetAddresses) == 0 && cfg.NumTargetAddresses == 0 {
		return fmt.Errorf("target address is required when num_wallets > 0 or wallets are provided")
	}

	if cfg.ERC20 && cfg.ERC20Address == "" && len(cfg.PrivateKeys) == 0 {
		return fmt.Errorf("erc20_address is required when erc20_test is true and private_keys is not set")
	}

	if cfg.ERC20 && numWallets == 0 && cfg.ERC20MintAmount.Cmp(big.NewInt(0)) == 1 {
		return fmt.Errorf("cannot transfer to erc20_contracts the erc20_mint_amount when num_wallets is 0")
	}

	if cfg.ERC20 && numWallets > 0 && cfg.ERC20MintAmount.Cmp(big.NewInt(0)) == 0 {
		return fmt.Errorf("erc20_mint_amount should be greater then 0 when num_wallets > 0")
	}

	if cfg.ERC721 && cfg.ERC721Address == "" && len(cfg.PrivateKeys) == 0 {
		return fmt.Errorf("erc721_address is required when erc721_test is true and private_keys is not set")
	}

	if cfg.ERC20 && numWallets > MaxWalletsNumContract {
		return fmt.Errorf("num_wallets should be less then %d when erc20_test is true", MaxWalletsNumContract)
	}

	if cfg.ERC721 && numWallets == 0 && cfg.ERC721Mint {
		return fmt.Errorf("cannot mint on erc721_contracts when num_wallets is 0")
	}

	if cfg.ERC721 && numWallets > MaxWalletsNumContract {
		return fmt.Errorf("num_wallets should be less then %d when erc721_test is true", MaxWalletsNumContract)
	}

	if cfg.DBPath != "" && filepath.Clean(cfg.DBPath) != cfg.DBPath && filepath.Ext(cfg.DBPath) != ".db" {
		return fmt.Errorf("invalid db path: %s", cfg.DBPath)
	}

	return nil
}
