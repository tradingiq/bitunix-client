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
	balanceSubscribers     map[chan model.BalanceChannelMessage]struct{}
	positionSubscribers    map[chan model.PositionChannelMessage]struct{}
	orderSubscribers       map[chan model.OrderChannelMessage]struct{}
	tpSlOrderSubscribers   map[chan model.TpSlOrderChannelMessage]struct{}
	orderSubscriberMtx     sync.Mutex
	positionSubscribersMtx sync.Mutex
	balanceSubscriberMtx   sync.Mutex
	tpSlOrderSubscriberMtx sync.Mutex
}

func NewPrivateWebsocket(ctx context.Context, apiKey, secretKey string) *privateWebsocketClient {
	ws := websocket.New(
		ctx,
		"wss://fapi.bitunix.com/private/",
		websocket.WithAuthentication(WebsocketSigner(apiKey, secretKey)),
		websocket.WithKeepAliveMonitor(30*time.Second, KeepAliveMonitor()),
	)

	return &privateWebsocketClient{
		websocketClient:        websocketClient{client: ws},
		orderSubscriberMtx:     sync.Mutex{},
		balanceSubscriberMtx:   sync.Mutex{},
		positionSubscribersMtx: sync.Mutex{},
		balanceSubscribers:     map[chan model.BalanceChannelMessage]struct{}{},
		positionSubscribers:    map[chan model.PositionChannelMessage]struct{}{},
		orderSubscribers:       map[chan model.OrderChannelMessage]struct{}{},
		tpSlOrderSubscribers:   map[chan model.TpSlOrderChannelMessage]struct{}{},
	}
}

func (ws *privateWebsocketClient) SubscribeBalance() <-chan model.BalanceChannelMessage {
	ws.balanceSubscriberMtx.Lock()
	defer ws.balanceSubscriberMtx.Unlock()
	c := make(chan model.BalanceChannelMessage)

	ws.balanceSubscribers[c] = struct{}{}

	return c
}

func (ws *privateWebsocketClient) UnsubscribeBalance(sub chan model.BalanceChannelMessage) {
	ws.balanceSubscriberMtx.Lock()
	defer ws.balanceSubscriberMtx.Unlock()

	if _, ok := ws.balanceSubscribers[sub]; ok {
		delete(ws.balanceSubscribers, sub)
	}
}

func (ws *privateWebsocketClient) SubscribePositions() <-chan model.PositionChannelMessage {
	ws.positionSubscribersMtx.Lock()
	defer ws.positionSubscribersMtx.Unlock()
	c := make(chan model.PositionChannelMessage)

	ws.positionSubscribers[c] = struct{}{}

	return c
}

func (ws *privateWebsocketClient) UnsubscribePosition(sub chan model.PositionChannelMessage) {
	ws.positionSubscribersMtx.Lock()
	defer ws.positionSubscribersMtx.Unlock()

	if _, ok := ws.positionSubscribers[sub]; ok {
		delete(ws.positionSubscribers, sub)
	}
}

func (ws *privateWebsocketClient) SubscribeOrders() <-chan model.OrderChannelMessage {
	ws.orderSubscriberMtx.Lock()
	defer ws.orderSubscriberMtx.Unlock()
	c := make(chan model.OrderChannelMessage)

	ws.orderSubscribers[c] = struct{}{}

	return c
}

func (ws *privateWebsocketClient) UnsubscribeOrders(sub chan model.OrderChannelMessage) {
	ws.orderSubscriberMtx.Lock()
	defer ws.orderSubscriberMtx.Unlock()

	if _, ok := ws.orderSubscribers[sub]; ok {
		delete(ws.orderSubscribers, sub)
	}
}

func (ws *privateWebsocketClient) SubscribeTpSlOrders() <-chan model.TpSlOrderChannelMessage {
	ws.tpSlOrderSubscriberMtx.Lock()
	defer ws.tpSlOrderSubscriberMtx.Unlock()
	c := make(chan model.TpSlOrderChannelMessage)

	ws.tpSlOrderSubscribers[c] = struct{}{}

	return c
}

func (ws *privateWebsocketClient) UnsubscribeTpSlOrders(sub chan model.TpSlOrderChannelMessage) {
	ws.tpSlOrderSubscriberMtx.Lock()
	defer ws.tpSlOrderSubscriberMtx.Unlock()

	if _, ok := ws.tpSlOrderSubscribers[sub]; ok {
		delete(ws.tpSlOrderSubscribers, sub)
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
					ws.populateBalanceResponse(bytes)
				case model.ChannelPosition:
					ws.populatePositionResponse(bytes)
				case model.ChannelOrder:
					ws.populateOrderResponse(bytes)
				case model.ChannelTpSl:
					ws.populateTpSlOrderResponse(bytes)
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

func (ws *privateWebsocketClient) populateTpSlOrderResponse(bytes []byte) {
	res := model.TpSlOrderChannelMessage{}
	if err := json.Unmarshal(bytes, &res); err != nil {
		log.WithError(err).Errorf("error unmarshaling tpslorder response JSON")
	}

	ws.tpSlOrderSubscriberMtx.Lock()
	defer ws.tpSlOrderSubscriberMtx.Unlock()
	for sub, _ := range ws.tpSlOrderSubscribers {
		sub <- res
	}
}

func (ws *privateWebsocketClient) populateOrderResponse(bytes []byte) {
	res := model.OrderChannelMessage{}
	if err := json.Unmarshal(bytes, &res); err != nil {
		log.WithError(err).Errorf("error unmarshaling position response JSON")
	}

	ws.orderSubscriberMtx.Lock()
	defer ws.orderSubscriberMtx.Unlock()
	for sub, _ := range ws.orderSubscribers {
		sub <- res
	}
}

func (ws *privateWebsocketClient) populatePositionResponse(bytes []byte) {
	res := model.PositionChannelMessage{}
	if err := json.Unmarshal(bytes, &res); err != nil {
		log.WithError(err).Errorf("error unmarshaling position response JSON")
	}

	ws.positionSubscribersMtx.Lock()
	defer ws.positionSubscribersMtx.Unlock()
	for sub, _ := range ws.positionSubscribers {
		sub <- res
	}
}

func (ws *privateWebsocketClient) populateBalanceResponse(bytes []byte) {
	res := model.BalanceChannelMessage{}
	if err := json.Unmarshal(bytes, &res); err != nil {
		log.WithError(err).Errorf("error unmarshaling balance response  JSON")
	}
	ws.balanceSubscriberMtx.Lock()
	defer ws.balanceSubscriberMtx.Unlock()
	for sub, _ := range ws.balanceSubscribers {
		sub <- res
	}
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
