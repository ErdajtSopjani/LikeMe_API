package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

// GenerateToken generates a secure random token
func GenerateToken() string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal("failed to generate token:", err)
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}
