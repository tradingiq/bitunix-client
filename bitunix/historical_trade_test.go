package bitunix

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/tradingiq/bitunix-client/model"
)

func TestGetTradeHistory(t *testing.T) {
	mockResponse := `{
		"code": 0,
		"msg": "success",
		"data": {
			"tradeList": [
				{
					"tradeId": "trade123",
					"orderId": "order123",
					"symbol": "BTCUSDT",
					"qty": "1.5",
					"positionMode": "ONE_WAY",
					"marginMode": "CROSS",
					"leverage": 10,
					"price": "50000",
					"side": "BUY",
					"orderType": "LIMIT",
					"effect": "GTC",
					"clientId": "client123",
					"reduceOnly": false,
					"fee": "5",
					"realizedPNL": "3000",
					"ctime": "1642521600000",
					"roleType": "TAKER"
				}
			],
			"total": "1"
		}
	}`

	mockAPI := NewMockAPI(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))

		if r.URL.Path != "/api/v1/futures/trade/get_history_trades" {
			t.Errorf("Expected request path /api/v1/futures/trade/get_history_trades, got %s", r.URL.Path)
		}

		if r.Method != http.MethodGet {
			t.Errorf("Expected request method GET, got %s", r.Method)
		}

		q := r.URL.Query()

		if q.Get("symbol") != "BTCUSDT" {
			t.Errorf("Expected symbol query param BTCUSDT, got %s", q.Get("symbol"))
		}

		if q.Get("orderId") != "order123" {
			t.Errorf("Expected orderId query param order123, got %s", q.Get("orderId"))
		}

		if q.Get("positionId") != "pos123" {
			t.Errorf("Expected positionId query param pos123, got %s", q.Get("positionId"))
		}

		if q.Get("limit") != "10" {
			t.Errorf("Expected limit query param 10, got %s", q.Get("limit"))
		}

		if q.Get("skip") != "5" {
			t.Errorf("Expected skip query param 5, got %s", q.Get("skip"))
		}

		if q.Get("startTime") == "" {
			t.Errorf("Expected startTime query param, got empty")
		}

		if q.Get("endTime") == "" {
			t.Errorf("Expected endTime query param, got empty")
		}
	})
	defer mockAPI.Close()

	startTime := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2022, 1, 31, 0, 0, 0, 0, time.UTC)

	params := model.TradeHistoryParams{
		Symbol:     "BTCUSDT",
		OrderID:    "order123",
		PositionID: "pos123",
		StartTime:  &startTime,
		EndTime:    &endTime,
		Skip:       5,
		Limit:      10,
	}

	response, err := mockAPI.client.GetTradeHistory(context.Background(), params)
	if err != nil {
		t.Fatalf("GetTradeHistory returned error: %v", err)
	}

	if response.Code != 0 {
		t.Errorf("Expected code 0, got %d", response.Code)
	}

	if response.Message != "success" {
		t.Errorf("Expected message 'success', got %s", response.Message)
	}

	if len(response.Data.Trades) != 1 {
		t.Fatalf("Expected 1 trade, got %d", len(response.Data.Trades))
	}

	trade := response.Data.Trades[0]

	if trade.TradeID != "trade123" {
		t.Errorf("Expected trade ID 'trade123', got %s", trade.TradeID)
	}

	if trade.OrderID != "order123" {
		t.Errorf("Expected order ID 'order123', got %s", trade.OrderID)
	}

	if trade.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol 'BTCUSDT', got %s", trade.Symbol)
	}

	if trade.Quantity != 1.5 {
		t.Errorf("Expected quantity 1.5, got %f", trade.Quantity)
	}

	if trade.PositionMode != model.PositionModeOneWay {
		t.Errorf("Expected positionMode 'ONE_WAY', got %s", trade.PositionMode)
	}

	if trade.MarginMode != model.MarginModeCross {
		t.Errorf("Expected marginMode 'CROSS', got %s", trade.MarginMode)
	}

	if trade.Leverage != 10 {
		t.Errorf("Expected leverage 10, got %d", trade.Leverage)
	}

	if trade.Price != 50000 {
		t.Errorf("Expected price 50000, got %f", trade.Price)
	}

	if trade.Side != model.TradeSideBuy {
		t.Errorf("Expected side 'OPEN', got %s", trade.Side)
	}

	if trade.OrderType != model.OrderTypeLimit {
		t.Errorf("Expected orderType 'LIMIT', got %s", trade.OrderType)
	}

	if trade.Effect != "GTC" {
		t.Errorf("Expected effect 'GTC', got %s", trade.Effect)
	}

	if trade.ClientID != "client123" {
		t.Errorf("Expected clientId 'client123', got %s", trade.ClientID)
	}

	if trade.ReduceOnly != false {
		t.Errorf("Expected reduceOnly false, got %v", trade.ReduceOnly)
	}

	if trade.Fee != 5 {
		t.Errorf("Expected fee 5, got %f", trade.Fee)
	}

	if trade.RealizedPNL != 3000 {
		t.Errorf("Expected realizedPNL 3000, got %f", trade.RealizedPNL)
	}

	expectedTime := time.Date(2022, 1, 18, 16, 0, 0, 0, time.UTC)
	if !trade.CreateTime.Equal(expectedTime) {
		t.Errorf("Expected createTime %v, got %v", expectedTime, trade.CreateTime)
	}

	if trade.RoleType != "TAKER" {
		t.Errorf("Expected roleType 'TAKER', got %s", trade.RoleType)
	}
}
