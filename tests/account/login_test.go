package account_test

import (
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"gorm.io/gorm"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/account"
	"github.com/ErdajtSopjani/LikeMe_API/tests"
)

var db *gorm.DB

// Test for RegisterUser handler
func TestLogin(t *testing.T) {
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

	// run the raw sql query located in this directory called loginTests.sql via gorm
	err = db.Exec(tests.ReadSQLFile("loginTests.sql")).Error

	testCases := []tests.TestCase{
		{
			Name:         "Empty Code",
			ReqBody:      map[string]string{},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Empty Code",
			QueryParams:  "",
		},
		{
			Name:         "Invalid Code",
			ReqBody:      map[string]string{},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Invalid Code",
			QueryParams:  "code=1234",
		},
		{
			Name:         "Code Expired",
			ReqBody:      map[string]string{},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Code Expired",
			QueryParams:  "code=162508",
		},
		{
			Name:         "No user with ID found",
			ReqBody:      map[string]string{},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Invalid User Record",
			QueryParams:  "code=112308",
		},
		{
			Name:         "Successful login",
			ReqBody:      map[string]string{},
			ExpectedCode: http.StatusOK,
			ExpectedBody: "",
			QueryParams:  "code=5692124",
		},
	}

	tests.RunTests(db, t, testCases, "/login", account.Login(db))

	// connect to test database
	db = tests.SetupTestDB()

	tests.CleanupTestDB(db) // cleanup database
}
