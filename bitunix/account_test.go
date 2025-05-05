package bitunix

import (
	"context"
	stderrors "errors"
	"github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
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

	bitunixClient, _ := NewApiClient("test-restClient-key", "test-restClient-secret", WithBaseURI(server.URL))

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
	bitunixClient, _ := NewApiClient("test-restClient-key", "test-restClient-secret", WithBaseURI("http://example.com"))

	params := model.AccountBalanceParams{}
	_, err := bitunixClient.GetAccountBalance(context.Background(), params)

	if err == nil {
		t.Fatal("expected error for missing marginCoin, got none")
	}

	expected := "validation error: field marginCoin: is required"
	if err.Error() != expected {
		t.Errorf("unexpected error message: %s", err.Error())
	}

	if !stderrors.Is(err, errors.ErrValidation) {
		t.Errorf("error should be a validation error")
	}
}
