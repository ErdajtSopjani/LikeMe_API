package account_test

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/gorm"

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

	testCases := []tests.TestCase{
		{
			Name: "Create valid user",
			ReqBody: map[string]string{
				"email":        "test@example.com",
				"country_code": "OK",
			},
			ExpectedCode: http.StatusCreated,
			ExpectedBody: "User created",
		},
		{
			Name: "Duplicate email",
			ReqBody: map[string]string{
				"email":        "test@example.com", // already exists after the first test
				"country_code": "BAD",
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Email already taken",
		},
		{
			Name: "Invalid email format",
			ReqBody: map[string]string{
				"email":        "invalidemail@gmail",
				"country_code": "BAD",
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Invalid Email",
		},
		{
			Name: "Missing country code",
			ReqBody: map[string]string{
				"email":        "test@example.com",
				"country_code": "",
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Country Code is required",
		},
	}

	tests.RunTests(db, t, testCases)

	tests.CleanupTestDB(db) // cleanup database
}
