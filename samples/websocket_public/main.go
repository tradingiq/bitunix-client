package main

import (
	"context"
	"errors"
	"time"

	"github.com/tradingiq/bitunix-client/bitunix"
	bitunix_errors "github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
	"go.uber.org/zap"
)

type subTest struct {
	interval  model.Interval
	symbol    model.Symbol
	priceType model.PriceType
	logger    *zap.Logger
}

func (s subTest) SubscribeKLine(message *model.KLineChannelMessage) {
	s.logger.Debug("KLine", zap.Any("message", message))
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
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	ctx := context.Background()

	// Create base public websocket client
	baseWs, err := bitunix.NewPublicWebsocket(ctx, bitunix.WithWebsocketLogLevel(model.LogLevelVeryAggressive))
	if err != nil {
		logger.Fatal("Failed to create base public websocket client", zap.Error(err))
	}

	// Create reconnecting public websocket client
	ws := bitunix.NewReconnectingPublicWebsocket(ctx, baseWs,
		bitunix.WithMaxReconnectAttempts(0),       // Infinite reconnect attempts
		bitunix.WithReconnectDelay(5*time.Second), // 5 second delay between attempts
		bitunix.WithReconnectLogger(logger),       // Use logger for reconnection events
	)
	defer ws.Disconnect()

	logger.Info("Connecting to public websocket with automatic reconnection enabled")
	if err := ws.Connect(); err != nil {
		switch {
		case errors.Is(err, bitunix_errors.ErrNetwork):
			logger.Fatal("Network error", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrTimeout):
			logger.Fatal("Timeout error", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrInternal):
			logger.Fatal("Internal error", zap.Error(err))
		default:
			var sktErr *bitunix_errors.WebsocketError
			if errors.As(err, &sktErr) {
				logger.Fatal("Websocket error", zap.String("operation", sktErr.Operation), zap.String("message", sktErr.Message))
			} else {
				logger.Fatal("Unexpected error", zap.Error(err))
			}
		}
	}

	interval, err := model.ParseInterval("1m")
	if err != nil {
		logger.Fatal("failed to parse interval", zap.Error(err))
	}

	symbol := model.ParseSymbol("BTCUSDT")

	priceType, err := model.ParsePriceType("mark")
	if err != nil {
		logger.Fatal("failed to parse price type", zap.Error(err))
	}

	sub := subTest{
		interval:  interval,
		symbol:    symbol,
		priceType: priceType,
		logger:    logger,
	}

	// Subscribe to KLine - this subscription will be automatically restored on reconnection
	err = ws.SubscribeKLine(sub)
	if err != nil {
		switch {
		case errors.Is(err, bitunix_errors.ErrValidation):
			logger.Fatal("Validation error", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrNetwork):
			logger.Fatal("Network error", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrInternal):
			logger.Fatal("Internal error", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrWorkgroupExhausted):
			logger.Info("workgroup is exhausted", zap.Error(err))
		default:
			var apiErr *bitunix_errors.WebsocketError
			if errors.As(err, &apiErr) {
				logger.Fatal("Websocket error", zap.String("operation", apiErr.Operation), zap.String("message", apiErr.Message))
			} else {
				logger.Fatal("Unexpected error", zap.Error(err))
			}
		}
	}
	logger.Info("Subscribed to KLine updates", zap.String("symbol", symbol.String()), zap.String("interval", interval.String()), zap.String("priceType", priceType.String()))

	logger.Info("Starting websocket stream with automatic reconnection - the connection will automatically reconnect if it fails")

	// Start streaming - this will automatically handle reconnections
	if err := ws.Stream(); err != nil {
		switch {
		case errors.Is(err, bitunix_errors.ErrConnectionClosed):
			logger.Info("connection closed, bye!")
		case errors.Is(err, bitunix_errors.ErrTimeout):
			logger.Fatal("timeout while streaming", zap.Error(err))
		default:
			var sktErr *bitunix_errors.WebsocketError
			if errors.As(err, &sktErr) {
				logger.Fatal("Websocket error", zap.String("operation", sktErr.Operation), zap.String("message", sktErr.Message))
			} else {
				logger.Fatal("Unexpected error", zap.Error(err))
			}
		}
	}
}
