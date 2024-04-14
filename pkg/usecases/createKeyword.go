package usecases

import (
	"fmt"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
	"github.com/google/uuid"
)

type CreateKeyword struct {
	UserStore stores.User
}

func (c *CreateKeyword) Exec(name string, userID uuid.UUID, subjectID uuid.UUID, unitID uuid.UUID) (*domain.Keyword, error) {
	user, err := c.UserStore.ReadById(userID)
	if err != nil {
		return nil, fmt.Errorf("in createKeyword: %w", err)
	}

	subject, err := user.FindSubject(subjectID)
	if err != nil {
		return nil, fmt.Errorf("in createKeyword: %w", err)
	}

	unit, err := subject.FindUnitByID(unitID)
	if err != nil {
		return nil, fmt.Errorf("in createKeyword: %w", err)
	}

	newKeyword, err := domain.NewKeyword(name, "")
	if err != nil {
		return nil, fmt.Errorf("in createKeyword: %w", err)
	}

	err = unit.AddKeyword(newKeyword)
	if err != nil {
		return nil, fmt.Errorf("in createKeyword: %w", err)
	}

	_, err = c.UserStore.Save(user)
	if err != nil {
		return nil, fmt.Errorf("in createKeyword: %w", err)
	}

	return newKeyword, nil
}
