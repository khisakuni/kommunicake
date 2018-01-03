package v1

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/khisakuni/kommunicake/api/middleware"
	"github.com/khisakuni/kommunicake/database"
	"github.com/khisakuni/kommunicake/models"

	"golang.org/x/crypto/bcrypt"
)

type userParams struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

// RegisterUser creates user record if valid
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var u userParams
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	if u.Password != u.PasswordConfirmation {
		err = errors.New("passwords do not match")
		errorResponse(w, err, http.StatusBadRequest)
		return
	}

	user, err := createUser(u, middleware.GetDBFromContext(r.Context()))
	if err != nil {
		err = errors.New("Duplicate email")
		errorResponse(w, err, http.StatusBadRequest)
	}

	token, err := user.GenerateToken()
	if err != nil {
		errorResponse(w, err, http.StatusInternalServerError)
		return
	}

	res := struct {
		User  models.User
		Token string
	}{
		User:  user,
		Token: token,
	}

	jsonResponse(w, res, http.StatusCreated)
}

func createUser(u userParams, db *database.DB) (models.User, error) {
	var user models.User
	hash, err := encryptPassword(u.Password)
	if err != nil {
		return user, err
	}
	user = models.User{Name: u.Name, Email: u.Email, EncryptedPassword: hash}
	if err = user.Create(db); err != nil {
		return user, err
	}
	return user, nil
}

func encryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
