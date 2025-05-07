package bitunix

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	bitunix_errors "github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/rest"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClientTimeout(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		time.Sleep(500 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":0,"msg":"OK"}`))
	}))
	defer server.Close()

	rest.WithDefaultTimeout(40 * time.Millisecond)

	client, err := rest.New(server.URL, rest.WithRequestSigner(func(req *http.Request, body []byte) error {
		return nil
	}), rest.WithDefaultTimeout(40*time.Millisecond))
	assert.NoError(t, err)

	bitunixClient := &apiClient{
		restClient: client,
		baseURI:    server.URL,
	}

	_, err = bitunixClient.GetAccountBalance(context.Background(), model.AccountBalanceParams{
		MarginCoin: "USDT",
	})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, bitunix_errors.ErrTimeout) || errors.Is(err, bitunix_errors.ErrNetwork),
		"Expected timeout or network error, got: %v", err)
}

func TestContextTimeoutCancellation(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		time.Sleep(300 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":0,"msg":"OK"}`))
	}))
	defer server.Close()

	client, err := rest.New(server.URL, rest.WithRequestSigner(func(req *http.Request, body []byte) error {
		return nil
	}))
	assert.NoError(t, err)

	bitunixClient := &apiClient{
		restClient: client,
		baseURI:    server.URL,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err = bitunixClient.GetAccountBalance(ctx, model.AccountBalanceParams{
		MarginCoin: "USDT",
	})

	assert.Error(t, err)
	assert.True(t, errors.Is(err, bitunix_errors.ErrTimeout) || errors.Is(err, context.DeadlineExceeded),
		"Expected timeout or context deadline error, got: %v", err)
}

func TestMultipleTimeoutScenarios(t *testing.T) {
	testCases := []struct {
		name           string
		serverDelay    time.Duration
		clientTimeout  time.Duration
		contextTimeout time.Duration
		expectTimeout  bool
	}{
		{
			name:           "No Timeout",
			serverDelay:    25 * time.Millisecond,
			clientTimeout:  200 * time.Millisecond,
			contextTimeout: 300 * time.Millisecond,
			expectTimeout:  false,
		},
		{
			name:           "Client Timeout",
			serverDelay:    300 * time.Millisecond,
			clientTimeout:  100 * time.Millisecond,
			contextTimeout: 500 * time.Millisecond,
			expectTimeout:  true,
		},
		{
			name:           "Context Timeout",
			serverDelay:    300 * time.Millisecond,
			clientTimeout:  500 * time.Millisecond,
			contextTimeout: 100 * time.Millisecond,
			expectTimeout:  true,
		},
		{
			name:           "Both Timeouts (Client triggers first)",
			serverDelay:    500 * time.Millisecond,
			clientTimeout:  100 * time.Millisecond,
			contextTimeout: 200 * time.Millisecond,
			expectTimeout:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(tc.serverDelay)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"code":0,"msg":"OK"}`))
			}))
			defer server.Close()

			client, err := rest.New(server.URL, rest.WithRequestSigner(func(req *http.Request, body []byte) error {
				return nil
			}), rest.WithDefaultTimeout(50*time.Millisecond))
			assert.NoError(t, err)

			bitunixClient := &apiClient{
				restClient: client,
				baseURI:    server.URL,
			}

			ctx, cancel := context.WithTimeout(context.Background(), tc.contextTimeout)
			defer cancel()

			_, err = bitunixClient.GetAccountBalance(ctx, model.AccountBalanceParams{
				MarginCoin: "USDT",
			})

			if tc.expectTimeout {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, bitunix_errors.ErrTimeout) || errors.Is(err, bitunix_errors.ErrNetwork) || errors.Is(err, context.DeadlineExceeded),
					"Expected timeout, network error, or deadline exceeded, got: %v", err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestOrderRequestTimeout(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":0,"message":"success","data":{"orderId":"12345","clientId":"client123"}}`))
	}))
	defer server.Close()

	client, err := rest.New(server.URL, rest.WithRequestSigner(func(req *http.Request, body []byte) error {
		return nil
	}), rest.WithDefaultTimeout(50*time.Millisecond))
	assert.NoError(t, err)

	bitunixClient := &apiClient{
		restClient: client,
		baseURI:    server.URL,
	}

	qty := 1.0
	price := 50000.0
	orderReq := &model.OrderRequest{
		Symbol:    "BTCUSDT",
		TradeSide: model.TradeSideBuy,
		Price:     &price,
		Qty:       qty,
		Side:      model.SideOpen,
		OrderType: model.OrderTypeLimit,
		ClientID:  "client123",
	}

	_, err = bitunixClient.PlaceOrder(context.Background(), orderReq)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, bitunix_errors.ErrTimeout) || errors.Is(err, bitunix_errors.ErrNetwork),
		"Expected timeout or network error, got: %v", err)
}

func TestCancelOrderRequestTimeout(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":0,"msg":"Success","data":{"successList":[{"orderId":"11111"}],"failureList":[]}}`))
	}))
	defer server.Close()

	client, err := rest.New(server.URL, rest.WithRequestSigner(func(req *http.Request, body []byte) error {
		return nil
	}), rest.WithDefaultTimeout(50*time.Millisecond))
	assert.NoError(t, err)

	bitunixClient := &apiClient{
		restClient: client,
		baseURI:    server.URL,
	}

	cancelOrderReq := &model.CancelOrderRequest{
		Symbol: "BTCUSDT",
		OrderList: []model.CancelOrderParam{
			{
				OrderID: "11111",
			},
		},
	}

	_, err = bitunixClient.CancelOrders(context.Background(), cancelOrderReq)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, bitunix_errors.ErrTimeout) || errors.Is(err, bitunix_errors.ErrNetwork),
		"Expected timeout or network error, got: %v", err)
}
