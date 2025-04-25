package main

import (
	"context"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/bitunix"
	"github.com/tradingiq/bitunix-client/websocket"
	"time"
)

func main() {
	ctx := context.Background()
	ws := websocket.New(
		ctx,
		"wss://fapi.bitunix.com/private/",
		websocket.WithAuthentication(bitunix.WebsocketSigner()),
		websocket.WithKeepAliveMonitor(15*time.Second, bitunix.KeepAliveMonitor()),
	)

	if err := ws.Connect(); err != nil {
		log.Fatalf("failed to connect to WebSocket: %v", err)
	}
	defer ws.Close()

	subReq := bitunix.SubscriptionMessage{
		Op: "subscribe",
		Args: []bitunix.SubscriptionParams{
			{
				Ch: "position",
			},
		},
	}
	bytes, err := json.Marshal(subReq)
	if err != nil {
		log.Fatalf("failed to marshal subscribe to position updates: %v", err)
	}

	if err := ws.Write(bytes); err != nil {
		log.Fatalf("failed to subscribe to position updates: %v", err)
	}

	if err := ws.Listen(func(data []byte) error {
		var message map[string]interface{}
		if err := json.Unmarshal(data, &message); err != nil {
			log.Printf("Error parsing message: %v", err)
			return err
		}

		return nil
	}); err != nil {
		log.Printf("WebSocket connection closed: %v", err)
	}
}
