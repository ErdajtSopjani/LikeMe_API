package profile_tests

import (
	"net/http"
	"testing"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/users/profiles"
	"github.com/ErdajtSopjani/LikeMe_API/tests"
)

func TestCreateProfile(t *testing.T) {
	// connect to test database
	db := tests.SetupTestDB(t)

	// setup the db with the required entries to run login tests
	tests.SetupDBEntries("createProfileTests.sql", db, t)

	testCases := []tests.TestCase{
		{
			Name: "Empty fields",
			ReqHeaders: map[string]string{
				"Authorization": "jUy2Iti6p3GqQxp0TjwrGA==",
				"Content-Type":  "application/json",
			},
			ReqBody: map[string]string{
				"username":        "",
				"last_name":       "",
				"profile_picture": "",
				"bio":             "",
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "All fields are required",
			QueryParams:  "",
			RequestType:  "POST",
		},
		{
			Name: "Successful Profile Creation",
			ReqHeaders: map[string]string{
				"Authorization": "jUy2Iti6p3GqQxp0TjwrGA==",
				"Content-Type":  "application/json",
			},
			ReqBody: map[string]string{
				"username":        "erdajttsopjani",
				"first_name":      "Erdajt",
				"last_name":       "Sopjani",
				"profile_picture": "base64_image_here",
				"bio":             "null",
			},
			ExpectedCode: http.StatusCreated,
			ExpectedBody: "User profile created",
			QueryParams:  "",
			RequestType:  "POST",
		},
		{
			Name: "User already has a profile",
			ReqHeaders: map[string]string{
				"Authorization": "jUy2Iti6p3GqQxp0TjwrGA==",
				"Content-Type":  "application/json",
			},
			ReqBody: map[string]string{
				"username":        "erdajttsopjani",
				"first_name":      "Erdajt",
				"last_name":       "Sopjani",
				"profile_picture": "base64_image_here",
				"bio":             "null",
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "User already has a profile, you can update it.",
			QueryParams:  "",
			RequestType:  "POST",
		},
		{
			Name: "Username already taken",
			ReqHeaders: map[string]string{
				"Authorization": "jUy2Iti6p3GqQxp0Tj1234==",
				"Content-Type":  "application/json",
			},
			ReqBody: map[string]string{
				"username":        "erdajttsopjani",
				"first_name":      "Erdajt",
				"last_name":       "Sopjani",
				"profile_picture": "base64_image_here",
				"bio":             "null",
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Username already taken",
			QueryParams:  "",
			RequestType:  "POST",
		},
	}

	tests.RunTests(db, t, testCases, "/api/v1/profile", profiles.ManageProfiles(db))
}
