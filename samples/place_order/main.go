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
	apiClient, err := api.New("https://fapi.bitunix.com/", api.WithDebug(), api.WithDefaultTimeout(30*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	bitunixClient := bitunix.New(apiClient, samples.Config.ApiKey, samples.Config.SecretKey)

	limitOrder := bitunix.NewOrderBuilder(
		"BTCUSDT",
		bitunix.TradeActionSell,
		bitunix.TradeSideOpen,
		0.002,
	).WithOrderType(bitunix.OrderTypeLimit).
		WithPrice(100000.0).
		WithTimeInForce(bitunix.TimeInForcePostOnly).
		Build()

	ctx := context.Background()
	response, err := bitunixClient.PlaceOrder(ctx, &limitOrder)
	if err != nil {
		log.Fatalf("Failed to place tpsl order: %v", err)
	}

	fmt.Printf("tpsl Order placed successfully: %+v\n", response)

}
