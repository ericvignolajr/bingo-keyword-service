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

type CreateSubjectPresenter interface {
	PresentCreateSubject(CreateSubjectResponse)
}

type CreateSubject struct {
	UserStore    stores.User
	SubjectStore stores.Subject
	Presenter    CreateSubjectPresenter
}

func (c *CreateSubject) Exec(req CreateSubjectRequest) CreateSubjectResponse {
	var result CreateSubjectResponse
	user, err := c.UserStore.ReadById(req.UserId)
	if err != nil {
		result = CreateSubjectResponse{
			false,
			err,
		}
	}

	subject, err := domain.NewSubject(req.SubjectName)
	if err != nil {
		result = CreateSubjectResponse{
			false,
			err,
		}
	}

	_, err = c.SubjectStore.Create(user.Id, subject)
	if err != nil {
		result = CreateSubjectResponse{
			false,
			err,
		}
	}

	result = CreateSubjectResponse{
		true,
		nil,
	}

	c.Presenter.PresentCreateSubject(result)
	return result
}
