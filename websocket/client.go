package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	log "github.com/sirupsen/logrus"
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
	ctx, cancel := context.WithCancel(context.Background())

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
		return fmt.Errorf("error parsing WebSocket URL: %w", err)
	}

	conn, _, err := websocket.Dial(ws.ctx, u.String(), &websocket.DialOptions{
		HTTPClient: http.DefaultClient,
	})
	if err != nil {
		return fmt.Errorf("error connecting to WebSocket: %w", err)
	}

	ws.conn = conn

	var initialMsg GenericMessage
	if err := wsjson.Read(ws.ctx, conn, &initialMsg); err != nil {
		conn.Close(websocket.StatusInternalError, "")
		return fmt.Errorf("error reading initial message: %w", err)
	}
	log.WithField("payload", initialMsg).Debug("received initial message")

	if ws.generateLoginMessage != nil {
		if err := ws.login(); err != nil {
			if err := conn.Close(websocket.StatusInternalError, ""); err != nil {
				log.WithError(err).Error("error closing websocket connection")
			}
			return fmt.Errorf("login failed: %w", err)
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
		return fmt.Errorf("not connected")
	}

	log.WithField("payload", string(bytes)).Debug("write to websocket")
	if err := ws.conn.Write(ws.ctx, websocket.MessageText, bytes); err != nil {
		return fmt.Errorf("error writing: %w", err)
	}

	return nil
}

type HandlerFunc func([]byte) error

func (ws *Client) Listen(handler HandlerFunc) error {
	if ws.conn == nil {
		return fmt.Errorf("not connected")
	}

	for {
		select {
		case <-ws.done:
			return nil
		default:
			var message json.RawMessage
			err := wsjson.Read(ws.ctx, ws.conn, &message)
			if err != nil {
				return fmt.Errorf("connection closed: %w", err)
			}

			log.WithField("payload", string(message)).Debug("received message")

			if handler != nil {
				if err := handler(message); err != nil {
					return fmt.Errorf("handler failed: %w", err)
				}
			}
		}
	}
}

func (ws *Client) login() error {
	loginReq, err := ws.generateLoginMessage()

	if err != nil {
		return fmt.Errorf("error generating nonce for login request: %w", err)
	}

	log.Debug("sending login message")
	if err := ws.Write(loginReq); err != nil {
		return fmt.Errorf("error sending login request: %w", err)
	}

	var loginResp GenericMessage
	if err := wsjson.Read(ws.ctx, ws.conn, &loginResp); err != nil {
		return fmt.Errorf("error reading login response: %w", err)
	}

	log.WithField("payload", loginResp).Debug("received login response")

	return nil
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
