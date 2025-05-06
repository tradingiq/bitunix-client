package main

import (
	"context"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/bitunix"
	bitunix_errors "github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/samples"
	"time"
)

func main() {
	log.SetLevel(log.DebugLevel)

	client, _ := bitunix.NewApiClient(samples.Config.ApiKey, samples.Config.SecretKey)

	startTime := time.Now().Add(-80 * time.Hour)

	params := model.OrderHistoryParams{
		Symbol:    model.ParseSymbol("BTCUSDT"),
		StartTime: &startTime,
		Limit:     10,
	}

	ctx := context.Background()
	response, err := client.GetOrderHistory(ctx, params)
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
		default:
			var apiErr *bitunix_errors.APIError
			if errors.As(err, &apiErr) {
				log.Fatalf("API error (code: %d): %s", apiErr.Code, apiErr.Message)
			} else {
				log.Fatalf("Unexpected error: %s", err.Error())
			}
		}
	}

	log.WithField("response", response).Info("got order history")
}
