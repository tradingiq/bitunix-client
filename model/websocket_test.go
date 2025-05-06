package model

import (
	"encoding/json"
	"fmt"
	"strings"
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

	var balance BalanceEvent
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

	var response BalanceChannelMessage
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal balance BaseResponse: %v", err)
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
		"ctime": "2025-04-26T17:06:13.155255000Z",
		"qty": "0.5",
		"entryValue": "20000.50",
		"realizedPNL": "100.25",
		"unrealizedPNL": "200.50",
		"funding": "1.25",
		"fee": "5.75"
	}`

	var position PositionEvent
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

	expectedTime, _ := time.Parse(time.RFC3339Nano, "2025-04-26T17:06:13.155255000Z")
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
		"data": 
			{
				"positionId": "12345",
				"symbol": "BTCUSDT",
				"event": "OPEN",
				"marginMode": "ISOLATION",
				"positionMode": "ONE_WAY",
				"side": "LONG",
				"leverage": "20",
				"margin": "500.25",
				"ctime": "2025-04-26T17:06:13.155255000Z",
				"qty": "0.5",
				"entryValue": "20000.50",
				"realizedPNL": "100.25",
				"unrealizedPNL": "200.50",
				"funding": "1.25",
				"fee": "5.75"
			}
	}`

	var subscription PositionChannelMessage
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

	if subscription.Data.PositionID != "12345" {
		t.Errorf("Expected first position ID '12345', got %s", subscription.Data.PositionID)
	}

	if subscription.Data.Symbol != "BTCUSDT" {
		t.Errorf("Expected first symbol 'BTCUSDT', got %s", subscription.Data.Symbol)
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

	var position PositionEvent
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
			var position PositionEvent
			err := json.Unmarshal([]byte(tc.jsonStr), &position)
			if err == nil {
				t.Fatalf("Expected error for invalid %s, got none", tc.field)
			}
		})
	}
}

func TestTPSLOrderUnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"positionId": "12345",
		"orderId": "67890",
		"symbol": "BTCUSDT",
		"event": "CREATE",
		"leverage": "20",
		"side": "BUY",
		"positionMode": "ONE_WAY",
		"status": "NEW",
		"ctime": "2025-04-26T17:06:13.155255000Z",
		"type": "POSITION_TPSL",
		"tpQty": "0.5",
		"slQty": "0.5",
		"tpStopType": "MARK_PRICE",
		"tpPrice": "50000.00",
		"tpOrderType": "LIMIT",
		"tpOrderPrice": "49900.00",
		"slStopType": "LAST_PRICE",
		"slPrice": "45000.00",
		"slOrderType": "MARKET",
		"slOrderPrice": "0.00"
	}`

	var tpslOrder TpSlOrderEvent
	err := json.Unmarshal([]byte(jsonStr), &tpslOrder)
	if err != nil {
		t.Fatalf("Failed to unmarshal TPSL order: %v", err)
	}

	if tpslOrder.PositionID != "12345" {
		t.Errorf("Expected positionId '12345', got %s", tpslOrder.PositionID)
	}

	if tpslOrder.OrderID != "67890" {
		t.Errorf("Expected orderId '67890', got %s", tpslOrder.OrderID)
	}

	if tpslOrder.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol 'BTCUSDT', got %s", tpslOrder.Symbol)
	}

	if tpslOrder.Event != TPSLEventCreate {
		t.Errorf("Expected event CREATE, got %s", tpslOrder.Event)
	}

	if tpslOrder.Leverage != 20 {
		t.Errorf("Expected leverage 20, got %d", tpslOrder.Leverage)
	}

	if tpslOrder.Side != TradeSideBuy {
		t.Errorf("Expected side BUY, got %s", tpslOrder.Side)
	}

	if tpslOrder.PositionMode != PositionModeOneWay {
		t.Errorf("Expected positionMode ONE_WAY, got %s", tpslOrder.PositionMode)
	}

	if tpslOrder.Status != OrderStatusNew {
		t.Errorf("Expected status NEW, got %s", tpslOrder.Status)
	}

	expectedTime, _ := time.Parse(time.RFC3339Nano, "2025-04-26T17:06:13.155255000Z")
	if !tpslOrder.CreateTime.Equal(expectedTime) {
		t.Errorf("Expected createTime %v, got %v", expectedTime, tpslOrder.CreateTime)
	}

	if tpslOrder.Type != TPSLTypeFull {
		t.Errorf("Expected type POSITION_TPSL, got %s", tpslOrder.Type)
	}

	if tpslOrder.TPQuantity != 0.5 {
		t.Errorf("Expected tpQty 0.5, got %f", tpslOrder.TPQuantity)
	}

	if tpslOrder.TPStopType != StopTypeMarkPrice {
		t.Errorf("Expected tpStopType MARK_PRICE, got %s", tpslOrder.TPStopType)
	}

	if tpslOrder.TPPrice != 50000.00 {
		t.Errorf("Expected tpPrice 50000.00, got %f", tpslOrder.TPPrice)
	}

	if tpslOrder.TPOrderType != OrderTypeLimit {
		t.Errorf("Expected tpOrderType LIMIT, got %s", tpslOrder.TPOrderType)
	}

	if tpslOrder.TPOrderPrice != 49900.00 {
		t.Errorf("Expected tpOrderPrice 49900.00, got %f", tpslOrder.TPOrderPrice)
	}

	if tpslOrder.SLQuantity != 0.5 {
		t.Errorf("Expected slQty 0.5, got %f", tpslOrder.SLQuantity)
	}

	if tpslOrder.SLStopType != StopTypeLastPrice {
		t.Errorf("Expected slStopType LAST_PRICE, got %s", tpslOrder.SLStopType)
	}

	if tpslOrder.SLPrice != 45000.00 {
		t.Errorf("Expected slPrice 45000.00, got %f", tpslOrder.SLPrice)
	}

	if tpslOrder.SLOrderType != OrderTypeMarket {
		t.Errorf("Expected slOrderType MARKET, got %s", tpslOrder.SLOrderType)
	}

	if tpslOrder.SLOrderPrice != 0.00 {
		t.Errorf("Expected slOrderPrice 0.00, got %f", tpslOrder.SLOrderPrice)
	}
}

func TestTPSLOrderChannelResponseUnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"ch": "tpsl",
		"ts": 1645123456000,
		"data": 
			{
				"positionId": "12345",
				"orderId": "67890",
				"symbol": "BTCUSDT",
				"event": "CREATE",
				"leverage": "20",
				"side": "BUY",
				"positionMode": "ONE_WAY",
				"status": "NEW",
				"ctime": "2025-04-26T17:06:13.155255000Z",
				"type": "POSITION_TPSL",
				"tpQty": "0.5",
				"slQty": "0.5",
				"tpStopType": "MARK_PRICE",
				"tpPrice": "50000.00",
				"tpOrderType": "LIMIT",
				"tpOrderPrice": "49900.00",
				"slStopType": "LAST_PRICE",
				"slPrice": "45000.00",
				"slOrderType": "MARKET",
				"slOrderPrice": "0.00"
			}
		
	}`

	var response TpSlOrderChannelMessage
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal TPSL channel BaseResponse: %v", err)
	}

	if response.Channel != "tpsl" {
		t.Errorf("Expected channel 'tpsl', got %s", response.Channel)
	}

	if response.Timestamp != 1645123456000 {
		t.Errorf("Expected timestamp 1645123456000, got %d", response.Timestamp)
	}

	tpslOrder := response.Data
	if tpslOrder.PositionID != "12345" {
		t.Errorf("Expected positionId '12345', got %s", tpslOrder.PositionID)
	}

	if tpslOrder.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol 'BTCUSDT', got %s", tpslOrder.Symbol)
	}
}

