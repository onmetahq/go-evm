package dex

import (
	"context"
	"log/slog"
	"math/big"
	"testing"
	"time"

	metahttp "github.com/onmetahq/meta-http/pkg/meta_http"
)

func Test0XFetchExactInQuote(t *testing.T) {
	c := metahttp.NewClient("", slog.Default(), 30*time.Second)

	oneClient := NewZeroX(c, map[uint64]string{
		137: "https://polygon.api.0x.org",
	})

	req := QuoteReq{
		Src:            TOKENB,
		Dst:            TOKENA,
		ChainId:        137,
		Amount:         big.NewInt(1000000),
		SkipValidation: true,
	}
	res, err := oneClient.FetchExactInQuote(context.Background(), req)
	if err != nil {
		t.Fatalf("quote err: %v", err)
	}

	if res.ToAmount.Cmp(big.NewInt(0)) < 1 {
		t.Fatalf("invalid quote, quote: %s", res.ToAmount)
	}
}

func Test0XFetchExactOutQuote(t *testing.T) {
	c := metahttp.NewClient("", slog.Default(), 30*time.Second)

	oneClient := NewZeroX(c, map[uint64]string{
		137: "https://polygon.api.0x.org",
	})

	req := QuoteReq{
		Src:            TOKENB,
		Dst:            TOKENA,
		ChainId:        137,
		Amount:         big.NewInt(1e18),
		SkipValidation: true,
	}
	res, err := oneClient.FetchExactOutQuote(context.Background(), req)
	if err != nil {
		t.Fatalf("quote err: %v", err)
	}

	if res.ToAmount.Cmp(big.NewInt(0)) < 1 {
		t.Fatalf("invalid quote, quote: %s", res.ToAmount)
	}

	// This will fail if 1 matic is more than 1 usd
	if res.ToAmount.Cmp(big.NewInt(1e6)) > 0 {
		t.Fatalf("invalid quote matic is more than 1 usd, quote: %s", res.ToAmount)
	}
}

func Test0XFetchExactInSwapCallData(t *testing.T) {
	c := metahttp.NewClient("", slog.Default(), 30*time.Second)

	oneClient := NewZeroX(c, map[uint64]string{
		137: "https://polygon.api.0x.org",
	})

	req := QuoteReq{
		Src:            TOKENB,
		Dst:            TOKENA,
		ChainId:        137,
		Amount:         big.NewInt(1e18),
		SkipValidation: true,
	}

	res, err := oneClient.FetchExactInSwapCallData(context.Background(), req)
	if err != nil {
		t.Fatalf("quote err: %v", err)
	}

	if len(res.Data) == 0 {
		t.Fatalf("invalid call data, res: %v", res)
	}
}
