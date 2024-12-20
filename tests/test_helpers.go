package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TestCase is a struct used to store expected and returned values when running tests
type TestCase struct {
	Name         string
	ReqHeaders   interface{}
	ReqBody      interface{}
	ExpectedCode int
	ExpectedBody string
	QueryParams  string
	RequestType  string
}

// InitTestDB makes a connection to the test database
func InitTestDB() *gorm.DB {
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

// SetupTestDB sets up the database env vars with the current test
func SetupTestDB(t *testing.T) *gorm.DB {
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

	db := InitTestDB()
	return db
}

// SetupDBEntries reads a SQL file required for tests and executes it
func SetupDBEntries(filePath string, db *gorm.DB, t *testing.T) {
	CleanupTestDB(db)

	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}

	err = db.Exec(string(content)).Error
	if err != nil {
		t.Fatalf("Failed to run SQL file: %v", err)
	}
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
		err := db.Exec("DELETE FROM " + table).Error
		if err != nil {
			log.Printf("Failed to delete records from table %s: %v", table, err)
		}
		err = db.Exec("ALTER SEQUENCE " + table + "_id_seq RESTART WITH 1").Error
		if err != nil {
			log.Printf("Failed to reset sequence for table %s: %v", table, err)
		}
	}

	// print each table with the number of rows to confirm that the tables are empty
	for _, table := range tables {
		var count int64
		db.Table(table).Count(&count)
		fmt.Printf("Table %s has %d rows\n", table, count)
	}
}

// RunTests iterates through the given test cases and runs them using the provided handler
func RunTests(db *gorm.DB, t *testing.T, testCases []TestCase, baseURL string, handler http.HandlerFunc) {
	for _, tt := range testCases {
		t.Run(tt.Name, func(t *testing.T) {
			// create request body
			reqBody, err := json.Marshal(tt.ReqBody)
			if err != nil {
				t.Fatalf("Failed to marshal request body: %v", err)
			}

			// build the full URL (including any query parameters)
			url := baseURL
			if tt.QueryParams != "" {
				url += "?" + tt.QueryParams
			}

			// create request with body and headers
			req, err := http.NewRequest(tt.RequestType, url, bytes.NewBuffer(reqBody))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			for key, value := range tt.ReqHeaders.(map[string]string) {
				req.Header.Set(key, value)
			}

			// create response recorder
			rr := httptest.NewRecorder()

			// serve http using the provided handler
			handler.ServeHTTP(rr, req)

			// check if status codes match
			assert.Equal(t, tt.ExpectedCode, rr.Code, "Status code should match")

			// expectedBody can be left empty on unpredicatable responses
			if tt.ExpectedBody == "" {
				return
			}

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
