package bitunix

import (
	"github.com/tradingiq/bitunix-client/rest"
	"github.com/tradingiq/bitunix-client/security"
	"time"
)

func generateTimestamp() int64 { return time.Now().UnixMilli() }

type client struct {
	restClient *rest.Client
}

func New(restClient *rest.Client, apiKey, apiSecret string) *client {
	client := &client{
		restClient: restClient,
	}

	restClient.SetOptions(rest.WithRequestSigner(RequestSigner(apiKey, apiSecret, generateTimestamp, security.GenerateNonce)))

	return client
}
