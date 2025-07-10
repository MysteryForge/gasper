package eth

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mysteryforge/gasper/bindings"
	"golang.org/x/sync/semaphore"
)

type ERC721 struct {
	Address  *common.Address
	Contract *bindings.ERC721
}

func DeployNewERC721Contract(ctx context.Context, c *Client, sponsors *WalletRegistry, minGasPrice uint64) (*ERC721, error) {
	sponsor := sponsors.GetAvailableWallet()
	if sponsor == nil {
		return nil, fmt.Errorf("no sponsor wallet for ERC721 contract")
	}
	defer sponsors.Unlock(sponsor.Address)

	tops, err := bind.NewKeyedTransactorWithChainID(sponsor.PrivateKey, c.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor ERC721: %w", err)
	}
	tops.Nonce = big.NewInt(int64(sponsor.Nonce))

	head, err := c.Ec.HeaderByNumber(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get head: %w", err)
	}

	if head.BaseFee == nil {
		// Increase gas price by 20% to avoid "price too low" errors
		gasPrice, err := c.Ec.SuggestGasPrice(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to suggest gas price: %w", err)
		}
		if gasPrice.Cmp(big.NewInt(0)) == 0 {
			gasPrice = big.NewInt(int64(minGasPrice)) // 1 gwei
		}
		tops.GasPrice = gasPrice
	} else {
		tipCap, err := c.Ec.SuggestGasTipCap(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to suggest gas tip cap: %w", err)
		}
		if tipCap.Cmp(big.NewInt(0)) == 0 {
			tipCap = big.NewInt(int64(minGasPrice)) // 1 gwei
		}
		feeCap := new(big.Int).Add(
			tipCap,
			new(big.Int).Mul(head.BaseFee, big.NewInt(2)),
		)
		tops.GasFeeCap = feeCap
		tops.GasTipCap = tipCap
	}

	addr, tx, contract, err := bindings.DeployERC721(tops, c.Ec)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy ERC721 contract: %w", err)
	}
	sponsor.Nonce++
	sponsor.OffsetNonce++
	hash := tx.Hash()

	if _, err := WaitUntilMined(ctx, c.Ec, hash, EthDefaultBlockTime, 10*time.Millisecond); err != nil {
		return nil, fmt.Errorf("failed to wait for ERC721 contract to be mined: %w", err)
	}

	return &ERC721{
		Address:  &addr,
		Contract: contract,
	}, nil
}

func InitExistingERC721Contract(ec *ethclient.Client, address common.Address) (*ERC721, error) {
	contract, err := bindings.NewERC721(address, ec)
	if err != nil {
		return nil, fmt.Errorf("failed to create ERC721 contract: %w", err)
	}

	return &ERC721{
		Address:  &address,
		Contract: contract,
	}, nil
}

func (e *ERC721) MintERC721Contract(ctx context.Context, c *Client, testers *WalletRegistry) error {
	if e.Address == nil {
		return fmt.Errorf("ERC721 contract is not initialized")
	}

	sem := semaphore.NewWeighted(MaxNumberOfCreatingWalletsAtOnce)
	errCh := make(chan error, testers.Ln())

	// ctx, cancel := context.WithTimeout(ctx, WalletsTimeout)
	// defer cancel()

	wallets := testers.All()
	for _, wallet := range wallets {
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

		go func(wallet *Wallet) {
			defer sem.Release(1)

			signFn := func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
				sig := types.LatestSignerForChainID(c.ChainID)
				return types.SignTx(tx, sig, wallet.PrivateKey)
			}
			tx, err := e.Contract.MintBatch(&bind.TransactOpts{
				From:    wallet.Address,
				Signer:  signFn,
				Context: ctx,
			}, wallet.Address, big.NewInt(int64(1)))
			if err != nil {
				errCh <- fmt.Errorf("failed to mint ERC721 contract amount for wallet %s: %w", wallet.Address.Hex(), err)
				return
			}
			wallet.Nonce++
			wallet.OffsetNonce++

			if _, err := WaitUntilMined(ctx, c.Ec, tx.Hash(), EthDefaultBlockTime, 10*time.Millisecond); err != nil {
				errCh <- fmt.Errorf("failed to wait for erc721 mint transaction to be mined: %w", err)
				return
			}
		}(wallet)
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

func (e *ERC721) Balance(ctx context.Context, wallet *Wallet) error {
	if err := RepeatWithTimeout(ctx, 3*time.Second, 10*time.Millisecond, func(ctx context.Context) error {
		_, err := e.Contract.BalanceOf(&bind.CallOpts{Context: ctx}, wallet.Address)
		if err != nil {
			return fmt.Errorf("failed to get balance of ERC721 contract: %w", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to check ERC721 balance for wallet %s: %w", wallet.Address.Hex(), err)
	}
	return nil
}
