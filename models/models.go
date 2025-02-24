package models

import (
	"time"

	"github.com/shopspring/decimal"
)

const (

	// Tool
	StatusPath       = "/status"
	HelloPath        = "/hello"
	AuthedHelloPath  = "/authedHello"
	V0GetBalancePath = "/v0/get_balance"

	// Markets
	V1MarketsPath = "/v1/markets"

	// Spot trading
	V1SpotOrdersPath      = "/v1/orders"
	V1SpotBatchOrdersPath = "/v1/orders/batch"
	V1SpotFillsPath       = "/v1/fills"
	V1SpotDepthPath       = "/v1/depth"

	V1SpotClientOrderIDPrefix = "client:"

	// Perps trading
	V1PerpsOrdersPath      = "/v1/perps/orders"
	V1PerpsBatchOrdersPath = "/v1/perps/orders/batch"

	// Cross
	V0PricePath = "/v0/price"
)

type V1PageRes[T any] struct {
	Result   []*T
	PageInfo APIPageInfo
}

type APIPageInfo struct {
	PrevCursor string `json:"prevCursor"`
	NextCursor string `json:"nextCursor"`
}

type AccountID string

type GetPublicStatusRes struct {
	MarketStatuses map[Market]string `json:"marketStatuses"`
}

type GenericResponse[T any] struct {
	Success bool   `json:"success"`
	Result  T      `json:"result"`
	Error   string `json:"error,omitempty"`
}

type V0GetBalanceRes struct {
	// the account ID of the user that made the request
	// example:5577006791947779410
	// required:true
	AccountId AccountID `json:"accountId"`

	// the coin for which the customer requested balance
	// example:AVAX
	// required:true
	Symbol Symbol `json:"symbol"`

	// the total balance of the coin
	// example:10000
	// required:true
	TotalBalance string `json:"totalBalance"`

	// the reserved balance of the coin, held in open orders
	// example:7000
	// required:true
	ReservedBalance string `json:"reservedBalance"`

	// the free balance of the coin
	// example:3000
	// required:true
	FreeBalance string `json:"freeBalance"`
}

type GetBalanceReq struct {
	// the coin for which the customer wants to get balance
	// example:AVAX
	// required:true
	Symbol Symbol `json:"symbol"`
}

type GetMarkPriceRes struct {
	// trading market pair
	// example:AVAX-USD.PERP
	// required:true
	Market    Market          `json:"market"`
	MarkPrice decimal.Decimal `json:"markPrice"`
}

type ApiPosition struct {
	Market                 Market           `json:"market"`
	Direction              string           `json:"direction"`
	NetQuantity            decimal.Decimal  `json:"netQuantity"`
	AverageEntryPrice      decimal.Decimal  `json:"averageEntryPrice"`
	UsedMargin             decimal.Decimal  `json:"usedMargin"`
	UnrealizedPnl          decimal.Decimal  `json:"unrealizedPnl"`
	MarkPrice              decimal.Decimal  `json:"markPrice"`
	LiquidationPrice       decimal.Decimal  `json:"liquidationPrice"`
	BankruptcyPrice        decimal.Decimal  `json:"bankruptcyPrice"`
	MaintenanceMargin      decimal.Decimal  `json:"maintenanceMargin"`
	NotionalValue          decimal.Decimal  `json:"notionalValue"`
	Leverage               decimal.Decimal  `json:"leverage"`
	NetFundingSinceNeutral decimal.Decimal  `json:"netFundingSinceNeutral"`
	StopLossTriggerPrice   *decimal.Decimal `json:"stopLossTriggerPrice,omitempty"`
	TakeProfitTriggerPrice *decimal.Decimal `json:"takeProfitTriggerPrice,omitempty"`
}

type ApiBookSnapshots []*ApiBookSnapshot

type ApiBookSnapshot struct {
	Market Market      `json:"market"`
	Time   time.Time   `json:"time"`
	Asks   []BookLevel `json:"asks"`
	Bids   []BookLevel `json:"bids"`
}

type GetPriceReq struct {
	Pair CurrencyPair `json:"pair"`
}

type V0GetPriceRes struct {
	Pair      CurrencyPair    `json:"pair"`
	Available bool            `json:"available"`
	Price     decimal.Decimal `json:"price"`
	QuotedAt  string          `json:"quotedAt"`
}
