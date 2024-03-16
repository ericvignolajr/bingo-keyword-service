package usecases

import (
	"fmt"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
	"github.com/google/uuid"
)

type CreateUnit struct {
	SubjectStore stores.Subject
	UserStore    stores.User
}

func (c *CreateUnit) Exec(name string, userID uuid.UUID, subjectID uuid.UUID) (*domain.Unit, error) {
	user, err := c.UserStore.ReadById(userID)
	if err != nil {
		return nil, fmt.Errorf("in createUnit: %w", err)
	}

	u, err := domain.NewUnit(name)
	if err != nil {
		return nil, err
	}

	s, err := user.FindSubject(subjectID)
	if err != nil {
		return nil, fmt.Errorf("in createUnit: %w", err)
	}

	err = s.AddUnit(u)
	if err != nil {
		return nil, fmt.Errorf("in createUnit: %w", err)
	}

	_, err = c.UserStore.Save(user)
	if err != nil {
		return nil, err
	}

	return u, nil
}
