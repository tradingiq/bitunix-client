package bitunix

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

type TPSLOrderRequest struct {
	Symbol       string    `json:"symbol"`
	PositionID   string    `json:"positionId"`
	TpPrice      *float64  `json:"-"`
	SlPrice      *float64  `json:"-"`
	TpStopType   StopType  `json:"tpStopType,omitempty"`
	SlStopType   StopType  `json:"slStopType,omitempty"`
	TpOrderType  OrderType `json:"tpOrderType,omitempty"`
	SlOrderType  OrderType `json:"slOrderType,omitempty"`
	TpOrderPrice *float64  `json:"-"`
	SlOrderPrice *float64  `json:"-"`
	TpQty        *float64  `json:"-"`
	SlQty        *float64  `json:"-"`
}

func (r TPSLOrderRequest) MarshalJSON() ([]byte, error) {
	type Alias TPSLOrderRequest

	aux := &struct {
		TpPrice      string `json:"tpPrice,omitempty"`
		SlPrice      string `json:"slPrice,omitempty"`
		TpOrderPrice string `json:"tpOrderPrice,omitempty"`
		SlOrderPrice string `json:"slOrderPrice,omitempty"`
		TpQty        string `json:"tpQty,omitempty"`
		SlQty        string `json:"slQty,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(&r),
	}

	if r.TpPrice != nil {
		aux.TpPrice = strconv.FormatFloat(*r.TpPrice, 'f', -1, 64)
	}

	if r.SlPrice != nil {
		aux.SlPrice = strconv.FormatFloat(*r.SlPrice, 'f', -1, 64)
	}

	if r.TpOrderPrice != nil {
		aux.TpOrderPrice = strconv.FormatFloat(*r.TpOrderPrice, 'f', -1, 64)
	}

	if r.SlOrderPrice != nil {
		aux.SlOrderPrice = strconv.FormatFloat(*r.SlOrderPrice, 'f', -1, 64)
	}

	if r.TpQty != nil {
		aux.TpQty = strconv.FormatFloat(*r.TpQty, 'f', -1, 64)
	}

	if r.SlQty != nil {
		aux.SlQty = strconv.FormatFloat(*r.SlQty, 'f', -1, 64)
	}

	return json.Marshal(aux)
}

type TPSLOrderResponse struct {
	Code    int                     `json:"code"`
	Message string                  `json:"msg"`
	Data    []TPSLOrderResponseData `json:"data"`
}

type TPSLOrderResponseData struct {
	OrderID  string `json:"orderId"`
	ClientId string `json:"clientId"`
}

type TPSLOrderBuilder struct {
	request TPSLOrderRequest
}

func NewTPSLOrderBuilder(symbol, positionID string) *TPSLOrderBuilder {
	return &TPSLOrderBuilder{
		request: TPSLOrderRequest{
			Symbol:     symbol,
			PositionID: positionID,
		},
	}
}

func (b *TPSLOrderBuilder) WithTakeProfit(price float64, qty float64, stopType StopType, orderType OrderType, orderPrice float64) *TPSLOrderBuilder {
	b.request.TpPrice = &price
	b.request.TpQty = &qty
	b.request.TpStopType = stopType
	b.request.TpOrderType = orderType
	b.request.TpOrderPrice = &orderPrice
	return b
}

func (b *TPSLOrderBuilder) WithStopLoss(price float64, qty float64, stopType StopType, orderType OrderType, orderPrice float64) *TPSLOrderBuilder {
	b.request.SlPrice = &price
	b.request.SlQty = &qty
	b.request.SlStopType = stopType
	b.request.SlOrderType = orderType
	b.request.SlOrderPrice = &orderPrice
	return b
}

func (b *TPSLOrderBuilder) Build() TPSLOrderRequest {
	return b.request
}

func (c *client) PlaceTpSlOrder(ctx context.Context, request *TPSLOrderRequest) (*TPSLOrderResponse, error) {
	marshaledRequest, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order request: %w", err)
	}

	responseBody, err := c.restClient.Post(ctx, "/api/v1/futures/tpsl/place_order", nil, marshaledRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to place order request: %w", err)
	}

	response := &TPSLOrderResponse{}
	if err := json.Unmarshal(responseBody, response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, err
}
