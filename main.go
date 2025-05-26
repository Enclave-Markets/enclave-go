package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Enclave-Markets/enclave-go/apiclient"
	"github.com/Enclave-Markets/enclave-go/models"
	"github.com/shopspring/decimal"
)

func main() {
	// the trading pair that will be used, and balance symbol to check
	symbol := models.Symbol("AVAX")
	market := models.Market("AVAX-USDC")

	client, err := apiclient.NewApiClientFromEnv("sandbox")
	if err != nil {
		fmt.Println("failed to create client", err)
		return
	}

	client.WithApiKey(
		os.Getenv("ENCLAVE_KEY"),
		os.Getenv("ENCLAVE_SECRET"),
	)

	client.WaitForEndpoint()
	_, err = client.AuthedHello()
	if err != nil {
		fmt.Println("authed-hello failed:", err)
		return
	}

	// get the balance of AVAX
	balanceResp, err := client.GetBalance(models.GetBalanceReq{Symbol: symbol})
	if err != nil {
		fmt.Println("failed to get balance:", err)
		return
	}
	fmt.Println(symbol, "balance:", balanceResp.Result.TotalBalance)

	// get the AVAX-USDC trading pair to find the min order sizes
	marketsResp, err := client.Markets()
	if err != nil {
		fmt.Println("failed to get markets", err, marketsResp.Success)
		return
	}

	var baseMin, quoteMin decimal.Decimal
	for _, pair := range marketsResp.Result.Spot.TradingPairs {
		if pair.Market != market {
			continue
		}
		baseMin, quoteMin = pair.BaseIncrement, pair.QuoteIncrement
	}
	fmt.Println("base-min:", baseMin, "quote-min:", quoteMin)

	// get top of book for avax usdc
	book, err := client.GetSpotDepthBook(market)
	if err != nil {
		fmt.Println("failed to get market depth:", err)
		return
	}

	if len(book.Result.Asks) == 0 {
		fmt.Println("no asks are resting on the book")
		return
	}
	bestAsk := book.Result.Asks[0]
	fmt.Println("best-ask-price:", bestAsk.Price, "best-ask-size:", bestAsk.Quantity)

	// place a sell limit order of the smallest size one tick above the top of book (so we don't get filled)
	orderResp, err := client.AddSpotOrder(
		models.AddOrderReq{
			Market: market,
			Side:   models.Ask,
			Price:  bestAsk.Price.Add(quoteMin),
			Size:   baseMin,
			Type:   models.OrderTypeLimit,
		},
	)
	if err != nil {
		fmt.Println("failed to place order:", err)
		return
	}
	fmt.Println("order placed, current state:", orderResp.Result.State)

	// cancel all orders in the market
	if err := client.CancelAllSpotOrders(); err != nil {
		fmt.Println("failed to cancel spot orders:", err)
		return
	}

	// get the order state
	orderResp, err = client.GetSpotOrder(orderResp.Result.OrderID)
	if err != nil {
		fmt.Println("failed to get spot order:", err)
		return
	}

	fmt.Println("order state:", orderResp.Result.State)

	var side models.BidAsk
	switch {
	case len(book.Result.Asks) != 0:
		side = models.Bid
	case len(book.Result.Bids) != 0:
		side = models.Ask
	default:
		fmt.Println("orderbook is empty -- cannot trigger a fill")
		return
	}

	clientOrderID := models.OrderID(strconv.FormatUint(uint64(time.Now().UnixNano()), 10))
	orderResp, err = client.AddSpotOrder(
		models.AddOrderReq{
			Market:        market,
			Side:          side,
			Size:          baseMin,
			Type:          models.OrderTypeMarket,
			ClientOrderID: clientOrderID,
		},
	)
	if err != nil {
		fmt.Println("failed to place market order:", err)
		return
	}

	if orderResp.Result.State != models.FullyFilled {
		fmt.Println("market order did not fill:", orderResp.Result.State)
		return
	}
	fmt.Println("market order placed, current state:", orderResp.Result.State)

	// lets find the fills
	fillsResp, err := client.GetSpotFillsByOrderID(orderResp.Result.OrderID)
	if err != nil {
		fmt.Println("unable to find fill:", err)
		return
	}
	if len(fillsResp.Result) == 0 {
		fmt.Println("could not find fills!")
	}
	fmt.Println("found n-fills for market order:", len(fillsResp.Result))

	// lets find the fills by client order ID this time
	fillsResp, err = client.GetSpotFillsByClientOrderID(orderResp.Result.ClientOrderID)
	if err != nil {
		fmt.Println("unable to find fill:", err)
		return
	}
	if len(fillsResp.Result) == 0 {
		fmt.Println("could not find fills!")
	}
	fmt.Println("found n-fills by client order ID for market order:", len(fillsResp.Result))

	// now find all fills
	allFillsResp, err := client.GetSpotFills(models.FillParams{})
	if err != nil {
		fmt.Println("unable to find all fills:", err)
		return
	}
	if len(allFillsResp.Result) == 0 {
		fmt.Println("could not find any fills!")
	}
	fmt.Println("found n-fills any orders:", len(allFillsResp.Result))

	targetMarket := models.Market("AVAX-USD.P")
	perpsResp, err := client.GetPerpsContracts()
	if err != nil {
		fmt.Printf("Failed to get perps contracts: %v\n", err)
	} else {
		contracts := perpsResp.Result
		fmt.Printf("Found %d perps contracts\n", len(contracts))

		// Display summary of all contracts
		fmt.Println("\nSummary of all contracts:")
		fmt.Println("----------------------------------------------------")
		fmt.Printf("%-12s | %-8s | %-8s | %-8s\n", "MARKET", "STATUS", "DISABLED", "LAST PRICE")
		fmt.Println("----------------------------------------------------")
		for _, contract := range contracts {
			status := "OPEN"
			if contract.IsClosed {
				status = "CLOSED"
			}
			fmt.Printf("%-12s | %-8s | %-8t | %-8s\n",
				contract.Market,
				status,
				contract.Disabled,
				contract.LastPrice)
		}

		// Check specific market
		fmt.Printf("\nLooking for target market: %s\n", targetMarket)
		var found bool
		for _, contract := range contracts {
			if contract.Market == targetMarket {
				found = true
				status := "OPEN"
				if contract.IsClosed {
					status = "CLOSED"
				}

				// Print detailed information about the target market
				fmt.Println("\nDetailed Information:")
				fmt.Println("----------------------------------------------------")
				fmt.Printf("Market:               %s\n", contract.Market)
				fmt.Printf("Status:               %s\n", status)
				fmt.Printf("Base/Quote:           %s/%s\n", contract.BaseCurrency, contract.QuoteCurrency)
				fmt.Printf("Contract Type:        %s\n", contract.ContractType)
				fmt.Printf("Last Price:           %s\n", contract.LastPrice)
				fmt.Printf("Bid/Ask:              %s/%s\n", contract.Bid, contract.Ask)
				fmt.Printf("24h High/Low:         %s/%s\n", contract.High, contract.Low)
				fmt.Printf("24h Change:           %s%%\n", contract.PriceChangePercent)
				fmt.Printf("Open Interest:        %s (%s USD)\n", contract.OpenInterest, contract.OpenInterestUsd)
				fmt.Printf("Funding Rate:         %s\n", contract.FundingRate)
				fmt.Printf("Next Funding:         %s\n", contract.NextFundingRateTimestamp)
				fmt.Printf("Maker/Taker Fee:      %s/%s\n", contract.MakerFee, contract.TakerFee)
				fmt.Printf("Tags:                 %v\n", contract.Tags)
				break
			}
		}

		if !found {
			fmt.Printf("\nTarget market %s not found in contracts\n", targetMarket)
		}
	}
}
