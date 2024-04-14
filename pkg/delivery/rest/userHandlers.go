package rest

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/sessions"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
)

func newShowHomePageHandler(readSubjectsUsecase usecases.ReadSubjects, viewsPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u, ok := ctx.Value(sessions.User).(*domain.User)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			fmt.Println("No user on context")
			return
		}

		subjects, err := readSubjectsUsecase.Exec(u.ID, nil)
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

		res, err := template.ParseFiles(viewsPath+"/base.tmpl", viewsPath+"/user.tmpl")
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res.Execute(w, data)
	}
}
