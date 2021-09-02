package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
	// "net/url"
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
	// {"search-availability", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2020-02-02"},
	// 	{key: "end", value: "2020-02-02"},
	// }, http.StatusOK},
	// {"search-availability-json", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "2020-02-02"},
	// 	{key: "end", value: "2020-02-02"},
	// }, http.StatusOK},
	// {"make-reservation", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "vishal"},
	// 	{key: "last_name", value: "tanwani"},
	// 	{key: "email", value: "vishal@gamil.com"},
	// 	{key: "phone", value: "444-4444-44"},
	// }, http.StatusOK},
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

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong code : wanted %d this get this %d", http.StatusTemporaryRedirect, rr.Code)
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

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong code : wanted %d this get this %d", http.StatusTemporaryRedirect, rr.Code)
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

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong code : wanted %d this get this %d", http.StatusTemporaryRedirect, rr.Code)
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

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post Reservation handler returned wrong code : wanted %d this get this %d", http.StatusTemporaryRedirect, rr.Code)
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

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post Reservation handler returned wrong code : wanted %d this get this %d", http.StatusTemporaryRedirect, rr.Code)
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

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post Reservation handler returned wrong code : wanted %d this get this %d", http.StatusTemporaryRedirect, rr.Code)
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

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong code : wanted %d this get this %d", http.StatusTemporaryRedirect, rr.Code)
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

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("post Reservation handler returned wrong code : wanted %d this get this %d", http.StatusTemporaryRedirect, rr.Code)
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

func getContext(r *http.Request) context.Context {
	ctx, err := session.Load(r.Context(), r.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
