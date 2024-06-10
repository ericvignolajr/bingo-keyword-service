package usecases

import (
	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
	"github.com/google/uuid"
)

type ReadKeywordsByUserID struct {
	UserStore stores.User
}

func (r *ReadKeywordsByUserID) Exec(userID uuid.UUID) ([]*domain.Keyword, error) {
	user, err := r.UserStore.ReadById(userID)
	if err != nil {
		return nil, err
	}

	keywords := make([]*domain.Keyword, 0)
	for _, s := range user.Subjects {
		for _, u := range s.Units {
			keywords = append(keywords, u.Keywords...)
		}
	}

	return keywords, nil
}
