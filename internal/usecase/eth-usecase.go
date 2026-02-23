package usecase

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthUC interface {
	ChainInfo(ctx context.Context) (chainID string, err error)
	Balance(ctx context.Context, address string) (wei string, eth string, err error)
}

type EthUsecase struct {
	client *ethclient.Client
}

func NewEthUsecase(rpcURL string) (*EthUsecase, error) {
	if rpcURL == "" {
		return nil, fmt.Errorf("rpc url is empty")
	}
	c, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to dial ethereum rpc: %w", err)
	}
	return &EthUsecase{client: c}, nil
}

func (e *EthUsecase) ChainInfo(ctx context.Context) (string, error) {
	chainID, err := e.client.ChainID(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get chain id: %w", err)
	}
	return chainID.String(), nil
}

func (e *EthUsecase) Balance(ctx context.Context, address string) (wei string, eth string, err error) {
	if !common.IsHexAddress(address) {
		return "", "", fmt.Errorf("invalid address")
	}

	addr := common.HexToAddress(address)

	// nil block number = latest
	balWei, err := e.client.BalanceAt(ctx, addr, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to get balance: %w", err)
	}

	// wei -> eth (as decimal string)
	balEth := weiToEthString(balWei)

	return balWei.String(), balEth, nil
}

func weiToEthString(wei *big.Int) string {
	// ETH = wei / 1e18
	f := new(big.Float).SetInt(wei)
	denom := new(big.Float).SetInt(big.NewInt(0).Exp(big.NewInt(10), big.NewInt(18), nil))
	f.Quo(f, denom)

	// 18 decimals max; trims trailing zeros reasonably when printed
	return f.Text('f', 18)
}
