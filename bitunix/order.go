package bitunix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
)

func (c *apiClient) PlaceOrder(ctx context.Context, request *model.OrderRequest) (*model.OrderResponse, error) {
	marshaledRequest, err := json.Marshal(request)
	if err != nil {
		return nil, errors.NewInternalError("failed to marshal order request", err)
	}

	endpoint := "/api/v1/futures/trade/place_order"
	responseBody, err := c.restClient.Post(ctx, endpoint, nil, marshaledRequest)
	if err != nil {
		return nil, err
	}

	response := &model.OrderResponse{}
	if err := handleAPIResponse(responseBody, endpoint, response); err != nil {
		return nil, err
	}

	return response, nil
}

func (c *apiClient) CancelOrders(ctx context.Context, request *model.CancelOrderRequest) (*model.CancelOrderResponse, error) {
	marshaledRequest, err := json.Marshal(request)
	if err != nil {
		return nil, errors.NewInternalError("failed to marshal cancel order request", err)
	}

	endpoint := "/api/v1/futures/trade/cancel_orders"
	responseBody, err := c.restClient.Post(ctx, endpoint, nil, marshaledRequest)
	if err != nil {
		return nil, err
	}

	response := &model.CancelOrderResponse{}
	if err := handleAPIResponse(responseBody, endpoint, response); err != nil {
		return nil, err
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
			Qty:        qtyVal,
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
		return model.OrderRequest{}, errors.NewValidationError("symbol", "is required", nil)
	}

	if b.request.TradeSide == "" {
		return model.OrderRequest{}, errors.NewValidationError("side", "trade action is required", nil)
	}

	if b.request.Side == "" {
		return model.OrderRequest{}, errors.NewValidationError("tradeSide", "is required", nil)
	}

	if b.request.Qty <= 0 {
		return model.OrderRequest{}, errors.NewValidationError("qty", "quantity is required and must be greater than zero", nil)
	}

	if b.request.Side == model.SideClose && b.request.PositionID == "" {
		return model.OrderRequest{}, errors.NewValidationError("positionId", "is required when tradeSide is CLOSE", nil)
	}

	if b.request.OrderType == "" {
		return model.OrderRequest{}, errors.NewValidationError("orderType", "is required", nil)
	}

	if b.request.OrderType == model.OrderTypeLimit && (b.request.Price == nil || *b.request.Price <= 0) {
		return model.OrderRequest{}, errors.NewValidationError("price", "is required for limit orders", nil)
	}

	if b.request.TpPrice != nil {
		if *b.request.TpPrice <= 0 {
			return model.OrderRequest{}, errors.NewValidationError("tpPrice", "take profit price must be greater than zero", nil)
		}

		if b.request.TpStopType == "" {
			return model.OrderRequest{}, errors.NewValidationError("tpStopType", "is required when setting take profit", nil)
		}

		if b.request.TpOrderType == "" {
			return model.OrderRequest{}, errors.NewValidationError("tpOrderType", "is required when setting take profit", nil)
		}

		if b.request.TpOrderType == model.OrderTypeLimit && (b.request.TpOrderPrice == nil || *b.request.TpOrderPrice <= 0) {
			return model.OrderRequest{}, errors.NewValidationError("tpOrderPrice", "is required when tpOrderType is LIMIT", nil)
		}
	}

	if b.request.SlPrice != nil {
		if *b.request.SlPrice <= 0 {
			return model.OrderRequest{}, errors.NewValidationError("slPrice", "stop loss price must be greater than zero", nil)
		}

		if b.request.SlStopType == "" {
			return model.OrderRequest{}, errors.NewValidationError("slStopType", "is required when setting stop loss", nil)
		}

		if b.request.SlOrderType == "" {
			return model.OrderRequest{}, errors.NewValidationError("slOrderType", "is required when setting stop loss", nil)
		}

		if b.request.SlOrderType == model.OrderTypeLimit && (b.request.SlOrderPrice == nil || *b.request.SlOrderPrice <= 0) {
			return model.OrderRequest{}, errors.NewValidationError("slOrderPrice", "is required when slOrderType is LIMIT", nil)
		}
	}

	return b.request, nil
}

type OrderDetailRequest struct {
	OrderID  string
	ClientID string
}

func (c *apiClient) GetOrderDetail(ctx context.Context, request *OrderDetailRequest) (*model.OrderDetailResponse, error) {
	params := url.Values{}

	// At least one of orderId or clientId is required
	if request.OrderID == "" && request.ClientID == "" {
		return nil, errors.NewValidationError("request", "either orderId or clientId is required", nil)
	}

	if request.OrderID != "" {
		params.Set("orderId", request.OrderID)
	}

	if request.ClientID != "" {
		params.Set("clientId", request.ClientID)
	}

	endpoint := "/api/v1/futures/trade/get_order_detail"
	responseBody, err := c.restClient.Get(ctx, endpoint, params)
	if err != nil {
		return nil, err
	}

	response := &model.OrderDetailResponse{}
	if err := handleAPIResponse(responseBody, endpoint, response); err != nil {
		return nil, err
	}

	return response, nil
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
		return model.CancelOrderRequest{}, errors.NewValidationError("symbol", "is required", nil)
	}

	if len(b.request.OrderList) == 0 {
		return model.CancelOrderRequest{}, errors.NewValidationError("orderList", "must contain at least one order to cancel", nil)
	}

	for i, order := range b.request.OrderList {
		if order.OrderID == "" && order.ClientID == "" {
			return model.CancelOrderRequest{}, errors.NewValidationError(
				fmt.Sprintf("orderList[%d]", i),
				"must have either orderId or clientId",
				nil,
			)
		}
	}

	return b.request, nil
}
