package bitunix

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestTPSLOrderBuilderCreation(t *testing.T) {
	symbol := "BTCUSDT"
	positionID := "position123"

	builder := NewTPSLOrderBuilder(symbol, positionID)
	order := builder.Build()

	if order.Symbol != symbol {
		t.Errorf("Expected symbol %s, got %s", symbol, order.Symbol)
	}

	if order.PositionID != positionID {
		t.Errorf("Expected positionID %s, got %s", positionID, order.PositionID)
	}
}

func TestTPSLOrderBuilderMethods(t *testing.T) {
	symbol := "BTCUSDT"
	positionID := "position123"

	builder := NewTPSLOrderBuilder(symbol, positionID)

	tpPrice := 55000.0
	tpQty := 1.0
	tpStopType := StopTypeLastPrice
	tpOrderType := OrderTypeLimit
	tpOrderPrice := 54500.0

	builder.WithTakeProfit(tpPrice, tpQty, tpStopType, tpOrderType, tpOrderPrice)
	order := builder.Build()

	if *order.TpPrice != tpPrice {
		t.Errorf("WithTakeProfit: Expected TP price %f, got %f", tpPrice, *order.TpPrice)
	}

	if *order.TpQty != tpQty {
		t.Errorf("WithTakeProfit: Expected TP qty %f, got %f", tpQty, *order.TpQty)
	}

	if order.TpStopType != tpStopType {
		t.Errorf("WithTakeProfit: Expected TP stop type %s, got %s", tpStopType, order.TpStopType)
	}

	if order.TpOrderType != tpOrderType {
		t.Errorf("WithTakeProfit: Expected TP order type %s, got %s", tpOrderType, order.TpOrderType)
	}

	if *order.TpOrderPrice != tpOrderPrice {
		t.Errorf("WithTakeProfit: Expected TP order price %f, got %f", tpOrderPrice, *order.TpOrderPrice)
	}

	slPrice := 45000.0
	slQty := 1.0
	slStopType := StopTypeMarkPrice
	slOrderType := OrderTypeMarket
	slOrderPrice := 45500.0

	builder.WithStopLoss(slPrice, slQty, slStopType, slOrderType, slOrderPrice)
	order = builder.Build()

	if *order.SlPrice != slPrice {
		t.Errorf("WithStopLoss: Expected SL price %f, got %f", slPrice, *order.SlPrice)
	}

	if *order.SlQty != slQty {
		t.Errorf("WithStopLoss: Expected SL qty %f, got %f", slQty, *order.SlQty)
	}

	if order.SlStopType != slStopType {
		t.Errorf("WithStopLoss: Expected SL stop type %s, got %s", slStopType, order.SlStopType)
	}

	if order.SlOrderType != slOrderType {
		t.Errorf("WithStopLoss: Expected SL order type %s, got %s", slOrderType, order.SlOrderType)
	}

	if *order.SlOrderPrice != slOrderPrice {
		t.Errorf("WithStopLoss: Expected SL order price %f, got %f", slOrderPrice, *order.SlOrderPrice)
	}
}

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

func TestPlaceTpSlOrder(t *testing.T) {
	mockResponse := `{
		"code": 0,
		"msg": "success",
		"data": [
			{
				"orderId": "tp123",
				"clientId": "tpclient123"
			},
			{
				"orderId": "sl123",
				"clientId": "slclient123"
			}
		]
	}`

	mockAPI := NewMockAPI(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))

		if r.URL.Path != "/api/v1/futures/tpsl/place_order" {
			t.Errorf("Expected request path /api/v1/futures/tpsl/place_order, got %s", r.URL.Path)
		}

		if r.Method != http.MethodPost {
			t.Errorf("Expected request method POST, got %s", r.Method)
		}
	})
	defer mockAPI.Close()

	tpPrice := 55000.0
	tpQty := 1.0
	tpOrderPrice := 54500.0
	slPrice := 45000.0
	slQty := 1.0
	slOrderPrice := 45500.0

	orderReq := &TPSLOrderRequest{
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

	response, err := mockAPI.client.PlaceTpSlOrder(context.Background(), orderReq)
	if err != nil {
		t.Fatalf("PlaceTpSlOrder returned error: %v", err)
	}

	if response.Code != 0 {
		t.Errorf("Expected code 0, got %d", response.Code)
	}

	if response.Message != "success" {
		t.Errorf("Expected message 'success', got %s", response.Message)
	}

	if len(response.Data) != 2 {
		t.Fatalf("Expected 2 orders in response, got %d", len(response.Data))
	}

	tpOrder := response.Data[0]
	if tpOrder.OrderID != "tp123" {
		t.Errorf("Expected TP order ID 'tp123', got %s", tpOrder.OrderID)
	}
	if tpOrder.ClientId != "tpclient123" {
		t.Errorf("Expected TP client ID 'tpclient123', got %s", tpOrder.ClientId)
	}

	slOrder := response.Data[1]
	if slOrder.OrderID != "sl123" {
		t.Errorf("Expected SL order ID 'sl123', got %s", slOrder.OrderID)
	}
	if slOrder.ClientId != "slclient123" {
		t.Errorf("Expected SL client ID 'slclient123', got %s", slOrder.ClientId)
	}
}
