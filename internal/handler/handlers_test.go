package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
	"net/url"
	"github.com/VishalTanwani/GolangWebApp/internal/driver"
	"github.com/VishalTanwani/GolangWebApp/internal/modals"
	"log"
	"reflect"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"general-quarters", "/general-quarters", "GET", http.StatusOK},
	{"major-suite", "/major-suite", "GET", http.StatusOK},
	{"search-availability", "/search-availability", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"something not found", "/something", "GET", http.StatusNotFound},
	{"login", "/user/login", "GET", http.StatusOK},
	{"logout", "/user/logout", "GET", http.StatusOK},
	{"admin-dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"admin-reservations-new", "/admin/reservations-new", "GET", http.StatusOK},
	{"admin-reservations-all", "/admin/reservations-all", "GET", http.StatusOK},
	{"admin-reservations-show", "/admin/reservations/new/1/show", "GET", http.StatusOK},
	{"admin-reservations-calender", "/admin/reservations-calendar", "GET", http.StatusOK},
	{"admin-reservations-calender", "/admin/reservations-calendar?y=2021&m=09", "GET", http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()
	for _, e := range theTests {
		resp, err := ts.Client().Get(ts.URL + e.url)
		if err != nil {
			t.Log(err)
			t.Fatal(err)
		}
		if resp.StatusCode != e.expectedStatusCode {
			t.Errorf("for %s expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		}

	}
}

func TestRepository_MakeReservation(t *testing.T) {
	reservation := modals.Reservation{
		RoomID: 1,
		Room: modals.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	//test where there is a session
	r, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getContext(r)
	r = r.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.MakeReservation)
	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong code : wanted %d this get this %d", http.StatusOK, rr.Code)
	}

	//test where there is no session
	r, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getContext(r)
	r = r.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.MakeReservation)
	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong code : wanted %d this get this %d", http.StatusSeeOther, rr.Code)
	}

	//test where there is a session and cant find a room
	r, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getContext(r)
	r = r.WithContext(ctx)

	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler = http.HandlerFunc(Repo.MakeReservation)
	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong code : wanted %d this get this %d", http.StatusSeeOther, rr.Code)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	sd := "2021-09-01"
	ed := "2021-09-03"

	layout := "2006-01-02"

	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)
	reservation := modals.Reservation{
		RoomID:    1,
		StartDate: startDate,
		EndDate:   endDate,
		Room: modals.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	reqBody := "first_name=vishal"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=tanwani")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=tanwani@vishal.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=987654321")

	r, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getContext(r)
	r = r.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("post Reservation handler returned wrong code : wanted %d this get this %d", http.StatusSeeOther, rr.Code)
	}

	//test where there is no session
	r, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getContext(r)
	r = r.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong code : wanted %d this get this %d", http.StatusSeeOther, rr.Code)
	}

	//missing request body
	r, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getContext(r)
	r = r.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("post Reservation handler returned wrong code : wanted %d this get this %d", http.StatusSeeOther, rr.Code)
	}

	//form is not valid
	reqBody = "first_name=v"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=tanwani")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=tanwani@vishal.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=987654321")

	r, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getContext(r)
	r = r.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("post Reservation handler returned wrong code : wanted %d this get this %d", http.StatusSeeOther, rr.Code)
	}

	//insert reservation db error
	reservation = modals.Reservation{
		RoomID:    2,
		StartDate: startDate,
		EndDate:   endDate,
		Room: modals.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}
	reqBody = "first_name=vishal"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=tanwani")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=tanwani@vishal.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=987654321")

	r, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getContext(r)
	r = r.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("post Reservation handler returned wrong code : wanted %d this get this %d", http.StatusSeeOther, rr.Code)
	}

	//insert restriction db error
	reservation = modals.Reservation{
		RoomID:    3,
		StartDate: startDate,
		EndDate:   endDate,
		Room: modals.Room{
			ID:       3,
			RoomName: "General's Quarters",
		},
	}
	reqBody = "first_name=vishal"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=tanwani")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=tanwani@vishal.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=987654321")

	r, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getContext(r)
	r = r.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("post Reservation handler returned wrong code : wanted %d this get this %d", http.StatusSeeOther, rr.Code)
	}

}

