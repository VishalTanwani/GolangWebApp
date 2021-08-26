package main

import (
	"github.com/justinas/nosurf"
	"net/http"
)

//Nosurf is our middle ware that add CSRF protection to all post request
func Nosurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

//SessionLoad is our session middleware which load and save session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
