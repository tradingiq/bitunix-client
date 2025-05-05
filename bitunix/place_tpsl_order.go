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

func (b *TPSLOrderBuilder) Build() (model.TPSLOrderRequest, error) {
	if b.request.Symbol == "" {
		return model.TPSLOrderRequest{}, fmt.Errorf("symbol is required")
	}

	if b.request.PositionID == "" {
		return model.TPSLOrderRequest{}, fmt.Errorf("positionId is required")
	}

	tpSet := b.request.TpPrice != nil
	slSet := b.request.SlPrice != nil

	if !tpSet && !slSet {
		return model.TPSLOrderRequest{}, fmt.Errorf("at least one of tpPrice or slPrice must be set")
	}

	tpQtySet := b.request.TpQty != nil && *b.request.TpQty > 0
	slQtySet := b.request.SlQty != nil && *b.request.SlQty > 0

	if !tpQtySet && !slQtySet {
		return model.TPSLOrderRequest{}, fmt.Errorf("at least one of tpQty or slQty must be set")
	}

	if tpSet {
		if *b.request.TpPrice <= 0 {
			return model.TPSLOrderRequest{}, fmt.Errorf("tpPrice must be greater than zero")
		}

		if !tpQtySet {
			return model.TPSLOrderRequest{}, fmt.Errorf("tpQty is required when setting tpPrice")
		}

		if b.request.TpOrderType == model.OrderTypeLimit {
			if b.request.TpOrderPrice == nil || *b.request.TpOrderPrice <= 0 {
				return model.TPSLOrderRequest{}, fmt.Errorf("tpOrderPrice is required when tpOrderType is LIMIT")
			}
		}
	}

	if slSet {
		if *b.request.SlPrice <= 0 {
			return model.TPSLOrderRequest{}, fmt.Errorf("slPrice must be greater than zero")
		}

		if !slQtySet {
			return model.TPSLOrderRequest{}, fmt.Errorf("slQty is required when setting slPrice")
		}

		if b.request.SlOrderType == model.OrderTypeLimit {
			if b.request.SlOrderPrice == nil || *b.request.SlOrderPrice <= 0 {
				return model.TPSLOrderRequest{}, fmt.Errorf("slOrderPrice is required when slOrderType is LIMIT")
			}
		}
	}

	return b.request, nil
}

func (c *apiClient) PlaceTpSlOrder(ctx context.Context, request *model.TPSLOrderRequest) (*model.TpSlOrderResponse, error) {
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

	return response, nil
}