func TestRepository_AvailabilityJSON(t *testing.T) {
	//rooms are not available
	reqBody := "start_date=2021-09-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2021-09-03")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	r, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getContext(r)
	r = r.WithContext(ctx)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, r)

	var j jsonResponse
	err := json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json", err)
	}

	if j.Ok {
		t.Error("Got availability room when rooms are not available")
	}

	//room is available
	reqBody = "start_date=2021-09-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2021-09-03")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	r, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getContext(r)
	r = r.WithContext(ctx)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, r)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json", err)
	}

	if !j.Ok {
		t.Error("Got not availability when room was available ")
	}

	//missing request body
	r, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getContext(r)
	r = r.WithContext(ctx)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, r)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json", err)
	}

	//database error
	reqBody = "start_date=2021-09-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2021-09-03")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=3")

	r, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getContext(r)
	r = r.WithContext(ctx)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, r)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)
	if err != nil {
		t.Error("failed to parse json", err)
	}

	if j.Ok && j.Message != "error connecting database" {
		t.Error("error connecting database")
	}

}

func TestRepository_PostAvailability(t *testing.T) {
	//rooms are not available
	reqBody := "start_date=2021-11-03"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2021-11-03")

	r, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getContext(r)
	r = r.WithContext(ctx)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong code : wanted %d this get this %d", http.StatusSeeOther, rr.Code)
	}

	//rooms are available
	reqBody = "start_date=2021-11-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2021-11-03")

	r, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getContext(r)
	r = r.WithContext(ctx)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong code : wanted %d this get this %d", http.StatusOK, rr.Code)
	}

	//rooms are available

	r, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getContext(r)
	r = r.WithContext(ctx)

	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("Reservation handler returned wrong code : wanted %d this get this %d", http.StatusSeeOther, rr.Code)
	}

}

func TestRepository_ChooseRoom(t *testing.T) {
	sd := "2021-09-01"
	ed := "2021-09-03"

	layout := "2006-01-02"

	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)
	reservation := modals.Reservation{
		RoomID:    1,
		StartDate: startDate,
		EndDate:   endDate,
		Room: modals.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	r, _ := http.NewRequest("GET", "/choose-room/1", nil)
	ctx := getContext(r)
	r = r.WithContext(ctx)
	r.RequestURI = "/choose-room/1"
	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("post Reservation handler returned wrong code : wanted %d this get this %d", http.StatusSeeOther, rr.Code)
	}
}

func TestRepository_ReservationSummary(t *testing.T) {
	sd := "2021-09-01"
	ed := "2021-09-03"

	layout := "2006-01-02"

	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)
	reservation := modals.Reservation{
		RoomID:    1,
		StartDate: startDate,
		EndDate:   endDate,
		Room: modals.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	r, _ := http.NewRequest("POST", "/reservation-summary", nil)
	ctx := getContext(r)
	r = r.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusOK {
		t.Errorf("post Reservation handler returned wrong code : wanted %d this get this %d", http.StatusOK, rr.Code)
	}

	//no session
	r, _ = http.NewRequest("POST", "/reservation-summary", nil)
	ctx = getContext(r)
	r = r.WithContext(ctx)

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.ReservationSummary)
	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("post Reservation handler returned wrong code : wanted %d this get this %d", http.StatusSeeOther, rr.Code)
	}
}

func TestRepository_BookRoom(t *testing.T) {
	sd := "2021-09-01"
	ed := "2021-09-03"

	layout := "2006-01-02"

	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)
	reservation := modals.Reservation{
		RoomID:    1,
		StartDate: startDate,
		EndDate:   endDate,
		Room: modals.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	r, _ := http.NewRequest("GET", "/book-room?s=2050-01-01&e=2050-01-02&id=1", nil)
	ctx := getContext(r)
	r = r.WithContext(ctx)
	rr := httptest.NewRecorder()

	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.BookRoom)
	handler.ServeHTTP(rr, r)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("post Reservation handler returned wrong code : wanted %d this get this %d", http.StatusSeeOther, rr.Code)
	}
}

