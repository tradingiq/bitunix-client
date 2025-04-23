package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/api"
	"github.com/tradingiq/bitunix-client/bitunix"
	"github.com/tradingiq/bitunix-client/samples"
	"time"
)

func main() {
	bitunixExample()
}

func bitunixExample() {
	log.SetLevel(log.DebugLevel)

	// Create API client
	apiClient, err := api.New("https://fapi.bitunix.com/", api.WithDebug(), api.WithDefaultTimeout(30*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	// Create Bitunix client
	client := bitunix.New(apiClient, samples.Config.ApiKey, samples.Config.SecretKey)

	// Set up parameters
	startTime := time.Now().Add(-80 * time.Hour) // Last 24 hours

	params := bitunix.OrderHistoryParams{
		Symbol:    "BTCUSDT",
		StartTime: &startTime,
		Limit:     10,
	}

	// Get order history
	ctx := context.Background()
	response, err := client.GetOrderHistory(ctx, params)
	if err != nil {
		log.Fatal(err)
	}

	// Log the full response for debugging
	log.WithField("response", response).Debug("Get Order History")

	// Display results if not in debug mode
	fmt.Printf("Found %d orders\n", len(response.Data.Orders))
	fmt.Printf("Total orders: %d\n\n", response.Data.Total)

	// Display order details
	for i, order := range response.Data.Orders {
		fmt.Printf("Order %d:\n", i+1)
		fmt.Printf("  OrderID: %s\n", order.OrderID)
		fmt.Printf("  Symbol: %s\n", order.Symbol)
		fmt.Printf("  Side: %s\n", order.Side)
		fmt.Printf("  Type: %s\n", order.OrderType)
		fmt.Printf("  Price: %s\n", order.Price)
		fmt.Printf("  Quantity: %.6f\n", order.Quantity)
		fmt.Printf("  Trade Quantity: %.6f\n", order.TradeQuantity)
		fmt.Printf("  Status: %s\n", order.Status)
		fmt.Printf("  Create Time: %s\n", order.CreateTime.Format(time.RFC3339))
		fmt.Printf("  Last Modified: %s\n", order.ModifyTime.Format(time.RFC3339))
		fmt.Printf("  Fee: %.6f\n", order.Fee)
		fmt.Printf("  Realized PNL: %.2f\n", order.RealizedPNL)
		fmt.Println()
	}

	// Wait to keep the program running for debugging
	time.Sleep(5 * time.Second)
}
