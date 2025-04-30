package bitunix

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/websocket"
	"strings"
	"time"
)

type wsClientInterface interface {
	Connect() error
	Listen(callback websocket.HandlerFunc) error
	Close()
	Write(bytes []byte) error
}

type websocketClient struct {
	client wsClientInterface
	uri    string
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

type publicWebsocketClient struct {
	websocketClient
	klineHandlers map[string]KLineHandler
}

type WebsocketClientOption func(*websocketClient)

func NewPublicWebsocket(ctx context.Context, options ...WebsocketClientOption) PublicWebsocketClient {
	wsc := websocketClient{
		uri: "wss://fapi.bitunix.com/public/",
	}
	for _, option := range options {
		option(&wsc)
	}

	wsc.client = websocket.New(
		ctx,
		wsc.uri,
		websocket.WithKeepAliveMonitor(30*time.Second, KeepAliveMonitor()),
	)

	return &publicWebsocketClient{
		websocketClient: wsc,
		klineHandlers:   make(map[string]KLineHandler),
	}
}

func (ws *publicWebsocketClient) SubscribeKLine(symbol string, interval string, priceType string, handler KLineHandler) error {
	if handler == nil {
		return fmt.Errorf("handler cannot be nil")
	}

	channelName := fmt.Sprintf("%s_kline_%s", priceType, interval)
	handlerKey := fmt.Sprintf("%s_%s", symbol, channelName)

	ws.klineHandlers[handlerKey] = handler

	req := SubscribeRequest{
		Op: "subscribe",
		Args: []interface{}{
			SubscribeKLineRequest{
				Symbol: symbol,
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

func (ws *publicWebsocketClient) UnsubscribeKLine(symbol string, interval string, priceType string) error {
	channelName := fmt.Sprintf("%s_kline_%s", priceType, interval)
	handlerKey := fmt.Sprintf("%s_%s", symbol, channelName)

	delete(ws.klineHandlers, handlerKey)

	req := SubscribeRequest{
		Op: "unsubscribe",
		Args: []interface{}{
			SubscribeKLineRequest{
				Symbol: symbol,
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

func (ws *publicWebsocketClient) Stream() error {
	err := ws.client.Listen(func(bytes []byte) error {
		go func() {
			var result map[string]interface{}

			err := json.Unmarshal(bytes, &result)
			if err != nil {
				log.WithError(fmt.Errorf("error unmarshaling JSON: %v", err)).Errorf("error unmarshaling JSON")
				return
			}

			if ch, ok := result["ch"].(string); ok {
				if symbol, symbolOk := result["symbol"].(string); symbolOk {
					if strings.Contains(ch, "kline") {
						handlerKey := fmt.Sprintf("%s_%s", symbol, ch)
						if handler, exists := ws.klineHandlers[handlerKey]; exists {
							var klineMsg model.KLineChannelMessage
							if err := json.Unmarshal(bytes, &klineMsg); err != nil {
								log.WithError(fmt.Errorf("error unmarshaling kline message: %v", err)).Errorf("error unmarshaling kline message")
								return
							}
							handler(&klineMsg)
						}
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

type KLineHandler func(msg *model.KLineChannelMessage)

type PublicWebsocketClient interface {
	Stream() error
	Connect() error
	Disconnect()
	SubscribeKLine(symbol string, interval string, priceType string, handler KLineHandler) error
	UnsubscribeKLine(symbol string, interval string, priceType string) error
}
