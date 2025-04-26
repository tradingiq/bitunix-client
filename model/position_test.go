package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestHistoricalPositionUnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"positionId": "pos123",
		"symbol": "BTCUSDT",
		"maxQty": "1.5",
		"entryPrice": "50000",
		"closePrice": "52000",
		"liqQty": "0",
		"side": "BUY",
		"positionMode": "ONE_WAY",
		"marginMode": "CROSS",
		"leverage": "10",
		"fee": "5",
		"funding": "1.25",
		"realizedPNL": "3000",
		"liqPrice": "45000",
		"ctime": "1642521600000",
		"mtime": "1642608000000"
	}`

	var pos HistoricalPosition
	err := json.Unmarshal([]byte(jsonStr), &pos)
	if err != nil {
		t.Fatalf("Failed to unmarshal historical position: %v", err)
	}

	if pos.PositionID != "pos123" {
		t.Errorf("Expected positionId 'pos123', got %s", pos.PositionID)
	}

	if pos.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol 'BTCUSDT', got %s", pos.Symbol)
	}

	if pos.MaxQty != 1.5 {
		t.Errorf("Expected maxQty 1.5, got %f", pos.MaxQty)
	}

	if pos.EntryPrice != 50000 {
		t.Errorf("Expected entryPrice 50000, got %f", pos.EntryPrice)
	}

	if pos.ClosePrice != 52000 {
		t.Errorf("Expected closePrice 52000, got %f", pos.ClosePrice)
	}

	if pos.Side != "BUY" {
		t.Errorf("Expected side 'BUY', got %s", pos.Side)
	}

	if pos.PositionMode != "ONE_WAY" {
		t.Errorf("Expected positionMode 'ONE_WAY', got %s", pos.PositionMode)
	}

	if pos.MarginMode != "CROSS" {
		t.Errorf("Expected marginMode 'CROSS', got %s", pos.MarginMode)
	}

	if pos.Leverage != "10" {
		t.Errorf("Expected leverage '10', got %s", pos.Leverage)
	}

	if pos.Fee != 5 {
		t.Errorf("Expected fee 5, got %f", pos.Fee)
	}

	if pos.Funding != 1.25 {
		t.Errorf("Expected funding 1.25, got %f", pos.Funding)
	}

	if pos.RealizedPNL != 3000 {
		t.Errorf("Expected realizedPNL 3000, got %f", pos.RealizedPNL)
	}

	if pos.LiqPrice != 45000 {
		t.Errorf("Expected liqPrice 45000, got %f", pos.LiqPrice)
	}

	expectedCtime := time.Date(2022, 1, 18, 16, 0, 0, 0, time.UTC)
	if !pos.Ctime.Equal(expectedCtime) {
		t.Errorf("Expected ctime %v, got %v", expectedCtime, pos.Ctime)
	}

	expectedMtime := time.Date(2022, 1, 19, 16, 0, 0, 0, time.UTC)
	if !pos.Mtime.Equal(expectedMtime) {
		t.Errorf("Expected mtime %v, got %v", expectedMtime, pos.Mtime)
	}
}
