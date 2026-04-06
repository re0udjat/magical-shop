package main

import (
	"context"
	"net/http"

	"github.com/re0udjat/magic-shop/internal/data"
)

type contextKey string

// Convert the string "user" to a contextKey type and assign it to the userContextKey
// constant
const userContextKey = contextKey("user")

// Returns a new copy of the request with the provided User struct added to the context
func (app *app) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// Retrieves the User struct from the request context
func (app *app) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
