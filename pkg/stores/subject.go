package stores

import (
	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/google/uuid"
)

type Subject interface {
	Read(UserID uuid.UUID) ([]*domain.Subject, error)
	ReadByName(UserId uuid.UUID, SubjectName string) (*domain.Subject, error)
	Create(UserId uuid.UUID, Subject *domain.Subject) (*domain.Subject, error)
}
