package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/tradingiq/bitunix-client/bitunix"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/samples"
)

func main() {
	client, err := bitunix.NewApiClient(samples.Config.ApiKey, samples.Config.SecretKey)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get all pending positions
	params := model.PendingPositionParams{}

	// Optionally filter by symbol or position ID
	if len(os.Args) > 1 {
		symbol := os.Args[1]
		params.Symbol = model.ParseSymbol(symbol)
		fmt.Printf("Getting pending positions for symbol: %s\n", symbol)
	} else {
		fmt.Println("Getting all pending positions")
	}

	resp, err := client.GetPendingPositions(ctx, params)
	if err != nil {
		fmt.Printf("Error getting pending positions: %v\n", err)
		return
	}

	fmt.Printf("Found %d pending positions\n", len(resp.Data))

	for i, position := range resp.Data {
		fmt.Printf("\nPosition %d:\n", i+1)
		fmt.Printf("  Position ID: %s\n", position.PositionID)
		fmt.Printf("  Symbol: %s\n", position.Symbol)
		fmt.Printf("  Side: %s\n", position.Side)
		fmt.Printf("  Quantity: %.4f\n", position.Qty)
		fmt.Printf("  Entry Value: %.2f\n", position.EntryValue)
		fmt.Printf("  Average Open Price: %.2f\n", position.AvgOpenPrice)
		fmt.Printf("  Leverage: %d\n", position.Leverage)
		fmt.Printf("  Margin Mode: %s\n", position.MarginMode)
		fmt.Printf("  Position Mode: %s\n", position.PositionMode)
		fmt.Printf("  Margin: %.2f\n", position.Margin)
		fmt.Printf("  Unrealized PNL: %.2f\n", position.UnrealizedPNL)
		fmt.Printf("  Realized PNL: %.2f\n", position.RealizedPNL)

		if position.LiqPrice > 0 {
			fmt.Printf("  Liquidation Price: %.2f\n", position.LiqPrice)
		} else {
			fmt.Println("  Liquidation Price: Low risk (no liquidation price)")
		}

		fmt.Printf("  Margin Rate: %.2f%%\n", position.MarginRate*100)
		fmt.Printf("  Fees: %.2f\n", position.Fees)
		fmt.Printf("  Funding: %.2f\n", position.Funding)
		fmt.Printf("  Created: %s\n", position.CreateTime.Format("2006-01-02 15:04:05"))
		fmt.Printf("  Modified: %s\n", position.ModifyTime.Format("2006-01-02 15:04:05"))
	}
}
