package inmemory

import (
	"fmt"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserStore struct {
	store []domain.User
}

func (u *UserStore) ReadById(id uuid.UUID) (*domain.User, error) {
	for _, v := range u.store {
		if v.Id == id {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("user %s could not be found", id)
}
func (u *UserStore) ReadByEmail(e string) (*domain.User, error) {
	for _, v := range u.store {
		if v.Email == e {
			return &v, nil
		}
	}

	return nil, nil
}

func (u *UserStore) Create(email string, password string) (uuid.UUID, error) {
	hashedP, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user, err := domain.NewUser(email, string(hashedP))
	if err != nil {
		return uuid.Nil, err
	}
	u.store = append(u.store, *user)

	return user.Id, nil
}

func (u *UserStore) CreateAccount(email string) error {
	user, err := domain.NewUser(email, "")
	if err != nil {
		return err
	}

	u.store = append(u.store, *user)
	return nil
}
