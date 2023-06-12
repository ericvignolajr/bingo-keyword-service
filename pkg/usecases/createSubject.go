package usecases

import (
	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
	"github.com/google/uuid"
)

type CreateSubjectRequest struct {
	UserId      uuid.UUID
	SubjectName string
}

type CreateSubjectResponse struct {
	Ok  bool
	Err error
}

func CreateSubject(req CreateSubjectRequest, userStore stores.User, subjectStore stores.Subject) CreateSubjectResponse {
	// get user from id
	user, err := userStore.ReadById(req.UserId)
	if err != nil {
		return CreateSubjectResponse{
			false,
			err,
		}
	}
	// create subject from name
	subject, err := domain.NewSubject(req.SubjectName)
	if err != nil {
		return CreateSubjectResponse{
			false,
			err,
		}
	}

	_, err = subjectStore.Create(user.Id, subject)
	if err != nil {
		return CreateSubjectResponse{
			false,
			err,
		}
	}

	// return result on success or error
	return CreateSubjectResponse{
		true,
		nil,
	}
}
