package rest

import (
	"github.com/ericvignolajr/bingo-keyword-service/pkg/sessions"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewServer() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("X-Frame-Options", "DENY"))

	r.Mount("/signin", userRoutes(&sessions.GoogleAuthenticator{}))
	r.Mount("/subjects", subjectRoutes())

	return r
}
