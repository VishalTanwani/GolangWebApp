package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"time"
	"webApp/pkg/config"
	"webApp/pkg/handler"
	"webApp/pkg/render"
)

const port = ":5000"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	app.InProduction = false

	session = scs.New()

	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplatesCache()
	if err != nil {
		log.Fatal("can not create template cache", err)
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handler.NewRepo(&app)
	handler.NewHandler(repo)
	render.NewTemplates(&app)

	server := &http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	fmt.Println("server is running in 5000 port")
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("error at running server", err)
	}
}
