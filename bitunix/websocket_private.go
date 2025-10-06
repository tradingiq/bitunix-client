package bitunix

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/security"
	"github.com/tradingiq/bitunix-client/websocket"
	"go.uber.org/zap"
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
	logger                 *zap.Logger
}

func NewPrivateWebsocket(ctx context.Context, apiKey, secretKey string, options ...WebsocketClientOption) (PrivateWebsocketClient, error) {
	wsc := &websocketClient{
		uri:      "wss://fapi.bitunix.com/private/",
		quit:     make(chan struct{}),
		logLevel: model.LogLevelNone,
	}
	for _, option := range options {
		option(wsc)
	}

	var wsOptions []websocket.ClientOption
	wsOptions = append(wsOptions, websocket.WithAuthentication(WebsocketSigner(apiKey, secretKey)))
	wsOptions = append(wsOptions, websocket.WithKeepAliveMonitor(30*time.Second, KeepAliveMonitor()))

	if wsc.logger != nil {
		wsOptions = append(wsOptions, websocket.WithLogger(wsc.logger))
	} else {
		wsOptions = append(wsOptions, websocket.WithLogLevel(wsc.logLevel))
	}

	wsc.client = websocket.New(
		ctx,
		wsc.uri,
		wsOptions...,
	)

	if wsc.workerPoolSize <= 0 {
		wsc.workerPoolSize = 10
	}

	if wsc.workerBufferSize > 0 {
		wsc.messageQueue = make(chan []byte, wsc.workerBufferSize)
	} else {
		wsc.messageQueue = make(chan []byte, 100)
	}

	var logger *zap.Logger
	if wsc.logger != nil {
		logger = wsc.logger
	} else {
		logger = createLoggerForLevel(wsc.logLevel)
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
		logger:                 logger,
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
		if ws.logger != nil {
			ws.logger.Error("failed to process websocket message", zap.Error(errors.NewInternalError("error unmarshaling websocket message", err)))
		}
		return
	}

	if errMsg, hasError := result["error"].(string); hasError && errMsg != "" {
		if ws.logger != nil {
			ws.logger.Error("received error from websocket server", zap.Error(errors.NewWebsocketError("message processing", errMsg, nil)))
		}
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
		if ws.logger != nil {
			ws.logger.Error("failed to process tp/sl order update", zap.Error(errors.NewInternalError("error unmarshaling tp/sl order response", err)))
		}
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
		if ws.logger != nil {
			ws.logger.Error("failed to process order update", zap.Error(errors.NewInternalError("error unmarshaling order response", err)))
		}
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
		if ws.logger != nil {
			ws.logger.Error("failed to process position update", zap.Error(errors.NewInternalError("error unmarshaling position response", err)))
		}
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
		if ws.logger != nil {
			ws.logger.Error("failed to process balance update", zap.Error(errors.NewInternalError("error unmarshaling balance response", err)))
		}
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

func WithWebsocketDebug(enabled bool) WebsocketClientOption {
	return func(ws *websocketClient) {
		if enabled {
			ws.logLevel = model.LogLevelAggressive
		} else {
			ws.logLevel = model.LogLevelNone
		}
	}
}

func WithWebsocketLogLevel(level model.LogLevel) WebsocketClientOption {
	return func(ws *websocketClient) {
		ws.logLevel = level
	}
}

func WithWebsocketLogger(logger *zap.Logger) WebsocketClientOption {
	return func(ws *websocketClient) {
		ws.logger = logger
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

type ReconnectingPrivateWebsocketClient struct {
	client               PrivateWebsocketClient
	ctx                  context.Context
	apiKey               string
	secretKey            string
	clientOptions        []WebsocketClientOption
	maxReconnectAttempts int
	reconnectDelay       time.Duration
	logger               *zap.Logger
	isConnected          bool
	mu                   sync.RWMutex
	stopReconnecting     chan struct{}
	balanceSubscribers   map[BalanceSubscriber]struct{}
	positionSubscribers  map[PositionSubscriber]struct{}
	orderSubscribers     map[OrderSubscriber]struct{}
	tpSlOrderSubscribers map[TpSlOrderSubscriber]struct{}
	subscriberMu         sync.RWMutex
}

type ReconnectingPrivateClientOption func(*ReconnectingPrivateWebsocketClient)

func WithPrivateMaxReconnectAttempts(attempts int) ReconnectingPrivateClientOption {
	return func(r *ReconnectingPrivateWebsocketClient) {
		r.maxReconnectAttempts = attempts
	}
}

func WithPrivateReconnectDelay(delay time.Duration) ReconnectingPrivateClientOption {
	return func(r *ReconnectingPrivateWebsocketClient) {
		r.reconnectDelay = delay
	}
}

func WithPrivateReconnectLogger(logger *zap.Logger) ReconnectingPrivateClientOption {
	return func(r *ReconnectingPrivateWebsocketClient) {
		r.logger = logger
	}
}

func NewReconnectingPrivateWebsocket(ctx context.Context, apiKey, secretKey string, clientOptions []WebsocketClientOption, options ...ReconnectingPrivateClientOption) (*ReconnectingPrivateWebsocketClient, error) {
	client, err := NewPrivateWebsocket(ctx, apiKey, secretKey, clientOptions...)
	if err != nil {
		return nil, err
	}

	r := &ReconnectingPrivateWebsocketClient{
		client:               client,
		ctx:                  ctx,
		apiKey:               apiKey,
		secretKey:            secretKey,
		clientOptions:        clientOptions,
		maxReconnectAttempts: 0,
		reconnectDelay:       5 * time.Second,
		logger:               zap.NewNop(),
		stopReconnecting:     make(chan struct{}),
		balanceSubscribers:   make(map[BalanceSubscriber]struct{}),
		positionSubscribers:  make(map[PositionSubscriber]struct{}),
		orderSubscribers:     make(map[OrderSubscriber]struct{}),
		tpSlOrderSubscribers: make(map[TpSlOrderSubscriber]struct{}),
	}

	for _, option := range options {
		option(r)
	}

	return r, nil
}

func (r *ReconnectingPrivateWebsocketClient) Connect() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	err := r.client.Connect()
	if err == nil {
		r.isConnected = true
	}
	return err
}

func (r *ReconnectingPrivateWebsocketClient) Disconnect() {
	r.mu.Lock()
	defer r.mu.Unlock()

	close(r.stopReconnecting)
	r.isConnected = false
	r.client.Disconnect()
}

func (r *ReconnectingPrivateWebsocketClient) SubscribeBalance(subscriber BalanceSubscriber) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	err := r.client.SubscribeBalance(subscriber)
	if err == nil {
		r.subscriberMu.Lock()
		r.balanceSubscribers[subscriber] = struct{}{}
		r.subscriberMu.Unlock()
	}
	return err
}

func (r *ReconnectingPrivateWebsocketClient) UnsubscribeBalance(subscriber BalanceSubscriber) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	err := r.client.UnsubscribeBalance(subscriber)
	if err == nil {
		r.subscriberMu.Lock()
		delete(r.balanceSubscribers, subscriber)
		r.subscriberMu.Unlock()
	}
	return err
}

func (r *ReconnectingPrivateWebsocketClient) SubscribePositions(subscriber PositionSubscriber) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	err := r.client.SubscribePositions(subscriber)
	if err == nil {
		r.subscriberMu.Lock()
		r.positionSubscribers[subscriber] = struct{}{}
		r.subscriberMu.Unlock()
	}
	return err
}

