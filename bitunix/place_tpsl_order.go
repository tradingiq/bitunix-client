package bitunix

import (
	"context"
	"encoding/json"

	"github.com/tradingiq/bitunix-client/errors"
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
		return model.TPSLOrderRequest{}, errors.NewValidationError("symbol", "is required", nil)
	}

	if b.request.PositionID == "" {
		return model.TPSLOrderRequest{}, errors.NewValidationError("positionId", "is required", nil)
	}

	tpSet := b.request.TpPrice != nil
	slSet := b.request.SlPrice != nil

	if !tpSet && !slSet {
		return model.TPSLOrderRequest{}, errors.NewValidationError("tpPrice/slPrice", "at least one of tpPrice or slPrice must be set", nil)
	}

	tpQtySet := b.request.TpQty != nil && *b.request.TpQty > 0
	slQtySet := b.request.SlQty != nil && *b.request.SlQty > 0

	if !tpQtySet && !slQtySet {
		return model.TPSLOrderRequest{}, errors.NewValidationError("tpQty/slQty", "at least one of tpQty or slQty must be set", nil)
	}

	if tpSet {
		if *b.request.TpPrice <= 0 {
			return model.TPSLOrderRequest{}, errors.NewValidationError("tpPrice", "must be greater than zero", nil)
		}

		if !tpQtySet {
			return model.TPSLOrderRequest{}, errors.NewValidationError("tpQty", "is required when setting tpPrice", nil)
		}

		if b.request.TpOrderType == model.OrderTypeLimit {
			if b.request.TpOrderPrice == nil || *b.request.TpOrderPrice <= 0 {
				return model.TPSLOrderRequest{}, errors.NewValidationError("tpOrderPrice", "is required when tpOrderType is LIMIT", nil)
			}
		}
	}

	if slSet {
		if *b.request.SlPrice <= 0 {
			return model.TPSLOrderRequest{}, errors.NewValidationError("slPrice", "must be greater than zero", nil)
		}

		if !slQtySet {
			return model.TPSLOrderRequest{}, errors.NewValidationError("slQty", "is required when setting slPrice", nil)
		}

		if b.request.SlOrderType == model.OrderTypeLimit {
			if b.request.SlOrderPrice == nil || *b.request.SlOrderPrice <= 0 {
				return model.TPSLOrderRequest{}, errors.NewValidationError("slOrderPrice", "is required when slOrderType is LIMIT", nil)
			}
		}
	}

	return b.request, nil
}

func (c *apiClient) PlaceTpSlOrder(ctx context.Context, request *model.TPSLOrderRequest) (*model.TpSlOrderResponse, error) {
	marshaledRequest, err := json.Marshal(request)
	if err != nil {
		return nil, errors.NewInternalError("failed to marshal tpsl order request", err)
	}

	endpoint := "/api/v1/futures/tpsl/place_order"
	responseBody, err := c.restClient.Post(ctx, endpoint, nil, marshaledRequest)
	if err != nil {
		return nil, err
	}

	response := &model.TpSlOrderResponse{}
	if err := handleAPIResponse(responseBody, endpoint, response); err != nil {
		return nil, err
	}

	return response, nil
}
