package sql

import (
	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLUserStore struct {
	DB *gorm.DB
}

func NewSQLUserStore() (*SQLUserStore, error) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	db.Migrator().AutoMigrate(domain.User{}, domain.Subject{})
	return &SQLUserStore{DB: db}, nil
}

func (s *SQLUserStore) ReadById(userID uuid.UUID) (*domain.User, error) {
	user := domain.User{
		ID: userID,
	}
	err := s.DB.Preload("Subjects").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
