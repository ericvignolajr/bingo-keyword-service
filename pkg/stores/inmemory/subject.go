package inmemory

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
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
		if strings.ToLower(v.Name) == strings.ToLower(SubjectName) {
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
		return subjectToCreate, errors.New(stores.ErrSubjectExists)
	}

	Subject.OwnerID = UserId
	subjects, ok := s.Store[UserId]
	if !ok {
		s.Store[UserId] = []*domain.Subject{Subject}
	} else {
		s.Store[UserId] = append(subjects, Subject)
	}

	return Subject, nil
}

func (s *SubjectStore) Update(Subject *domain.Subject) (*domain.Subject, error) {
	subjectToUpdate, _ := s.ReadByID(Subject.Id)
	if subjectToUpdate == nil {
		newSubject, err := s.Create(Subject.OwnerID, Subject)
		if err != nil {
			return nil, err
		}
		return newSubject, nil
	}

	subjects, ok := s.Store[Subject.OwnerID]
	if !ok {
		return nil, errors.New(fmt.Sprintf("could not find record in in-memory subject store for user with id: %s", Subject.OwnerID.String()))
	}
	for i, v := range subjects {
		if v.Id == Subject.Id {
			subjects[i] = Subject
		}
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
