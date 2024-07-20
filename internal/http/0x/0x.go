package zerox

import (
	"context"
	"fmt"
	"math/big"
	uri "net/url"
	"os"

	"github.com/onmetahq/go-evm.git/internal/http/common"
	metahttp "github.com/onmetahq/meta-http/pkg/meta_http"
)

type zeroX struct {
	client      metahttp.Requests
	chainUrlMap map[uint64]string
}

func NewZeroX(client metahttp.Requests, chainUrlMap map[uint64]string) *zeroX {
	return &zeroX{
		client:      client,
		chainUrlMap: chainUrlMap,
	}
}

func (o *zeroX) FetchSupportedTokens(ctx context.Context, chainId uint64) (any, error) {
	return "", fmt.Errorf("operation not supported")
}

func (o *zeroX) FetchExactInQuote(ctx context.Context, req common.QuoteReq) (common.QuoteRes, error) {
	v := uri.Values{}
	v.Add("buyToken", req.Dst)
	v.Add("sellToken", req.Src)
	v.Add("sellAmount", req.Amount.String())
	v.Add("slippagePercentage", "0.01")
	v.Add("skipValidation", "false")

	if req.From != "" {
		v.Add("takerAddress", req.From)
	}

	if req.SkipValidation {
		v.Add("skipValidation", "true")
	}

	res, err := o.price(ctx, req.ChainId, v.Encode())
	if err != nil {
		return common.QuoteRes{}, err
	}
	res.OutAmount = res.BuyAmount
	return parseZeroxResponse(req, res)
}

func (o *zeroX) FetchExactOutQuote(ctx context.Context, req common.QuoteReq) (common.QuoteRes, error) {
	v := uri.Values{}
	v.Add("buyToken", req.Dst)
	v.Add("sellToken", req.Src)
	v.Add("buyAmount", req.Amount.String())
	v.Add("slippagePercentage", "0.01")
	v.Add("skipValidation", "false")

	if req.From != "" {
		v.Add("takerAddress", req.From)
	}

	if req.SkipValidation {
		v.Add("skipValidation", "true")
	}

	res, err := o.price(ctx, req.ChainId, v.Encode())
	if err != nil {
		return common.QuoteRes{}, err
	}
	res.OutAmount = res.SellAmount
	return parseZeroxResponse(req, res)
}

func (o *zeroX) FetchExactInSwapCallData(ctx context.Context, req common.QuoteReq) (ZeroXSwapResponse, error) {
	v := uri.Values{}
	v.Add("buyToken", req.Dst)
	v.Add("sellToken", req.Src)
	v.Add("sellAmount", req.Amount.String())
	v.Add("slippagePercentage", "0.01")

	if req.From != "" {
		v.Add("takerAddress", req.From)
	}

	if req.SkipValidation {
		v.Add("skipValidation", "true")
	} else {
		v.Add("skipValidation", "false")
	}

	return o.quote(ctx, req.ChainId, v.Encode())
}

func (o *zeroX) FetchExactOutSwapCallData(ctx context.Context, req common.QuoteReq) (ZeroXSwapResponse, error) {
	v := uri.Values{}
	v.Add("buyToken", req.Dst)
	v.Add("sellToken", req.Src)
	v.Add("buyAmount", req.Amount.String())
	v.Add("slippagePercentage", "0.01")

	if req.From != "" {
		v.Add("takerAddress", req.From)
	}

	if req.SkipValidation {
		v.Add("skipValidation", "true")
	} else {
		v.Add("skipValidation", "false")
	}

	return o.quote(ctx, req.ChainId, v.Encode())
}

func (o *zeroX) quote(ctx context.Context, chainId uint64, queryParams string) (ZeroXSwapResponse, error) {
	base, ok := o.chainUrlMap[chainId]
	if !ok {
		return ZeroXSwapResponse{}, fmt.Errorf("unsupported chainId %d", chainId)
	}

	url := fmt.Sprintf("%s/swap/v1/quote?%s", base, queryParams)
	var res ZeroXSwapResponse

	_, err := o.client.Get(ctx, url, map[string]string{
		"0x-api-key": os.Getenv("0X_KEY"),
	}, &res)

	if err != nil {
		return ZeroXSwapResponse{}, fmt.Errorf("unable to fetch 0x quote, err: %v", err)
	}

	return res, nil
}

func (o *zeroX) price(ctx context.Context, chainId uint64, queryParams string) (ZeroXQuoteResponse, error) {
	base, ok := o.chainUrlMap[chainId]
	if !ok {
		return ZeroXQuoteResponse{}, fmt.Errorf("unsupported chainId %d", chainId)
	}

	url := fmt.Sprintf("%s/swap/v1/price?%s", base, queryParams)
	var res ZeroXQuoteResponse

	_, err := o.client.Get(ctx, url, map[string]string{
		"0x-api-key": os.Getenv("0X_KEY"),
	}, &res)

	if err != nil {
		return ZeroXQuoteResponse{}, fmt.Errorf("unable to fetch 0x quote, err: %v", err)
	}

	return res, nil
}

