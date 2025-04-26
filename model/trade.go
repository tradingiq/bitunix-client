package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type TradeHistoryParams struct {
	Symbol     string
	OrderID    string
	PositionID string
	StartTime  *time.Time
	EndTime    *time.Time
	Skip       int64
	Limit      int64
}

type TradeHistoryResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    struct {
		Trades []HistoricalTrade `json:"tradeList"`
		Total  string            `json:"total"`
	} `json:"data"`
}

type HistoricalTrade struct {
	TradeID      string            `json:"tradeId"`
	OrderID      string            `json:"orderId"`
	Symbol       string            `json:"symbol"`
	Quantity     float64           `json:"-"`
	PositionMode TradePositionMode `json:"positionMode"`
	MarginMode   MarginMode        `json:"marginMode"`
	Leverage     int               `json:"leverage"`
	Price        float64           `json:"-"`
	Side         TradeSide         `json:"side"`
	OrderType    OrderType         `json:"orderType"`
	Effect       string            `json:"effect"`
	ClientID     string            `json:"clientId"`
	ReduceOnly   bool              `json:"reduceOnly"`
	Fee          float64           `json:"-"`
	RealizedPNL  float64           `json:"-"`
	CreateTime   time.Time         `json:"-"`
	RoleType     string            `json:"roleType"`
}

func (t *HistoricalTrade) UnmarshalJSON(data []byte) error {
	type Alias HistoricalTrade
	aux := &struct {
		Quantity    string `json:"qty"`
		Price       string `json:"price"`
		Fee         string `json:"fee"`
		RealizedPNL string `json:"realizedPNL"`
		CreateTime  string `json:"ctime"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	quantity, err := strconv.ParseFloat(aux.Quantity, 64)
	if err == nil {
		t.Quantity = quantity
	} else {
		return fmt.Errorf("invalid quantity: %w", err)
	}

	price, err := strconv.ParseFloat(aux.Price, 64)
	if err == nil {
		t.Price = price
	} else {
		return fmt.Errorf("invalid price: %w", err)
	}

	fee, err := strconv.ParseFloat(aux.Fee, 64)
	if err == nil {
		t.Fee = fee
	} else {
		return fmt.Errorf("invalid fee: %w", err)
	}

	realizedPNL, err := strconv.ParseFloat(aux.RealizedPNL, 64)
	if err == nil {
		t.RealizedPNL = realizedPNL
	} else {
		return fmt.Errorf("invalid realized pnl: %w", err)
	}

	createTime, err := strconv.ParseInt(aux.CreateTime, 10, 64)
	if err == nil {
		t.CreateTime = time.Unix(0, createTime*1000000)
	} else {
		return fmt.Errorf("invalid create time: %w", err)
	}

	return nil
}
