package inmemory

import (
	"fmt"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserStore struct {
	store []domain.User
}

func (u *UserStore) ReadById(id uuid.UUID) (*domain.User, error) {
	for i, v := range u.store {
		if v.Id == id {
			return &u.store[i], nil
		}
	}

	return nil, fmt.Errorf("%s, userID: %s", stores.ErrUserDoesNotExist, id)
}
func (u *UserStore) ReadByEmail(e string) (*domain.User, error) {
	for i, v := range u.store {
		if v.Email == e {
			return &u.store[i], nil
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

func (u *UserStore) Save(user *domain.User) (*domain.User, error) {
	userToUpdate, err := u.ReadById(user.Id)
	if err != nil {
		return nil, err
	}

	userToUpdate = user
	userToUpdate.Subjects = user.Subjects
	return userToUpdate, nil
}
