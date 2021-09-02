package main

import (
	"encoding/gob"
	"fmt"
	"github.com/VishalTanwani/GolangWebApp/internal/config"
	"github.com/VishalTanwani/GolangWebApp/internal/driver"
	"github.com/VishalTanwani/GolangWebApp/internal/handler"
	"github.com/VishalTanwani/GolangWebApp/internal/helpers"
	"github.com/VishalTanwani/GolangWebApp/internal/modals"
	"github.com/VishalTanwani/GolangWebApp/internal/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"time"
)

const port = ":5000"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	db, err := run()
	if err != nil {
		log.Println("error at run in main", err)
		return
	}

	defer db.SQL.Close()

	defer close(app.MailChan)

	//listening for mails to send
	go listenForMail()

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

func run() (*driver.DB, error) {
	//what i am going to put in session
	gob.Register(modals.Reservation{})
	gob.Register(modals.User{})
	gob.Register(modals.Room{})
	gob.Register(modals.Restriction{})

	mailChan := make(chan modals.MailData)
	app.MailChan = mailChan
	//change this to true in production
	app.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()

	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//connect to database
	fmt.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=GolangWebApp user=vishal password=")
	if err != nil {
		log.Fatal("cannot connect to database ", err)
	}

	tc, err := render.CreateTemplatesCache()
	if err != nil {
		log.Fatal("can not create template cache", err)
		return nil, err
	}
	app.TemplateCache = tc
	app.UseCache = false

	repo := handler.NewRepo(&app, db)
	handler.NewHandler(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
