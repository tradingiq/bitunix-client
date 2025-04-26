package bitunix

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/security"
	"github.com/tradingiq/bitunix-client/websocket"
	"sync"
	"time"
)

type websocketClient struct {
	client *websocket.Client
}

func (ws *websocketClient) Connect() error {
	if err := ws.client.Connect(); err != nil {
		return err
	}

	return nil
}

func (ws *websocketClient) Disconnect() {
	ws.client.Close()
}

type privateWebsocketClient struct {
	websocketClient
	balanceSubscribers map[chan model.BalanceResponse]struct{}
	subMtx             sync.Mutex
}

func NewPrivateWebsocket(ctx context.Context, apiKey, secretKey string) *privateWebsocketClient {
	ws := websocket.New(
		ctx,
		"wss://fapi.bitunix.com/private/",
		websocket.WithAuthentication(WebsocketSigner(apiKey, secretKey)),
		websocket.WithKeepAliveMonitor(30*time.Second, KeepAliveMonitor()),
	)

	return &privateWebsocketClient{
		websocketClient:    websocketClient{client: ws},
		subMtx:             sync.Mutex{},
		balanceSubscribers: map[chan model.BalanceResponse]struct{}{},
	}
}

func (ws *privateWebsocketClient) SubscribeBalance() <-chan model.BalanceResponse {
	ws.subMtx.Lock()
	defer ws.subMtx.Unlock()
	c := make(chan model.BalanceResponse)

	ws.balanceSubscribers[c] = struct{}{}

	return c
}

func (ws *privateWebsocketClient) UnsubscribeBalance(sub chan model.BalanceResponse) {
	ws.subMtx.Lock()
	defer ws.subMtx.Unlock()

	if _, ok := ws.balanceSubscribers[sub]; ok {
		delete(ws.balanceSubscribers, sub)
	}
}

func (ws *privateWebsocketClient) Stream() error {
	err := ws.client.Listen(func(bytes []byte) error {
		go func() {
			var result map[string]interface{}

			err := json.Unmarshal(bytes, &result)
			if err != nil {
				log.WithError(fmt.Errorf("error unmarshaling JSON: %v", err)).Errorf("error unmarshaling JSON")
			}

			if ch, ok := result["ch"].(string); ok {
				switch ch {
				case model.ChannelBalance:
					balance := model.BalanceResponse{}
					if err := json.Unmarshal(bytes, &balance); err != nil {
						log.WithError(fmt.Errorf("error unmarshaling balance response JSON: %v", err)).Errorf("error unmarshaling balance response  JSON")
					}

					for sub, _ := range ws.balanceSubscribers {
						sub <- balance
					}
				}

			}
		}()

		return nil
	})
	if err != nil {
		return fmt.Errorf("stream error: %v", err)
	}

	return nil
}

type heartbeatMessage struct {
	Op   string `json:"op"`
	Ping int64  `json:"ping"`
}

type SubscriptionParams struct {
	Ch string `json:"ch"`
}

type loginMessage struct {
	Op   string        `json:"op"`
	Args []loginParams `json:"args"`
}

type loginParams struct {
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
		msg := heartbeatMessage{
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

		loginReq := loginMessage{
			Op: "login",
			Args: []loginParams{
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
