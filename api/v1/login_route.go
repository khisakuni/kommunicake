package v1

import (
	"fmt"
	"errors"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/khisakuni/kommunicake/api/middleware"
	"github.com/khisakuni/kommunicake/models"
	helpers "github.com/khisakuni/kommunicake/api/helpers"
)

type loginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "cool")
}

// Login authenticates user and returns token
func Login(w http.ResponseWriter, r *http.Request) {
	var p loginParams
	if ok := decodeJSON(w, r.Body, &p); !ok {
		return
	}

	validate := func() error {
		var err error
		if p.Email == "" || p.Password == "" {
			err = errors.New("Please provide valid email and password")
		}
		return err
	}
	if ok := validateParams(w, validate); !ok {
		return
	}

	db := middleware.GetDBFromContext(r.Context())
	var user models.User
	db.Conn.Where("email = ?", p.Email).First(&user)
	if user.Email == "" {
		err := errors.New("User with that email does not exist")
		helpers.ErrorResponse(w, err, http.StatusNotFound)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(p.Password))
	if err != nil {
		err = errors.New("Password is incorrect")
		helpers.ErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	token, err := user.GenerateToken(db)
	if err != nil {
		helpers.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	res := userResponse{User: user, Token: token.Value}
	jsonResponse(w, res, http.StatusOK)
}
