package bitunix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

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

func (c *client) GetOrderHistory(ctx context.Context, params OrderHistoryParams) (*OrderHistoryResponse, error) {
	queryParams := url.Values{}

	if params.Symbol != "" {
		queryParams.Add("symbol", params.Symbol)
	}
	if params.OrderID != "" {
		queryParams.Add("orderId", params.OrderID)
	}
	if params.ClientID != "" {
		queryParams.Add("clientId", params.ClientID)
	}
	if params.Status != "" {
		queryParams.Add("status", params.Status)
	}
	if params.Type != "" {
		queryParams.Add("type", params.Type)
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

	responseBody, err := c.restClient.Get(ctx, "/api/v1/futures/trade/get_history_orders", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get order history: %w", err)
	}

	response := &OrderHistoryResponse{}
	if err := json.Unmarshal(responseBody, response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
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
	OrderID       string            `json:"orderId"`
	Symbol        string            `json:"symbol"`
	Quantity      float64           `json:"-"`
	TradeQuantity float64           `json:"-"`
	PositionMode  TradePositionMode `json:"positionMode"`
	MarginMode    MarginMode        `json:"marginMode"`
	Leverage      int               `json:"leverage"`
	Price         string            `json:"price"`
	Side          TradeAction       `json:"side"`
	OrderType     OrderType         `json:"orderType"`
	Effect        TimeInForce       `json:"effect"`
	ClientID      string            `json:"clientId"`
	ReduceOnly    bool              `json:"reduceOnly"`
	Status        string            `json:"status"`
	Fee           float64           `json:"-"`
	RealizedPNL   float64           `json:"-"`
	TpPrice       float64           `json:"-"`
	TpStopType    StopType          `json:"tpStopType"`
	TpOrderType   OrderType         `json:"tpOrderType"`
	TpOrderPrice  float64           `json:"-"`
	SlPrice       float64           `json:"-"`
	SlStopType    StopType          `json:"slStopType"`
	SlOrderType   OrderType         `json:"slOrderType"`
	SlOrderPrice  float64           `json:"-"`
	CreateTime    time.Time         `json:"-"`
	ModifyTime    time.Time         `json:"-"`
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

	return nil
}
