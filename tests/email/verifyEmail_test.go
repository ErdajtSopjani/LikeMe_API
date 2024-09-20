package email_test

import (
	"net/http"
	"testing"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/email"
	"github.com/ErdajtSopjani/LikeMe_API/tests"
)

func TestVerifyEmail(t *testing.T) {
	db := tests.SetupTestDB(t) // connect to db

	// setup the db with the required entries to run login tests
	tests.SetupDBEntries("verifyEmailTests.sql", db, t)

	testCases := []tests.TestCase{
		{
			Name:         "Empty Token",
			ReqHeaders:   map[string]string{},
			ReqBody:      map[string]string{},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Invalid Token",
			QueryParams:  "",
		},
		{
			Name:         "Invalid Token",
			ReqHeaders:   map[string]string{},
			ReqBody:      map[string]string{},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Invalid Token",
			QueryParams:  "token=1234",
		},
		{
			Name:         "Successful Verification",
			ReqHeaders:   map[string]string{},
			ReqBody:      map[string]string{},
			ExpectedCode: http.StatusOK,
			ExpectedBody: "Email Verified",
			QueryParams:  "token=123456",
		},
	}

	tests.RunTests(db, t, testCases, "/api/v1/email/verify", email.VerifyEmail(db))
}
