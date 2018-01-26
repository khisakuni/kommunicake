package v1

import (
	"net/http"

	helpers "github.com/khisakuni/kommunicake/api/helpers"
	middleware "github.com/khisakuni/kommunicake/api/middleware"
	"github.com/khisakuni/kommunicake/models"
)

// MessageProvidersIndex responds with a list of message providers
func MessageProvidersIndex(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := middleware.GetUserFromContext(ctx)
	db := middleware.GetDBFromContext(ctx)

	var messageProviders []models.UserMessageProvider
	db.Conn.Where("user_id = ?", user.ID).Find(&messageProviders)
	if db.Conn.Error != nil {
		helpers.ErrorResponse(w, db.Conn.Error, http.StatusInternalServerError)
		return
	}

	jsonResponse(w, messageProviders, http.StatusOK)
}
