package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tradingiq/bitunix-client/bitunix"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/samples"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	client, _ := bitunix.NewApiClient(samples.Config.ApiKey, samples.Config.SecretKey)

	params := model.PendingTPSLOrderParams{
		Symbol:     model.ParseSymbol("BTCUSDT"),
		Limit:      10,
		PositionID: "5122906478792974765",
	}

	response, err := client.GetPendingTPSLOrder(context.Background(), params)
	if err != nil {
		log.Fatal("Failed to get pending TP/SL orders:", err)
	}

	if response.Code != 0 {
		log.Fatalf("API error: %s (code: %d)", response.Message, response.Code)
	}

	fmt.Printf("Found %d pending TP/SL orders:\n", len(response.Data))
	for i, order := range response.Data {
		fmt.Printf("\nOrder %d:\n", i+1)
		fmt.Printf("  Order ID: %s\n", order.ID)
		fmt.Printf("  Position ID: %s\n", order.PositionID)
		fmt.Printf("  Symbol: %s (%s/%s)\n", order.Symbol, order.Base, order.Quote)

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
			if order.TpQty != nil {
				fmt.Printf("    - Quantity: %.8f\n", *order.TpQty)
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
			if order.SlQty != nil {
				fmt.Printf("    - Quantity: %.8f\n", *order.SlQty)
			}
		}
	}
}
