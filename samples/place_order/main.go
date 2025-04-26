package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/bitunix"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/rest"
	"github.com/tradingiq/bitunix-client/samples"
	"time"
)

func main() {
	bitunixExample()
}

func bitunixExample() {
	log.SetLevel(log.DebugLevel)
	apiClient, err := rest.New("https://fapi.bitunix.com/", rest.WithDebug(), rest.WithDefaultTimeout(30*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	bitunixClient := bitunix.New(apiClient, samples.Config.ApiKey, samples.Config.SecretKey)

	limitOrder := bitunix.NewOrderBuilder(
		"BTCUSDT",
		model.TradeSideSell,
		model.SideOpen,
		0.002,
	).WithOrderType(model.OrderTypeLimit).
		WithPrice(100000.0).
		WithTimeInForce(model.TimeInForcePostOnly).
		Build()

	ctx := context.Background()
	response, err := bitunixClient.PlaceOrder(ctx, &limitOrder)
	if err != nil {
		log.Fatalf("Failed to place tpsl order: %v", err)
	}

	fmt.Printf("tpsl Order placed successfully: %+v\n", response)

}
