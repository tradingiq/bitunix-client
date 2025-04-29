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

	params := model.PositionHistoryParams{
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
