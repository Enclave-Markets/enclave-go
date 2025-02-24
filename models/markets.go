package models

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type (
	Market string
	Symbol string
)

type V1GetMarketsResult struct {
	// Spot markets that the user is allowed to trade in
	Spot            SpotMarkets             `json:"spot"`
	PerpetualFuture *PerpetualFutureMarkets `json:"perps,omitempty"`
}

type PerpsConfig struct {
	Market           Market          `json:"market"`
	UnderlyingMarket Market          `json:"underlyingMarket"`
	Pair             CurrencyPair    `json:"pair"`
	BaseIncrement    decimal.Decimal `json:"baseIncrement"`
	QuoteIncrement   decimal.Decimal `json:"quoteIncrement"`
}
type PerpetualFutureMarkets struct {
	TradingPairs []PerpsConfig `json:"tradingPairs"`
}
type SpotMarkets struct {
	TradingPairs []V1SpotMarketsResult `json:"tradingPairs"`
}

type CurrencyPair struct {
	Base  string `json:"base"`
	Quote string `json:"quote"`
}

type V1SpotMarketsResult struct {
	Market         Market          `json:"market"`
	BaseIncrement  decimal.Decimal `json:"baseIncrement"`
	Pair           *CurrencyPair   `json:"pair"`
	QuoteIncrement decimal.Decimal `json:"quoteIncrement"`
	Disabled       bool            `json:"disabled,omitempty"`
}

// swagger:model
type BookSnapshot struct {
	// best n bids in the market
	// required:true
	// example:[["21.05", "0.34"], ["21.02", "1.25"]]
	Bids []BookLevel `json:"bids"`

	// best n asks in the market
	// required:true
	// example:[["21.11", "1.74"], ["21.13", "0.23"]]
	Asks []BookLevel `json:"asks"`
}

type BookLevel struct {
	Price    decimal.Decimal `json:"price"`
	Quantity decimal.Decimal `json:"size"`
}

// Output in the standard format `[price, quantity]` instead of `{"price": price, "quantity": quantity}`
func (b BookLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{b.Price, b.Quantity})
}

func (b *BookLevel) UnmarshalJSON(data []byte) error {
	var v []decimal.Decimal
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	if len(v) != 2 {
		return fmt.Errorf("expected 2 elements, got %d", len(v))
	}
	b.Price = v[0]
	b.Quantity = v[1]
	return nil
}

func (m Market) AsPair() (CurrencyPair, error) {
	parts := strings.Split(string(m), "-")
	if len(parts) != 2 {
		return CurrencyPair{}, fmt.Errorf("failed CurrencyPairFromString, expected 2, got %d parts.  Input: '%s'", len(parts), string(m))
	}
	return CurrencyPair{
		Base:  parts[0],
		Quote: parts[1],
	}, nil
}
