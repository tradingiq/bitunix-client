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
	response, err := bitunixClient.GetPositionHistory(ctx)
	if err != nil {
		fmt.Println(err)
	}

	log.WithField("response", response).Debug("Get HistoricalPosition History")

	time.Sleep(5000 * time.Second)
}
