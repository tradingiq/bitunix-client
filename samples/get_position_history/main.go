package main

import (
	"context"
	"errors"
	"time"

	"github.com/tradingiq/bitunix-client/bitunix"
	bitunix_errors "github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/samples"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	bitunixClient, _ := bitunix.NewApiClient(samples.Config.ApiKey, samples.Config.SecretKey)

	params := model.PositionHistoryParams{
		Limit: 100,
	}

	ctx := context.Background()
	response, err := bitunixClient.GetPositionHistory(ctx, params)
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

	logger.Info("got HistoricalPosition History", zap.Any("response", response))

	var totalPnL float64
	for _, position := range response.Data.Positions {
		totalPnL += position.RealizedPNL
	}

	logger.Info("Summary",
		zap.Int("positions_count", len(response.Data.Positions)),
		zap.Float64("total_pnl", totalPnL),
	)

	time.Sleep(5000 * time.Second)
}
