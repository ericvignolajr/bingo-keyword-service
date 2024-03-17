package sql_test

import (
	"errors"
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	store "github.com/ericvignolajr/bingo-keyword-service/pkg/stores/sql"
)

func TestDeleteSubject(t *testing.T) {
	user, _ := domain.NewUser("foo@example.com", "")

	subject, _ := domain.NewSubject("science", user.ID)

	unit, _ := domain.NewUnit("electro magnets")
	unit.SubjectID = subject.ID
	keyword, _ := domain.NewKeyword("tesla coil", "a tesla coil")
	unit.AddKeyword(keyword)
	subject.AddUnit(unit)

	unit2, _ := domain.NewUnit("geology")
	unit2.SubjectID = subject.ID
	keyword2, _ := domain.NewKeyword("rock", "a rock")
	unit2.AddKeyword(keyword2)
	subject.AddUnit(unit2)

	user.AddSubject(subject)

	uStore, _ := store.NewSQLUserStore()
	uStore.Save(user)

	userFromDB, _ := uStore.ReadById(user.ID)

	subjectFromDB, err := userFromDB.FindSubject(subject.ID)
	if err != nil {
		if errors.Is(err, domain.ErrSubjectDoesNotExist) {
			t.Errorf("expected user from database to have subject %v", subject)
		}
	}

	unitFromDB, err := subjectFromDB.FindUnitByID(unit.ID)
	if err != nil {
		t.Errorf("expected subject from database to have unit %v", unit)
	}

	keyword, err = unitFromDB.FindKeyword(keyword.Id)
	if err != nil {
		t.Errorf("expected unit from database to have keyword %v", keyword.Name)
	}

	unitFromDB2, err := subjectFromDB.FindUnitByID(unit2.ID)
	if err != nil {
		t.Errorf("expected subject from database to have unit %v", unit)
	}

	keyword2, err = unitFromDB2.FindKeyword(keyword2.Id)
	if err != nil {
		t.Errorf("expected unit from database to have keyword %v", keyword2.Name)
	}

	sStore, _ := store.NewSQLSubjectStore()
	sStore.Delete(user.ID, subject.ID)

	userAfterDelete, _ := uStore.ReadById(user.ID)

	subjectFromDBAfterDelete, err := userAfterDelete.FindSubject(subject.ID)
	if err != nil {
		if !errors.Is(err, domain.ErrSubjectDoesNotExist) {
			t.Errorf("subject should be deleted and user should not longer have subject %s", subjectFromDBAfterDelete.Name)
		}
	}

}
