package bitunix

import (
	"encoding/base64"
	"fmt"
	"github.com/tradingiq/bitunix-client/rest"
	"github.com/tradingiq/bitunix-client/security"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func generateTimestamp() int64 { return time.Now().UnixMilli() }

type API struct {
	restClient *rest.Client
}

func New(restClient *rest.Client, apiKey, apiSecret string) *API {
	client := &API{
		restClient: restClient,
	}

	restClient.SetOptions(rest.WithRequestSigner(createRequestSigner(apiKey, apiSecret, generateTimestamp, security.GenerateNonce)))

	return client
}

func generateRequestSignature(apiKey, apiSecret, queryParams, bodyStr string, timestamp int64, nonceBytes []byte) (string, string, string, error) {
	timestampStr := strconv.FormatInt(timestamp, 10)

	nonce := base64.StdEncoding.EncodeToString(nonceBytes)

	queryParams = strings.ReplaceAll(queryParams, "&", "")
	queryParams = strings.ReplaceAll(queryParams, "=", "")

	digestInput := nonce + timestampStr + apiKey + queryParams + bodyStr
	digest := security.Sha256Hex(digestInput)

	signInput := digest + apiSecret
	signature := security.Sha256Hex(signInput)

	return signature, timestampStr, nonce, nil
}

func createRequestSigner(apiKey string, apiSecret string, timestampGenerationFunc func() int64, nonceGenerationFunc func(int) ([]byte, error)) func(req *http.Request, body []byte) error {
	return func(req *http.Request, body []byte) error {
		ts := timestampGenerationFunc()

		randomBytes, err := nonceGenerationFunc(32)
		if err != nil {
			return fmt.Errorf("failed to generate nonce: %w", err)
		}

		signature, timestamp, nonce, err := generateRequestSignature(
			apiKey,
			apiSecret,
			req.URL.RawQuery,
			string(body),
			ts,
			randomBytes,
		)
		if err != nil {
			return fmt.Errorf("failed to generate signature: %w", err)
		}

		req.Header.Add("Api-Key", apiKey)
		req.Header.Add("Sign", signature)
		req.Header.Add("Timestamp", timestamp)
		req.Header.Add("Nonce", nonce)
		req.Header.Add("Language", "en-US")

		return nil
	}
}
