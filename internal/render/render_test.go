package render

import (
	"github.com/VishalTanwani/GolangWebApp/internal/modals"
	"net/http"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var td modals.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "123")
	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}

}

func TestTemplates(t *testing.T) {
	pathToTemplates = "./../../templates"

	tc, err := CreateTemplatesCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	var w myWriter
	err = Templates(&w, r, "home.page.tmpl", &modals.TemplateData{})
	if err != nil {
		t.Error("error writing template to browser")
	}

	err = Templates(&w, r, "not-available.page.tmpl", &modals.TemplateData{})
	if err == nil {
		t.Error("render template that does not exist")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplatesCache(t *testing.T) {
	pathToTemplates = "./../../template"
	_, err := CreateTemplatesCache()
	if err != nil {
		t.Error(err)
	}
}
