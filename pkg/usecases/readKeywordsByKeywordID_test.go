package usecases_test

import (
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/sql"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/google/go-cmp/cmp"
)

func TestReadKeywordsByKeywordID(t *testing.T) {
	uStore, _ := sql.NewSQLUserStore()
	user, _ := domain.NewUser("foo@example.com", "foobar")

	scienceSubj, _ := domain.NewSubject("science", user.ID)
	electroMagnetsUnit, _ := domain.NewUnit("electro-magnets")
	magnetKeyword, _ := domain.NewKeyword("magnet", "a magnetized object")
	rocksKeyword, _ := domain.NewKeyword("rocks", "a rock")
	geologyUnit, _ := domain.NewUnit("geology")
	batteryKeyword, _ := domain.NewKeyword("battery", "powers devices")
	electroMagnetsUnit.AddKeyword(magnetKeyword)
	electroMagnetsUnit.AddKeyword(rocksKeyword)
	geologyUnit.AddKeyword(batteryKeyword)
	scienceSubj.AddUnit(electroMagnetsUnit)
	scienceSubj.AddUnit(geologyUnit)
	user.AddSubject(scienceSubj)

	mathSubj, _ := domain.NewSubject("math", user.ID)
	geometryUnit, _ := domain.NewUnit("geometry")
	triangleKeyword, _ := domain.NewKeyword("triangle", "three-sided shape")
	geometryUnit.AddKeyword(triangleKeyword)
	mathSubj.AddUnit(geometryUnit)
	user.AddSubject(mathSubj)

	uStore.Save(user)

	readKeywordsByKeywordID := usecases.ReadKeywordsByKeywordID{
		UserStore: uStore,
	}

	keywordsToReadFromDB := []*domain.Keyword{
		magnetKeyword,
		rocksKeyword,
		batteryKeyword,
		triangleKeyword,
	}

	for _, k := range keywordsToReadFromDB {
		keywordFromDB, err := readKeywordsByKeywordID.Exec(user.ID, k.ID)
		if err != nil {
			t.Error(err)
		}

		/*
			have to wrap k in a slice because usecases.ReadKeywordsByKeywordID
			returns a slice containing a keyword NOT the keyword directly
		*/
		isEqual := cmp.Equal([]*domain.Keyword{k}, keywordFromDB)
		if !isEqual {
			t.Error(cmp.Diff(k, keywordFromDB))
		}
	}
}
