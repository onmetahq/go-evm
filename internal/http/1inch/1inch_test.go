package oneinch

import (
	"context"
	"log/slog"
	"math/big"
	"testing"
	"time"

	"github.com/onmetahq/go-evm/internal/http/common"
	metahttp "github.com/onmetahq/meta-http/pkg/meta_http"
)

const TOKENA = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
const TOKENB = "0x2791bca1f2de4661ed88a30c99a7a9449aa84174"

func TestFetchSupportedTokens(t *testing.T) {
	c := metahttp.NewClient("https://api.1inch.dev/swap/v5.2", slog.Default(), 30*time.Second)

	oneClient := NewOneInch(c)
	res, err := oneClient.FetchSupportedTokens(context.Background(), 137)
	if err != nil {
		t.Fatalf("tokens err: %v", err)
	}

	if len(res) == 0 {
		t.Fatalf("received empty token list")
	}
}

func TestFetchExactInQuote(t *testing.T) {
	c := metahttp.NewClient("https://api.1inch.dev/swap/v5.2", slog.Default(), 30*time.Second)

	oneClient := NewOneInch(c)
	req := common.QuoteReq{
		Src:     TOKENB,
		Dst:     TOKENA,
		ChainId: 137,
		Amount:  big.NewInt(1000000),
	}
	res, err := oneClient.FetchExactInQuote(context.Background(), req)
	if err != nil {
		t.Fatalf("quote err: %v", err)
	}

	if res.ToAmount.Cmp(big.NewInt(0)) < 1 {
		t.Fatalf("invalid quote, quote: %s", res.ToAmount)
	}
}

func TestFetchExactInSwapCallData(t *testing.T) {
	c := metahttp.NewClient("https://api.1inch.dev/swap/v5.2", slog.Default(), 30*time.Second)

	oneClient := NewOneInch(c)
	req := common.QuoteReq{
		ChainId:            137,
		Src:                TOKENB,
		Dst:                TOKENA,
		Amount:             big.NewInt(1000000),
		From:               "0x15Ba05723b04785C3E21157171810892A4FB795c",
		SlippagePercentage: 1,
		Referrer:           "0x15Ba05723b04785C3E21157171810892A4FB795c",
		SkipValidation:     true,
	}

	res, err := oneClient.FetchExactInSwapCallData(context.Background(), req)
	if err != nil {
		t.Fatalf("swap err: %v", err)
	}

	if res.FromToken.Address != TOKENB {
		t.Fatalf("invalid from token, token: %s", res.FromToken.Address)
	}

	if len(res.Tx.Data) < 1 {
		t.Fatalf("invalid txn data, tx: %s", res.Tx.Data)
	}
}
