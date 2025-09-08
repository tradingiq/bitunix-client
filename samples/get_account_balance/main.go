package main

import (
	"context"
	"errors"

	"github.com/tradingiq/bitunix-client/bitunix"
	bitunix_errors "github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/samples"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Example: Use a custom logger configuration
	// You can either use WithLogLevel for predefined log levels or WithLogger for custom loggers

	// Option 1: Use predefined log levels
	client, _ := bitunix.NewApiClient(samples.Config.ApiKey, samples.Config.SecretKey, bitunix.WithLogLevel(model.LogLevelVeryAggressive))

	// Option 2: Use custom logger (commented out for this example)
	// customLogger, _ := zap.NewProduction() // or any other zap logger configuration
	// client, _ := bitunix.NewApiClient(samples.Config.ApiKey, samples.Config.SecretKey, bitunix.WithLogger(customLogger))

	params := model.AccountBalanceParams{
		MarginCoin: "USDT",
	}

	ctx := context.Background()
	response, err := client.GetAccountBalance(ctx, params)
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

	logger.Info("get account balance", zap.Any("response", response))
}
