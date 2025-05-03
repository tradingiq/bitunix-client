package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/bitunix"
	"github.com/tradingiq/bitunix-client/model"
)

type subTest struct {
	interval  model.Interval
	symbol    model.Symbol
	priceType model.PriceType
}

func (s subTest) Handle(message *model.KLineChannelMessage) {
	log.WithField("message", message).Debug("KLine")
}

func (s subTest) Interval() model.Interval {
	return s.interval
}

func (s subTest) Symbol() model.Symbol {
	return s.symbol
}

func (s subTest) PriceType() model.PriceType {
	return s.priceType
}

func main() {
	ctx := context.Background()
	ws := bitunix.NewPublicWebsocket(ctx)
	defer ws.Disconnect()
	log.SetLevel(log.DebugLevel)

	if err := ws.Connect(); err != nil {
		log.Fatalf("failed to connect to WebSocket: %v", err)
	}

	sub := subTest{
		interval:  "1min",
		symbol:    "BTCUSDT",
		priceType: "market",
	}
	err := ws.SubscribeKLine(sub)
	if err != nil {
		log.Fatalf("failed to subscribe: %v", err)
	}

	if err := ws.Stream(); err != nil {
		log.WithError(err).Fatal("failed to stream")
	}
}
