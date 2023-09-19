package usecases_test

import (
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/inmemory"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDeleteSubject(t *testing.T) {
	uID, _ := uuid.NewUUID()
	subject1, _ := domain.NewSubject("Science")
	subject2, _ := domain.NewSubject("Physics")

	sStore := inmemory.NewSubjectStore()
	sStore.Create(uID, subject1)
	sStore.Create(uID, subject2)
	deleteSubject := usecases.DeleteSubject{
		SubjectStore: sStore,
	}

	deleteSubject.Exec(uID, subject1.Id)
	expected := []*domain.Subject{subject2}
	assert.Equal(t, expected, sStore.Store[uID])
}
