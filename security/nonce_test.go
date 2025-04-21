package security

import (
	"encoding/base64"
	"testing"
)

func TestGenerateNonce(t *testing.T) {
	tests := []struct {
		name   string
		length int
	}{
		{
			name:   "32 bytes",
			length: 32,
		},
		{
			name:   "64 bytes",
			length: 64,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			nonce, err := GenerateNonce(tc.length)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if len(nonce) != tc.length {
				t.Errorf("Expected length %d, got %d", tc.length, len(nonce))
			}

			nonce2, err := GenerateNonce(tc.length)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			nonceBase64 := base64.StdEncoding.EncodeToString(nonce)
			nonce2Base64 := base64.StdEncoding.EncodeToString(nonce2)

			if nonceBase64 == nonce2Base64 {
				t.Error("Expected different nonces, but got identical values")
			}
		})
	}
}
