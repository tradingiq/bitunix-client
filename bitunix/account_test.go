package bitunix

import (
	"context"
	stderrors "errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
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

func TestGetAccountBalance_APIErrors(t *testing.T) {
	testCases := []struct {
		name           string
		mockResponse   string
		responseStatus int
		expectError    bool
		errorCode      int
	}{
		{
			name: "Insufficient Permissions",
			mockResponse: `{
				"code": 10003,
				"message": "Unauthorized access"
			}`,
			responseStatus: http.StatusUnauthorized,
			expectError:    true,
			errorCode:      10003,
		},
		{
			name: "Invalid Margin Coin",
			mockResponse: `{
				"code": 10002,
				"message": "Invalid margin coin"
			}`,
			responseStatus: http.StatusOK,
			expectError:    true,
			errorCode:      10002,
		},
		{
			name: "Rate Limit Exceeded",
			mockResponse: `{
				"code": 10005,
				"message": "Rate limit exceeded"
			}`,
			responseStatus: http.StatusTooManyRequests,
			expectError:    true,
			errorCode:      10005,
		},
		{
			name: "Network Error",
			mockResponse: `{
				"code": 10001,
				"message": "Network error"
			}`,
			responseStatus: http.StatusInternalServerError,
			expectError:    true,
			errorCode:      10001,
		},
		{
			name:           "Invalid JSON Response",
			mockResponse:   `{ invalid json }`,
			responseStatus: http.StatusOK,
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.responseStatus)
				_, _ = w.Write([]byte(tc.mockResponse))
			}))
			defer server.Close()

			bitunixClient, _ := NewApiClient("test-restClient-key", "test-restClient-secret", WithBaseURI(server.URL))

			params := model.AccountBalanceParams{
				MarginCoin: "USDT",
			}

			response, err := bitunixClient.GetAccountBalance(context.Background(), params)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}

				if tc.errorCode > 0 {
					apiErr, ok := err.(*errors.APIError)
					if !ok {
						t.Errorf("Expected *errors.APIError, got %T", err)
						return
					}

					if apiErr.Code != tc.errorCode {
						t.Errorf("Expected error code %d, got %d", tc.errorCode, apiErr.Code)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
					return
				}

				if response == nil {
					t.Errorf("Expected response but got nil")
				}
			}
		})
	}
}

func TestGetAccountBalance_RequestError(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Skip("Hijacking not supported")
			return
		}
		conn, _, err := hj.Hijack()
		if err != nil {
			t.Errorf("Failed to hijack connection: %v", err)
			return
		}
		conn.Close()
	}))
	defer server.Close()

	bitunixClient, _ := NewApiClient("test-restClient-key", "test-restClient-secret", WithBaseURI(server.URL))

	params := model.AccountBalanceParams{
		MarginCoin: "USDT",
	}

	_, err := bitunixClient.GetAccountBalance(context.Background(), params)
	if err == nil {
		t.Errorf("Expected network error, got nil")
	}
}

