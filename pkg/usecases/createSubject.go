package usecases

import (
	"fmt"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
	"github.com/google/uuid"
)

type CreateSubject struct {
	UserStore    stores.User
	SubjectStore stores.Subject
}

func (c *CreateSubject) Exec(name string, userID uuid.UUID) (*domain.Subject, error) {
	user, err := c.UserStore.ReadById(userID)
	if err != nil {
		return nil, fmt.Errorf("in createSubject: %w", err)
	}

	subject, err := domain.NewSubject(name, user.ID)
	if err != nil {
		return nil, fmt.Errorf("in createSubject: %w", err)
	}

	_, err = user.AddSubject(*subject)
	if err != nil {
		return nil, fmt.Errorf("in createSubject: %w", err)
	}

	_, err = c.UserStore.Save(user)
	if err != nil {
		return nil, fmt.Errorf("in createSubject: %w", err)
	}

	return subject, nil
}
