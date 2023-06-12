package stores

import (
	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/google/uuid"
)

type User interface {
	ReadById(uuid.UUID) (*domain.User, error)
}
