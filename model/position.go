package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type PositionHistoryParams struct {
	Symbol     Symbol
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
	Symbol       Symbol       `json:"symbol"`
	MaxQty       float64      `json:"-"`
	EntryPrice   float64      `json:"-"`
	ClosePrice   float64      `json:"-"`
	LiqQty       float64      `json:"-"`
	Side         TradeSide    `json:"-"`
	PositionMode PositionMode `json:"-"`
	MarginMode   MarginMode   `json:"-"`
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
		Fee          string `json:"fee"`
		RealizedPNL  string `json:"realizedPNL"`
		LiqPrice     string `json:"liqPrice"`
		Funding      string `json:"funding"`
		Ctime        string `json:"ctime"`
		Mtime        string `json:"mtime"`
		MaxQty       string `json:"maxQty"`
		EntryPrice   string `json:"entryPrice"`
		ClosePrice   string `json:"closePrice"`
		LiqQty       string `json:"liqQty"`
		Side         string `json:"side"`
		PositionMode string `json:"positionMode"`
		MarginMode   string `json:"marginMode"`
		Symbol       string `json:"symbol"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Parse symbol
	p.Symbol = ParseSymbol(aux.Symbol)

	if aux.Ctime != "" {
		ctime, err := strconv.ParseInt(aux.Ctime, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse ctime: %w", err)
		}
		p.Ctime = time.Unix(0, ctime*1000000)
	}

	if aux.Mtime != "" {
		mtime, err := strconv.ParseInt(aux.Mtime, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse mtime: %w", err)
		}
		p.Mtime = time.Unix(0, mtime*1000000)
	}

	if aux.Fee != "" {
		feeFloat, err := strconv.ParseFloat(aux.Fee, 64)
		if err == nil {
			p.Fee = feeFloat
		} else {
			return fmt.Errorf("failed to parse fee: %w", err)
		}
	}

	if aux.Funding != "" {
		funding, err := strconv.ParseFloat(aux.Funding, 64)
		if err == nil {
			p.Funding = funding
		} else {
			return fmt.Errorf("failed to parse funding: %w", err)
		}
	}

	if aux.RealizedPNL != "" {
		realizedPNL, err := strconv.ParseFloat(aux.RealizedPNL, 64)
		if err == nil {
			p.RealizedPNL = realizedPNL
		} else {
			return fmt.Errorf("failed to parse realizedPNL: %w", err)
		}
	}

	if aux.LiqPrice != "" {
		liqPrice, err := strconv.ParseFloat(aux.LiqPrice, 64)
		if err == nil {
			p.LiqPrice = liqPrice
		} else {
			return fmt.Errorf("failed to parse liqPrice: %w", err)
		}
	}

	if aux.MaxQty != "" {
		maxQty, err := strconv.ParseFloat(aux.MaxQty, 64)
		if err == nil {
			p.MaxQty = maxQty
		} else {
			return fmt.Errorf("failed to parse maxQty: %w", err)
		}
	}

	if aux.EntryPrice != "" {
		entryPrice, err := strconv.ParseFloat(aux.EntryPrice, 64)
		if err == nil {
			p.EntryPrice = entryPrice
		} else {
			return fmt.Errorf("failed to parse entryPrice: %w", err)
		}
	}

	if aux.ClosePrice != "" {
		closePrice, err := strconv.ParseFloat(aux.ClosePrice, 64)
		if err == nil {
			p.ClosePrice = closePrice
		} else {
			return fmt.Errorf("failed to parse closePrice: %w", err)
		}
	}

	if aux.LiqQty != "" {
		liqQty, err := strconv.ParseFloat(aux.LiqQty, 64)
		if err == nil {
			p.LiqQty = liqQty
		} else {
			return fmt.Errorf("failed to parse liqQty: %w", err)
		}
	}

	side, err := ParseTradeSide(aux.Side)
	if err != nil {
		return fmt.Errorf("invalid side: %w", err)
	}
	p.Side = side

	posMode, err := ParsePositionMode(aux.PositionMode)
	if err != nil {
		return fmt.Errorf("invalid position mode: %w", err)
	}
	p.PositionMode = posMode

	marginMode, err := ParseMarginMode(aux.MarginMode)
	if err != nil {
		return fmt.Errorf("invalid margin mode: %w", err)
	}
	p.MarginMode = marginMode

	return nil
}
