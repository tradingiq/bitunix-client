package bitunix

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"github.com/tradingiq/bitunix-client/security"
	"testing"
)

func TestGenerateWebsocketSignature(t *testing.T) {
	apiKey := "test_api_key"
	apiSecret := "test_api_secret"
	timestamp := int64(1234567890)
	nonceBytes := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	expectedNonceHex := "0102030405060708"
	expectedPreSign := expectedNonceHex + "1234567890" + apiKey
	expectedPreSignHash := security.Sha256Hex(expectedPreSign)
	expectedSign := security.Sha256Hex(expectedPreSignHash + apiSecret)

	actualSign, actualTimestamp := GenerateWebsocketSignature(apiKey, apiSecret, timestamp, nonceBytes)

	assert.Equal(t, expectedSign, actualSign, "The signature should match the expected value")
	assert.Equal(t, timestamp, actualTimestamp, "The timestamp should be returned unchanged")
}

func TestGenerateWebsocketSignatureFormat(t *testing.T) {
	apiKey := "test_api_key"
	apiSecret := "test_api_secret"
	timestamp := int64(1234567890)
	nonceBytes := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	sign, _ := GenerateWebsocketSignature(apiKey, apiSecret, timestamp, nonceBytes)

	assert.Equal(t, 64, len(sign), "Signature should be 64 characters long (32 bytes in hex)")

	_, err := hex.DecodeString(sign)
	assert.NoError(t, err, "Signature should be a valid hex string")
}

func TestGenerateWebsocketSignatureDifferentInputs(t *testing.T) {
	tests := []struct {
		name       string
		apiKey     string
		apiSecret  string
		timestamp  int64
		nonceBytes []byte
	}{
		{
			name:       "Empty inputs",
			apiKey:     "",
			apiSecret:  "",
			timestamp:  0,
			nonceBytes: []byte{},
		},
		{
			name:       "Different API key",
			apiKey:     "different_api_key",
			apiSecret:  "test_api_secret",
			timestamp:  1234567890,
			nonceBytes: []byte{0x01, 0x02, 0x03, 0x04},
		},
		{
			name:       "Different API secret",
			apiKey:     "test_api_key",
			apiSecret:  "different_api_secret",
			timestamp:  1234567890,
			nonceBytes: []byte{0x01, 0x02, 0x03, 0x04},
		},
		{
			name:       "Different timestamp",
			apiKey:     "test_api_key",
			apiSecret:  "test_api_secret",
			timestamp:  9876543210,
			nonceBytes: []byte{0x01, 0x02, 0x03, 0x04},
		},
		{
			name:       "Different nonce",
			apiKey:     "test_api_key",
			apiSecret:  "test_api_secret",
			timestamp:  1234567890,
			nonceBytes: []byte{0x08, 0x07, 0x06, 0x05},
		},
	}

	signatures := make(map[string]bool)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sign, returnedTimestamp := GenerateWebsocketSignature(tc.apiKey, tc.apiSecret, tc.timestamp, tc.nonceBytes)

			assert.Equal(t, tc.timestamp, returnedTimestamp)

			assert.Equal(t, 64, len(sign))
			_, err := hex.DecodeString(sign)
			assert.NoError(t, err)

			signatures[sign] = true
		})
	}

	assert.GreaterOrEqual(t, len(signatures), len(tests)-1,
		"Different inputs should produce different signatures in most cases")
}
