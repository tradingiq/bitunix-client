package bitunix

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/websocket"
	"go.uber.org/zap"
	"strings"
	"sync"
	"time"
)

type wsClientInterface interface {
	Connect() error
	Listen(callback websocket.HandlerFunc) error
	Close()
	Write(bytes []byte) error
}

type websocketClient struct {
	client           wsClientInterface
	uri              string
	workerPoolSize   int
	workerBufferSize int
	messageQueue     chan []byte
	quit             chan struct{}
	processFunc      func(bytes []byte)
	logLevel         model.LogLevel
	logger           *zap.Logger
}

func (ws *websocketClient) Connect() error {
	if err := ws.client.Connect(); err != nil {
		return errors.NewWebsocketError("connect", "client failed to connect", err)
	}

	return nil
}

func (ws *websocketClient) Stream() error {
	err := ws.client.Listen(func(bytes []byte) error {
		select {
		case ws.messageQueue <- bytes:
		default:
			return errors.NewWorkgroupExhaustedError("stream", "workgroup exhausted", nil)
		}
		return nil
	})

	if err != nil {
		return errors.NewWebsocketError("stream", "client failed to listen", err)
	}

	return nil
}

func (ws *websocketClient) Disconnect() {
	ws.client.Close()

	close(ws.quit)
	close(ws.messageQueue)
}

func (ws *websocketClient) startWorkerPool(ctx context.Context) error {
	if ws.processFunc == nil {
		return errors.NewInternalError("processFunc is nil", nil)
	}

	for i := 0; i < ws.workerPoolSize; i++ {
		go ws.worker(ctx)
	}

	return nil
}

func (ws *websocketClient) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ws.quit:
			return
		case msg := <-ws.messageQueue:
			ws.processFunc(msg)
		}
	}
}

type publicWebsocketClient struct {
	*websocketClient
	subscriberMtx sync.Mutex
	klineHandlers map[KLineSubscriber]struct{}
	logger        *zap.Logger
}

type WebsocketClientOption func(*websocketClient)

func createLoggerForLevel(level model.LogLevel) *zap.Logger {
	switch level {
	case model.LogLevelNone:
		return zap.NewNop()
	case model.LogLevelAggressive:
		logger, _ := zap.NewDevelopment()
		return logger
	case model.LogLevelVeryAggressive:
		config := zap.NewDevelopmentConfig()
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		config.Development = true
		config.DisableCaller = false
		config.DisableStacktrace = false
		logger, _ := config.Build()
		return logger
	default:
		logger, _ := zap.NewDevelopment()
		return logger
	}
}

func NewPublicWebsocket(ctx context.Context, options ...WebsocketClientOption) (PublicWebsocketClient, error) {
	wsc := &websocketClient{
		uri:      "wss://fapi.bitunix.com/public/",
		quit:     make(chan struct{}),
		logLevel: model.LogLevelNone,
	}
	for _, option := range options {
		option(wsc)
	}

	var wsOptions []websocket.ClientOption
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

	client := &publicWebsocketClient{
		websocketClient: wsc,
		subscriberMtx:   sync.Mutex{},
		klineHandlers:   make(map[KLineSubscriber]struct{}),
		logger:          logger,
	}
	wsc.processFunc = client.processMessage

	if err := client.startWorkerPool(ctx); err != nil {
		return nil, errors.NewWebsocketError("initialize", "failed to start worker pool", err)
	}

	return client, nil
}

func (ws *publicWebsocketClient) SubscribeKLine(subscriber KLineSubscriber) error {
	if subscriber == nil {
		return errors.NewValidationError("subscriber", "cannot be nil", nil)
	}

	symbol := subscriber.SubscribeSymbol().Normalize()
	priceType := subscriber.SubscribePriceType().Normalize()
	interval := subscriber.SubscribeInterval().Normalize()

	channelName := fmt.Sprintf("%s_kline_%s", priceType, interval)

	ws.subscriberMtx.Lock()
	defer ws.subscriberMtx.Unlock()

	needsSubscription := true
	for existingSubscriber := range ws.klineHandlers {
		if existingSubscriber.SubscribeSymbol().Normalize() == symbol &&
			existingSubscriber.SubscribeInterval().Normalize() == interval &&
			existingSubscriber.SubscribePriceType().Normalize() == priceType {
			needsSubscription = false
			break
		}
	}

	ws.klineHandlers[subscriber] = struct{}{}

	if needsSubscription {
		req := SubscribeRequest{
			Op: "subscribe",
			Args: []interface{}{
				SubscribeKLineRequest{
					Symbol: symbol.String(),
					Ch:     channelName,
				},
			},
		}

		bytes, err := json.Marshal(req)
		if err != nil {
			return errors.NewInternalError("failed to marshal subscription request", err)
		}

		if err := ws.client.Write(bytes); err != nil {
			return errors.NewWebsocketError("subscribe", "failed to send subscription request", err)
		}
	}

	return nil
}

