package usecases

import (
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
	"github.com/google/uuid"
)

type DeleteSubject struct {
	SubjectStore stores.Subject
}

func (d *DeleteSubject) Exec(userID uuid.UUID, subjectID uuid.UUID) error {
	d.SubjectStore.Delete(userID, subjectID)
	return nil
}
