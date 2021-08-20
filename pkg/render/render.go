package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"webApp/pkg/config"
	"webApp/pkg/modals"
)

var functions = template.FuncMap{}

var app *config.AppConfig

//NewTemplates set the template for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func addDefaultData(td *modals.TemplateData) *modals.TemplateData {
	return td
}

//Templates is for rendering html files
func Templates(w http.ResponseWriter, tmpl string, td *modals.TemplateData) {
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

	td = addDefaultData(td)
	_ = t.Execute(buf, td)
	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("error at writing in buf", err)
	}
}

//CreateTemplatesCache will create our our applications cache
func CreateTemplatesCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil

}
