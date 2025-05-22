package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client interface {
	Get(ctx context.Context, path string, query url.Values) ([]byte, error)
	Post(ctx context.Context, path string, query url.Values, body []byte) ([]byte, error)
}

type client struct {
	httpClient  *http.Client
	signRequest func(req *http.Request, body []byte) error
	baseUri     *url.URL
	logger      *zap.Logger
	logLevel    model.LogLevel
}

type ClientOption func(*client)

func WithRequestSigner(requestSigner func(req *http.Request, body []byte) error) ClientOption {
	return func(c *client) {
		c.signRequest = requestSigner
	}
}

func WithDefaultTimeout(timeout time.Duration) ClientOption {
	return func(c *client) {
		c.httpClient.Timeout = timeout
	}
}

func WithDebug(enabled bool) ClientOption {
	return func(c *client) {
		if enabled {
			c.logLevel = model.LogLevelAggressive
		} else {
			c.logLevel = model.LogLevelNone
		}
	}
}

func WithLogLevel(level model.LogLevel) ClientOption {
	return func(c *client) {
		c.logLevel = level
	}
}

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

func New(baseUri string, options ...ClientOption) (Client, error) {
	uri, err := url.Parse(baseUri)
	if err != nil {
		return nil, errors.NewInternalError(
			fmt.Sprintf("error parsing base uri %q", baseUri),
			err,
		)
	}

	c := &client{
		httpClient: &http.Client{},
		baseUri:    uri,
		logLevel:   model.LogLevelAggressive,
	}
	for _, option := range options {
		option(c)
	}

	c.logger = createLoggerForLevel(c.logLevel)

	return c, nil
}

type APIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message,omitempty"`
	Msg     string          `json:"msg,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (c *client) request(ctx context.Context, method, path string, query url.Values, bodyBytes []byte) ([]byte, error) {
	if c.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
		c.logger.Debug("initiating HTTP request",
			zap.String("method", method),
			zap.String("path", path),
			zap.Int("body_size", len(bodyBytes)))
	}

	reqURL := *c.baseUri
	reqURL.Path = path

	if query != nil {
		reqURL.RawQuery = query.Encode()
		if c.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
			c.logger.Debug("request query parameters", zap.String("query", reqURL.RawQuery))
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, errors.NewInternalError("failed to create request", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if c.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
		c.logger.Debug("signing request with authentication")
	}

	if err := c.signRequest(req, bodyBytes); err != nil {
		return nil, errors.NewAuthenticationError("failed to sign request", err)
	}

	c.logRequest(req, bodyBytes)

	if c.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
		c.logger.Debug("sending HTTP request", zap.String("url", req.URL.String()))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if c.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
			c.logger.Error("HTTP request failed", zap.Error(err))
		}

		if ctx.Err() != nil {
			return nil, errors.NewTimeoutError("HTTP request", "", ctx.Err())
		}
		return nil, errors.NewNetworkError("HTTP request", "failed to send request", err)
	}
	defer resp.Body.Close()

	if c.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
		c.logger.Debug("HTTP request completed", zap.Int("status_code", resp.StatusCode))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewNetworkError("reading response", "failed to read response body", err)
	}

	if c.logLevel.ShouldLog(model.LogLevelVeryAggressive) {
		c.logger.Debug("response body read", zap.Int("response_size", len(respBody)))
	}

	c.logResponse(resp, respBody)

	if resp.StatusCode >= 400 {

		var apiResp APIResponse
		if err := json.Unmarshal(respBody, &apiResp); err == nil {

			message := apiResp.Message
			if message == "" {
				message = apiResp.Msg
			}

			return nil, errors.NewAPIError(
				apiResp.Code,
				message,
				path,
				fmt.Errorf("HTTP status: %s (%d)", resp.Status, resp.StatusCode),
			)
		}

		return nil, errors.NewAPIError(
			resp.StatusCode,
			string(respBody),
			path,
			fmt.Errorf("HTTP status: %s", resp.Status),
		)
	}

	return respBody, nil
}

func (c *client) Get(ctx context.Context, path string, query url.Values) ([]byte, error) {
	respBody, err := c.request(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (c *client) Post(ctx context.Context, path string, query url.Values, body []byte) ([]byte, error) {
	respBody, err := c.request(ctx, http.MethodPost, path, query, body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (c *client) logRequest(req *http.Request, body []byte) {
	fields := []zap.Field{
		zap.String("method", req.Method),
		zap.String("uri", req.URL.String()),
	}

	for k, v := range req.Header {
		fields = append(fields, zap.String(k, strings.Join(v, ",")))
	}
	if len(body) > 0 {
		fields = append(fields, zap.String("body", string(body)))
	}

	c.logger.Debug("request", fields...)
}

func (c *client) logResponse(resp *http.Response, body []byte) {
	fields := []zap.Field{
		zap.Int("status_code", resp.StatusCode),
	}

	for k, v := range resp.Header {
		fields = append(fields, zap.String(k, strings.Join(v, ",")))
	}
	if len(body) > 0 {
		fields = append(fields, zap.String("body", string(body)))
	}

	c.logger.Debug("response", fields...)
}
