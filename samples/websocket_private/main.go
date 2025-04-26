package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/bitunix"
	"github.com/tradingiq/bitunix-client/samples"
)

func main() {
	ctx := context.Background()
	ws := bitunix.NewPrivateWebsocket(ctx, samples.Config.ApiKey, samples.Config.SecretKey)
	defer ws.Disconnect()

	if err := ws.Connect(); err != nil {
		log.Fatalf("failed to connect to WebSocket: %v", err)
	}

	balance := ws.SubscribeBalance()
	go func() {
		for balanceResponse := range balance {
			log.WithField("balance", balanceResponse).Debug("got balance")
		}
	}()

	positions := ws.SubscribePositions()
	go func() {
		for positionResponse := range positions {
			log.WithField("position", positionResponse).Debug("got position")
		}
	}()

	orders := ws.SubscribeOrders()
	go func() {
		for res := range orders {
			log.WithField("order", res).Debug("got order")
		}
	}()

	if err := ws.Stream(); err != nil {
		log.WithError(err).Fatal("failed to stream")
	}

}
