package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/bitunix"
	"github.com/tradingiq/bitunix-client/model"
)

func main() {
	ctx := context.Background()
	ws := bitunix.NewPublicWebsocket(ctx)
	defer ws.Disconnect()
	log.SetLevel(log.DebugLevel)

	if err := ws.Connect(); err != nil {
		log.Fatalf("failed to connect to WebSocket: %v", err)
	}

	err := ws.SubscribeKLine("BTCUSDT", "1min", "market", func(msg *model.KLineChannelMessage) {
		log.WithField("candle", msg).Debug("KLine")
	})
	if err != nil {
		log.Fatalf("failed to subscribe: %v", err)
	}

	if err := ws.Stream(); err != nil {
		log.WithError(err).Fatal("failed to stream")
	}
}
