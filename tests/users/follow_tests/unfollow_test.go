package users_test

import (
	"net/http"
	"testing"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/social/follows"
	"github.com/ErdajtSopjani/LikeMe_API/tests"
)

func TestUnfollow(t *testing.T) {
	db := tests.SetupTestDB(t) // connect to test database

	tests.SetupDBEntries("followTests.sql", db, t) // setup the db with the required entries to run tests

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
			RequestType:  "POST",
		},
		{
			Name: "Token and User don't match",
			ReqHeaders: map[string]string{
				"Authorization": "token1",
				"Content-Type":  "application/json",
			},
			ReqBody: map[string]interface{}{
				"follower_id":  2,
				"following_id": 3,
			},
			ExpectedCode: http.StatusUnauthorized,
			ExpectedBody: "Unauthorized",
			QueryParams:  "",
			RequestType:  "POST",
		}, {
			Name: "Follow doesn't exist",
			ReqHeaders: map[string]string{
				"Authorization": "token1",
				"Content-Type":  "application/json",
			},
			ReqBody: map[string]int{
				"follower_id":  1,
				"following_id": 3,
			},
			ExpectedCode: http.StatusBadRequest,
			ExpectedBody: "Follow does not exist",
			QueryParams:  "",
			RequestType:  "POST",
		}, {
			Name: "Successful Unfollow",
			ReqHeaders: map[string]string{
				"Authorization": "token2",
				"Content-Type":  "application/json",
			},
			ReqBody: map[string]int{
				"follower_id":  2,
				"following_id": 3,
			},
			ExpectedCode: http.StatusOK,
			ExpectedBody: "Unfollowed successfully",
			QueryParams:  "",
			RequestType:  "POST",
		},
	}

	tests.RunTests(db, t, testCases, "/api/v1/follow", follows.UnfollowAccount(db))
}
