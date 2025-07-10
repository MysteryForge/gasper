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

type ERC20 struct {
	Address  *common.Address
	Contract *bindings.ERC20
}

func DeployNewERC20Contract(ctx context.Context, c *Client, sponsors *WalletRegistry, minGasPrice uint64) (*ERC20, error) {
	sponsor := sponsors.GetAvailableWallet()
	if sponsor == nil {
		return nil, fmt.Errorf("no sponsor wallet for ERC20 contract")
	}
	defer sponsors.Unlock(sponsor.Address)

	tops, err := bind.NewKeyedTransactorWithChainID(sponsor.PrivateKey, c.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor ERC20: %w", err)
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

	addr, tx, contract, err := bindings.DeployERC20(tops, c.Ec)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy ERC20 contract: %w", err)
	}
	sponsor.Nonce++
	sponsor.OffsetNonce++
	hash := tx.Hash()

	if _, err := WaitUntilMined(ctx, c.Ec, hash, EthDefaultBlockTime, 10*time.Millisecond); err != nil {
		return nil, fmt.Errorf("failed to wait for ERC20 contract to be mined: %w", err)
	}

	return &ERC20{
		Address:  &addr,
		Contract: contract,
	}, nil
}

func InitExistingERC20Contract(ec *ethclient.Client, address common.Address) (*ERC20, error) {
	contract, err := bindings.NewERC20(address, ec)
	if err != nil {
		return nil, fmt.Errorf("failed to create ERC20 contract: %w", err)
	}

	return &ERC20{
		Address:  &address,
		Contract: contract,
	}, nil
}

func (e *ERC20) MintERC20Contract(ctx context.Context, c *Client, testers *WalletRegistry, amount big.Int) error {
	if e.Address == nil {
		return fmt.Errorf("ERC20 contract is not initialized")
	}
	if amount.Cmp(big.NewInt(0)) == 0 {
		return nil
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

			balance, err := e.Balance(ctx, wallet.Address)
			if err != nil {
				errCh <- err
				return
			}
			if balance.Cmp(&amount) >= 0 {
				return
			}

			diffAmount := new(big.Int).Sub(&amount, balance)

			signFn := func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
				sig := types.LatestSignerForChainID(c.ChainID)
				return types.SignTx(tx, sig, wallet.PrivateKey)
			}
			tx, err := e.Contract.Mint(&bind.TransactOpts{
				From:    wallet.Address,
				Signer:  signFn,
				Context: ctx,
			}, diffAmount)
			if err != nil {
				errCh <- err
				return
			}
			wallet.Nonce++
			wallet.OffsetNonce++

			if _, err := WaitUntilMined(ctx, c.Ec, tx.Hash(), EthDefaultBlockTime, 10*time.Millisecond); err != nil {
				errCh <- fmt.Errorf("failed to wait for erc20 mint transaction to be mined: %w", err)
				return
			}

			balance, err = e.Balance(ctx, wallet.Address)
			if err != nil {
				errCh <- err
				return
			}
			if balance.Cmp(&amount) < 0 {
				errCh <- fmt.Errorf("ERC20 contract balance is less than required amount")
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

func (e *ERC20) Balance(ctx context.Context, addr common.Address) (*big.Int, error) {
	var balance *big.Int
	if err := RepeatWithTimeout(ctx, 3*time.Second, 10*time.Millisecond, func(ctx context.Context) error {
		b, err := e.Contract.BalanceOf(&bind.CallOpts{Context: ctx}, addr)
		if err != nil {
			return err
		}
		balance = b

		return nil
	}); err != nil {
		return balance, fmt.Errorf("failed to check ERC20 balance for wallet %s: %w", addr.Hex(), err)
	}
	return balance, nil
}
