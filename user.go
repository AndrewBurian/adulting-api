package main

import (
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
}
