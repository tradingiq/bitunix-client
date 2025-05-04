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

type privateWebsocketClient struct {
	websocketClient
	balanceSubscribers     map[BalanceSubscriber]struct{}
	positionSubscribers    map[PositionSubscriber]struct{}
	orderSubscribers       map[OrderSubscriber]struct{}
	tpSlOrderSubscribers   map[TpSlOrderSubscriber]struct{}
	orderSubscriberMtx     sync.Mutex
	positionSubscribersMtx sync.Mutex
	balanceSubscriberMtx   sync.Mutex
	tpSlOrderSubscriberMtx sync.Mutex
}

func NewPrivateWebsocket(ctx context.Context, apiKey, secretKey string, options ...WebsocketClientOption) PrivateWebsocketClient {
	wsc := websocketClient{
		uri: "wss://fapi.bitunix.com/private/",
	}
	for _, option := range options {
		option(&wsc)
	}

	wsc.client = websocket.New(
		ctx,
		wsc.uri,
		websocket.WithAuthentication(WebsocketSigner(apiKey, secretKey)),
		websocket.WithKeepAliveMonitor(30*time.Second, KeepAliveMonitor()),
	)

	return &privateWebsocketClient{
		websocketClient:        wsc,
		orderSubscriberMtx:     sync.Mutex{},
		balanceSubscriberMtx:   sync.Mutex{},
		positionSubscribersMtx: sync.Mutex{},
		tpSlOrderSubscriberMtx: sync.Mutex{},
		balanceSubscribers:     map[BalanceSubscriber]struct{}{},
		positionSubscribers:    map[PositionSubscriber]struct{}{},
		orderSubscribers:       map[OrderSubscriber]struct{}{},
		tpSlOrderSubscribers:   map[TpSlOrderSubscriber]struct{}{},
	}
}

func (ws *privateWebsocketClient) SubscribeBalance(subscriber BalanceSubscriber) error {
	if subscriber == nil {
		return fmt.Errorf("balance subscriber cannot be nil")
	}

	ws.balanceSubscriberMtx.Lock()
	defer ws.balanceSubscriberMtx.Unlock()

	ws.balanceSubscribers[subscriber] = struct{}{}
	return nil
}

func (ws *privateWebsocketClient) UnsubscribeBalance(subscriber BalanceSubscriber) error {
	if subscriber == nil {
		return fmt.Errorf("balance subscriber cannot be nil")
	}

	ws.balanceSubscriberMtx.Lock()
	defer ws.balanceSubscriberMtx.Unlock()

	if _, ok := ws.balanceSubscribers[subscriber]; ok {
		delete(ws.balanceSubscribers, subscriber)
	}
	return nil
}

func (ws *privateWebsocketClient) SubscribePositions(subscriber PositionSubscriber) error {
	if subscriber == nil {
		return fmt.Errorf("position subscriber cannot be nil")
	}

	ws.positionSubscribersMtx.Lock()
	defer ws.positionSubscribersMtx.Unlock()

	ws.positionSubscribers[subscriber] = struct{}{}
	return nil
}

func (ws *privateWebsocketClient) UnsubscribePosition(subscriber PositionSubscriber) error {
	if subscriber == nil {
		return fmt.Errorf("position subscriber cannot be nil")
	}

	ws.positionSubscribersMtx.Lock()
	defer ws.positionSubscribersMtx.Unlock()

	if _, ok := ws.positionSubscribers[subscriber]; ok {
		delete(ws.positionSubscribers, subscriber)
	}
	return nil
}

type BalanceSubscriber interface {
	Handle(*model.BalanceChannelMessage)
}

type PositionSubscriber interface {
	Handle(*model.PositionChannelMessage)
}

type OrderSubscriber interface {
	Handle(*model.OrderChannelMessage)
}

type TpSlOrderSubscriber interface {
	Handle(*model.TpSlOrderChannelMessage)
}

