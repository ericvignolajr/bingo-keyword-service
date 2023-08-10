package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func subjectRoutes() chi.Router {
	r := chi.NewRouter()

	r.Get("/{subjectID}", readSubject)
	r.Put("/", createSubject)

	return r
}

func createSubject(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}

func readSubject(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
	return
}
