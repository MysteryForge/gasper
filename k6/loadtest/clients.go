package loadtest

import (
	"context"
	"fmt"
	"math/big"
	"sync"

	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-logr/logr"
	"github.com/mysteryforge/gasper/k6/eth"
	"go.k6.io/k6/js/modules"
)

type Clients struct {
	list []Client
	log  logr.Logger
}

func NewClients(ctx context.Context, cfgs []*clientConfig, uid string, dbs map[string]*eth.PebbleDb, log logr.Logger) (*Clients, error) {

	clients := &Clients{
		list: make([]Client, 0, len(cfgs)),
		log:  log,
	}

	for _, cfg := range cfgs {
		db := dbs[cfg.DBPath]

		c, err := NewClient(ctx, cfg, uid, db, log)
		if err != nil {
			return nil, fmt.Errorf("failed to create client: %w", err)
		}
		clients.Add(c)
	}
	return clients, nil
}

type Result struct {
	Err  error `json:"err,omitempty"`
	Data any   `json:"data,omitempty"`
}

func (cs *Clients) Export() modules.Exports {
	return modules.Exports{}
}

func (cs *Clients) Add(c Client) {
	cs.list = append(cs.list, c)
}

type SharedWallet struct {
	UID     string `json:"uid"`
	Address string `json:"address"`
}

func (cs *Clients) RequestSharedWallet() map[string]Result {
	return executeOnAllClients(cs, func(c Client) (*SharedWallet, error) {
		w, err := c.RequestSharedWallet()
		if err != nil {
			return nil, err
		}
		// pkBytes := crypto.FromECDSA(w.PrivateKey)
		// pkHex := hex.EncodeToString(pkBytes)
		return &SharedWallet{
			UID:     c.UID(),
			Address: w.Address.Hex(),
		}, nil
	})
}

func (cs *Clients) ReleaseSharedWallet(address common.Address) {
	for _, c := range cs.list {
		c.ReleaseSharedWallet(address)
	}
}

func (cs *Clients) ChainID(vu modules.VU, metrics *EthMetrics) map[string]Result {
	return executeOnAllClients(cs, func(c Client) (*big.Int, error) {
		return c.ChainID(vu, metrics)
	})
}

func (cs *Clients) TxPoolStatus(vu modules.VU, metrics *EthMetrics) map[string]Result {
	return executeOnAllClients(cs, func(c Client) (*eth.PoolStatus, error) {
		return c.TxPoolStatus(vu, metrics)
	})
}

func (cs *Clients) ReportBlockMetrics(vu modules.VU, metrics *EthMetrics) map[string]Result {
	return executeOnAllClients(cs, func(c Client) (any, error) {
		return nil, c.ReportBlockMetrics(vu, metrics)
	})
}

func (cs *Clients) SendTransaction(vu modules.VU, metrics *EthMetrics, params map[string]interface{}) map[string]Result {
	return executeOnAllClients(cs, func(c Client) (string, error) {
		options := parseSendTransactionParams(params, c.UID())
		hash, err := c.SendTransaction(vu, metrics, options...)
		if err != nil {
			return "", err
		}
		return hash.Hex(), nil
	})
}

func (cs *Clients) SendERC20Transaction(vu modules.VU, metrics *EthMetrics, params map[string]interface{}) map[string]Result {
	return executeOnAllClients(cs, func(c Client) (string, error) {
		options := parseSendTransactionParams(params, c.UID())
		hash, err := c.SendERC20Transaction(vu, metrics, options...)
		if err != nil {
			return "", err
		}
		return hash.Hex(), nil
	})
}

func (cs *Clients) SendOffsetTransaction(vu modules.VU, metrics *EthMetrics, params map[string]interface{}) map[string]Result {
	return executeOnAllClients(cs, func(c Client) (string, error) {
		options := parseSendTransactionParams(params, c.UID())
		hash, err := c.SendTransaction(vu, metrics, options...)
		if err != nil {
			return "", err
		}
		return hash.Hex(), nil
	})
}

func (cs *Clients) SendERC721Transaction(vu modules.VU, metrics *EthMetrics, params map[string]interface{}) map[string]Result {
	return executeOnAllClients(cs, func(c Client) (string, error) {
		options := parseSendTransactionParams(params, c.UID())
		hash, err := c.SendERC721Transaction(vu, metrics, options...)
		if err != nil {
			return "", err
		}
		return hash.Hex(), nil
	})
}

type ContractDeploymentResponse struct {
	TransactionHash string
	ContractAddress string
	OwnerWallet     string
}

