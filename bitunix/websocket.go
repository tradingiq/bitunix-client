package bitunix

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/websocket"
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
}

func (ws *websocketClient) Connect() error {
	if err := ws.client.Connect(); err != nil {
		return fmt.Errorf("websocket client failed to connect: %w", err)
	}

	return nil
}

func (ws *websocketClient) Stream() error {
	err := ws.client.Listen(func(bytes []byte) error {
		select {
		case ws.messageQueue <- bytes:
		default:
			log.Error("message queue is full, dropping message")
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("websocket client failed to listen: %w", err)
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
		return fmt.Errorf("processFunc is nil")
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
}

type WebsocketClientOption func(*websocketClient)

func NewPublicWebsocket(ctx context.Context, options ...WebsocketClientOption) (PublicWebsocketClient, error) {
	wsc := &websocketClient{
		uri:  "wss://fapi.bitunix.com/public/",
		quit: make(chan struct{}),
	}
	for _, option := range options {
		option(wsc)
	}

	wsc.client = websocket.New(
		ctx,
		wsc.uri,
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

	client := &publicWebsocketClient{
		websocketClient: wsc,
		subscriberMtx:   sync.Mutex{},
		klineHandlers:   make(map[KLineSubscriber]struct{}),
	}
	wsc.processFunc = client.processMessage

	if err := client.startWorkerPool(ctx); err != nil {
		return nil, fmt.Errorf("public websocket client failed to start worker pool: %w", err)
	}

	return client, nil
}

func (ws *publicWebsocketClient) SubscribeKLine(subscriber KLineSubscriber) error {
	if subscriber == nil {
		return fmt.Errorf("subscriber cannot be nil")
	}

	symbol := subscriber.SubscribeSymbol().Normalize()
	priceType := subscriber.SubscribePriceType().Normalize()
	interval := subscriber.SubscribeInterval().Normalize()

	channelName := fmt.Sprintf("%s_kline_%s", priceType, interval)

	ws.subscriberMtx.Lock()
	defer ws.subscriberMtx.Unlock()

	ws.klineHandlers[subscriber] = struct{}{}

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
		return fmt.Errorf("failed to marshal subscription request: %w", err)
	}

	if err := ws.client.Write(bytes); err != nil {
		return fmt.Errorf("failed to send subscription request: %w", err)
	}

	return nil
}

func (ws *publicWebsocketClient) UnsubscribeKLine(subscriber KLineSubscriber) error {
	if subscriber == nil {
		return fmt.Errorf("subscriber cannot be nil")
	}

	symbol := subscriber.SubscribeSymbol().Normalize()
	priceType := subscriber.SubscribePriceType().Normalize()
	interval := subscriber.SubscribeInterval().Normalize()

	channelName := fmt.Sprintf("%s_kline_%s", priceType, interval)

	ws.subscriberMtx.Lock()
	defer ws.subscriberMtx.Unlock()

	delete(ws.klineHandlers, subscriber)

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
		return fmt.Errorf("failed to marshal unsubscription request: %w", err)
	}

	if err := ws.client.Write(bytes); err != nil {
		return fmt.Errorf("failed to send unsubscription request: %w", err)
	}

	return nil
}

func parseChannel(channelStr string) (model.Interval, model.Channel, model.PriceType, error) {
	parts := strings.Split(channelStr, "_")

	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("invalid channel format: expected priceType_type_interval, got %s", channelStr)
	}

	priceTypeStr := parts[0]
	channelStr = parts[1]
	intervalStr := parts[2]

	priceType, err := model.ParsePriceType(priceTypeStr)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse price type: %w", err)
	}

	interval, err := model.ParseInterval(intervalStr)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse interval: %w", err)
	}

	channel, err := model.ParseChannel(channelStr)
	if err != nil {
		return "", "", "", fmt.Errorf("failed to parse channel: %w", err)
	}

	return interval, channel, priceType, nil
}

func (ws *publicWebsocketClient) processMessage(bytes []byte) {
	var result map[string]interface{}

	err := json.Unmarshal(bytes, &result)
	if err != nil {
		log.WithError(fmt.Errorf("error unmarshaling JSON: %v", err)).Errorf("error unmarshaling JSON")
		return
	}

	if ch, ok := result["ch"].(string); ok {
		if sym, symbolOk := result["symbol"].(string); symbolOk {
			interval, channel, priceType, err := parseChannel(ch)
			if err != nil {
				log.WithError(err).Errorf("error parsing channel")
				return
			}

			if channel == model.ChannelKline {
				symbol := model.ParseSymbol(sym).Normalize()

				for subscriber := range ws.klineHandlers {

					if subscriber.SubscribeInterval().Normalize() == interval &&
						subscriber.SubscribeSymbol().Normalize() == symbol &&
						subscriber.SubscribePriceType().Normalize() == priceType {
						var klineMsg model.KLineChannelMessage
						if err := json.Unmarshal(bytes, &klineMsg); err != nil {
							log.WithError(fmt.Errorf("error unmarshaling kline message: %v", err)).Errorf("error unmarshaling kline message")
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
