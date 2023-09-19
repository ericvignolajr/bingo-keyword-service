package inmemory

import (
	"errors"
	"fmt"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/google/uuid"
)

type SubjectStore struct {
	Store map[uuid.UUID][]*domain.Subject
}

func NewSubjectStore() *SubjectStore {
	store := make(map[uuid.UUID][]*domain.Subject)
	return &SubjectStore{
		Store: store,
	}
}

func (s *SubjectStore) Read(UserID uuid.UUID) ([]*domain.Subject, error) {
	subjects, ok := s.Store[UserID]
	if !ok {
		return nil, nil
	}

	subjectPointers := make([]*domain.Subject, len(subjects))
	for i := range subjects {
		subjectPointers[i] = subjects[i]
	}

	return subjectPointers, nil
}

func (s *SubjectStore) ReadByName(UserId uuid.UUID, SubjectName string) (*domain.Subject, error) {
	for _, v := range s.Store[UserId] {
		if v.Name == SubjectName {
			return v, nil
		}
	}

	return nil, fmt.Errorf("subject %s for user %s could not be found", SubjectName, UserId)
}

func (s *SubjectStore) ReadByID(subjectID uuid.UUID) (*domain.Subject, error) {
	for _, user := range s.Store {
		for _, subject := range user {
			if subject.Id == subjectID {
				return subject, nil
			}
		}
	}

	return nil, nil
}

func (s *SubjectStore) Create(UserId uuid.UUID, Subject *domain.Subject) (*domain.Subject, error) {
	subjectToCreate, _ := s.ReadByName(UserId, Subject.Name)
	if subjectToCreate != nil {
		return subjectToCreate, errors.New("subject already exists")
	}

	subjects, ok := s.Store[UserId]
	if !ok {
		s.Store[UserId] = []*domain.Subject{Subject}
	} else {
		s.Store[UserId] = append(subjects, Subject)
	}

	return Subject, nil
}

func (s *SubjectStore) Delete(userID uuid.UUID, subjectID uuid.UUID) error {
	var subjectSlice = s.Store[userID]
	var nullIndex struct {
		index int
		found bool // found is true if the index was set
	}
	for idx, subject := range subjectSlice {
		if subject.Id == subjectID {
			nullIndex = struct {
				index int
				found bool
			}{
				index: idx,
				found: true,
			}
		}
	}

	if nullIndex.found == false {
		return nil
	}

	length := len(subjectSlice)
	subjectSlice[nullIndex.index] = subjectSlice[length-1]
	subjectSlice[length-1] = nil
	s.Store[userID] = subjectSlice[:length-1]
	return nil
}