func (cs *Clients) DeployContract(vu modules.VU, params map[string]interface{}) map[string]Result {
	return executeOnAllClients(cs, func(c Client) (*ContractDeploymentResponse, error) {
		parsed, err := eth.ParseDeployContractParams(params)
		if err != nil {
			return nil, err
		}
		wallet, hash, addr, err := c.DeployContract(vu, parsed)
		if err != nil {
			return nil, err
		}
		return &ContractDeploymentResponse{
			TransactionHash: hash.Hex(),
			ContractAddress: addr.Hex(),
			OwnerWallet:     wallet.Hex(),
		}, nil
	})
}

func (cs *Clients) TxContract(vu modules.VU, metrics *EthMetrics, params map[string]interface{}) map[string]Result {
	return executeOnAllClients(cs, func(c Client) (string, error) {
		parsed, err := eth.ParseFnContractParams(params)
		if err != nil {
			return "", err
		}
		hash, err := c.TxContract(vu, metrics, parsed)
		if err != nil {
			return "", err
		}
		return hash.Hex(), err
	})
}

func (cs *Clients) CallContract(vu modules.VU, metrics *EthMetrics, params map[string]interface{}) map[string]Result {
	return executeOnAllClients(cs, func(c Client) (interface{}, error) {
		parsed, err := eth.ParseFnContractParams(params)
		if err != nil {
			return nil, err
		}
		return c.CallContract(vu, metrics, parsed)
	})
}

func (cs *Clients) TxInfoByHash(vu modules.VU, metrics *EthMetrics, txHash string) map[string]Result {
	return executeOnAllClients(cs, func(c Client) (*eth.TransactionInfo, error) {
		if txHash == "" {
			return nil, fmt.Errorf("tx_hash is empty")
		}
		txHash = strings.TrimPrefix(txHash, "0x")
		if len(txHash) != 64 {
			return nil, fmt.Errorf("tx_hash does not look like a valid hash")
		}

		hash := common.HexToHash(txHash)
		if hash == (common.Hash{}) {
			return nil, fmt.Errorf("tx_hash is empty")
		}
		res, err := c.GetTxInfo(vu, metrics, hash)
		if err != nil {
			return nil, err
		}
		return res, nil
	})
}

func (cs *Clients) Call(vu modules.VU, metrics *EthMetrics, method string, params []interface{}) map[string]Result {
	return executeOnAllClients(cs, func(c Client) (interface{}, error) {
		return c.Call(vu, metrics, method, params...)
	})
}

func executeOnAllClients[T any](cs *Clients, fn func(Client) (T, error)) map[string]Result {
	res := make(map[string]Result)
	wg := sync.WaitGroup{}
	wg.Add(len(cs.list))

	for _, c := range cs.list {
		ch := make(chan Result, 1)
		go func(client Client) {
			defer wg.Done()
			data, err := fn(client)
			if err != nil {
				ch <- Result{Err: err}
				return
			}
			ch <- Result{Data: data}
		}(c)
		res[c.UID()] = <-ch
		close(ch)
	}
	wg.Wait()
	return res
}

func parseSendTransactionParams(params map[string]interface{}, clientUid string) []TransactionOption {
	var options []TransactionOption

	if noSend, ok := params["no_send"].(bool); ok {
		options = append(options, WithTransactionNoSend(noSend))
	}

	if delay, ok := params["confirmation_delay"].(int64); ok {
		options = append(options, WithTransactionConfirmationDelay(uint64(delay)))
	}

	if offset, ok := params["nonce_offset"].(int64); ok {
		options = append(options, WithTransactionOffsetNonce(uint64(offset)))
	}

	if txCount, ok := params["tx_count"].(int64); ok {
		options = append(options, WithTransactionTxCount(uint64(txCount)))
	}

	if gasPriceMultiplier, ok := params["gas_price_multiplier"].(int64); ok {
		options = append(options, WithTransactionGasPriceMultiplier(uint64(gasPriceMultiplier)))
	}

	if wallets, ok := params["wallets"].([]interface{}); ok {
		if len(wallets) > 0 {
			addresses := make([]*common.Address, 0, len(wallets))
			for _, w := range wallets {
				mapWallet, ok := w.(map[string]interface{})
				if !ok {
					continue
				}
				uid, ok := mapWallet["uid"].(string)
				if !ok {
					continue
				}
				if uid != clientUid {
					continue
				}
				addr, ok := mapWallet["address"].(string)
				if !ok {
					continue
				}
				hex := common.HexToAddress(addr)
				addresses = append(addresses, &hex)
			}
			options = append(options, WithTransactionWalletAddress(addresses))
		}
	}

	return options
}
