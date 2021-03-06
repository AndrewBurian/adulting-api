package main

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/AndrewBurian/adulting-api/data"
	"github.com/AndrewBurian/powermux"
	"github.com/sirupsen/logrus"
)

// PasswordAuthReq is the incoming auth request
type PasswordAuthReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthReply is an outgoing authentication response
type AuthReply struct {
	Token string `json:"token"`
}

// AuthHandler is the router for handling authentication
type AuthHandler struct {
	db     data.UserDAL
	logger *logrus.Entry
}

// Setup creates routes for the auth handler
func (h *AuthHandler) Setup(r *powermux.Route) {
	r.Route("/password").PostFunc(h.PasswordAuth)
	r.Route("/logout").PostFunc(h.Logout)
}

// PasswordAuth processes a password attempt
func (h *AuthHandler) PasswordAuth(w http.ResponseWriter, r *http.Request) {

	log := h.logger

	// Parse request
	req := &PasswordAuthReq{}
	if err := DecodeJSON(r, req); err != nil {
		log.WithError(err).Warn("Bad Request")
		http.Error(w, "Could not parse request", http.StatusBadRequest)
		return
	}

	// Fetch user password
	user := &data.User{
		Username: req.Username,
	}
	log = log.WithField("user", req.Username)

	err := h.db.GetUser(user)
	if err == data.ErrNotFound {
		log.Debug("User not found")
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		log.WithError(err).Error("Database error")
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	// Check password
	if bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password)) != nil {
		log.Debug("Wrong Password")
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}

	log.Debug("User authenticated")

	// Send token
	token := &AuthReply{
		Token: "no", //TODO
	}
	if err = WriteResponse(w, r, token, http.StatusOK); err != nil {
		log.WithError(err).Error("Unable to send response")
	}
}

// Logout invalidates the current access token
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	//TODO
	http.Error(w, "Work in progress", http.StatusNotImplemented)
	return
}
