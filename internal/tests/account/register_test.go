package account_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/account"
	"github.com/ErdajtSopjani/LikeMe_API/internal/tests"
)

var db *gorm.DB

// Set up before running tests
func setup() {
	db = tests.SetupTestDB() // Connect to the test database
}

// Clean up after tests
func teardown() {
	tests.CleanupTestDB(db) // Clean up the test DB
}

// Test for RegisterUser handler
func TestRegisterUser(t *testing.T) {
	// Run the setup before each test
	setup()

	// Define a request payload for the test
	reqBody := map[string]string{
		"email":        "test@example.com",
		"country_code": "+1",
	}
	body, _ := json.Marshal(reqBody)

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	// Use httptest to create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Call the RegisterUser handler
	handler := account.RegisterUser(db)
	handler.ServeHTTP(rr, req)

	// Check if the status code is 201 Created
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Check the response body
	expected := "User created"
	assert.Equal(t, expected, rr.Body.String())

	// Verify that the user has been added to the database
	var count int64
	err = db.Table("users").Where("email = ?", "test@example.com").Count(&count).Error
	if err != nil {
		t.Fatalf("Failed to check the database for user: %v", err)
	}

	// Ensure a user was added
	assert.Equal(t, int64(1), count)

	// Clean up after the test
	teardown()
}
