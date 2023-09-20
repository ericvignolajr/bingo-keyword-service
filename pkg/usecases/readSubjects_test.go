package usecases_test

import (
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/inmemory"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/stretchr/testify/assert"
)

func TestReadSubjects(t *testing.T) {
	userStore := inmemory.UserStore{}
	uid, err := userStore.Create("foobar@example.com", "supersecret")
	if err != nil {
		panic(err)
	}

	subjectStore := inmemory.NewSubjectStore()
	s1, err := domain.NewSubject("Math")
	if err != nil {
		panic(err)
	}
	subjectStore.Create(uid, s1)

	s2, err := domain.NewSubject("Science")
	if err != nil {
		panic(err)
	}
	subjectStore.Create(uid, s2)

	readSubjectsRequest := usecases.ReadSubjectsRequest{uid}
	res := usecases.ReadSubjects(readSubjectsRequest, subjectStore)
	if res.Err != nil {
		panic(res.Err)
	}

	expected := usecases.ReadSubjectsResponse{
		[]usecases.SubjectOutput{
			{s1.Id, s1.Name},
			{s2.Id, s2.Name},
		},
		nil,
	}
	assert.Equal(t, expected, res)
}
