package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateNonce(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func Sha256(message string) []byte {
	hash := sha256.New()
	hash.Write([]byte(message))
	return hash.Sum(nil)
}

func Sha256Hex(message string) string {
	return hex.EncodeToString(Sha256(message))
}
