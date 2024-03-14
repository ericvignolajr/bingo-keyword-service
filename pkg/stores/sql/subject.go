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

func NewSQLSubjectStore() *SubjectStore {
	return &SubjectStore{DB: db}
}

// func (s *SubjectStore) Read(UserID uuid.UUID) ([]*domain.Subject, error) {
// 	subjects, ok := s.Store[UserID]
// 	if !ok {
// 		return nil, nil
// 	}

// 	subjectPointers := make([]*domain.Subject, len(subjects))
// 	for i := range subjects {
// 		subjectPointers[i] = subjects[i]
// 	}

// 	return subjectPointers, nil
// }

// func (s *SubjectStore) ReadByName(UserId uuid.UUID, SubjectName string) (*domain.Subject, error) {
// 	for _, v := range s.Store[UserId] {
// 		if strings.ToLower(v.Name) == strings.ToLower(SubjectName) {
// 			return v, nil
// 		}
// 	}

// 	return nil, fmt.Errorf("subject %s for user %s could not be found", SubjectName, UserId)
// }

// func (s *SubjectStore) ReadByID(subjectID uuid.UUID) (*domain.Subject, error) {
// 	for _, user := range s.Store {
// 		for _, subject := range user {
// 			if subject.ID == subjectID {
// 				return subject, nil
// 			}
// 		}
// 	}

// 	return nil, nil
// }

// func (s *SubjectStore) Create(UserId uuid.UUID, Subject *domain.Subject) (*domain.Subject, error) {
// 	subjectToCreate, _ := s.ReadByName(UserId, Subject.Name)
// 	if subjectToCreate != nil {
// 		return subjectToCreate, errors.New(stores.ErrSubjectExists)
// 	}

// 	Subject.UserID = UserId
// 	subjects, ok := s.Store[UserId]
// 	if !ok {
// 		s.Store[UserId] = []*domain.Subject{Subject}
// 	} else {
// 		s.Store[UserId] = append(subjects, Subject)
// 	}

// 	return Subject, nil
// }

// func (s *SubjectStore) Update(Subject *domain.Subject) (*domain.Subject, error) {
// 	subjectToUpdate, _ := s.ReadByID(Subject.ID)
// 	if subjectToUpdate == nil {
// 		newSubject, err := s.Create(Subject.UserID, Subject)
// 		if err != nil {
// 			return nil, err
// 		}
// 		return newSubject, nil
// 	}

// 	subjects, ok := s.Store[Subject.UserID]
// 	if !ok {
// 		return nil, errors.New(fmt.Sprintf("could not find record in in-memory subject store for user with id: %s", Subject.UserID.String()))
// 	}
// 	for i, v := range subjects {
// 		if v.ID == Subject.ID {
// 			subjects[i] = Subject
// 		}
// 	}

// 	return Subject, nil
// }

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
