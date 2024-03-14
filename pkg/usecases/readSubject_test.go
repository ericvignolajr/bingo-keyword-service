package usecases_test

import (
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/inmemory"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/stretchr/testify/assert"
)

func TestReadSubject(t *testing.T) {
	userStore := inmemory.UserStore{}
	subjectStore := inmemory.NewSubjectStore()

	err := userStore.CreateAccount("foo@example.com")
	if err != nil {
		t.Error(err)
	}
	user, err := userStore.ReadByEmail("foo@example.com")
	if err != nil {
		t.Error(err)
	}

	subject, _ := domain.NewSubject("Science", user.ID)
	user.AddSubject(subject)
	userStore.Save(user)

	readSubject := usecases.ReadSubject{
		UserStore:    &userStore,
		SubjectStore: subjectStore,
	}
	readRes, err := readSubject.Exec(subject.Name, user.ID)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "Science", readRes.Name)
}
