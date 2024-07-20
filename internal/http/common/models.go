package common

import "math/big"

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
