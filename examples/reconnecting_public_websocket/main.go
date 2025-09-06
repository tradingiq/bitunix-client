package main

import (
	"context"
	"log"
	"time"

	"github.com/tradingiq/bitunix-client/bitunix"
	"github.com/tradingiq/bitunix-client/model"
	"go.uber.org/zap"
)

type ExampleKLineSubscriber struct {
	symbol    model.Symbol
	interval  model.Interval
	priceType model.PriceType
}

func (s *ExampleKLineSubscriber) SubscribeKLine(msg *model.KLineChannelMessage) {
	log.Printf("Received KLine: %+v", msg)
}

func (s *ExampleKLineSubscriber) SubscribeSymbol() model.Symbol {
	return s.symbol
}

func (s *ExampleKLineSubscriber) SubscribeInterval() model.Interval {
	return s.interval
}

func (s *ExampleKLineSubscriber) SubscribePriceType() model.PriceType {
	return s.priceType
}

func main() {
	ctx := context.Background()
	
	// Create logger
	logger, _ := zap.NewDevelopment()
	
	// Create base websocket client
	baseClient, err := bitunix.NewPublicWebsocket(ctx, 
		bitunix.WithWebsocketLogger(logger),
	)
	if err != nil {
		log.Fatalf("Failed to create base client: %v", err)
	}
	
	// Create reconnecting websocket client
	reconnectingClient := bitunix.NewReconnectingPublicWebsocket(ctx, baseClient,
		bitunix.WithMaxReconnectAttempts(0), // Infinite reconnect attempts
		bitunix.WithReconnectDelay(10*time.Second),
		bitunix.WithReconnectLogger(logger),
	)
	
	// Connect
	if err := reconnectingClient.Connect(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	
	// Create subscriber
	subscriber := &ExampleKLineSubscriber{
		symbol:    model.ParseSymbol("BTCUSDT"),
		interval:  model.Interval1Min,
		priceType: model.PriceTypeMark,
	}
	
	// Subscribe to KLine
	if err := reconnectingClient.SubscribeKLine(subscriber); err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}
	
	// Start streaming (this will automatically reconnect on failures)
	go func() {
		if err := reconnectingClient.Stream(); err != nil {
			log.Printf("Stream ended with error: %v", err)
		}
	}()
	
	// Keep running for demo
	select {}
}