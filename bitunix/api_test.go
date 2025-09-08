package bitunix

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/tradingiq/bitunix-client/security"
)

func MockNonceGenerator(bytes []byte) func(int) ([]byte, error) {
	return func(length int) ([]byte, error) {
		return bytes[:length], nil
	}
}

func MockFailingNonceGenerator(err error) func(int) ([]byte, error) {
	return func(length int) ([]byte, error) {
		return nil, err
	}
}

func MockMillisecondTimestampGenerator(timestamp int64) func() int64 {
	return func() int64 {
		return timestamp
	}
}

func TestGenerateSignature(t *testing.T) {
	fixedNonceBytes := make([]byte, 32)
	for i := range fixedNonceBytes {
		fixedNonceBytes[i] = byte(i)
	}
	var fixedTimestamp int64 = 1744918230067

	apiKey := "test-api-key"
	apiSecret := "test-api-secret"
	queryParamsRaw := "Alter=schwede&alter=schwede&param=value"
	bodyStr := `{"test":"data"}`

	signature, timestamp, nonce, err := generateRequestSignature(
		apiKey,
		apiSecret,
		queryParamsRaw,
		bodyStr,
		fixedTimestamp,
		fixedNonceBytes,
	)

	if err != nil {
		t.Fatalf("Signature generation failed: %v", err)
	}

	expectedNonce := base64.StdEncoding.EncodeToString(fixedNonceBytes)
	expectedTimestamp := strconv.FormatInt(fixedTimestamp, 10)

	queryParams := strings.ReplaceAll(queryParamsRaw, "&", "")
	queryParams = strings.ReplaceAll(queryParams, "=", "")

	digestInput := expectedNonce + expectedTimestamp + apiKey + queryParams + bodyStr
	digest := security.Sha256Hex(digestInput)
	signInput := digest + apiSecret
	expectedSignature := security.Sha256Hex(signInput)

	if signature != expectedSignature {
		t.Errorf("Expected signature %s, got %s", expectedSignature, signature)
	}

	if timestamp != expectedTimestamp {
		t.Errorf("Expected timestamp %s, got %s", expectedTimestamp, timestamp)
	}

	if nonce != expectedNonce {
		t.Errorf("Expected nonce %s, got %s", expectedNonce, nonce)
	}
}

func TestRequestSigner(t *testing.T) {
	fixedNonceBytes := make([]byte, 32)
	for i := range fixedNonceBytes {
		fixedNonceBytes[i] = byte(i)
	}
	fixedNonce := base64.StdEncoding.EncodeToString(fixedNonceBytes)
	var fixedTimestamp int64 = 1744918230067

	requestBody := []byte(`{"test":"data"}`)
	query := url.Values{}
	query.Set("Alter", "schwede")
	query.Set("alter", "schwede")
	query.Set("param", "value")

	uri, _ := url.Parse("https://openapidoc.bitunix.com/")
	uri.RawQuery = query.Encode()
	uri.Path = "/test"

	req, _ := http.NewRequestWithContext(context.Background(), "POST", uri.String(), bytes.NewReader(requestBody))

	requestSigner := RequestSigner("test-api-key", "test-api-secret", MockMillisecondTimestampGenerator(fixedTimestamp), MockNonceGenerator(fixedNonceBytes))
	err := requestSigner(req, requestBody)
	if err != nil {
		t.Errorf("Unable to sign request")
	}

	queryParams := "Alterschwedealterschwedeparamvalue"
	bodyStr := `{"test":"data"}`
	digestInput := fixedNonce + strconv.FormatInt(fixedTimestamp, 10) + "test-api-key" + queryParams + bodyStr
	digest := security.Sha256Hex(digestInput)
	signInput := digest + "test-api-secret"
	expectedSignature := security.Sha256Hex(signInput)

	if req.Header.Get("sign") != expectedSignature {
		t.Errorf("Expected signature %s, got %s", expectedSignature, req.Header.Get("sign"))
	}

	if req.Header.Get("nonce") != fixedNonce {
		t.Errorf("Expected nonce %s, got %s", fixedNonce, req.Header.Get("nonce"))
	}

	if req.Header.Get("timestamp") != strconv.FormatInt(fixedTimestamp, 10) {
		t.Errorf("Expected timestamp %d, got %s", fixedTimestamp, req.Header.Get("timestamp"))
	}
}

