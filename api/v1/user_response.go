package v1

import "github.com/khisakuni/kommunicake/models"

type userResponse struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}
