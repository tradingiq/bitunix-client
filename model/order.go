package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
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

type OrderRequest struct {
	Symbol       string      `json:"symbol"`
	TradeAction  TradeSide   `json:"side"`
	Price        *float64    `json:"-"`
	Qty          *float64    `json:"-"`
	PositionID   string      `json:"positionId,omitempty"`
	TradeSide    Side        `json:"tradeSide"`
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

type CancelOrderResponse struct {
	Code    int                     `json:"code"`
	Message string                  `json:"msg"`
	Data    CancelOrderResponseData `json:"data"`
}

type CancelOrderResponseData struct {
	SuccessList []CancelOrderResult  `json:"successList"`
	FailureList []CancelOrderFailure `json:"failureList"`
}

type CancelOrderResult struct {
	OrderId  string `json:"orderId"`
	ClientId string `json:"clientId"`
}

type CancelOrderFailure struct {
	OrderId   string `json:"orderId"`
	ClientId  string `json:"clientId"`
	ErrorMsg  string `json:"errorMsg"`
	ErrorCode string `json:"errorCode"`
}

type CancelOrderParam struct {
	OrderID  string `json:"orderId,omitempty"`
	ClientID string `json:"clientId,omitempty"`
}

type CancelOrderRequest struct {
	Symbol    string             `json:"symbol"`
	OrderList []CancelOrderParam `json:"orderList"`
}

type OrderHistoryParams struct {
	Symbol    string
	OrderID   string
	ClientID  string
	Status    string
	Type      string
	StartTime *time.Time
	EndTime   *time.Time
	Skip      int64
	Limit     int64
}

type OrderHistoryResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    struct {
		Orders []HistoricalOrder `json:"orderList"`
		Total  string            `json:"total"`
	} `json:"data"`
}

type HistoricalOrder struct {
	OrderID       string       `json:"orderId"`
	Symbol        string       `json:"symbol"`
	Quantity      float64      `json:"-"`
	TradeQuantity float64      `json:"-"`
	PositionMode  PositionMode `json:"-"`
	MarginMode    MarginMode   `json:"-"`
	Leverage      int          `json:"leverage"`
	Price         string       `json:"price"`
	Side          TradeSide    `json:"-"`
	OrderType     OrderType    `json:"-"`
	Effect        TimeInForce  `json:"-"`
	ClientID      string       `json:"clientId"`
	ReduceOnly    bool         `json:"reduceOnly"`
	Status        OrderStatus  `json:"-"`
	Fee           float64      `json:"-"`
	RealizedPNL   float64      `json:"-"`
	TpPrice       float64      `json:"-"`
	TpOrderPrice  float64      `json:"-"`
	SlPrice       float64      `json:"-"`
	TpStopType    StopType     `json:"-"`
	TpOrderType   OrderType    `json:"-"`
	SlStopType    StopType     `json:"-"`
	SlOrderType   OrderType    `json:"-"`
	SlOrderPrice  float64      `json:"-"`
	CreateTime    time.Time    `json:"-"`
	ModifyTime    time.Time    `json:"-"`
}

