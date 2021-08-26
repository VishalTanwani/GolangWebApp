package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/VishalTanwani/GolangWebApp/internal/config"
	"github.com/VishalTanwani/GolangWebApp/internal/forms"
	"github.com/VishalTanwani/GolangWebApp/internal/modals"
	"github.com/VishalTanwani/GolangWebApp/internal/render"
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

	render.Templates(w, r, "home.page.tmpl", &modals.TemplateData{})
}

//About is about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "hello again!"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP
	render.Templates(w, r, "about.page.tmpl", &modals.TemplateData{
		StringMap: stringMap,
	})
}

//Reservation is for make reservatin and display forms page handler
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "make-reservation.page.tmpl", &modals.TemplateData{})
}

//Generals is room page handler
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "generals.page.tmpl", &modals.TemplateData{})
}

//Majors is room page handler
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "majors.page.tmpl", &modals.TemplateData{})
}

//Availability is room availability page handler
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "search-availability.page.tmpl", &modals.TemplateData{})
}

type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

//AvailabilityJSON handle request and send json response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		Ok:      true,
		Message: "Available",
	}

	out, err := json.MarshalIndent(resp, "", "  ")
	if err != nil {
		fmt.Println("error at converting obj to json", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

//PostAvailability is room availability page handler
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hellokasjdfb"))
}

//Contact is contact page handler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "contact.page.tmpl", &modals.TemplateData{})
}

//MakeReservation is make reservation form page handler
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation modals.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.Templates(w, r, "make-reservation.page.tmpl", &modals.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

//PostReservation is for handling reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error at parsing the fofm", err)
		return
	}

	reservation := modals.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	// form.Has("first_name", r)
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.IsEmail("email", r)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Templates(w, r, "make-reservation.page.tmpl", &modals.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

//ReservationSummary is Reservation-Summary page handler
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation,ok := m.App.Session.Get(r.Context(),"reservation").(modals.Reservation)
	if !ok {
		fmt.Println("cannot get data from session")
		m.App.Session.Put(r.Context(),"error","cant get reservation rom session")
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Templates(w, r, "reservation-summary.page.tmpl", &modals.TemplateData{
		Data: data,
	})
}