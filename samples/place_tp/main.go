package main

import (
	"bitunix-client/api"
	"bitunix-client/bitunix"
	"bitunix-client/samples"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
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

	ctx := context.Background()

	request := bitunix.NewTPSLOrderBuilder("BTCUSDT", "1988277607969914382").
		WithTakeProfit(100000, 0.001, bitunix.StopTypeLastPrice, bitunix.OrderTypeLimit, 100050).
		Build()

	response, err := bitunixClient.PlaceTpSlOrder(ctx, &request)
	if err != nil {
		fmt.Println(err)
	}

	log.WithField("response", response).Debug("place tpsl order")

	time.Sleep(5000 * time.Second)
}
