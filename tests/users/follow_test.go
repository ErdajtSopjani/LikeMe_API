package users_test

import (
	"log"
	"net/http"
	"testing"

	"github.com/ErdajtSopjani/LikeMe_API/internal/handlers"
	_users "github.com/ErdajtSopjani/LikeMe_API/internal/handlers/users"
	"github.com/ErdajtSopjani/LikeMe_API/tests"
)

func TestFollow(t *testing.T) {
	// connect to test database
	db := tests.SetupTestDB(t)

	// setup the db with the required entries to run tests
	tests.SetupDBEntries("followTests.sql", db, t)

	var users []handlers.User
	db.Find(&users)
	log.Printf("\n\n\n\n\n\n\nUsers after SetupDBEntries: %+v\n\n\n\n\n\n\n\n\n", users)

	var user_tokens []handlers.UserToken
	db.Find(&user_tokens)
	log.Printf("\n\n\n\n\n\n\nUser Tokens after SetupDBEntries: %+v\n\n\n\n\n\n\n\n\n", user_tokens)

	var follows []handlers.Follow
	db.Find(&follows)
	log.Printf("\n\n\n\n\n\n\nFollows after SetupDBEntries: %+v\n\n\n\n\n\n\n\n\n", follows)

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
				"Authorization": "token1",
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
			ExpectedBody: "",
			QueryParams:  "",
		},
	}

	tests.RunTests(db, t, testCases, "/api/v1/follow", _users.FollowAccount(db))
}
