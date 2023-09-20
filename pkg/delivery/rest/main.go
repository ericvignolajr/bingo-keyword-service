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
	outputports "github.com/ericvignolajr/bingo-keyword-service/pkg/output_ports"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/sessions"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/inmemory"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/google/uuid"

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
	subjectStore := inmemory.NewSubjectStore()

	deleteSubject := usecases.DeleteSubject{
		SubjectStore: subjectStore,
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
			subjects := usecases.ReadSubjects(req, subjectStore)
			if subjects.Err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			data := struct {
				User     domain.User
				Subjects []usecases.SubjectOutput
			}{
				User:     *u,
				Subjects: subjects.Subjects,
			}

			res, _ := template.ParseFiles(viewsPath+"/base.tmpl", viewsPath+"/user.tmpl")
			if err != nil {
				fmt.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			res.Execute(w, data)
		})
	})

	r.Route("/subject", func(r chi.Router) {
		r.Use(requireSession)
		r.Use(sessions.AddUserToContext)
		r.Get("/create", func(w http.ResponseWriter, r *http.Request) {
			tmpl, _ := template.ParseFiles(viewsPath + "/createSubject.tmpl")
			tmpl.Execute(w, nil)
		})
		r.Post("/create", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			u, ok := ctx.Value(sessions.User).(*domain.User)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				fmt.Println("No user on context")
				return
			}

			r.ParseForm()
			s := r.PostFormValue("subjectName")
			if s == "" {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Println("subjectName is required to create a new subject")
				return
			}
			subject, _ := domain.NewSubject(s)

			createSubject := usecases.CreateSubject{
				UserStore:    &sessions.UserStore,
				SubjectStore: subjectStore,
				Presenter:    &outputports.MockPresenter{},
			}

			req := usecases.CreateSubjectRequest{
				UserId:      u.Id,
				SubjectName: s,
			}
			res := createSubject.Exec(req)
			data := new(struct {
				ID   uuid.UUID
				Name string
			})
			if res.Err != nil && res.Err.Error() == "subject already exists" {
				data = nil
			} else {
				data.ID = subject.Id
				data.Name = subject.Name
			}

			tmpl, _ := template.ParseFiles(viewsPath + "/user.tmpl")
			tmpl.ExecuteTemplate(w, "subject-element", data)
		})
		r.Get("/{subjectID}/edit", func(w http.ResponseWriter, r *http.Request) {
			subjectID, err := uuid.Parse(chi.URLParam(r, "subjectID"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("Bad Request"))
				fmt.Println("could not parse subject ID into UUID")
				return
			}

			ctx := r.Context()
			u, ok := ctx.Value(sessions.User).(*domain.User)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				fmt.Println("No user on context")
				return
			}

			readSubject := usecases.ReadSubjectByID{
				SubjectStore: subjectStore,
			}
			s, err := readSubject.ReadSubjectByID(subjectID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("We've hit an issue, please retry your request"))
				fmt.Println(err)
				return
			}

			if s.OwnerID != u.Id {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				fmt.Println("Subject owner ID does not match ID of current user")
				return
			}

			tmpl, _ := template.ParseFiles(viewsPath + "/editSubject.tmpl")
			tmpl.Execute(w, s)
		})
		r.Put("/{subjectID}", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(chi.URLParam(r, "subjectID"))
		})
		r.Delete("/{subjectID}", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			u, ok := ctx.Value(sessions.User).(*domain.User)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				fmt.Println("No user on context")
				return
			}

			subjectID, err := uuid.Parse(chi.URLParam(r, "subjectID"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
			}
			err = deleteSubject.Exec(u.Id, subjectID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		})
	})

	return r
}
