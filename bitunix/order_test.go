package bitunix

import (
	"context"
	"github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockAPI struct {
	server *httptest.Server
	client ApiClient
}

func NewMockAPI(handler http.HandlerFunc) *MockAPI {
	server := httptest.NewServer(handler)
	client, _ := NewApiClient("", "", WithBaseURI(server.URL))
	return &MockAPI{
		server: server,
		client: client,
	}
}

func (m *MockAPI) Close() {
	m.server.Close()
}

func TestOrderBuilderCreation(t *testing.T) {
	symbol := model.ParseSymbol("BTCUSDT")
	side := model.TradeSideBuy
	tradeSide := model.SideOpen
	qty := 1.0

	builder := NewOrderBuilder(symbol, side, tradeSide, qty)
	order, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build order: %v", err)
	}

	if order.Symbol != symbol {
		t.Errorf("Expected symbol %s, got %s", symbol, order.Symbol)
	}
	if order.TradeSide != side {
		t.Errorf("Expected side %s, got %s", side, order.TradeSide)
	}
	if order.Side != tradeSide {
		t.Errorf("Expected tradeSide %s, got %s", tradeSide, order.Side)
	}
	if order.Qty != qty {
		t.Errorf("Expected qty %f, got %f", qty, order.Qty)
	}
	if order.OrderType != model.OrderTypeMarket {
		t.Errorf("Expected order type %s, got %s", model.OrderTypeMarket, order.OrderType)
	}
	if order.ReduceOnly != false {
		t.Errorf("Expected reduceOnly %v, got %v", false, order.ReduceOnly)
	}
	if order.Effect != model.TimeInForceGTC {
		t.Errorf("Expected effect %s, got %s", model.TimeInForceGTC, order.Effect)
	}
}

func TestOrderBuilderMethods(t *testing.T) {
	symbol := model.ParseSymbol("BTCUSDT")
	side := model.TradeSideBuy
	tradeSide := model.SideOpen
	qty := 1.0

	builder := NewOrderBuilder(symbol, side, tradeSide, qty)

	orderType := model.OrderTypeLimit
	builder.WithOrderType(orderType)
	builder.WithPrice(10)
	order, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build order: %v", err)
	}
	if order.OrderType != orderType {
		t.Errorf("WithOrderType: Expected order type %s, got %s", orderType, order.OrderType)
	}

	price := 50000.0
	builder.WithPrice(price)
	order, err = builder.Build()
	if err != nil {
		t.Fatalf("Failed to build order: %v", err)
	}
	if *order.Price != price {
		t.Errorf("WithPrice: Expected price %f, got %f", price, *order.Price)
	}

	positionID := "position123"
	builder.WithPositionID(positionID)
	order, err = builder.Build()
	if err != nil {
		t.Fatalf("Failed to build order: %v", err)
	}
	if order.PositionID != positionID {
		t.Errorf("WithPositionID: Expected position ID %s, got %s", positionID, order.PositionID)
	}

	reduceOnly := true
	builder.WithReduceOnly(reduceOnly)
	order, err = builder.Build()
	if err != nil {
		t.Fatalf("Failed to build order: %v", err)
	}
	if order.ReduceOnly != reduceOnly {
		t.Errorf("WithReduceOnly: Expected reduce only %v, got %v", reduceOnly, order.ReduceOnly)
	}

	tif := model.TimeInForceIOC
	builder.WithTimeInForce(tif)
	order, err = builder.Build()
	if err != nil {
		t.Fatalf("Failed to build order: %v", err)
	}
	if order.Effect != tif {
		t.Errorf("WithTimeInForce: Expected time in force %s, got %s", tif, order.Effect)
	}

	clientID := "client123"
	builder.WithClientID(clientID)
	order, err = builder.Build()
	if err != nil {
		t.Fatalf("Failed to build order: %v", err)
	}
	if order.ClientID != clientID {
		t.Errorf("WithClientID: Expected apiClient ID %s, got %s", clientID, order.ClientID)
	}

	tpPrice := 55000.0
	tpStopType := model.StopTypeLastPrice
	tpOrderType := model.OrderTypeLimit
	builder.WithTakeProfit(tpPrice, tpStopType, tpOrderType)
	builder.WithTakeProfitPrice(10)
	order, err = builder.Build()
	if err != nil {
		t.Fatalf("Failed to build order: %v", err)
	}
	if *order.TpPrice != tpPrice {
		t.Errorf("WithTakeProfit: Expected TP price %f, got %f", tpPrice, *order.TpPrice)
	}
	if order.TpStopType != tpStopType {
		t.Errorf("WithTakeProfit: Expected TP stop type %s, got %s", tpStopType, order.TpStopType)
	}
	if order.TpOrderType != tpOrderType {
		t.Errorf("WithTakeProfit: Expected TP order type %s, got %s", tpOrderType, order.TpOrderType)
	}

	tpOrderPrice := 54500.0
	builder.WithTakeProfitPrice(tpOrderPrice)
	order, err = builder.Build()
	if err != nil {
		t.Fatalf("Failed to build order: %v", err)
	}
	if *order.TpOrderPrice != tpOrderPrice {
		t.Errorf("WithTakeProfitPrice: Expected TP order price %f, got %f", tpOrderPrice, *order.TpOrderPrice)
	}

	slPrice := 45000.0
	slStopType := model.StopTypeMarkPrice
	slOrderType := model.OrderTypeMarket
	builder.WithStopLoss(slPrice, slStopType, slOrderType)
	order, err = builder.Build()
	if err != nil {
		t.Fatalf("Failed to build order: %v", err)
	}
	if *order.SlPrice != slPrice {
		t.Errorf("WithStopLoss: Expected SL price %f, got %f", slPrice, *order.SlPrice)
	}
	if order.SlStopType != slStopType {
		t.Errorf("WithStopLoss: Expected SL stop type %s, got %s", slStopType, order.SlStopType)
	}
	if order.SlOrderType != slOrderType {
		t.Errorf("WithStopLoss: Expected SL order type %s, got %s", slOrderType, order.SlOrderType)
	}

	slOrderPrice := 45500.0
	builder.WithStopLossPrice(slOrderPrice)
	order, err = builder.Build()
	if err != nil {
		t.Fatalf("Failed to build order: %v", err)
	}
	if *order.SlOrderPrice != slOrderPrice {
		t.Errorf("WithStopLossPrice: Expected SL order price %f, got %f", slOrderPrice, *order.SlOrderPrice)
	}
}

