package main

import (
	"net/http"

	"github.com/AndrewBurian/adulting-api/middlewares"
	"github.com/AndrewBurian/powermux"
	"github.com/go-pg/pg"
)

// UserHandler is the router for user requests
type UserHandler struct {
	db *pg.DB
}

// Setup creates routes for the user handler
func (h *UserHandler) Setup(r *powermux.Route) {
	r.PostFunc(h.CreateUser)
	r.Route("/:name").
		MiddlewareFunc(middlewares.RequireAuth).
		GetFunc(h.GetUser)
}

// CreateUser sets up a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Work in Progress (create)", http.StatusNotImplemented)
}

// GetUser gets publically accessible profile information for a user
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	user := powermux.PathParam(r, "name")
	if user == "me" {
		h.GetSelf(w, r)
		return
	}

	http.Error(w, "Work in Progress (get)", http.StatusNotImplemented)
}

func (h *UserHandler) GetSelf(w http.ResponseWriter, r *http.Request) {

	WriteResponse(w, r, middlewares.GetUser(r), http.StatusOK)
}
