package handler

import (
	"encoding/json"
	"fmt"
	"github.com/VishalTanwani/GolangWebApp/internal/config"
	"github.com/VishalTanwani/GolangWebApp/internal/driver"
	"github.com/VishalTanwani/GolangWebApp/internal/forms"
	"github.com/VishalTanwani/GolangWebApp/internal/helpers"
	"github.com/VishalTanwani/GolangWebApp/internal/modals"
	"github.com/VishalTanwani/GolangWebApp/internal/render"
	"github.com/VishalTanwani/GolangWebApp/internal/repository"
	"github.com/VishalTanwani/GolangWebApp/internal/repository/dbrepo"
	"net/http"
	"strconv"
	"time"
)

//Repository is repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

//Repo used by the handlers
var Repo *Repository

//NewRepo creates new Repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

//NewHandler sets the repository with handlers
func NewHandler(r *Repository) {
	Repo = r
}

//Home is home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "home.page.tmpl", &modals.TemplateData{})
}

//About is about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "about.page.tmpl", &modals.TemplateData{})
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
		helpers.ServerError(w, err)
		return
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
		helpers.ServerError(w, err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")
	fmt.Println(sd, ed)
	// go has defernt format for date like this Mon Jan 2 15:04:05 -0700 MST 2006
	// so we have to provide ours
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := modals.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	form := forms.New(r.PostForm)

	//form validation
	// form.Has("first_name")
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Templates(w, r, "make-reservation.page.tmpl", &modals.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	//sending data to db
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	restriction := modals.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	//writing reservatoin data to session
	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

//ReservationSummary is Reservation-Summary page handler
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(modals.Reservation)
	if !ok {
		m.App.ErrorLog.Println("cant get error from session")
		m.App.Session.Put(r.Context(), "error", "cant get reservation rom session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Templates(w, r, "reservation-summary.page.tmpl", &modals.TemplateData{
		Data: data,
	})
}
