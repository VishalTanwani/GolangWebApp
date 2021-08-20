package handler

import (
	"fmt"
	"net/http"
	"webApp/pkg/config"
	"webApp/pkg/modals"
	"webApp/pkg/render"
)

//Repository is repository type
type Repository struct {
	App *config.AppConfig
}

//Repo used by the handlers
var Repo *Repository

//NewRepo creates new Repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

//NewHandler sets the repository with handlers
func NewHandler(r *Repository) {
	Repo = r
}

//Home is home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	fmt.Println(remoteIP)
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Templates(w, "home.page.tmpl", &modals.TemplateData{})
}

//About is about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "hello again!"
	
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.Templates(w, "about.page.tmpl", &modals.TemplateData{
		StringMap: stringMap,
	})
}
