package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/account"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestCase is a struct used to store expected and returned values when running tests
type TestCase struct {
	Name         string
	ReqBody      interface{}
	ExpectedCode int
	ExpectedBody string
}

// SetupTestDB sets up a connection to the test database
func SetupTestDB() *gorm.DB {
	// get test-db variables from env
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_TEST_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	// connect to the test database
	db, err := gorm.Open(
		postgres.Open(
			fmt.Sprintf(
				"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
				dbHost, dbPort, dbUser, dbName, dbPassword, dbSSLMode)),
		&gorm.Config{})

	// connect to the test database
	if err != nil {
		log.Fatalf("Failed to connect to the test database: %v", err)
	}

	// create tables in the test database
	err = runMigrations(db)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// return test DB connection
	return db
}

func runMigrations(db *gorm.DB) error {
	// use the already defined table models on the test database
	return db.AutoMigrate(
		&handlers.User{},
		&handlers.UserProfile{},
		&handlers.UserToken{},
		&handlers.Tag{},
		&handlers.UserInterest{},
		&handlers.Follow{},
		&handlers.BlockedUser{},
		&handlers.Post{},
		&handlers.Like{},
		&handlers.Comment{},
		&handlers.Message{},
		&handlers.TwoFactor{},
		&handlers.VerificationToken{},
	)
}

// CleanupTestDB clears the database after each test if needed
func CleanupTestDB(db *gorm.DB) {
	tables := []string{
		"users",
		"user_profiles",
		"user_tokens",
		"tags",
		"user_interests",
		"follows",
		"blocked_users",
		"posts",
		"likes",
		"comments",
		"messages",
		"two_factors",
		"verification_tokens",
	}

	// Truncate all tables and reset the primary key sequences
	for _, table := range tables {
		err := db.Exec("TRUNCATE TABLE " + table + " RESTART IDENTITY CASCADE").Error
		if err != nil {
			log.Fatalf("Failed to truncate table %s: %v", table, err)
		}
	}
}

// RunTests interates through the given testCase/s and runs them
func RunTests(db *gorm.DB, t *testing.T, testCases []TestCase) {
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			// create request body
			reqBody, err := json.Marshal(tt.ReqBody)
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
			assert.Equal(t, tt.ExpectedCode, rr.Code, "Status code should match")

			// check response body as JSON
			var response map[string]string
			if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
				t.Fatalf("Response body is not valid JSON: %v", err)
			}

			// check if the "message" field matches the expected body
			assert.Equal(t, tt.ExpectedBody, response["message"], "Response message should match")
		})
	}
}
