package bitunix

import (
	"context"
	"github.com/tradingiq/bitunix-client/model"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestGetOrderHistory(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path != "/api/v1/futures/trade/get_history_orders" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("symbol") != "BTCUSDT" {
			t.Errorf("unexpected symbol: %s", query.Get("symbol"))
		}

		if query.Get("orderId") != "12345" {
			t.Errorf("unexpected orderId: %s", query.Get("orderId"))
		}

		if query.Get("limit") != "5" {
			t.Errorf("unexpected limit: %s", query.Get("limit"))
		}

		response := `{
			"code": 0,
			"msg": "Success",
			"data": {
				"orderList": [
					{
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
					}
				],
				"total": "1"
			}
		}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	bitunixClient, _ := NewApiClient("test-api-key", "test-api-secret", WithBaseURI(server.URL))

	startTime := time.Date(2022, 7, 29, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2022, 7, 30, 0, 0, 0, 0, time.UTC)

	params := model.OrderHistoryParams{
		Symbol:    "BTCUSDT",
		OrderID:   "12345",
		StartTime: &startTime,
		EndTime:   &endTime,
		Limit:     5,
	}

	response, err := bitunixClient.GetOrderHistory(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.Code != 0 {
		t.Errorf("unexpected response code: %d", response.Code)
	}

	if response.Message != "Success" {
		t.Errorf("unexpected response message: %s", response.Message)
	}

	if len(response.Data.Orders) != 1 {
		t.Fatalf("unexpected number of orders: %d", len(response.Data.Orders))
	}

	order := response.Data.Orders[0]
	if order.OrderID != "12345" {
		t.Errorf("unexpected orderId: %s", order.OrderID)
	}

	if order.Symbol != "BTCUSDT" {
		t.Errorf("unexpected symbol: %s", order.Symbol)
	}

	if order.Quantity != 1.5 {
		t.Errorf("unexpected quantity: %f", order.Quantity)
	}

	if order.TradeQuantity != 1.0 {
		t.Errorf("unexpected trade quantity: %f", order.TradeQuantity)
	}

	if order.PositionMode != model.PositionModeHedge {
		t.Errorf("unexpected position mode: %s", order.PositionMode)
	}

	if order.MarginMode != model.MarginModeCross {
		t.Errorf("unexpected margin mode: %s", order.MarginMode)
	}

	if order.Leverage != 10 {
		t.Errorf("unexpected leverage: %d", order.Leverage)
	}

	if order.Price != "MARKET" {
		t.Errorf("unexpected price: %s", order.Price)
	}

	if order.Side != model.TradeSideBuy {
		t.Errorf("unexpected side: %s", order.Side)
	}

	if order.OrderType != model.OrderTypeLimit {
		t.Errorf("unexpected order type: %s", order.OrderType)
	}

	if order.Effect != model.TimeInForceGTC {
		t.Errorf("unexpected time in force: %s", order.Effect)
	}

	if order.ClientID != "client_123" {
		t.Errorf("unexpected apiClient ID: %s", order.ClientID)
	}

	if order.ReduceOnly != false {
		t.Errorf("unexpected reduce only flag: %v", order.ReduceOnly)
	}

	if order.Status != "FILLED" {
		t.Errorf("unexpected status: %s", order.Status)
	}

	if order.Fee != 0.6 {
		t.Errorf("unexpected fee: %f", order.Fee)
	}

	if order.RealizedPNL != 100 {
		t.Errorf("unexpected realized PNL: %f", order.RealizedPNL)
	}

	if order.TpPrice != 45000 {
		t.Errorf("unexpected TP price: %f", order.TpPrice)
	}

	if order.TpStopType != model.StopTypeMarkPrice {
		t.Errorf("unexpected TP stop type: %s", order.TpStopType)
	}

	if order.TpOrderType != model.OrderTypeLimit {
		t.Errorf("unexpected TP order type: %s", order.TpOrderType)
	}

	if order.TpOrderPrice != 45000 {
		t.Errorf("unexpected TP order price: %f", order.TpOrderPrice)
	}

	if order.SlPrice != 35000 {
		t.Errorf("unexpected SL price: %f", order.SlPrice)
	}

	if order.SlStopType != model.StopTypeMarkPrice {
		t.Errorf("unexpected SL stop type: %s", order.SlStopType)
	}

	if order.SlOrderType != model.OrderTypeMarket {
		t.Errorf("unexpected SL order type: %s", order.SlOrderType)
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
