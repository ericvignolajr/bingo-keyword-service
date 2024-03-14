package rest

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/sessions"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/sql"
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
		r.Get("/user", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			u, ok := ctx.Value(sessions.User).(*domain.User)
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				fmt.Println("No user on context")
				return
			}

			subjects, err := readSubjects.Exec(u.ID, nil)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			data := struct {
				User     domain.User
				Subjects []domain.Subject
			}{
				User:     *u,
				Subjects: subjects,
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

			errTmpl, _ := template.ParseFiles(viewsPath + "/createSubject.tmpl")
			tmplData := new(struct {
				ID     uuid.UUID
				Name   string
				Errors []error
			})

			_, err := createSubject.Exec(s, u.ID)
			if err != nil && errors.Is(err, domain.ErrSubjectNameEmpty) {
				w.Header().Set("HX-Retarget", "#create-subject-form")
				tmplData.Errors = append(tmplData.Errors, domain.ErrSubjectNameEmpty)
				errTmpl.Execute(w, tmplData)
				return
			}

			http.Redirect(w, r, "/user", http.StatusFound)
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

			s, err := readSubjects.Exec(u.ID, &subjectID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("We've hit an issue, please retry your request"))
				fmt.Println(err)
				return
			}

			if len(s) > 0 && s[0].UserID != u.ID {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				fmt.Println("Subject owner ID does not match ID of current user")
				fmt.Printf("UserID: %s\n", u.ID)
				fmt.Printf("Subject Owner ID: %s\n", s[0].UserID)
				return
			}

			tmplData := struct {
				ID      uuid.UUID
				Subject string
				Units   []*domain.Unit
			}{
				ID:      subjectID,
				Subject: s[0].Name,
				Units:   s[0].Units,
			}
			tmpl, _ := template.ParseFiles(viewsPath + "/subject.tmpl")
			tmpl.Execute(w, tmplData)
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
			err = deleteSubject.Exec(u.ID, subjectID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		})

		r.Route("/{subjectID}/unit", func(r chi.Router) {
			r.Use(requireSession)
			r.Use(sessions.AddUserToContext)
			r.Get("/create", func(w http.ResponseWriter, r *http.Request) {
				sID, err := uuid.Parse(chi.URLParam(r, "subjectID"))
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("Bad Request"))
					fmt.Println("could not parse subject ID into UUID")
					return
				}
				tmplData := struct {
					SubjectID uuid.UUID
					Errors    []error
				}{
					SubjectID: sID,
					Errors:    nil,
				}
				fmt.Println(tmplData.SubjectID)
				tmpl, err := template.ParseFiles(viewsPath + "/createUnit.tmpl")
				if err != nil {
					fmt.Println(err)
				}
				tmpl.Execute(w, tmplData)
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

				sID, err := uuid.Parse(chi.URLParam(r, "subjectID"))
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("Bad Request"))
					fmt.Printf("could not parse subject ID \"%s\" into UUID", sID)
					return
				}

				err = r.ParseForm()
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				unitName := r.PostFormValue("unitName")
				_, err = createUnit.Exec(unitName, u.ID, sID)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Println(err)
					return
				}

				subject, err := readSubject.ReadSubjectByID(u.ID, sID)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Println(err)
					return
				}

				if subject.UserID != u.ID {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Unauthorized"))
					fmt.Println("Subject owner ID does not match ID of current user")
					fmt.Printf("UserID: %s\n", u.ID)
					fmt.Printf("Subject Owner ID: %s\n", subject.UserID)
					return
				}

				tmplData := struct {
					Subject string
					ID      uuid.UUID
					Units   []*domain.Unit
				}{
					Subject: subject.Name,
					ID:      subject.ID,
					Units:   subject.Units,
				}

				fmt.Println(tmplData)
				for _, v := range tmplData.Units {
					fmt.Println(v.Name)
				}

				http.Redirect(w, r, fmt.Sprintf("/subject/%s/edit", sID), http.StatusFound)
			})
		})

	})

	return r
}
