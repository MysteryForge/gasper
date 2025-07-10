package eth

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mysteryforge/gasper/bindings"
	"golang.org/x/sync/semaphore"
)

type BatchFunder struct {
	Address  *common.Address
	Contract *bindings.BatchFunder
}

func DeployNewBatchFunder(ctx context.Context, c *Client, sponsors *WalletRegistry, minGasPrice uint64) (*BatchFunder, error) {
	sponsor := sponsors.GetAvailableWallet()
	if sponsor == nil {
		return nil, fmt.Errorf("no sponsor wallet for batch funder contract")
	}
	defer sponsors.Unlock(sponsor.Address)

	tops, err := bind.NewKeyedTransactorWithChainID(sponsor.PrivateKey, c.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactor BatchFunder: %w", err)
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

	addr, tx, contract, err := bindings.DeployBatchFunder(tops, c.Ec)
	if err != nil {
		return nil, fmt.Errorf("failed to deploy BatchFunder contract: %w", err)
	}

	sponsor.Nonce++
	sponsor.OffsetNonce++

	hash := tx.Hash()
	if _, err := WaitUntilMined(ctx, c.Ec, hash, EthDefaultBlockTime, 10*time.Millisecond); err != nil {
		return nil, fmt.Errorf("failed to wait for BatchFunder contract to be mined: %w", err)
	}

	return &BatchFunder{
		Address:  &addr,
		Contract: contract,
	}, nil
}

func InitExistingBatchFunder(c *Client, addr common.Address) (*BatchFunder, error) {
	contract, err := bindings.NewBatchFunder(addr, c.Ec)
	if err != nil {
		return nil, fmt.Errorf("failed to create batch funder contract: %w", err)
	}

	return &BatchFunder{
		Address:  &addr,
		Contract: contract,
	}, nil
}

func (b *BatchFunder) FundWallets(ctx context.Context, c *Client, sponsors *WalletRegistry, recipients *WalletRegistry, amount *big.Int) error {
	if b.Address == nil {
		return fmt.Errorf("batch funder contract is not initialized")
	}

	lnOfSponsors := int64(sponsors.Ln())
	lnOfRecipients := recipients.Ln()

	sem := semaphore.NewWeighted(lnOfSponsors)
	wg := sync.WaitGroup{}
	errCh := make(chan error, lnOfRecipients)

	addresses := make([]common.Address, 0, MaxNumberOfFundingWalletsAtOnce)
	added := 0
	indx := 0
	wallets := recipients.All()
	for _, wallet := range wallets {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return err
		default:
		}

		addresses = append(addresses, wallet.Address)
		added++
		if added == MaxNumberOfFundingWalletsAtOnce || indx == lnOfRecipients-1 {
			if err := sem.Acquire(ctx, 1); err != nil {
				return err
			}

			go func(addrs []common.Address) {
				defer sem.Release(1)

				ticker := time.NewTicker(250 * time.Millisecond)
				defer ticker.Stop()

				var sponsor *Wallet
				for ctx.Err() == nil {
					select {
					case <-ctx.Done():
						return
					case <-errCh:
						return
					default:
					}
					sponsor = sponsors.GetAvailableWallet()
					if sponsor != nil {
						break
					}
					<-ticker.C
				}
				if sponsor == nil {
					errCh <- fmt.Errorf("no sponsor wallet for funding wallets")
					return
				}
				defer sponsors.Unlock(sponsor.Address)

				signFn := func(addr common.Address, tx *types.Transaction) (*types.Transaction, error) {
					sig := types.LatestSignerForChainID(c.ChainID)
					return types.SignTx(tx, sig, sponsor.PrivateKey)
				}

				totalAmount := new(big.Int).Mul(amount, big.NewInt(int64(len(addrs))))

				tx, err := b.Contract.BatchSend(&bind.TransactOpts{
					From:    sponsor.Address,
					Signer:  signFn,
					Context: ctx,
					Value:   totalAmount,
					Nonce:   big.NewInt(int64(sponsor.Nonce)),
				}, addrs, amount)
				if err != nil {
					errCh <- fmt.Errorf("failed to batch fund wallets: %w", err)
					return
				}
				sponsor.Nonce++
				sponsor.OffsetNonce++
				wg.Add(1)

				go func() {
					defer wg.Done()
					if _, err := WaitUntilMined(ctx, c.Ec, tx.Hash(), WalletsTimeout, WalletsInterval); err != nil {
						errCh <- fmt.Errorf("failed to wait for batch funding transaction to be mined: %w", err)
						return
					}
				}()
			}(addresses)
			addresses = make([]common.Address, 0, MaxNumberOfFundingWalletsAtOnce)
			added = 0
		}
		indx++
	}

	if err := sem.Acquire(ctx, lnOfSponsors); err != nil {
		return err
	}

	wg.Wait()

	close(errCh)
	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil

}
