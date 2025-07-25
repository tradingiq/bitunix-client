package bitunix

import (
	"context"
	"net/http"
	"testing"

	"github.com/tradingiq/bitunix-client/model"
)

func TestGetPendingTPSLOrder(t *testing.T) {
	mockResponse := `{
		"code": 0,
		"msg": "Success",
		"data": [
			{
				"id": "123",
				"positionId": "12345678",
				"symbol": "BTCUSDT",
				"base": "BTC",
				"quote": "USDT",
				"tpPrice": "50000",
				"tpStopType": "LAST_PRICE",
				"slPrice": "70000",
				"slStopType": "LAST_PRICE",
				"tpOrderType": "LIMIT",
				"tpOrderPrice": "50000",
				"slOrderType": "LIMIT",
				"slOrderPrice": "70000",
				"tpQty": "0.01",
				"slQty": "0.01"
			}
		]
	}`

	mockAPI := NewMockAPI(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/futures/tpsl/get_pending_orders" {
			t.Errorf("Expected path '/api/v1/futures/tpsl/get_pending_orders', got %s", r.URL.Path)
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

	params := model.PendingTPSLOrderParams{
		Symbol: model.ParseSymbol("BTCUSDT"),
		Limit:  10,
	}

	response, err := mockAPI.client.GetPendingTPSLOrder(context.Background(), params)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if response.Code != 0 {
		t.Errorf("Expected code 0, got %d", response.Code)
	}

	if len(response.Data) != 1 {
		t.Fatalf("Expected 1 order, got %d", len(response.Data))
	}

	order := response.Data[0]
	if order.ID != "123" {
		t.Errorf("Expected ID '123', got %s", order.ID)
	}
	if order.PositionID != "12345678" {
		t.Errorf("Expected PositionID '12345678', got %s", order.PositionID)
	}
	if order.Symbol != model.ParseSymbol("BTCUSDT") {
		t.Errorf("Expected Symbol 'BTCUSDT', got %s", order.Symbol)
	}
	if order.TpPrice == nil || *order.TpPrice != 50000 {
		t.Errorf("Expected TpPrice 50000, got %v", order.TpPrice)
	}
	if order.SlPrice == nil || *order.SlPrice != 70000 {
		t.Errorf("Expected SlPrice 70000, got %v", order.SlPrice)
	}
	if order.TpStopType == nil || *order.TpStopType != model.StopTypeLastPrice {
		t.Errorf("Expected TpStopType LAST_PRICE, got %v", order.TpStopType)
	}
	if order.SlStopType == nil || *order.SlStopType != model.StopTypeLastPrice {
		t.Errorf("Expected SlStopType LAST_PRICE, got %v", order.SlStopType)
	}
	if order.TpOrderType == nil || *order.TpOrderType != model.OrderTypeLimit {
		t.Errorf("Expected TpOrderType LIMIT, got %v", order.TpOrderType)
	}
	if order.SlOrderType == nil || *order.SlOrderType != model.OrderTypeLimit {
		t.Errorf("Expected SlOrderType LIMIT, got %v", order.SlOrderType)
	}
	if order.TpQty == nil || *order.TpQty != 0.01 {
		t.Errorf("Expected TpQty 0.01, got %v", order.TpQty)
	}
	if order.SlQty == nil || *order.SlQty != 0.01 {
		t.Errorf("Expected SlQty 0.01, got %v", order.SlQty)
	}
}

func TestGetPendingTPSLOrder_Error(t *testing.T) {
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
				"code": 40001,
				"msg": "Invalid parameter"
			}`,
			responseStatus: http.StatusOK,
			expectError:    true,
			errorCode:      40001,
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

			params := model.PendingTPSLOrderParams{
				Symbol: model.ParseSymbol("BTCUSDT"),
			}

			_, err := mockAPI.client.GetPendingTPSLOrder(context.Background(), params)
			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
