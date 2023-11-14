package usecases_test

import (
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/inmemory"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/stretchr/testify/assert"
)

func TestCreateUnit(t *testing.T) {
	uStore := inmemory.UserStore{}
	sStore := inmemory.NewSubjectStore()
	createUnit := usecases.CreateUnit{
		UserStore:    &uStore,
		SubjectStore: sStore,
	}

	err := uStore.CreateAccount("example@bingboard.com")
	if err != nil {
		t.Error(err)
	}

	u, err := uStore.ReadByEmail("example@bingboard.com")
	if err != nil {
		t.Error(err)
	}

	s, err := domain.NewSubject("Science", u.Id)
	if err != nil {
		t.Error(err)
	}
	u.AddSubject(*s)
	uStore.Save(u)

	newUnit, err := createUnit.Exec("Electro magnets", u.Id, s.Id)
	if err != nil {
		t.Error(err)
	}

	assert.Contains(t, u.Subjects[0].Units, newUnit)
}
