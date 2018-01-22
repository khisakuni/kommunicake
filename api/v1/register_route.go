package v1

import (
	"errors"
	"net/http"

	"github.com/khisakuni/kommunicake/api/middleware"
	"github.com/khisakuni/kommunicake/database"
	"github.com/khisakuni/kommunicake/models"

	helpers "github.com/khisakuni/kommunicake/api/helpers"
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
	if ok := decodeJSON(w, r.Body, &u); !ok {
		return
	}

	validate := func() error {
		var err error
		if u.Password != u.PasswordConfirmation {
			err = errors.New("passwords do not match")
		}
		return err
	}
	if ok := validateParams(w, validate); !ok {
		return
	}

	db := middleware.GetDBFromContext(r.Context())
	user, err := createUser(u, db)
	if err != nil {
		err = errors.New("Duplicate email")
		helpers.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	token, err := user.GenerateToken(db)
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	res := userResponse{User: user, Token: token.Value}
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
