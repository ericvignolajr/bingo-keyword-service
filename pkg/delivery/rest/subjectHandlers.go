package rest

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/sessions"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func newCreateSubjectFormHandler(viewsPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, _ := template.ParseFiles(viewsPath + "/createSubject.tmpl")
		tmpl.Execute(w, nil)
	}
}

func newCreateSubjectHandler(createSubjectUsecase usecases.CreateSubject, viewsPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		_, err := createSubjectUsecase.Exec(s, u.ID)
		if err != nil && errors.Is(err, domain.ErrSubjectNameEmpty) {
			w.Header().Set("HX-Retarget", "#create-subject-form")
			tmplData.Errors = append(tmplData.Errors, domain.ErrSubjectNameEmpty)
			errTmpl.Execute(w, tmplData)
			return
		}

		http.Redirect(w, r, "/user", http.StatusFound)
	}
}

func newEditSubjectHandler(readSubjectsUsecase usecases.ReadSubjects, viewsPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		s, err := readSubjectsUsecase.Exec(u.ID, &subjectID)
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
	}
}

func newUpdateSubjectFormHandler(readSubjectsUsecase usecases.ReadSubjects, viewsPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		subjects, err := readSubjectsUsecase.Exec(u.ID, &subjectID)
		if err != nil {
			if errors.Is(err, domain.ErrSubjectDoesNotExist) {
				w.WriteHeader(http.StatusNotFound)
				fmt.Println(err)
				http.Redirect(w, r, "/user", http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		tmpl, _ := template.ParseFiles(viewsPath + "/updateSubject.tmpl")
		tmplData := struct {
			ID     uuid.UUID
			Name   string
			Errors []error
		}{
			ID:     subjects[0].ID,
			Name:   subjects[0].Name,
			Errors: []error{},
		}
		tmpl.Execute(w, tmplData)
	}
}

func newUpdateSubjectHandler(readSubjectsUsecase usecases.ReadSubjects, updateSubjectUsecase usecases.UpdateSubject, viewsPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		err = r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		newName := r.PostFormValue("subjectName")

		subjects, err := readSubjectsUsecase.Exec(u.ID, &subjectID)
		if err != nil {
			if errors.Is(err, domain.ErrSubjectDoesNotExist) {
				fmt.Printf("subject with id %s does not exist for user %s", subjectID, u.ID)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		if subjects[0].UserID != u.ID {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			fmt.Println("Subject owner ID does not match ID of current user")
			fmt.Printf("UserID: %s\n", u.ID)
			fmt.Printf("Subject Owner ID: %s\n", subjects[0].UserID)
			return
		}

		err = updateSubjectUsecase.Exec(u.ID, subjectID, newName)
		tmpl, _ := template.ParseFiles(viewsPath + "/updateSubject.tmpl")
		tmplData := struct {
			ID     uuid.UUID
			Name   string
			Errors []error
		}{
			ID:     subjects[0].ID,
			Name:   subjects[0].Name,
			Errors: []error{},
		}
		if err != nil {
			if errors.Is(err, domain.ErrSubjectNameEmpty) {
				w.Header().Set("HX-Retarget", "#update-subject-form")
				tmplData.Errors = append(tmplData.Errors, domain.ErrSubjectNameEmpty)
				tmpl.Execute(w, tmplData)
				return
			}

			if errors.Is(err, domain.ErrSubjectIsDuplicate) {
				w.Header().Set("HX-Retarget", "#update-subject-form")
				tmplData.Errors = append(tmplData.Errors, domain.ErrSubjectIsDuplicate)
				tmpl.Execute(w, tmplData)
				return
			}

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Bad Request"))
			fmt.Println(fmt.Errorf("in %s with http method %s: %w", r.URL.Path, r.Method, err))
			return
		}
		http.Redirect(w, r, "/user", http.StatusFound)
	}
}

func newDeleteSubjectHandler(deleteSubjectUsecase usecases.DeleteSubject) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		err = deleteSubjectUsecase.Exec(u.ID, subjectID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