func TestOrderBuilder_ValidationErrors(t *testing.T) {
	testCases := []struct {
		name          string
		setupBuilder  func() *OrderBuilder
		expectedError string
	}{
		{
			name: "Missing Symbol",
			setupBuilder: func() *OrderBuilder {
				builder := NewOrderBuilder("", model.TradeSideBuy, model.SideOpen, 1.0)
				return builder
			},
			expectedError: "validation error: field symbol: is required",
		},
		{
			name: "Missing TradeSide",
			setupBuilder: func() *OrderBuilder {
				builder := NewOrderBuilder("BTCUSDT", "", model.SideOpen, 1.0)
				return builder
			},
			expectedError: "validation error: field side: trade action is required",
		},
		{
			name: "Missing Side",
			setupBuilder: func() *OrderBuilder {
				builder := NewOrderBuilder("BTCUSDT", model.TradeSideBuy, "", 1.0)
				return builder
			},
			expectedError: "validation error: field tradeSide: is required",
		},
		{
			name: "Zero Quantity",
			setupBuilder: func() *OrderBuilder {
				builder := NewOrderBuilder("BTCUSDT", model.TradeSideBuy, model.SideOpen, 0)
				return builder
			},
			expectedError: "validation error: field qty: quantity is required and must be greater than zero",
		},
		{
			name: "Negative Quantity",
			setupBuilder: func() *OrderBuilder {
				builder := NewOrderBuilder("BTCUSDT", model.TradeSideBuy, model.SideOpen, -1.0)
				return builder
			},
			expectedError: "validation error: field qty: quantity is required and must be greater than zero",
		},
		{
			name: "Missing PositionID for CLOSE Side",
			setupBuilder: func() *OrderBuilder {
				builder := NewOrderBuilder("BTCUSDT", model.TradeSideBuy, model.SideClose, 1.0)
				return builder
			},
			expectedError: "validation error: field positionId: is required when tradeSide is CLOSE",
		},
		{
			name: "Missing Price for Limit Order",
			setupBuilder: func() *OrderBuilder {
				builder := NewOrderBuilder("BTCUSDT", model.TradeSideBuy, model.SideOpen, 1.0)
				builder.WithOrderType(model.OrderTypeLimit)
				return builder
			},
			expectedError: "validation error: field price: is required for limit orders",
		},
		{
			name: "Zero Price for Limit Order",
			setupBuilder: func() *OrderBuilder {
				builder := NewOrderBuilder("BTCUSDT", model.TradeSideBuy, model.SideOpen, 1.0)
				builder.WithOrderType(model.OrderTypeLimit)
				price := 0.0
				builder.WithPrice(price)
				return builder
			},
			expectedError: "validation error: field price: is required for limit orders",
		},
		{
			name: "Missing TP StopType with TP Price",
			setupBuilder: func() *OrderBuilder {
				builder := NewOrderBuilder("BTCUSDT", model.TradeSideBuy, model.SideOpen, 1.0)
				tpPrice := 50000.0
				builder.request.TpPrice = &tpPrice
				return builder
			},
			expectedError: "validation error: field tpStopType: is required when setting take profit",
		},
		{
			name: "Missing TP OrderType with TP Price",
			setupBuilder: func() *OrderBuilder {
				builder := NewOrderBuilder("BTCUSDT", model.TradeSideBuy, model.SideOpen, 1.0)
				tpPrice := 50000.0
				builder.request.TpPrice = &tpPrice
				builder.request.TpStopType = model.StopTypeLastPrice
				return builder
			},
			expectedError: "validation error: field tpOrderType: is required when setting take profit",
		},
		{
			name: "Missing TP OrderPrice with Limit TP",
			setupBuilder: func() *OrderBuilder {
				builder := NewOrderBuilder("BTCUSDT", model.TradeSideBuy, model.SideOpen, 1.0)
				tpPrice := 50000.0
				builder.request.TpPrice = &tpPrice
				builder.request.TpStopType = model.StopTypeLastPrice
				builder.request.TpOrderType = model.OrderTypeLimit
				return builder
			},
			expectedError: "validation error: field tpOrderPrice: is required when tpOrderType is LIMIT",
		},
		{
			name: "Missing SL StopType with SL Price",
			setupBuilder: func() *OrderBuilder {
				builder := NewOrderBuilder("BTCUSDT", model.TradeSideBuy, model.SideOpen, 1.0)
				slPrice := 40000.0
				builder.request.SlPrice = &slPrice
				return builder
			},
			expectedError: "validation error: field slStopType: is required when setting stop loss",
		},
		{
			name: "Missing SL OrderType with SL Price",
			setupBuilder: func() *OrderBuilder {
				builder := NewOrderBuilder("BTCUSDT", model.TradeSideBuy, model.SideOpen, 1.0)
				slPrice := 40000.0
				builder.request.SlPrice = &slPrice
				builder.request.SlStopType = model.StopTypeMarkPrice
				return builder
			},
			expectedError: "validation error: field slOrderType: is required when setting stop loss",
		},
		{
			name: "Missing SL OrderPrice with Limit SL",
			setupBuilder: func() *OrderBuilder {
				builder := NewOrderBuilder("BTCUSDT", model.TradeSideBuy, model.SideOpen, 1.0)
				slPrice := 40000.0
				builder.request.SlPrice = &slPrice
				builder.request.SlStopType = model.StopTypeMarkPrice
				builder.request.SlOrderType = model.OrderTypeLimit
				return builder
			},
			expectedError: "validation error: field slOrderPrice: is required when slOrderType is LIMIT",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := tc.setupBuilder()
			_, err := builder.Build()

			if err == nil {
				t.Error("Expected validation error, got nil")
				return
			}

			if err.Error() != tc.expectedError {
				t.Errorf("Expected error message '%s', got '%s'", tc.expectedError, err.Error())
			}
		})
	}
}

