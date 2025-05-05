package main

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/bitunix"
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
		log.Fatalf("failed to connect to WebSocket: %v", err)
	}

	balanceHandler := &BalanceHandler{}
	if err := ws.SubscribeBalance(balanceHandler); err != nil {
		log.WithError(err).Fatal("failed to subscribe to balance")
	}

	positionHandler := &PositionHandler{}
	if err := ws.SubscribePositions(positionHandler); err != nil {
		log.WithError(err).Fatal("failed to subscribe to positions")
	}

	orderHandler := &OrderHandler{}
	if err := ws.SubscribeOrders(orderHandler); err != nil {
		log.WithError(err).Fatal("failed to subscribe to orders")
	}

	tpslHandler := &TpSlOrderHandler{}
	if err := ws.SubscribeTpSlOrders(tpslHandler); err != nil {
		log.WithError(err).Fatal("failed to subscribe to TP/SL orders")
	}

	if err := ws.Stream(); err != nil {
		log.WithError(err).Fatal("failed to stream")
	}
}
