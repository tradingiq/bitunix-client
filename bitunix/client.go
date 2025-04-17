package bitunix

import (
	"bitunix-client/api"
	"bitunix-client/security"
	"bitunix-client/util"
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Client struct {
	api                          *api.Client
	generateNonce                func(int) ([]byte, error)
	generateMillisecondTimestamp func() int64
}

func New(rest *api.Client, apiKey, apiString string) *Client {
	client := &Client{
		api:                          rest,
		generateNonce:                security.GenerateNonce,
		generateMillisecondTimestamp: util.CurrentTimestampMillis,
	}

	rest.SetOptions(api.WithRequestSigner(client.requestSigner(apiKey, apiString)))

	return client
}

func (client *Client) requestSigner(apiKey string, apiSecret string) func(req *http.Request, body []byte) error {
	return func(req *http.Request, body []byte) error {
		ts := client.generateMillisecondTimestamp()
		timestamp := strconv.FormatInt(ts, 10)

		randomBytes, err := client.generateNonce(32)
		if err != nil {
			return fmt.Errorf("failed to generate nonce: %w", err)
		}
		nonce := base64.StdEncoding.EncodeToString(randomBytes)

		queryParams := req.URL.RawQuery
		queryParams = strings.ReplaceAll(queryParams, "&", "")
		queryParams = strings.ReplaceAll(queryParams, "=", "")

		bodyStr := string(body)
		digestInput := nonce + timestamp + apiKey + queryParams + bodyStr

		digest := security.Sha256Hex(digestInput)

		signInput := digest + apiSecret
		signature := security.Sha256Hex(signInput)

		req.Header.Set("api-key", apiKey)
		req.Header.Set("sign", signature)
		req.Header.Set("timestamp", timestamp)
		req.Header.Set("nonce", nonce)
		req.Header.Set("language", "en-US")

		return nil
	}
}
