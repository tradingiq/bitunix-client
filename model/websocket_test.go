package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestBalanceDetailUnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"coin": "USDT",
		"available": "1000.50",
		"frozen": "100.25",
		"isolationFrozen": "20.75",
		"crossFrozen": "30.25",
		"margin": "200.50",
		"isolationMargin": "50.75",
		"crossMargin": "150.25",
		"expMoney": "1200.75"
	}`

	var balance BalanceDetail
	err := json.Unmarshal([]byte(jsonStr), &balance)
	if err != nil {
		t.Fatalf("Failed to unmarshal balance detail: %v", err)
	}

	if balance.Coin != "USDT" {
		t.Errorf("Expected coin 'USDT', got %s", balance.Coin)
	}

	if balance.Available != 1000.50 {
		t.Errorf("Expected available 1000.50, got %f", balance.Available)
	}

	if balance.Frozen != 100.25 {
		t.Errorf("Expected frozen 100.25, got %f", balance.Frozen)
	}

	if balance.IsolationFrozen != 20.75 {
		t.Errorf("Expected isolationFrozen 20.75, got %f", balance.IsolationFrozen)
	}

	if balance.CrossFrozen != 30.25 {
		t.Errorf("Expected crossFrozen 30.25, got %f", balance.CrossFrozen)
	}

	if balance.Margin != 200.50 {
		t.Errorf("Expected margin 200.50, got %f", balance.Margin)
	}

	if balance.IsolationMargin != 50.75 {
		t.Errorf("Expected isolationMargin 50.75, got %f", balance.IsolationMargin)
	}

	if balance.CrossMargin != 150.25 {
		t.Errorf("Expected crossMargin 150.25, got %f", balance.CrossMargin)
	}

	if balance.ExpMoney != 1200.75 {
		t.Errorf("Expected expMoney 1200.75, got %f", balance.ExpMoney)
	}

	invalidJSON := `{
		"coin": "USDT",
		"available": "not-a-number",
		"frozen": "100.25",
		"isolationFrozen": "20.75",
		"crossFrozen": "30.25",
		"margin": "200.50",
		"isolationMargin": "50.75",
		"crossMargin": "150.25",
		"expMoney": "1200.75"
	}`

	err = json.Unmarshal([]byte(invalidJSON), &balance)
	if err == nil {
		t.Fatal("Expected error for invalid number format, got none")
	}
}

func TestBalanceResponseUnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"ch": "balance",
		"ts": 1642521600000,
		"data": {
			"coin": "USDT",
			"available": "1000.50",
			"frozen": "100.25",
			"isolationFrozen": "20.75",
			"crossFrozen": "30.25",
			"margin": "200.50",
			"isolationMargin": "50.75",
			"crossMargin": "150.25",
			"expMoney": "1200.75"
		}
	}`

	var response BalanceResponse
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal balance response: %v", err)
	}

	if response.Ch != "balance" {
		t.Errorf("Expected ch 'balance', got %s", response.Ch)
	}

	if response.Ts != 1642521600000 {
		t.Errorf("Expected ts 1642521600000, got %d", response.Ts)
	}

	if response.Data.Coin != "USDT" {
		t.Errorf("Expected data.coin 'USDT', got %s", response.Data.Coin)
	}

	if response.Data.Available != 1000.50 {
		t.Errorf("Expected data.available 1000.50, got %f", response.Data.Available)
	}

	if response.Data.Frozen != 100.25 {
		t.Errorf("Expected data.frozen 100.25, got %f", response.Data.Frozen)
	}
}

func TestPositionDataUnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"positionId": "12345",
		"symbol": "BTCUSDT",
		"event": "OPEN",
		"marginMode": "ISOLATION",
		"positionMode": "ONE_WAY",
		"side": "LONG",
		"leverage": "20",
		"margin": "500.25",
		"ctime": "1645123456789",
		"qty": "0.5",
		"entryValue": "20000.50",
		"realizedPNL": "100.25",
		"unrealizedPNL": "200.50",
		"funding": "1.25",
		"fee": "5.75"
	}`

	var position PositionData
	err := json.Unmarshal([]byte(jsonStr), &position)
	if err != nil {
		t.Fatalf("Failed to unmarshal position data: %v", err)
	}

	if position.PositionID != "12345" {
		t.Errorf("Expected positionId '12345', got %s", position.PositionID)
	}

	if position.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol 'BTCUSDT', got %s", position.Symbol)
	}

	if position.Event != PositionEventOpen {
		t.Errorf("Expected event 'OPEN', got %s", position.Event)
	}

	if position.MarginMode != MarginModeIsolation {
		t.Errorf("Expected marginMode 'ISOLATION', got %s", position.MarginMode)
	}

	if position.PositionMode != PositionModeOneWay {
		t.Errorf("Expected positionMode 'ONE_WAY', got %s", position.PositionMode)
	}

	if position.Side != PositionSideLong {
		t.Errorf("Expected side 'LONG', got %s", position.Side)
	}

	if position.Leverage != 20 {
		t.Errorf("Expected leverage 20, got %d", position.Leverage)
	}

	expectedTime := time.Unix(0, 1645123456789*1000000)
	if !position.CreateTime.Equal(expectedTime) {
		t.Errorf("Expected createTime %v, got %v", expectedTime, position.CreateTime)
	}

	if position.Margin != 500.25 {
		t.Errorf("Expected margin 500.25, got %f", position.Margin)
	}

	if position.Quantity != 0.5 {
		t.Errorf("Expected quantity 0.5, got %f", position.Quantity)
	}

	if position.EntryValue != 20000.50 {
		t.Errorf("Expected entryValue 20000.50, got %f", position.EntryValue)
	}

	if position.RealizedPNL != 100.25 {
		t.Errorf("Expected realizedPNL 100.25, got %f", position.RealizedPNL)
	}

	if position.UnrealizedPNL != 200.50 {
		t.Errorf("Expected unrealizedPNL 200.50, got %f", position.UnrealizedPNL)
	}

	if position.Funding != 1.25 {
		t.Errorf("Expected funding 1.25, got %f", position.Funding)
	}

	if position.Fee != 5.75 {
		t.Errorf("Expected fee 5.75, got %f", position.Fee)
	}
}

func TestPositionChannelSubscriptionUnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"ch": "position",
		"ts": 1645123456000,
		"data": [
			{
				"positionId": "12345",
				"symbol": "BTCUSDT",
				"event": "OPEN",
				"marginMode": "ISOLATION",
				"positionMode": "ONE_WAY",
				"side": "LONG",
				"leverage": "20",
				"margin": "500.25",
				"ctime": "1645123456789",
				"qty": "0.5",
				"entryValue": "20000.50",
				"realizedPNL": "100.25",
				"unrealizedPNL": "200.50",
				"funding": "1.25",
				"fee": "5.75"
			},
			{
				"positionId": "67890",
				"symbol": "ETHUSDT",
				"event": "UPDATE",
				"marginMode": "CROSS",
				"positionMode": "HEDGE",
				"side": "SHORT",
				"leverage": "10",
				"margin": "300.75",
				"ctime": "1645123456999",
				"qty": "2.0",
				"entryValue": "3000.25",
				"realizedPNL": "50.75",
				"unrealizedPNL": "75.50",
				"funding": "0.75",
				"fee": "2.25"
			}
		]
	}`

	var subscription PositionChannelSubscription
	err := json.Unmarshal([]byte(jsonStr), &subscription)
	if err != nil {
		t.Fatalf("Failed to unmarshal position channel subscription: %v", err)
	}

	if subscription.Channel != "position" {
		t.Errorf("Expected channel 'position', got %s", subscription.Channel)
	}

	if subscription.TimeStamp != 1645123456000 {
		t.Errorf("Expected timestamp 1645123456000, got %d", subscription.TimeStamp)
	}

	if len(subscription.Data) != 2 {
		t.Fatalf("Expected 2 position data entries, got %d", len(subscription.Data))
	}

	if subscription.Data[0].PositionID != "12345" {
		t.Errorf("Expected first position ID '12345', got %s", subscription.Data[0].PositionID)
	}

	if subscription.Data[0].Symbol != "BTCUSDT" {
		t.Errorf("Expected first symbol 'BTCUSDT', got %s", subscription.Data[0].Symbol)
	}

	if subscription.Data[1].PositionID != "67890" {
		t.Errorf("Expected second position ID '67890', got %s", subscription.Data[1].PositionID)
	}

	if subscription.Data[1].Symbol != "ETHUSDT" {
		t.Errorf("Expected second symbol 'ETHUSDT', got %s", subscription.Data[1].Symbol)
	}
}

