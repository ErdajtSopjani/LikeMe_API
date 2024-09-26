package email_test

import (
	"net/http"
	"testing"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/email/verify"
	"github.com/ErdajtSopjani/LikeMe_API/tests"
)

func TestResendVerificationEmail(t *testing.T) {
	// connect to test database
	db := tests.SetupTestDB(t)

	// setup the db with the required entries to run tests
	tests.SetupDBEntries("resendVerificationTests.sql", db, t)

	testCases := []tests.TestCase{
		{
			Name:       "Email not found",
			ReqHeaders: map[string]string{},
			ReqBody: map[string]string{
				"email": "invalidmail@mailmail.com",
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Email not found",
			QueryParams:  "",
			RequestType:  "POST",
		},
		{
			Name:       "Already Verified",
			ReqHeaders: map[string]string{},
			ReqBody: map[string]string{
				"email": "verified-email@gmail.com",
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "User already verified",
			QueryParams:  "",
			RequestType:  "POST",
		},
		{
			Name:       "Successful Resend",
			ReqHeaders: map[string]string{},
			ReqBody: map[string]string{
				"email": "erdajtsopjani.tech@gmail.com",
			},
			ExpectedCode: http.StatusOK,
			ExpectedBody: "Email sent",
			QueryParams:  "",
			RequestType:  "POST",
		},
	}

	tests.RunTests(db, t, testCases, "/email/resend/register", verify.ResendVerificationEmail(db))
}
