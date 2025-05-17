package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestPendingPositionUnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"positionId": "12345678",
		"symbol": "BTCUSDT",
		"qty": "0.5",
		"entryValue": "30000",
		"side": "BUY",
		"positionMode": "HEDGE",
		"marginMode": "ISOLATION",
		"leverage": 100,
		"fees": "1.5",
		"funding": "0.25",
		"realizedPNL": "150.75",
		"margin": "300",
		"unrealizedPNL": "75.5",
		"liqPrice": "29100.5",
		"marginRate": "0.1",
		"avgOpenPrice": "60000",
		"ctime": "1659076670000",
		"mtime": "1659086670000"
	}`

	var pos PendingPosition
	err := json.Unmarshal([]byte(jsonStr), &pos)
	if err != nil {
		t.Fatalf("Failed to unmarshal pending position: %v", err)
	}

	if pos.PositionID != "12345678" {
		t.Errorf("Expected positionId '12345678', got %s", pos.PositionID)
	}

	if pos.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol 'BTCUSDT', got %s", pos.Symbol)
	}

	if pos.Qty != 0.5 {
		t.Errorf("Expected qty 0.5, got %f", pos.Qty)
	}

	if pos.EntryValue != 30000 {
		t.Errorf("Expected entryValue 30000, got %f", pos.EntryValue)
	}

	if pos.Side != TradeSideBuy {
		t.Errorf("Expected side BUY, got %s", pos.Side)
	}

	if pos.PositionMode != PositionModeHedge {
		t.Errorf("Expected position mode HEDGE, got %s", pos.PositionMode)
	}

	if pos.MarginMode != MarginModeIsolation {
		t.Errorf("Expected margin mode ISOLATION, got %s", pos.MarginMode)
	}

	if pos.Leverage != 100 {
		t.Errorf("Expected leverage 100, got %d", pos.Leverage)
	}

	if pos.Fees != 1.5 {
		t.Errorf("Expected fees 1.5, got %f", pos.Fees)
	}

	if pos.Funding != 0.25 {
		t.Errorf("Expected funding 0.25, got %f", pos.Funding)
	}

	if pos.RealizedPNL != 150.75 {
		t.Errorf("Expected realizedPNL 150.75, got %f", pos.RealizedPNL)
	}

	if pos.Margin != 300 {
		t.Errorf("Expected margin 300, got %f", pos.Margin)
	}

	if pos.UnrealizedPNL != 75.5 {
		t.Errorf("Expected unrealizedPNL 75.5, got %f", pos.UnrealizedPNL)
	}

	if pos.LiqPrice != 29100.5 {
		t.Errorf("Expected liqPrice 29100.5, got %f", pos.LiqPrice)
	}

	if pos.MarginRate != 0.1 {
		t.Errorf("Expected marginRate 0.1, got %f", pos.MarginRate)
	}

	if pos.AvgOpenPrice != 60000 {
		t.Errorf("Expected avgOpenPrice 60000, got %f", pos.AvgOpenPrice)
	}

	expectedCtime := time.Unix(0, 1659076670000*1000000)
	if !pos.CreateTime.Equal(expectedCtime) {
		t.Errorf("Expected ctime %v, got %v", expectedCtime, pos.CreateTime)
	}

	expectedMtime := time.Unix(0, 1659086670000*1000000)
	if !pos.ModifyTime.Equal(expectedMtime) {
		t.Errorf("Expected mtime %v, got %v", expectedMtime, pos.ModifyTime)
	}
}

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

	if pos.Leverage != 10 {
		t.Errorf("Expected leverage '10', got %d", pos.Leverage)
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
