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
	GetPendingPositions(ctx context.Context, params model.PendingPositionParams) (*model.PendingPositionResponse, error)
	GetOrderDetail(ctx context.Context, request *OrderDetailRequest) (*model.OrderDetailResponse, error)
}

func generateTimestamp() int64 { return time.Now().UnixMilli() }

type apiClient struct {
	restClient rest.Client
	baseURI    string
	logLevel   model.LogLevel
}

type ClientOption func(*apiClient)

func WithBaseURI(uri string) ClientOption {
	return func(c *apiClient) {
		c.baseURI = uri
	}
}

func WithLogLevel(level model.LogLevel) ClientOption {
	return func(c *apiClient) {
		c.logLevel = level
	}
}

func NewApiClient(apiKey, apiSecret string, option ...ClientOption) (ApiClient, error) {
	client := &apiClient{
		baseURI:  "https://fapi.bitunix.com/",
		logLevel: model.LogLevelNone,
	}
	for _, option := range option {
		option(client)
	}

	restOptions := []rest.ClientOption{
		rest.WithRequestSigner(RequestSigner(apiKey, apiSecret, generateTimestamp, security.GenerateNonce)),
		rest.WithLogLevel(client.logLevel),
	}

	restClient, err := rest.New(client.baseURI, restOptions...)
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
			case response.Code == 10001:
				underlyingErr = errors.ErrNetwork
			case response.Code == 10002:
				underlyingErr = errors.ErrParameterError
			case response.Code == 10003:
				underlyingErr = errors.ErrAuthentication
			case response.Code == 10004:
				underlyingErr = errors.ErrIPNotAllowed
			case response.Code == 10005 || response.Code == 10006:
				underlyingErr = errors.ErrRateLimitExceeded
			case response.Code == 10007:
				underlyingErr = errors.ErrSignatureError
			case response.Code == 10008:
				underlyingErr = errors.ErrInvalidValue

			case response.Code == 20001:
				underlyingErr = errors.ErrMarketNotExists
			case response.Code == 20002:
				underlyingErr = errors.ErrPositionLimitExceeded
			case response.Code == 20003 || response.Code == 20008:
				underlyingErr = errors.ErrInsufficientBalance
			case response.Code == 20004:
				underlyingErr = errors.ErrInsufficientTrader
			case response.Code == 20005:
				underlyingErr = errors.ErrInvalidLeverage
			case response.Code == 20006:
				underlyingErr = errors.ErrOpenOrdersExist
			case response.Code == 20007:
				underlyingErr = errors.ErrOrderNotFound
			case response.Code == 20009:
				underlyingErr = errors.ErrPositionsModeChange
			case response.Code == 20010:
				underlyingErr = errors.ErrInsufficientBalance
			case response.Code == 20011:
				underlyingErr = errors.ErrAccountNotAllowed
			case response.Code == 20012 || response.Code == 20015:
				underlyingErr = errors.ErrFuturesNotSupported
			case response.Code == 20013 || response.Code == 20014:
				underlyingErr = errors.ErrAccountInactive

			case response.Code >= 30001 && response.Code <= 30003:
				underlyingErr = errors.ErrOrderPriceIssue
			case response.Code == 30004:
				underlyingErr = errors.ErrPositionNotExist
			case response.Code >= 30005 && response.Code <= 30038:
				underlyingErr = errors.ErrTPSLOrderError
			case response.Code == 30039:
				underlyingErr = errors.ErrOrderQuantityIssue
			case response.Code == 30041:
				underlyingErr = errors.ErrTriggerPriceInvalid
			case response.Code == 30042:
				underlyingErr = errors.ErrDuplicateClientID

			case response.Code >= 40001 && response.Code <= 40004:
				underlyingErr = errors.ErrLeadTrading
			case response.Code >= 40005 && response.Code <= 40008:
				underlyingErr = errors.ErrSubAccountIssue

			default:
				underlyingErr = errors.UnknownAPIError
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
