package bitunix

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/websocket"
	"time"
)

type websocketClient struct {
	client *websocket.Client
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
}

type WebsocketClientOption func(*websocketClient)

func NewPublicWebsocket(ctx context.Context, apiKey, secretKey string, options ...WebsocketClientOption) PublicWebsocketClient {
	wsc := websocketClient{
		uri: "wss://fapi.bitunix.com/public/",
	}
	for _, option := range options {
		option(&wsc)
	}

	wsc.client = newWebsocket(ctx, apiKey, secretKey, wsc.uri)

	return &publicWebsocketClient{
		websocketClient: wsc,
	}
}

func newWebsocket(ctx context.Context, apiKey string, secretKey string, uri string) *websocket.Client {
	ws := websocket.New(
		ctx,
		uri,
		websocket.WithAuthentication(WebsocketSigner(apiKey, secretKey)),
		websocket.WithKeepAliveMonitor(30*time.Second, KeepAliveMonitor()),
	)
	return ws
}

func (ws *publicWebsocketClient) Stream() error {
	err := ws.client.Listen(func(bytes []byte) error {
		go func() {
			var result map[string]interface{}

			err := json.Unmarshal(bytes, &result)
			if err != nil {
				log.WithError(fmt.Errorf("error unmarshaling JSON: %v", err)).Errorf("error unmarshaling JSON")
			}

			if _, ok := result["ch"].(string); ok {

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

type PublicWebsocketClient interface {
	Stream() error
	Connect() error
	Disconnect()
}
