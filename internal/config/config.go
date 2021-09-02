package config

import (
	"github.com/VishalTanwani/GolangWebApp/internal/modals"
	"github.com/alexedwards/scs/v2"
	"html/template"
	"log"
)

//AppConfig hold the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan modals.MailData
}
