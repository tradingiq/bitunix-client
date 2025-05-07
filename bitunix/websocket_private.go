package bitunix

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/security"
	"github.com/tradingiq/bitunix-client/websocket"
	"sync"
	"time"
)

type privateWebsocketClient struct {
	*websocketClient
	balanceSubscribers     map[BalanceSubscriber]struct{}
	positionSubscribers    map[PositionSubscriber]struct{}
	orderSubscribers       map[OrderSubscriber]struct{}
	tpSlOrderSubscribers   map[TpSlOrderSubscriber]struct{}
	orderSubscriberMtx     sync.Mutex
	positionSubscribersMtx sync.Mutex
	balanceSubscriberMtx   sync.Mutex
	tpSlOrderSubscriberMtx sync.Mutex
}

func NewPrivateWebsocket(ctx context.Context, apiKey, secretKey string, options ...WebsocketClientOption) (PrivateWebsocketClient, error) {
	wsc := &websocketClient{
		uri:  "wss://fapi.bitunix.com/private/",
		quit: make(chan struct{}),
	}
	for _, option := range options {
		option(wsc)
	}

	wsc.client = websocket.New(
		ctx,
		wsc.uri,
		websocket.WithAuthentication(WebsocketSigner(apiKey, secretKey)),
		websocket.WithKeepAliveMonitor(30*time.Second, KeepAliveMonitor()),
	)

	if wsc.workerPoolSize <= 0 {
		wsc.workerPoolSize = 10
	}

	if wsc.workerBufferSize > 0 {
		wsc.messageQueue = make(chan []byte, wsc.workerBufferSize)
	} else {
		wsc.messageQueue = make(chan []byte, 100)
	}

	client := &privateWebsocketClient{
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

	wsc.processFunc = client.processMessage

	if err := client.startWorkerPool(ctx); err != nil {
		return nil, errors.NewWebsocketError("initialize private websocket", "failed to start worker pool", err)
	}

	return client, nil
}

func (ws *privateWebsocketClient) SubscribeBalance(subscriber BalanceSubscriber) error {
	if subscriber == nil {
		return errors.NewValidationError("subscriber", "balance subscriber cannot be nil", nil)
	}

	ws.balanceSubscriberMtx.Lock()
	defer ws.balanceSubscriberMtx.Unlock()

	ws.balanceSubscribers[subscriber] = struct{}{}
	return nil
}

func (ws *privateWebsocketClient) UnsubscribeBalance(subscriber BalanceSubscriber) error {
	if subscriber == nil {
		return errors.NewValidationError("subscriber", "balance subscriber cannot be nil", nil)
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
		return errors.NewValidationError("subscriber", "position subscriber cannot be nil", nil)
	}

	ws.positionSubscribersMtx.Lock()
	defer ws.positionSubscribersMtx.Unlock()

	ws.positionSubscribers[subscriber] = struct{}{}
	return nil
}

func (ws *privateWebsocketClient) UnsubscribePositions(subscriber PositionSubscriber) error {
	if subscriber == nil {
		return errors.NewValidationError("subscriber", "position subscriber cannot be nil", nil)
	}

	ws.positionSubscribersMtx.Lock()
	defer ws.positionSubscribersMtx.Unlock()

	if _, ok := ws.positionSubscribers[subscriber]; ok {
		delete(ws.positionSubscribers, subscriber)
	}
	return nil
}

type BalanceSubscriber interface {
	SubscribeBalance(*model.BalanceChannelMessage)
}

type PositionSubscriber interface {
	SubscribePosition(*model.PositionChannelMessage)
}

type OrderSubscriber interface {
	SubscribeOrder(*model.OrderChannelMessage)
}

type TpSlOrderSubscriber interface {
	SubscribeTpSlOrder(*model.TpSlOrderChannelMessage)
}

func (ws *privateWebsocketClient) SubscribeOrders(subscriber OrderSubscriber) error {
	if subscriber == nil {
		return errors.NewValidationError("subscriber", "order subscriber cannot be nil", nil)
	}

	ws.orderSubscriberMtx.Lock()
	defer ws.orderSubscriberMtx.Unlock()

	ws.orderSubscribers[subscriber] = struct{}{}
	return nil
}

func (ws *privateWebsocketClient) UnsubscribeOrders(subscriber OrderSubscriber) error {
	if subscriber == nil {
		return errors.NewValidationError("subscriber", "order subscriber cannot be nil", nil)
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
		return errors.NewValidationError("subscriber", "tp/sl order subscriber cannot be nil", nil)
	}

	ws.tpSlOrderSubscriberMtx.Lock()
	defer ws.tpSlOrderSubscriberMtx.Unlock()

	ws.tpSlOrderSubscribers[subscriber] = struct{}{}
	return nil
}

func (ws *privateWebsocketClient) UnsubscribeTpSlOrders(subscriber TpSlOrderSubscriber) error {
	if subscriber == nil {
		return errors.NewValidationError("subscriber", "tp/sl order subscriber cannot be nil", nil)
	}

	ws.tpSlOrderSubscriberMtx.Lock()
	defer ws.tpSlOrderSubscriberMtx.Unlock()

	if _, ok := ws.tpSlOrderSubscribers[subscriber]; ok {
		delete(ws.tpSlOrderSubscribers, subscriber)
	}
	return nil
}

func (ws *privateWebsocketClient) processMessage(bytes []byte) {
	var result map[string]interface{}

	err := json.Unmarshal(bytes, &result)
	if err != nil {
		log.WithError(
			errors.NewInternalError("error unmarshaling websocket message", err),
		).Error("failed to process websocket message")
		return
	}

	if errMsg, hasError := result["error"].(string); hasError && errMsg != "" {
		log.WithError(
			errors.NewWebsocketError("message processing", errMsg, nil),
		).Error("received error from websocket server")
		return
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
}

func (ws *privateWebsocketClient) populateTpSlOrderResponse(bytes []byte) {
	res := model.TpSlOrderChannelMessage{}
	if err := json.Unmarshal(bytes, &res); err != nil {
		log.WithError(
			errors.NewInternalError("error unmarshaling tp/sl order response", err),
		).Error("failed to process tp/sl order update")
		return
	}

	ws.tpSlOrderSubscriberMtx.Lock()
	defer ws.tpSlOrderSubscriberMtx.Unlock()
	for sub := range ws.tpSlOrderSubscribers {
		sub.SubscribeTpSlOrder(&res)
	}
}

func (ws *privateWebsocketClient) populateOrderResponse(bytes []byte) {
	res := model.OrderChannelMessage{}
	if err := json.Unmarshal(bytes, &res); err != nil {
		log.WithError(
			errors.NewInternalError("error unmarshaling order response", err),
		).Error("failed to process order update")
		return
	}

	ws.orderSubscriberMtx.Lock()
	defer ws.orderSubscriberMtx.Unlock()
	for sub := range ws.orderSubscribers {
		sub.SubscribeOrder(&res)
	}
}

func (ws *privateWebsocketClient) populatePositionResponse(bytes []byte) {
	res := model.PositionChannelMessage{}
	if err := json.Unmarshal(bytes, &res); err != nil {
		log.WithError(
			errors.NewInternalError("error unmarshaling position response", err),
		).Error("failed to process position update")
		return
	}

	ws.positionSubscribersMtx.Lock()
	defer ws.positionSubscribersMtx.Unlock()
	for sub := range ws.positionSubscribers {
		sub.SubscribePosition(&res)
	}
}

func (ws *privateWebsocketClient) populateBalanceResponse(bytes []byte) {
	res := model.BalanceChannelMessage{}
	if err := json.Unmarshal(bytes, &res); err != nil {
		log.WithError(
			errors.NewInternalError("error unmarshaling balance response", err),
		).Error("failed to process balance update")
		return
	}

	ws.balanceSubscriberMtx.Lock()
	defer ws.balanceSubscriberMtx.Unlock()
	for sub := range ws.balanceSubscribers {
		sub.SubscribeBalance(&res)
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

func WithWorkerPoolSize(size int) WebsocketClientOption {
	return func(ws *websocketClient) {
		ws.workerPoolSize = size
	}
}

func WithWorkerBufferSize(size int) WebsocketClientOption {
	return func(ws *websocketClient) {
		ws.workerBufferSize = size
	}
}

func WebsocketSigner(apiKey, apiSecret string) func() ([]byte, error) {
	return func() ([]byte, error) {
		nonce, err := security.GenerateNonce(32)
		if err != nil {
			return nil, errors.NewAuthenticationError("failed to generate nonce for websocket authentication", err)
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
			return nil, errors.NewInternalError("failed to marshal login request", err)
		}

		return bytes, nil
	}
}

type PrivateWebsocketClient interface {
	SubscribeBalance(subscriber BalanceSubscriber) error
	UnsubscribeBalance(subscriber BalanceSubscriber) error
	SubscribePositions(subscriber PositionSubscriber) error
	UnsubscribePositions(subscriber PositionSubscriber) error
	SubscribeOrders(subscriber OrderSubscriber) error
	UnsubscribeOrders(subscriber OrderSubscriber) error
	SubscribeTpSlOrders(subscriber TpSlOrderSubscriber) error
	UnsubscribeTpSlOrders(subscriber TpSlOrderSubscriber) error
	Stream() error
	Connect() error
	Disconnect()
}
