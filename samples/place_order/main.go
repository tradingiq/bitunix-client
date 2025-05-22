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

	bitunixClient, _ := bitunix.NewApiClient(samples.Config.ApiKey, samples.Config.SecretKey)

	orderBuilder := bitunix.NewOrderBuilder(
		model.ParseSymbol("BTCUSDT"),
		model.TradeSideSell,
		model.SideOpen,
		0.002,
	).WithOrderType(model.OrderTypeLimit).
		WithPrice(100000.0).
		WithTimeInForce(model.TimeInForcePostOnly)

	limitOrder, err := orderBuilder.Build()
	if err != nil {
		logger.Fatal("Failed to build order", zap.Error(err))
	}

	ctx := context.Background()
	response, err := bitunixClient.PlaceOrder(ctx, &limitOrder)
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
		case errors.Is(err, bitunix_errors.ErrAccountNotAllowed):
			logger.Fatal("Account not allowed error", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrDuplicateClientID):
			logger.Fatal("Duplicated client id error", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrInsufficientBalance):
			logger.Fatal("Insufficient balance error", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrMarketNotExists):
			logger.Fatal("Market not exists error", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrTriggerPriceInvalid):
			logger.Fatal("Trigger price invalid error", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrOrderPriceIssue):
			logger.Fatal("Order Price Issue invalid error", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrOrderQuantityIssue):
			logger.Fatal("Order Quantity Issue invalid error", zap.Error(err))
		case errors.Is(err, bitunix_errors.ErrPositionLimitExceeded):
			logger.Fatal("Position limit exceeded invalid error", zap.Error(err))

			// and so on...
		default:
			var apiErr *bitunix_errors.APIError
			if errors.As(err, &apiErr) {
				logger.Fatal("API error", zap.Int("code", apiErr.Code), zap.String("message", apiErr.Message))
			} else {
				logger.Fatal("Unexpected error", zap.Error(err))
			}
		}
	}

	logger.Info("placed order", zap.Any("response", response))
}
