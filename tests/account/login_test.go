package account_test

import (
	"net/http"
	"testing"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/account"
	"github.com/ErdajtSopjani/LikeMe_API/tests"
)

// Test for RegisterUser handler
func TestLogin(t *testing.T) {
	db := tests.SetupTestDB(t) // connect to db

	// setup the db with the required entries to run login tests
	tests.SetupDBEntries("loginTests.sql", db, t)

	testCases := []tests.TestCase{
		{
			Name:         "Empty Code",
			ReqHeaders:   map[string]string{},
			ReqBody:      map[string]string{},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Empty Code",
			QueryParams:  "",
		},
		{
			Name:         "Invalid Code",
			ReqHeaders:   map[string]string{},
			ReqBody:      map[string]string{},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Invalid Code",
			QueryParams:  "code=1234",
		},
		{
			Name:         "Code Expired",
			ReqHeaders:   map[string]string{},
			ReqBody:      map[string]string{},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Code Expired",
			QueryParams:  "code=162508",
		},
		{
			Name:         "No user with ID found",
			ReqHeaders:   map[string]string{},
			ReqBody:      map[string]string{},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Invalid User Record",
			QueryParams:  "code=112308",
		},
		{
			Name:         "Successful login",
			ReqHeaders:   map[string]string{},
			ReqBody:      map[string]string{},
			ExpectedCode: http.StatusOK,
			ExpectedBody: "",
			QueryParams:  "code=5692124",
		},
	}

	tests.RunTests(db, t, testCases, "/login", account.Login(db))
}
