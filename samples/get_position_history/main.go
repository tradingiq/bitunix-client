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

	params := &bitunix.PositionHistoryParams{
		Limit: 2,
	}

	ctx := context.Background()
	response, err := bitunixClient.GetPositionHistory(ctx, params)
	if err != nil {
		fmt.Println(err)
	}

	log.WithField("response", response).Debug("Get HistoricalPosition History")

	time.Sleep(5000 * time.Second)
}
