package handler

import (
	"encoding/json"
	// "errors"
	"fmt"
	"github.com/VishalTanwani/GolangWebApp/internal/config"
	"github.com/VishalTanwani/GolangWebApp/internal/driver"
	"github.com/VishalTanwani/GolangWebApp/internal/forms"
	"github.com/VishalTanwani/GolangWebApp/internal/helpers"
	"github.com/VishalTanwani/GolangWebApp/internal/modals"
	"github.com/VishalTanwani/GolangWebApp/internal/render"
	"github.com/VishalTanwani/GolangWebApp/internal/repository"
	"github.com/VishalTanwani/GolangWebApp/internal/repository/dbrepo"
	// "github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
	"strings"
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

//NewTestRepo creates new Repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
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
	Ok        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

//AvailabilityJSON handle request and send json response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			Ok:      false,
			Message: "error at parsing form",
		}

		out, _ := json.MarshalIndent(resp, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")
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

	availability, err := m.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		resp := jsonResponse{
			Ok:      false,
			Message: "error connecting database",
		}

		out, _ := json.MarshalIndent(resp, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	resp := jsonResponse{
		Ok:        availability,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomID),
	}

	out, _ := json.MarshalIndent(resp, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

//PostAvailability is room availability page handler
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot parse the form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")
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

	rooms, err := m.DB.SearchAvailbilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No rooms available")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	reservation := modals.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	render.Templates(w, r, "choose-room.page.tmpl", &modals.TemplateData{
		Data: data,
	})
}

//Contact is contact page handler
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "contact.page.tmpl", &modals.TemplateData{})
}

