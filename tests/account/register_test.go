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

	// define request payload for test
	reqBody := map[string]string{
		"email":        "test@example.com",
		"country_code": "+1",
	}
	body, _ := json.Marshal(reqBody)

	// make new POST request
	req, err := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// use httptest to create a ResponseRecorder
	rr := httptest.NewRecorder()

	// call the RegisterUser handler
	handler := account.RegisterUser(db)
	handler.ServeHTTP(rr, req)

	// check if the status code is 201 Created
	assert.Equal(t, http.StatusCreated, rr.Code)

	// check the response body
	expected := "User created"
	assert.Equal(t, expected, rr.Body.String())

	// verify that user has been added to the database
	var count int64
	err = db.Table("users").Where("email = ?", "test@example.com").Count(&count).Error
	if err != nil {
		t.Fatalf("Failed to check the database for user: %v", err)
	}

	// ensure user was added
	assert.Equal(t, int64(1), count)

	tests.CleanupTestDB(db) // cleanup database
}