func TestTPSLOrderEmptyFieldsHandling(t *testing.T) {
	jsonStr := `{
		"positionId": "12345",
		"orderId": "67890",
		"symbol": "BTCUSDT",
		"event": "CREATE",
		"leverage": "",
		"side": "BUY",
		"positionMode": "ONE_WAY",
		"status": "NEW",
		"ctime": "",
		"type": "POSITION_TPSL",
		"tpQty": "",
		"slQty": "",
		"tpStopType": "",
		"tpPrice": "",
		"tpOrderType": "",
		"tpOrderPrice": "",
		"slStopType": "",
		"slPrice": "",
		"slOrderType": "",
		"slOrderPrice": ""
	}`

	var tpslOrder TpSlOrderEvent
	err := json.Unmarshal([]byte(jsonStr), &tpslOrder)
	if err != nil {
		t.Fatalf("Failed to handle empty TPSL fields: %v", err)
	}

	if tpslOrder.Leverage != 0 {
		t.Errorf("Expected leverage to default to 0 for empty string, got %d", tpslOrder.Leverage)
	}

	if !tpslOrder.CreateTime.IsZero() {
		t.Errorf("Expected create time to default to zero time for empty string, got %v", tpslOrder.CreateTime)
	}

	if tpslOrder.TPQuantity != 0 {
		t.Errorf("Expected TPQuantity to default to 0 for empty string, got %f", tpslOrder.TPQuantity)
	}

	if tpslOrder.SLQuantity != 0 {
		t.Errorf("Expected SLQuantity to default to 0 for empty string, got %f", tpslOrder.SLQuantity)
	}

	if tpslOrder.TPPrice != 0 {
		t.Errorf("Expected TPPrice to default to 0 for empty string, got %f", tpslOrder.TPPrice)
	}

	if tpslOrder.TPOrderPrice != 0 {
		t.Errorf("Expected TPOrderPrice to default to 0 for empty string, got %f", tpslOrder.TPOrderPrice)
	}

	if tpslOrder.SLPrice != 0 {
		t.Errorf("Expected SLPrice to default to 0 for empty string, got %f", tpslOrder.SLPrice)
	}

	if tpslOrder.SLOrderPrice != 0 {
		t.Errorf("Expected SLOrderPrice to default to 0 for empty string, got %f", tpslOrder.SLOrderPrice)
	}
}

func TestTPSLOrderWithPartialFields(t *testing.T) {

	jsonStrTPOnly := `{
		"positionId": "12345",
		"orderId": "67890",
		"symbol": "BTCUSDT",
		"event": "CREATE",
		"leverage": "20",
		"side": "BUY",
		"positionMode": "ONE_WAY",
		"status": "NEW",
		"ctime": "2025-04-26T17:06:13.155255000Z",
		"type": "POSITION_TPSL",
		"tpQty": "0.5",
		"tpStopType": "MARK_PRICE",
		"tpPrice": "50000.00",
		"tpOrderType": "LIMIT",
		"tpOrderPrice": "49900.00"
	}`

	var tpslOrderTPOnly TpSlOrderEvent
	err := json.Unmarshal([]byte(jsonStrTPOnly), &tpslOrderTPOnly)
	if err != nil {
		t.Fatalf("Failed to unmarshal TPSL order with TP fields only: %v", err)
	}

	if tpslOrderTPOnly.TPQuantity != 0.5 {
		t.Errorf("Expected tpQty 0.5, got %f", tpslOrderTPOnly.TPQuantity)
	}

	if tpslOrderTPOnly.TPStopType != StopTypeMarkPrice {
		t.Errorf("Expected tpStopType MARK_PRICE, got %s", tpslOrderTPOnly.TPStopType)
	}

	if tpslOrderTPOnly.SLQuantity != 0 {
		t.Errorf("Expected slQty to be 0 when not provided, got %f", tpslOrderTPOnly.SLQuantity)
	}

	jsonStrSLOnly := `{
		"positionId": "12345",
		"orderId": "67890",
		"symbol": "BTCUSDT",
		"event": "CREATE",
		"leverage": "20",
		"side": "BUY",
		"positionMode": "ONE_WAY",
		"status": "NEW",
		"ctime": "2025-04-26T17:06:13.155255000Z",
		"type": "POSITION_TPSL",
		"slQty": "0.5",
		"slStopType": "LAST_PRICE",
		"slPrice": "45000.00",
		"slOrderType": "MARKET",
		"slOrderPrice": "0.00"
	}`

	var tpslOrderSLOnly TpSlOrderEvent
	err = json.Unmarshal([]byte(jsonStrSLOnly), &tpslOrderSLOnly)
	if err != nil {
		t.Fatalf("Failed to unmarshal TPSL order with SL fields only: %v", err)
	}

	if tpslOrderSLOnly.SLQuantity != 0.5 {
		t.Errorf("Expected slQty 0.5, got %f", tpslOrderSLOnly.SLQuantity)
	}

	if tpslOrderSLOnly.SLStopType != StopTypeLastPrice {
		t.Errorf("Expected slStopType LAST_PRICE, got %s", tpslOrderSLOnly.SLStopType)
	}

	if tpslOrderSLOnly.TPQuantity != 0 {
		t.Errorf("Expected tpQty to be 0 when not provided, got %f", tpslOrderSLOnly.TPQuantity)
	}
}

