package account_test

import (
	"net/http"
	"testing"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/account"
	"github.com/ErdajtSopjani/LikeMe_API/tests"
)

// Test for RegisterUser handler
func TestRegisterUser(t *testing.T) {
	db := tests.SetupTestDB(t) // connect to db

	// setup db with the required entries to run tests
	tests.SetupDBEntries("registerTests.sql", db, t)

	testCases := []tests.TestCase{
		{
			Name:       "Create valid user",
			ReqHeaders: map[string]string{},
			ReqBody: map[string]string{
				"email":        "erdajttest@gmail.com",
				"country_code": "RKS",
			},
			ExpectedCode: http.StatusCreated,
			ExpectedBody: "User created",
			QueryParams:  "",
		},
		{
			Name:       "Duplicate email",
			ReqHeaders: map[string]string{},
			ReqBody: map[string]string{
				"email":        "erdajtsopjani.tech@gmail.com",
				"country_code": "BAD",
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Email already taken",
		},
		{
			Name:       "Invalid email format",
			ReqHeaders: map[string]string{},
			ReqBody: map[string]string{
				"email":        "invalidemail@gmail",
				"country_code": "BAD",
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Invalid Email",
		},
		{
			Name:       "Missing country code",
			ReqHeaders: map[string]string{},
			ReqBody: map[string]string{
				"email":        "test@example.com",
				"country_code": "",
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Country Code is required",
		},
	}

	tests.RunTests(db, t, testCases, "/register", account.RegisterUser(db))
}
