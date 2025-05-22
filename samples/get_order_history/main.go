package main

import (
	"context"
	"errors"
	"github.com/tradingiq/bitunix-client/bitunix"
	bitunix_errors "github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/samples"
	"go.uber.org/zap"
	"time"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	client, _ := bitunix.NewApiClient(samples.Config.ApiKey, samples.Config.SecretKey, bitunix.WithLogLevel(model.LogLevelVeryAggressive))

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
			logger.Fatal("Authentication failed", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrSignatureError):
			logger.Fatal("Signature error", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrNetwork):
			logger.Fatal("Network error", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrIPNotAllowed):
			logger.Fatal("IP restriction error", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrRateLimitExceeded):
			logger.Fatal("Rate limit exceeded", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrParameterError):
			logger.Fatal("Parameter error", zap.Error(err))
		default:
			var apiErr *bitunix_errors.APIError
			if errors.As(err, &apiErr) {
				logger.Fatal("API error", zap.Int("code", apiErr.Code), zap.String("message", apiErr.Message))
			} else {
				logger.Fatal("Unexpected error", zap.Error(err))
			}
		}
	}

	logger.Info("got order history", zap.Any("response", response))
}
