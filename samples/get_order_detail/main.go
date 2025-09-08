package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/tradingiq/bitunix-client/bitunix"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/samples"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	client, err := bitunix.NewApiClient(samples.Config.ApiKey, samples.Config.SecretKey)
	if err != nil {
		logger.Fatal("Failed to create API client", zap.Error(err))
	}

	ctx := context.Background()

	// Step 1: Fetch order history first
	fmt.Println("Fetching order history...")
	fmt.Println("========================")

	// Configure order history parameters
	startTime := time.Now().Add(-24 * 7 * time.Hour)
	endTime := time.Now().Add(-24 * 4 * time.Hour)
	historyParams := model.OrderHistoryParams{
		Symbol:    model.ParseSymbol("BTCUSDT"),
		StartTime: &startTime,
		EndTime:   &endTime,
		Type:      model.OrderTypeMarket,
		Limit:     5,
	}

	// Get order history
	historyResponse, err := client.GetOrderHistory(ctx, historyParams)
	if err != nil {
		logger.Fatal("Error getting order history", zap.Error(err))
	}

	fmt.Printf("Found %d orders in history\n", len(historyResponse.Data.Orders))
	fmt.Printf("Total orders matching criteria: %s\n\n", historyResponse.Data.Total)

	// Step 2: Fetch details for each order
	fmt.Println("Fetching details for each order...")
	fmt.Println("=================================")

	for i, order := range historyResponse.Data.Orders {
		fmt.Printf("\nOrder %d/%d:\n", i+1, len(historyResponse.Data.Orders))
		fmt.Printf("Order ID: %s\n", order.OrderID)
		fmt.Printf("Symbol: %s\n", order.Symbol)
		fmt.Printf("Status: %s\n", order.Status)
		fmt.Printf("Created: %s\n", order.CreateTime.Format(time.RFC3339))

		// Fetch detailed order information
		detailRequest := &bitunix.OrderDetailRequest{
			OrderID: order.OrderID,
		}

		detailResponse, err := client.GetOrderDetail(ctx, detailRequest)
		if err != nil {
			logger.Error("Error getting details for order", zap.String("orderID", order.OrderID), zap.Error(err))
			continue
		}

		// Display detailed order information
		fmt.Println("\nDetailed Information:")
		fmt.Printf("  Order ID: %s\n", detailResponse.Data.OrderID)
		fmt.Printf("  Symbol: %s\n", detailResponse.Data.Symbol)
		fmt.Printf("  Quantity: %f\n", detailResponse.Data.Quantity)
		fmt.Printf("  Trade Quantity: %f\n", detailResponse.Data.TradeQuantity)
		fmt.Printf("  Price: %f\n", detailResponse.Data.Price)
		fmt.Printf("  Side: %s\n", detailResponse.Data.Side)
		fmt.Printf("  Status: %s\n", detailResponse.Data.Status)
		fmt.Printf("  Leverage: %d\n", detailResponse.Data.Leverage)
		fmt.Printf("  Order Type: %s\n", detailResponse.Data.OrderType)
		fmt.Printf("  Position Mode: %s\n", detailResponse.Data.PositionMode)
		fmt.Printf("  Margin Mode: %s\n", detailResponse.Data.MarginMode)

		if detailResponse.Data.ClientID != "" {
			fmt.Printf("  Client ID: %s\n", detailResponse.Data.ClientID)
		}

		fmt.Printf("  Fee: %f\n", detailResponse.Data.Fee)
		fmt.Printf("  Realized PNL: %f\n", detailResponse.Data.RealizedPNL)
		fmt.Printf("  Created: %s\n", detailResponse.Data.CreateTime.Format(time.RFC3339))
		fmt.Printf("  Modified: %s\n", detailResponse.Data.ModifyTime.Format(time.RFC3339))

		// Check if order has take profit or stop loss
		if detailResponse.Data.TpPrice != nil {
			fmt.Println("\n  Take Profit Settings:")
			fmt.Printf("    TP Price: %f\n", *detailResponse.Data.TpPrice)
			fmt.Printf("    TP Stop Type: %s\n", *detailResponse.Data.TpStopType)
			fmt.Printf("    TP Order Type: %s\n", *detailResponse.Data.TpOrderType)
			if detailResponse.Data.TpOrderPrice != nil {
				fmt.Printf("    TP Order Price: %f\n", *detailResponse.Data.TpOrderPrice)
			}
		}

		if detailResponse.Data.SlPrice != nil {
			fmt.Println("\n  Stop Loss Settings:")
			fmt.Printf("    SL Price: %f\n", *detailResponse.Data.SlPrice)
			fmt.Printf("    SL Stop Type: %s\n", *detailResponse.Data.SlStopType)
			fmt.Printf("    SL Order Type: %s\n", *detailResponse.Data.SlOrderType)
			if detailResponse.Data.SlOrderPrice != nil {
				fmt.Printf("    SL Order Price: %f\n", *detailResponse.Data.SlOrderPrice)
			}
		}

		fmt.Println("\n" + strings.Repeat("-", 50))
	}

	fmt.Println("\nOrder history and detail fetching completed!")
}