type ZeroXQuoteResponse struct {
	AllowanceTarget      string `json:"allowanceTarget"`
	BuyAmount            string `json:"buyAmount"`
	BuyTokenAddress      string `json:"buyTokenAddress"`
	BuyTokenToEthRate    string `json:"buyTokenToEthRate"`
	ChainID              uint64 `json:"chainId"`
	EstimatedGas         string `json:"estimatedGas"`
	EstimatedPriceImpact string `json:"estimatedPriceImpact"`
	Fees                 struct {
		ZeroExFee struct {
			BillingType string `json:"billingType"`
			FeeAmount   string `json:"feeAmount"`
			FeeToken    string `json:"feeToken"`
			FeeType     string `json:"feeType"`
		} `json:"zeroExFee"`
	} `json:"fees"`
	Gas                string `json:"gas"`
	GasPrice           string `json:"gasPrice"`
	GrossBuyAmount     string `json:"grossBuyAmount"`
	GrossPrice         string `json:"grossPrice"`
	GrossSellAmount    string `json:"grossSellAmount"`
	MinimumProtocolFee string `json:"minimumProtocolFee"`
	Price              string `json:"price"`
	ProtocolFee        string `json:"protocolFee"`
	SellAmount         string `json:"sellAmount"`
	SellTokenAddress   string `json:"sellTokenAddress"`
	SellTokenToEthRate string `json:"sellTokenToEthRate"`
	Sources            []struct {
		Name       string `json:"name"`
		Proportion string `json:"proportion"`
	} `json:"sources"`
	Value     string `json:"value"`
	OutAmount string `json:"outAmount"`
}

type ZeroXSwapResponse struct {
	AllowanceTarget      string `json:"allowanceTarget"`
	BuyAmount            string `json:"buyAmount"`
	BuyTokenAddress      string `json:"buyTokenAddress"`
	BuyTokenToEthRate    string `json:"buyTokenToEthRate"`
	ChainID              int    `json:"chainId"`
	Data                 string `json:"data"`
	DecodedUniqueID      string `json:"decodedUniqueId"`
	EstimatedGas         string `json:"estimatedGas"`
	EstimatedPriceImpact string `json:"estimatedPriceImpact"`
	Fees                 struct {
		ZeroExFee struct {
			BillingType string `json:"billingType"`
			FeeAmount   string `json:"feeAmount"`
			FeeToken    string `json:"feeToken"`
			FeeType     string `json:"feeType"`
		} `json:"zeroExFee"`
	} `json:"fees"`
	Gas                string `json:"gas"`
	GasPrice           string `json:"gasPrice"`
	GrossBuyAmount     string `json:"grossBuyAmount"`
	GrossPrice         string `json:"grossPrice"`
	GrossSellAmount    string `json:"grossSellAmount"`
	GuaranteedPrice    string `json:"guaranteedPrice"`
	MinimumProtocolFee string `json:"minimumProtocolFee"`
	Orders             []struct {
		Fill struct {
			AdjustedOutput string `json:"adjustedOutput"`
			Gas            int    `json:"gas"`
			Input          string `json:"input"`
			Output         string `json:"output"`
		} `json:"fill"`
		FillData struct {
			Router           string   `json:"router"`
			TokenAddressPath []string `json:"tokenAddressPath"`
		} `json:"fillData,omitempty"`
		MakerAmount string `json:"makerAmount"`
		MakerToken  string `json:"makerToken"`
		Source      string `json:"source"`
		TakerAmount string `json:"takerAmount"`
		TakerToken  string `json:"takerToken"`
		Type        int    `json:"type"`
	} `json:"orders"`
	Price              string `json:"price"`
	ProtocolFee        string `json:"protocolFee"`
	SellAmount         string `json:"sellAmount"`
	SellTokenAddress   string `json:"sellTokenAddress"`
	SellTokenToEthRate string `json:"sellTokenToEthRate"`
	Sources            []struct {
		Name       string `json:"name"`
		Proportion string `json:"proportion"`
	} `json:"sources"`
	To    string `json:"to"`
	Value string `json:"value"`
}

func parseZeroxResponse(req common.QuoteReq, quote ZeroXQuoteResponse) (common.QuoteRes, error) {
	outAmount, ok := new(big.Int).SetString(quote.OutAmount, 0)
	if !ok {
		return common.QuoteRes{}, fmt.Errorf("invalid out amount from 0x, amount: %v", quote.OutAmount)
	}

	gas, ok := new(big.Int).SetString(quote.Gas, 0)
	if !ok {
		return common.QuoteRes{}, fmt.Errorf("invalid out amount from 0x, amount: %v", quote.OutAmount)
	}

	gasPrice, ok := new(big.Int).SetString(quote.GasPrice, 0)
	if !ok {
		return common.QuoteRes{}, fmt.Errorf("invalid out amount from 0x, amount: %v", quote.OutAmount)
	}

	return common.QuoteRes{
		ChainId:    req.ChainId,
		Src:        req.Src,
		Dst:        req.Dst,
		FromAmount: req.Amount,
		ToAmount:   outAmount,
		Gas:        gas,
		GasPrice:   gasPrice,
	}, nil
}
