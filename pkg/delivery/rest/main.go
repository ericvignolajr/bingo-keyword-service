package rest

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"

	"github.com/clerkinc/clerk-sdk-go/clerk"
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

	r.Get("/signin", signin)

	var clerkAuth = sessions.NewClerkAuthenticator()
	requireSession := clerk.RequireSessionV2(clerkAuth.Client)
	r.Group(func(r chi.Router) {
		r.Use(requireSession)
		r.Get("/dashboard", dashboard(clerkAuth.Client))
	})

	return r
}

func signin(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	path := filepath.Join(wd, "/pkg", "/delivery", "/signin.html")

	http.ServeFile(w, r, path)
}

func dashboard(clerkAuth clerk.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		c, _ := clerk.SessionFromContext(ctx)
		u, _ := clerkAuth.Users().Read(c.Subject)
		json.NewEncoder(w).Encode(u)
	}
}
