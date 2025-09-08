package bitunix

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tradingiq/bitunix-client/model"
)

func TestGetTPSLOrderHistory(t *testing.T) {
	mockResponse := `{
		"code": 0,
		"data": {
			"orderList": [
				{
					"id": "123456",
					"positionId": "78910",
					"symbol": "BTCUSDT",
					"base": "BTC",
					"quote": "USDT",
					"tpPrice": "50000",
					"tpStopType": "LAST_PRICE",
					"slPrice": "40000",
					"slStopType": "MARK_PRICE",
					"tpOrderType": "LIMIT",
					"tpOrderPrice": "50100",
					"slOrderType": "MARKET",
					"tpQty": "0.5",
					"slQty": "0.5",
					"status": "TRIGGERED",
					"ctime": "1659076670000",
					"triggerTime": "1659076680000"
				}
			],
			"total": "1"
		}
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/v1/futures/tpsl/get_history_orders", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "BTCUSDT", r.URL.Query().Get("symbol"))
		assert.Equal(t, "LONG", r.URL.Query().Get("side"))
		assert.Equal(t, "HEDGE", r.URL.Query().Get("positionMode"))
		assert.Equal(t, "1659076600000", r.URL.Query().Get("startTime"))
		assert.Equal(t, "1659076700000", r.URL.Query().Get("endTime"))
		assert.Equal(t, "", r.URL.Query().Get("skip"))
		assert.Equal(t, "10", r.URL.Query().Get("limit"))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client, err := NewApiClient("test-key", "test-secret", WithBaseURI(server.URL+"/"))
	assert.NoError(t, err)

	startTime := time.Unix(1659076600, 0)
	endTime := time.Unix(1659076700, 0)

	params := model.TPSLOrderHistoryParams{
		Symbol:       model.Symbol("BTCUSDT"),
		Side:         model.PositionSideLong,
		PositionMode: model.PositionModeHedge,
		StartTime:    &startTime,
		EndTime:      &endTime,
		Skip:         0,
		Limit:        10,
	}

	response, err := client.GetTPSLOrderHistory(context.Background(), params)
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, 0, response.Code)
	assert.Equal(t, int64(1), response.Data.Total)
	assert.Len(t, response.Data.OrderList, 1)

	order := response.Data.OrderList[0]
	assert.Equal(t, "123456", order.ID)
	assert.Equal(t, "78910", order.PositionID)
	assert.Equal(t, model.Symbol("BTCUSDT"), order.Symbol)
	assert.Equal(t, "BTC", order.Base)
	assert.Equal(t, "USDT", order.Quote)
	assert.NotNil(t, order.TpPrice)
	assert.Equal(t, 50000.0, *order.TpPrice)
	assert.NotNil(t, order.TpStopType)
	assert.Equal(t, model.StopTypeLastPrice, *order.TpStopType)
	assert.NotNil(t, order.SlPrice)
	assert.Equal(t, 40000.0, *order.SlPrice)
	assert.NotNil(t, order.SlStopType)
	assert.Equal(t, model.StopTypeMarkPrice, *order.SlStopType)
	assert.NotNil(t, order.TpOrderType)
	assert.Equal(t, model.OrderTypeLimit, *order.TpOrderType)
	assert.NotNil(t, order.TpOrderPrice)
	assert.Equal(t, 50100.0, *order.TpOrderPrice)
	assert.NotNil(t, order.SlOrderType)
	assert.Equal(t, model.OrderTypeMarket, *order.SlOrderType)
	assert.NotNil(t, order.TpQty)
	assert.Equal(t, 0.5, *order.TpQty)
	assert.NotNil(t, order.SlQty)
	assert.Equal(t, 0.5, *order.SlQty)
	assert.Equal(t, "TRIGGERED", order.Status)
	assert.Equal(t, time.Time(time.Date(2022, time.July, 29, 8, 37, 50, 0, time.Local)), order.Ctime)
	assert.NotNil(t, order.TriggerTime)
	assert.Equal(t, time.Time(time.Date(2022, time.July, 29, 8, 38, 0, 0, time.Local)), *order.TriggerTime)
}

func TestGetTPSLOrderHistoryError(t *testing.T) {
	mockResponse := `{
		"code": 20007,
		"message": "Order not found"
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	client, err := NewApiClient("test-key", "test-secret", WithBaseURI(server.URL+"/"))
	assert.NoError(t, err)

	params := model.TPSLOrderHistoryParams{
		Symbol: model.Symbol("BTCUSDT"),
	}

	response, err := client.GetTPSLOrderHistory(context.Background(), params)
	assert.Error(t, err)
	assert.Nil(t, response)
	assert.Contains(t, err.Error(), "Order not found")
}
