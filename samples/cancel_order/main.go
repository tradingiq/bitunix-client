package main

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/bitunix"
	bitunix_errors "github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/samples"
)

func main() {
	log.SetLevel(log.DebugLevel)

	bitunixClient, _ := bitunix.NewApiClient(samples.Config.ApiKey, samples.Config.SecretKey)

	cancelRequestBuilder := bitunix.NewCancelOrderBuilder(model.ParseSymbol("BTCUSDT")).
		WithOrderID("1915122868439269376")

	cancelRequest, err := cancelRequestBuilder.Build()
	if err != nil {
		log.Fatalf("Failed to build cancel order request: %v", err)
	}

	ctx := context.Background()
	response, err := bitunixClient.CancelOrders(ctx, &cancelRequest)
	if err != nil {
		switch {
		case errors.Is(err, bitunix_errors.ErrAuthentication):
			log.Fatalf("Authentication failed: %s", err.Error())
		case errors.Is(err, bitunix_errors.ErrSignatureError):
			log.Fatalf("Signature error: %s", err.Error())
		case errors.Is(err, bitunix_errors.ErrNetwork):
			log.Fatalf("Network error: %s", err.Error())
		case errors.Is(err, bitunix_errors.ErrIPNotAllowed):
			log.Fatalf("IP restriction error: %s", err.Error())
		case errors.Is(err, bitunix_errors.ErrRateLimitExceeded):
			log.Fatalf("Rate limit exceeded: %s", err.Error())
		case errors.Is(err, bitunix_errors.ErrParameterError):
			log.Fatalf("Parameter error: %s", err.Error())
		case errors.Is(err, bitunix_errors.ErrAccountNotAllowed):
			log.Fatalf("Account not allowed error: %s", err.Error())
		case errors.Is(err, bitunix_errors.ErrOrderNotFound):
			log.Fatalf("Order not found error: %s", err.Error())

			// and so on...
		default:
			var apiErr *bitunix_errors.APIError
			if errors.As(err, &apiErr) {
				log.Fatalf("API error (code: %d): %s", apiErr.Code, apiErr.Message)
			} else {
				log.Fatalf("Unexpected error: %s", err.Error())
			}
		}
	}

	log.WithField("response", response).Info("got cancellation response")
}
