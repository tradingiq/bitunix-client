package bitunix

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/security"
	"sync"
	"time"
)

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

func NewPrivateWebsocket(ctx context.Context, apiKey, secretKey string, options ...WebsocketClientOption) PrivateWebsocketClient {
	wsc := websocketClient{
		uri: "wss://fapi.bitunix.com/private/",
	}
	for _, option := range options {
		option(&wsc)
	}

	wsc.client = newWebsocket(ctx, apiKey, secretKey, wsc.uri)

	return &privateWebsocketClient{
		websocketClient:        wsc,
		orderSubscriberMtx:     sync.Mutex{},
		balanceSubscriberMtx:   sync.Mutex{},
		positionSubscribersMtx: sync.Mutex{},
		tpSlOrderSubscriberMtx: sync.Mutex{},
		balanceSubscribers:     map[chan model.BalanceChannelMessage]struct{}{},
		positionSubscribers:    map[chan model.PositionChannelMessage]struct{}{},
		orderSubscribers:       map[chan model.OrderChannelMessage]struct{}{},
		tpSlOrderSubscribers:   map[chan model.TpSlOrderChannelMessage]struct{}{},
	}
}

func (ws *privateWebsocketClient) SubscribeBalance() <-chan model.BalanceChannelMessage {
	ws.balanceSubscriberMtx.Lock()
	defer ws.balanceSubscriberMtx.Unlock()
	c := make(chan model.BalanceChannelMessage, 10)

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
	c := make(chan model.PositionChannelMessage, 10)

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
	c := make(chan model.OrderChannelMessage, 10)

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
	c := make(chan model.TpSlOrderChannelMessage, 10)

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
	for sub := range ws.tpSlOrderSubscribers {
		select {
		case sub <- res:
		default:
			log.Errorf("tp sl channel poisoned, data might be not recent")
		}
	}
}

func (ws *privateWebsocketClient) populateOrderResponse(bytes []byte) {
	res := model.OrderChannelMessage{}
	if err := json.Unmarshal(bytes, &res); err != nil {
		log.WithError(err).Errorf("error unmarshaling position response JSON")
	}

	ws.orderSubscriberMtx.Lock()
	defer ws.orderSubscriberMtx.Unlock()
	for sub := range ws.orderSubscribers {
		select {
		case sub <- res:
		default:
			log.Errorf("order channel poisoned, data might be not recent")
		}
	}
}

func (ws *privateWebsocketClient) populatePositionResponse(bytes []byte) {
	res := model.PositionChannelMessage{}
	if err := json.Unmarshal(bytes, &res); err != nil {
		log.WithError(err).Errorf("error unmarshaling position response JSON")
	}

	ws.positionSubscribersMtx.Lock()
	defer ws.positionSubscribersMtx.Unlock()
	for sub := range ws.positionSubscribers {
		select {
		case sub <- res:
		default:
			log.Errorf("position channel poisoned, data might be not recent")
		}
	}
}

func (ws *privateWebsocketClient) populateBalanceResponse(bytes []byte) {
	res := model.BalanceChannelMessage{}
	if err := json.Unmarshal(bytes, &res); err != nil {
		log.WithError(err).Errorf("error unmarshaling balance response  JSON")
	}
	ws.balanceSubscriberMtx.Lock()
	defer ws.balanceSubscriberMtx.Unlock()
	for sub := range ws.balanceSubscribers {
		select {
		case sub <- res:
		default:
			log.Errorf("balance channel poisoned, data might be not recent")
		}
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

type PrivateWebsocketClient interface {
	SubscribeBalance() <-chan model.BalanceChannelMessage
	UnsubscribeBalance(sub chan model.BalanceChannelMessage)
	SubscribePositions() <-chan model.PositionChannelMessage
	UnsubscribePosition(sub chan model.PositionChannelMessage)
	SubscribeOrders() <-chan model.OrderChannelMessage
	UnsubscribeOrders(sub chan model.OrderChannelMessage)
	SubscribeTpSlOrders() <-chan model.TpSlOrderChannelMessage
	UnsubscribeTpSlOrders(sub chan model.TpSlOrderChannelMessage)
	Stream() error
	Connect() error
	Disconnect()
}
