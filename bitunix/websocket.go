package bitunix

import (
	"fmt"
	"github.com/tradingiq/bitunix-client/security"
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

func GenerateWebsocketSignature(apiKey, apiSecret string, timestamp int64, nonceBytes []byte) (string, int64) {
	preSign := fmt.Sprintf("%x%d%s", nonceBytes, timestamp, apiKey)

	preSign = security.Sha256Hex(preSign)
	sign := security.Sha256Hex(preSign + apiSecret)

	return sign, timestamp
}
