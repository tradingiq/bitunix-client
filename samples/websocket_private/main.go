package main

import (
	"context"
	"errors"
	"time"

	"github.com/tradingiq/bitunix-client/bitunix"
	bitunix_errors "github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/samples"
	"go.uber.org/zap"
)

type BalanceHandler struct {
	logger *zap.Logger
}

func (h *BalanceHandler) SubscribeBalance(balance *model.BalanceChannelMessage) {
	h.logger.Debug("got balance", zap.Any("balance", balance))
}

type PositionHandler struct {
	logger *zap.Logger
}

func (h *PositionHandler) SubscribePosition(position *model.PositionChannelMessage) {
	h.logger.Debug("got position", zap.Any("position", position))
}

type OrderHandler struct {
	logger *zap.Logger
}

func (h *OrderHandler) SubscribeOrder(order *model.OrderChannelMessage) {
	h.logger.Debug("got order", zap.Any("order", order))
}

type TpSlOrderHandler struct {
	logger *zap.Logger
}

func (h *TpSlOrderHandler) SubscribeTpSlOrder(order *model.TpSlOrderChannelMessage) {
	h.logger.Debug("got tpsl order", zap.Any("tpslorder", order))
}

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	ctx := context.Background()

	// Create reconnecting private websocket client
	ws, err := bitunix.NewReconnectingPrivateWebsocket(ctx, samples.Config.ApiKey, samples.Config.SecretKey,
		bitunix.WithPrivateWebsocketOptions(bitunix.WithWebsocketLogLevel(model.LogLevelVeryAggressive)),
		bitunix.WithPrivateMaxReconnectAttempts(0),       // Infinite reconnect attempts
		bitunix.WithPrivateReconnectDelay(5*time.Second), // 5 second delay between attempts
		bitunix.WithPrivateReconnectLogger(logger),       // Use logger for reconnection events
	)
	if err != nil {
		logger.Fatal("Failed to create reconnecting private websocket client", zap.Error(err))
	}
	defer ws.Disconnect()

	logger.Info("Connecting to private websocket with automatic reconnection enabled")
	if err := ws.Connect(); err != nil {
		switch {
		case errors.Is(err, bitunix_errors.ErrAuthentication):
			logger.Fatal("Authentication failed", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrSignatureError):
			logger.Fatal("Signature error", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrNetwork):
			logger.Fatal("Network error", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrIPNotAllowed):
			logger.Fatal("IP restriction error", zap.Error(err))
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

	// Subscribe to all private channels - these will be automatically restored on reconnection
	balanceHandler := &BalanceHandler{logger: logger}
	if err := ws.SubscribeBalance(balanceHandler); err != nil {
		handleError(err, "balance", logger)
	}
	logger.Info("Subscribed to balance updates")

	positionHandler := &PositionHandler{logger: logger}
	if err := ws.SubscribePositions(positionHandler); err != nil {
		handleError(err, "positions", logger)
	}
	logger.Info("Subscribed to position updates")

	orderHandler := &OrderHandler{logger: logger}
	if err := ws.SubscribeOrders(orderHandler); err != nil {
		handleError(err, "orders", logger)
	}
	logger.Info("Subscribed to order updates")

	tpslHandler := &TpSlOrderHandler{logger: logger}
	if err := ws.SubscribeTpSlOrders(tpslHandler); err != nil {
		handleError(err, "tpsl orders", logger)
	}
	logger.Info("Subscribed to TP/SL order updates")

	logger.Info("Starting websocket stream with automatic reconnection - the connection will automatically reconnect if it fails")

	// Start streaming - this will automatically handle reconnections
	if err := ws.Stream(); err != nil {
		switch {
		case errors.Is(err, bitunix_errors.ErrTimeout):
			logger.Fatal("timeout while streaming", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrConnectionClosed):
			logger.Info("connection closed", zap.Error(err))
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
}

func handleError(err error, name string, logger *zap.Logger) {
	switch {
	case errors.Is(err, bitunix_errors.ErrValidation):
		logger.Fatal("failed to subscribe", zap.String("name", name), zap.Error(err))
	default:
		var apiErr *bitunix_errors.WebsocketError
		if errors.As(err, &apiErr) {
			logger.Fatal("Websocket error", zap.String("operation", apiErr.Operation), zap.String("message", apiErr.Message))
		} else {
			logger.Fatal("Unexpected error", zap.Error(err))
		}
	}
}
