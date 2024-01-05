package sql_test

import (
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/sql"
	"github.com/stretchr/testify/assert"
)

func TestReadById(t *testing.T) {
	uStore, err := sql.NewSQLUserStore()
	if err != nil {
		t.Error(err)
	}

	user, _ := domain.NewUser("foo@example.com", "foobar")
	subj, _ := domain.NewSubject("electro-magnets", user.ID)
	user.AddSubject(*subj)
	uStore.DB.Create(user)

	subjFromDB := &domain.Subject{}
	err = uStore.DB.First(&subjFromDB).Error
	if err != nil {
		t.Error(err)
	}
	userFromDB, err := uStore.ReadById(user.ID)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, user, userFromDB)
}