func TestGetAccountBalance_InvalidUnmarshal(t *testing.T) {
	testCases := []struct {
		name         string
		mockResponse string
		fieldToTest  string
		expectError  bool
	}{
		{
			name: "Invalid Available Amount",
			mockResponse: `{
				"code": 0,
				"data": {
					"marginCoin": "USDT",
					"available": "not_a_number",
					"frozen": "0",
					"margin": "10",
					"transfer": "1000",
					"positionMode": "HEDGE",
					"crossUnrealizedPNL": "2",
					"isolationUnrealizedPNL": "0",
					"bonus": "0"
				},
				"msg": "Success"
			}`,
			fieldToTest: "available",
			expectError: true,
		},
		{
			name: "Invalid Frozen Amount",
			mockResponse: `{
				"code": 0,
				"data": {
					"marginCoin": "USDT",
					"available": "1000",
					"frozen": "not_a_number",
					"margin": "10",
					"transfer": "1000",
					"positionMode": "HEDGE",
					"crossUnrealizedPNL": "2",
					"isolationUnrealizedPNL": "0",
					"bonus": "0"
				},
				"msg": "Success"
			}`,
			fieldToTest: "frozen",
			expectError: true,
		},
		{
			name: "Invalid Margin Amount",
			mockResponse: `{
				"code": 0,
				"data": {
					"marginCoin": "USDT",
					"available": "1000",
					"frozen": "0",
					"margin": "not_a_number",
					"transfer": "1000",
					"positionMode": "HEDGE",
					"crossUnrealizedPNL": "2",
					"isolationUnrealizedPNL": "0",
					"bonus": "0"
				},
				"msg": "Success"
			}`,
			fieldToTest: "margin",
			expectError: true,
		},
		{
			name: "Invalid Transfer Amount",
			mockResponse: `{
				"code": 0,
				"data": {
					"marginCoin": "USDT",
					"available": "1000",
					"frozen": "0",
					"margin": "10",
					"transfer": "not_a_number",
					"positionMode": "HEDGE",
					"crossUnrealizedPNL": "2",
					"isolationUnrealizedPNL": "0",
					"bonus": "0"
				},
				"msg": "Success"
			}`,
			fieldToTest: "transfer",
			expectError: true,
		},
		{
			name: "Invalid CrossUnrealizedPNL Amount",
			mockResponse: `{
				"code": 0,
				"data": {
					"marginCoin": "USDT",
					"available": "1000",
					"frozen": "0",
					"margin": "10",
					"transfer": "1000",
					"positionMode": "HEDGE",
					"crossUnrealizedPNL": "not_a_number",
					"isolationUnrealizedPNL": "0",
					"bonus": "0"
				},
				"msg": "Success"
			}`,
			fieldToTest: "crossUnrealizedPNL",
			expectError: true,
		},
		{
			name: "Invalid IsolationUnrealizedPNL Amount",
			mockResponse: `{
				"code": 0,
				"data": {
					"marginCoin": "USDT",
					"available": "1000",
					"frozen": "0",
					"margin": "10",
					"transfer": "1000",
					"positionMode": "HEDGE",
					"crossUnrealizedPNL": "2",
					"isolationUnrealizedPNL": "not_a_number",
					"bonus": "0"
				},
				"msg": "Success"
			}`,
			fieldToTest: "isolationUnrealizedPNL",
			expectError: true,
		},
		{
			name: "Invalid Bonus Amount",
			mockResponse: `{
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
					"bonus": "not_a_number"
				},
				"msg": "Success"
			}`,
			fieldToTest: "bonus",
			expectError: true,
		},
		{
			name: "Invalid Position Mode",
			mockResponse: `{
				"code": 0,
				"data": {
					"marginCoin": "USDT",
					"available": "1000",
					"frozen": "0",
					"margin": "10",
					"transfer": "1000",
					"positionMode": "INVALID_MODE",
					"crossUnrealizedPNL": "2",
					"isolationUnrealizedPNL": "0",
					"bonus": "0"
				},
				"msg": "Success"
			}`,
			fieldToTest: "positionMode",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(tc.mockResponse))
			}))
			defer server.Close()

			bitunixClient, _ := NewApiClient("test-restClient-key", "test-restClient-secret", WithBaseURI(server.URL))

			params := model.AccountBalanceParams{
				MarginCoin: "USDT",
			}

			_, err := bitunixClient.GetAccountBalance(context.Background(), params)

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for invalid %s, got nil", tc.fieldToTest)
					return
				}

				if !stderrors.Is(err, errors.ErrInternal) {
					t.Errorf("Expected internal error, got: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

func TestGetAccountBalance_AlternatePositionMode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `{
			"code": 0,
			"data": {
				"marginCoin": "BTC",
				"available": "10.5",
				"frozen": "2.3",
				"margin": "1.2",
				"transfer": "15.7",
				"positionMode": "ONE_WAY",
				"crossUnrealizedPNL": "0.5",
				"isolationUnrealizedPNL": "0.3",
				"bonus": "1.1"
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
		MarginCoin: "BTC",
	}

	response, err := bitunixClient.GetAccountBalance(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	balance := response.Data
	if balance.MarginCoin != "BTC" {
		t.Errorf("unexpected marginCoin: %s", balance.MarginCoin)
	}

	if balance.Available != 10.5 {
		t.Errorf("unexpected available amount: %f", balance.Available)
	}

	if balance.Frozen != 2.3 {
		t.Errorf("unexpected frozen amount: %f", balance.Frozen)
	}

	if balance.PositionMode != model.PositionModeOneWay {
		t.Errorf("unexpected position mode: %s", balance.PositionMode)
	}

	if balance.CrossUnrealizedPNL != 0.5 {
		t.Errorf("unexpected cross unrealized PNL: %f", balance.CrossUnrealizedPNL)
	}

	if balance.Bonus != 1.1 {
		t.Errorf("unexpected bonus: %f", balance.Bonus)
	}
}