//MakeReservation is make reservation form page handler
func (m *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(modals.Reservation)
	if !ok {
		// helpers.ServerError(w, errors.New("Cannot get reservation out from sessiion"))
		m.App.Session.Put(r.Context(), "error", "Cannot get reservation out from sessiion")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	room, err := m.DB.GetRoomByID(reservation.RoomID)
	if err != nil {
		// helpers.ServerError(w, err)
		m.App.Session.Put(r.Context(), "error", "Cant find a room")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	reservation.Room = room

	m.App.Session.Put(r.Context(), "reservation", reservation)
	//converting time.time go's format to given date format because we are going to show this in html page
	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Templates(w, r, "make-reservation.page.tmpl", &modals.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

//PostReservation is for handling reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(modals.Reservation)
	if !ok {
		// helpers.ServerError(w, errors.New("Cannot get reservation out from sessiion"))
		m.App.Session.Put(r.Context(), "error", "Cannot get reservation out from sessiion")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot parse the form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")
	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	form := forms.New(r.PostForm)

	//form validation
	// form.Has("first_name")
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Error(w, "form is not valid", http.StatusSeeOther)
		render.Templates(w, r, "make-reservation.page.tmpl", &modals.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	//sending data to db
	newReservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot insert the reservation data")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "reservation", reservation)

	restriction := modals.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot insert room restriction")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	//send mail to guest
	htmlMessage := fmt.Sprintf(`
	<h1>happy bang</h1>
		<strong>Reservation Confirmation</strong><br>
		Dear %s, <br>
		this is to confirm ur reservation from %s to %s
	`, reservation.FirstName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))
	msg := modals.MailData{
		From:     "jhon@cena.com",
		To:       reservation.Email,
		Subject:  "Reservation Confirmation",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	m.App.MailChan <- msg

	//send mail to owner
	htmlMessage = fmt.Sprintf(`
		<strong>Reservation Confirmation</strong><br>
		Dear jhon cena, <br>
		there is a one reservation in %s from %s to %s
	`, reservation.Room.RoomName, reservation.StartDate.Format("2006-01-02"), reservation.EndDate.Format("2006-01-02"))
	msg = modals.MailData{
		From:    "reservation@no-reply.com",
		To:      "jhon@cena.com",
		Subject: "Reservation Confirmation",
		Content: htmlMessage,
	}

	m.App.MailChan <- msg

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
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")
	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Templates(w, r, "reservation-summary.page.tmpl", &modals.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

//ChooseRoom is our chose room page
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	// changed to this, so we can test it more easily
	// split the URL up by /, and grab the 3rd element
	exploded := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(exploded[2])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(modals.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	reservation.RoomID = roomID
	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

//BookRoom to take user to make reservation room and some data to session
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")
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
	var reservation modals.Reservation
	reservation.RoomID = roomID
	reservation.StartDate = startDate
	reservation.EndDate = endDate

	room, err := m.DB.GetRoomByID(reservation.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation.Room = room

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

//ShowLogin will show login page
func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "login.page.tmpl", &modals.TemplateData{
		Form: forms.New(nil),
	})
}

//PostShowLogin handels the login form
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		render.Templates(w, r, "login.page.tmpl", &modals.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "invalid email and password")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

//Logout logout out user
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

//AdminDashboard will open admin dash board
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Templates(w, r, "admin-dashboard.page.tmpl", &modals.TemplateData{})
}

//AdminNewReservations will open admin dash board which will show all new reervations
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.GetAllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.Templates(w, r, "admin-new-reservations.page.tmpl", &modals.TemplateData{
		Data: data,
	})
}

//AdminAllReservations will open admin dash board which will show all reervations
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.GetAllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.Templates(w, r, "admin-all-reservations.page.tmpl", &modals.TemplateData{
		Data: data,
	})
}

//AdminCalendarReservations will open admin dash board
func (m *Repository) AdminCalendarReservations(w http.ResponseWriter, r *http.Request) {
	//assume that there is no month/year specified
	now := time.Now()
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if r.URL.Query().Get("y") != "" {
		year, err := strconv.Atoi(r.URL.Query().Get("y"))
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		month, err := strconv.Atoi(r.URL.Query().Get("m"))
		if err != nil {
			helpers.ServerError(w, err)
			return
		}
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	data := make(map[string]interface{})
	data["now"] = now

	nextMonth := next.Format("01")
	nextYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastYear := last.Format("2006")

	stringMap := make(map[string]string)

	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastYear
	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	//get the first and last days of the month
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastOfMonth.Day()

	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["rooms"] = rooms

	for _, x := range rooms {
		//create maps
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}

		//get all restrictions
		restrictions, err := m.DB.GetRestrictionsForRoomByDate(x.ID, firstOfMonth, lastOfMonth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		for _, y := range restrictions {
			if y.ReservationID > 0 {
				//reservations
				// fmt.Println(y.ReservationID)
				for d := y.StartDate; d.After(y.EndDate) == false; d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = y.ReservationID
				}
				// fmt.Println(reservationMap)
			} else {
				blockMap[y.StartDate.Format("2006-01-2")] = y.ID

			}
		}

		data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap

		m.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap)
	}

	render.Templates(w, r, "admin-calendar-reservations.page.tmpl", &modals.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}

//AdminShowReservations will open admin dash board
func (m *Repository) AdminShowReservations(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(url[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := url[3]

	stringMap := make(map[string]string)
	stringMap["src"] = src

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	stringMap["year"] = year
	stringMap["month"] = month

	reservation, err := m.DB.GetReservationByID(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Templates(w, r, "admin-show-reservations.page.tmpl", &modals.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil),
	})
}

//AdminPostShowReservations will open admin dash board
func (m *Repository) AdminPostShowReservations(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(url[4])
	if err != nil {
		helpers.ServerError(w, err)
		fmt.Println(err)
		return
	}

	src := url[3]

	reservation, err := m.DB.GetReservationByID(id)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	err = m.DB.UpdateReservation(reservation)
	if err != nil {
		fmt.Println(err)
		helpers.ServerError(w, err)
		return
	}

	year := r.Form.Get("year")
	month := r.Form.Get("month")

	m.App.Session.Put(r.Context(), "flash", "Reservation updated")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

//AdminProcessReservation will open process a reservation
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(url[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := url[3]

	err = m.DB.UpdateProcssedForReservation(id, 1)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

//AdminDeleteReservation will delete a reservation
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	url := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(url[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := url[3]

	err = m.DB.UpdateProcssedForReservation(id, 1)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")

	m.App.Session.Put(r.Context(), "flash", "Reservation deleted")
	if year == "" {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s", year, month), http.StatusSeeOther)
	}
}

//AdminPostCalendarReservations will handle form submit of admin calender page
func (m *Repository) AdminPostCalendarReservations(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	year, _ := strconv.Atoi(r.Form.Get("y"))
	month, _ := strconv.Atoi(r.Form.Get("m"))

	//process blocks
	rooms, err := m.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	for _, x := range rooms {
		curMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", x.ID)).(map[string]int)
		for name, value := range curMap {
			if val, ok := curMap[name]; ok {
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", x.ID, name)) {
						err = m.DB.DeleteBlockByID(value)
						if err != nil {
							helpers.ServerError(w, err)
							return
						}
					}
				}
			}
		}
	}

	//now handle new blocks
	formValues := r.PostForm
	for name, _ := range formValues {
		if strings.HasPrefix(name, "add_block") {
			explode := strings.Split(name, "_")
			id, _ := strconv.Atoi(explode[2])
			date, _ := time.Parse("2006-01-2", explode[3])
			err = m.DB.InsertBlockForRoom(id, date)
			if err != nil {
				helpers.ServerError(w, err)
				return
			}
		}
	}

	m.App.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d", year, month), http.StatusSeeOther)
}
