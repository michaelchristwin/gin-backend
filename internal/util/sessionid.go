package util

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateSessionID() (string, error) {
	b := make([]byte, 32) // 256 bits of entropy
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
