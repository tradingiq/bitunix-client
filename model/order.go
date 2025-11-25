package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type TPSLOrderRequest struct {
	Symbol       Symbol    `json:"symbol"`
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

func (r *TPSLOrderRequest) MarshalJSON() ([]byte, error) {
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
		Alias: (*Alias)(r),
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

type TpSlOrderResponse struct {
	BaseResponse
	Data []TPSLOrderResponseData `json:"data"`
}

type TPSLOrderResponseData struct {
	OrderID  string `json:"orderId"`
	ClientId string `json:"clientId"`
}

type OrderRequest struct {
	Symbol       Symbol      `json:"symbol"`
	TradeSide    TradeSide   `json:"side"`
	Price        *float64    `json:"-"`
	Qty          float64     `json:"-"`
	PositionID   string      `json:"positionId,omitempty"`
	Side         Side        `json:"tradeSide"`
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

	aux.Qty = strconv.FormatFloat(r.Qty, 'f', -1, 64)

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
	BaseResponse
	Data CancelOrderResponseData `json:"data"`
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
	Symbol    Symbol             `json:"symbol"`
	OrderList []CancelOrderParam `json:"orderList"`
}

type OrderHistoryParams struct {
	Symbol    Symbol
	OrderID   string
	ClientID  string
	Status    OrderStatus
	Type      OrderType
	StartTime *time.Time
	EndTime   *time.Time
	Skip      int64
	Limit     int64
}

type OrderHistoryResponse struct {
	BaseResponse
	Data struct {
		Orders []HistoricalOrder `json:"orderList"`
		Total  string            `json:"total"`
	} `json:"data"`
}

type OrderDetailResponse struct {
	BaseResponse
	Data *OrderDetail `json:"data"`
}

type OrderDetail struct {
	OrderID       string       `json:"orderId"`
	Symbol        Symbol       `json:"symbol"`
	Quantity      float64      `json:"-"`
	TradeQuantity float64      `json:"-"`
	PositionMode  PositionMode `json:"-"`
	MarginMode    MarginMode   `json:"-"`
	Leverage      int          `json:"leverage"`
	Price         float64      `json:"-"`
	Side          TradeSide    `json:"-"`
	OrderType     OrderType    `json:"-"`
	Effect        TimeInForce  `json:"-"`
	ClientID      string       `json:"clientId"`
	ReduceOnly    bool         `json:"reduceOnly"`
	Status        OrderStatus  `json:"-"`
	Fee           float64      `json:"-"`
	RealizedPNL   float64      `json:"-"`
	TpPrice       *float64     `json:"-"`
	TpOrderPrice  *float64     `json:"-"`
	SlPrice       *float64     `json:"-"`
	TpStopType    *StopType    `json:"-"`
	TpOrderType   *OrderType   `json:"-"`
	SlStopType    *StopType    `json:"-"`
	SlOrderType   *OrderType   `json:"-"`
	SlOrderPrice  *float64     `json:"-"`
	CreateTime    time.Time    `json:"-"`
	ModifyTime    time.Time    `json:"-"`
}

func (o *OrderDetail) UnmarshalJSON(data []byte) error {
	type Alias OrderDetail
	aux := &struct {
		Quantity      string  `json:"qty"`
		TradeQuantity string  `json:"tradeQty"`
		Fee           string  `json:"fee"`
		RealizedPNL   string  `json:"realizedPNL"`
		TpPrice       *string `json:"tpPrice"`
		Price         string  `json:"price"`
		TpOrderPrice  *string `json:"tpOrderPrice"`
		SlPrice       *string `json:"slPrice"`
		SlOrderPrice  *string `json:"slOrderPrice"`
		CreateTime    string  `json:"ctime"`
		ModifyTime    string  `json:"mtime"`
		Status        string  `json:"status"`
		PositionMode  string  `json:"positionMode"`
		MarginMode    string  `json:"marginMode"`
		Side          string  `json:"side"`
		OrderType     string  `json:"orderType"`
		Effect        string  `json:"effect"`
		TpStopType    *string `json:"tpStopType"`
		TpOrderType   *string `json:"tpOrderType"`
		SlStopType    *string `json:"slStopType"`
		SlOrderType   *string `json:"slOrderType"`
		Symbol        string  `json:"symbol"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	o.Symbol = ParseSymbol(aux.Symbol)

	price, err := strconv.ParseFloat(aux.Price, 64)
	if err == nil {
		o.Price = price
	} else {
		return fmt.Errorf("invalid price: %w", err)
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

	if aux.TpPrice != nil {
		tpPrice, err := strconv.ParseFloat(*aux.TpPrice, 64)
		if err == nil {
			o.TpPrice = &tpPrice
		} else {
			return fmt.Errorf("invalid tp price: %w", err)
		}
	}

	if aux.TpOrderPrice != nil {
		tpOrderPrice, err := strconv.ParseFloat(*aux.TpOrderPrice, 64)
		if err == nil {
			o.TpOrderPrice = &tpOrderPrice
		} else {
			return fmt.Errorf("invalid tp order price: %w", err)
		}
	}

	if aux.SlPrice != nil {
		slPrice, err := strconv.ParseFloat(*aux.SlPrice, 64)
		if err == nil {
			o.SlPrice = &slPrice
		} else {
			return fmt.Errorf("invalid sl price: %w", err)
		}
	}

	if aux.SlOrderPrice != nil {
		slOrderPrice, err := strconv.ParseFloat(*aux.SlOrderPrice, 64)
		if err == nil {
			o.SlOrderPrice = &slOrderPrice
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

	if aux.SlStopType != nil {
		slStopType, err := ParseStopType(*aux.SlStopType)
		if err != nil {
			return fmt.Errorf("invalid sl stop type: %w", err)
		}
		o.SlStopType = &slStopType
	}

	if aux.TpStopType != nil {
		tpStopType, err := ParseStopType(*aux.TpStopType)
		if err != nil {
			return fmt.Errorf("invalid tp stop type: %w", err)
		}
		o.TpStopType = &tpStopType
	}

	if aux.TpOrderType != nil {
		tpOrderType, err := ParseOrderType(*aux.TpOrderType)
		if err != nil {
			return fmt.Errorf("invalid tp order type: %w", err)
		}
		o.TpOrderType = &tpOrderType
	}

	if aux.SlOrderType != nil {
		slOrderType, err := ParseOrderType(*aux.SlOrderType)
		if err != nil {
			return fmt.Errorf("invalid sl order type: %w", err)
		}
		o.SlOrderType = &slOrderType
	}

	return nil
}

type HistoricalOrder struct {
	OrderID       string       `json:"orderId"`
	Symbol        Symbol       `json:"symbol"`
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
	TpPrice       *float64     `json:"-"`
	TpOrderPrice  *float64     `json:"-"`
	SlPrice       *float64     `json:"-"`
	TpStopType    *StopType    `json:"-"`
	TpOrderType   *OrderType   `json:"-"`
	SlStopType    *StopType    `json:"-"`
	SlOrderType   *OrderType   `json:"-"`
	SlOrderPrice  *float64     `json:"-"`
	CreateTime    time.Time    `json:"-"`
	ModifyTime    time.Time    `json:"-"`
}

func (o *HistoricalOrder) UnmarshalJSON(data []byte) error {
	type Alias HistoricalOrder
	aux := &struct {
		Quantity      string  `json:"qty"`
		TradeQuantity string  `json:"tradeQty"`
		Fee           string  `json:"fee"`
		RealizedPNL   string  `json:"realizedPNL"`
		TpPrice       *string `json:"tpPrice"`
		TpOrderPrice  *string `json:"tpOrderPrice"`
		SlPrice       *string `json:"slPrice"`
		SlOrderPrice  *string `json:"slOrderPrice"`
		CreateTime    string  `json:"ctime"`
		ModifyTime    string  `json:"mtime"`
		Status        string  `json:"status"`
		PositionMode  string  `json:"positionMode"`
		MarginMode    string  `json:"marginMode"`
		Side          string  `json:"side"`
		OrderType     string  `json:"orderType"`
		Effect        string  `json:"effect"`
		TpStopType    *string `json:"tpStopType"`
		TpOrderType   *string `json:"tpOrderType"`
		SlStopType    *string `json:"slStopType"`
		SlOrderType   *string `json:"slOrderType"`
		Symbol        string  `json:"symbol"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	o.Symbol = ParseSymbol(aux.Symbol)

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

	if aux.TpPrice != nil {
		tpPrice, err := strconv.ParseFloat(*aux.TpPrice, 64)
		if err == nil {
			o.TpPrice = &tpPrice
		} else {
			return fmt.Errorf("invalid tp price: %w", err)
		}
	}

	if aux.TpOrderPrice != nil {
		tpOrderPrice, err := strconv.ParseFloat(*aux.TpOrderPrice, 64)
		if err == nil {
			o.TpOrderPrice = &tpOrderPrice
		} else {
			return fmt.Errorf("invalid tp order price: %w", err)
		}
	}

	if aux.SlPrice != nil {
		slPrice, err := strconv.ParseFloat(*aux.SlPrice, 64)
		if err == nil {
			o.SlPrice = &slPrice
		} else {
			return fmt.Errorf("invalid sl price: %w", err)
		}
	}

	if aux.SlOrderPrice != nil {
		slOrderPrice, err := strconv.ParseFloat(*aux.SlOrderPrice, 64)
		if err == nil {
			o.SlOrderPrice = &slOrderPrice
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

	if aux.SlStopType != nil {
		slStopType, err := ParseStopType(*aux.SlStopType)
		if err != nil {
			return fmt.Errorf("invalid sl stop type: %w", err)
		}
		o.SlStopType = &slStopType
	}

	if aux.TpStopType != nil {
		tpStopType, err := ParseStopType(*aux.TpStopType)
		if err != nil {
			return fmt.Errorf("invalid tp stop type: %w", err)
		}
		o.TpStopType = &tpStopType
	}

	if aux.TpOrderType != nil {
		tpOrderType, err := ParseOrderType(*aux.TpOrderType)
		if err != nil {
			return fmt.Errorf("invalid tp order type: %w", err)
		}
		o.TpOrderType = &tpOrderType
	}

	if aux.SlOrderType != nil {
		slOrderType, err := ParseOrderType(*aux.SlOrderType)
		if err != nil {
			return fmt.Errorf("invalid sl order type: %w", err)
		}
		o.SlOrderType = &slOrderType
	}

	return nil
}

type PendingTPSLOrderParams struct {
	Symbol       Symbol
	PositionID   string
	Side         TradeSide
	PositionMode PositionMode
	Skip         int64
	Limit        int64
}

type PendingTPSLOrderResponse struct {
	BaseResponse
	Data []PendingTPSLOrder `json:"data"`
}

type PendingTPSLOrder struct {
	ID           string     `json:"id"`
	PositionID   string     `json:"positionId"`
	Symbol       Symbol     `json:"-"`
	Base         string     `json:"base"`
	Quote        string     `json:"quote"`
	TpPrice      *float64   `json:"-"`
	TpStopType   *StopType  `json:"-"`
	SlPrice      *float64   `json:"-"`
	SlStopType   *StopType  `json:"-"`
	TpOrderType  *OrderType `json:"-"`
	TpOrderPrice *float64   `json:"-"`
	SlOrderType  *OrderType `json:"-"`
	SlOrderPrice *float64   `json:"-"`
	TpQty        *float64   `json:"-"`
	SlQty        *float64   `json:"-"`
}

func (o *PendingTPSLOrder) UnmarshalJSON(data []byte) error {
	type Alias PendingTPSLOrder
	aux := &struct {
		Symbol       string  `json:"symbol"`
		TpPrice      *string `json:"tpPrice"`
		TpStopType   *string `json:"tpStopType"`
		SlPrice      *string `json:"slPrice"`
		SlStopType   *string `json:"slStopType"`
		TpOrderType  *string `json:"tpOrderType"`
		TpOrderPrice *string `json:"tpOrderPrice"`
		SlOrderType  *string `json:"slOrderType"`
		SlOrderPrice *string `json:"slOrderPrice"`
		TpQty        *string `json:"tpQty"`
		SlQty        *string `json:"slQty"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	o.Symbol = ParseSymbol(aux.Symbol)

	if aux.TpPrice != nil && *aux.TpPrice != "" {
		tpPrice, err := strconv.ParseFloat(*aux.TpPrice, 64)
		if err != nil {
			return fmt.Errorf("invalid tp price: %w", err)
		}
		o.TpPrice = &tpPrice
	}

	if aux.SlPrice != nil && *aux.SlPrice != "" {
		slPrice, err := strconv.ParseFloat(*aux.SlPrice, 64)
		if err != nil {
			return fmt.Errorf("invalid sl price: %w", err)
		}
		o.SlPrice = &slPrice
	}

	if aux.TpOrderPrice != nil && *aux.TpOrderPrice != "" {
		tpOrderPrice, err := strconv.ParseFloat(*aux.TpOrderPrice, 64)
		if err != nil {
			return fmt.Errorf("invalid tp order price: %w", err)
		}
		o.TpOrderPrice = &tpOrderPrice
	}

	if aux.SlOrderPrice != nil && *aux.SlOrderPrice != "" {
		slOrderPrice, err := strconv.ParseFloat(*aux.SlOrderPrice, 64)
		if err != nil {
			return fmt.Errorf("invalid sl order price: %w", err)
		}
		o.SlOrderPrice = &slOrderPrice
	}

	if aux.TpQty != nil && *aux.TpQty != "" {
		tpQty, err := strconv.ParseFloat(*aux.TpQty, 64)
		if err != nil {
			return fmt.Errorf("invalid tp qty: %w", err)
		}
		o.TpQty = &tpQty
	}

	if aux.SlQty != nil && *aux.SlQty != "" {
		slQty, err := strconv.ParseFloat(*aux.SlQty, 64)
		if err != nil {
			return fmt.Errorf("invalid sl qty: %w", err)
		}
		o.SlQty = &slQty
	}

	if aux.TpStopType != nil && *aux.TpStopType != "" {
		tpStopType, err := ParseStopType(*aux.TpStopType)
		if err != nil {
			return fmt.Errorf("invalid tp stop type: %w", err)
		}
		o.TpStopType = &tpStopType
	}

	if aux.SlStopType != nil && *aux.SlStopType != "" {
		slStopType, err := ParseStopType(*aux.SlStopType)
		if err != nil {
			return fmt.Errorf("invalid sl stop type: %w", err)
		}
		o.SlStopType = &slStopType
	}

	if aux.TpOrderType != nil && *aux.TpOrderType != "" {
		tpOrderType, err := ParseOrderType(*aux.TpOrderType)
		if err != nil {
			return fmt.Errorf("invalid tp order type: %w", err)
		}
		o.TpOrderType = &tpOrderType
	}

	if aux.SlOrderType != nil && *aux.SlOrderType != "" {
		slOrderType, err := ParseOrderType(*aux.SlOrderType)
		if err != nil {
			return fmt.Errorf("invalid sl order type: %w", err)
		}
		o.SlOrderType = &slOrderType
	}

	return nil
}

type PendingOrderParams struct {
	Symbol    Symbol
	OrderID   string
	ClientID  string
	Status    OrderStatus
	StartTime *time.Time
	EndTime   *time.Time
	Skip      int64
	Limit     int64
}

type PendingOrderResponse struct {
	BaseResponse
	Data struct {
		OrderList []PendingOrder `json:"orderList"`
		Total     int64          `json:"-"`
	} `json:"data"`
}

func (r *PendingOrderResponse) UnmarshalJSON(data []byte) error {
	type Alias PendingOrderResponse
	aux := &struct {
		Data struct {
			OrderList []PendingOrder `json:"orderList"`
			Total     string         `json:"total"`
		} `json:"data"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	r.Data.OrderList = aux.Data.OrderList

	if aux.Data.Total != "" {
		total, err := strconv.ParseInt(aux.Data.Total, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid total: %w", err)
		}
		r.Data.Total = total
	}

	return nil
}

type PendingOrder struct {
	OrderID       string       `json:"orderId"`
	Symbol        Symbol       `json:"-"`
	Quantity      float64      `json:"-"`
	TradeQuantity float64      `json:"-"`
	PositionMode  PositionMode `json:"-"`
	MarginMode    MarginMode   `json:"-"`
	Leverage      int          `json:"leverage"`
	Price         float64      `json:"-"`
	Side          TradeSide    `json:"-"`
	OrderType     OrderType    `json:"-"`
	Effect        TimeInForce  `json:"-"`
	ClientID      string       `json:"clientId"`
	ReduceOnly    bool         `json:"reduceOnly"`
	Status        OrderStatus  `json:"-"`
	Fee           float64      `json:"-"`
	RealizedPNL   float64      `json:"-"`
	TpPrice       *float64     `json:"-"`
	TpStopType    *StopType    `json:"-"`
	TpOrderType   *OrderType   `json:"-"`
	TpOrderPrice  *float64     `json:"-"`
	SlPrice       *float64     `json:"-"`
	SlStopType    *StopType    `json:"-"`
	SlOrderType   *OrderType   `json:"-"`
	SlOrderPrice  *float64     `json:"-"`
	CreateTime    time.Time    `json:"-"`
	ModifyTime    time.Time    `json:"-"`
}

func (o *PendingOrder) UnmarshalJSON(data []byte) error {
	type Alias PendingOrder
	aux := &struct {
		Quantity      string  `json:"qty"`
		TradeQuantity string  `json:"tradeQty"`
		Fee           string  `json:"fee"`
		RealizedPNL   string  `json:"realizedPNL"`
		TpPrice       *string `json:"tpPrice"`
		Price         string  `json:"price"`
		TpOrderPrice  *string `json:"tpOrderPrice"`
		SlPrice       *string `json:"slPrice"`
		SlOrderPrice  *string `json:"slOrderPrice"`
		CreateTime    string  `json:"ctime"`
		ModifyTime    string  `json:"mtime"`
		Status        string  `json:"status"`
		PositionMode  string  `json:"positionMode"`
		MarginMode    string  `json:"marginMode"`
		Side          string  `json:"side"`
		OrderType     string  `json:"orderType"`
		Effect        string  `json:"effect"`
		TpStopType    *string `json:"tpStopType"`
		TpOrderType   *string `json:"tpOrderType"`
		SlStopType    *string `json:"slStopType"`
		SlOrderType   *string `json:"slOrderType"`
		Symbol        string  `json:"symbol"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	o.Symbol = ParseSymbol(aux.Symbol)

	price, err := strconv.ParseFloat(aux.Price, 64)
	if err == nil {
		o.Price = price
	} else {
		return fmt.Errorf("invalid price: %w", err)
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

	if aux.TpPrice != nil && *aux.TpPrice != "" {
		tpPrice, err := strconv.ParseFloat(*aux.TpPrice, 64)
		if err == nil {
			o.TpPrice = &tpPrice
		} else {
			return fmt.Errorf("invalid tp price: %w", err)
		}
	}

	if aux.TpOrderPrice != nil && *aux.TpOrderPrice != "" {
		tpOrderPrice, err := strconv.ParseFloat(*aux.TpOrderPrice, 64)
		if err == nil {
			o.TpOrderPrice = &tpOrderPrice
		} else {
			return fmt.Errorf("invalid tp order price: %w", err)
		}
	}

	if aux.SlPrice != nil && *aux.SlPrice != "" {
		slPrice, err := strconv.ParseFloat(*aux.SlPrice, 64)
		if err == nil {
			o.SlPrice = &slPrice
		} else {
			return fmt.Errorf("invalid sl price: %w", err)
		}
	}

	if aux.SlOrderPrice != nil && *aux.SlOrderPrice != "" {
		slOrderPrice, err := strconv.ParseFloat(*aux.SlOrderPrice, 64)
		if err == nil {
			o.SlOrderPrice = &slOrderPrice
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

	if aux.SlStopType != nil && *aux.SlStopType != "" {
		slStopType, err := ParseStopType(*aux.SlStopType)
		if err != nil {
			return fmt.Errorf("invalid sl stop type: %w", err)
		}
		o.SlStopType = &slStopType
	}

	if aux.TpStopType != nil && *aux.TpStopType != "" {
		tpStopType, err := ParseStopType(*aux.TpStopType)
		if err != nil {
			return fmt.Errorf("invalid tp stop type: %w", err)
		}
		o.TpStopType = &tpStopType
	}

	if aux.TpOrderType != nil && *aux.TpOrderType != "" {
		tpOrderType, err := ParseOrderType(*aux.TpOrderType)
		if err != nil {
			return fmt.Errorf("invalid tp order type: %w", err)
		}
		o.TpOrderType = &tpOrderType
	}

	if aux.SlOrderType != nil && *aux.SlOrderType != "" {
		slOrderType, err := ParseOrderType(*aux.SlOrderType)
		if err != nil {
			return fmt.Errorf("invalid sl order type: %w", err)
		}
		o.SlOrderType = &slOrderType
	}

	return nil
}

type TPSLOrderHistoryParams struct {
	Symbol       Symbol
	Side         PositionSide
	PositionMode PositionMode
	StartTime    *time.Time
	EndTime      *time.Time
	Skip         int64
	Limit        int64
}

type TPSLOrderHistoryResponse struct {
	BaseResponse
	Data struct {
		OrderList []HistoricalTPSLOrder `json:"orderList"`
		Total     int64                 `json:"-"`
	} `json:"data"`
}

func (r *TPSLOrderHistoryResponse) UnmarshalJSON(data []byte) error {
	type Alias TPSLOrderHistoryResponse
	aux := &struct {
		Data struct {
			OrderList []HistoricalTPSLOrder `json:"orderList"`
			Total     string                `json:"total"`
		} `json:"data"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Copy the data
	r.Data.OrderList = aux.Data.OrderList

	// Parse Total from string to int64
	if aux.Data.Total != "" {
		total, err := strconv.ParseInt(aux.Data.Total, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid total: %w", err)
		}
		r.Data.Total = total
	}

	return nil
}

type HistoricalTPSLOrder struct {
	ID           string     `json:"id"`
	PositionID   string     `json:"positionId"`
	Symbol       Symbol     `json:"-"`
	Base         string     `json:"base"`
	Quote        string     `json:"quote"`
	TpPrice      *float64   `json:"-"`
	TpStopType   *StopType  `json:"-"`
	SlPrice      *float64   `json:"-"`
	SlStopType   *StopType  `json:"-"`
	TpOrderType  *OrderType `json:"-"`
	TpOrderPrice *float64   `json:"-"`
	SlOrderType  *OrderType `json:"-"`
	SlOrderPrice *float64   `json:"-"`
	TpQty        *float64   `json:"-"`
	SlQty        *float64   `json:"-"`
	Status       string     `json:"status"`
	Ctime        time.Time  `json:"-"`
	TriggerTime  *time.Time `json:"-"`
}

func (o *HistoricalTPSLOrder) UnmarshalJSON(data []byte) error {
	aux := struct {
		ID           string  `json:"id"`
		PositionID   string  `json:"positionId"`
		Symbol       string  `json:"symbol"`
		Base         string  `json:"base"`
		Quote        string  `json:"quote"`
		TpPrice      *string `json:"tpPrice"`
		TpStopType   *string `json:"tpStopType"`
		SlPrice      *string `json:"slPrice"`
		SlStopType   *string `json:"slStopType"`
		TpOrderType  *string `json:"tpOrderType"`
		TpOrderPrice *string `json:"tpOrderPrice"`
		SlOrderType  *string `json:"slOrderType"`
		SlOrderPrice *string `json:"slOrderPrice"`
		TpQty        *string `json:"tpQty"`
		SlQty        *string `json:"slQty"`
		Status       string  `json:"status"`
		Ctime        string  `json:"ctime"`
		TriggerTime  *string `json:"triggerTime"`
	}{}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	o.ID = aux.ID
	o.PositionID = aux.PositionID
	o.Base = aux.Base
	o.Quote = aux.Quote
	o.Status = aux.Status

	// Parse Ctime from string to time.Time
	if aux.Ctime != "" {
		ctimeInt, err := strconv.ParseInt(aux.Ctime, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid ctime: %w", err)
		}
		o.Ctime = time.Unix(0, ctimeInt*1000000) // Convert milliseconds to nanoseconds
	}

	// Parse TriggerTime from string to time.Time (can be null)
	if aux.TriggerTime != nil && *aux.TriggerTime != "" {
		triggerTimeInt, err := strconv.ParseInt(*aux.TriggerTime, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid triggerTime: %w", err)
		}
		triggerTime := time.Unix(0, triggerTimeInt*1000000) // Convert milliseconds to nanoseconds
		o.TriggerTime = &triggerTime
	}

	if aux.Symbol != "" {
		o.Symbol = Symbol(aux.Symbol)
	} else if aux.Base != "" && aux.Quote != "" {
		o.Symbol = Symbol(aux.Base + aux.Quote)
	}

	if aux.TpPrice != nil {
		tpPrice, err := strconv.ParseFloat(*aux.TpPrice, 64)
		if err == nil {
			o.TpPrice = &tpPrice
		}
	}

	if aux.TpStopType != nil && *aux.TpStopType != "" {
		tpStopType, err := ParseStopType(*aux.TpStopType)
		if err == nil {
			o.TpStopType = &tpStopType
		}
	}

	if aux.SlPrice != nil {
		slPrice, err := strconv.ParseFloat(*aux.SlPrice, 64)
		if err == nil {
			o.SlPrice = &slPrice
		}
	}

	if aux.SlStopType != nil && *aux.SlStopType != "" {
		slStopType, err := ParseStopType(*aux.SlStopType)
		if err == nil {
			o.SlStopType = &slStopType
		}
	}

	if aux.TpOrderType != nil && *aux.TpOrderType != "" {
		tpOrderType, err := ParseOrderType(*aux.TpOrderType)
		if err == nil {
			o.TpOrderType = &tpOrderType
		}
	}

	if aux.TpOrderPrice != nil {
		tpOrderPrice, err := strconv.ParseFloat(*aux.TpOrderPrice, 64)
		if err == nil {
			o.TpOrderPrice = &tpOrderPrice
		}
	}

	if aux.SlOrderType != nil && *aux.SlOrderType != "" {
		slOrderType, err := ParseOrderType(*aux.SlOrderType)
		if err == nil {
			o.SlOrderType = &slOrderType
		}
	}

	if aux.SlOrderPrice != nil {
		slOrderPrice, err := strconv.ParseFloat(*aux.SlOrderPrice, 64)
		if err == nil {
			o.SlOrderPrice = &slOrderPrice
		}
	}

	if aux.TpQty != nil {
		tpQty, err := strconv.ParseFloat(*aux.TpQty, 64)
		if err == nil {
			o.TpQty = &tpQty
		}
	}

	if aux.SlQty != nil {
		slQty, err := strconv.ParseFloat(*aux.SlQty, 64)
		if err == nil {
			o.SlQty = &slQty
		}
	}

	return nil
}
