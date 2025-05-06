package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type TradeHistoryParams struct {
	Symbol     Symbol
	OrderID    string
	PositionID string
	StartTime  *time.Time
	EndTime    *time.Time
	Skip       int64
	Limit      int64
}

type TradeHistoryResponse struct {
	BaseResponse
	Data struct {
		Trades []HistoricalTrade `json:"tradeList"`
		Total  string            `json:"total"`
	} `json:"data"`
}

type HistoricalTrade struct {
	TradeID      string        `json:"tradeId"`
	OrderID      string        `json:"orderId"`
	Symbol       Symbol        `json:"symbol"`
	Quantity     float64       `json:"-"`
	PositionMode PositionMode  `json:"-"`
	MarginMode   MarginMode    `json:"-"`
	Leverage     int           `json:"leverage"`
	Price        float64       `json:"-"`
	Side         TradeSide     `json:"-"`
	OrderType    OrderType     `json:"-"`
	Effect       string        `json:"effect"`
	ClientID     string        `json:"clientId"`
	ReduceOnly   bool          `json:"reduceOnly"`
	Fee          float64       `json:"-"`
	RealizedPNL  float64       `json:"-"`
	CreateTime   time.Time     `json:"-"`
	RoleType     TradeRoleType `json:"-"`
}

func (t *HistoricalTrade) UnmarshalJSON(data []byte) error {
	type Alias HistoricalTrade
	aux := &struct {
		Quantity     string `json:"qty"`
		Price        string `json:"price"`
		Fee          string `json:"fee"`
		RealizedPNL  string `json:"realizedPNL"`
		CreateTime   string `json:"ctime"`
		PositionMode string `json:"positionMode"`
		Side         string `json:"side"`
		OrderType    string `json:"orderType"`
		MarginMode   string `json:"marginMode"`
		RoleType     string `json:"roleType"`
		Symbol       string `json:"symbol"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	t.Symbol = ParseSymbol(aux.Symbol)

	if aux.Quantity != "" {
		quantity, err := strconv.ParseFloat(aux.Quantity, 64)
		if err == nil {
			t.Quantity = quantity
		} else {
			return fmt.Errorf("invalid quantity: %w", err)
		}
	}

	if aux.Price != "" {
		price, err := strconv.ParseFloat(aux.Price, 64)
		if err == nil {
			t.Price = price
		} else {
			return fmt.Errorf("invalid price: %w", err)
		}
	}

	if aux.Fee != "" {
		fee, err := strconv.ParseFloat(aux.Fee, 64)
		if err == nil {
			t.Fee = fee
		} else {
			return fmt.Errorf("invalid fee: %w", err)
		}
	}

	if aux.RealizedPNL != "" {
		realizedPNL, err := strconv.ParseFloat(aux.RealizedPNL, 64)
		if err == nil {
			t.RealizedPNL = realizedPNL
		} else {
			return fmt.Errorf("invalid realized pnl: %w", err)
		}
	}

	if aux.CreateTime != "" {
		createTime, err := strconv.ParseInt(aux.CreateTime, 10, 64)
		if err == nil {
			t.CreateTime = time.Unix(0, createTime*1000000)
		} else {
			return fmt.Errorf("invalid create time: %w", err)
		}
	}

	side, err := ParseTradeSide(aux.Side)
	if err != nil {
		return fmt.Errorf("invalid side: %w", err)
	}
	t.Side = side

	orderType, err := ParseOrderType(aux.OrderType)
	if err != nil {
		return fmt.Errorf("invalid order type: %w", err)
	}
	t.OrderType = orderType

	posMode, err := ParsePositionMode(aux.PositionMode)
	if err != nil {
		return fmt.Errorf("invalid position mode: %w", err)
	}
	t.PositionMode = posMode

	marginMode, err := ParseMarginMode(aux.MarginMode)
	if err != nil {
		return fmt.Errorf("invalid margin mode: %w", err)
	}
	t.MarginMode = marginMode

	roleType, err := ParseTradeRoleType(aux.RoleType)
	if err != nil {
		return fmt.Errorf("invalid role type: %w", err)
	}
	t.RoleType = roleType

	return nil
}
