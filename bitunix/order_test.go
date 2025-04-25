package bitunix

import (
	"context"
	"encoding/json"
	"github.com/tradingiq/bitunix-client/rest"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockAPI struct {
	server *httptest.Server
	client *API
}

func NewMockAPI(handler http.HandlerFunc) *MockAPI {
	server := httptest.NewServer(handler)
	apiClient, _ := NewTestClient(server.URL)
	client := &API{restClient: apiClient}
	return &MockAPI{
		server: server,
		client: client,
	}
}

func (m *MockAPI) Close() {
	m.server.Close()
}

func TestOrderBuilderCreation(t *testing.T) {
	symbol := "BTCUSDT"
	side := TradeActionBuy
	tradeSide := TradeSideOpen
	qty := 1.0

	builder := NewOrderBuilder(symbol, side, tradeSide, qty)
	order := builder.Build()

	if order.Symbol != symbol {
		t.Errorf("Expected symbol %s, got %s", symbol, order.Symbol)
	}
	if order.TradeAction != side {
		t.Errorf("Expected side %s, got %s", side, order.TradeAction)
	}
	if order.TradeSide != tradeSide {
		t.Errorf("Expected tradeSide %s, got %s", tradeSide, order.TradeSide)
	}
	if *order.Qty != qty {
		t.Errorf("Expected qty %f, got %f", qty, *order.Qty)
	}
	if order.OrderType != OrderTypeMarket {
		t.Errorf("Expected order type %s, got %s", OrderTypeMarket, order.OrderType)
	}
	if order.ReduceOnly != false {
		t.Errorf("Expected reduceOnly %v, got %v", false, order.ReduceOnly)
	}
	if order.Effect != TimeInForceGTC {
		t.Errorf("Expected effect %s, got %s", TimeInForceGTC, order.Effect)
	}
}

func TestOrderBuilderMethods(t *testing.T) {
	symbol := "BTCUSDT"
	side := TradeActionBuy
	tradeSide := TradeSideOpen
	qty := 1.0

	builder := NewOrderBuilder(symbol, side, tradeSide, qty)

	orderType := OrderTypeLimit
	builder.WithOrderType(orderType)
	order := builder.Build()
	if order.OrderType != orderType {
		t.Errorf("WithOrderType: Expected order type %s, got %s", orderType, order.OrderType)
	}

	price := 50000.0
	builder.WithPrice(price)
	order = builder.Build()
	if *order.Price != price {
		t.Errorf("WithPrice: Expected price %f, got %f", price, *order.Price)
	}

	positionID := "position123"
	builder.WithPositionID(positionID)
	order = builder.Build()
	if order.PositionID != positionID {
		t.Errorf("WithPositionID: Expected position ID %s, got %s", positionID, order.PositionID)
	}

	reduceOnly := true
	builder.WithReduceOnly(reduceOnly)
	order = builder.Build()
	if order.ReduceOnly != reduceOnly {
		t.Errorf("WithReduceOnly: Expected reduce only %v, got %v", reduceOnly, order.ReduceOnly)
	}

	tif := TimeInForceIOC
	builder.WithTimeInForce(tif)
	order = builder.Build()
	if order.Effect != tif {
		t.Errorf("WithTimeInForce: Expected time in force %s, got %s", tif, order.Effect)
	}

	clientID := "client123"
	builder.WithClientID(clientID)
	order = builder.Build()
	if order.ClientID != clientID {
		t.Errorf("WithClientID: Expected client ID %s, got %s", clientID, order.ClientID)
	}

	tpPrice := 55000.0
	tpStopType := StopTypeLastPrice
	tpOrderType := OrderTypeLimit
	builder.WithTakeProfit(tpPrice, tpStopType, tpOrderType)
	order = builder.Build()
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
	order = builder.Build()
	if *order.TpOrderPrice != tpOrderPrice {
		t.Errorf("WithTakeProfitPrice: Expected TP order price %f, got %f", tpOrderPrice, *order.TpOrderPrice)
	}

	slPrice := 45000.0
	slStopType := StopTypeMarkPrice
	slOrderType := OrderTypeMarket
	builder.WithStopLoss(slPrice, slStopType, slOrderType)
	order = builder.Build()
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
	order = builder.Build()
	if *order.SlOrderPrice != slOrderPrice {
		t.Errorf("WithStopLossPrice: Expected SL order price %f, got %f", slOrderPrice, *order.SlOrderPrice)
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
		TradeAction:  TradeActionBuy,
		Price:        &price,
		Qty:          &qty,
		PositionID:   "position123",
		TradeSide:    TradeSideOpen,
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

	orderReq := &OrderRequest{
		Symbol:      "BTCUSDT",
		TradeAction: TradeActionBuy,
		Price:       &price,
		Qty:         &qty,
		TradeSide:   TradeSideOpen,
		OrderType:   OrderTypeLimit,
		ClientID:    "client123",
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
		t.Errorf("Expected client ID 'client123', got %s", response.Data.ClientId)
	}
}

func TestCancelOrderBuilderCreation(t *testing.T) {
	symbol := "BTCUSDT"

	builder := NewCancelOrderBuilder(symbol)
	cancelOrder := builder.Build()

	if cancelOrder.Symbol != symbol {
		t.Errorf("Expected symbol %s, got %s", symbol, cancelOrder.Symbol)
	}
	if len(cancelOrder.OrderList) != 0 {
		t.Errorf("Expected empty order list, got %d items", len(cancelOrder.OrderList))
	}
}

func TestCancelOrderBuilderMethods(t *testing.T) {
	symbol := "BTCUSDT"
	builder := NewCancelOrderBuilder(symbol)

	orderID := "order123"
	builder.WithOrderID(orderID)
	cancelOrder := builder.Build()

	if len(cancelOrder.OrderList) != 1 {
		t.Fatalf("Expected 1 item in order list, got %d", len(cancelOrder.OrderList))
	}
	if cancelOrder.OrderList[0].OrderID != orderID {
		t.Errorf("WithOrderID: Expected order ID %s, got %s", orderID, cancelOrder.OrderList[0].OrderID)
	}
	if cancelOrder.OrderList[0].ClientID != "" {
		t.Errorf("WithOrderID: Expected empty client ID, got %s", cancelOrder.OrderList[0].ClientID)
	}

	clientID := "client123"
	builder.WithClientID(clientID)
	cancelOrder = builder.Build()

	if len(cancelOrder.OrderList) != 2 {
		t.Fatalf("Expected 2 items in order list, got %d", len(cancelOrder.OrderList))
	}
	if cancelOrder.OrderList[1].ClientID != clientID {
		t.Errorf("WithClientID: Expected client ID %s, got %s", clientID, cancelOrder.OrderList[1].ClientID)
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

	cancelOrderReq := &CancelOrderRequest{
		Symbol: "BTCUSDT",
		OrderList: []CancelOrderParam{
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
		t.Errorf("Expected success client ID '22222', got %s", response.Data.SuccessList[0].ClientId)
	}

	if len(response.Data.FailureList) != 1 {
		t.Fatalf("Expected 1 item in failure list, got %d", len(response.Data.FailureList))
	}
	if response.Data.FailureList[0].OrderId != "33333" {
		t.Errorf("Expected failure order ID '33333', got %s", response.Data.FailureList[0].OrderId)
	}
	if response.Data.FailureList[0].ClientId != "44444" {
		t.Errorf("Expected failure client ID '44444', got %s", response.Data.FailureList[0].ClientId)
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
