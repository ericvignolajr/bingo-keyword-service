package inmemory

import (
	"fmt"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/google/uuid"
)

type SubjectStore struct {
	Store map[uuid.UUID][]domain.Subject
}

func NewSubjectStore() *SubjectStore {
	store := make(map[uuid.UUID][]domain.Subject)
	return &SubjectStore{
		Store: store,
	}
}

func (s *SubjectStore) ReadByName(UserId uuid.UUID, SubjectName string) (*domain.Subject, error) {
	for _, v := range s.Store[UserId] {
		if v.Name == SubjectName {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("subject %s for user %s could not be found", SubjectName, UserId)
}

func (s *SubjectStore) Create(UserId uuid.UUID, Subject *domain.Subject) (*domain.Subject, error) {
	subjectToCreate, _ := s.ReadByName(UserId, Subject.Name)
	if subjectToCreate != nil {
		return subjectToCreate, nil
	}

	subjects, ok := s.Store[UserId]
	if !ok {
		s.Store[UserId] = []domain.Subject{*Subject}
	} else {
		subjects = append(subjects, *Subject)
	}

	return Subject, nil
}