func TestPlaceOrder(t *testing.T) {
	mockResponse := `{
		"code": 0,
		"message": "success",
		"data": {
			"orderId": "12345",
			"clientId": "client123"
		}
	}`

	mockAPI := NewMockAPI(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))

		if r.URL.Path != "/api/v1/futures/trade/place_order" {
			t.Errorf("Expected request path /api/v1/futures/trade/place_order, got %s", r.URL.Path)
		}

		if r.Method != http.MethodPost {
			t.Errorf("Expected request method POST, got %s", r.Method)
		}
	})
	defer mockAPI.Close()

	qty := 1.0
	price := 50000.0

	orderReq := &model.OrderRequest{
		Symbol:    "BTCUSDT",
		TradeSide: model.TradeSideBuy,
		Price:     &price,
		Qty:       qty,
		Side:      model.SideOpen,
		OrderType: model.OrderTypeLimit,
		ClientID:  "client123",
	}

	response, err := mockAPI.client.PlaceOrder(context.Background(), orderReq)
	if err != nil {
		t.Fatalf("PlaceOrder returned error: %v", err)
	}

	if response.Code != 0 {
		t.Errorf("Expected code 0, got %d", response.Code)
	}
	if response.Message != "success" {
		t.Errorf("Expected message 'success', got %s", response.Message)
	}
	if response.Data.OrderId != "12345" {
		t.Errorf("Expected order ID '12345', got %s", response.Data.OrderId)
	}
	if response.Data.ClientId != "client123" {
		t.Errorf("Expected apiClient ID 'client123', got %s", response.Data.ClientId)
	}
}

