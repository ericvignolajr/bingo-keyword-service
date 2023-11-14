package usecases_test

import (
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/inmemory"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/stretchr/testify/assert"
)

func TestDeleteSubject(t *testing.T) {
	uStore := inmemory.UserStore{}
	sStore := inmemory.NewSubjectStore()

	uStore.CreateAccount("foo@example.com")
	user, err := uStore.ReadByEmail("foo@example.com")
	if err != nil {
		t.Error(err)
	}

	subject1, _ := domain.NewSubject("Science", user.Id)
	subject2, _ := domain.NewSubject("Physics", user.Id)
	user.AddSubject(*subject1)
	user.AddSubject(*subject2)
	uStore.Save(user)

	deleteSubject := usecases.DeleteSubject{
		SubjectStore: sStore,
		UserStore:    &uStore,
	}

	updatedUser, err := uStore.ReadById(user.Id)
	if err != nil {
		t.Error(err)
	}

	deleteSubject.Exec(user.Id, subject1.Id)
	expected := []*domain.Subject{subject2}

	assert.Equal(t, expected, updatedUser.Subjects)
}
