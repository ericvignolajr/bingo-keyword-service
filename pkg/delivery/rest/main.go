package rest

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/sessions"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/sql"
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
	subjectStore, _ := sql.NewSQLSubjectStore()

	createSubject := usecases.CreateSubject{
		UserStore:    sessions.UserStore,
		SubjectStore: subjectStore,
	}

	readSubject := usecases.ReadSubjectByID{
		UserStore:    sessions.UserStore,
		SubjectStore: subjectStore,
	}

	readSubjects := usecases.ReadSubjects{
		UserStore: sessions.UserStore,
	}

	updateSubject := usecases.UpdateSubject{
		UserStore: sessions.UserStore,
	}

	deleteSubject := usecases.DeleteSubject{
		SubjectStore: subjectStore,
		UserStore:    sessions.UserStore,
	}

	createUnit := usecases.CreateUnit{
		SubjectStore: subjectStore,
		UserStore:    sessions.UserStore,
	}

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
		r.Get("/user", newShowHomePageHandler(readSubjects, viewsPath))
	})

	r.Route("/subject", func(r chi.Router) {
		r.Use(requireSession)
		r.Use(sessions.AddUserToContext)
		r.Get("/create", newCreateSubjectFormHandler(viewsPath))
		r.Post("/create", newCreateSubjectHandler(createSubject, viewsPath))
		r.Get("/{subjectID}/edit", newEditSubjectHandler(readSubjects, viewsPath))
		r.Get("/{subjectID}/update", newUpdateSubjectFormHandler(readSubjects, viewsPath))
		r.Post("/{subjectID}/update", newUpdateSubjectHandler(readSubjects, updateSubject, viewsPath))
		r.Delete("/{subjectID}", newDeleteSubjectHandler(deleteSubject))

		r.Route("/{subjectID}/unit", func(r chi.Router) {
			r.Use(requireSession)
			r.Use(sessions.AddUserToContext)
			r.Get("/create", newCreateUnitFormHandler(viewsPath))
			r.Post("/create", newCreateUnitHandler(createUnit, readSubject))
		})

	})

	return r
}
