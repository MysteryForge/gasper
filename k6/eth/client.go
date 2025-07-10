package eth

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Client struct {
	uri     string
	Rc      *rpc.Client
	Ec      *ethclient.Client
	ChainID *big.Int
}

func NewClient(ctx context.Context, uri string) (*Client, error) {
	rc, ec, chainID, err := InitEthClient(ctx, uri)
	if err != nil {
		return nil, fmt.Errorf("failed to create eth client for %s: %w", SanitizeURL(uri), err)
	}

	c := &Client{
		uri:     uri,
		Rc:      rc,
		Ec:      ec,
		ChainID: chainID,
	}
	return c, nil
}

func (c *Client) Close() {
	c.Ec.Close()
	c.Rc.Close()
}

func (c *Client) SlimBlockByNumber(ctx context.Context, number *big.Int) (*SlimBlock, error) {
	var result SlimBlock
	if err := c.Rc.CallContext(ctx, &result, "eth_getBlockByNumber", hexutil.EncodeBig(number), false); err != nil {
		return nil, err
	}
	return &result, nil
}
