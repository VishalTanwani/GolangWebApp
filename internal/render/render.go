package render

import (
	"bytes"
	"fmt"
	"github.com/justinas/nosurf"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"github.com/VishalTanwani/GolangWebApp/internal/config"
	"github.com/VishalTanwani/GolangWebApp/internal/modals"
)

var functions = template.FuncMap{}

var app *config.AppConfig
var pathToTemplates = "./templates"

//NewTemplates set the template for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func addDefaultData(td *modals.TemplateData, r *http.Request) *modals.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

//Templates is for rendering html files
func Templates(w http.ResponseWriter, r *http.Request, tmpl string, td *modals.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplatesCache()
	}
	t, ok := tc[tmpl]

	if !ok {
		log.Fatal("cannot get template from template cache")
	}

	buf := new(bytes.Buffer)

	td = addDefaultData(td, r)
	_ = t.Execute(buf, td)
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error at writing in buf", err)
	}
}

//CreateTemplatesCache will create our our applications cache
func CreateTemplatesCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl",pathToTemplates))
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl",pathToTemplates))
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl",pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil

}
