package users_test

import (
	"net/http"
	"testing"

	_users "github.com/ErdajtSopjani/LikeMe_API/internal/handlers/users"
	"github.com/ErdajtSopjani/LikeMe_API/tests"
)

func TestFollow(t *testing.T) {
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
			ReqBody: map[string]int{
				"follower_id":  4,
				"following_id": 5,
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Invalid follower or following ID",
			QueryParams:  "",
		},
		{
			Name: "Token and User don't match",
			ReqHeaders: map[string]string{
				"Authorization": "token1",
				"Content-Type":  "application/json",
			},
			ReqBody: map[string]interface{}{
				"follower_id":  3,
				"following_id": 2,
			},
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: "Unauthorized",
			QueryParams:  "",
		}, {
			Name: "Follow already exists",
			ReqHeaders: map[string]string{
				"Authorization": "token2",
				"Content-Type":  "application/json",
			},
			ReqBody: map[string]int{
				"follower_id":  2,
				"following_id": 3,
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
			ReqBody: map[string]int{
				"follower_id":  1,
				"following_id": 2,
			},
			ExpectedCode: http.StatusCreated,
			ExpectedBody: "Follow successfully created",
			QueryParams:  "",
		},
	}

	tests.RunTests(db, t, testCases, "/api/v1/follow", _users.FollowAccount(db))
}
