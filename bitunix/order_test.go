package bitunix

import (
	"context"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/rest"
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
	if *order.Qty != qty {
		t.Errorf("Expected qty %f, got %f", qty, *order.Qty)
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
		Qty:       &qty,
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

func NewTestClient(baseURL string) (*rest.Client, error) {
	apiClient, err := rest.New(baseURL)
	if err != nil {
		return nil, err
	}

	apiClient.SetOptions(rest.WithRequestSigner(func(req *http.Request, body []byte) error {
		return nil
	}))

	return apiClient, nil
}
