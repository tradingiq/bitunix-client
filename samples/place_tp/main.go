package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/bitunix"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/samples"
	"time"
)

func main() {
	bitunixExample()
}

func bitunixExample() {
	log.SetLevel(log.DebugLevel)

	bitunixClient, _ := bitunix.NewApiClient(samples.Config.ApiKey, samples.Config.SecretKey)

	ctx := context.Background()

	requestBuilder := bitunix.NewTPSLOrderBuilder("BTCUSDT", "50000").
		WithTakeProfit(80000, 0.001, model.StopTypeLastPrice, model.OrderTypeLimit, 70000)
		
	request, err := requestBuilder.Build()
	if err != nil {
		log.Fatalf("Failed to build TPSL order: %v", err)
	}

	response, err := bitunixClient.PlaceTpSlOrder(ctx, &request)
	if err != nil {
		fmt.Println(err)
	}

	log.WithField("response", response).Debug("place tpsl order")

	time.Sleep(5000 * time.Second)
}
