package delivery

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/viewers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jws"
)

type arbitraryJSON struct {
	data interface{}
}

func NewHttpServer(rs usecases.ReadSubject, rv *viewers.WebViewer) chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.SetHeader("X-Frame-Options", "DENY"))

	r.Route("/signin", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			wd, err := os.Getwd()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			path := filepath.Join(wd, "/pkg", "/delivery", "/signin.html")

			f, err := os.Open(path)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			io.Copy(w, f)
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			err := r.ParseForm()
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
			str, ok := r.PostForm["credential"]
			if !ok {
				w.WriteHeader(http.StatusBadRequest)
			}
			src := str[0]
			token, _ := jws.Parse([]byte(src))
			w.Header().Set("Content-Type", "application/json")
			w.Write(token.Payload())
		})
	})

	r.Route("/subjects", func(r chi.Router) {
		r.Put("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotImplemented)
			return
		})
		r.Get("/{subjectID}", func(w http.ResponseWriter, r *http.Request) {
			s := chi.URLParam(r, "subjectID")
			uid, err := uuid.Parse(s)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error() + "\n"))
				w.Write([]byte(s))
				return
			}
			req := usecases.ReadSubjectRequest{
				UserId:      uid,
				SubjectName: "Science",
			}
			res := rs.Exec(req)
			if res.Err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("<h1>There was an error</h1>"))
			}

			w.Header().Add("Content-Type", "text/html")
			w.Write([]byte(rv.View()))
		})
	})

	return r
}
