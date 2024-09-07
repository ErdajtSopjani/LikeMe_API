package handlers

import (
	"fmt"
	"log"

	"crypto/rand"
	"encoding/base64"
	"math/big"
)

// GenerateToken generates a secure random token
func GenerateToken() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Panicf("\n\nERROR\n\tFailed to generate token\n\tERROR: %v\n\n", err)
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
}

// GenerateCode generates a 6-digit code
func GenerateCode() (string, error) {
	// Define the range for a 6-digit code (100000 to 999999)
	max := big.NewInt(900000)

	// generate random number between 0-899999, then add 100000 to get 6 digit num
	num, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	code := num.Int64() + 100000
	return fmt.Sprintf("%06d", code), nil
}
