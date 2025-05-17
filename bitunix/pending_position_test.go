package bitunix

import (
	"context"
	"github.com/tradingiq/bitunix-client/model"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestGetPendingPositions(t *testing.T) {
	mockResponse := `{
		"code": 0,
		"data": [
			{
				"positionId": "12345678",
				"symbol": "BTCUSDT",
				"qty": "0.5",
				"entryValue": "30000",
				"side": "BUY",
				"positionMode": "HEDGE",
				"marginMode": "ISOLATION",
				"leverage": 100,
				"fees": "1.5",
				"funding": "0.25",
				"realizedPNL": "150.75",
				"margin": "300",
				"unrealizedPNL": "75.5",
				"liqPrice": "29100.5",
				"marginRate": "0.1",
				"avgOpenPrice": "60000",
				"ctime": "1659076670000",
				"mtime": "1659086670000"
			}
		]
	}`

	mockAPI := NewMockAPI(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))

		if r.URL.Path != "/api/v1/futures/position/get_pending_positions" {
			t.Errorf("Expected request path /api/v1/futures/position/get_pending_positions, got %s", r.URL.Path)
		}

		if r.Method != http.MethodGet {
			t.Errorf("Expected request method GET, got %s", r.Method)
		}

		symbol := r.URL.Query().Get("symbol")
		if symbol != "" && symbol != "BTCUSDT" {
			t.Errorf("Expected symbol BTCUSDT, got %s", symbol)
		}

		positionId := r.URL.Query().Get("positionId")
		if positionId != "" && positionId != "12345678" {
			t.Errorf("Expected positionId 12345678, got %s", positionId)
		}
	})
	defer mockAPI.Close()

	client := mockAPI.client

	t.Run("Success", func(t *testing.T) {
		params := model.PendingPositionParams{
			Symbol:     "BTCUSDT",
			PositionID: "12345678",
		}

		resp, err := client.GetPendingPositions(context.Background(), params)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resp.Code != 0 {
			t.Errorf("Expected code 0, got %d", resp.Code)
		}

		if len(resp.Data) != 1 {
			t.Fatalf("Expected 1 position, got %d", len(resp.Data))
		}

		position := resp.Data[0]
		if position.PositionID != "12345678" {
			t.Errorf("Expected position ID 12345678, got %s", position.PositionID)
		}

		if position.Symbol != "BTCUSDT" {
			t.Errorf("Expected symbol BTCUSDT, got %s", position.Symbol)
		}

		if position.Qty != 0.5 {
			t.Errorf("Expected qty 0.5, got %f", position.Qty)
		}

		if position.EntryValue != 30000 {
			t.Errorf("Expected entryValue 30000, got %f", position.EntryValue)
		}

		if position.Side != model.TradeSideBuy {
			t.Errorf("Expected side BUY, got %s", position.Side)
		}

		if position.PositionMode != model.PositionModeHedge {
			t.Errorf("Expected position mode HEDGE, got %s", position.PositionMode)
		}

		if position.MarginMode != model.MarginModeIsolation {
			t.Errorf("Expected margin mode ISOLATION, got %s", position.MarginMode)
		}

		if position.Leverage != 100 {
			t.Errorf("Expected leverage 100, got %d", position.Leverage)
		}

		if position.Fees != 1.5 {
			t.Errorf("Expected fees 1.5, got %f", position.Fees)
		}

		if position.Funding != 0.25 {
			t.Errorf("Expected funding 0.25, got %f", position.Funding)
		}

		if position.RealizedPNL != 150.75 {
			t.Errorf("Expected realizedPNL 150.75, got %f", position.RealizedPNL)
		}

		if position.Margin != 300 {
			t.Errorf("Expected margin 300, got %f", position.Margin)
		}

		if position.UnrealizedPNL != 75.5 {
			t.Errorf("Expected unrealizedPNL 75.5, got %f", position.UnrealizedPNL)
		}

		if position.LiqPrice != 29100.5 {
			t.Errorf("Expected liqPrice 29100.5, got %f", position.LiqPrice)
		}

		if position.MarginRate != 0.1 {
			t.Errorf("Expected marginRate 0.1, got %f", position.MarginRate)
		}

		if position.AvgOpenPrice != 60000 {
			t.Errorf("Expected avgOpenPrice 60000, got %f", position.AvgOpenPrice)
		}

		expectedCtime := time.Unix(0, 1659076670000*1000000)
		if !position.CreateTime.Equal(expectedCtime) {
			t.Errorf("Expected ctime %v, got %v", expectedCtime, position.CreateTime)
		}

		expectedMtime := time.Unix(0, 1659086670000*1000000)
		if !position.ModifyTime.Equal(expectedMtime) {
			t.Errorf("Expected mtime %v, got %v", expectedMtime, position.ModifyTime)
		}
	})

	t.Run("Error response", func(t *testing.T) {
		errorAPI := NewMockAPI(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"code": 400, "msg": "Invalid request"}`))
		})
		defer errorAPI.Close()

		errorClient := errorAPI.client

		params := model.PendingPositionParams{}
		resp, err := errorClient.GetPendingPositions(context.Background(), params)
		if err == nil {
			t.Fatal("Expected error, got none")
		}
		if resp != nil {
			t.Error("Expected nil response on error")
		}
	})

	t.Run("Empty parameters", func(t *testing.T) {
		params := model.PendingPositionParams{}
		resp, err := client.GetPendingPositions(context.Background(), params)
		if err != nil {
			t.Fatalf("Expected no error with empty params, got %v", err)
		}
		if resp == nil {
			t.Error("Expected response, got nil")
		}
	})
}

func TestPendingPositionUnmarshal(t *testing.T) {
	jsonData := `{
		"positionId": "12345678",
		"symbol": "BTCUSDT",
		"qty": "0.5",
		"entryValue": "30000",
		"side": "BUY",
		"positionMode": "HEDGE",
		"marginMode": "ISOLATION",
		"leverage": 100,
		"fees": "1.5",
		"funding": "0.25",
		"realizedPNL": "150.75",
		"margin": "300",
		"unrealizedPNL": "75.5",
		"liqPrice": "29100.5",
		"marginRate": "0.1",
		"avgOpenPrice": "60000",
		"ctime": "1659076670000",
		"mtime": "1659086670000"
	}`

	var position model.PendingPosition
	err := position.UnmarshalJSON([]byte(jsonData))
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if position.PositionID != "12345678" {
		t.Errorf("Expected position ID 12345678, got %s", position.PositionID)
	}

	if position.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol BTCUSDT, got %s", position.Symbol)
	}

	if position.Qty != 0.5 {
		t.Errorf("Expected qty 0.5, got %f", position.Qty)
	}

	if position.EntryValue != 30000 {
		t.Errorf("Expected entryValue 30000, got %f", position.EntryValue)
	}

	if position.Side != model.TradeSideBuy {
		t.Errorf("Expected side BUY, got %s", position.Side)
	}

	if position.PositionMode != model.PositionModeHedge {
		t.Errorf("Expected position mode HEDGE, got %s", position.PositionMode)
	}

	if position.MarginMode != model.MarginModeIsolation {
		t.Errorf("Expected margin mode ISOLATION, got %s", position.MarginMode)
	}

	if position.Leverage != 100 {
		t.Errorf("Expected leverage 100, got %d", position.Leverage)
	}

	if position.Fees != 1.5 {
		t.Errorf("Expected fees 1.5, got %f", position.Fees)
	}

	if position.Funding != 0.25 {
		t.Errorf("Expected funding 0.25, got %f", position.Funding)
	}

	if position.RealizedPNL != 150.75 {
		t.Errorf("Expected realizedPNL 150.75, got %f", position.RealizedPNL)
	}

	if position.Margin != 300 {
		t.Errorf("Expected margin 300, got %f", position.Margin)
	}

	if position.UnrealizedPNL != 75.5 {
		t.Errorf("Expected unrealizedPNL 75.5, got %f", position.UnrealizedPNL)
	}

	if position.LiqPrice != 29100.5 {
		t.Errorf("Expected liqPrice 29100.5, got %f", position.LiqPrice)
	}

	if position.MarginRate != 0.1 {
		t.Errorf("Expected marginRate 0.1, got %f", position.MarginRate)
	}

	if position.AvgOpenPrice != 60000 {
		t.Errorf("Expected avgOpenPrice 60000, got %f", position.AvgOpenPrice)
	}

	expectedCtime := time.Unix(0, 1659076670000*1000000)
	if !position.CreateTime.Equal(expectedCtime) {
		t.Errorf("Expected ctime %v, got %v", expectedCtime, position.CreateTime)
	}

	expectedMtime := time.Unix(0, 1659086670000*1000000)
	if !position.ModifyTime.Equal(expectedMtime) {
		t.Errorf("Expected mtime %v, got %v", expectedMtime, position.ModifyTime)
	}
}

func TestPendingPositionUnmarshalErrors(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		errMsg   string
	}{
		{
			name:     "Invalid JSON",
			jsonData: `{invalid json}`,
			errMsg:   "invalid character",
		},
		{
			name:     "Invalid qty",
			jsonData: `{"positionId": "123", "qty": "invalid"}`,
			errMsg:   "failed to parse qty",
		},
		{
			name:     "Invalid side",
			jsonData: `{"positionId": "123", "side": "INVALID"}`,
			errMsg:   "invalid side",
		},
		{
			name:     "Invalid margin mode",
			jsonData: `{"positionId": "123", "side": "BUY", "marginMode": "INVALID"}`,
			errMsg:   "invalid margin mode",
		},
		{
			name:     "Invalid position mode",
			jsonData: `{"positionId": "123", "side": "BUY", "marginMode": "CROSS", "positionMode": "INVALID"}`,
			errMsg:   "invalid position mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var position model.PendingPosition
			err := position.UnmarshalJSON([]byte(tt.jsonData))
			if err == nil {
				t.Fatal("Expected error, got none")
			}
			if !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("Expected error containing '%s', got '%s'", tt.errMsg, err.Error())
			}
		})
	}
}
