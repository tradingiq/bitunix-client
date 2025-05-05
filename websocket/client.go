package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/errors"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	conn                     *websocket.Conn
	wsURL                    string
	done                     chan struct{}
	ctx                      context.Context
	cancel                   context.CancelFunc
	heartBeatInterval        time.Duration
	generateHeartbeatMessage func() ([]byte, error)
	generateLoginMessage     func() ([]byte, error)
}

type ClientOption func(*Client)

func WithAuthentication(loginMessageGenerator func() ([]byte, error)) ClientOption {
	return func(ws *Client) {
		ws.generateLoginMessage = loginMessageGenerator

	}
}

func WithKeepAliveMonitor(interval time.Duration, messageGenerator func() ([]byte, error)) ClientOption {
	return func(ws *Client) {
		ws.heartBeatInterval = interval
		ws.generateHeartbeatMessage = messageGenerator
	}
}

type GenericMessage map[string]interface{}

func New(ctx context.Context, uri string, options ...ClientOption) *Client {
	ctx, cancel := context.WithCancel(ctx)

	ws := &Client{
		wsURL:  uri,
		done:   make(chan struct{}),
		ctx:    ctx,
		cancel: cancel,
	}

	for _, option := range options {
		option(ws)
	}

	return ws
}

func (ws *Client) Connect() error {
	u, err := url.Parse(ws.wsURL)
	if err != nil {
		return errors.NewInternalError("error parsing WebSocket URL", err)
	}

	conn, _, err := websocket.Dial(ws.ctx, u.String(), &websocket.DialOptions{
		HTTPClient: http.DefaultClient,
	})
	if err != nil {
		if ws.ctx.Err() != nil {
			return errors.NewTimeoutError("websocket connection", "", ws.ctx.Err())
		}
		return errors.NewWebsocketError("connect", "error connecting to WebSocket", err)
	}

	ws.conn = conn

	var initialMsg GenericMessage
	if err := wsjson.Read(ws.ctx, conn, &initialMsg); err != nil {
		conn.Close(websocket.StatusInternalError, "")
		return errors.NewWebsocketError("initial handshake", "error reading initial message", err)
	}
	log.WithField("payload", initialMsg).Debug("received initial message")

	if ws.generateLoginMessage != nil {
		if err := ws.login(); err != nil {
			closeErr := conn.Close(websocket.StatusInternalError, "")
			if closeErr != nil {
				log.WithError(closeErr).Error("error closing websocket connection")
				return errors.NewWebsocketError(
					"login and connection closure",
					"login failed and connection could not be closed properly",
					err,
				)
			}
			return err
		}
	}

	if ws.heartBeatInterval > 0 {
		go ws.sendHeartbeat()
	}

	return nil
}

func (ws *Client) Close() {
	if ws.conn != nil {
		log.Error("closing websocket connection")

		select {
		case <-ws.done:

		default:
			close(ws.done)
		}

		ws.cancel()
		ws.conn.Close(websocket.StatusNormalClosure, "")
	}
}

func (ws *Client) Write(bytes []byte) error {
	if ws.conn == nil {
		return errors.NewWebsocketError("write", "connection not established", nil)
	}

	log.WithField("payload", string(bytes)).Debug("write to websocket")
	if err := ws.conn.Write(ws.ctx, websocket.MessageText, bytes); err != nil {
		return errors.NewWebsocketError("write", "error writing to websocket", err)
	}

	return nil
}

type HandlerFunc func([]byte) error

func (ws *Client) Listen(handler HandlerFunc) error {
	if ws.conn == nil {
		return errors.NewWebsocketError("listen", "connection not established", nil)
	}

	for {
		select {
		case <-ws.done:
			return nil
		default:
			var message json.RawMessage
			err := wsjson.Read(ws.ctx, ws.conn, &message)
			if err != nil {
				if ws.ctx.Err() != nil {
					return errors.NewTimeoutError("websocket read", "", ws.ctx.Err())
				}
				return errors.NewWebsocketError("listen", "connection closed", err)
			}

			log.WithField("payload", string(message)).Debug("received message")

			if handler != nil {
				if err := handler(message); err != nil {
					return errors.NewWebsocketError("message handling", "handler failed", err)
				}
			}
		}
	}
}

func (ws *Client) login() error {
	loginReq, err := ws.generateLoginMessage()

	if err != nil {
		return errors.NewAuthenticationError("error generating nonce for login request", err)
	}

	log.Debug("sending login message")
	if err := ws.Write(loginReq); err != nil {
		return errors.NewWebsocketError("login", "error sending login request", err)
	}

	var loginResp GenericMessage
	if err := wsjson.Read(ws.ctx, ws.conn, &loginResp); err != nil {
		return errors.NewWebsocketError("login", "error reading login response", err)
	}
	if op, ok := loginResp["op"].(string); ok && op == "login" {
		data := loginResp["data"].(map[string]interface{})
		if result, ok := data["result"].(bool); ok && result == true {

			log.WithField("payload", loginResp).Debug("received login response")
			return nil
		}
	}

	return errors.NewAuthenticationError(fmt.Sprintf("authentication failed"), nil)
}

func (ws *Client) sendHeartbeat() {
	ticker := time.NewTicker(ws.heartBeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			heartbeat, err := ws.generateHeartbeatMessage()
			if err != nil {
				log.WithField("error", err).Error("error generating heartbeat message")
				ws.Close()
				return
			}
			log.Debug("sending ping message")

			err = ws.Write(heartbeat)
			if err != nil {
				log.WithField("error", err).Error("writing heartbeat message")
				ws.Close()
				return
			}

		case <-ws.done:
			return
		}
	}
}
