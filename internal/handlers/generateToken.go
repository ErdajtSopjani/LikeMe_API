package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

// GenerateToken generates a secure random token
func GenerateToken() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatalf("\n\nERROR\n\tFailed to generate token\n\tERROR: %v\n\n", err)
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
