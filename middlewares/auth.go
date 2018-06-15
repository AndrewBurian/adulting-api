package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/AndrewBurian/adulting-api/data"
	log "github.com/sirupsen/logrus"
)

type AuthDetection struct {
	DB data.UserDAL
}

type authCtxKey string

var userKey = authCtxKey("user")

// ServeHTTPMiddleware uses the Authorization header if present to associate a user with this request.
// It does NOT block requests with no Authorization
func (h *AuthDetection) ServeHTTPMiddleware(w http.ResponseWriter, r *http.Request, next func(http.ResponseWriter, *http.Request)) {
	auth := r.Header.Get("Authorization")

	// skip if no auth token
	if auth == "" {
		next(w, r)
		return
	}

	parts := strings.Split(auth, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		http.Error(w, "Invalid Authorization header", http.StatusBadRequest)
		log.Warn("Bad authorization header")
		return
	}

	auth = parts[1]

	user, err := h.DB.GetUserByToken(auth)
	if err == data.ErrNotFound {
		http.Error(w, "Invalid Authorization header", http.StatusBadRequest)
		log.Warn("Unknown auth token")
		return
	}
	if err != nil {
		http.Error(w, "Error processing Authorization header", http.StatusInternalServerError)
		log.WithError(err).Error("Database Error")
		return
	}

	r = r.WithContext(context.WithValue(r.Context(), userKey, user))
	next(w, r)

}

// GetUser pulls the user from the auth middleware
func GetUser(r *http.Request) *data.User {
	user := r.Context().Value(userKey)
	if user == nil {
		return nil
	}
	return user.(*data.User)
}

// RequireAuth prevents a request that hasn't detected a user from passing
func RequireAuth(w http.ResponseWriter, r *http.Request, next func(http.ResponseWriter, *http.Request)) {
	if GetUser(r) == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Warn("Unauthorized request")
		return
	}

	next(w, r)
}
