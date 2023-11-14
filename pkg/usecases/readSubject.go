package usecases

import (
	"fmt"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
	"github.com/google/uuid"
)

type ReadSubjectResponse struct {
	Subject *domain.Subject
	Err     error
}

type ReadSubject struct {
	UserStore    stores.User
	SubjectStore stores.Subject
}

func (r *ReadSubject) Exec(name string, userID uuid.UUID) (*domain.Subject, error) {
	user, err := r.UserStore.ReadById(userID)
	if err != nil {
		return nil, fmt.Errorf("in readSubject: %w", err)
	}

	subject, err := user.FindSubjectByName(name)
	if err != nil {
		return nil, fmt.Errorf("in readSubject: %w", err)
	}

	return subject, nil
}

type ReadSubjectByID struct {
	UserStore    stores.User
	SubjectStore stores.Subject
}

func (r *ReadSubjectByID) ReadSubjectByID(userID uuid.UUID, subjectID uuid.UUID) (*domain.Subject, error) {
	user, err := r.UserStore.ReadById(userID)
	if err != nil {
		return nil, fmt.Errorf("in readSubject: %w", err)
	}

	s, err := user.FindSubject(subjectID)
	if err != nil {
		return nil, fmt.Errorf("in readSubject: %w", err)
	}
	return s, err
}
