package usecases_test

import (
	"testing"

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
	usecases.CreateSubject(createReq, &userStore, subjectStore)

	readReq := usecases.ReadSubjectRequest{
		uid,
		createReq.SubjectName,
	}

	readRes := usecases.ReadSubject(readReq, subjectStore)
	if readRes.Err != nil {
		t.Error(readRes.Err)
	}

	assert.Equal(t, readReq.SubjectName, readRes.Subject.Name)
}
