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
	// StartTime usage of StartTime only has effect if EndTime is provided too
	StartTime *time.Time
	// EndTime usage of EndTime only has effect if StartTime is provided too
	EndTime *time.Time
	Skip    int64
	Limit   int64
}

type PositionHistoryResponse struct {
	BaseResponse
	Data struct {
		Positions []HistoricalPosition `json:"positionList"`
		Total     string               `json:"total"`
	} `json:"data"`
}

type HistoricalPosition struct {
	PositionID   string       `json:"positionId"`
	Symbol       Symbol       `json:"-"`
	MaxQty       float64      `json:"-"`
	EntryPrice   float64      `json:"-"`
	ClosePrice   float64      `json:"-"`
	LiqQty       float64      `json:"-"`
	Side         TradeSide    `json:"-"`
	PositionMode PositionMode `json:"-"`
	MarginMode   MarginMode   `json:"-"`
	Leverage     int          `json:"-"`
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
		Leverage     string `json:"leverage"`
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

	p.Symbol = ParseSymbol(aux.Symbol)

	if aux.Leverage != "" {
		lev, err := strconv.Atoi(aux.Leverage)
		if err == nil {
			p.Leverage = lev
		} else {
			return fmt.Errorf("failed to parse leverage: %w", err)
		}
	}

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

type PendingPositionParams struct {
	Symbol     Symbol
	PositionID string
}

type PendingPositionResponse struct {
	BaseResponse
	Data []PendingPosition `json:"data"`
}

type PendingPosition struct {
	PositionID    string       `json:"positionId"`
	Symbol        Symbol       `json:"-"`
	Qty           float64      `json:"-"`
	EntryValue    float64      `json:"-"`
	Side          TradeSide    `json:"-"`
	MarginMode    MarginMode   `json:"-"`
	PositionMode  PositionMode `json:"-"`
	Leverage      int          `json:"-"`
	Fees          float64      `json:"-"`
	Funding       float64      `json:"-"`
	RealizedPNL   float64      `json:"-"`
	Margin        float64      `json:"-"`
	UnrealizedPNL float64      `json:"-"`
	LiqPrice      float64      `json:"-"`
	MarginRate    float64      `json:"-"`
	AvgOpenPrice  float64      `json:"-"`
	CreateTime    time.Time    `json:"-"`
	ModifyTime    time.Time    `json:"-"`
}

func (p *PendingPosition) UnmarshalJSON(data []byte) error {
	type Alias PendingPosition
	aux := &struct {
		Symbol        string `json:"symbol"`
		Qty           string `json:"qty"`
		EntryValue    string `json:"entryValue"`
		Side          string `json:"side"`
		MarginMode    string `json:"marginMode"`
		PositionMode  string `json:"positionMode"`
		Leverage      int32  `json:"leverage"`
		Fees          string `json:"fees"`
		Funding       string `json:"funding"`
		RealizedPNL   string `json:"realizedPNL"`
		Margin        string `json:"margin"`
		UnrealizedPNL string `json:"unrealizedPNL"`
		LiqPrice      string `json:"liqPrice"`
		MarginRate    string `json:"marginRate"`
		AvgOpenPrice  string `json:"avgOpenPrice"`
		CreateTime    string `json:"ctime"`
		ModifyTime    string `json:"mtime"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	p.Symbol = ParseSymbol(aux.Symbol)
	p.Leverage = int(aux.Leverage)

	if aux.CreateTime != "" {
		createTime, err := strconv.ParseInt(aux.CreateTime, 10, 64)
		if err == nil {
			p.CreateTime = time.Unix(0, createTime*1000000)
		} else {
			return fmt.Errorf("invalid create time: %w", err)
		}
	}

	if aux.ModifyTime != "" {
		modifyTime, err := strconv.ParseInt(aux.ModifyTime, 10, 64)
		if err == nil {
			p.ModifyTime = time.Unix(0, modifyTime*1000000)
		} else {
			return fmt.Errorf("invalid modify time: %w", err)
		}
	}
	if aux.Qty != "" {
		qty, err := strconv.ParseFloat(aux.Qty, 64)
		if err != nil {
			return fmt.Errorf("failed to parse qty: %w", err)
		}
		p.Qty = qty
	}

	if aux.EntryValue != "" {
		entryValue, err := strconv.ParseFloat(aux.EntryValue, 64)
		if err != nil {
			return fmt.Errorf("failed to parse entryValue: %w", err)
		}
		p.EntryValue = entryValue
	}

	if aux.Fees != "" {
		fees, err := strconv.ParseFloat(aux.Fees, 64)
		if err != nil {
			return fmt.Errorf("failed to parse fees: %w", err)
		}
		p.Fees = fees
	}

	if aux.Funding != "" {
		funding, err := strconv.ParseFloat(aux.Funding, 64)
		if err != nil {
			return fmt.Errorf("failed to parse funding: %w", err)
		}
		p.Funding = funding
	}

	if aux.RealizedPNL != "" {
		realizedPNL, err := strconv.ParseFloat(aux.RealizedPNL, 64)
		if err != nil {
			return fmt.Errorf("failed to parse realizedPNL: %w", err)
		}
		p.RealizedPNL = realizedPNL
	}

	if aux.Margin != "" {
		margin, err := strconv.ParseFloat(aux.Margin, 64)
		if err != nil {
			return fmt.Errorf("failed to parse margin: %w", err)
		}
		p.Margin = margin
	}

	if aux.UnrealizedPNL != "" {
		unrealizedPNL, err := strconv.ParseFloat(aux.UnrealizedPNL, 64)
		if err != nil {
			return fmt.Errorf("failed to parse unrealizedPNL: %w", err)
		}
		p.UnrealizedPNL = unrealizedPNL
	}

	if aux.LiqPrice != "" {
		liqPrice, err := strconv.ParseFloat(aux.LiqPrice, 64)
		if err != nil {
			return fmt.Errorf("failed to parse liqPrice: %w", err)
		}
		p.LiqPrice = liqPrice
	}

	if aux.MarginRate != "" {
		marginRate, err := strconv.ParseFloat(aux.MarginRate, 64)
		if err != nil {
			return fmt.Errorf("failed to parse marginRate: %w", err)
		}
		p.MarginRate = marginRate
	}

	if aux.AvgOpenPrice != "" {
		avgOpenPrice, err := strconv.ParseFloat(aux.AvgOpenPrice, 64)
		if err != nil {
			return fmt.Errorf("failed to parse avgOpenPrice: %w", err)
		}
		p.AvgOpenPrice = avgOpenPrice
	}

	side, err := ParseTradeSide(aux.Side)
	if err != nil {
		return fmt.Errorf("invalid side: %w", err)
	}
	p.Side = side

	marginMode, err := ParseMarginMode(aux.MarginMode)
	if err != nil {
		return fmt.Errorf("invalid margin mode: %w", err)
	}
	p.MarginMode = marginMode

	positionMode, err := ParsePositionMode(aux.PositionMode)
	if err != nil {
		return fmt.Errorf("invalid position mode: %w", err)
	}
	p.PositionMode = positionMode

	return nil
}
