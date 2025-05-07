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
		TradeSide:    TradeSideBuy,
		Price:        &price,
		Qty:          qty,
		PositionID:   "position123",
		Side:         SideOpen,
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

func TestOrderResponseDataParsing(t *testing.T) {
	jsonStr := `{
		"orderId": "12345",
		"clientId": "client_123"
	}`

	var data OrderResponseData
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal OrderResponseData: %v", err)
	}

	if data.OrderId != "12345" {
		t.Errorf("Expected OrderId '12345', got %s", data.OrderId)
	}

	if data.ClientId != "client_123" {
		t.Errorf("Expected ClientId 'client_123', got %s", data.ClientId)
	}
}

func TestOrderResponseParsing(t *testing.T) {
	jsonStr := `{
		"code": 0,
		"message": "success",
		"data": {
			"orderId": "12345",
			"clientId": "client_123"
		}
	}`

	var response OrderResponse
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal OrderResponse: %v", err)
	}

	if response.Code != 0 {
		t.Errorf("Expected Code 0, got %d", response.Code)
	}

	if response.Message != "success" {
		t.Errorf("Expected Message 'success', got %s", response.Message)
	}

	if response.Data.OrderId != "12345" {
		t.Errorf("Expected Data.OrderId '12345', got %s", response.Data.OrderId)
	}

	if response.Data.ClientId != "client_123" {
		t.Errorf("Expected Data.ClientId 'client_123', got %s", response.Data.ClientId)
	}
}

func TestCancelOrderResponseParsing(t *testing.T) {
	jsonStr := `{
		"code": 0,
		"msg": "success",
		"data": {
			"successList": [
				{
					"orderId": "order123",
					"clientId": "client123"
				},
				{
					"orderId": "order456",
					"clientId": "client456"
				}
			],
			"failureList": [
				{
					"orderId": "order789",
					"clientId": "client789",
					"errorMsg": "Order not found",
					"errorCode": "40004"
				}
			]
		}
	}`

	var response CancelOrderResponse
	err := json.Unmarshal([]byte(jsonStr), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal CancelOrderResponse: %v", err)
	}

	if response.Code != 0 {
		t.Errorf("Expected Code 0, got %d", response.Code)
	}

	if response.Message != "success" {
		t.Errorf("Expected Message 'success', got %s", response.Message)
	}

	if len(response.Data.SuccessList) != 2 {
		t.Errorf("Expected 2 items in successList, got %d", len(response.Data.SuccessList))
	} else {
		if response.Data.SuccessList[0].OrderId != "order123" {
			t.Errorf("Expected first success OrderId 'order123', got %s", response.Data.SuccessList[0].OrderId)
		}
		if response.Data.SuccessList[0].ClientId != "client123" {
			t.Errorf("Expected first success ClientId 'client123', got %s", response.Data.SuccessList[0].ClientId)
		}

		if response.Data.SuccessList[1].OrderId != "order456" {
			t.Errorf("Expected second success OrderId 'order456', got %s", response.Data.SuccessList[1].OrderId)
		}
		if response.Data.SuccessList[1].ClientId != "client456" {
			t.Errorf("Expected second success ClientId 'client456', got %s", response.Data.SuccessList[1].ClientId)
		}
	}

	if len(response.Data.FailureList) != 1 {
		t.Errorf("Expected 1 item in failureList, got %d", len(response.Data.FailureList))
	} else {
		if response.Data.FailureList[0].OrderId != "order789" {
			t.Errorf("Expected failure OrderId 'order789', got %s", response.Data.FailureList[0].OrderId)
		}
		if response.Data.FailureList[0].ClientId != "client789" {
			t.Errorf("Expected failure ClientId 'client789', got %s", response.Data.FailureList[0].ClientId)
		}
		if response.Data.FailureList[0].ErrorMsg != "Order not found" {
			t.Errorf("Expected failure ErrorMsg 'Order not found', got %s", response.Data.FailureList[0].ErrorMsg)
		}
		if response.Data.FailureList[0].ErrorCode != "40004" {
			t.Errorf("Expected failure ErrorCode '40004', got %s", response.Data.FailureList[0].ErrorCode)
		}
	}
}

func TestCancelOrderRequestMarshaling(t *testing.T) {
	req := CancelOrderRequest{
		Symbol: "BTCUSDT",
		OrderList: []CancelOrderParam{
			{
				OrderID:  "order123",
				ClientID: "client123",
			},
			{
				OrderID:  "order456",
				ClientID: "",
			},
			{
				OrderID:  "",
				ClientID: "client789",
			},
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal CancelOrderRequest: %v", err)
	}

	var m map[string]interface{}
	err = json.Unmarshal(data, &m)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if m["symbol"] != "BTCUSDT" {
		t.Errorf("Expected symbol to be 'BTCUSDT', got %v", m["symbol"])
	}

	orderList, ok := m["orderList"].([]interface{})
	if !ok {
		t.Fatalf("Expected orderList to be an array, got %T", m["orderList"])
	}

	if len(orderList) != 3 {
		t.Errorf("Expected 3 items in orderList, got %d", len(orderList))
	} else {
		item1, ok := orderList[0].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected orderList[0] to be an object, got %T", orderList[0])
		}
		if item1["orderId"] != "order123" {
			t.Errorf("Expected orderId to be 'order123', got %v", item1["orderId"])
		}
		if item1["clientId"] != "client123" {
			t.Errorf("Expected clientId to be 'client123', got %v", item1["clientId"])
		}

		item2, ok := orderList[1].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected orderList[1] to be an object, got %T", orderList[1])
		}
		if item2["orderId"] != "order456" {
			t.Errorf("Expected orderId to be 'order456', got %v", item2["orderId"])
		}

		clientID2, exists := item2["clientId"]
		if exists {

			if clientID2 != "" {
				t.Errorf("Expected clientId to be empty string, got %v", clientID2)
			}
		}

		item3, ok := orderList[2].(map[string]interface{})
		if !ok {
			t.Fatalf("Expected orderList[2] to be an object, got %T", orderList[2])
		}

		orderID3, exists := item3["orderId"]
		if exists {

			if orderID3 != "" {
				t.Errorf("Expected orderId to be empty string, got %v", orderID3)
			}
		}
		if item3["clientId"] != "client789" {
			t.Errorf("Expected clientId to be 'client789', got %v", item3["clientId"])
		}
	}
}
