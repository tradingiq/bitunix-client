package rest

import (
	"bytes"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	httpClient  *http.Client
	signRequest func(req *http.Request, body []byte) error
	baseUri     *url.URL
}

type ClientOption func(*Client)

func WithRequestSigner(requestSigner func(req *http.Request, body []byte) error) ClientOption {
	return func(c *Client) {
		c.signRequest = requestSigner
	}
}

func WithDefaultTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

func (c *Client) SetOptions(options ...ClientOption) {
	for _, option := range options {
		option(c)
	}
}

func New(baseUri string, options ...ClientOption) (*Client, error) {
	uri, err := url.Parse(baseUri)
	if err != nil {
		return nil, err
	}

	client := &Client{
		httpClient: &http.Client{},
		baseUri:    uri,
	}

	client.SetOptions(options...)

	return client, nil
}

func (c *Client) Request(ctx context.Context, method, path string, query url.Values, bodyBytes []byte) ([]byte, error) {
	reqURL := *c.baseUri
	reqURL.Path = path

	if query != nil {
		reqURL.RawQuery = query.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, reqURL.String(), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if err := c.signRequest(req, bodyBytes); err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	c.logRequest(req, bodyBytes)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	c.logResponse(resp, respBody)

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error: %s (%d): %s", resp.Status, resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (c *Client) Get(ctx context.Context, path string, query url.Values) ([]byte, error) {
	respBody, err := c.Request(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (c *Client) Post(ctx context.Context, path string, query url.Values, body []byte) ([]byte, error) {
	respBody, err := c.Request(ctx, http.MethodPost, path, query, body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (c *Client) logRequest(req *http.Request, body []byte) {
	logging := log.WithField("method", req.Method).WithField("uri", req.URL.String())

	for k, v := range req.Header {
		logging.WithField(k, strings.Join(v, ","))
	}
	if len(body) > 0 {
		logging.WithField("body", string(body))
	}

	logging.Debug("request")
}

func (c *Client) logResponse(resp *http.Response, body []byte) {
	logging := log.WithField("status_code", resp.StatusCode)

	for k, v := range resp.Header {
		logging.WithField(k, strings.Join(v, ","))
	}
	if len(body) > 0 {
		logging.WithField("body", string(body))
	}

	logging.Debug("response")
}
