package handlers

/*import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm"
)*/

type User struct {
	Email     string `json:"email"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Bio       string `json:"bio"`
	Token     string `json:"token"`
}
