package usecases_test

import (
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	outputports "github.com/ericvignolajr/bingo-keyword-service/pkg/output_ports"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/inmemory"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/viewers"
	"github.com/stretchr/testify/assert"
)

func TestCreateSubject(t *testing.T) {
	userStore := inmemory.UserStore{}
	uid, _ := userStore.Create("foo@example.com", "supersecret")

	subjectStore := inmemory.NewSubjectStore()

	mockPresenter := outputports.MockPresenter{
		Viewer: &viewers.MockViewer{},
	}

	req := usecases.CreateSubjectRequest{
		UserId:      uid,
		SubjectName: "Science",
	}

	createSubject := &usecases.CreateSubject{
		UserStore:    &userStore,
		SubjectStore: subjectStore,
		Presenter:    &mockPresenter,
	}

	res := createSubject.Exec(req)

	assert.Equal(t, true, res.Ok)
}

func TestReadByName(t *testing.T) {
	userStore := inmemory.UserStore{}
	uid, _ := userStore.Create("foo@example.com", "supersecret")

	subjectToAdd, _ := domain.NewSubject("Science")

	subjectStore := inmemory.NewSubjectStore()
	subjectStore.Create(uid, subjectToAdd)

	subjectToRead, err := subjectStore.ReadByName(uid, subjectToAdd.Name)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, subjectToAdd.Id, subjectToRead.Id)
}
