package sql

import (
	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SQLUserStore struct {
	DB *gorm.DB
}

func NewSQLUserStore() (*SQLUserStore, error) {
	return &SQLUserStore{DB: db}, nil
}

func (s *SQLUserStore) ReadById(userID uuid.UUID) (*domain.User, error) {
	user := domain.User{
		ID: userID,
	}
	err := s.DB.Preload("Subjects.Units.Keywords").First(&user, user.ID).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &stores.RecordNotFoundError{Err: err}
		}
		return nil, err
	}
	return &user, nil
}

func (s *SQLUserStore) ReadByEmail(email string) (*domain.User, error) {
	user := domain.User{
		Email: email,
	}

	err := s.DB.Where("email = ?", user.Email).Preload("Subjects.Units.Keywords").First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &stores.RecordNotFoundError{Err: err}
		}
		return nil, err
	}

	return &user, nil
}

func (s *SQLUserStore) Save(User *domain.User) (*domain.User, error) {
	if err := s.DB.Save(User).Error; err != nil {
		return nil, err
	}

	return User, nil
}

func (s *SQLUserStore) Create(email string, password string) (uuid.UUID, error) {
	hashedP, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user, err := domain.NewUser(email, string(hashedP))
	if err != nil {
		return uuid.Nil, err
	}

	_, err = s.Save(user)
	if err != nil {
		return uuid.Nil, err
	}
	return user.ID, nil
}

func (s *SQLUserStore) CreateAccount(email string) error {
	_, err := s.Create(email, "")
	if err != nil {
		return err
	}
	return nil
}
