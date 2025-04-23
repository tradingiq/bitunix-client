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
	cancelOrderExample()
}

func cancelOrderExample() {
	log.SetLevel(log.DebugLevel)
	apiClient, err := api.New("https://fapi.bitunix.com/", api.WithDebug(), api.WithDefaultTimeout(30*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	bitunixClient := bitunix.New(apiClient, samples.Config.ApiKey, samples.Config.SecretKey)

	// Create a request to cancel orders by providing either order IDs or client IDs
	cancelRequest := bitunix.NewCancelOrderBuilder("BTCUSDT").
		WithOrderID("11111").        // Cancel by order ID
		WithClientID("client22222"). // Cancel by client ID
		Build()

	ctx := context.Background()
	response, err := bitunixClient.CancelOrders(ctx, &cancelRequest)
	if err != nil {
		log.Fatalf("Failed to cancel orders: %v", err)
	}

	// Print successful cancellations
	fmt.Println("Successfully canceled orders:")
	for _, success := range response.Data.SuccessList {
		fmt.Printf("- Order ID: %s, Client ID: %s\n", success.OrderId, success.ClientId)
	}

	// Print failed cancellations
	fmt.Println("Failed to cancel orders:")
	for _, failure := range response.Data.FailureList {
		fmt.Printf("- Order ID: %s, Client ID: %s, Error: %s (Code: %s)\n",
			failure.OrderId, failure.ClientId, failure.ErrorMsg, failure.ErrorCode)
	}
}
