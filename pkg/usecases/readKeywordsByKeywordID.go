package usecases

import (
	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
	"github.com/google/uuid"
)

type ReadKeywordsByKeywordID struct {
	UserStore stores.User
}

func (r *ReadKeywordsByKeywordID) Exec(userID uuid.UUID, keywordID uuid.UUID) ([]*domain.Keyword, error) {
	user, err := r.UserStore.ReadById(userID)
	if err != nil {
		return nil, err
	}

	keywords := make([]*domain.Keyword, 1)
	for _, s := range user.Subjects {
		for _, u := range s.Units {
			for _, k := range u.Keywords {
				if k.ID == keywordID {
					keywords[0] = k
					break
				}
			}
		}
	}

	return keywords, nil
}
