package bitunix

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestGetPositionHistory(t *testing.T) {
	mockResponse := `{
		"code": 0,
		"msg": "success",
		"data": {
			"positionList": [
				{
					"positionId": "pos123",
					"symbol": "BTCUSDT",
					"maxQty": "1.5",
					"entryPrice": "50000",
					"closePrice": "52000",
					"liqQty": "0",
					"side": "BUY",
					"positionMode": "ONE_WAY",
					"marginMode": "CROSS",
					"leverage": "10",
					"fee": "5",
					"funding": "1.25",
					"realizedPNL": "3000",
					"liqPrice": "45000",
					"ctime": "1642521600000",
					"mtime": "1642608000000"
				}
			],
			"total": "1"
		}
	}`

	mockAPI := NewMockAPI(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))

		if r.URL.Path != "/api/v1/futures/position/get_history_positions" {
			t.Errorf("Expected request path /api/v1/futures/position/get_history_positions, got %s", r.URL.Path)
		}
		
		if r.Method != http.MethodGet {
			t.Errorf("Expected request method GET, got %s", r.Method)
		}

		q := r.URL.Query()

		if q.Get("symbol") != "BTCUSDT" {
			t.Errorf("Expected symbol query param BTCUSDT, got %s", q.Get("symbol"))
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

	params := &PositionHistoryParams{
		Symbol:     "BTCUSDT",
		PositionID: "pos123",
		StartTime:  &startTime,
		EndTime:    &endTime,
		Skip:       5,
		Limit:      10,
	}

	response, err := mockAPI.client.GetPositionHistory(context.Background(), params)
	if err != nil {
		t.Fatalf("GetPositionHistory returned error: %v", err)
	}

	if response.Code != 0 {
		t.Errorf("Expected code 0, got %d", response.Code)
	}

	if response.Message != "success" {
		t.Errorf("Expected message 'success', got %s", response.Message)
	}

	if len(response.Data.Positions) != 1 {
		t.Fatalf("Expected 1 position, got %d", len(response.Data.Positions))
	}

	pos := response.Data.Positions[0]

	if pos.PositionID != "pos123" {
		t.Errorf("Expected position ID 'pos123', got %s", pos.PositionID)
	}

	if pos.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol 'BTCUSDT', got %s", pos.Symbol)
	}

	if pos.MaxQty != 1.5 {
		t.Errorf("Expected maxQty 1.5, got %f", pos.MaxQty)
	}

	if pos.EntryPrice != 50000 {
		t.Errorf("Expected entryPrice 50000, got %f", pos.EntryPrice)
	}

	if pos.ClosePrice != 52000 {
		t.Errorf("Expected closePrice 52000, got %f", pos.ClosePrice)
	}

	if pos.Side != "BUY" {
		t.Errorf("Expected side 'BUY', got %s", pos.Side)
	}

	if pos.PositionMode != "ONE_WAY" {
		t.Errorf("Expected positionMode 'ONE_WAY', got %s", pos.PositionMode)
	}

	if pos.MarginMode != "CROSS" {
		t.Errorf("Expected marginMode 'CROSS', got %s", pos.MarginMode)
	}

	if pos.Leverage != "10" {
		t.Errorf("Expected leverage '10', got %s", pos.Leverage)
	}

	if pos.Fee != 5 {
		t.Errorf("Expected fee 5, got %f", pos.Fee)
	}

	if pos.Funding != 1.25 {
		t.Errorf("Expected funding 1.25, got %f", pos.Funding)
	}

	if pos.RealizedPNL != 3000 {
		t.Errorf("Expected realizedPNL 3000, got %f", pos.RealizedPNL)
	}

	if pos.LiqPrice != 45000 {
		t.Errorf("Expected liqPrice 45000, got %f", pos.LiqPrice)
	}

	expectedCtime := time.Date(2022, 1, 18, 16, 0, 0, 0, time.UTC)
	if !pos.Ctime.Equal(expectedCtime) {
		t.Errorf("Expected ctime %v, got %v", expectedCtime, pos.Ctime)
	}

	expectedMtime := time.Date(2022, 1, 19, 16, 0, 0, 0, time.UTC)
	if !pos.Mtime.Equal(expectedMtime) {
		t.Errorf("Expected mtime %v, got %v", expectedMtime, pos.Mtime)
	}
}

func TestHistoricalPositionUnmarshalJSON(t *testing.T) {
	jsonStr := `{
		"positionId": "pos123",
		"symbol": "BTCUSDT",
		"maxQty": "1.5",
		"entryPrice": "50000",
		"closePrice": "52000",
		"liqQty": "0",
		"side": "BUY",
		"positionMode": "ONE_WAY",
		"marginMode": "CROSS",
		"leverage": "10",
		"fee": "5",
		"funding": "1.25",
		"realizedPNL": "3000",
		"liqPrice": "45000",
		"ctime": "1642521600000",
		"mtime": "1642608000000"
	}`

	var pos HistoricalPosition
	err := json.Unmarshal([]byte(jsonStr), &pos)
	if err != nil {
		t.Fatalf("Failed to unmarshal historical position: %v", err)
	}

	if pos.PositionID != "pos123" {
		t.Errorf("Expected positionId 'pos123', got %s", pos.PositionID)
	}

	if pos.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol 'BTCUSDT', got %s", pos.Symbol)
	}

	if pos.MaxQty != 1.5 {
		t.Errorf("Expected maxQty 1.5, got %f", pos.MaxQty)
	}

	if pos.EntryPrice != 50000 {
		t.Errorf("Expected entryPrice 50000, got %f", pos.EntryPrice)
	}

	if pos.ClosePrice != 52000 {
		t.Errorf("Expected closePrice 52000, got %f", pos.ClosePrice)
	}

	if pos.Side != "BUY" {
		t.Errorf("Expected side 'BUY', got %s", pos.Side)
	}

	if pos.PositionMode != "ONE_WAY" {
		t.Errorf("Expected positionMode 'ONE_WAY', got %s", pos.PositionMode)
	}

	if pos.MarginMode != "CROSS" {
		t.Errorf("Expected marginMode 'CROSS', got %s", pos.MarginMode)
	}

	if pos.Leverage != "10" {
		t.Errorf("Expected leverage '10', got %s", pos.Leverage)
	}

	if pos.Fee != 5 {
		t.Errorf("Expected fee 5, got %f", pos.Fee)
	}

	if pos.Funding != 1.25 {
		t.Errorf("Expected funding 1.25, got %f", pos.Funding)
	}

	if pos.RealizedPNL != 3000 {
		t.Errorf("Expected realizedPNL 3000, got %f", pos.RealizedPNL)
	}

	if pos.LiqPrice != 45000 {
		t.Errorf("Expected liqPrice 45000, got %f", pos.LiqPrice)
	}

	expectedCtime := time.Date(2022, 1, 18, 16, 0, 0, 0, time.UTC)
	if !pos.Ctime.Equal(expectedCtime) {
		t.Errorf("Expected ctime %v, got %v", expectedCtime, pos.Ctime)
	}

	expectedMtime := time.Date(2022, 1, 19, 16, 0, 0, 0, time.UTC)
	if !pos.Mtime.Equal(expectedMtime) {
		t.Errorf("Expected mtime %v, got %v", expectedMtime, pos.Mtime)
	}
}
