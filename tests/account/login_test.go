package account_test

import (
	"net/http"
	"testing"

	"gorm.io/gorm"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/account"
	"github.com/ErdajtSopjani/LikeMe_API/tests"
)

var db *gorm.DB

// Test for RegisterUser handler
func TestLogin(t *testing.T) {
	// connect to test database
	db = tests.SetupTestDB(t)

	// setup the db with the required entries to run login tests
	err := db.Exec(tests.ReadSQLFile("loginTests.sql")).Error
	if err != nil {
		t.Fatalf("Failed to run loginTests.sql: %v", err)
	}

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
}