func TestMalformedTPSLOrderData(t *testing.T) {
	testCases := []struct {
		name    string
		jsonStr string
		field   string
	}{
		{
			name: "Invalid event",
			jsonStr: `{
				"positionId": "12345",
				"orderId": "67890",
				"symbol": "BTCUSDT",
				"event": "INVALID_EVENT",
				"side": "BUY",
				"positionMode": "ONE_WAY",
				"status": "NEW"
			}`,
			field: "event",
		},
		{
			name: "Invalid side",
			jsonStr: `{
				"positionId": "12345",
				"orderId": "67890",
				"symbol": "BTCUSDT",
				"event": "CREATE",
				"side": "INVALID_SIDE",
				"positionMode": "ONE_WAY",
				"status": "NEW"
			}`,
			field: "side",
		},
		{
			name: "Invalid position mode",
			jsonStr: `{
				"positionId": "12345",
				"orderId": "67890",
				"symbol": "BTCUSDT",
				"event": "CREATE",
				"side": "BUY",
				"positionMode": "INVALID_MODE",
				"status": "NEW"
			}`,
			field: "position mode",
		},
		{
			name: "Invalid order status",
			jsonStr: `{
				"positionId": "12345",
				"orderId": "67890",
				"symbol": "BTCUSDT",
				"event": "CREATE",
				"side": "BUY",
				"positionMode": "ONE_WAY",
				"status": "INVALID_STATUS"
			}`,
			field: "order status",
		},
		{
			name: "Invalid TPSL type",
			jsonStr: `{
				"positionId": "12345",
				"orderId": "67890",
				"symbol": "BTCUSDT",
				"event": "CREATE",
				"side": "BUY",
				"positionMode": "ONE_WAY",
				"status": "NEW",
				"type": "INVALID_TYPE"
			}`,
			field: "TPSL type",
		},
		{
			name: "Invalid leverage format",
			jsonStr: `{
				"positionId": "12345",
				"orderId": "67890",
				"symbol": "BTCUSDT",
				"event": "CREATE",
				"side": "BUY",
				"positionMode": "ONE_WAY",
				"status": "NEW",
				"leverage": "not-a-number"
			}`,
			field: "leverage",
		},
		{
			name: "Invalid TP quantity format",
			jsonStr: `{
				"positionId": "12345",
				"orderId": "67890",
				"symbol": "BTCUSDT",
				"event": "CREATE",
				"side": "BUY",
				"positionMode": "ONE_WAY",
				"status": "NEW",
				"tpQty": "not-a-number"
			}`,
			field: "TP quantity",
		},
		{
			name: "Invalid create time format",
			jsonStr: `{
				"positionId": "12345",
				"orderId": "67890",
				"symbol": "BTCUSDT",
				"event": "CREATE",
				"side": "BUY",
				"positionMode": "ONE_WAY",
				"status": "NEW",
				"ctime": "not-a-timestamp"
			}`,
			field: "create time",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var tpslOrder TpSlOrderEvent
			err := json.Unmarshal([]byte(tc.jsonStr), &tpslOrder)
			if err == nil {
				t.Fatalf("Expected error for invalid %s, got none", tc.field)
			}
		})
	}
}

