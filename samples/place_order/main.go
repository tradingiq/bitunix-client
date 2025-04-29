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
	bitunixExample()
}

func bitunixExample() {
	log.SetLevel(log.DebugLevel)

	bitunixClient, _ := bitunix.NewApiClient(samples.Config.ApiKey, samples.Config.SecretKey)

	limitOrder := bitunix.NewOrderBuilder(
		model.ParseSymbol("BTCUSDT"),
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