func (ws *publicWebsocketClient) UnsubscribeKLine(subscriber KLineSubscriber) error {
	if subscriber == nil {
		return errors.NewValidationError("subscriber", "cannot be nil", nil)
	}

	symbol := subscriber.SubscribeSymbol().Normalize()
	priceType := subscriber.SubscribePriceType().Normalize()
	interval := subscriber.SubscribeInterval().Normalize()

	channelName := fmt.Sprintf("%s_kline_%s", priceType, interval)

	ws.subscriberMtx.Lock()
	defer ws.subscriberMtx.Unlock()

	delete(ws.klineHandlers, subscriber)

	hasRemainingSubscribers := false
	for remainingSubscriber := range ws.klineHandlers {
		if remainingSubscriber.SubscribeSymbol().Normalize() == symbol &&
			remainingSubscriber.SubscribeInterval().Normalize() == interval &&
			remainingSubscriber.SubscribePriceType().Normalize() == priceType {
			hasRemainingSubscribers = true
			break
		}
	}

	if !hasRemainingSubscribers {
		req := SubscribeRequest{
			Op: "unsubscribe",
			Args: []interface{}{
				SubscribeKLineRequest{
					Symbol: symbol.String(),
					Ch:     channelName,
				},
			},
		}

		bytes, err := json.Marshal(req)
		if err != nil {
			return errors.NewInternalError("failed to marshal unsubscription request", err)
		}

		if err := ws.client.Write(bytes); err != nil {
			return errors.NewWebsocketError("unsubscribe", "failed to send unsubscription request", err)
		}
	}

	return nil
}

func parseChannel(channelStr string) (model.Interval, model.Channel, model.PriceType, error) {
	parts := strings.Split(channelStr, "_")

	if len(parts) != 3 {
		return "", "", "", errors.NewValidationError(
			"channel",
			fmt.Sprintf("invalid channel format: expected priceType_type_interval, got %s", channelStr),
			nil,
		)
	}

	priceTypeStr := parts[0]
	channelStr = parts[1]
	intervalStr := parts[2]

	priceType, err := model.ParsePriceType(priceTypeStr)
	if err != nil {
		return "", "", "", errors.NewValidationError("priceType", "failed to parse price type", err)
	}

	interval, err := model.ParseInterval(intervalStr)
	if err != nil {
		return "", "", "", errors.NewValidationError("interval", "failed to parse interval", err)
	}

	channel, err := model.ParseChannel(channelStr)
	if err != nil {
		return "", "", "", errors.NewValidationError("channel", "failed to parse channel", err)
	}

	return interval, channel, priceType, nil
}

func (ws *publicWebsocketClient) processMessage(bytes []byte) {
	var result map[string]interface{}

	err := json.Unmarshal(bytes, &result)
	if err != nil {
		if ws.logger != nil {
			ws.logger.Error("failed to process message", zap.Error(errors.NewInternalError("error unmarshaling JSON", err)))
		}
		return
	}

	if ch, ok := result["ch"].(string); ok {
		if sym, symbolOk := result["symbol"].(string); symbolOk {
			interval, channel, priceType, err := parseChannel(ch)
			if err != nil {
				if ws.logger != nil {
					ws.logger.Error("error parsing channel", zap.Error(err))
				}
				return
			}

			if channel == model.ChannelKline {
				symbol := model.ParseSymbol(sym).Normalize()
				ws.subscriberMtx.Lock()
				defer ws.subscriberMtx.Unlock()
				for subscriber := range ws.klineHandlers {
					if subscriber.SubscribeInterval().Normalize() == interval &&
						subscriber.SubscribeSymbol().Normalize() == symbol &&
						subscriber.SubscribePriceType().Normalize() == priceType {
						var klineMsg model.KLineChannelMessage
						if err := json.Unmarshal(bytes, &klineMsg); err != nil {
							if ws.logger != nil {
								ws.logger.Error("failed to unmarshal kline message", zap.Error(errors.NewInternalError("error unmarshaling kline message", err)))
							}
							return
						}

						subscriber.SubscribeKLine(&klineMsg)
					}
				}
			}
		}
	}
}

type heartbeatMessage struct {
	Op   string `json:"op"`
	Ping int64  `json:"ping"`
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

type SubscribeKLineRequest struct {
	Symbol string `json:"symbol"`
	Ch     string `json:"ch"`
}

type SubscribeRequest struct {
	Op   string        `json:"op"`
	Args []interface{} `json:"args"`
}

type KLineSubscriber interface {
	SubscribeKLine(*model.KLineChannelMessage)
	SubscribeInterval() model.Interval
	SubscribeSymbol() model.Symbol
	SubscribePriceType() model.PriceType
}

type PublicWebsocketClient interface {
	Stream() error
	Connect() error
	Disconnect()
	SubscribeKLine(subscriber KLineSubscriber) error
	UnsubscribeKLine(subscriber KLineSubscriber) error
}
