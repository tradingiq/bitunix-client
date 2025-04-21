package bitunix

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type OrderRequest struct {
	Symbol       string      `json:"symbol"`
	TradeAction  TradeAction `json:"side"`
	Price        *float64    `json:"-"`
	Qty          *float64    `json:"-"`
	PositionID   string      `json:"positionId,omitempty"`
	TradeSide    TradeSide   `json:"tradeSide"`
	OrderType    OrderType   `json:"orderType"`
	ReduceOnly   bool        `json:"reduceOnly"`
	Effect       TimeInForce `json:"effect,omitempty"`
	ClientID     string      `json:"clientId,omitempty"`
	TpPrice      *float64    `json:"-"`
	TpStopType   StopType    `json:"tpStopType,omitempty"`
	TpOrderType  OrderType   `json:"tpOrderType,omitempty"`
	TpOrderPrice *float64    `json:"-"`
	SlPrice      *float64    `json:"-"`
	SlStopType   StopType    `json:"slStopType,omitempty"`
	SlOrderType  OrderType   `json:"slOrderType,omitempty"`
	SlOrderPrice *float64    `json:"-"`
}

func (r *OrderRequest) MarshalJSON() ([]byte, error) {
	type Alias OrderRequest

	aux := &struct {
		Price        string `json:"price,omitempty"`
		Qty          string `json:"qty"`
		TpPrice      string `json:"tpPrice,omitempty"`
		TpOrderPrice string `json:"tpOrderPrice,omitempty"`
		SlPrice      string `json:"slPrice,omitempty"`
		SlOrderPrice string `json:"slOrderPrice,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if r.Price != nil {
		aux.Price = strconv.FormatFloat(*r.Price, 'f', -1, 64)
	}

	if r.Qty != nil {
		aux.Qty = strconv.FormatFloat(*r.Qty, 'f', -1, 64)
	}

	if r.TpPrice != nil {
		aux.TpPrice = strconv.FormatFloat(*r.TpPrice, 'f', -1, 64)
	}

	if r.TpOrderPrice != nil {
		aux.TpOrderPrice = strconv.FormatFloat(*r.TpOrderPrice, 'f', -1, 64)
	}

	if r.SlPrice != nil {
		aux.SlPrice = strconv.FormatFloat(*r.SlPrice, 'f', -1, 64)
	}

	if r.SlOrderPrice != nil {
		aux.SlOrderPrice = strconv.FormatFloat(*r.SlOrderPrice, 'f', -1, 64)
	}

	return json.Marshal(aux)
}

type OrderResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    OrderResponseData `json:"data"`
}

type OrderResponseData struct {
	OrderId  string `json:"orderId"`
	ClientId string `json:"clientId"`
}

func (client *Client) PlaceOrder(ctx context.Context, request *OrderRequest) (*OrderResponse, error) {
	marshaledRequest, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal order request: %w", err)
	}

	responseBody, err := client.api.Post(ctx, "/api/v1/futures/trade/place_order", nil, marshaledRequest)
	response := &OrderResponse{}
	if err := json.Unmarshal(responseBody, response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, err
}

type OrderBuilder struct {
	request OrderRequest
	errors  []string
}

func NewOrderBuilder(symbol string, side TradeAction, tradeSide TradeSide, qty float64) *OrderBuilder {
	qtyVal := qty

	return &OrderBuilder{
		request: OrderRequest{
			Symbol:      symbol,
			OrderType:   OrderTypeMarket,
			TradeAction: side,
			TradeSide:   tradeSide,
			Qty:         &qtyVal,
			ReduceOnly:  false,
			ClientID:    fmt.Sprintf("client_%d", time.Now().UnixNano()),
			Effect:      TimeInForceGTC,
		},
	}
}

func (b *OrderBuilder) WithOrderType(orderType OrderType) *OrderBuilder {
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

func (b *OrderBuilder) WithTimeInForce(tif TimeInForce) *OrderBuilder {
	b.request.Effect = tif
	return b
}

func (b *OrderBuilder) WithClientID(clientID string) *OrderBuilder {
	b.request.ClientID = clientID
	return b
}

func (b *OrderBuilder) WithTakeProfit(price float64, stopType StopType, orderType OrderType) *OrderBuilder {
	b.request.TpPrice = &price
	b.request.TpStopType = stopType
	b.request.TpOrderType = orderType
	return b
}

func (b *OrderBuilder) WithTakeProfitPrice(orderPrice float64) *OrderBuilder {
	b.request.TpOrderPrice = &orderPrice
	return b
}

func (b *OrderBuilder) WithStopLoss(price float64, stopType StopType, orderType OrderType) *OrderBuilder {
	b.request.SlPrice = &price
	b.request.SlStopType = stopType
	b.request.SlOrderType = orderType
	return b
}

func (b *OrderBuilder) WithStopLossPrice(orderPrice float64) *OrderBuilder {
	b.request.SlOrderPrice = &orderPrice
	return b
}

func (b *OrderBuilder) Build() OrderRequest {
	return b.request
}
