package usecases

import (
	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
	"github.com/google/uuid"
)

type ReadKeywordsByUnitID struct {
	UserStore stores.User
}

func (r *ReadKeywordsByUnitID) Exec(userID uuid.UUID, unitID uuid.UUID) ([]*domain.Keyword, error) {
	user, err := r.UserStore.ReadById(userID)
	if err != nil {
		return nil, err
	}

	var unit domain.Unit
	for _, s := range user.Subjects {
		for _, u := range s.Units {
			if u.ID == unitID {
				unit = *u
				break
			}
		}
	}

	keywords := make([]*domain.Keyword, len(unit.Keywords))
	copy(keywords, unit.Keywords)
	return keywords, nil
}
