package bitunix

import (
	"encoding/hex"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tradingiq/bitunix-client/security"
	"testing"
	"time"
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

	actualSign, actualTimestamp := generateWebsocketSignature(apiKey, apiSecret, timestamp, nonceBytes)

	assert.Equal(t, expectedSign, actualSign, "The signature should match the expected value")
	assert.Equal(t, timestamp, actualTimestamp, "The timestamp should be returned unchanged")
}

func TestGenerateWebsocketSignatureFormat(t *testing.T) {
	apiKey := "test_api_key"
	apiSecret := "test_api_secret"
	timestamp := int64(1234567890)
	nonceBytes := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	sign, _ := generateWebsocketSignature(apiKey, apiSecret, timestamp, nonceBytes)

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
			name:       "Different apiClient key",
			apiKey:     "different_api_key",
			apiSecret:  "test_api_secret",
			timestamp:  1234567890,
			nonceBytes: []byte{0x01, 0x02, 0x03, 0x04},
		},
		{
			name:       "Different apiClient secret",
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
			sign, returnedTimestamp := generateWebsocketSignature(tc.apiKey, tc.apiSecret, tc.timestamp, tc.nonceBytes)

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

func TestKeepAliveMonitor(t *testing.T) {
	heartbeatGenerator := KeepAliveMonitor()
	require.NotNil(t, heartbeatGenerator)

	bytes, err := heartbeatGenerator()
	require.NoError(t, err)
	require.NotNil(t, bytes)

	var message heartbeatMessage
	err = json.Unmarshal(bytes, &message)
	require.NoError(t, err)

	assert.Equal(t, "ping", message.Op)
	assert.NotZero(t, message.Ping)
	assert.LessOrEqual(t, message.Ping, time.Now().Unix())
	assert.GreaterOrEqual(t, message.Ping, time.Now().Unix()-5)
}

func TestWebsocketSigner(t *testing.T) {
	apiKey := "test_api_key"
	apiSecret := "test_api_secret"

	signer := WebsocketSigner(apiKey, apiSecret)
	require.NotNil(t, signer)

	// Test the generator function
	bytes, err := signer()
	require.NoError(t, err)
	require.NotNil(t, bytes)

	// Parse the login message
	var message loginMessage
	err = json.Unmarshal(bytes, &message)
	require.NoError(t, err)

	// Verify the structure
	assert.Equal(t, "login", message.Op)
	assert.Len(t, message.Args, 1)

	loginParams := message.Args[0]
	assert.Equal(t, apiKey, loginParams.ApiKey, "The apiClient key in the message should match the provided key")
	assert.NotEmpty(t, loginParams.Nonce)
	assert.NotEmpty(t, loginParams.Sign)
	assert.NotZero(t, loginParams.Timestamp)

	// Verify the nonce is a valid hex string
	nonceBytes, err := hex.DecodeString(loginParams.Nonce)
	assert.NoError(t, err)
	assert.Equal(t, 32, len(nonceBytes))

	// Verify signature is valid hex of expected length
	_, err = hex.DecodeString(loginParams.Sign)
	assert.NoError(t, err)
	assert.Equal(t, 64, len(loginParams.Sign))

	// Test with different apiClient keys
	differentApiKey := "different_api_key"
	differentApiSecret := "different_api_secret"

	differentSigner := WebsocketSigner(differentApiKey, differentApiSecret)
	differentBytes, err := differentSigner()
	require.NoError(t, err)

	var differentMessage loginMessage
	err = json.Unmarshal(differentBytes, &differentMessage)
	require.NoError(t, err)

	// Check that the different apiClient key is actually used
	assert.Equal(t, differentApiKey, differentMessage.Args[0].ApiKey)

	// Create and verify signatures with known inputs
	knownTimestamp := int64(1234567890)
	knownNonce := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}

	sign1, _ := generateWebsocketSignature(apiKey, apiSecret, knownTimestamp, knownNonce)
	sign2, _ := generateWebsocketSignature(differentApiKey, differentApiSecret, knownTimestamp, knownNonce)

	// Different credentials should produce different signatures with the same inputs
	assert.NotEqual(t, sign1, sign2,
		"Different apiClient credentials should produce different signatures")
}

func TestLoginMessage(t *testing.T) {
	message := loginMessage{
		Op: "login",
		Args: []loginParams{
			{
				ApiKey:    "test_key",
				Timestamp: 1234567890,
				Nonce:     "0102030405060708",
				Sign:      "test_signature",
			},
		},
	}

	bytes, err := json.Marshal(message)
	require.NoError(t, err)

	var decoded loginMessage
	err = json.Unmarshal(bytes, &decoded)
	require.NoError(t, err)

	assert.Equal(t, "login", decoded.Op)
	assert.Len(t, decoded.Args, 1)
	assert.Equal(t, "test_key", decoded.Args[0].ApiKey)
	assert.Equal(t, int64(1234567890), decoded.Args[0].Timestamp)
	assert.Equal(t, "0102030405060708", decoded.Args[0].Nonce)
	assert.Equal(t, "test_signature", decoded.Args[0].Sign)
}

func TestHeartbeatMessage(t *testing.T) {
	timestamp := time.Now().Unix()
	message := heartbeatMessage{
		Op:   "ping",
		Ping: timestamp,
	}

	bytes, err := json.Marshal(message)
	require.NoError(t, err)

	var decoded heartbeatMessage
	err = json.Unmarshal(bytes, &decoded)
	require.NoError(t, err)

	assert.Equal(t, "ping", decoded.Op)
	assert.Equal(t, timestamp, decoded.Ping)
}
