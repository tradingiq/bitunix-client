package bitunix

import (
	"bitunix-client/api"
	"bitunix-client/security"
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
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

func TestDeterministicSignature(t *testing.T) {

	var capturedRequest *http.Request
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedRequest = r.Clone(context.Background()) // Save the request for verification
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	fixedNonceBytes := make([]byte, 32)
	for i := range fixedNonceBytes {
		fixedNonceBytes[i] = byte(i)
	}
	fixedNonce := base64.StdEncoding.EncodeToString(fixedNonceBytes)
	var fixedTimestamp int64 = 1744918230067

	apiClient, _ := api.New(server.URL)
	client := New(apiClient, "test-api-key", "test-api-secret")
	client.generateNonce = MockNonceGenerator(fixedNonceBytes)
	client.generateMillisecondTimestamp = MockMillisecondTimestampGenerator(fixedTimestamp)

	ctx := context.Background()
	requestBody := []byte(`{"test":"data"}`)
	query := url.Values{}
	query.Set("Alter", "schwede")
	query.Set("alter", "schwede")
	query.Set("param", "value")

	_, err := client.api.Post(ctx, "/test", query, requestBody)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}

	// ascii ordered!
	queryParams := "Alterschwedealterschwedeparamvalue"
	bodyStr := `{"test":"data"}`
	digestInput := fixedNonce + strconv.FormatInt(fixedTimestamp, 10) + "test-api-key" + queryParams + bodyStr
	digest := security.Sha256Hex(digestInput)
	signInput := digest + "test-api-secret"
	expectedSignature := security.Sha256Hex(signInput)

	if capturedRequest.Header.Get("sign") != expectedSignature {
		t.Errorf("Expected signature %s, got %s", expectedSignature, capturedRequest.Header.Get("sign"))
	}

	if capturedRequest.Header.Get("nonce") != fixedNonce {
		t.Errorf("Expected nonce %s, got %s", fixedNonce, capturedRequest.Header.Get("nonce"))
	}

	if capturedRequest.Header.Get("timestamp") != strconv.FormatInt(fixedTimestamp, 10) {
		t.Errorf("Expected timestamp %d, got %s", fixedTimestamp, capturedRequest.Header.Get("timestamp"))
	}
}
