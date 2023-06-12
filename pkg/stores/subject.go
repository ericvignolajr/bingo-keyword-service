package stores

import (
	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/google/uuid"
)

type Subject interface {
	ReadByName(UserId uuid.UUID, SubjectName string) (*domain.Subject, error)
	Create(UserId uuid.UUID, Subject *domain.Subject) (*domain.Subject, error)
}
