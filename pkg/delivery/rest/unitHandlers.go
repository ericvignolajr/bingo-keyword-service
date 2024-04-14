package rest

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/sessions"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func newCreateUnitFormHandler(viewsPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func newCreateUnitHandler(createUnitUsecase usecases.CreateUnit, readSubjectUsecase usecases.ReadSubjectByID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		_, err = createUnitUsecase.Exec(unitName, u.ID, sID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
			return
		}

		subject, err := readSubjectUsecase.ReadSubjectByID(u.ID, sID)
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
	}
}