func TestNewRepo(t *testing.T) {
	var db driver.DB
	repo := NewRepo(&app, &db)

	if reflect.TypeOf(repo).String() != "*handler.Repository" {
		t.Errorf("Did not get correct type from NewRepo: got %s, wanted *Repository", reflect.TypeOf(repo).String())
	}
}

var loginTest = []struct{
	name string
	email string
	expectedStatusCode int
	expectedHTML string
	expectedLocation string
}{
	{
		"valid-credentials",
		"admin@admin.com",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid-credentials",
		"amin@admin.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},
	{
		"invalid-data",
		"aa",
		http.StatusOK,
		`action="/user/login"`,
		"",
	},
}

func TestLogin(t *testing.T){
	//range thorugh all tests
	for _,e := range loginTest{
		postedData := url.Values{}
		postedData.Add("email",e.email)
		postedData.Add("password","01090109")

		r, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		ctx := getContext(r)
		r = r.WithContext(ctx)
		r.Header.Set("Content-Type","application/x-www-form-urlencoded")
		
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.PostShowLogin)
		handler.ServeHTTP(rr, r)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("fialed %s: expected code %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {
			//get the url from test
			realLocation,_ := rr.Result().Location()
			if realLocation.String() != e.expectedLocation {
				t.Errorf("fialed %s: expected location %s but got %s", e.name, e.expectedLocation, realLocation.String())
			}
		}

		//checking for expected html
		if e.expectedHTML != "" {
			html := rr.Body.String()
			if !strings.Contains(html,e.expectedHTML) {
				t.Errorf("fialed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

var adminPostShowReservationTests = []struct {
	name                 string
	url                  string
	postedData           url.Values
	expectedResponseCode int
	expectedLocation     string
	expectedHTML         string
}{
	{
		name: "valid-data-from-new",
		url:  "/admin/reservations/new/1/show",
		postedData: url.Values{
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"john@smith.com"},
			"phone":      {"555-555-5555"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/admin/reservations-new",
		expectedHTML:         "",
	},
	{
		name: "valid-data-from-all",
		url:  "/admin/reservations/all/1/show",
		postedData: url.Values{
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"john@smith.com"},
			"phone":      {"555-555-5555"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/admin/reservations-all",
		expectedHTML:         "",
	},
	{
		name: "valid-data-from-cal",
		url:  "/admin/reservations/cal/1/show",
		postedData: url.Values{
			"first_name": {"John"},
			"last_name":  {"Smith"},
			"email":      {"john@smith.com"},
			"phone":      {"555-555-5555"},
			"year":       {"2022"},
			"month":      {"01"},
		},
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "/admin/reservations-calendar?y=2022&m=01",
		expectedHTML:         "",
	},
}

func TestAdminPostShowReservation(t *testing.T) {
	//range thorugh all tests
	for _,e := range adminPostShowReservationTests{

		r, _ := http.NewRequest("POST", e.url, strings.NewReader(e.postedData.Encode()))
		ctx := getContext(r)
		r = r.WithContext(ctx)
		r.RequestURI = e.url

		r.Header.Set("Content-Type","application/x-www-form-urlencoded")
		
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(Repo.AdminPostShowReservations)
		handler.ServeHTTP(rr, r)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("fialed %s: expected code %d but got %d", e.name, e.expectedResponseCode, rr.Code)
		}

		if e.expectedLocation != "" {
			//get the url from test
			realLocation,_ := rr.Result().Location()
			if realLocation.String() != e.expectedLocation {
				t.Errorf("fialed %s: expected location %s but got %s", e.name, e.expectedLocation, realLocation.String())
			}
		}

		//checking for expected html
		if e.expectedHTML != "" {
			html := rr.Body.String()
			if !strings.Contains(html,e.expectedHTML) {
				t.Errorf("fialed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

var adminPostReservationCalendarTests = []struct {
	name                 string
	postedData           url.Values
	expectedResponseCode int
	expectedLocation     string
	expectedHTML         string
	blocks               int
	reservations         int
}{
	{
		name: "cal",
		postedData: url.Values{
			"y":  {time.Now().Format("2006")},
			"m": {time.Now().Format("01")},
			fmt.Sprintf("add_block_1_%s", time.Now().AddDate(0, 0, 2).Format("2006-01-2")): {"1"},
		},
		expectedResponseCode: http.StatusSeeOther,
	},
	{
		name:                 "cal-blocks",
		postedData:           url.Values{},
		expectedResponseCode: http.StatusSeeOther,
		blocks:               1,
	},
	{
		name:                 "cal-res",
		postedData:           url.Values{},
		expectedResponseCode: http.StatusSeeOther,
		reservations:         1,
	},
}

func TestPostReservationCalendar(t *testing.T) {
	for _, e := range adminPostReservationCalendarTests {
		var req *http.Request
		if e.postedData != nil {
			req, _ = http.NewRequest("POST", "/admin/reservations-calendar", strings.NewReader(e.postedData.Encode()))
		} else {
			req, _ = http.NewRequest("POST", "/admin/reservations-calendar", nil)
		}
		ctx := getContext(req)
		req = req.WithContext(ctx)

		now := time.Now()
		bm := make(map[string]int)
		rm := make(map[string]int)

		currentYear, currentMonth, _ := now.Date()
		currentLocation := now.Location()

		firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

		for d := firstOfMonth; d.After(lastOfMonth) == false; d = d.AddDate(0, 0, 1) {
			rm[d.Format("2006-01-2")] = 0
			bm[d.Format("2006-01-2")] = 0
		}

		if e.blocks > 0 {
			bm[firstOfMonth.Format("2006-01-2")] = e.blocks
		}

		if e.reservations > 0 {
			rm[lastOfMonth.Format("2006-01-2")] = e.reservations
		}

		session.Put(ctx, "block_map_1", bm)
		session.Put(ctx, "reservation_map_1", rm)

		// set the header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.AdminPostCalendarReservations)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedResponseCode {
			t.Errorf("fialed %s: expected code %d but got %d", e.name, e.expectedResponseCode, rr.Code)
		}

		if e.expectedLocation != "" {
			//get the url from test
			realLocation,_ := rr.Result().Location()
			if realLocation.String() != e.expectedLocation {
				t.Errorf("fialed %s: expected location %s but got %s", e.name, e.expectedLocation, realLocation.String())
			}
		}

		//checking for expected html
		if e.expectedHTML != "" {
			html := rr.Body.String()
			if !strings.Contains(html,e.expectedHTML) {
				t.Errorf("fialed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}

	}
}

var adminProcessReservationTests = []struct {
	name                 string
	queryParams          string
	expectedResponseCode int
	expectedLocation     string
}{
	{
		name:                 "process-reservation",
		queryParams:          "",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
	{
		name:                 "process-reservation-back-to-cal",
		queryParams:          "?y=2021&m=12",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
}

func TestAdminProcessReservation(t *testing.T) {
	for _, e := range adminProcessReservationTests {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/process-reservation/cal/1/do%s", e.queryParams), nil)
		ctx := getContext(req)
		req.RequestURI = fmt.Sprintf("/admin/process-reservation/cal/1/do%s", e.queryParams)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminProcessReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusSeeOther {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}
	}
}

var adminDeleteReservationTests = []struct {
	name                 string
	queryParams          string
	expectedResponseCode int
	expectedLocation     string
}{
	{
		name:                 "delete-reservation",
		queryParams:          "",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
	{
		name:                 "delete-reservation-back-to-cal",
		queryParams:          "?y=2021&m=12",
		expectedResponseCode: http.StatusSeeOther,
		expectedLocation:     "",
	},
}

func TestAdminDeleteReservation(t *testing.T) {
	for _, e := range adminDeleteReservationTests {
		req, _ := http.NewRequest("GET", fmt.Sprintf("/admin/process-reservation/cal/1/do%s", e.queryParams), nil)
		ctx := getContext(req)
		req = req.WithContext(ctx)
		req.RequestURI = fmt.Sprintf("/admin/process-reservation/cal/1/do%s", e.queryParams)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(Repo.AdminDeleteReservation)
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusSeeOther {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedResponseCode, rr.Code)
		}
	}
}

func getContext(r *http.Request) context.Context {
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