func (ws *privateWebsocketClient) SubscribeOrders(subscriber OrderSubscriber) error {
	if subscriber == nil {
		return fmt.Errorf("order subscriber cannot be nil")
	}

	ws.orderSubscriberMtx.Lock()
	defer ws.orderSubscriberMtx.Unlock()

	ws.orderSubscribers[subscriber] = struct{}{}
	return nil
}

func (ws *privateWebsocketClient) UnsubscribeOrders(subscriber OrderSubscriber) error {
	if subscriber == nil {
		return fmt.Errorf("order subscriber cannot be nil")
	}

	ws.orderSubscriberMtx.Lock()
	defer ws.orderSubscriberMtx.Unlock()

	if _, ok := ws.orderSubscribers[subscriber]; ok {
		delete(ws.orderSubscribers, subscriber)
	}
	return nil
}

func (ws *privateWebsocketClient) SubscribeTpSlOrders(subscriber TpSlOrderSubscriber) error {
	if subscriber == nil {
		return fmt.Errorf("tp/sl order subscriber cannot be nil")
	}

	ws.tpSlOrderSubscriberMtx.Lock()
	defer ws.tpSlOrderSubscriberMtx.Unlock()

	ws.tpSlOrderSubscribers[subscriber] = struct{}{}
	return nil
}

func (ws *privateWebsocketClient) UnsubscribeTpSlOrders(subscriber TpSlOrderSubscriber) error {
	if subscriber == nil {
		return fmt.Errorf("tp/sl order subscriber cannot be nil")
	}

	ws.tpSlOrderSubscriberMtx.Lock()
	defer ws.tpSlOrderSubscriberMtx.Unlock()

	if _, ok := ws.tpSlOrderSubscribers[subscriber]; ok {
		delete(ws.tpSlOrderSubscribers, subscriber)
	}
	return nil
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
		return
	}

	ws.tpSlOrderSubscriberMtx.Lock()
	defer ws.tpSlOrderSubscriberMtx.Unlock()
	for sub := range ws.tpSlOrderSubscribers {
		sub.Handle(&res)
	}
}

func (ws *privateWebsocketClient) populateOrderResponse(bytes []byte) {
	res := model.OrderChannelMessage{}
	if err := json.Unmarshal(bytes, &res); err != nil {
		log.WithError(err).Errorf("error unmarshaling order response JSON")
		return
	}

	ws.orderSubscriberMtx.Lock()
	defer ws.orderSubscriberMtx.Unlock()
	for sub := range ws.orderSubscribers {
		sub.Handle(&res)
	}
}

func (ws *privateWebsocketClient) populatePositionResponse(bytes []byte) {
	res := model.PositionChannelMessage{}
	if err := json.Unmarshal(bytes, &res); err != nil {
		log.WithError(err).Errorf("error unmarshaling position response JSON")
		return
	}

	ws.positionSubscribersMtx.Lock()
	defer ws.positionSubscribersMtx.Unlock()
	for sub := range ws.positionSubscribers {
		sub.Handle(&res)
	}
}

func (ws *privateWebsocketClient) populateBalanceResponse(bytes []byte) {
	res := model.BalanceChannelMessage{}
	if err := json.Unmarshal(bytes, &res); err != nil {
		log.WithError(err).Errorf("error unmarshaling balance response JSON")
		return
	}

	ws.balanceSubscriberMtx.Lock()
	defer ws.balanceSubscriberMtx.Unlock()
	for sub := range ws.balanceSubscribers {
		sub.Handle(&res)
	}
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
					ApiKey:    apiKey,
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

type PrivateWebsocketClient interface {
	SubscribeBalance(subscriber BalanceSubscriber) error
	UnsubscribeBalance(subscriber BalanceSubscriber) error
	SubscribePositions(subscriber PositionSubscriber) error
	UnsubscribePosition(subscriber PositionSubscriber) error
	SubscribeOrders(subscriber OrderSubscriber) error
	UnsubscribeOrders(subscriber OrderSubscriber) error
	SubscribeTpSlOrders(subscriber TpSlOrderSubscriber) error
	UnsubscribeTpSlOrders(subscriber TpSlOrderSubscriber) error
	Stream() error
	Connect() error
	Disconnect()
}
