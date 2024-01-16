package sql_test

import (
	"fmt"
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/sql"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestReadById(t *testing.T) {
	uStore, err := sql.NewSQLUserStore()
	if err != nil {
		t.Error(err)
	}

	user, _ := domain.NewUser("foo@example.com", "foobar")
	subj, _ := domain.NewSubject("science", user.ID)
	unit, _ := domain.NewUnit("electro-magnets")
	subj.AddUnit(*unit)
	user.AddSubject(*subj)
	uStore.DB.Create(user)

	userFromDB, err := uStore.ReadById(user.ID)
	if err != nil {
		t.Error(err)
	}

	isEqual := cmp.Equal(user, userFromDB)
	assert.Equal(t, true, isEqual)
}

func TestReadByEmail(t *testing.T) {
	uStore, err := sql.NewSQLUserStore()
	if err != nil {
		t.Error(err)
	}

	userEmail := "test@example.com"
	user, _ := domain.NewUser(userEmail, "foobaz")
	uStore.DB.Create(user)

	userFromDBByEmail, err := uStore.ReadByEmail(userEmail)
	if err != nil {
		t.Error(err)
	}

	if cmp.Equal(user, userFromDBByEmail) == false {
		fmt.Printf(cmp.Diff(user, userFromDBByEmail))
		t.FailNow()
	}
}

func TestSave(t *testing.T) {
	user, _ := domain.NewUser("foo@example.com", "baz")
	uStore, _ := sql.NewSQLUserStore()

	uStore.DB.Create(user)
	userFromDB, _ := uStore.ReadById(user.ID)
	newEmail := "updated@example.com"
	userFromDB.Email = newEmail

	uStore.Save(userFromDB)

	userAfterUpdate, _ := uStore.ReadById(user.ID)

	assert.Equal(t, newEmail, userAfterUpdate.Email)
}
