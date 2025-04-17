package security

import (
	"encoding/base64"
	"encoding/hex"
	"testing"
)

func TestSha256Hex(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:     "hello world",
			input:    "hello world",
			expected: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		},
		{
			name:     "complex string",
			input:    "nonce123456timestamp123apiKey123params123body123",
			expected: "da89a660059394f78b8729241ceec3955950be0d28140f0e7d5a8ec17da30b2b",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Sha256Hex(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, result)
			}
		})
	}
}

func TestSha256(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		},
		{
			name:     "hello world",
			input:    "hello world",
			expected: "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := Sha256(tc.input)
			hexResult := hex.EncodeToString(result)
			if hexResult != tc.expected {
				t.Errorf("Expected %s, got %s", tc.expected, hexResult)
			}
		})
	}
}

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
