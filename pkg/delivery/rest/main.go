package rest

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/sessions"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/inmemory"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"

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
	var clerkAuth = sessions.NewClerkAuthenticator()
	requireSession := clerk.RequireSessionV2(clerkAuth.Client, clerk.WithLeeway(30*time.Second))
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	viewsPath := filepath.Join(wd, "/pkg/delivery/rest/views")
	subjectStore := inmemory.SubjectStore{}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/user", http.StatusFound)
	})
	r.Get("/signin", func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(viewsPath + "/signin.html")

		http.ServeFile(w, r, path)
	})

	r.Group(func(r chi.Router) {
		r.Use(requireSession)
		r.Use(sessions.AddUserToContext)
		r.Get("/user", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			u, ok := ctx.Value(sessions.User).(*domain.User)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				fmt.Println("No user on context")
				return
			}

			req := usecases.ReadSubjectsRequest{
				UserID: u.Id,
			}
			data := usecases.ReadSubjects(req, &subjectStore)
			if data.Err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			res, _ := template.ParseFiles(viewsPath+"/base.tmpl", viewsPath+"/user.tmpl")
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			res.Execute(w, u)
		})
	})

	return r
}
