package bitunix

import (
	"context"
	"encoding/json"
	"errors"
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
	klineHandlers map[KLineSubscriber]struct{}
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
		klineHandlers:   make(map[KLineSubscriber]struct{}),
	}
}

func (ws *publicWebsocketClient) SubscribeKLine(subscriber KLineSubscriber) error {
	if subscriber == nil {
		return fmt.Errorf("subscriber cannot be nil")
	}

	symbol := subscriber.Symbol().Normalize()
	priceType := subscriber.PriceType().Normalize()
	interval := subscriber.Interval().Normalize()

	channelName := fmt.Sprintf("%s_kline_%s", priceType, interval)

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

	symbol := subscriber.Symbol().Normalize()
	priceType := subscriber.PriceType().Normalize()
	interval := subscriber.Interval().Normalize()

	channelName := fmt.Sprintf("%s_kline_%s", priceType, interval)

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

func parseChannel(channelStr string) (model.Interval, model.Channel, error) {
	parts := strings.Split(channelStr, "_")

	if len(parts) != 3 {
		return "", "", errors.New("invalid channel format: expected priceType_type_interval")
	}

	channelStr = parts[1]
	intervalStr := parts[2]

	interval, err := model.ParseInterval(intervalStr)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse interval: %w", err)
	}

	channel, err := model.ParseChannel(channelStr)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse channel: %w", err)
	}

	return interval, channel, nil
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
				if sym, symbolOk := result["symbol"].(string); symbolOk {
					if strings.Contains(ch, "kline") {
						symbol := model.ParseSymbol(sym).Normalize()
						interval, _, err := parseChannel(ch)
						if err != nil {
							log.WithError(err).Errorf("error parsing channel")
							return
						}
						for subscriber := range ws.klineHandlers {
							if subscriber.Interval() == interval && subscriber.Symbol() == symbol {
								var klineMsg model.KLineChannelMessage
								if err := json.Unmarshal(bytes, &klineMsg); err != nil {
									log.WithError(fmt.Errorf("error unmarshaling kline message: %v", err)).Errorf("error unmarshaling kline message")
									return
								}
								subscriber.Handle(&klineMsg)
							}
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

type KLineSubscriber interface {
	Handle(*model.KLineChannelMessage)
	Interval() model.Interval
	Symbol() model.Symbol
	PriceType() model.PriceType
}

type PublicWebsocketClient interface {
	Stream() error
	Connect() error
	Disconnect()
	SubscribeKLine(handler KLineSubscriber) error
	UnsubscribeKLine(subscriber KLineSubscriber) error
}
