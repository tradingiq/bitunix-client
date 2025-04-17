package bitunix

import (
	"bitunix-client/api"
	"bitunix-client/security"
	"bitunix-client/util"
	"encoding/base64"
	"fmt"
	"net/http"
)

type Client struct {
	api *api.Client
}

func New(rest *api.Client, apiKey, apiString string) *Client {
	client := &Client{
		api: rest,
	}

	rest.SetOptions(api.WithRequestSigner(client.requestSigner(apiKey, apiString)))

	return client
}

func (client *Client) requestSigner(apiKey string, apiSecret string) func(req *http.Request, body []byte) error {
	return func(req *http.Request, body []byte) error {
		timestamp := util.CurrentTimestampMillisString()

		randomBytes, err := security.GenerateNonce(32)
		if err != nil {
			return fmt.Errorf("failed to generate nonce: %w", err)
		}
		nonce := base64.StdEncoding.EncodeToString(randomBytes)

		queryParams := req.URL.RawQuery

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
