package account_test

import (
	"net/http"
	"testing"

	"github.com/ErdajtSopjani/LikeMe_API/internal/config"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/account"
	"github.com/ErdajtSopjani/LikeMe_API/tests"
)

var app config.AppConfig

// Test for RegisterUser handler
func TestRegisterUser(t *testing.T) {
	db := tests.SetupTestDB(t) // connect to db
	tests.CleanupTestDB(db)

	// setup db with the required entries to run tests
	tests.SetupDBEntries("registerTests.sql", db, t)

	testCases := []tests.TestCase{
		{
			Name:       "Create valid user",
			ReqHeaders: map[string]string{},
			ReqBody: map[string]string{
				"email":        "validuser@gmail.com",
				"country_code": "RKS",
			},
			ExpectedCode: http.StatusCreated,
			ExpectedBody: "User created",
			QueryParams:  "",
			RequestType:  "POST",
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
			QueryParams:  "",
			RequestType:  "POST",
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
			QueryParams:  "",
			RequestType:  "POST",
		},
		{
			Name:       "Missing country code",
			ReqHeaders: map[string]string{},
			ReqBody: map[string]string{
				"email":        "verify.likeme.dev@gmail.com",
				"country_code": "",
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Country Code is required",
			QueryParams:  "",
			RequestType:  "POST",
		},
	}

	app.IsProd = false
	app.IsTest = true
	tests.RunTests(db, t, testCases, "/register", account.RegisterUser(db, app))
}
