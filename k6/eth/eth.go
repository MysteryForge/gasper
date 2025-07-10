package eth

import (
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

const (
	EthDefaultBlockTimeNum = 12
	EthDefaultBlockTime    = EthDefaultBlockTimeNum * time.Second

	WalletsTimeout                   = 10 * time.Minute
	WalletsInterval                  = 1 * time.Second
	MaxNumberOfCreatingWalletsAtOnce = 500
	MaxNumberOfFundingWalletsAtOnce  = 700
)

type TransactionType string

const (
	TransactionTypeETH    TransactionType = "EIP155"
	TransactionTypeERC20  TransactionType = "ERC20"
	TransactionTypeERC721 TransactionType = "ERC721"
)

type TransactionInfo struct {
	TxHash  string `json:"tx_hash"`
	Status  uint64 `json:"status"`
	GasUsed uint64 `json:"gas_used"`
}

type SlimBlock struct {
	Number       *big.Int `json:"number"`
	Timestamp    uint64   `json:"timestamp"`
	GasUsed      uint64   `json:"gasUsed"`
	Transactions []string `json:"transactions"` // hashes only, since we set 'false'
}

func (b *SlimBlock) MarshalJSON() ([]byte, error) {
	type SlimBlock struct {
		Number       *hexutil.Big   `json:"number"`
		Timestamp    hexutil.Uint64 `json:"timestamp"`
		GasUsed      hexutil.Uint64 `json:"gasUsed"`
		Transactions []string       `json:"transactions"`
	}

	return json.Marshal(SlimBlock{
		Number:       (*hexutil.Big)(b.Number),
		Timestamp:    hexutil.Uint64(b.Timestamp),
		GasUsed:      hexutil.Uint64(b.GasUsed),
		Transactions: b.Transactions,
	})
}

func (b *SlimBlock) UnmarshalJSON(input []byte) error {
	type SlimBlock struct {
		Number       *hexutil.Big    `json:"number"`
		Timestamp    *hexutil.Uint64 `json:"timestamp"`
		GasUsed      *hexutil.Uint64 `json:"gasUsed"`
		Transactions []string        `json:"transactions"`
	}

	var dec SlimBlock
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	b.Transactions = dec.Transactions

	if dec.Number == nil {
		return errors.New("missing required field 'number' for Header")
	}
	b.Number = (*big.Int)(dec.Number)
	if dec.GasUsed == nil {
		return errors.New("missing required field 'gasUsed' for Header")
	}
	b.GasUsed = uint64(*dec.GasUsed)
	if dec.Timestamp == nil {
		return errors.New("missing required field 'timestamp' for Header")
	}
	b.Timestamp = uint64(*dec.Timestamp)

	return nil
}
