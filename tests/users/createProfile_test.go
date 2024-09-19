package users_test

import (
	"net/http"
	"testing"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/users"
	"github.com/ErdajtSopjani/LikeMe_API/tests"
)

func TestCreateProfile(t *testing.T) {
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
		},
		{
			Name: "Username taken",
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
			ExpectedBody: "Username already taken",
			QueryParams:  "",
		},
	}

	tests.RunTests(db, t, testCases, "/api/v1/profile", users.CreateProfile(db))
}
