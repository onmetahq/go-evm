package dex

import (
	"fmt"
	"math/big"
)

type QuoteReq struct {
	ChainId            uint64
	Src                string
	Dst                string
	Amount             *big.Int
	From               string
	SlippagePercentage uint8
	SkipValidation     bool
	Referrer           string
}

type QuoteRes struct {
	ChainId    uint64
	Src        string
	Dst        string
	FromAmount *big.Int
	ToAmount   *big.Int
	Gas        *big.Int `json:"gas"`
	GasPrice   *big.Int `json:"gasPrice"`
}

func parseZeroxResponse(req QuoteReq, quote ZeroXQuoteResponse) (QuoteRes, error) {
	outAmount, ok := new(big.Int).SetString(quote.OutAmount, 0)
	if !ok {
		return QuoteRes{}, fmt.Errorf("invalid out amount from 0x, amount: %v", quote.OutAmount)
	}

	gas, ok := new(big.Int).SetString(quote.Gas, 0)
	if !ok {
		return QuoteRes{}, fmt.Errorf("invalid out amount from 0x, amount: %v", quote.OutAmount)
	}

	gasPrice, ok := new(big.Int).SetString(quote.GasPrice, 0)
	if !ok {
		return QuoteRes{}, fmt.Errorf("invalid out amount from 0x, amount: %v", quote.OutAmount)
	}

	return QuoteRes{
		ChainId:    req.ChainId,
		Src:        req.Src,
		Dst:        req.Dst,
		FromAmount: req.Amount,
		ToAmount:   outAmount,
		Gas:        gas,
		GasPrice:   gasPrice,
	}, nil
}

func parse1inchResponse(req QuoteReq, quote OneInchQuoteResponse) (QuoteRes, error) {
	outAmount, ok := new(big.Int).SetString(quote.ToAmount, 0)
	if !ok {
		return QuoteRes{}, fmt.Errorf("invalid out amount from 0x, amount: %v", quote.ToAmount)
	}

	return QuoteRes{
		ChainId:    req.ChainId,
		Src:        req.Src,
		Dst:        req.Dst,
		FromAmount: req.Amount,
		ToAmount:   outAmount,
		Gas:        big.NewInt(quote.Gas),
	}, nil
}
