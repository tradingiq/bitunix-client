package bitunix

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/security"
	"github.com/tradingiq/bitunix-client/websocket"
	"strconv"
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
	balanceSubscribers map[chan BalanceResponse]struct{}
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
		balanceSubscribers: map[chan BalanceResponse]struct{}{},
	}
}

type BalanceResponse struct {
	Ch   string        `json:"ch"`
	Ts   int64         `json:"ts"`
	Data BalanceDetail `json:"data"`
}

type BalanceDetail struct {
	Coin            string  `json:"coin"`
	Available       float64 `json:"-"`
	Frozen          float64 `json:"-"`
	IsolationFrozen float64 `json:"-"`
	CrossFrozen     float64 `json:"-"`
	Margin          float64 `json:"-"`
	IsolationMargin float64 `json:"-"`
	CrossMargin     float64 `json:"-"`
	ExpMoney        float64 `json:"-"`
}

func (p *BalanceDetail) UnmarshalJSON(data []byte) error {
	type Alias BalanceDetail
	aux := &struct {
		Available       string `json:"available"`
		Frozen          string `json:"frozen"`
		IsolationFrozen string `json:"isolationFrozen"`
		CrossFrozen     string `json:"crossFrozen"`
		Margin          string `json:"margin"`
		IsolationMargin string `json:"isolationMargin"`
		CrossMargin     string `json:"crossMargin"`
		ExpMoney        string `json:"expMoney"`
		*Alias
	}{
		Alias: (*Alias)(p),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	available, err := strconv.ParseFloat(aux.Available, 64)
	if err == nil {
		p.Available = available
	} else {
		return fmt.Errorf("failed to parse available: %w", err)
	}

	frozen, err := strconv.ParseFloat(aux.Frozen, 64)
	if err == nil {
		p.Frozen = frozen
	} else {
		return fmt.Errorf("failed to parse frozen: %w", err)
	}

	isolationFrozen, err := strconv.ParseFloat(aux.IsolationFrozen, 64)
	if err == nil {
		p.IsolationFrozen = isolationFrozen
	} else {
		return fmt.Errorf("failed to parse IisolationFrozen: %w", err)
	}

	crossFrozen, err := strconv.ParseFloat(aux.CrossFrozen, 64)
	if err == nil {
		p.CrossFrozen = crossFrozen
	} else {
		return fmt.Errorf("failed to parse crossFrozen: %w", err)
	}

	Margin, err := strconv.ParseFloat(aux.Margin, 64)
	if err == nil {
		p.Margin = Margin
	} else {
		return fmt.Errorf("failed to parse Margin: %w", err)
	}

	isolationMargin, err := strconv.ParseFloat(aux.IsolationMargin, 64)
	if err == nil {
		p.IsolationMargin = isolationMargin
	} else {
		return fmt.Errorf("failed to parse isolationMargin: %w", err)
	}

	crossMargin, err := strconv.ParseFloat(aux.CrossMargin, 64)
	if err == nil {
		p.CrossMargin = crossMargin
	} else {
		return fmt.Errorf("failed to parse crossMargin: %w", err)
	}

	expMoney, err := strconv.ParseFloat(aux.ExpMoney, 64)
	if err == nil {
		p.ExpMoney = expMoney
	} else {
		return fmt.Errorf("failed to parse expMoney: %w", err)
	}

	return nil
}

func (ws *privateWebsocketClient) SubscribeBalance() <-chan BalanceResponse {
	ws.subMtx.Lock()
	defer ws.subMtx.Unlock()
	c := make(chan BalanceResponse)

	ws.balanceSubscribers[c] = struct{}{}

	return c
}

func (ws *privateWebsocketClient) UnsubscribeBalance(sub chan BalanceResponse) {
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
				case ChannelBalance:
					balance := BalanceResponse{}
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

type SubscriptionMessage struct {
	Op   string               `json:"op"`
	Args []SubscriptionParams `json:"args"`
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
