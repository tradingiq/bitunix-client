package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestHistoricalTradeUnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"tradeId": "trade123",
		"orderId": "order123",
		"symbol": "BTCUSDT",
		"qty": "1.5",
		"positionMode": "ONE_WAY",
		"marginMode": "CROSS",
		"leverage": 10,
		"price": "50000",
		"side": "BUY",
		"orderType": "LIMIT",
		"effect": "GTC",
		"clientId": "client123",
		"reduceOnly": false,
		"fee": "5",
		"realizedPNL": "3000",
		"ctime": "1642521600000",
		"roleType": "TAKER"
	}`

	var trade HistoricalTrade
	err := json.Unmarshal([]byte(jsonStr), &trade)
	if err != nil {
		t.Fatalf("Failed to unmarshal historical trade: %v", err)
	}

	if trade.TradeID != "trade123" {
		t.Errorf("Expected tradeId 'trade123', got %s", trade.TradeID)
	}

	if trade.OrderID != "order123" {
		t.Errorf("Expected orderId 'order123', got %s", trade.OrderID)
	}

	if trade.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol 'BTCUSDT', got %s", trade.Symbol)
	}

	if trade.Quantity != 1.5 {
		t.Errorf("Expected quantity 1.5, got %f", trade.Quantity)
	}

	if trade.PositionMode != PositionModeOneWay {
		t.Errorf("Expected positionMode 'ONE_WAY', got %s", trade.PositionMode)
	}

	if trade.MarginMode != MarginModeCross {
		t.Errorf("Expected marginMode 'CROSS', got %s", trade.MarginMode)
	}

	if trade.Leverage != 10 {
		t.Errorf("Expected leverage 10, got %d", trade.Leverage)
	}

	if trade.Price != 50000 {
		t.Errorf("Expected price 50000, got %f", trade.Price)
	}

	if trade.Side != TradeSideBuy {
		t.Errorf("Expected side 'OPEN', got %s", trade.Side)
	}

	if trade.OrderType != OrderTypeLimit {
		t.Errorf("Expected orderType 'LIMIT', got %s", trade.OrderType)
	}

	if trade.Effect != "GTC" {
		t.Errorf("Expected effect 'GTC', got %s", trade.Effect)
	}

	if trade.ClientID != "client123" {
		t.Errorf("Expected clientId 'client123', got %s", trade.ClientID)
	}

	if trade.ReduceOnly != false {
		t.Errorf("Expected reduceOnly false, got %v", trade.ReduceOnly)
	}

	if trade.Fee != 5 {
		t.Errorf("Expected fee 5, got %f", trade.Fee)
	}

	if trade.RealizedPNL != 3000 {
		t.Errorf("Expected realizedPNL 3000, got %f", trade.RealizedPNL)
	}

	expectedTime := time.Date(2022, 1, 18, 16, 0, 0, 0, time.UTC)
	if !trade.CreateTime.Equal(expectedTime) {
		t.Errorf("Expected createTime %v, got %v", expectedTime, trade.CreateTime)
	}

	if trade.RoleType != TradeRoleTypeTaker {
		t.Errorf("Expected roleType 'TAKER', got %s", trade.RoleType)
	}
}