func TestEmptyStringHandling(t *testing.T) {
	jsonStr := `{
		"positionId": "12345",
		"symbol": "BTCUSDT",
		"event": "OPEN",
		"marginMode": "ISOLATION",
		"positionMode": "ONE_WAY",
		"side": "LONG",
		"leverage": "",
		"margin": "",
		"ctime": "",
		"qty": "",
		"entryValue": "",
		"realizedPNL": "",
		"unrealizedPNL": "",
		"funding": "",
		"fee": ""
	}`

	var position PositionData
	err := json.Unmarshal([]byte(jsonStr), &position)
	if err != nil {
		t.Fatalf("Failed to handle empty strings: %v", err)
	}

	if position.Leverage != 0 {
		t.Errorf("Expected leverage to default to 0 for empty string, got %d", position.Leverage)
	}

	if !position.CreateTime.IsZero() {
		t.Errorf("Expected create time to default to zero time for empty string, got %v", position.CreateTime)
	}

	if position.Margin != 0 {
		t.Errorf("Expected margin to default to 0 for empty string, got %f", position.Margin)
	}

	if position.Quantity != 0 {
		t.Errorf("Expected quantity to default to 0 for empty string, got %f", position.Quantity)
	}

	if position.EntryValue != 0 {
		t.Errorf("Expected entryValue to default to 0 for empty string, got %f", position.EntryValue)
	}
}

func TestMalformedPositionData(t *testing.T) {
	testCases := []struct {
		name    string
		jsonStr string
		field   string
	}{
		{
			name: "Invalid event",
			jsonStr: `{
				"positionId": "12345",
				"symbol": "BTCUSDT",
				"event": "INVALID_EVENT",
				"marginMode": "ISOLATION",
				"positionMode": "ONE_WAY",
				"side": "LONG"
			}`,
			field: "event",
		},
		{
			name: "Invalid margin mode",
			jsonStr: `{
				"positionId": "12345",
				"symbol": "BTCUSDT",
				"event": "OPEN",
				"marginMode": "INVALID_MODE",
				"positionMode": "ONE_WAY",
				"side": "LONG"
			}`,
			field: "margin mode",
		},
		{
			name: "Invalid position mode",
			jsonStr: `{
				"positionId": "12345",
				"symbol": "BTCUSDT",
				"event": "OPEN",
				"marginMode": "ISOLATION",
				"positionMode": "INVALID_MODE",
				"side": "LONG"
			}`,
			field: "position mode",
		},
		{
			name: "Invalid side",
			jsonStr: `{
				"positionId": "12345",
				"symbol": "BTCUSDT",
				"event": "OPEN",
				"marginMode": "ISOLATION",
				"positionMode": "ONE_WAY",
				"side": "INVALID_SIDE"
			}`,
			field: "side",
		},
		{
			name: "Invalid leverage format",
			jsonStr: `{
				"positionId": "12345",
				"symbol": "BTCUSDT",
				"event": "OPEN",
				"marginMode": "ISOLATION",
				"positionMode": "ONE_WAY",
				"side": "LONG",
				"leverage": "not-a-number"
			}`,
			field: "leverage",
		},
		{
			name: "Invalid create time format",
			jsonStr: `{
				"positionId": "12345",
				"symbol": "BTCUSDT",
				"event": "OPEN",
				"marginMode": "ISOLATION",
				"positionMode": "ONE_WAY",
				"side": "LONG",
				"ctime": "not-a-timestamp"
			}`,
			field: "time",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var position PositionData
			err := json.Unmarshal([]byte(tc.jsonStr), &position)
			if err == nil {
				t.Fatalf("Expected error for invalid %s, got none", tc.field)
			}
		})
	}
}
