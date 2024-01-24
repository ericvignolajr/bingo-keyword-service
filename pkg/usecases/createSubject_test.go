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

	createSubject := usecases.CreateSubject{
		UserStore:    &userStore,
		SubjectStore: subjectStore,
	}

	res, err := createSubject.Exec("Science", uid)
	if err != nil {
		t.Error(err)
	}

	assert.NotEqual(t, nil, res)
	assert.Equal(t, "Science", res.Name)
}

func TestReadByName(t *testing.T) {
	userStore := inmemory.UserStore{}
	uid, _ := userStore.Create("foo@example.com", "supersecret")

	subjectToAdd, _ := domain.NewSubject("Science", uid)

	subjectStore := inmemory.NewSubjectStore()
	subjectStore.Create(uid, subjectToAdd)

	subjectToRead, err := subjectStore.ReadByName(uid, subjectToAdd.Name)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, subjectToAdd.ID, subjectToRead.ID)
}
