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
