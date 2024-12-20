package users_test

import (
	"net/http"
	"testing"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers/social/follows"
	"github.com/ErdajtSopjani/LikeMe_API/tests"
)

func TestFollow(t *testing.T) {
	// connect to test database
	db := tests.SetupTestDB(t)
	tests.CleanupTestDB(db)

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
			RequestType:  "POST",
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
			RequestType:  "POST",
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
			RequestType:  "POST",
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
			RequestType:  "POST",
		},
	}

	tests.RunTests(db, t, testCases, "/api/v1/follow", follows.FollowAccount(db))
}
