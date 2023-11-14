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
	user, err := userStore.ReadById(uid)
	if err != nil {
		t.Error(err)
	}

	s1, err := domain.NewSubject("Math", uid)
	if err != nil {
		panic(err)
	}

	s2, err := domain.NewSubject("Science", uid)
	if err != nil {
		panic(err)
	}

	user.AddSubject(*s1)
	user.AddSubject(*s2)

	readSubjects := usecases.ReadSubjects{
		UserStore: &userStore,
	}

	res := readSubjects.Exec(uid)
	if res.Err != nil {
		t.Error(res.Err.Error())
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
