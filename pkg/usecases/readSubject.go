package usecases

import (
	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
	"github.com/google/uuid"
)

type ReadSubjectRequest struct {
	UserId      uuid.UUID
	SubjectName string
}

type ReadSubjectResponse struct {
	Subject *domain.Subject
	Err     error
}

func ReadSubject(req ReadSubjectRequest, subjectStore stores.Subject) ReadSubjectResponse {
	subject, err := subjectStore.ReadByName(req.UserId, req.SubjectName)
	if err != nil {
		return ReadSubjectResponse{
			nil,
			err,
		}
	}

	return ReadSubjectResponse{
		subject,
		nil,
	}
}
