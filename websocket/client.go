package websocket

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	bitunix_errors "github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
	"go.uber.org/zap"
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
	logger                   *zap.Logger
	logLevel                 model.LogLevel
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

func WithDebug(enabled bool) ClientOption {
	return func(ws *Client) {
		if enabled {
			ws.logLevel = model.LogLevelAggressive
		} else {
			ws.logLevel = model.LogLevelNone
		}
	}
}

func WithLogLevel(level model.LogLevel) ClientOption {
	return func(ws *Client) {
		ws.logLevel = level
	}
}

type GenericMessage map[string]interface{}

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

func New(ctx context.Context, uri string, options ...ClientOption) *Client {
	ctx, cancel := context.WithCancel(ctx)

	ws := &Client{
		wsURL:    uri,
		done:     make(chan struct{}),
		ctx:      ctx,
		cancel:   cancel,
		logLevel: model.LogLevelAggressive,
	}

	for _, option := range options {
		option(ws)
	}

	ws.logger = createLoggerForLevel(ws.logLevel)

	return ws
}

func (ws *Client) Connect() error {
	if ws.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
		ws.logger.Debug("initiating websocket connection", zap.String("url", ws.wsURL))
	}

	u, err := url.Parse(ws.wsURL)
	if err != nil {
		return bitunix_errors.NewInternalError("error parsing WebSocket URL", err)
	}

	if ws.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
		ws.logger.Debug("dialing websocket", zap.String("parsed_url", u.String()))
	}

	conn, _, err := websocket.Dial(ws.ctx, u.String(), &websocket.DialOptions{
		HTTPClient: http.DefaultClient,
	})
	if err != nil {
		switch {
		case errors.Is(ws.ctx.Err(), context.Canceled):
			return bitunix_errors.NewConnectionClosedError("listen", "context cancelled", ws.ctx.Err())
		case errors.Is(ws.ctx.Err(), context.DeadlineExceeded):
			return bitunix_errors.NewTimeoutError("websocket connection", "", ws.ctx.Err())
		}

		return bitunix_errors.NewWebsocketError("connect", "error connecting to WebSocket", err)
	}

	if ws.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
		ws.logger.Debug("websocket connection established successfully")
	}

	ws.conn = conn

	var initialMsg GenericMessage
	if err := wsjson.Read(ws.ctx, conn, &initialMsg); err != nil {
		conn.Close(websocket.StatusInternalError, "")
		return bitunix_errors.NewWebsocketError("initial handshake", "error reading initial message", err)
	}

	if ws.logLevel.ShouldLog(model.LogLevelAggressive) {
		ws.logger.Debug("received initial message", zap.Any("payload", initialMsg))
	}
	if ws.generateLoginMessage != nil {
		if ws.logLevel.ShouldLog(model.LogLevelAggressive) {
			ws.logger.Debug("authentication required, initiating login sequence")
		}
		if err := ws.login(); err != nil {
			closeErr := conn.Close(websocket.StatusInternalError, "")
			if closeErr != nil {
				return bitunix_errors.NewWebsocketError(
					"login and connection closure",
					"login failed and connection could not be closed properly",
					err,
				)
			}
			return err
		}
	}

	if ws.heartBeatInterval > 0 {
		if ws.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
			ws.logger.Debug("starting heartbeat routine", zap.Duration("interval", ws.heartBeatInterval))
		}
		go ws.sendHeartbeat()
	}

	return nil
}

func (ws *Client) Close() {
	if ws.conn != nil {
		ws.logger.Error("closing websocket connection")

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
		return bitunix_errors.NewWebsocketError("write", "connection not established", nil)
	}

	if ws.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
		ws.logger.Debug("write to websocket", zap.Int("message_size", len(bytes)), zap.String("payload", string(bytes)))
	}

	if err := ws.conn.Write(ws.ctx, websocket.MessageText, bytes); err != nil {
		ws.logger.Error("failed to write to websocket", zap.Error(err), zap.Int("message_size", len(bytes)))
		return bitunix_errors.NewWebsocketError("write", "error writing to websocket", err)
	}

	if ws.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
		ws.logger.Debug("message written to websocket successfully")
	}

	return nil
}

type HandlerFunc func([]byte) error

func (ws *Client) Listen(handler HandlerFunc) error {
	if ws.conn == nil {
		return bitunix_errors.NewWebsocketError("listen", "connection not established", nil)
	}

	if ws.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
		ws.logger.Debug("starting message listening loop")
	}

	for {
		select {
		case <-ws.done:
			if ws.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
				ws.logger.Debug("listen loop terminated via done channel")
			}
			return nil
		default:
			if ws.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
				ws.logger.Debug("waiting for incoming message")
			}

			var message json.RawMessage
			err := wsjson.Read(ws.ctx, ws.conn, &message)
			if err != nil {
				if ws.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
					ws.logger.Debug("error reading from websocket", zap.Error(err))
				}
				switch {
				case errors.Is(ws.ctx.Err(), context.Canceled):
					return bitunix_errors.NewConnectionClosedError("listen", "context cancelled", ws.ctx.Err())
				case errors.Is(ws.ctx.Err(), context.DeadlineExceeded):
					return bitunix_errors.NewTimeoutError("websocket connection", "", ws.ctx.Err())
				}

				return bitunix_errors.NewWebsocketError("listen", "connection closed", err)
			}

			if ws.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
				ws.logger.Debug("received message from websocket", zap.Int("message_size", len(message)))
			}

			if handler != nil {
				if ws.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
					ws.logger.Debug("invoking message handler")
				}
				if err := handler(message); err != nil {
					if ws.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
						ws.logger.Error("message handler failed", zap.Error(err))
					}
					return bitunix_errors.NewWebsocketError("message handling", "handler failed", err)
				}
				if ws.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
					ws.logger.Debug("message handler completed successfully")
				}
			}
		}
	}
}

func (ws *Client) login() error {
	loginReq, err := ws.generateLoginMessage()

	if err != nil {
		return bitunix_errors.NewAuthenticationError("error generating nonce for login request", err)
	}

	ws.logger.Debug("sending login message")
	if err := ws.Write(loginReq); err != nil {
		return bitunix_errors.NewWebsocketError("login", "error sending login request", err)
	}

	var loginResp GenericMessage
	if err := wsjson.Read(ws.ctx, ws.conn, &loginResp); err != nil {
		return bitunix_errors.NewWebsocketError("login", "error reading login response", err)
	}
	if op, ok := loginResp["op"].(string); ok && op == "login" {
		data := loginResp["data"].(map[string]interface{})
		if result, ok := data["result"].(bool); ok && result == true {

			ws.logger.Debug("received login response", zap.Any("payload", loginResp))
			return nil
		}
	}

	return bitunix_errors.NewAuthenticationError(fmt.Sprintf("authentication failed"), nil)
}

func (ws *Client) sendHeartbeat() {
	ticker := time.NewTicker(ws.heartBeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			heartbeat, err := ws.generateHeartbeatMessage()
			if err != nil {
				ws.logger.Error("error generating heartbeat message", zap.Error(err))
				ws.Close()
				return
			}

			if ws.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
				ws.logger.Debug("sending ping message")
			}

			err = ws.Write(heartbeat)
			if err != nil {
				ws.logger.Error("writing heartbeat message", zap.Error(err))
				ws.Close()
				return
			}

		case <-ws.done:
			return
		}
	}
}
