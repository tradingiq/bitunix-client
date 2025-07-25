package model

import (
	"encoding/json"
	"testing"
)

func TestPendingTPSLOrderUnmarshal(t *testing.T) {
	jsonData := `{
		"id": "123",
		"positionId": "12345678",
		"symbol": "BTCUSDT",
		"base": "BTC",
		"quote": "USDT",
		"tpPrice": "50000",
		"tpStopType": "LAST_PRICE",
		"slPrice": "45000",
		"slStopType": "MARK_PRICE",
		"tpOrderType": "LIMIT",
		"tpOrderPrice": "50100",
		"slOrderType": "MARKET",
		"slOrderPrice": "",
		"tpQty": "0.01",
		"slQty": "0.02"
	}`

	var order PendingTPSLOrder
	err := json.Unmarshal([]byte(jsonData), &order)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if order.ID != "123" {
		t.Errorf("Expected ID '123', got %s", order.ID)
	}
	if order.PositionID != "12345678" {
		t.Errorf("Expected PositionID '12345678', got %s", order.PositionID)
	}
	if order.Symbol != "BTCUSDT" {
		t.Errorf("Expected Symbol 'BTCUSDT', got %s", order.Symbol)
	}
	if order.Base != "BTC" {
		t.Errorf("Expected Base 'BTC', got %s", order.Base)
	}
	if order.Quote != "USDT" {
		t.Errorf("Expected Quote 'USDT', got %s", order.Quote)
	}

	if order.TpPrice == nil || *order.TpPrice != 50000 {
		t.Errorf("Expected TpPrice 50000, got %v", order.TpPrice)
	}
	if order.TpStopType == nil || *order.TpStopType != StopTypeLastPrice {
		t.Errorf("Expected TpStopType LAST_PRICE, got %v", order.TpStopType)
	}
	if order.TpOrderType == nil || *order.TpOrderType != OrderTypeLimit {
		t.Errorf("Expected TpOrderType LIMIT, got %v", order.TpOrderType)
	}
	if order.TpOrderPrice == nil || *order.TpOrderPrice != 50100 {
		t.Errorf("Expected TpOrderPrice 50100, got %v", order.TpOrderPrice)
	}
	if order.TpQty == nil || *order.TpQty != 0.01 {
		t.Errorf("Expected TpQty 0.01, got %v", order.TpQty)
	}

	if order.SlPrice == nil || *order.SlPrice != 45000 {
		t.Errorf("Expected SlPrice 45000, got %v", order.SlPrice)
	}
	if order.SlStopType == nil || *order.SlStopType != StopTypeMarkPrice {
		t.Errorf("Expected SlStopType MARK_PRICE, got %v", order.SlStopType)
	}
	if order.SlOrderType == nil || *order.SlOrderType != OrderTypeMarket {
		t.Errorf("Expected SlOrderType MARKET, got %v", order.SlOrderType)
	}
	if order.SlOrderPrice != nil {
		t.Errorf("Expected SlOrderPrice nil for market order, got %v", order.SlOrderPrice)
	}
	if order.SlQty == nil || *order.SlQty != 0.02 {
		t.Errorf("Expected SlQty 0.02, got %v", order.SlQty)
	}
}

func TestPendingTPSLOrderUnmarshalPartial(t *testing.T) {

	jsonData := `{
		"id": "456",
		"positionId": "87654321",
		"symbol": "ETHUSDT",
		"base": "ETH",
		"quote": "USDT",
		"tpPrice": "3000",
		"tpStopType": "MARK_PRICE",
		"tpOrderType": "MARKET",
		"tpQty": "1.5"
	}`

	var order PendingTPSLOrder
	err := json.Unmarshal([]byte(jsonData), &order)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if order.TpPrice == nil || *order.TpPrice != 3000 {
		t.Errorf("Expected TpPrice 3000, got %v", order.TpPrice)
	}
	if order.TpStopType == nil || *order.TpStopType != StopTypeMarkPrice {
		t.Errorf("Expected TpStopType MARK_PRICE, got %v", order.TpStopType)
	}

	if order.SlPrice != nil {
		t.Errorf("Expected SlPrice nil, got %v", order.SlPrice)
	}
	if order.SlStopType != nil {
		t.Errorf("Expected SlStopType nil, got %v", order.SlStopType)
	}
}

func TestPendingTPSLOrderUnmarshalErrors(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		expectedErr string
	}{
		{
			name:        "Invalid JSON",
			jsonData:    `{"id": "123"`,
			expectedErr: "unexpected end of JSON input",
		},
		{
			name: "Invalid tpPrice",
			jsonData: `{
				"id": "123",
				"symbol": "BTCUSDT",
				"tpPrice": "invalid"
			}`,
			expectedErr: "invalid tp price",
		},
		{
			name: "Invalid slPrice",
			jsonData: `{
				"id": "123",
				"symbol": "BTCUSDT",
				"slPrice": "not-a-number"
			}`,
			expectedErr: "invalid sl price",
		},
		{
			name: "Invalid tpQty",
			jsonData: `{
				"id": "123",
				"symbol": "BTCUSDT",
				"tpQty": "abc"
			}`,
			expectedErr: "invalid tp qty",
		},
		{
			name: "Invalid tpStopType",
			jsonData: `{
				"id": "123",
				"symbol": "BTCUSDT",
				"tpStopType": "INVALID_TYPE"
			}`,
			expectedErr: "invalid tp stop type",
		},
		{
			name: "Invalid tpOrderType",
			jsonData: `{
				"id": "123",
				"symbol": "BTCUSDT",
				"tpOrderType": "INVALID_ORDER"
			}`,
			expectedErr: "invalid tp order type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var order PendingTPSLOrder
			err := json.Unmarshal([]byte(tt.jsonData), &order)
			if err == nil {
				t.Error("Expected error but got none")
			} else if !contains(err.Error(), tt.expectedErr) {
				t.Errorf("Expected error containing '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
		len(s) > len(substr) && contains(s[1:], substr)
}
