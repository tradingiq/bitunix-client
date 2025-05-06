package main

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/bitunix"
	bitunix_errors "github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
)

type subTest struct {
	interval  model.Interval
	symbol    model.Symbol
	priceType model.PriceType
}

func (s subTest) SubscribeKLine(message *model.KLineChannelMessage) {
	log.WithField("message", message).Debug("KLine")
}

func (s subTest) SubscribeInterval() model.Interval {
	return s.interval
}

func (s subTest) SubscribeSymbol() model.Symbol {
	return s.symbol
}

func (s subTest) SubscribePriceType() model.PriceType {
	return s.priceType
}

func main() {
	ctx := context.Background()
	ws, _ := bitunix.NewPublicWebsocket(ctx)
	defer ws.Disconnect()
	log.SetLevel(log.DebugLevel)

	if err := ws.Connect(); err != nil {
		switch {
		case errors.Is(err, bitunix_errors.ErrNetwork):
			log.Fatalf("Network error: %s", err.Error())
		case errors.Is(err, bitunix_errors.ErrTimeout):
			log.Fatalf("Timeout error: %s", err.Error())
		case errors.Is(err, bitunix_errors.ErrInternal):
			log.Fatalf("Internal error: %s", err.Error())
		default:
			var sktErr *bitunix_errors.WebsocketError
			if errors.As(err, &sktErr) {
				log.Fatalf("Websocket error (code: %s): %s", sktErr.Operation, sktErr.Message)
			} else {
				log.Fatalf("Unexpected error: %s", err.Error())
			}
		}
	}

	interval, err := model.ParseInterval("1m")
	if err != nil {
		log.Fatalf("failed to parse interval: %v", err)
	}

	symbol := model.ParseSymbol("BTcUSDT")

	priceType, err := model.ParsePriceType("mark")
	if err != nil {
		log.Fatalf("failed to parse price type: %v", err)
	}

	sub := subTest{
		interval:  interval,
		symbol:    symbol,
		priceType: priceType,
	}

	err = ws.SubscribeKLine(sub)
	if err != nil {
		switch {
		case errors.Is(err, bitunix_errors.ErrValidation):
			log.WithError(err).Fatalf("Validation error: %s", err.Error())
		case errors.Is(err, bitunix_errors.ErrNetwork):
			log.WithError(err).Fatalf("Network error: %s", err.Error())
		case errors.Is(err, bitunix_errors.ErrInternal):
			log.WithError(err).Fatalf("Internal error: %s", err.Error())
		default:
			var apiErr *bitunix_errors.WebsocketError
			if errors.As(err, &apiErr) {
				log.Fatalf("Websocket error (code: %s): %s", apiErr.Operation, apiErr.Message)
			} else {
				log.Fatalf("Unexpected error: %s", err.Error())
			}
		}
	}

	if err := ws.Stream(); err != nil {
		switch {
		case errors.Is(err, bitunix_errors.ErrConnectionClosed):
			log.Info("connection closed, bye!")
		case errors.Is(err, bitunix_errors.ErrTimeout):
			log.WithError(err).Fatalf("timeout while streaming")
		default:
			var sktErr *bitunix_errors.WebsocketError
			if errors.As(err, &sktErr) {
				log.Fatalf("Websocket error (code: %s): %s", sktErr.Operation, sktErr.Message)
			} else {
				log.Fatalf("Unexpected error: %s", err.Error())
			}
		}
	}
}
