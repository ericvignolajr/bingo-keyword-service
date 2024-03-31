package usecases

import (
	"fmt"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
	"github.com/google/uuid"
)

type UpdateSubject struct {
	UserStore stores.User
}

func (c *UpdateSubject) Exec(UserID uuid.UUID, SubjectID uuid.UUID, NewName string) error {
	user, err := c.UserStore.ReadById(UserID)
	if err != nil {
		return fmt.Errorf("error in UpdateSubject: %w", err)
	}

	subject, err := user.FindSubject(SubjectID)
	if err != nil {
		return fmt.Errorf("error in UpdateSubject: %w", err)
	}

	isDuplicate := user.IsDuplicateSubject(domain.Subject{ID: subject.ID, Name: NewName})
	if isDuplicate {
		return fmt.Errorf("error in UpdateSubject: %w", domain.ErrSubjectIsDuplicate)
	}

	err = subject.UpdateSubjectName(NewName)
	if err != nil {
		return fmt.Errorf("error in UpdateSubject: %w", err)
	}

	_, err = c.UserStore.Save(user)
	if err != nil {
		return fmt.Errorf("error in UpdateSubject: %w", err)
	}

	return nil
}
