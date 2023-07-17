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

type ReadSubjectPresenter interface {
	PresentReadSubject(ReadSubjectResponse) error
}

type ReadSubject struct {
	SubjectStore stores.Subject
	Presenter    ReadSubjectPresenter
}

func (r *ReadSubject) Exec(req ReadSubjectRequest) ReadSubjectResponse {
	var result ReadSubjectResponse
	subject, err := r.SubjectStore.ReadByName(req.UserId, req.SubjectName)
	if err != nil {
		result = ReadSubjectResponse{
			nil,
			err,
		}
	}

	result = ReadSubjectResponse{
		subject,
		nil,
	}

	r.Presenter.PresentReadSubject(result)
	return result
}
