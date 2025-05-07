package main

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/bitunix"
	bitunix_errors "github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/samples"
)

type BalanceHandler struct{}

func (h *BalanceHandler) SubscribeBalance(balance *model.BalanceChannelMessage) {
	log.WithField("balance", balance).Debug("got balance")
}

type PositionHandler struct{}

func (h *PositionHandler) SubscribePosition(position *model.PositionChannelMessage) {
	log.WithField("position", position).Debug("got position")
}

type OrderHandler struct{}

func (h *OrderHandler) SubscribeOrder(order *model.OrderChannelMessage) {
	log.WithField("order", order).Debug("got order")
}

type TpSlOrderHandler struct{}

func (h *TpSlOrderHandler) SubscribeTpSlOrder(order *model.TpSlOrderChannelMessage) {
	log.WithField("tpslorder", order).Debug("got tpsl order")
}

func main() {
	ctx := context.Background()
	ws, _ := bitunix.NewPrivateWebsocket(ctx, samples.Config.ApiKey, samples.Config.SecretKey)
	defer ws.Disconnect()

	if err := ws.Connect(); err != nil {
		switch {
		case errors.Is(err, bitunix_errors.ErrAuthentication):
			log.Fatalf("Authentication failed: %s", err.Error())
		case errors.Is(err, bitunix_errors.ErrSignatureError):
			log.Fatalf("Signature error: %s", err.Error())
		case errors.Is(err, bitunix_errors.ErrNetwork):
			log.Fatalf("Network error: %s", err.Error())
		case errors.Is(err, bitunix_errors.ErrIPNotAllowed):
			log.Fatalf("IP restriction error: %s", err.Error())
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

	balanceHandler := &BalanceHandler{}
	if err := ws.SubscribeBalance(balanceHandler); err != nil {
		handleError(err, "balance")
	}

	positionHandler := &PositionHandler{}
	if err := ws.SubscribePositions(positionHandler); err != nil {
		handleError(err, "positions")
	}

	orderHandler := &OrderHandler{}
	if err := ws.SubscribeOrders(orderHandler); err != nil {
		handleError(err, "orders")
	}

	tpslHandler := &TpSlOrderHandler{}
	if err := ws.SubscribeTpSlOrders(tpslHandler); err != nil {
		handleError(err, "tpsl orders")
	}

	if err := ws.Stream(); err != nil {
		switch {
		case errors.Is(err, bitunix_errors.ErrTimeout):
			log.WithError(err).Fatalf("timeout while streaming")
		case errors.Is(err, bitunix_errors.ErrConnectionClosed):
			log.WithError(err).Info("connection closed")
		case errors.Is(err, bitunix_errors.ErrWorkgroupExhausted):
			log.WithError(err).Info("workgroup is exhausted")
		default:
			var apiErr *bitunix_errors.WebsocketError
			if errors.As(err, &apiErr) {
				log.Fatalf("Websocket error (code: %s): %s", apiErr.Operation, apiErr.Message)
			} else {
				log.Fatalf("Unexpected error: %s", err.Error())
			}
		}
	}
}

func handleError(err error, name string) {
	switch {
	case errors.Is(err, bitunix_errors.ErrValidation):
		log.WithError(err).Fatalf("failed to subscribe to %s", name)
	default:
		var apiErr *bitunix_errors.WebsocketError
		if errors.As(err, &apiErr) {
			log.Fatalf("Websocket error (code: %s): %s", apiErr.Operation, apiErr.Message)
		} else {
			log.Fatalf("Unexpected error: %s", err.Error())
		}
	}
}
