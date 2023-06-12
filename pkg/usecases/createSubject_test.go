package usecases_test

import (
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/inmemory"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/stretchr/testify/assert"
)

func TestCreateSubject(t *testing.T) {
	userStore := inmemory.UserStore{}
	uid, _ := userStore.Create("foo@example.com", "supersecret")

	subjectStore := inmemory.NewSubjectStore()

	req := usecases.CreateSubjectRequest{
		UserId:      uid,
		SubjectName: "Science",
	}

	res := usecases.CreateSubject(req, &userStore, subjectStore)

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
