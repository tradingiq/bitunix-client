package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/tradingiq/bitunix-client/bitunix"
	bitunix_errors "github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/samples"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	client, _ := bitunix.NewApiClient(samples.Config.ApiKey, samples.Config.SecretKey, bitunix.WithLogLevel(model.LogLevelVeryAggressive))

	params := model.PendingOrderParams{
		Symbol: model.ParseSymbol("BTCUSDT"),
		Limit:  100,
	}

	ctx := context.Background()
	response, err := client.GetPendingOrder(ctx, params)
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

	fmt.Printf("\nTotal pending orders: %d\n", response.Data.Total)
	fmt.Printf("Showing %d orders:\n\n", len(response.Data.OrderList))

	for i, order := range response.Data.OrderList {
		fmt.Printf("Order %d:\n", i+1)
		fmt.Printf("  Order ID: %s\n", order.OrderID)
		if order.ClientID != "" {
			fmt.Printf("  Client ID: %s\n", order.ClientID)
		}
		fmt.Printf("  Symbol: %s\n", order.Symbol)
		fmt.Printf("  Side: %s\n", order.Side)
		fmt.Printf("  Type: %s\n", order.OrderType)
		fmt.Printf("  Status: %s\n", order.Status)
		fmt.Printf("  Price: %.8f\n", order.Price)
		fmt.Printf("  Quantity: %.8f\n", order.Quantity)
		fmt.Printf("  Filled Quantity: %.8f\n", order.TradeQuantity)
		fmt.Printf("  Position Mode: %s\n", order.PositionMode)
		fmt.Printf("  Margin Mode: %s\n", order.MarginMode)
		fmt.Printf("  Leverage: %dx\n", order.Leverage)
		fmt.Printf("  Time in Force: %s\n", order.Effect)
		fmt.Printf("  Reduce Only: %t\n", order.ReduceOnly)
		fmt.Printf("  Fee: %.8f\n", order.Fee)
		fmt.Printf("  Realized PnL: %.8f\n", order.RealizedPNL)
		fmt.Printf("  Created: %s\n", order.CreateTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Modified: %s\n", order.ModifyTime.Format("2006-01-02 15:04:05"))

		// Display TP/SL information if present
		if order.TpPrice != nil {
			fmt.Printf("  Take Profit:\n")
			fmt.Printf("    - Trigger Price: %.8f\n", *order.TpPrice)
			if order.TpStopType != nil {
				fmt.Printf("    - Stop Type: %s\n", *order.TpStopType)
			}
			if order.TpOrderType != nil {
				fmt.Printf("    - Order Type: %s\n", *order.TpOrderType)
			}
			if order.TpOrderPrice != nil {
				fmt.Printf("    - Order Price: %.8f\n", *order.TpOrderPrice)
			}
		}

		if order.SlPrice != nil {
			fmt.Printf("  Stop Loss:\n")
			fmt.Printf("    - Trigger Price: %.8f\n", *order.SlPrice)
			if order.SlStopType != nil {
				fmt.Printf("    - Stop Type: %s\n", *order.SlStopType)
			}
			if order.SlOrderType != nil {
				fmt.Printf("    - Order Type: %s\n", *order.SlOrderType)
			}
			if order.SlOrderPrice != nil {
				fmt.Printf("    - Order Price: %.8f\n", *order.SlOrderPrice)
			}
		}

		fmt.Println()
	}

	logger.Info("Successfully retrieved pending orders", zap.Int("count", len(response.Data.OrderList)))
}
