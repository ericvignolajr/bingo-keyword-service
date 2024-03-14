package sql

import (
	"errors"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubjectStore struct {
	DB *gorm.DB
}

func NewSQLSubjectStore() (*SubjectStore, error) {
	return &SubjectStore{DB: db}, nil
}

func (s *SubjectStore) Delete(userID uuid.UUID, subjectID uuid.UUID) error {
	subject := domain.Subject{ID: subjectID}
	err := s.DB.Transaction(func(tx *gorm.DB) error {
		var units []domain.Unit
		err := tx.Where("subject_id = ?", subjectID).Find(&units).Error
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		for _, unit := range units {
			err := tx.Where("unit_id = ?", unit.ID).Delete(domain.Keyword{}).Error
			if err != nil {
				return err
			}

			err = tx.Delete(&unit).Error
			if err != nil {
				return err
			}
		}

		err = tx.Delete(&subject).Error
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
