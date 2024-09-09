package account_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/account"
	"github.com/ErdajtSopjani/LikeMe_API/tests"
)

var db *gorm.DB

// Test for RegisterUser handler
func TestRegisterUser(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	t.Setenv("DB_HOST", os.Getenv("DB_HOST"))
	t.Setenv("DB_PORT", os.Getenv("DB_PORT"))
	t.Setenv("DB_USER", os.Getenv("DB_USER"))
	t.Setenv("DB_PASSWORD", os.Getenv("DB_PASSWORD"))
	t.Setenv("DB_TEST_NAME", os.Getenv("DB_TEST_NAME"))
	t.Setenv("DB_SSLMODE", os.Getenv("DB_SSLMODE"))

	// connect to test database
	db = tests.SetupTestDB()

	testCases := []struct {
		name         string
		reqBody      map[string]string
		expectedCode int
		expectedBody string
	}{
		{
			name: "Create valid user",
			reqBody: map[string]string{
				"email":        "test@example.com",
				"country_code": "OK",
			},
			expectedCode: http.StatusCreated,
			expectedBody: "User created",
		},
		{
			name: "Duplicate email",
			reqBody: map[string]string{
				"email":        "test@example.com", // already exists after the first test
				"country_code": "BAD",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Email already taken",
		},
		{
			name: "Invalid email format",
			reqBody: map[string]string{
				"email":        "invalidemail@gmail",
				"country_code": "BAD",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid Email",
		},
		{
			name: "Missing country code",
			reqBody: map[string]string{
				"email":        "test@example.com",
				"country_code": "",
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: "Country Code is required",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// create request body
			reqBody, err := json.Marshal(tt.reqBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			// create request
			req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// create response recorder
			rr := httptest.NewRecorder()

			// create handler
			handler := http.HandlerFunc(account.RegisterUser(db))

			// serve http
			handler.ServeHTTP(rr, req)

			// check status code
			assert.Equal(t, tt.expectedCode, rr.Code)

			// check response body as JSON
			var response map[string]string
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Fatalf("Response body is not valid JSON: %v", err)
			}

			// check if the "message" field matches the expected body
			assert.Equal(t, tt.expectedBody, response["message"])
		})
	}

	tests.CleanupTestDB(db) // cleanup database
}
