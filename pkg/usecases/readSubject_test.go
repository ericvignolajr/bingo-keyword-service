package usecases_test

import (
	"testing"

	outputports "github.com/ericvignolajr/bingo-keyword-service/pkg/output_ports"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/inmemory"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/stretchr/testify/assert"
)

func TestReadSubject(t *testing.T) {
	userStore := inmemory.UserStore{}
	uid, _ := userStore.Create("foo@example.com", "supersecret")

	subjectStore := inmemory.NewSubjectStore()
	createReq := usecases.CreateSubjectRequest{
		uid,
		"Science",
	}

	mockPresenter := outputports.MockPresenter{}

	createSubjectUsecase := usecases.CreateSubject{
		UserStore:    &userStore,
		SubjectStore: subjectStore,
		Presenter:    &mockPresenter,
	}
	createSubjectUsecase.Exec(createReq)

	readReq := usecases.ReadSubjectRequest{
		uid,
		createReq.SubjectName,
	}

	readSubject := usecases.ReadSubject{
		SubjectStore: subjectStore,
		Presenter:    &mockPresenter,
	}
	readRes := readSubject.Exec(readReq)
	if readRes.Err != nil {
		t.Error(readRes.Err)
	}

	assert.Equal(t, readReq.SubjectName, readRes.Subject.Name)
}
