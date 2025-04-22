package bitunix

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

func (client *Client) GetPositionHistory(ctx context.Context, params PositionHistoryParams) (*PositionHistoryResponse, error) {
	queryParams := url.Values{}

	if params.Symbol != "" {
		queryParams.Add("symbol", params.Symbol)
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

	responseBody, err := client.api.Get(ctx, "/api/v1/futures/position/get_history_positions", queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get position history: %w", err)
	}

	response := &PositionHistoryResponse{}
	if err := json.Unmarshal(responseBody, response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}

type PositionHistoryParams struct {
	Symbol     string
	PositionID string
	StartTime  *time.Time
	EndTime    *time.Time
	Skip       int64
	Limit      int64
}

type PositionHistoryResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Data    struct {
		Positions []HistoricalPosition `json:"positionList"`
		Total     string               `json:"total"`
	} `json:"data"`
}

type HistoricalPosition struct {
	PositionID   string    `json:"positionId"`
	Symbol       string    `json:"symbol"`
	MaxQty       float64   `json:"-"`
	EntryPrice   float64   `json:"-"`
	ClosePrice   float64   `json:"-"`
	LiqQty       float64   `json:"-"`
	Side         string    `json:"side"`
	PositionMode string    `json:"positionMode"`
	MarginMode   string    `json:"marginMode"`
	Leverage     string    `json:"leverage"`
	Fee          float64   `json:"-"`
	Funding      float64   `json:"-"`
	RealizedPNL  float64   `json:"-"`
	LiqPrice     float64   `json:"-"`
	Ctime        time.Time `json:"-"`
	Mtime        time.Time `json:"-"`
}

func (p *HistoricalPosition) UnmarshalJSON(data []byte) error {
	type Alias HistoricalPosition
	aux := &struct {
		Fee         string `json:"fee"`
		RealizedPNL string `json:"realizedPNL"`
		LiqPrice    string `json:"liqPrice"`
		Funding     string `json:"funding"`
		Ctime       string `json:"ctime"`
		Mtime       string `json:"mtime"`
		MaxQty      string `json:"maxQty"`
		EntryPrice  string `json:"entryPrice"`
		ClosePrice  string `json:"closePrice"`
		LiqQty      string `json:"liqQty"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	ctime, err := strconv.ParseInt(aux.Ctime, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse ctime: %w", err)
	}

	p.Ctime = time.Unix(0, ctime*1000000)

	mtime, err := strconv.ParseInt(aux.Mtime, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse mtime: %w", err)
	}
	p.Mtime = time.Unix(0, mtime*1000000)

	feeFloat, err := strconv.ParseFloat(aux.Fee, 64)
	if err == nil {
		p.Fee = feeFloat
	} else {
		return fmt.Errorf("failed to parse fee: %w", err)
	}

	funding, err := strconv.ParseFloat(aux.Funding, 64)
	if err == nil {
		p.Funding = funding
	} else {
		return fmt.Errorf("failed to parse funding: %w", err)
	}

	realizedPNL, err := strconv.ParseFloat(aux.RealizedPNL, 64)
	if err == nil {
		p.RealizedPNL = realizedPNL
	} else {
		return fmt.Errorf("failed to parse realizedPNL: %w", err)
	}

	liqPrice, err := strconv.ParseFloat(aux.LiqPrice, 64)
	if err == nil {
		p.LiqPrice = liqPrice
	} else {
		return fmt.Errorf("failed to parse liqPrice: %w", err)
	}

	maxQty, err := strconv.ParseFloat(aux.MaxQty, 64)
	if err == nil {
		p.MaxQty = maxQty
	} else {
		return fmt.Errorf("failed to parse maxQty: %w", err)
	}

	entryPrice, err := strconv.ParseFloat(aux.EntryPrice, 64)
	if err == nil {
		p.EntryPrice = entryPrice
	} else {
		return fmt.Errorf("failed to parse entryPrice: %w", err)
	}

	closePrice, err := strconv.ParseFloat(aux.ClosePrice, 64)
	if err == nil {
		p.ClosePrice = closePrice
	} else {
		return fmt.Errorf("failed to parse closePrice: %w", err)
	}

	liqQty, err := strconv.ParseFloat(aux.LiqQty, 64)
	if err == nil {
		p.LiqQty = liqQty
	} else {
		return fmt.Errorf("failed to parse liqQty: %w", err)
	}

	return nil
}
