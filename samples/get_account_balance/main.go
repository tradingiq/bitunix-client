package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/bitunix"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/rest"
	"github.com/tradingiq/bitunix-client/samples"
	"time"
)

func main() {
	bitunixExample()
}

func bitunixExample() {
	log.SetLevel(log.DebugLevel)

	apiClient, err := rest.New("https://fapi.bitunix.com/", rest.WithDebug(), rest.WithDefaultTimeout(30*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	client := bitunix.New(apiClient, samples.Config.ApiKey, samples.Config.SecretKey)

	params := model.AccountBalanceParams{
		MarginCoin: "USDT",
	}

	ctx := context.Background()
	response, err := client.GetAccountBalance(ctx, params)
	if err != nil {
		log.Fatal(err)
	}

	log.WithField("response", response).Debug("Get Account Balance")

	fmt.Printf("Account Balance Information:\n\n")
	fmt.Printf("  Margin Coin: %s\n", response.Data.MarginCoin)
	fmt.Printf("  Available: %.6f\n", response.Data.Available)
	fmt.Printf("  Frozen: %.6f\n", response.Data.Frozen)
	fmt.Printf("  Margin: %.6f\n", response.Data.Margin)
	fmt.Printf("  Transfer: %.6f\n", response.Data.Transfer)
	fmt.Printf("  Position Mode: %s\n", response.Data.PositionMode)
	fmt.Printf("  Cross Unrealized PNL: %.6f\n", response.Data.CrossUnrealizedPNL)
	fmt.Printf("  Isolation Unrealized PNL: %.6f\n", response.Data.IsolationUnrealizedPNL)
	fmt.Printf("  Bonus: %.6f\n", response.Data.Bonus)
	fmt.Println()
}
