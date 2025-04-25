package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/bitunix"
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

	ctx := context.Background()

	request := bitunix.NewTPSLOrderBuilder("BTCUSDT", "50000").
		WithTakeProfit(80000, 0.001, bitunix.StopTypeLastPrice, bitunix.OrderTypeLimit, 70000).
		Build()

	response, err := bitunixClient.PlaceTpSlOrder(ctx, &request)
	if err != nil {
		fmt.Println(err)
	}

	log.WithField("response", response).Debug("place tpsl order")

	time.Sleep(5000 * time.Second)
}
