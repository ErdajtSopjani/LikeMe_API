package users_test

import (
	"net/http"
	"testing"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/email"
	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/users"
	"github.com/ErdajtSopjani/LikeMe_API/tests"
)

func TestResendVerificationEmail(t *testing.T) {
	// connect to test database
	db := tests.SetupTestDB(t)

	// setup the db with the required entries to run tests
	tests.SetupDBEntries("followTests.sql", db, t)

	testCases := []tests.TestCase{
		{
			Name: "Invalid Users",
			ReqHeaders: map[string]string{
				"Authorization": "token1",
				"Content-Type":  "application/json",
			},
			ReqBody: map[string]string{
				"follower_id": "4",
				"followed_id": "5",
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Invalid following or follower id",
			QueryParams:  "",
		},
		{
			Name: "Token and User don't match",
			ReqHeaders: map[string]string{
				"Authorization": "token2",
				"Content-Type":  "application/json",
			},
			ReqBody: map[string]string{
				"follower_id": "1",
				"followed_id": "2",
			},
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: "Invalid user or token",
			QueryParams:  "",
		}, {
			Name: "Follow already exists",
			ReqHeaders: map[string]string{
				"Authorization": "token1",
				"Content-Type":  "application/json",
			},
			ReqBody: map[string]string{
				"follower_id": "2",
				"followed_id": "3",
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Follow already exists",
			QueryParams:  "",
		}, {
			Name: "Successful Follow",
			ReqHeaders: map[string]string{
				"Authorization": "token1",
				"Content-Type":  "application/json",
			},
			ReqBody: map[string]string{
				"follower_id": "1",
				"followed_id": "2",
			},
			ExpectedCode: http.StatusCreated,
			ExpectedBody: "",
			QueryParams:  "",
		},
	}

	tests.RunTests(db, t, testCases, "/api/v1/follow", users.FollowAccount(db))
}
