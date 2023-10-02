package usecases

import (
	"sort"
	"strings"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores"
	"github.com/google/uuid"
)

type ReadSubjectsRequest struct {
	UserID uuid.UUID
}

type SubjectOutput struct {
	ID   uuid.UUID
	Name string
}
type ReadSubjectsResponse struct {
	Subjects []SubjectOutput
	Err      error
}

func ReadSubjects(req ReadSubjectsRequest, subjectStore stores.Subject) ReadSubjectsResponse {
	subjects, err := subjectStore.Read(req.UserID)
	if err != nil {
		return ReadSubjectsResponse{
			nil,
			err,
		}
	}

	subjectOutput := make([]SubjectOutput, len(subjects))
	for i, v := range subjects {
		subjectOutput[i] = SubjectOutput{v.Id, v.Name}
	}

	sort.Slice(subjectOutput, func(i, j int) bool {
		return strings.ToLower(subjectOutput[i].Name) < strings.ToLower(subjectOutput[j].Name)
	})

	return ReadSubjectsResponse{
		subjectOutput,
		nil,
	}
}
