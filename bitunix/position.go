package bitunix

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

func (client *Client) GetPositionHistory(ctx context.Context) (*PositionHistoryResponse, error) {
	responseBody, err := client.api.Get(ctx, "/api/v1/futures/position/get_history_positions", nil)
	if err != nil {
		return nil, err
	}

	response := &PositionHistoryResponse{}
	if err := json.Unmarshal(responseBody, response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, err
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
	LiqQty       float64   `json:"-"`            // Liquidated quantity
	Side         string    `json:"side"`         // LONG or SHORT
	PositionMode string    `json:"positionMode"` // ONE_WAY or HEDGE
	MarginMode   string    `json:"marginMode"`   // ISOLATION or CROSS
	Leverage     string    `json:"leverage"`
	Fee          float64   `json:"-"` // Deducted transaction fees
	Funding      float64   `json:"-"` // Total funding fee during the position
	RealizedPNL  float64   `json:"-"` // Realized PnL (excludes funding fee and transaction fee)
	LiqPrice     float64   `json:"-"` // Estimated liquidation price
	Ctime        time.Time `json:"-"` // Create timestamp
	Mtime        time.Time `json:"-"` // Latest modify timestamp
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
		return err
	}

	p.Ctime = time.Unix(0, ctime*1000000)

	mtime, err := strconv.ParseInt(aux.Mtime, 10, 64)
	if err != nil {
		return err
	}
	p.Mtime = time.Unix(0, mtime*1000000)

	feeFloat, err := strconv.ParseFloat(aux.Fee, 64)
	if err == nil {
		p.Fee = feeFloat
	}

	funding, err := strconv.ParseFloat(aux.Fee, 64)
	if err == nil {
		p.Funding = funding
	}

	realizedPNL, err := strconv.ParseFloat(aux.RealizedPNL, 64)
	if err == nil {
		p.RealizedPNL = realizedPNL
	}

	liqPrice, err := strconv.ParseFloat(aux.LiqPrice, 64)
	if err == nil {
		p.LiqPrice = liqPrice
	}

	maxQty, err := strconv.ParseFloat(aux.MaxQty, 64)
	if err == nil {
		p.MaxQty = maxQty
	}

	entryPrice, err := strconv.ParseFloat(aux.EntryPrice, 64)
	if err == nil {
		p.EntryPrice = entryPrice
	}

	closePrice, err := strconv.ParseFloat(aux.ClosePrice, 64)
	if err == nil {
		p.ClosePrice = closePrice
	}

	liqQty, err := strconv.ParseFloat(aux.LiqQty, 64)
	if err == nil {
		p.LiqQty = liqQty
	}

	return nil
}
