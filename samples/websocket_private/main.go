package main

import (
	"bitunix-client/samples"
	"bitunix-client/security"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/coder/websocket/wsjson"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"

	"github.com/coder/websocket"
)

type LoginRequest struct {
	Op   string        `json:"op"`
	Args []LoginParams `json:"args"`
}

type LoginParams struct {
	ApiKey    string `json:"apiKey"`
	Timestamp int64  `json:"timestamp"`
	Nonce     string `json:"nonce"`
	Sign      string `json:"sign"`
}

type SubscriptionRequest struct {
	Op   string               `json:"op"`
	Args []SubscriptionParams `json:"args"`
}

type SubscriptionParams struct {
	Ch string `json:"ch"`
}

type HeartbeatMessage struct {
	Op   string `json:"op"`
	Ping int64  `json:"ping"`
}

type GenericMessage map[string]interface{}

func generateSignature(apiKey, secretKey string, nonce []byte) (string, int64) {
	timestamp := time.Now().Unix()
	preSign := fmt.Sprintf("%x%d%s", nonce, timestamp, apiKey)

	preSign = security.Sha256Hex(preSign)
	sign := security.Sha256Hex(preSign + secretKey)

	return sign, timestamp
}

func sendHeartbeat(ctx context.Context, conn *websocket.Conn, done chan struct{}) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			heartbeat := HeartbeatMessage{
				Op:   "ping",
				Ping: time.Now().Unix(),
			}

			err := wsjson.Write(ctx, conn, heartbeat)
			if err != nil {
				log.Printf("Error sending heartbeat: %v", err)
				return
			}
			log.Println("Heartbeat sent")

		case <-done:
			return
		}
	}
}

func monitorBalance(apiKey, secretKey string) {
	wsURL := "wss://fapi.bitunix.com/private/"
	u, err := url.Parse(wsURL)
	if err != nil {
		log.Fatalf("Error parsing WebSocket URL: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*10)
	defer cancel()

	log.Printf("Connecting to %s", u.String())
	conn, _, err := websocket.Dial(ctx, u.String(), &websocket.DialOptions{
		HTTPClient: http.DefaultClient,
	})
	if err != nil {
		log.Fatalf("Error connecting to WebSocket: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	var initialMsg GenericMessage
	err = wsjson.Read(ctx, conn, &initialMsg)
	if err != nil {
		log.Fatalf("Error reading initial message: %v", err)
	}
	log.Printf("Initial message: %v", initialMsg)

	nonce, _ := security.GenerateNonce(32)
	sign, timestamp := generateSignature(apiKey, secretKey, nonce)

	loginReq := LoginRequest{
		Op: "login",
		Args: []LoginParams{
			{
				ApiKey:    apiKey,
				Timestamp: timestamp,
				Nonce:     hex.EncodeToString(nonce),
				Sign:      sign,
			},
		},
	}

	err = wsjson.Write(ctx, conn, loginReq)
	if err != nil {
		log.Fatalf("Error sending login request: %v", err)
	}

	var loginResp GenericMessage
	err = wsjson.Read(ctx, conn, &loginResp)
	if err != nil {
		log.Fatalf("Error reading login response: %v", err)
	}
	log.Printf("Login response: %v", loginResp)

	done := make(chan struct{})
	go sendHeartbeat(ctx, conn, done)

	subReq := SubscriptionRequest{
		Op: "subscribe",
		Args: []SubscriptionParams{
			{
				Ch: "position",
			},
		},
	}

	err = wsjson.Write(ctx, conn, subReq)
	if err != nil {
		log.Fatalf("Error sending subscription request: %v", err)
	}

	var subResp GenericMessage
	err = wsjson.Read(ctx, conn, &subResp)
	if err != nil {
		log.Fatalf("Error reading subscription response: %v", err)
	}
	log.Printf("Subscription response: %v", subResp)

	for {
		var message json.RawMessage
		err = wsjson.Read(ctx, conn, &message)
		if err != nil {
			log.Printf("Connection closed: %v", err)
			break
		}

		log.Printf("Balance update: %s", message)
	}

	close(done)
}

func main() {
	monitorBalance(samples.Config.ApiKey, samples.Config.SecretKey)
}
