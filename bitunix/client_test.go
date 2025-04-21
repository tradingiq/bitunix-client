package bitunix

import (
	"bitunix-client/security"
	"bytes"
	"context"
	"encoding/base64"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

// MockNonceGenerator returns a predetermined nonce for testing
func MockNonceGenerator(bytes []byte) func(int) ([]byte, error) {
	return func(length int) ([]byte, error) {
		return bytes[:length], nil
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

	requestSigner := createRequestSigner("test-api-key", "test-api-secret", MockMillisecondTimestampGenerator(fixedTimestamp), MockNonceGenerator(fixedNonceBytes))
	err := requestSigner(req, requestBody)
	if err != nil {

		t.Errorf("Unable to sign request")
	}

	// ascii ordered!
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