func (r *ReconnectingPrivateWebsocketClient) UnsubscribePositions(subscriber PositionSubscriber) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	err := r.client.UnsubscribePositions(subscriber)
	if err == nil {
		r.subscriberMu.Lock()
		delete(r.positionSubscribers, subscriber)
		r.subscriberMu.Unlock()
	}
	return err
}

func (r *ReconnectingPrivateWebsocketClient) SubscribeOrders(subscriber OrderSubscriber) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	err := r.client.SubscribeOrders(subscriber)
	if err == nil {
		r.subscriberMu.Lock()
		r.orderSubscribers[subscriber] = struct{}{}
		r.subscriberMu.Unlock()
	}
	return err
}

func (r *ReconnectingPrivateWebsocketClient) UnsubscribeOrders(subscriber OrderSubscriber) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	err := r.client.UnsubscribeOrders(subscriber)
	if err == nil {
		r.subscriberMu.Lock()
		delete(r.orderSubscribers, subscriber)
		r.subscriberMu.Unlock()
	}
	return err
}

func (r *ReconnectingPrivateWebsocketClient) SubscribeTpSlOrders(subscriber TpSlOrderSubscriber) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	err := r.client.SubscribeTpSlOrders(subscriber)
	if err == nil {
		r.subscriberMu.Lock()
		r.tpSlOrderSubscribers[subscriber] = struct{}{}
		r.subscriberMu.Unlock()
	}
	return err
}

