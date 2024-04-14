package usecases

import (
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/sql"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestCreateKeyword(t *testing.T) {
	uStore, _ := sql.NewSQLUserStore()
	err := uStore.CreateAccount("example@bingboard.com")
	if err != nil {
		t.Error(err)
	}

	u, err := uStore.ReadByEmail("example@bingboard.com")
	if err != nil {
		t.Error(err)
	}

	s, err := domain.NewSubject("Science", u.ID)
	if err != nil {
		t.Error(err)
	}

	unit, err := domain.NewUnit("electro-magnets")
	if err != nil {
		t.Error(err)
	}
	s.AddUnit(unit)
	u.AddSubject(s)
	uStore.Save(u)

	createKeywordUsecase := CreateKeyword{
		UserStore: uStore,
	}

	keyword, err := createKeywordUsecase.Exec("magnet", u.ID, s.ID, unit.ID)
	if err != nil {
		t.Error(err)
	}

	userAfterCreate, err := uStore.ReadById(u.ID)
	if err != nil {
		t.Error(err)
	}
	subjectAfterCreate, err := userAfterCreate.FindSubject(s.ID)
	if err != nil {
		t.Error(err)
	}
	unitAfterCreate, err := subjectAfterCreate.FindUnitByID(unit.ID)
	if err != nil {
		t.Error(err)
	}
	keywordFromDB, err := unitAfterCreate.FindKeyword(keyword.ID)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, true, cmp.Equal(keyword, keywordFromDB))
}
