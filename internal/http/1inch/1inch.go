package oneinch

import (
	"context"
	"fmt"
	"math/big"
	uri "net/url"
	"os"
	"strconv"

	"github.com/onmetahq/go-evm.git/internal/http/common"
	metahttp "github.com/onmetahq/meta-http/pkg/meta_http"
)

type oneInch struct {
	client metahttp.Requests
}

func NewOneInch(client metahttp.Requests) *oneInch {
	return &oneInch{
		client: client,
	}
}

func (o *oneInch) FetchSupportedTokens(ctx context.Context, chainId uint64) ([]OneInchToken, error) {
	var res OneInchTokens
	url := fmt.Sprintf("/%d/tokens", chainId)
	_, err := o.client.Get(ctx, url, map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", os.Getenv("1INCH_KEY")),
	}, &res)

	if err != nil {
		return []OneInchToken{}, fmt.Errorf("unable to fetch all tokens from 1inch, err: %v", err)
	}

	var out []OneInchToken
	for _, v := range res.Tokens {
		out = append(out, v)
	}

	return out, nil
}

func (o *oneInch) FetchExactInQuote(ctx context.Context, req common.QuoteReq) (common.QuoteRes, error) {
	v := uri.Values{}
	v.Add("src", req.Src)
	v.Add("dst", req.Dst)
	v.Add("amount", req.Amount.String())
	v.Add("includeTokensInfo", "true")
	v.Add("includeGas", "true")
	query := v.Encode()
	url := fmt.Sprintf("/%d/quote?%s", req.ChainId, query)
	var res OneInchQuoteResponse
	_, err := o.client.Get(ctx, url, map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", os.Getenv("1INCH_KEY")),
	}, &res)
	if err != nil {
		return common.QuoteRes{}, fmt.Errorf("unable to fetch 1inch quote, err: %v", err)
	}

	return parse1inchResponse(req, res)
}

func (o *oneInch) FetchExactInSwapCallData(ctx context.Context, req common.QuoteReq) (OneInchSwapResponse, error) {
	v := uri.Values{}
	v.Add("src", req.Src)
	v.Add("dst", req.Dst)
	v.Add("amount", req.Amount.String())
	v.Add("from", req.From)
	v.Add("origin", req.From)
	v.Add("slippage", strconv.Itoa(int(req.SlippagePercentage)))
	v.Add("includeTokensInfo", "true")
	v.Add("includeGas", "true")
	v.Add("disableEstimate", "false")

	if req.SkipValidation {
		v.Add("disableEstimate", "true")
	}

	if len(req.Referrer) > 0 {
		v.Add("referrer", req.Referrer)
	}
	query := v.Encode()
	url := fmt.Sprintf("/%d/swap?%s", req.ChainId, query)
	var res OneInchSwapResponse
	_, err := o.client.Get(ctx, url, map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", os.Getenv("1INCH_KEY")),
	}, &res)
	if err != nil {
		return OneInchSwapResponse{}, fmt.Errorf("unable to fetch 1inch quote, err: %v", err)
	}

	return res, nil
}

func (o *oneInch) FetchExactOutQuote(ctx context.Context, req common.QuoteReq) (common.QuoteRes, error) {
	return common.QuoteRes{}, fmt.Errorf("operation exact out is not supported")
}

func (o *oneInch) FetchExactOutSwapCallData(ctx context.Context, req common.QuoteReq) (OneInchSwapResponse, error) {
	return OneInchSwapResponse{}, fmt.Errorf("operation not supported")
}

type OneInchToken struct {
	Address  string   `json:"address"`
	Symbol   string   `json:"symbol"`
	Decimals int      `json:"decimals"`
	Name     string   `json:"name"`
	LogoURI  string   `json:"logoURI"`
	Eip2612  bool     `json:"eip2612"`
	Tags     []string `json:"tags"`
}

type OneInchTokens struct {
	Tokens map[string]OneInchToken `json:"tokens"`
}

type OneInchQuoteResponse struct {
	FromToken OneInchToken `json:"fromToken"`
	Gas       int64        `json:"gas"`
	ToAmount  string       `json:"toAmount"`
	ToToken   OneInchToken `json:"toToken"`
}

type OneInchSwapResponse struct {
	FromToken OneInchToken `json:"fromToken"`
	ToAmount  string       `json:"toAmount"`
	ToToken   OneInchToken `json:"toToken"`
	Tx        struct {
		Data     string `json:"data"`
		From     string `json:"from"`
		Gas      int    `json:"gas"`
		GasPrice string `json:"gasPrice"`
		To       string `json:"to"`
		Value    string `json:"value"`
	} `json:"tx"`
}

func parse1inchResponse(req common.QuoteReq, quote OneInchQuoteResponse) (common.QuoteRes, error) {
	outAmount, ok := new(big.Int).SetString(quote.ToAmount, 0)
	if !ok {
		return common.QuoteRes{}, fmt.Errorf("invalid out amount from 0x, amount: %v", quote.ToAmount)
	}

	return common.QuoteRes{
		ChainId:    req.ChainId,
		Src:        req.Src,
		Dst:        req.Dst,
		FromAmount: req.Amount,
		ToAmount:   outAmount,
		Gas:        big.NewInt(quote.Gas),
	}, nil
}