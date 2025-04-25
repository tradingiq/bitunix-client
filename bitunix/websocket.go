package bitunix

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tradingiq/bitunix-client/security"
	"time"
)

type HeartbeatMessage struct {
	Op   string `json:"op"`
	Ping int64  `json:"ping"`
}

type SubscriptionMessage struct {
	Op   string               `json:"op"`
	Args []SubscriptionParams `json:"args"`
}

type SubscriptionParams struct {
	Ch string `json:"ch"`
}

type LoginMessage struct {
	Op   string        `json:"op"`
	Args []LoginParams `json:"args"`
}

type LoginParams struct {
	ApiKey    string `json:"apiKey"`
	Timestamp int64  `json:"timestamp"`
	Nonce     string `json:"nonce"`
	Sign      string `json:"sign"`
}

func generateWebsocketSignature(apiKey, apiSecret string, timestamp int64, nonceBytes []byte) (string, int64) {
	preSign := fmt.Sprintf("%x%d%s", nonceBytes, timestamp, apiKey)

	preSign = security.Sha256Hex(preSign)
	sign := security.Sha256Hex(preSign + apiSecret)

	return sign, timestamp
}

func KeepAliveMonitor() func() ([]byte, error) {
	return func() ([]byte, error) {
		msg := HeartbeatMessage{
			Op:   "ping",
			Ping: time.Now().Unix(),
		}

		bytes, err := json.Marshal(msg)
		if err != nil {
			return bytes, err
		}
		return bytes, err
	}
}

func WebsocketSigner(apiKey, apiSecret string) func() ([]byte, error) {
	return func() ([]byte, error) {

		nonce, err := security.GenerateNonce(32)
		if err != nil {
			return nil, fmt.Errorf("failed to generate nonce: %w", err)
		}

		sign, timestamp := generateWebsocketSignature(apiKey, apiSecret, time.Now().Unix(), nonce)

		loginReq := LoginMessage{
			Op: "login",
			Args: []LoginParams{
				{
					ApiKey:    apiKey, // Use the provided apiKey parameter
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

	}
}
