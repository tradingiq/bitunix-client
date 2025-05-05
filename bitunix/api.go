package bitunix

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/tradingiq/bitunix-client/errors"
	"github.com/tradingiq/bitunix-client/model"
	"github.com/tradingiq/bitunix-client/rest"
	"github.com/tradingiq/bitunix-client/security"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ApiClient interface {
	PlaceOrder(ctx context.Context, request *model.OrderRequest) (*model.OrderResponse, error)
	CancelOrders(ctx context.Context, request *model.CancelOrderRequest) (*model.CancelOrderResponse, error)
	GetTradeHistory(ctx context.Context, params model.TradeHistoryParams) (*model.TradeHistoryResponse, error)
	GetOrderHistory(ctx context.Context, params model.OrderHistoryParams) (*model.OrderHistoryResponse, error)
	GetPositionHistory(ctx context.Context, params model.PositionHistoryParams) (*model.PositionHistoryResponse, error)
	PlaceTpSlOrder(ctx context.Context, request *model.TPSLOrderRequest) (*model.TpSlOrderResponse, error)
	GetAccountBalance(ctx context.Context, params model.AccountBalanceParams) (*model.AccountBalanceResponse, error)
}

func generateTimestamp() int64 { return time.Now().UnixMilli() }

type apiClient struct {
	restClient *rest.Client
	baseURI    string
}

type ClientOption func(*apiClient)

func WithBaseURI(uri string) ClientOption {
	return func(c *apiClient) {
		c.baseURI = uri
	}
}

func NewApiClient(apiKey, apiSecret string, option ...ClientOption) (ApiClient, error) {
	client := &apiClient{
		baseURI: "https://fapi.bitunix.com/",
	}
	for _, option := range option {
		option(client)
	}

	restClient, err := rest.New(client.baseURI, rest.WithRequestSigner(RequestSigner(apiKey, apiSecret, generateTimestamp, security.GenerateNonce)))
	if err != nil {
		return nil, errors.NewInternalError("creating rest client", err)
	}

	client.restClient = restClient

	return client, nil
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

func handleAPIResponse(responseBody []byte, endpoint string, result interface{}) error {

	if err := json.Unmarshal(responseBody, result); err != nil {
		return errors.NewInternalError(fmt.Sprintf("failed to unmarshal response from %s", endpoint), err)
	}

	response := struct {
		Code    int    `json:"code"`
		Message string `json:"message,omitempty"`
		Msg     string `json:"msg,omitempty"`
	}{}

	if err := json.Unmarshal(responseBody, &response); err == nil {

		if response.Code != 0 {
			message := response.Message
			if message == "" {
				message = response.Msg
			}

			var underlyingErr error
			switch {
			case response.Code == 10002:
				underlyingErr = errors.ErrParameterError
			case response.Code == 10005 || response.Code == 10006:
				underlyingErr = errors.ErrRateLimitExceeded
			case response.Code == 10007:
				underlyingErr = errors.ErrSignatureError
			case response.Code == 20001:
				underlyingErr = errors.ErrMarketNotExists
			case response.Code == 20003 || response.Code == 20008:
				underlyingErr = errors.ErrInsufficientBalance
			case response.Code == 20005:
				underlyingErr = errors.ErrInvalidLeverage
			case response.Code == 20007:
				underlyingErr = errors.ErrOrderNotFound
			case response.Code == 20011:
				underlyingErr = errors.ErrAccountNotAllowed
			case response.Code == 30004:
				underlyingErr = errors.ErrPositionNotExist
			case response.Code >= 30005 && response.Code <= 30038:
				underlyingErr = errors.ErrTPSLOrderError
			case response.Code == 30042:
				underlyingErr = errors.ErrDuplicateClientID
			default:
				underlyingErr = errors.ErrAPI
			}

			return errors.NewAPIError(response.Code, message, endpoint, underlyingErr)
		}
	}

	return nil
}

func RequestSigner(apiKey string, apiSecret string, timestampGenerationFunc func() int64, nonceGenerationFunc func(int) ([]byte, error)) func(req *http.Request, body []byte) error {
	return func(req *http.Request, body []byte) error {
		ts := timestampGenerationFunc()

		randomBytes, err := nonceGenerationFunc(32)
		if err != nil {
			return errors.NewAuthenticationError("failed to generate nonce", err)
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
			return errors.NewAuthenticationError("failed to generate signature", err)
		}

		req.Header.Add("Api-Key", apiKey)
		req.Header.Add("Sign", signature)
		req.Header.Add("Timestamp", timestamp)
		req.Header.Add("Nonce", nonce)
		req.Header.Add("Language", "en-US")

		return nil
	}
}
