package main

import (
	"net/http"

	"github.com/andrewburian/powermux"
	"github.com/go-pg/pg"
)

// User is an application user and all their data
type User struct {
	Username string
	Points   uint
	Password []byte
}

// UserHandler is the router for user requests
type UserHandler struct {
	db *pg.DB
}

// Setup creates routes for the user handler
func (h *UserHandler) Setup(r *powermux.Route) {
	r.PostFunc(h.CreateUser)
	r.Route("/:name").GetFunc(h.GetUser)
}

// CreateUser sets up a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Work in Progress", http.StatusNotImplemented)
}

// GetUser gets publically accessible profile information for a user
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Work in Progress", http.StatusNotImplemented)
}
