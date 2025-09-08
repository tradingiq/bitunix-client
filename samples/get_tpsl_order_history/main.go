package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tradingiq/bitunix-client/samples"

	"github.com/tradingiq/bitunix-client/bitunix"
	"github.com/tradingiq/bitunix-client/model"
)

func main() {

	client, _ := bitunix.NewApiClient(samples.Config.ApiKey, samples.Config.SecretKey)

	params := model.TPSLOrderHistoryParams{
		Limit: 10,
	}

	response, err := client.GetTPSLOrderHistory(context.Background(), params)
	if err != nil {
		log.Fatalf("Failed to get TP/SL order history: %v", err)
	}

	fmt.Printf("Total TP/SL orders: %d\n", response.Data.Total)
	fmt.Println("\nTP/SL Order History:")
	fmt.Println("===================")

	for _, order := range response.Data.OrderList {
		fmt.Printf("\nOrder ID: %s\n", order.ID)
		fmt.Printf("Position ID: %s\n", order.PositionID)
		fmt.Printf("Symbol: %s\n", order.Symbol)
		fmt.Printf("Status: %s\n", order.Status)
		fmt.Printf("Created: %s\n", order.Ctime.Format(time.RFC3339))

		if order.TpPrice != nil {
			fmt.Printf("Take Profit Price: %.2f\n", *order.TpPrice)
			if order.TpStopType != nil {
				fmt.Printf("Take Profit Stop Type: %s\n", *order.TpStopType)
			}
			if order.TpOrderType != nil {
				fmt.Printf("Take Profit Order Type: %s\n", *order.TpOrderType)
			}
			if order.TpOrderPrice != nil {
				fmt.Printf("Take Profit Order Price: %.2f\n", *order.TpOrderPrice)
			}
			if order.TpQty != nil {
				fmt.Printf("Take Profit Quantity: %.4f\n", *order.TpQty)
			}
		}

		if order.SlPrice != nil {
			fmt.Printf("Stop Loss Price: %.2f\n", *order.SlPrice)
			if order.SlStopType != nil {
				fmt.Printf("Stop Loss Stop Type: %s\n", *order.SlStopType)
			}
			if order.SlOrderType != nil {
				fmt.Printf("Stop Loss Order Type: %s\n", *order.SlOrderType)
			}
			if order.SlOrderPrice != nil {
				fmt.Printf("Stop Loss Order Price: %.2f\n", *order.SlOrderPrice)
			}
			if order.SlQty != nil {
				fmt.Printf("Stop Loss Quantity: %.4f\n", *order.SlQty)
			}
		}

		if order.TriggerTime != nil {
			fmt.Printf("Triggered: %s\n", (*order.TriggerTime).Format(time.RFC3339))
		}
		fmt.Println("---")
	}
}
