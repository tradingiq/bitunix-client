package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type BalanceResponse struct {
	Ch   string        `json:"ch"`
	Ts   int64         `json:"ts"`
	Data BalanceDetail `json:"data"`
}

type BalanceDetail struct {
	Coin            string  `json:"coin"`
	Available       float64 `json:"-"`
	Frozen          float64 `json:"-"`
	IsolationFrozen float64 `json:"-"`
	CrossFrozen     float64 `json:"-"`
	Margin          float64 `json:"-"`
	IsolationMargin float64 `json:"-"`
	CrossMargin     float64 `json:"-"`
	ExpMoney        float64 `json:"-"`
}

func (p *BalanceDetail) UnmarshalJSON(data []byte) error {
	type Alias BalanceDetail
	aux := &struct {
		Available       string `json:"available"`
		Frozen          string `json:"frozen"`
		IsolationFrozen string `json:"isolationFrozen"`
		CrossFrozen     string `json:"crossFrozen"`
		Margin          string `json:"margin"`
		IsolationMargin string `json:"isolationMargin"`
		CrossMargin     string `json:"crossMargin"`
		ExpMoney        string `json:"expMoney"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	available, err := strconv.ParseFloat(aux.Available, 64)
	if err == nil {
		p.Available = available
	} else {
		return fmt.Errorf("failed to parse available: %w", err)
	}

	frozen, err := strconv.ParseFloat(aux.Frozen, 64)
	if err == nil {
		p.Frozen = frozen
	} else {
		return fmt.Errorf("failed to parse frozen: %w", err)
	}

	isolationFrozen, err := strconv.ParseFloat(aux.IsolationFrozen, 64)
	if err == nil {
		p.IsolationFrozen = isolationFrozen
	} else {
		return fmt.Errorf("failed to parse IisolationFrozen: %w", err)
	}

	crossFrozen, err := strconv.ParseFloat(aux.CrossFrozen, 64)
	if err == nil {
		p.CrossFrozen = crossFrozen
	} else {
		return fmt.Errorf("failed to parse crossFrozen: %w", err)
	}

	Margin, err := strconv.ParseFloat(aux.Margin, 64)
	if err == nil {
		p.Margin = Margin
	} else {
		return fmt.Errorf("failed to parse Margin: %w", err)
	}

	isolationMargin, err := strconv.ParseFloat(aux.IsolationMargin, 64)
	if err == nil {
		p.IsolationMargin = isolationMargin
	} else {
		return fmt.Errorf("failed to parse isolationMargin: %w", err)
	}

	crossMargin, err := strconv.ParseFloat(aux.CrossMargin, 64)
	if err == nil {
		p.CrossMargin = crossMargin
	} else {
		return fmt.Errorf("failed to parse crossMargin: %w", err)
	}

	expMoney, err := strconv.ParseFloat(aux.ExpMoney, 64)
	if err == nil {
		p.ExpMoney = expMoney
	} else {
		return fmt.Errorf("failed to parse expMoney: %w", err)
	}

	return nil
}

type PositionData struct {
	Event         PositionEvent `json:"-"`
	PositionID    string        `json:"positionId"`
	MarginMode    MarginMode    `json:"-"`
	PositionMode  PositionMode  `json:"-"`
	Side          PositionSide  `json:"-"`
	Leverage      int           `json:"-"`
	Margin        float64       `json:"-"`
	CreateTime    time.Time     `json:"-"`
	Quantity      float64       `json:"-"`
	EntryValue    float64       `json:"-"`
	Symbol        string        `json:"symbol"`
	RealizedPNL   float64       `json:"-"`
	UnrealizedPNL float64       `json:"-"`
	Funding       float64       `json:"-"`
	Fee           float64       `json:"-"`
}

func (p *PositionData) UnmarshalJSON(data []byte) error {
	type Alias PositionData
	aux := &struct {
		Event         string `json:"event"`
		MarginMode    string `json:"marginMode"`
		PositionMode  string `json:"positionMode"`
		Side          string `json:"side"`
		Leverage      string `json:"leverage"`
		Margin        string `json:"margin"`
		CreateTime    string `json:"ctime"`
		Quantity      string `json:"qty"`
		EntryValue    string `json:"entryValue"`
		RealizedPNL   string `json:"realizedPNL"`
		UnrealizedPNL string `json:"unrealizedPNL"`
		Funding       string `json:"funding"`
		Fee           string `json:"fee"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	event, err := ParsePositionEvent(aux.Event)
	if err != nil {
		return fmt.Errorf("invalid position event: %w", err)
	}
	p.Event = event

	side, err := ParsePositionSide(aux.Side)
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

	if aux.CreateTime != "" {
		createTime, err := strconv.ParseInt(aux.CreateTime, 10, 64)
		if err == nil {
			p.CreateTime = time.Unix(0, createTime*1000000)
		} else {
			return fmt.Errorf("invalid create time: %w", err)
		}
	}

	if aux.Leverage != "" {
		val, err := strconv.Atoi(aux.Leverage)
		if err != nil {
			return fmt.Errorf("failed to parse leverage: %w", err)
		}
		p.Leverage = val
	}

	if aux.Margin != "" {
		margin, err := strconv.ParseFloat(aux.Margin, 64)
		if err == nil {
			p.Margin = margin
		} else {
			return fmt.Errorf("failed to parse margin: %w", err)
		}
	}

	if aux.Quantity != "" {
		qty, err := strconv.ParseFloat(aux.Quantity, 64)
		if err == nil {
			p.Quantity = qty
		} else {
			return fmt.Errorf("failed to parse qty: %w", err)
		}
	}

	if aux.EntryValue != "" {
		val, err := strconv.ParseFloat(aux.EntryValue, 64)
		if err == nil {
			p.EntryValue = val
		} else {
			return fmt.Errorf("failed to parse EntryValue: %w", err)
		}
	}

	if aux.RealizedPNL != "" {
		val, err := strconv.ParseFloat(aux.RealizedPNL, 64)
		if err == nil {
			p.RealizedPNL = val
		} else {
			return fmt.Errorf("failed to parse RealizedPNL: %w", err)
		}
	}

	if aux.UnrealizedPNL != "" {
		val, err := strconv.ParseFloat(aux.UnrealizedPNL, 64)
		if err == nil {
			p.UnrealizedPNL = val
		} else {
			return fmt.Errorf("failed to parse UnrealizedPNL: %w", err)
		}
	}

	if aux.Funding != "" {
		val, err := strconv.ParseFloat(aux.Funding, 64)
		if err == nil {
			p.Funding = val
		} else {
			return fmt.Errorf("failed to parse Funding: %w", err)
		}
	}

	if aux.Fee != "" {
		val, err := strconv.ParseFloat(aux.Fee, 64)
		if err == nil {
			p.Fee = val
		} else {
			return fmt.Errorf("failed to parse Fee: %w", err)
		}
	}

	return nil
}

type PositionChannelSubscription struct {
	Channel   string         `json:"ch"`
	TimeStamp int64          `json:"ts"`
	Data      []PositionData `json:"data"`
}