func TestMarginCoinString(t *testing.T) {
	testCases := []struct {
		name     string
		input    MarginCoin
		expected string
	}{
		{
			name:     "Standard coin",
			input:    MarginCoin("USDT"),
			expected: "USDT",
		},
		{
			name:     "Lowercase coin",
			input:    MarginCoin("usdt"),
			expected: "usdt",
		},
		{
			name:     "Mixed case coin",
			input:    MarginCoin("uSdT"),
			expected: "uSdT",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.String()
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestMarginCoinNormalize(t *testing.T) {
	testCases := []struct {
		name     string
		input    MarginCoin
		expected MarginCoin
	}{
		{
			name:     "Already uppercase",
			input:    MarginCoin("USDT"),
			expected: MarginCoin("USDT"),
		},
		{
			name:     "Lowercase to uppercase",
			input:    MarginCoin("usdt"),
			expected: MarginCoin("USDT"),
		},
		{
			name:     "Mixed case to uppercase",
			input:    MarginCoin("uSdT"),
			expected: MarginCoin("USDT"),
		},
		{
			name:     "Trim spaces",
			input:    MarginCoin(" USDT "),
			expected: MarginCoin("USDT"),
		},
		{
			name:     "Trim spaces and convert case",
			input:    MarginCoin(" uSdT "),
			expected: MarginCoin("USDT"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.input.Normalize()
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestParseMarginCoin(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected MarginCoin
	}{
		{
			name:     "Standard coin format",
			input:    "USDT",
			expected: MarginCoin("USDT"),
		},
		{
			name:     "Lowercase input",
			input:    "usdt",
			expected: MarginCoin("USDT"),
		},
		{
			name:     "Mixed case input",
			input:    "uSdT",
			expected: MarginCoin("USDT"),
		},
		{
			name:     "With spaces",
			input:    " USDT ",
			expected: MarginCoin("USDT"),
		},
		{
			name:     "With spaces and mixed case",
			input:    " uSdT ",
			expected: MarginCoin("USDT"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := ParseMarginCoin(tc.input)
			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestMarginCoinWithCommonCoins(t *testing.T) {
	commonCoins := []string{
		"USDT", "BTC", "ETH", "SOL", "USDC", "BNB", "XRP", "ADA", "DOGE", "DOT",
	}

	for _, coin := range commonCoins {
		t.Run(coin, func(t *testing.T) {

			lowerCoin := strings.ToLower(coin)
			result := ParseMarginCoin(lowerCoin)
			if string(result) != coin {
				t.Errorf("Expected '%s', got '%s'", coin, result)
			}

			spacedCoin := " " + coin + " "
			result = ParseMarginCoin(spacedCoin)
			if string(result) != coin {
				t.Errorf("Expected '%s', got '%s'", coin, result)
			}
		})
	}
}

func TestKLineEventUnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"o": "50000.50",
		"h": "51000.25",
		"l": "49500.75",
		"c": "50800.00",
		"b": "100.50",
		"q": "5100000.25"
	}`

	var kline KLineEvent
	err := json.Unmarshal([]byte(jsonStr), &kline)
	if err != nil {
		t.Fatalf("Failed to unmarshal KLineEvent: %v", err)
	}

	if kline.OpenPrice != 50000.50 {
		t.Errorf("Expected OpenPrice 50000.50, got %f", kline.OpenPrice)
	}

	if kline.HighPrice != 51000.25 {
		t.Errorf("Expected HighPrice 51000.25, got %f", kline.HighPrice)
	}

	if kline.LowPrice != 49500.75 {
		t.Errorf("Expected LowPrice 49500.75, got %f", kline.LowPrice)
	}

	if kline.ClosePrice != 50800.00 {
		t.Errorf("Expected ClosePrice 50800.00, got %f", kline.ClosePrice)
	}

	if kline.BaseVolume != 100.50 {
		t.Errorf("Expected BaseVolume 100.50, got %f", kline.BaseVolume)
	}

	if kline.QuoteVolume != 5100000.25 {
		t.Errorf("Expected QuoteVolume 5100000.25, got %f", kline.QuoteVolume)
	}

	invalidJSON := `{
		"o": "not-a-number",
		"h": "51000.25",
		"l": "49500.75",
		"c": "50800.00",
		"b": "100.50",
		"q": "5100000.25"
	}`

	err = json.Unmarshal([]byte(invalidJSON), &kline)
	if err == nil {
		t.Fatal("Expected error for invalid number format, got none")
	}
}

func TestKLineChannelMessageUnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"ch": "kline",
		"symbol": "BTCUSDT",
		"ts": 1645123456000,
		"data": {
			"o": "50000.50",
			"h": "51000.25",
			"l": "49500.75",
			"c": "50800.00",
			"b": "100.50",
			"q": "5100000.25"
		}
	}`

	var message KLineChannelMessage
	err := json.Unmarshal([]byte(jsonStr), &message)
	if err != nil {
		t.Fatalf("Failed to unmarshal KLineChannelMessage: %v", err)
	}

	if message.Channel != "kline" {
		t.Errorf("Expected Channel 'kline', got %s", message.Channel)
	}

	if message.Symbol != "BTCUSDT" {
		t.Errorf("Expected Symbol 'BTCUSDT', got %s", message.Symbol)
	}

	if message.Ts != 1645123456000 {
		t.Errorf("Expected Ts 1645123456000, got %d", message.Ts)
	}

	if message.Data.OpenPrice != 50000.50 {
		t.Errorf("Expected OpenPrice 50000.50, got %f", message.Data.OpenPrice)
	}

	if message.Data.HighPrice != 51000.25 {
		t.Errorf("Expected HighPrice 51000.25, got %f", message.Data.HighPrice)
	}

	if message.Data.LowPrice != 49500.75 {
		t.Errorf("Expected LowPrice 49500.75, got %f", message.Data.LowPrice)
	}

	if message.Data.ClosePrice != 50800.00 {
		t.Errorf("Expected ClosePrice 50800.00, got %f", message.Data.ClosePrice)
	}

	invalidJSON := `{
		"ch": "kline",
		"symbol": 12345,
		"ts": 1645123456000,
		"data": {
			"o": "50000.50",
			"h": "51000.25",
			"l": "49500.75",
			"c": "50800.00",
			"b": "100.50",
			"q": "5100000.25"
		}
	}`

	err = json.Unmarshal([]byte(invalidJSON), &message)
	if err == nil {
		t.Fatal("Expected error for invalid symbol format, got none")
	}
}

func TestMarginCoinConsistencyWithSymbol(t *testing.T) {
	testCases := []struct {
		symbolInput     string
		marginCoinInput string
	}{
		{
			symbolInput:     "BTCUSDT",
			marginCoinInput: "USDT",
		},
		{
			symbolInput:     "ETHUSDT",
			marginCoinInput: "USDT",
		},
		{
			symbolInput:     "btcusdt",
			marginCoinInput: "usdt",
		},
		{
			symbolInput:     " BtCuSdT ",
			marginCoinInput: " uSdT ",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s-%s", tc.symbolInput, tc.marginCoinInput), func(t *testing.T) {
			parsedSymbol := ParseSymbol(tc.symbolInput)
			parsedCoin := ParseMarginCoin(tc.marginCoinInput)

			if !strings.HasSuffix(string(parsedSymbol), string(parsedCoin)) {
				t.Errorf("SubscribeSymbol '%s' should end with coin '%s' after normalization", parsedSymbol, parsedCoin)
			}
		})
	}
}
