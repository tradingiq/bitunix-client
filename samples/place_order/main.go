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
		log.Fatalf("Failed to build order: %v", err)
	}

	ctx := context.Background()
	response, err := bitunixClient.PlaceOrder(ctx, &limitOrder)
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
		case errors.Is(err, bitunix_errors.ErrDuplicateClientID):
			log.Fatalf("Duplicated client id error: %s", err.Error())
		case errors.Is(err, bitunix_errors.ErrInsufficientBalance):
			log.Fatalf("Insufficient balance error: %s", err.Error())
		case errors.Is(err, bitunix_errors.ErrMarketNotExists):
			log.Fatalf("Market not exists error: %s", err.Error())
		case errors.Is(err, bitunix_errors.ErrTriggerPriceInvalid):
			log.Fatalf("Trigger price invalid error: %s", err.Error())
		case errors.Is(err, bitunix_errors.ErrOrderPriceIssue):
			log.Fatalf("Order Price Issue invalid error: %s", err.Error())
		case errors.Is(err, bitunix_errors.ErrOrderQuantityIssue):
			log.Fatalf("Order Quantity Issue invalid error: %s", err.Error())
		case errors.Is(err, bitunix_errors.ErrPositionLimitExceeded):
			log.Fatalf("Position limit exceeded invalid error: %s", err.Error())

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

	log.WithField("response", response).Info("placed order")
}
