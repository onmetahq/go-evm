package gas

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/lmittmann/w3"
	w3eth "github.com/lmittmann/w3/module/eth"
)

var DEFAULT = big.NewInt(0)

func PriceEstimate(ctx context.Context, client *w3.Client) (*big.Int, error) {
	var (
		chainID  uint64
		gasPrice big.Int
		errs     w3.CallErrors
	)

	if err := client.CallCtx(ctx,
		w3eth.ChainID().Returns(&chainID),
		w3eth.GasPrice().Returns(&gasPrice),
	); errors.As(err, &errs) {
		if errs[0] != nil {
			return DEFAULT, fmt.Errorf("failed to get chain ID: %s", err)
		} else if errs[1] != nil {
			return DEFAULT, fmt.Errorf("failed to get gas price: %s", errs[1])
		}
	} else if err != nil {
		return DEFAULT, fmt.Errorf("failed RPC request: %s", err)
	}

	return &gasPrice, nil
}

func EIP1559Estimate(ctx context.Context, client *w3.Client) (*big.Int, *big.Int, error) {
	var (
		chainID     uint64
		baseFee     big.Int
		priorityFee big.Int
		errs        w3.CallErrors
	)

	if err := client.CallCtx(ctx,
		w3eth.ChainID().Returns(&chainID),
		w3eth.GasPrice().Returns(&baseFee),
		w3eth.GasTipCap().Returns(&priorityFee),
	); errors.As(err, &errs) {
		if errs[0] != nil {
			return DEFAULT, DEFAULT, fmt.Errorf("failed to get chain ID: %s", err)
		} else if errs[1] != nil {
			return DEFAULT, DEFAULT, fmt.Errorf("failed to get gas price: %s", errs[1])
		}
	} else if err != nil {
		return DEFAULT, DEFAULT, fmt.Errorf("failed RPC request: %s", err)
	}

	return &baseFee, &priorityFee, nil
}
