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

	params := model.TradeHistoryParams{
		PositionID: "8623335020855713906",
	}

	ctx := context.Background()
	response, err := bitunixClient.GetTradeHistory(ctx, params)
	if err != nil {
		log.Fatalf("Failed to place tpsl order: %v", err)
	}

	fmt.Printf("tpsl Order placed successfully: %+v\n", response)

}
