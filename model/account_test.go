package model

import (
	"encoding/json"
	"testing"
)

func TestAccountBalanceEntry_UnmarshalJSON(t *testing.T) {

	jsonData := `{
		"marginCoin": "USDT",
		"available": "1000.50",
		"frozen": "100.25",
		"margin": "50.75",
		"transfer": "900.25",
		"positionMode": "HEDGE",
		"crossUnrealizedPNL": "25.50",
		"isolationUnrealizedPNL": "10.25",
		"bonus": "5.75"
	}`

	var entry AccountBalanceEntry
	err := json.Unmarshal([]byte(jsonData), &entry)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling account balance: %v", err)
	}

	if entry.MarginCoin != "USDT" {
		t.Errorf("unexpected marginCoin: %s", entry.MarginCoin)
	}

	if entry.Available != 1000.50 {
		t.Errorf("unexpected available amount: %f", entry.Available)
	}

	if entry.Frozen != 100.25 {
		t.Errorf("unexpected frozen amount: %f", entry.Frozen)
	}

	if entry.Margin != 50.75 {
		t.Errorf("unexpected margin amount: %f", entry.Margin)
	}

	if entry.Transfer != 900.25 {
		t.Errorf("unexpected transfer amount: %f", entry.Transfer)
	}

	if entry.PositionMode != TradePositionModeHedge {
		t.Errorf("unexpected position mode: %s", entry.PositionMode)
	}

	if entry.CrossUnrealizedPNL != 25.50 {
		t.Errorf("unexpected cross unrealized PNL: %f", entry.CrossUnrealizedPNL)
	}

	if entry.IsolationUnrealizedPNL != 10.25 {
		t.Errorf("unexpected isolation unrealized PNL: %f", entry.IsolationUnrealizedPNL)
	}

	if entry.Bonus != 5.75 {
		t.Errorf("unexpected bonus: %f", entry.Bonus)
	}

	invalidJSON := `{
		"marginCoin": "USDT",
		"available": "not-a-number",
		"frozen": "0",
		"margin": "0",
		"transfer": "0",
		"positionMode": "HEDGE",
		"crossUnrealizedPNL": "0",
		"isolationUnrealizedPNL": "0",
		"bonus": "0"
	}`

	err = json.Unmarshal([]byte(invalidJSON), &entry)
	if err == nil {
		t.Fatal("expected error for invalid number format, got none")
	}
}
