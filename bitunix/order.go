package bitunix

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tradingiq/bitunix-client/model"
	"time"
)

func (c *apiClient) PlaceOrder(ctx context.Context, request *model.OrderRequest) (*model.OrderResponse, error) {
	marshaledRequest, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order request: %w", err)
	}

	responseBody, err := c.restClient.Post(ctx, "/api/v1/futures/trade/place_order", nil, marshaledRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to place order request: %w", err)
	}

	response := &model.OrderResponse{}
	if err := json.Unmarshal(responseBody, response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}

func (c *apiClient) CancelOrders(ctx context.Context, request *model.CancelOrderRequest) (*model.CancelOrderResponse, error) {
	marshaledRequest, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal cancel order request: %w", err)
	}

	responseBody, err := c.restClient.Post(ctx, "/api/v1/futures/trade/cancel_orders", nil, marshaledRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to cancel orders request: %w", err)
	}

	response := &model.CancelOrderResponse{}
	if err := json.Unmarshal(responseBody, response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}

type OrderBuilder struct {
	request model.OrderRequest
	errors  []string
}

func NewOrderBuilder(symbol model.Symbol, side model.TradeSide, tradeSide model.Side, qty float64) *OrderBuilder {
	qtyVal := qty

	return &OrderBuilder{
		request: model.OrderRequest{
			Symbol:     symbol,
			OrderType:  model.OrderTypeMarket,
			TradeSide:  side,
			Side:       tradeSide,
			Qty:        &qtyVal,
			ReduceOnly: false,
			ClientID:   fmt.Sprintf("client_%d", time.Now().UnixNano()),
			Effect:     model.TimeInForceGTC,
		},
	}
}

func (b *OrderBuilder) WithOrderType(orderType model.OrderType) *OrderBuilder {
	b.request.OrderType = orderType
	return b
}

func (b *OrderBuilder) WithPrice(price float64) *OrderBuilder {
	b.request.Price = &price
	return b
}

func (b *OrderBuilder) WithPositionID(positionID string) *OrderBuilder {
	b.request.PositionID = positionID
	return b
}

func (b *OrderBuilder) WithReduceOnly(reduceOnly bool) *OrderBuilder {
	b.request.ReduceOnly = reduceOnly
	return b
}

func (b *OrderBuilder) WithTimeInForce(tif model.TimeInForce) *OrderBuilder {
	b.request.Effect = tif
	return b
}

func (b *OrderBuilder) WithClientID(clientID string) *OrderBuilder {
	b.request.ClientID = clientID
	return b
}

func (b *OrderBuilder) WithTakeProfit(price float64, stopType model.StopType, orderType model.OrderType) *OrderBuilder {
	b.request.TpPrice = &price
	b.request.TpStopType = stopType
	b.request.TpOrderType = orderType
	return b
}

func (b *OrderBuilder) WithTakeProfitPrice(orderPrice float64) *OrderBuilder {
	b.request.TpOrderPrice = &orderPrice
	return b
}

func (b *OrderBuilder) WithStopLoss(price float64, stopType model.StopType, orderType model.OrderType) *OrderBuilder {
	b.request.SlPrice = &price
	b.request.SlStopType = stopType
	b.request.SlOrderType = orderType
	return b
}

func (b *OrderBuilder) WithStopLossPrice(orderPrice float64) *OrderBuilder {
	b.request.SlOrderPrice = &orderPrice
	return b
}

func (b *OrderBuilder) Build() (model.OrderRequest, error) {
	if b.request.Symbol == "" {
		return model.OrderRequest{}, fmt.Errorf("symbol is required")
	}

	if b.request.TradeSide == "" {
		return model.OrderRequest{}, fmt.Errorf("side (trade action) is required")
	}

	if b.request.Side == "" {
		return model.OrderRequest{}, fmt.Errorf("tradeSide is required")
	}

	if b.request.Qty == nil || *b.request.Qty <= 0 {
		return model.OrderRequest{}, fmt.Errorf("qty (quantity) is required and must be greater than zero")
	}

	if b.request.Side == model.SideClose && b.request.PositionID == "" {
		return model.OrderRequest{}, fmt.Errorf("positionId is required when tradeSide is CLOSE")
	}

	if b.request.OrderType == "" {
		return model.OrderRequest{}, fmt.Errorf("orderType is required")
	}

	if b.request.OrderType == model.OrderTypeLimit && (b.request.Price == nil || *b.request.Price <= 0) {
		return model.OrderRequest{}, fmt.Errorf("price is required for limit orders")
	}

	if b.request.TpPrice != nil {
		if *b.request.TpPrice <= 0 {
			return model.OrderRequest{}, fmt.Errorf("take profit price must be greater than zero")
		}

		if b.request.TpStopType == "" {
			return model.OrderRequest{}, fmt.Errorf("tpStopType is required when setting take profit")
		}

		if b.request.TpOrderType == "" {
			return model.OrderRequest{}, fmt.Errorf("tpOrderType is required when setting take profit")
		}

		if b.request.TpOrderType == model.OrderTypeLimit && (b.request.TpOrderPrice == nil || *b.request.TpOrderPrice <= 0) {
			return model.OrderRequest{}, fmt.Errorf("tpOrderPrice is required when tpOrderType is LIMIT")
		}
	}

	if b.request.SlPrice != nil {
		if *b.request.SlPrice <= 0 {
			return model.OrderRequest{}, fmt.Errorf("stop loss price must be greater than zero")
		}

		if b.request.SlStopType == "" {
			return model.OrderRequest{}, fmt.Errorf("slStopType is required when setting stop loss")
		}

		if b.request.SlOrderType == "" {
			return model.OrderRequest{}, fmt.Errorf("slOrderType is required when setting stop loss")
		}

		if b.request.SlOrderType == model.OrderTypeLimit && (b.request.SlOrderPrice == nil || *b.request.SlOrderPrice <= 0) {
			return model.OrderRequest{}, fmt.Errorf("slOrderPrice is required when slOrderType is LIMIT")
		}
	}

	return b.request, nil
}

type CancelOrderBuilder struct {
	request model.CancelOrderRequest
}

func NewCancelOrderBuilder(symbol model.Symbol) *CancelOrderBuilder {
	return &CancelOrderBuilder{
		request: model.CancelOrderRequest{
			Symbol:    symbol,
			OrderList: make([]model.CancelOrderParam, 0),
		},
	}
}

func (b *CancelOrderBuilder) WithOrderID(orderID string) *CancelOrderBuilder {
	b.request.OrderList = append(b.request.OrderList, model.CancelOrderParam{
		OrderID: orderID,
	})
	return b
}

func (b *CancelOrderBuilder) WithClientID(clientID string) *CancelOrderBuilder {
	b.request.OrderList = append(b.request.OrderList, model.CancelOrderParam{
		ClientID: clientID,
	})
	return b
}

func (b *CancelOrderBuilder) Build() (model.CancelOrderRequest, error) {
	if b.request.Symbol == "" {
		return model.CancelOrderRequest{}, fmt.Errorf("symbol is required")
	}

	if len(b.request.OrderList) == 0 {
		return model.CancelOrderRequest{}, fmt.Errorf("orderList must contain at least one order to cancel")
	}

	for i, order := range b.request.OrderList {
		if order.OrderID == "" && order.ClientID == "" {
			return model.CancelOrderRequest{}, fmt.Errorf("order at index %d must have either orderId or clientId", i)
		}
	}

	return b.request, nil
}
