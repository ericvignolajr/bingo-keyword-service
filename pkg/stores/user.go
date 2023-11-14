package stores

import (
	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/google/uuid"
)

const (
	ErrUserDoesNotExist = "could not find user in database"
)

type User interface {
	ReadById(uuid.UUID) (*domain.User, error)
	ReadByEmail(string) (*domain.User, error)
	CreateAccount(string) error
	Save(User *domain.User) (*domain.User, error)
}
