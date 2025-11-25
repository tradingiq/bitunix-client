package bitunix

import (
	"context"
	"net/http"
	"testing"

	"github.com/tradingiq/bitunix-client/model"
)

func TestGetPendingOrder(t *testing.T) {
	mockResponse := `{
		"code": 0,
		"msg": "Success",
		"data": {
			"orderList": [
				{
					"orderId": "11111",
					"qty": "1",
					"tradeQty": "0.5",
					"price": "60000",
					"symbol": "BTCUSDT",
					"positionMode": "HEDGE",
					"marginMode": "ISOLATION",
					"leverage": 15,
					"status": "NEW_",
					"fee": "0.01",
					"realizedPNL": "1.78",
					"orderType": "LIMIT",
					"effect": "GTC",
					"reduceOnly": false,
					"clientId": "22222",
					"side": "BUY",
					"tpPrice": "61000",
					"tpStopType": "MARK_PRICE",
					"tpOrderType": "LIMIT",
					"tpOrderPrice": "61000.1",
					"slPrice": "59000",
					"slStopType": "MARK_PRICE",
					"slOrderType": "LIMIT",
					"slOrderPrice": "59000.1",
					"ctime": "1597026383085",
					"mtime": "1597026383085"
				}
			],
			"total": "10"
		}
	}`

	mockAPI := NewMockAPI(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/futures/trade/get_pending_orders" {
			t.Errorf("Expected path '/api/v1/futures/trade/get_pending_orders', got %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("Expected method GET, got %s", r.Method)
		}

		query := r.URL.Query()
		if symbol := query.Get("symbol"); symbol != "BTCUSDT" {
			t.Errorf("Expected symbol 'BTCUSDT', got %s", symbol)
		}
		if limit := query.Get("limit"); limit != "10" {
			t.Errorf("Expected limit '10', got %s", limit)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	})
	defer mockAPI.Close()

	params := model.PendingOrderParams{
		Symbol: model.ParseSymbol("BTCUSDT"),
		Limit:  10,
	}

	response, err := mockAPI.client.GetPendingOrder(context.Background(), params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if response.Code != 0 {
		t.Errorf("Expected code 0, got %d", response.Code)
	}

	if len(response.Data.OrderList) != 1 {
		t.Fatalf("Expected 1 order, got %d", len(response.Data.OrderList))
	}

	if response.Data.Total != 10 {
		t.Errorf("Expected total 10, got %d", response.Data.Total)
	}

	order := response.Data.OrderList[0]
	if order.OrderID != "11111" {
		t.Errorf("Expected OrderID '11111', got %s", order.OrderID)
	}
	if order.ClientID != "22222" {
		t.Errorf("Expected ClientID '22222', got %s", order.ClientID)
	}
	if order.Symbol != model.ParseSymbol("BTCUSDT") {
		t.Errorf("Expected Symbol 'BTCUSDT', got %s", order.Symbol)
	}
	if order.Price != 60000 {
		t.Errorf("Expected Price 60000, got %f", order.Price)
	}
	if order.Quantity != 1 {
		t.Errorf("Expected Quantity 1, got %f", order.Quantity)
	}
	if order.TradeQuantity != 0.5 {
		t.Errorf("Expected TradeQuantity 0.5, got %f", order.TradeQuantity)
	}
	if order.PositionMode != model.PositionModeHedge {
		t.Errorf("Expected PositionMode HEDGE, got %s", order.PositionMode)
	}
	if order.MarginMode != model.MarginModeIsolation {
		t.Errorf("Expected MarginMode ISOLATION, got %s", order.MarginMode)
	}
	if order.Leverage != 15 {
		t.Errorf("Expected Leverage 15, got %d", order.Leverage)
	}
	if order.Status != model.OrderStatusNew {
		t.Errorf("Expected Status NEW, got %s", order.Status)
	}
	if order.Fee != 0.01 {
		t.Errorf("Expected Fee 0.01, got %f", order.Fee)
	}
	if order.RealizedPNL != 1.78 {
		t.Errorf("Expected RealizedPNL 1.78, got %f", order.RealizedPNL)
	}
	if order.OrderType != model.OrderTypeLimit {
		t.Errorf("Expected OrderType LIMIT, got %s", order.OrderType)
	}
	if order.Effect != model.TimeInForceGTC {
		t.Errorf("Expected Effect GTC, got %s", order.Effect)
	}
	if order.ReduceOnly != false {
		t.Errorf("Expected ReduceOnly false, got %t", order.ReduceOnly)
	}
	if order.Side != model.TradeSideBuy {
		t.Errorf("Expected Side BUY, got %s", order.Side)
	}
	if order.TpPrice == nil || *order.TpPrice != 61000 {
		t.Errorf("Expected TpPrice 61000, got %v", order.TpPrice)
	}
	if order.SlPrice == nil || *order.SlPrice != 59000 {
		t.Errorf("Expected SlPrice 59000, got %v", order.SlPrice)
	}
	if order.TpStopType == nil || *order.TpStopType != model.StopTypeMarkPrice {
		t.Errorf("Expected TpStopType MARK_PRICE, got %v", order.TpStopType)
	}
	if order.SlStopType == nil || *order.SlStopType != model.StopTypeMarkPrice {
		t.Errorf("Expected SlStopType MARK_PRICE, got %v", order.SlStopType)
	}
	if order.TpOrderType == nil || *order.TpOrderType != model.OrderTypeLimit {
		t.Errorf("Expected TpOrderType LIMIT, got %v", order.TpOrderType)
	}
	if order.SlOrderType == nil || *order.SlOrderType != model.OrderTypeLimit {
		t.Errorf("Expected SlOrderType LIMIT, got %v", order.SlOrderType)
	}
	if order.TpOrderPrice == nil || *order.TpOrderPrice != 61000.1 {
		t.Errorf("Expected TpOrderPrice 61000.1, got %v", order.TpOrderPrice)
	}
	if order.SlOrderPrice == nil || *order.SlOrderPrice != 59000.1 {
		t.Errorf("Expected SlOrderPrice 59000.1, got %v", order.SlOrderPrice)
	}
}

func TestGetPendingOrder_Error(t *testing.T) {
	testCases := []struct {
		name           string
		mockResponse   string
		responseStatus int
		expectError    bool
		errorCode      int
	}{
		{
			name: "API Error",
			mockResponse: `{
				"code": 10002,
				"msg": "Invalid parameter"
			}`,
			responseStatus: http.StatusOK,
			expectError:    true,
			errorCode:      10002,
		},
		{
			name:           "HTTP Error",
			mockResponse:   "Internal Server Error",
			responseStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAPI := NewMockAPI(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.responseStatus)
				w.Write([]byte(tc.mockResponse))
			})
			defer mockAPI.Close()

			params := model.PendingOrderParams{
				Symbol: model.ParseSymbol("BTCUSDT"),
			}

			_, err := mockAPI.client.GetPendingOrder(context.Background(), params)
			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}