package security

import (
	"crypto/rand"
	"fmt"
)

func GenerateNonce(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("error generating nonce: %v", err)
	}
	return b, nil
}
