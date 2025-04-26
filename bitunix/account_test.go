package bitunix

import (
	"context"
	"github.com/tradingiq/bitunix-client/model"
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

	params := model.AccountBalanceParams{
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

	if balance.PositionMode != model.PositionModeHedge {
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

	client, err := rest.New("http://example.com")
	if err != nil {
		t.Fatalf("failed to create restClient client: %v", err)
	}
	bitunixClient := New(client, "test-restClient-key", "test-restClient-secret")

	params := model.AccountBalanceParams{}
	_, err = bitunixClient.GetAccountBalance(context.Background(), params)

	if err == nil {
		t.Fatal("expected error for missing marginCoin, got none")
	}

	if err.Error() != "marginCoin is required" {
		t.Errorf("unexpected error message: %s", err.Error())
	}
}