func TestRequestSigner_NonceError(t *testing.T) {
	mockError := errors.New("nonce generation failed")

	requestBody := []byte(`{"test":"data"}`)
	query := url.Values{}
	query.Set("param", "value")

	uri, _ := url.Parse("https://openapidoc.bitunix.com/")
	uri.RawQuery = query.Encode()
	uri.Path = "/test"

	req, _ := http.NewRequestWithContext(context.Background(), "POST", uri.String(), bytes.NewReader(requestBody))

	requestSigner := RequestSigner(
		"test-api-key",
		"test-api-secret",
		MockMillisecondTimestampGenerator(1744918230067),
		MockFailingNonceGenerator(mockError),
	)

	err := requestSigner(req, requestBody)

	if err == nil {
		t.Error("Expected an error when nonce generation fails, but got nil")
	}
}

func TestHandleAPIResponse_Success(t *testing.T) {
	responseBody := []byte(`{
		"code": 0,
		"message": "success",
		"data": {
			"field1": "value1",
			"field2": 123
		}
	}`)

	type TestResult struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Field1 string `json:"field1"`
			Field2 int    `json:"field2"`
		} `json:"data"`
	}

	var result TestResult
	err := handleAPIResponse(responseBody, "/test/endpoint", &result)

	if err != nil {
		t.Errorf("Expected no error for successful response, got: %v", err)
	}

	if result.Code != 0 {
		t.Errorf("Expected code 0, got %d", result.Code)
	}

	if result.Message != "success" {
		t.Errorf("Expected message 'success', got %s", result.Message)
	}

	if result.Data.Field1 != "value1" {
		t.Errorf("Expected data.field1 'value1', got %s", result.Data.Field1)
	}

	if result.Data.Field2 != 123 {
		t.Errorf("Expected data.field2 123, got %d", result.Data.Field2)
	}
}

func TestHandleAPIResponse_Error(t *testing.T) {

	testCases := []struct {
		name           string
		responseBody   string
		expectedErrMsg string
		endpoint       string
	}{
		{
			name: "Authentication Error",
			responseBody: `{
				"code": 10003,
				"message": "Authentication failed"
			}`,
			expectedErrMsg: "API error on /auth: code=10003, message=Authentication failed",
			endpoint:       "/auth",
		},
		{
			name: "Rate Limit Exceeded",
			responseBody: `{
				"code": 10005,
				"message": "Rate limit exceeded"
			}`,
			expectedErrMsg: "API error on /rate: code=10005, message=Rate limit exceeded",
			endpoint:       "/rate",
		},
		{
			name: "Insufficient Balance",
			responseBody: `{
				"code": 20003,
				"message": "Insufficient balance"
			}`,
			expectedErrMsg: "API error on /balance: code=20003, message=Insufficient balance",
			endpoint:       "/balance",
		},
		{
			name: "Order Not Found",
			responseBody: `{
				"code": 20007,
				"message": "Order not found"
			}`,
			expectedErrMsg: "API error on /orders: code=20007, message=Order not found",
			endpoint:       "/orders",
		},
		{
			name: "Alternative Message Field",
			responseBody: `{
				"code": 10001,
				"msg": "Network error occurred"
			}`,
			expectedErrMsg: "API error on /network: code=10001, message=Network error occurred",
			endpoint:       "/network",
		},
	}

	type TestResult struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Msg     string `json:"msg"`
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var result TestResult
			err := handleAPIResponse([]byte(tc.responseBody), tc.endpoint, &result)

			if err == nil {
				t.Errorf("Expected error for response with code > 0, got nil")
				return
			}

			if err.Error() != tc.expectedErrMsg {
				t.Errorf("Expected error message '%s', got '%s'", tc.expectedErrMsg, err.Error())
			}
		})
	}
}

func TestHandleAPIResponse_UnmarshalError(t *testing.T) {

	responseBody := []byte(`{ invalid json }`)

	type TestResult struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	var result TestResult
	err := handleAPIResponse(responseBody, "/test", &result)

	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestNewApiClient_WithOptions(t *testing.T) {
	apiKey := "test-key"
	apiSecret := "test-secret"
	customURI := "https://custom-api.example.com/"

	client, err := NewApiClient(apiKey, apiSecret, WithBaseURI(customURI))

	if err != nil {
		t.Fatalf("Failed to create API client: %v", err)
	}

	apiClient, ok := client.(*apiClient)
	if !ok {
		t.Fatal("Expected *apiClient type")
	}

	if apiClient.baseURI != customURI {
		t.Errorf("Expected baseURI to be %s, got %s", customURI, apiClient.baseURI)
	}
}
