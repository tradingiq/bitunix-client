package model

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type AccountBalanceParams struct {
	MarginCoin string
}

type AccountBalanceResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"msg"`
	Data    AccountBalanceEntry `json:"data"`
}

type AccountBalanceEntry struct {
	MarginCoin             string       `json:"marginCoin"`
	Available              float64      `json:"-"`
	Frozen                 float64      `json:"-"`
	Margin                 float64      `json:"-"`
	Transfer               float64      `json:"-"`
	PositionMode           PositionMode `json:"-"`
	CrossUnrealizedPNL     float64      `json:"-"`
	IsolationUnrealizedPNL float64      `json:"-"`
	Bonus                  float64      `json:"-"`
}

func (a *AccountBalanceEntry) UnmarshalJSON(data []byte) error {
	type Alias AccountBalanceEntry
	aux := &struct {
		Available              string `json:"available"`
		Frozen                 string `json:"frozen"`
		Margin                 string `json:"margin"`
		Transfer               string `json:"transfer"`
		CrossUnrealizedPNL     string `json:"crossUnrealizedPNL"`
		IsolationUnrealizedPNL string `json:"isolationUnrealizedPNL"`
		Bonus                  string `json:"bonus"`
		PositionMode           string `json:"positionMode"`
		*Alias
	}{
		Alias: (*Alias)(a),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Available != "" {
		available, err := strconv.ParseFloat(aux.Available, 64)
		if err == nil {
			a.Available = available
		} else {
			return fmt.Errorf("invalid available amount: %w", err)
		}
	}

	if aux.Frozen != "" {
		frozen, err := strconv.ParseFloat(aux.Frozen, 64)
		if err == nil {
			a.Frozen = frozen
		} else {
			return fmt.Errorf("invalid frozen amount: %w", err)
		}
	}

	if aux.Margin != "" {
		margin, err := strconv.ParseFloat(aux.Margin, 64)
		if err == nil {
			a.Margin = margin
		} else {
			return fmt.Errorf("invalid margin amount: %w", err)
		}
	}

	if aux.Transfer != "" {
		transfer, err := strconv.ParseFloat(aux.Transfer, 64)
		if err == nil {
			a.Transfer = transfer
		} else {
			return fmt.Errorf("invalid transfer amount: %w", err)
		}
	}

	if aux.CrossUnrealizedPNL != "" {
		crossPNL, err := strconv.ParseFloat(aux.CrossUnrealizedPNL, 64)
		if err == nil {
			a.CrossUnrealizedPNL = crossPNL
		} else {
			return fmt.Errorf("invalid cross unrealized PNL: %w", err)
		}
	}

	if aux.IsolationUnrealizedPNL != "" {
		isolationPNL, err := strconv.ParseFloat(aux.IsolationUnrealizedPNL, 64)
		if err == nil {
			a.IsolationUnrealizedPNL = isolationPNL
		} else {
			return fmt.Errorf("invalid isolation unrealized PNL: %w", err)
		}
	}

	if aux.Bonus != "" {
		bonus, err := strconv.ParseFloat(aux.Bonus, 64)
		if err == nil {
			a.Bonus = bonus
		} else {
			return fmt.Errorf("invalid bonus amount: %w", err)
		}
	}

	positionMode, err := ParsePositionMode(aux.PositionMode)
	if err != nil {
		return fmt.Errorf("invalid position mode: %w", err)
	}
	a.PositionMode = positionMode

	return nil
}
