package config

import (
	"log"
	"github.com/alexedwards/scs/v2"
	"html/template"
)

//AppConfig hold the application config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog *log.Logger
	ErrorLog *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
}
