package model

import (
	"encoding/json"
	"fmt"
	"strconv"
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
