package usecases

import (
	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
	"github.com/google/uuid"
)

type ReadKeywordsBySubjectID struct {
	UserStore stores.User
}

func (r *ReadKeywordsBySubjectID) Exec(userID uuid.UUID, subjectID uuid.UUID) ([]*domain.Keyword, error) {
	user, err := r.UserStore.ReadById(userID)
	if err != nil {
		return nil, err
	}

	subj, err := user.FindSubject(subjectID)
	if err != nil {
		return nil, err
	}

	keywords := make([]*domain.Keyword, 0)
	for _, u := range subj.Units {
		keywords = append(keywords, u.Keywords...)
	}

	return keywords, nil
}
