package usecases_test

import (
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/sql"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/stretchr/testify/assert"
)

func TestUpdateSubject(t *testing.T) {
	user, err := domain.NewUser("foo@example.com", "foo")
	if err != nil {
		t.Error(err)
	}
	subject, err := domain.NewSubject("electro magnets", user.ID)
	if err != nil {
		t.Error(err)
	}

	err = user.AddSubject(subject)
	if err != nil {
		t.Error(err)
	}

	uStore, _ := sql.NewSQLUserStore()
	uStore.Save(user)

	updateSubject := usecases.UpdateSubject{
		UserStore: uStore,
	}

	err = updateSubject.Exec(user.ID, subject.ID, "a new name")
	if err != nil {
		t.Error(err)
	}

	userFromDB, _ := uStore.ReadById(user.ID)
	subjectFromDB, _ := userFromDB.FindSubject(subject.ID)
	assert.Equal(t, "A new name", subjectFromDB.Name)
}
