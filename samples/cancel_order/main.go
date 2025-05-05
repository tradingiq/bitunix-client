package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/bitunix"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/samples"
)

func main() {
	cancelOrderExample()
}

func cancelOrderExample() {
	log.SetLevel(log.DebugLevel)

	bitunixClient, _ := bitunix.NewApiClient(samples.Config.ApiKey, samples.Config.SecretKey)

	cancelRequestBuilder := bitunix.NewCancelOrderBuilder(model.ParseSymbol("BTCUSDT")).
		WithOrderID("1915122868439269376")
		
	cancelRequest, err := cancelRequestBuilder.Build()
	if err != nil {
		log.Fatalf("Failed to build cancel order request: %v", err)
	}

	ctx := context.Background()
	response, err := bitunixClient.CancelOrders(ctx, &cancelRequest)
	if err != nil {
		log.Fatalf("Failed to cancel orders: %v", err)
	}

	fmt.Println("Successfully canceled orders:")
	for _, success := range response.Data.SuccessList {
		fmt.Printf("- Order ID: %s, client ID: %s\n", success.OrderId, success.ClientId)
	}

	fmt.Println("Failed to cancel orders:")
	for _, failure := range response.Data.FailureList {
		fmt.Printf("- Order ID: %s, client ID: %s, Error: %s (Code: %s)\n",
			failure.OrderId, failure.ClientId, failure.ErrorMsg, failure.ErrorCode)
	}
}
