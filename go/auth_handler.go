package main

import (
	"net/http"
)

// Auth struct
type Auth struct {
}

// NewAuth retusn Auth
func NewAuth() *Auth {
	return &Auth{}
}

// ServerHTTP is negroni.Handler
func (a *Auth) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	if r.Method == `GET` {
		next(rw, r)
		return
	}

	token := r.Header.Get("X-Mauth-Token")
	if len(token) == 0 {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
		return
	}

	next(rw, r)
}