func (o *HistoricalOrder) UnmarshalJSON(data []byte) error {
	type Alias HistoricalOrder
	aux := &struct {
		Quantity      string `json:"qty"`
		TradeQuantity string `json:"tradeQty"`
		Fee           string `json:"fee"`
		RealizedPNL   string `json:"realizedPNL"`
		TpPrice       string `json:"tpPrice"`
		TpOrderPrice  string `json:"tpOrderPrice"`
		SlPrice       string `json:"slPrice"`
		SlOrderPrice  string `json:"slOrderPrice"`
		CreateTime    string `json:"ctime"`
		ModifyTime    string `json:"mtime"`
		Status        string `json:"status"`
		PositionMode  string `json:"positionMode"`
		MarginMode    string `json:"marginMode"`
		Side          string `json:"side"`
		OrderType     string `json:"orderType"`
		Effect        string `json:"effect"`
		TpStopType    string `json:"tpStopType"`
		TpOrderType   string `json:"tpOrderType"`
		SlStopType    string `json:"slStopType"`
		SlOrderType   string `json:"slOrderType"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Quantity != "" {
		quantity, err := strconv.ParseFloat(aux.Quantity, 64)
		if err == nil {
			o.Quantity = quantity
		} else {
			return fmt.Errorf("invalid quantity: %w", err)
		}
	}

	if aux.TradeQuantity != "" {
		tradeQty, err := strconv.ParseFloat(aux.TradeQuantity, 64)
		if err == nil {
			o.TradeQuantity = tradeQty
		} else {
			return fmt.Errorf("invalid trade quantity: %w", err)
		}
	}

	if aux.Fee != "" {
		fee, err := strconv.ParseFloat(aux.Fee, 64)
		if err == nil {
			o.Fee = fee
		} else {
			return fmt.Errorf("invalid fee: %w", err)
		}
	}

	if aux.RealizedPNL != "" {
		realizedPNL, err := strconv.ParseFloat(aux.RealizedPNL, 64)
		if err == nil {
			o.RealizedPNL = realizedPNL
		} else {
			return fmt.Errorf("invalid realized pnl: %w", err)
		}
	}

	if aux.TpPrice != "" {
		tpPrice, err := strconv.ParseFloat(aux.TpPrice, 64)
		if err == nil {
			o.TpPrice = tpPrice
		} else {
			return fmt.Errorf("invalid tp price: %w", err)
		}
	}

	if aux.TpOrderPrice != "" {
		tpOrderPrice, err := strconv.ParseFloat(aux.TpOrderPrice, 64)
		if err == nil {
			o.TpOrderPrice = tpOrderPrice
		} else {
			return fmt.Errorf("invalid tp order price: %w", err)
		}
	}

	if aux.SlPrice != "" {
		slPrice, err := strconv.ParseFloat(aux.SlPrice, 64)
		if err == nil {
			o.SlPrice = slPrice
		} else {
			return fmt.Errorf("invalid sl price: %w", err)
		}
	}

	if aux.SlOrderPrice != "" {
		slOrderPrice, err := strconv.ParseFloat(aux.SlOrderPrice, 64)
		if err == nil {
			o.SlOrderPrice = slOrderPrice
		} else {
			return fmt.Errorf("invalid sl order price: %w", err)
		}
	}

	if aux.CreateTime != "" {
		createTime, err := strconv.ParseInt(aux.CreateTime, 10, 64)
		if err == nil {
			o.CreateTime = time.Unix(0, createTime*1000000)
		} else {
			return fmt.Errorf("invalid create time: %w", err)
		}
	}

	if aux.ModifyTime != "" {
		modifyTime, err := strconv.ParseInt(aux.ModifyTime, 10, 64)
		if err == nil {
			o.ModifyTime = time.Unix(0, modifyTime*1000000)
		} else {
			return fmt.Errorf("invalid modify time: %w", err)
		}
	}

	status, err := ParseOrderStatus(aux.Status)
	if err != nil {
		return fmt.Errorf("invalid order status: %w", err)
	}
	o.Status = status

	marginMode, err := ParseMarginMode(aux.MarginMode)
	if err != nil {
		return fmt.Errorf("invalid margin mode: %w", err)
	}
	o.MarginMode = marginMode

	positionMode, err := ParsePositionMode(aux.PositionMode)
	if err != nil {
		return fmt.Errorf("invalid position mode: %w", err)
	}
	o.PositionMode = positionMode

	side, err := ParseTradeSide(aux.Side)
	if err != nil {
		return fmt.Errorf("invalid side: %w", err)
	}
	o.Side = side

	orderType, err := ParseOrderType(aux.OrderType)
	if err != nil {
		return fmt.Errorf("invalid order type: %w", err)
	}
	o.OrderType = orderType

	effect, err := ParseTimeInForce(aux.Effect)
	if err != nil {
		return fmt.Errorf("invalid effect: %w", err)
	}
	o.Effect = effect

	slStopType, err := ParseStopType(aux.SlStopType)
	if err != nil {
		return fmt.Errorf("invalid sl stop type: %w", err)
	}
	o.SlStopType = slStopType

	tpStopType, err := ParseStopType(aux.TpStopType)
	if err != nil {
		return fmt.Errorf("invalid tp stop type: %w", err)
	}
	o.TpStopType = tpStopType

	tpOrderType, err := ParseOrderType(aux.TpOrderType)
	if err != nil {
		return fmt.Errorf("invalid tp order type: %w", err)
	}
	o.TpOrderType = tpOrderType

	slOrderType, err := ParseOrderType(aux.SlOrderType)
	if err != nil {
		return fmt.Errorf("invalid sl order type: %w", err)
	}
	o.SlOrderType = slOrderType

	return nil
}
