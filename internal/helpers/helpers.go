package helpers

import (
	"fmt"
	"github.com/VishalTanwani/GolangWebApp/internal/config"
	"net/http"
	"runtime/debug"
)

var app *config.AppConfig

//NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

//ClientError will handle client errors
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("client error with the status of ", status)
	http.Error(w, http.StatusText(status), status)
}

//ServerError will handle server errors
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

//IsAuthenticate will check is user exist in our system or not
func IsAuthenticate(r *http.Request) bool {
	exist := app.Session.Exists(r.Context(), "user_id")
	return exist
}
