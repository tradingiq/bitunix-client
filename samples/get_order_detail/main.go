package main

import (
	"context"
	"fmt"
	"github.com/tradingiq/bitunix-client/bitunix"
	"github.com/tradingiq/bitunix-client/samples"
	"log"
)

func main() {
	client, err := bitunix.NewApiClient(samples.Config.ApiKey, samples.Config.SecretKey)
	if err != nil {
		log.Fatalf("Failed to create API client: %v", err)
	}

	// Example 1: Get order detail by order ID
	request := &bitunix.OrderDetailRequest{
		OrderID: "1923822405555130369",
	}

	response, err := client.GetOrderDetail(context.Background(), request)
	if err != nil {
		log.Printf("Error getting order detail by order ID: %v", err)
	} else {
		fmt.Printf("Order Detail:\n")
		fmt.Printf("Order ID: %s\n", response.Data.OrderID)
		fmt.Printf("Symbol: %s\n", response.Data.Symbol)
		fmt.Printf("Quantity: %f\n", response.Data.Quantity)
		fmt.Printf("Trade Quantity: %f\n", response.Data.TradeQuantity)
		fmt.Printf("Price: %f\n", *response.Data.Price)
		fmt.Printf("Side: %s\n", response.Data.Side)
		fmt.Printf("Status: %s\n", response.Data.Status)
		fmt.Printf("Leverage: %d\n", response.Data.Leverage)
		fmt.Printf("Order Type: %s\n", response.Data.OrderType)

		if response.Data.ClientID != "" {
			fmt.Printf("Client ID: %s\n", response.Data.ClientID)
		}

		fmt.Printf("Fee: %f\n", response.Data.Fee)
		fmt.Printf("Realized PNL: %f\n", response.Data.RealizedPNL)
		fmt.Printf("\n")
	}

	// Example 2: Get order detail by client ID
	request2 := &bitunix.OrderDetailRequest{
		ClientID: "my-client-order-123",
	}

	response2, err := client.GetOrderDetail(context.Background(), request2)
	if err != nil {
		log.Printf("Error getting order detail by client ID: %v", err)
	} else {
		fmt.Printf("Order Detail by Client ID:\n")
		fmt.Printf("Order ID: %s\n", response2.Data.OrderID)
		fmt.Printf("Client ID: %s\n", response2.Data.ClientID)
		fmt.Printf("Symbol: %s\n", response2.Data.Symbol)
		fmt.Printf("Status: %s\n", response2.Data.Status)
		fmt.Printf("\n")
	}

	// Example 3: Using both order ID and client ID
	request3 := &bitunix.OrderDetailRequest{
		OrderID:  "123456789",
		ClientID: "my-client-order-123",
	}

	response3, err := client.GetOrderDetail(context.Background(), request3)
	if err != nil {
		log.Printf("Error getting order detail with both IDs: %v", err)
	} else {
		fmt.Printf("Order Detail with both IDs:\n")
		fmt.Printf("Order ID: %s\n", response3.Data.OrderID)
		fmt.Printf("Client ID: %s\n", response3.Data.ClientID)
		fmt.Printf("Symbol: %s\n", response3.Data.Symbol)
		fmt.Printf("\n")

		// Check if order has take profit or stop loss
		if response3.Data.TpPrice != nil {
			fmt.Printf("Take Profit Price: %f\n", *response3.Data.TpPrice)
			fmt.Printf("Take Profit Stop Type: %s\n", *response3.Data.TpStopType)
			fmt.Printf("Take Profit Order Type: %s\n", *response3.Data.TpOrderType)
		}

		if response3.Data.SlPrice != nil {
			fmt.Printf("Stop Loss Price: %f\n", *response3.Data.SlPrice)
			fmt.Printf("Stop Loss Stop Type: %s\n", *response3.Data.SlStopType)
			fmt.Printf("Stop Loss Order Type: %s\n", *response3.Data.SlOrderType)
		}
	}
}
