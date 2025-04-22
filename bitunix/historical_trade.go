package bitunix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
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

func (client *Client) GetTradeHistory(ctx context.Context, params *TradeHistoryParams) (*TradeHistoryResponse, error) {
	queryParams := url.Values{}

	if params != nil {
		if params.Symbol != "" {
			queryParams.Add("symbol", params.Symbol)
		}
		if params.OrderID != "" {
			queryParams.Add("orderId", params.OrderID)
		}
		if params.PositionID != "" {
			queryParams.Add("positionId", params.PositionID)
		}
		if params.StartTime != nil {
			queryParams.Add("startTime", strconv.FormatInt(params.StartTime.UnixMilli(), 10))
		}
		if params.EndTime != nil {
			queryParams.Add("endTime", strconv.FormatInt(params.EndTime.UnixMilli(), 10))
		}
		if params.Skip > 0 {
			queryParams.Add("skip", strconv.FormatInt(params.Skip, 10))
		}
		if params.Limit > 0 {
			queryParams.Add("limit", strconv.FormatInt(params.Limit, 10))
		}
	}

	responseBody, err := client.api.Get(ctx, "/api/v1/futures/trade/get_history_trades", queryParams)
	if err != nil {
		return nil, err
	}

	response := &TradeHistoryResponse{}
	if err := json.Unmarshal(responseBody, response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, err
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
		return err
	}

	price, err := strconv.ParseFloat(aux.Price, 64)
	if err == nil {
		t.Price = price
	} else {
		return err
	}

	fee, err := strconv.ParseFloat(aux.Fee, 64)
	if err == nil {
		t.Fee = fee
	} else {
		return err
	}

	realizedPNL, err := strconv.ParseFloat(aux.RealizedPNL, 64)
	if err == nil {
		t.RealizedPNL = realizedPNL
	} else {
		return err
	}

	createTime, err := strconv.ParseInt(aux.CreateTime, 10, 64)
	if err == nil {
		t.CreateTime = time.Unix(0, createTime*1000000)
	} else {
		return err
	}

	return nil
}
