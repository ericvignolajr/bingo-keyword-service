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
	db.Migrator().AutoMigrate(domain.User{}, domain.Subject{})
	return &SQLUserStore{DB: db}, nil
}

func (s *SQLUserStore) ReadById(userID uuid.UUID) (*domain.User, error) {
	user := domain.User{
		ID: userID,
	}
	result := s.DB.Model(&user).First(&user)
	result.Row().Scan(user)
	return nil, nil
}
