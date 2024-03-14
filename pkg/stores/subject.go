package stores

import (
	"github.com/google/uuid"
)

const (
	ErrSubjectExists = "subject already exists"
)

type Subject interface {
	Delete(userID, subjectID uuid.UUID) error
}