func TestPlaceOrder_Error(t *testing.T) {
	testCases := []struct {
		name           string
		mockResponse   string
		responseStatus int
		expectError    bool
		errorCode      int
	}{
		{
			name: "Order Price Issue",
			mockResponse: `{
				"code": 30001,
				"message": "Invalid price"
			}`,
			responseStatus: http.StatusOK,
			expectError:    true,
			errorCode:      30001,
		},
		{
			name: "Insufficient Balance",
			mockResponse: `{
				"code": 20003,
				"message": "Insufficient balance"
			}`,
			responseStatus: http.StatusOK,
			expectError:    true,
			errorCode:      20003,
		},
		{
			name: "Authentication Error",
			mockResponse: `{
				"code": 10003,
				"message": "Invalid signature"
			}`,
			responseStatus: http.StatusUnauthorized,
			expectError:    true,
			errorCode:      10003,
		},
		{
			name: "Network Error",
			mockResponse: `{
				"code": 10001,
				"message": "Internal server error"
			}`,
			responseStatus: http.StatusInternalServerError,
			expectError:    true,
			errorCode:      10001,
		},
		{
			name:           "HTTP Error",
			mockResponse:   `Internal Server Error`,
			responseStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAPI := NewMockAPI(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.responseStatus)
				w.Write([]byte(tc.mockResponse))
			})
			defer mockAPI.Close()

			qty := 1.0
			price := 50000.0

			orderReq := &model.OrderRequest{
				Symbol:    "BTCUSDT",
				TradeSide: model.TradeSideBuy,
				Price:     &price,
				Qty:       qty,
				Side:      model.SideOpen,
				OrderType: model.OrderTypeLimit,
				ClientID:  "client123",
			}

			response, err := mockAPI.client.PlaceOrder(context.Background(), orderReq)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}

				if tc.errorCode > 0 {
					apiErr, ok := err.(*errors.APIError)
					if !ok {
						t.Errorf("Expected *errors.APIError, got %T", err)
						return
					}

					if apiErr.Code != tc.errorCode {
						t.Errorf("Expected error code %d, got %d", tc.errorCode, apiErr.Code)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
					return
				}

				if response == nil {
					t.Errorf("Expected response but got nil")
				}
			}
		})
	}
}