func (r *ReconnectingPrivateWebsocketClient) UnsubscribeTpSlOrders(subscriber TpSlOrderSubscriber) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	err := r.client.UnsubscribeTpSlOrders(subscriber)
	if err == nil {
		r.subscriberMu.Lock()
		delete(r.tpSlOrderSubscribers, subscriber)
		r.subscriberMu.Unlock()
	}
	return err
}

func (r *ReconnectingPrivateWebsocketClient) Stream() error {
	attempt := 0
	for {
		r.mu.RLock()
		if !r.isConnected {
			r.mu.RUnlock()
			return errors.NewWebsocketError("stream", "not connected", nil)
		}
		r.mu.RUnlock()

		err := r.client.Stream()
		if err == nil {
			return nil
		}

		r.logger.Error("private websocket stream error", zap.Error(err), zap.Int("attempt", attempt+1))

		r.mu.Lock()
		r.isConnected = false
		r.mu.Unlock()

		if r.maxReconnectAttempts > 0 && attempt >= r.maxReconnectAttempts {
			r.logger.Error("max reconnect attempts reached for private websocket", zap.Int("attempts", attempt))
			return errors.NewWebsocketError("stream", "max reconnect attempts reached", err)
		}

		select {
		case <-r.ctx.Done():
			return errors.NewConnectionClosedError("stream", "context cancelled during reconnection", r.ctx.Err())
		case <-r.stopReconnecting:
			return errors.NewWebsocketError("stream", "reconnection stopped", nil)
		case <-time.After(r.reconnectDelay):
		}

		r.logger.Info("attempting to reconnect private websocket", zap.Int("attempt", attempt+1))

		if reconnectErr := r.connectWithResubscription(); reconnectErr != nil {
			r.logger.Error("private websocket reconnection failed", zap.Error(reconnectErr), zap.Int("attempt", attempt+1))
			attempt++
			continue
		}

		r.logger.Info("private websocket reconnected successfully", zap.Int("attempt", attempt+1))
		attempt = 0
	}
}

func (r *ReconnectingPrivateWebsocketClient) connectWithResubscription() error {
	// Disconnect old client first
	r.client.Disconnect()

	// Create a new client instance
	newClient, err := NewPrivateWebsocket(r.ctx, r.apiKey, r.secretKey, r.clientOptions...)
	if err != nil {
		return err
	}

	// Connect the new client
	err = newClient.Connect()
	if err != nil {
		return err
	}

	// Replace the old client with the new one
	r.client = newClient

	r.mu.Lock()
	r.isConnected = true
	r.mu.Unlock()

	return r.resubscribeAll()
}

func (r *ReconnectingPrivateWebsocketClient) resubscribeAll() error {
	r.subscriberMu.RLock()
	defer r.subscriberMu.RUnlock()

	for subscriber := range r.balanceSubscribers {
		err := r.client.SubscribeBalance(subscriber)
		if err != nil {
			r.logger.Error("failed to resubscribe to balance updates", zap.Error(err))
			return err
		}
		r.logger.Debug("resubscribed to balance updates successfully")
	}

	for subscriber := range r.positionSubscribers {
		err := r.client.SubscribePositions(subscriber)
		if err != nil {
			r.logger.Error("failed to resubscribe to position updates", zap.Error(err))
			return err
		}
		r.logger.Debug("resubscribed to position updates successfully")
	}

	for subscriber := range r.orderSubscribers {
		err := r.client.SubscribeOrders(subscriber)
		if err != nil {
			r.logger.Error("failed to resubscribe to order updates", zap.Error(err))
			return err
		}
		r.logger.Debug("resubscribed to order updates successfully")
	}

	for subscriber := range r.tpSlOrderSubscribers {
		err := r.client.SubscribeTpSlOrders(subscriber)
		if err != nil {
			r.logger.Error("failed to resubscribe to tp/sl order updates", zap.Error(err))
			return err
		}
		r.logger.Debug("resubscribed to tp/sl order updates successfully")
	}

	return nil
}
