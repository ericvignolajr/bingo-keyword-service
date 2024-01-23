package usecases_test

import (
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/inmemory"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/google/go-cmp/cmp"
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

	res, err := readSubjects.Exec(uid, nil)
	if err != nil {
		t.Error(err)
	}
	res, err = readSubjects.Exec(uid, nil)
	if err != nil {
		t.Error(err)
	}

	expected := []domain.Subject{
		{ID: s1.ID, Name: s1.Name, UserID: uid},
		{ID: s2.ID, Name: s2.Name, UserID: uid},
	}

	isEqual := cmp.Equal(expected, res)
	if isEqual != true {
		diff := cmp.Diff(expected, res)
		t.Errorf("\non TestReadSubjects%v", diff)
	}
	assert.Equal(t, true, isEqual)
}