func TestGetOrderDetail(t *testing.T) {
	tests := []struct {
		name          string
		request       *OrderDetailRequest
		expectedURL   string
		responseBody  string
		expectedCode  int
		expectedError error
	}{
		{
			name: "valid request with order ID",
			request: &OrderDetailRequest{
				OrderID: "123456",
			},
			expectedURL: "/api/v1/futures/trade/get_order_detail?orderId=123456",
			responseBody: `{
				"code": 0,
				"data": {
					"orderId": "123456",
					"symbol": "BTCUSDT",
					"qty": "1.0",
					"tradeQty": "0.5",
					"price": "50000",
					"side": "BUY",
					"leverage": 10,
					"orderType": "LIMIT",
					"status": "PART_FILLED",
					"clientId": "",
					"positionMode": "ONE_WAY",
					"marginMode": "ISOLATION",
					"effect": "GTC",
					"reduceOnly": false,
					"fee": "0.05",
					"realizedPNL": "0",
					"ctime": "1700000000000",
					"mtime": "1700000000000"
				}
			}`,
			expectedCode: 0,
		},
		{
			name: "valid request with client ID",
			request: &OrderDetailRequest{
				ClientID: "client123",
			},
			expectedURL: "/api/v1/futures/trade/get_order_detail?clientId=client123",
			responseBody: `{
				"code": 0,
				"data": {
					"orderId": "123456",
					"clientId": "client123",
					"symbol": "BTCUSDT",
					"qty": "1.0",
					"price": "50000",
					"side": "BUY",
					"leverage": 10,
					"orderType": "LIMIT",
					"status": "NEW",
					"positionMode": "ONE_WAY",
					"marginMode": "ISOLATION",
					"effect": "GTC",
					"reduceOnly": false,
					"fee": "0",
					"realizedPNL": "0",
					"ctime": "1700000000000",
					"mtime": "1700000000000"
				}
			}`,
			expectedCode: 0,
		},
		{
			name: "request with both order ID and client ID",
			request: &OrderDetailRequest{
				OrderID:  "123456",
				ClientID: "client123",
			},
			expectedURL: "/api/v1/futures/trade/get_order_detail?clientId=client123&orderId=123456",
			responseBody: `{
				"code": 0,
				"data": {
					"orderId": "123456",
					"clientId": "client123",
					"symbol": "BTCUSDT",
					"qty": "1.0",
					"tradeQty": "0",
					"price": "50000",
					"side": "BUY",
					"leverage": 10,
					"orderType": "LIMIT",
					"status": "NEW",
					"positionMode": "ONE_WAY",
					"marginMode": "ISOLATION",
					"effect": "GTC",
					"reduceOnly": false,
					"fee": "0",
					"realizedPNL": "0",
					"ctime": "1700000000000",
					"mtime": "1700000000000"
				}
			}`,
			expectedCode: 0,
		},
		{
			name:          "error - missing both IDs",
			request:       &OrderDetailRequest{},
			expectedError: errors.NewValidationError("request", "either orderId or clientId is required", nil),
		},
		{
			name: "API error - order not found",
			request: &OrderDetailRequest{
				OrderID: "nonexistent",
			},
			expectedURL: "/api/v1/futures/trade/get_order_detail?orderId=nonexistent",
			responseBody: `{
				"code": 20007,
				"message": "Order not found"
			}`,
			expectedCode:  20007,
			expectedError: errors.NewAPIError(20007, "Order not found", "/api/v1/futures/trade/get_order_detail", errors.ErrOrderNotFound),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var capturedURL string
			mockAPI := NewMockAPI(func(w http.ResponseWriter, r *http.Request) {
				capturedURL = r.URL.String()
				w.Header().Set("Content-Type", "application/json")
				if tt.responseBody != "" {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(tt.responseBody))
				} else {
					w.WriteHeader(http.StatusBadRequest)
				}
			})
			defer mockAPI.Close()

			response, err := mockAPI.client.GetOrderDetail(context.Background(), tt.request)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError.Error() {
					t.Errorf("Expected error %v, got %v", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tt.expectedURL != "" && capturedURL != tt.expectedURL {
				t.Errorf("Expected URL %s, got %s", tt.expectedURL, capturedURL)
			}

			if response.Code != tt.expectedCode {
				t.Errorf("Expected code %d, got %d", tt.expectedCode, response.Code)
			}

			// Validate response data if successful
			if tt.expectedCode == 0 && response.Data.OrderID != "" {
				t.Logf("Received order detail: %+v", response.Data)
			}
		})
	}
}

func TestCancelOrderBuilderCreation(t *testing.T) {
	symbol := model.ParseSymbol("BTCUSDT")

	builder := NewCancelOrderBuilder(symbol)
	cancelOrder, err := builder.Build()

	if err == nil {
		t.Fatal("Expected error due to no orders specified for cancellation, but got nil")
	}

	builder.WithOrderID("test-order-id")
	cancelOrder, err = builder.Build()
	if err != nil {
		t.Fatalf("Failed to build cancel order: %v", err)
	}

	if cancelOrder.Symbol != symbol {
		t.Errorf("Expected symbol %s, got %s", symbol, cancelOrder.Symbol)
	}
	if len(cancelOrder.OrderList) != 1 {
		t.Errorf("Expected order list with 1 item, got %d items", len(cancelOrder.OrderList))
	}
}

