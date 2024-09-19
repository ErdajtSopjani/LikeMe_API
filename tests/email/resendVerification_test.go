package email_test

import (
	"net/http"
	"testing"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/email"
	"github.com/ErdajtSopjani/LikeMe_API/tests"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestResendVerificationEmail(t *testing.T) {
	// connect to test database
	db = tests.SetupTestDB(t)

	// setup the db with the required entries to run tests
	err := db.Exec(tests.ReadSQLFile("resendVerificationTests.sql")).Error
	if err != nil {
		t.Fatalf("Failed to run resendVerificationTests.sql: %v", err)
	}

	testCases := []tests.TestCase{
		{
			Name: "Email not found",
			ReqBody: map[string]string{
				"email": "invalidmail@mailmail.com",
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Email not found",
			QueryParams:  "",
		},
		{
			Name: "Already Verified",
			ReqBody: map[string]string{
				"email": "verified-email@gmail.com",
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "User already verified",
			QueryParams:  "",
		},
		{
			Name: "Successful Resend",
			ReqBody: map[string]string{
				"email": "erdajtsopjani.tech@gmail.com",
			},
			ExpectedCode: http.StatusOK,
			ExpectedBody: "Email sent",
			QueryParams:  "",
		},
	}

	tests.RunTests(db, t, testCases, "/email/resend/register", email.ResendVerificationEmail(db))
}
