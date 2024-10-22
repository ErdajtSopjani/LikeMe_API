package profile_tests

import (
	"net/http"
	"testing"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/social/profiles"
	"github.com/ErdajtSopjani/LikeMe_API/tests"
)

func TestUpdateProfile(t *testing.T) {
	// connect to test database
	db := tests.SetupTestDB(t)
	tests.CleanupTestDB(db)

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
			RequestType:  "PUT",
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
			ExpectedBody: "Profile Created",
			QueryParams:  "",
			RequestType:  "PUT",
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
			ExpectedCode: http.StatusOK,
			ExpectedBody: "Profile Updated",
			QueryParams:  "",
			RequestType:  "PUT",
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
			RequestType:  "PUT",
		},
	}

	tests.RunTests(db, t, testCases, "/api/v1/profile", profiles.ManageProfiles(db))
}
