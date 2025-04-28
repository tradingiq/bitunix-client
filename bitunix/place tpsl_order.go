package bitunix

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tradingiq/bitunix-client/model"
)

type TPSLOrderBuilder struct {
	request model.TPSLOrderRequest
}

func NewTPSLOrderBuilder(symbol model.Symbol, positionID string) *TPSLOrderBuilder {
	return &TPSLOrderBuilder{
		request: model.TPSLOrderRequest{
			Symbol:     symbol,
			PositionID: positionID,
		},
	}
}

func (b *TPSLOrderBuilder) WithTakeProfit(price float64, qty float64, stopType model.StopType, orderType model.OrderType, orderPrice float64) *TPSLOrderBuilder {
	b.request.TpPrice = &price
	b.request.TpQty = &qty
	b.request.TpStopType = stopType
	b.request.TpOrderType = orderType
	b.request.TpOrderPrice = &orderPrice
	return b
}

func (b *TPSLOrderBuilder) WithStopLoss(price float64, qty float64, stopType model.StopType, orderType model.OrderType, orderPrice float64) *TPSLOrderBuilder {
	b.request.SlPrice = &price
	b.request.SlQty = &qty
	b.request.SlStopType = stopType
	b.request.SlOrderType = orderType
	b.request.SlOrderPrice = &orderPrice
	return b
}

func (b *TPSLOrderBuilder) Build() model.TPSLOrderRequest {
	return b.request
}

func (c *client) PlaceTpSlOrder(ctx context.Context, request *model.TPSLOrderRequest) (*model.TpSlOrderResponse, error) {
	marshaledRequest, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order request: %w", err)
	}

	responseBody, err := c.restClient.Post(ctx, "/api/v1/futures/tpsl/place_order", nil, marshaledRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to place order request: %w", err)
	}

	response := &model.TpSlOrderResponse{}
	if err := json.Unmarshal(responseBody, response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, err
}
