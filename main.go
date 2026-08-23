package main

import (
	"fmt"
	"log"
	"os"
	"github.com/blaskojan/coinmate-api-client/coinmate"
	"github.com/blaskojan/coinmate-api-client/coinmate/public"
	"github.com/blaskojan/coinmate-api-client/coinmate/secure"
)

func main() {
	// Get API credentials from environment variables
	clientID := os.Getenv("COINMATE_CLIENT_ID")
	apiKey := os.Getenv("COINMATE_API_KEY")
	privateKey := os.Getenv("COINMATE_PRIVATE_KEY")

	// Create client
	client := coinmate.GetCoinmateClient(clientID, apiKey, privateKey)

	fmt.Println("🚀 Coinmate API Client Demo")
	fmt.Println("============================")

	// Test public endpoints
	fmt.Println("\n📊 Testing Public Endpoints:")
	fmt.Println("-----------------------------")

	// Test ticker
	ticker := &public.Ticker{Client: client}
	tickerResponse, err := ticker.GetTicker("BTC_EUR")
	if err != nil {
		log.Printf("Ticker error: %v", err)
	} else {
		fmt.Printf("✅ Ticker: Last price: %.2f EUR\n", tickerResponse.Data.Last)
	}

	// Test trading pairs
	tradingPairs := &public.TradingPairs{Client: client}
	pairsResponse, err := tradingPairs.GetTradingPairs()
	if err != nil {
		log.Printf("Trading pairs error: %v", err)
	} else {
		fmt.Printf("✅ Trading pairs: Found %d pairs\n", len(pairsResponse.Data))
	}

	// Test order book
	orderBook := &public.OrderBook{Client: client}
	orderBookResponse, err := orderBook.GetOrderBook("BTC_EUR", false)
	if err != nil {
		log.Printf("Order book error: %v", err)
	} else {
		fmt.Printf("✅ Order book: %d asks, %d bids\n",
			len(orderBookResponse.Data.Asks),
			len(orderBookResponse.Data.Bids))
	}

	// Test transactions
	transactions := &public.Transactions{Client: client}
	transactionsResponse, err := transactions.GetTransactions("BTC_EUR", 60)
	if err != nil {
		log.Printf("Transactions error: %v", err)
	} else {
		fmt.Printf("✅ Transactions: Found %d recent transactions\n", len(transactionsResponse.Data))
	}

	// Test secure endpoints (only if credentials are provided)
	if clientID != "" && apiKey != "" && privateKey != "" {
		fmt.Println("\n🔒 Testing Secure Endpoints:")
		fmt.Println("----------------------------")

		// Test balances
		balances := &secure.Balances{Client: client}
		balancesResponse, err := balances.GetBalances()
		if err != nil {
			log.Printf("Balances error: %v", err)
		} else {
			fmt.Printf("✅ Balances: Found %d currencies\n", len(balancesResponse.Data))
			for currency, balance := range balancesResponse.Data {
				fmt.Printf("   %s: %.8f (available: %.8f)\n",
					currency, balance.Balance, balance.Available)
			}
		}

		// Test order history
		order := &secure.Order{Client: client}
		orderHistoryResponse, err := order.GetHistory("BTC_EUR", 10)
		if err != nil {
			log.Printf("Order history error: %v", err)
		} else {
			fmt.Printf("✅ Order history: Found %d orders\n", len(orderHistoryResponse.Data))
		}

		// Test open orders
		openOrdersResponse, err := order.GetOpenOrders("BTC_EUR")
		if err != nil {
			log.Printf("Open orders error: %v", err)
		} else {
			fmt.Printf("✅ Open orders: Found %d open orders\n", len(openOrdersResponse.Data))
		}
	} else {
		fmt.Println("\n⚠️  Skipping secure endpoints (no API credentials provided)")
		fmt.Println("   Set COINMATE_CLIENT_ID, COINMATE_API_KEY, and COINMATE_PRIVATE_KEY environment variables")
	}

	fmt.Println("\n✅ Demo completed!")
}
