package sql

import (
	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SQLUserStore struct {
	DB *gorm.DB
}

func NewSQLUserStore() (*SQLUserStore, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	db.Migrator().AutoMigrate(
		domain.User{},
		domain.Subject{},
		domain.Unit{},
		domain.Translation{},
	)
	return &SQLUserStore{DB: db}, nil
}

func (s *SQLUserStore) ReadById(userID uuid.UUID) (*domain.User, error) {
	user := domain.User{
		ID: userID,
	}
	err := s.DB.Preload("Subjects").Preload("Subjects.Units").First(&user, user.ID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *SQLUserStore) ReadByEmail(email string) (*domain.User, error) {
	user := domain.User{
		Email: email,
	}

	err := s.DB.Where("email = ?", user.Email).Preload("Subjects").First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}
