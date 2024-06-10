package rest

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/delivery/rest/internal/templates"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/sessions"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/google/uuid"
)

func newReadKeywordsHandler(
	readKeywordsByUserID usecases.ReadKeywordsByUserID,
	readKeywordsBySubjectID usecases.ReadKeywordsBySubjectID,
	readKeywordsByUnitID usecases.ReadKeywordsByUnitID,
	readKeywordsByKeywordID usecases.ReadKeywordsByKeywordID,
	templates *templates.Templates) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		u, ok := ctx.Value(sessions.User).(*domain.User)
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			fmt.Println("No user on context")
			return
		}

		query, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, "no keywords match the query")
			message := fmt.Sprintf("could not parse query, %s", r.URL.RawQuery)
			fmt.Println(message, err)
			return
		}

		keywordID := query.Get(keywordIDQueryString)
		if keywordID != "" {
			keywordIDAsUUID, err := uuid.Parse(keywordID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Printf("could not parse keyword ID: %s into UUID\n", keywordID)
				return
			}
			keywords, err := readKeywordsByKeywordID.Exec(u.ID, keywordIDAsUUID)
			if err != nil {
				fmt.Println(err)
			}
			templates.Render(w, "keyword.tmpl", keywords)
			return
		}

		unitID := query.Get(unitIDQueryString)
		if unitID != "" {
			unitIDAsUUID, err := uuid.Parse(unitID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Printf("could not parse unit ID: %s into UUID\n", unitID)
				return
			}
			keywords, _ := readKeywordsByUnitID.Exec(u.ID, unitIDAsUUID)
			templates.Render(w, "keyword.tmpl", keywords)
			return
		}

		subjectID := query.Get(subjectIDQueryString)
		if subjectID != "" {
			subjectIDAsUUID, err := uuid.Parse(subjectID)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Printf("could not parse subject ID: %s into UUID\n", subjectID)
				return
			}
			keywords, _ := readKeywordsBySubjectID.Exec(u.ID, subjectIDAsUUID)
			templates.Render(w, "keyword.tmpl", keywords)
			return
		}

		keywords, _ := readKeywordsByUserID.Exec(u.ID)
		templates.Render(w, "keyword.tmpl", keywords)
	}
}
