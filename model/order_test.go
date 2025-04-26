package model

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTPSLOrderRequestMarshalJSON(t *testing.T) {
	tpPrice := 55000.0
	tpQty := 1.0
	tpOrderPrice := 54500.0
	slPrice := 45000.0
	slQty := 1.0
	slOrderPrice := 45500.0

	req := TPSLOrderRequest{
		Symbol:       "BTCUSDT",
		PositionID:   "position123",
		TpPrice:      &tpPrice,
		TpQty:        &tpQty,
		TpStopType:   StopTypeLastPrice,
		TpOrderType:  OrderTypeLimit,
		TpOrderPrice: &tpOrderPrice,
		SlPrice:      &slPrice,
		SlQty:        &slQty,
		SlStopType:   StopTypeMarkPrice,
		SlOrderType:  OrderTypeMarket,
		SlOrderPrice: &slOrderPrice,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal TPSL order request: %v", err)
	}

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if m["tpPrice"] != "55000" {
		t.Errorf("Expected tpPrice to be marshaled as string \"55000\", got %v", m["tpPrice"])
	}

	if m["tpQty"] != "1" {
		t.Errorf("Expected tpQty to be marshaled as string \"1\", got %v", m["tpQty"])
	}

	if m["tpOrderPrice"] != "54500" {
		t.Errorf("Expected tpOrderPrice to be marshaled as string \"54500\", got %v", m["tpOrderPrice"])
	}

	if m["slPrice"] != "45000" {
		t.Errorf("Expected slPrice to be marshaled as string \"45000\", got %v", m["slPrice"])
	}

	if m["slQty"] != "1" {
		t.Errorf("Expected slQty to be marshaled as string \"1\", got %v", m["slQty"])
	}

	if m["slOrderPrice"] != "45500" {
		t.Errorf("Expected slOrderPrice to be marshaled as string \"45500\", got %v", m["slOrderPrice"])
	}
}

func TestOrderRequestMarshalJSON(t *testing.T) {
	qty := 1.0
	price := 50000.0
	tpPrice := 55000.0
	tpOrderPrice := 54500.0
	slPrice := 45000.0
	slOrderPrice := 45500.0

	req := &OrderRequest{
		Symbol:       "BTCUSDT",
		TradeAction:  TradeSideBuy,
		Price:        &price,
		Qty:          &qty,
		PositionID:   "position123",
		TradeSide:    SideOpen,
		OrderType:    OrderTypeLimit,
		ReduceOnly:   false,
		Effect:       TimeInForceGTC,
		ClientID:     "client123",
		TpPrice:      &tpPrice,
		TpStopType:   StopTypeLastPrice,
		TpOrderType:  OrderTypeLimit,
		TpOrderPrice: &tpOrderPrice,
		SlPrice:      &slPrice,
		SlStopType:   StopTypeMarkPrice,
		SlOrderType:  OrderTypeMarket,
		SlOrderPrice: &slOrderPrice,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal order request: %v", err)
	}

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if m["price"] != "50000" {
		t.Errorf("Expected price to be marshaled as string \"50000\", got %v", m["price"])
	}
	if m["qty"] != "1" {
		t.Errorf("Expected qty to be marshaled as string \"1\", got %v", m["qty"])
	}
	if m["tpPrice"] != "55000" {
		t.Errorf("Expected tpPrice to be marshaled as string \"55000\", got %v", m["tpPrice"])
	}
	if m["tpOrderPrice"] != "54500" {
		t.Errorf("Expected tpOrderPrice to be marshaled as string \"54500\", got %v", m["tpOrderPrice"])
	}
	if m["slPrice"] != "45000" {
		t.Errorf("Expected slPrice to be marshaled as string \"45000\", got %v", m["slPrice"])
	}
	if m["slOrderPrice"] != "45500" {
		t.Errorf("Expected slOrderPrice to be marshaled as string \"45500\", got %v", m["slOrderPrice"])
	}
}

func TestOrderHistoryParamsMarshaling(t *testing.T) {

	jsonStr := `{
		"orderId": "12345",
		"symbol": "BTCUSDT",
		"qty": "1.5",
		"tradeQty": "1.0",
		"positionMode": "HEDGE",
		"marginMode": "CROSS",
		"leverage": 10,
		"price": "MARKET",
		"side": "BUY",
		"orderType": "LIMIT",
		"effect": "GTC",
		"clientId": "client_123",
		"reduceOnly": false,
		"status": "FILLED",
		"fee": "0.6",
		"realizedPNL": "100",
		"tpPrice": "45000",
		"tpStopType": "MARK_PRICE",
		"tpOrderType": "LIMIT",
		"tpOrderPrice": "45000",
		"slPrice": "35000",
		"slStopType": "MARK_PRICE",
		"slOrderType": "MARKET",
		"slOrderPrice": "0",
		"ctime": "1659076670000",
		"mtime": "1659076680000"
	}`

	order := HistoricalOrder{}
	err := json.Unmarshal([]byte(jsonStr), &order)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling order: %v", err)
	}

	if order.OrderID != "12345" {
		t.Errorf("unexpected orderId: %s", order.OrderID)
	}

	if order.Quantity != 1.5 {
		t.Errorf("unexpected quantity: %f", order.Quantity)
	}

	if order.TradeQuantity != 1.0 {
		t.Errorf("unexpected trade quantity: %f", order.TradeQuantity)
	}

	if order.Price != "MARKET" {
		t.Errorf("unexpected price: %s", order.Price)
	}

	if order.Fee != 0.6 {
		t.Errorf("unexpected fee: %f", order.Fee)
	}

	if order.RealizedPNL != 100 {
		t.Errorf("unexpected realized PNL: %f", order.RealizedPNL)
	}

	expectedCreateTime := time.Unix(0, 1659076670000*1000000)
	if !order.CreateTime.Equal(expectedCreateTime) {
		t.Errorf("unexpected create time: %v, expected: %v", order.CreateTime, expectedCreateTime)
	}

	expectedModifyTime := time.Unix(0, 1659076680000*1000000)
	if !order.ModifyTime.Equal(expectedModifyTime) {
		t.Errorf("unexpected modify time: %v, expected: %v", order.ModifyTime, expectedModifyTime)
	}
}
