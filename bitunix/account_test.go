package bitunix

import (
	"context"
	"encoding/json"
	"github.com/tradingiq/bitunix-client/rest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccountBalance(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/futures/account" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		query := r.URL.Query()
		if query.Get("marginCoin") != "USDT" {
			t.Errorf("unexpected marginCoin: %s", query.Get("marginCoin"))
		}

		response := `{
			"code": 0,
			"data": {
					"marginCoin": "USDT",
					"available": "1000",
					"frozen": "0",
					"margin": "10",
					"transfer": "1000",
					"positionMode": "HEDGE",
					"crossUnrealizedPNL": "2",
					"isolationUnrealizedPNL": "0",
					"bonus": "0"
			},
			"msg": "Success"
		}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	serverURL := server.URL

	client, err := rest.New(serverURL)
	if err != nil {
		t.Fatalf("failed to create restClient client: %v", err)
	}
	bitunixClient := New(client, "test-restClient-key", "test-restClient-secret")

	params := AccountBalanceParams{
		MarginCoin: "USDT",
	}

	response, err := bitunixClient.GetAccountBalance(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if response.Code != 0 {
		t.Errorf("unexpected response code: %d", response.Code)
	}

	if response.Message != "Success" {
		t.Errorf("unexpected response message: %s", response.Message)
	}

	balance := response.Data
	if balance.MarginCoin != "USDT" {
		t.Errorf("unexpected marginCoin: %s", balance.MarginCoin)
	}

	if balance.Available != 1000 {
		t.Errorf("unexpected available amount: %f", balance.Available)
	}

	if balance.Frozen != 0 {
		t.Errorf("unexpected frozen amount: %f", balance.Frozen)
	}

	if balance.Margin != 10 {
		t.Errorf("unexpected margin amount: %f", balance.Margin)
	}

	if balance.Transfer != 1000 {
		t.Errorf("unexpected transfer amount: %f", balance.Transfer)
	}

	if balance.PositionMode != TradePositionModeHedge {
		t.Errorf("unexpected position mode: %s", balance.PositionMode)
	}

	if balance.CrossUnrealizedPNL != 2 {
		t.Errorf("unexpected cross unrealized PNL: %f", balance.CrossUnrealizedPNL)
	}

	if balance.IsolationUnrealizedPNL != 0 {
		t.Errorf("unexpected isolation unrealized PNL: %f", balance.IsolationUnrealizedPNL)
	}

	if balance.Bonus != 0 {
		t.Errorf("unexpected bonus: %f", balance.Bonus)
	}
}

func TestAccountBalanceParamsValidation(t *testing.T) {
	// Test missing marginCoin parameter
	client, err := rest.New("http://example.com")
	if err != nil {
		t.Fatalf("failed to create restClient client: %v", err)
	}
	bitunixClient := New(client, "test-restClient-key", "test-restClient-secret")

	params := AccountBalanceParams{}
	_, err = bitunixClient.GetAccountBalance(context.Background(), params)

	if err == nil {
		t.Fatal("expected error for missing marginCoin, got none")
	}

	if err.Error() != "marginCoin is required" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}

func TestAccountBalanceEntry_UnmarshalJSON(t *testing.T) {
	// Test valid JSON parsing
	jsonData := `{
		"marginCoin": "USDT",
		"available": "1000.50",
		"frozen": "100.25",
		"margin": "50.75",
		"transfer": "900.25",
		"positionMode": "HEDGE",
		"crossUnrealizedPNL": "25.50",
		"isolationUnrealizedPNL": "10.25",
		"bonus": "5.75"
	}`

	var entry AccountBalanceEntry
	err := json.Unmarshal([]byte(jsonData), &entry)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling account balance: %v", err)
	}

	if entry.MarginCoin != "USDT" {
		t.Errorf("unexpected marginCoin: %s", entry.MarginCoin)
	}

	if entry.Available != 1000.50 {
		t.Errorf("unexpected available amount: %f", entry.Available)
	}

	if entry.Frozen != 100.25 {
		t.Errorf("unexpected frozen amount: %f", entry.Frozen)
	}

	if entry.Margin != 50.75 {
		t.Errorf("unexpected margin amount: %f", entry.Margin)
	}

	if entry.Transfer != 900.25 {
		t.Errorf("unexpected transfer amount: %f", entry.Transfer)
	}

	if entry.PositionMode != TradePositionModeHedge {
		t.Errorf("unexpected position mode: %s", entry.PositionMode)
	}

	if entry.CrossUnrealizedPNL != 25.50 {
		t.Errorf("unexpected cross unrealized PNL: %f", entry.CrossUnrealizedPNL)
	}

	if entry.IsolationUnrealizedPNL != 10.25 {
		t.Errorf("unexpected isolation unrealized PNL: %f", entry.IsolationUnrealizedPNL)
	}

	if entry.Bonus != 5.75 {
		t.Errorf("unexpected bonus: %f", entry.Bonus)
	}

	// Test invalid JSON
	invalidJSON := `{
		"marginCoin": "USDT",
		"available": "not-a-number",
		"frozen": "0",
		"margin": "0",
		"transfer": "0",
		"positionMode": "HEDGE",
		"crossUnrealizedPNL": "0",
		"isolationUnrealizedPNL": "0",
		"bonus": "0"
	}`

	err = json.Unmarshal([]byte(invalidJSON), &entry)
	if err == nil {
		t.Fatal("expected error for invalid number format, got none")
	}
}
