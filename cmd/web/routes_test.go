package main

import (
	"fmt"
	"github.com/VishalTanwani/GolangWebApp/internal/config"
	"github.com/go-chi/chi/v5"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
	default:
		t.Error(fmt.Sprintf("this is not a *chi.Mux but %T", v))
	}
}
