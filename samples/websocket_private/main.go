package main

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/bitunix"
	"github.com/tradingiq/bitunix-client/samples"
	"github.com/tradingiq/bitunix-client/security"
	"github.com/tradingiq/bitunix-client/websocket"
	"time"
)

func main() {
	ctx := context.Background()
	ws := websocket.New(
		ctx,
		"wss://fapi.bitunix.com/private/",
		websocket.WithAuthentication(func() ([]byte, error) {
			nonce, err := security.GenerateNonce(32)
			if err != nil {
				return nil, fmt.Errorf("failed to generate nonce: %w", err)
			}

			sign, timestamp := bitunix.GenerateWebsocketSignature(samples.Config.ApiKey, samples.Config.SecretKey, time.Now().Unix(), nonce)

			loginReq := bitunix.LoginMessage{
				Op: "login",
				Args: []bitunix.LoginParams{
					{
						ApiKey:    samples.Config.ApiKey,
						Timestamp: timestamp,
						Nonce:     hex.EncodeToString(nonce),
						Sign:      sign,
					},
				},
			}
			bytes, err := json.Marshal(loginReq)
			if err != nil {
				return bytes, err
			}

			return bytes, nil
		}),
		websocket.WithHeartbeat(15*time.Second, func() ([]byte, error) {
			msg := bitunix.HeartbeatMessage{
				Op:   "ping",
				Ping: time.Now().Unix(),
			}

			bytes, err := json.Marshal(msg)
			if err != nil {
				return bytes, err
			}
			return bytes, err
		}),
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

	// Listen for messages until an error occurs or connection is closed
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
