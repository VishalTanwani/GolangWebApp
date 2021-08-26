package main

import (
	"github.com/VishalTanwani/GolangWebApp/internal/config"
	"github.com/VishalTanwani/GolangWebApp/internal/handler"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(Nosurf)
	mux.Use(SessionLoad)

	mux.Get("/", handler.Repo.Home)
	mux.Get("/about", handler.Repo.About)
	mux.Get("/general-quarters", handler.Repo.Generals)
	mux.Get("/major-suite", handler.Repo.Majors)

	mux.Get("/search-availability", handler.Repo.Availability)
	mux.Post("/search-availability", handler.Repo.PostAvailability)
	mux.Post("/search-availability-json", handler.Repo.AvailabilityJSON)

	mux.Get("/make-reservation", handler.Repo.MakeReservation)
	mux.Post("/make-reservation", handler.Repo.PostReservation)
	mux.Get("/contact", handler.Repo.Contact)
	mux.Get("/reservation-summary", handler.Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
