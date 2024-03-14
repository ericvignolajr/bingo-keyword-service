package usecases

import (
	"fmt"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
	"github.com/google/uuid"
)

type DeleteSubject struct {
	SubjectStore stores.Subject
	UserStore    stores.User
}

func (d *DeleteSubject) Exec(userID uuid.UUID, subjectID uuid.UUID) error {
	user, err := d.UserStore.ReadById(userID)
	if err != nil {
		return fmt.Errorf("in deleteSubject: %w", err)
	}

	err = user.DeleteSubject(subjectID)
	if err != nil {
		return fmt.Errorf("in deleteSubject: %w", err)
	}

	d.SubjectStore.Delete(userID, subjectID)
	// d.UserStore.Save(user)
	return nil
}
