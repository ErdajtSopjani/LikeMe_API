package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func IsEmailValid(email string) (bool, error) {
	// make the email validation request
	resp, err := http.Get(fmt.Sprintf("https://emailvalidation.abstractapi.com/v1/?api_key=%s&email=%s", os.Getenv("ABSTRACT_API_KEY"), email))
	if err != nil {
		return false, err
	}

	// read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil || string(body) == "" {
		return false, err
	}

	// unmarshal response
	jsonResp := make(map[string]interface{})
	if err := json.Unmarshal(body, &jsonResp); err != nil {
		return false, err
	}

	println(string(body))
	// check if email is valid
	if jsonResp["deliverability"] == "DELIVERABLE" {
		return true, nil
	}

	return false, nil
}
