package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tradingiq/bitunix-client/errors"
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
	}
	for _, option := range options {
		option(c)
	}

	return c, nil
}

type APIResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message,omitempty"`
	Msg     string          `json:"msg,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}

func (c *client) request(ctx context.Context, method, path string, query url.Values, bodyBytes []byte) ([]byte, error) {
	reqURL := *c.baseUri
	reqURL.Path = path

	if query != nil {
		reqURL.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, errors.NewInternalError("failed to create request", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if err := c.signRequest(req, bodyBytes); err != nil {
		return nil, errors.NewAuthenticationError("failed to sign request", err)
	}

	c.logRequest(req, bodyBytes)

	resp, err := c.httpClient.Do(req)
	if err != nil {

		if ctx.Err() != nil {
			return nil, errors.NewTimeoutError("HTTP request", "", ctx.Err())
		}
		return nil, errors.NewNetworkError("HTTP request", "failed to send request", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewNetworkError("reading response", "failed to read response body", err)
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
	logging := log.WithField("method", req.Method).WithField("uri", req.URL.String())

	for k, v := range req.Header {
		logging.WithField(k, strings.Join(v, ","))
	}
	if len(body) > 0 {
		logging.WithField("body", string(body))
	}

	logging.Debug("request")
}

func (c *client) logResponse(resp *http.Response, body []byte) {
	logging := log.WithField("status_code", resp.StatusCode)

	for k, v := range resp.Header {
		logging.WithField(k, strings.Join(v, ","))
	}
	if len(body) > 0 {
		logging.WithField("body", string(body))
	}

	logging.Debug("response")
}
