package security

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256(message string) []byte {
	hash := sha256.New()
	hash.Write([]byte(message))
	return hash.Sum(nil)
}

func Sha256Hex(message string) string {
	return hex.EncodeToString(Sha256(message))
}