func TestCancelOrderBuilderMethods(t *testing.T) {
	symbol := model.ParseSymbol("BTCUSDT")
	builder := NewCancelOrderBuilder(symbol)

	orderID := "order123"
	builder.WithOrderID(orderID)
	cancelOrder, err := builder.Build()
	if err != nil {
		t.Fatalf("Failed to build cancel order: %v", err)
	}

	if len(cancelOrder.OrderList) != 1 {
		t.Fatalf("Expected 1 item in order list, got %d", len(cancelOrder.OrderList))
	}
	if cancelOrder.OrderList[0].OrderID != orderID {
		t.Errorf("WithOrderID: Expected order ID %s, got %s", orderID, cancelOrder.OrderList[0].OrderID)
	}
	if cancelOrder.OrderList[0].ClientID != "" {
		t.Errorf("WithOrderID: Expected empty apiClient ID, got %s", cancelOrder.OrderList[0].ClientID)
	}

	clientID := "client123"
	builder.WithClientID(clientID)
	cancelOrder, err = builder.Build()
	if err != nil {
		t.Fatalf("Failed to build cancel order: %v", err)
	}

	if len(cancelOrder.OrderList) != 2 {
		t.Fatalf("Expected 2 items in order list, got %d", len(cancelOrder.OrderList))
	}
	if cancelOrder.OrderList[1].ClientID != clientID {
		t.Errorf("WithClientID: Expected apiClient ID %s, got %s", clientID, cancelOrder.OrderList[1].ClientID)
	}
	if cancelOrder.OrderList[1].OrderID != "" {
		t.Errorf("WithClientID: Expected empty order ID, got %s", cancelOrder.OrderList[1].OrderID)
	}
}

func TestCancelOrderBuilder_ValidationErrors(t *testing.T) {
	testCases := []struct {
		name          string
		setupBuilder  func() *CancelOrderBuilder
		expectedError string
	}{
		{
			name: "Missing Symbol",
			setupBuilder: func() *CancelOrderBuilder {
				builder := NewCancelOrderBuilder("")
				builder.WithOrderID("order123")
				return builder
			},
			expectedError: "validation error: field symbol: is required",
		},
		{
			name: "Empty OrderList",
			setupBuilder: func() *CancelOrderBuilder {
				builder := NewCancelOrderBuilder("BTCUSDT")
				return builder
			},
			expectedError: "validation error: field orderList: must contain at least one order to cancel",
		},
		{
			name: "Order with Neither OrderID nor ClientID",
			setupBuilder: func() *CancelOrderBuilder {
				builder := NewCancelOrderBuilder("BTCUSDT")

				builder.request.OrderList = append(builder.request.OrderList, model.CancelOrderParam{
					OrderID:  "",
					ClientID: "",
				})
				return builder
			},
			expectedError: "validation error: field orderList[0]: must have either orderId or clientId",
		},
		{
			name: "Multiple Orders with One Invalid",
			setupBuilder: func() *CancelOrderBuilder {
				builder := NewCancelOrderBuilder("BTCUSDT")
				builder.WithOrderID("order123")

				builder.request.OrderList = append(builder.request.OrderList, model.CancelOrderParam{
					OrderID:  "",
					ClientID: "",
				})
				return builder
			},
			expectedError: "validation error: field orderList[1]: must have either orderId or clientId",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			builder := tc.setupBuilder()
			_, err := builder.Build()

			if err == nil {
				t.Error("Expected validation error, got nil")
				return
			}

			if err.Error() != tc.expectedError {
				t.Errorf("Expected error message '%s', got '%s'", tc.expectedError, err.Error())
			}
		})
	}
}

