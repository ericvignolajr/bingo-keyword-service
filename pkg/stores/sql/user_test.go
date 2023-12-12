package sql_test

import (
	"fmt"
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/sql"
)

func TestReadById(t *testing.T) {
	uStore, err := sql.NewSQLUserStore()
	if err != nil {
		t.Error(err)
	}

	user, _ := domain.NewUser("foo@example.com", "foobar")
	subj, _ := domain.NewSubject("electro-magnets", user.ID)
	user.AddSubject(*subj)
	uStore.DB.Model(user).Create(user)
	subjFromDB := &domain.Subject{}
	subjResult := uStore.DB.Model(subj).First(subjFromDB)
	subjResult.Scan(subjFromDB)
	fmt.Println(subjFromDB)
	uStore.ReadById(user.ID)
}
