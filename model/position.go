package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

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
	PositionID   string       `json:"positionId"`
	Symbol       string       `json:"symbol"`
	MaxQty       float64      `json:"-"`
	EntryPrice   float64      `json:"-"`
	ClosePrice   float64      `json:"-"`
	LiqQty       float64      `json:"-"`
	Side         TradeSide    `json:"side"`
	PositionMode PositionMode `json:"positionMode"`
	MarginMode   MarginMode   `json:"marginMode"`
	Leverage     string       `json:"leverage"`
	Fee          float64      `json:"-"`
	Funding      float64      `json:"-"`
	RealizedPNL  float64      `json:"-"`
	LiqPrice     float64      `json:"-"`
	Ctime        time.Time    `json:"-"`
	Mtime        time.Time    `json:"-"`
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

	if aux.LiqQty != "" {
		liqQty, err := strconv.ParseFloat(aux.LiqQty, 64)
		if err == nil {
			p.LiqQty = liqQty
		} else {
			return fmt.Errorf("failed to parse liqQty: %w", err)
		}
	}

	return nil
}