func TestCancelOrders(t *testing.T) {
	mockResponse := `{
		"code": 0,
		"msg": "Success",
		"data": {
			"successList": [
				{
					"orderId": "11111",
					"clientId": "22222"
				}
			],
			"failureList": [
				{
					"orderId": "33333",
					"clientId": "44444",
					"errorMsg": "Order not found",
					"errorCode": "10001"
				}
			]
		}
	}`

	mockAPI := NewMockAPI(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))

		if r.URL.Path != "/api/v1/futures/trade/cancel_orders" {
			t.Errorf("Expected request path /api/v1/futures/trade/cancel_orders, got %s", r.URL.Path)
		}

		if r.Method != http.MethodPost {
			t.Errorf("Expected request method POST, got %s", r.Method)
		}
	})
	defer mockAPI.Close()

	cancelOrderReq := &model.CancelOrderRequest{
		Symbol: "BTCUSDT",
		OrderList: []model.CancelOrderParam{
			{
				OrderID: "11111",
			},
			{
				ClientID: "22223",
			},
		},
	}

	response, err := mockAPI.client.CancelOrders(context.Background(), cancelOrderReq)
	if err != nil {
		t.Fatalf("CancelOrders returned error: %v", err)
	}

	if response.Code != 0 {
		t.Errorf("Expected code 0, got %d", response.Code)
	}
	if response.Message != "Success" {
		t.Errorf("Expected message 'Success', got %s", response.Message)
	}

	if len(response.Data.SuccessList) != 1 {
		t.Fatalf("Expected 1 item in success list, got %d", len(response.Data.SuccessList))
	}
	if response.Data.SuccessList[0].OrderId != "11111" {
		t.Errorf("Expected success order ID '11111', got %s", response.Data.SuccessList[0].OrderId)
	}
	if response.Data.SuccessList[0].ClientId != "22222" {
		t.Errorf("Expected success apiClient ID '22222', got %s", response.Data.SuccessList[0].ClientId)
	}

	if len(response.Data.FailureList) != 1 {
		t.Fatalf("Expected 1 item in failure list, got %d", len(response.Data.FailureList))
	}
	if response.Data.FailureList[0].OrderId != "33333" {
		t.Errorf("Expected failure order ID '33333', got %s", response.Data.FailureList[0].OrderId)
	}
	if response.Data.FailureList[0].ClientId != "44444" {
		t.Errorf("Expected failure apiClient ID '44444', got %s", response.Data.FailureList[0].ClientId)
	}
	if response.Data.FailureList[0].ErrorMsg != "Order not found" {
		t.Errorf("Expected failure error message 'Order not found', got %s", response.Data.FailureList[0].ErrorMsg)
	}
	if response.Data.FailureList[0].ErrorCode != "10001" {
		t.Errorf("Expected failure error code '10001', got %s", response.Data.FailureList[0].ErrorCode)
	}
}

func TestCancelOrders_Error(t *testing.T) {
	testCases := []struct {
		name           string
		mockResponse   string
		responseStatus int
		expectError    bool
		errorCode      int
	}{
		{
			name: "Authentication Error",
			mockResponse: `{
				"code": 10003,
				"message": "Invalid API key"
			}`,
			responseStatus: http.StatusUnauthorized,
			expectError:    true,
			errorCode:      10003,
		},
		{
			name: "Order Not Found",
			mockResponse: `{
				"code": 20007,
				"message": "Order not found"
			}`,
			responseStatus: http.StatusOK,
			expectError:    true,
			errorCode:      20007,
		},
		{
			name: "All Orders Successful",
			mockResponse: `{
				"code": 0,
				"msg": "Success",
				"data": {
					"successList": [
						{
							"orderId": "11111",
							"clientId": "22222"
						}
					],
					"failureList": []
				}
			}`,
			responseStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name: "All Orders Failed",
			mockResponse: `{
				"code": 0,
				"msg": "Success",
				"data": {
					"successList": [],
					"failureList": [
						{
							"orderId": "33333",
							"clientId": "44444",
							"errorMsg": "Order not found",
							"errorCode": "10001"
						}
					]
				}
			}`,
			responseStatus: http.StatusOK,
			expectError:    false,
		},
		{
			name:           "Invalid JSON Response",
			mockResponse:   `{ invalid json }`,
			responseStatus: http.StatusOK,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockAPI := NewMockAPI(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.responseStatus)
				w.Write([]byte(tc.mockResponse))
			})
			defer mockAPI.Close()

			cancelOrderReq := &model.CancelOrderRequest{
				Symbol: "BTCUSDT",
				OrderList: []model.CancelOrderParam{
					{
						OrderID: "11111",
					},
					{
						ClientID: "22223",
					},
				},
			}

			response, err := mockAPI.client.CancelOrders(context.Background(), cancelOrderReq)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}

				if tc.errorCode > 0 {
					apiErr, ok := err.(*errors.APIError)
					if !ok {
						t.Errorf("Expected *errors.APIError, got %T", err)
						return
					}

					if apiErr.Code != tc.errorCode {
						t.Errorf("Expected error code %d, got %d", tc.errorCode, apiErr.Code)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
					return
				}

				if response == nil {
					t.Errorf("Expected response but got nil")
				}
			}
		})
	}
}
