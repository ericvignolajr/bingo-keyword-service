package usecases_test

import (
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/stores/sql"
	"github.com/ericvignolajr/bingo-keyword-service/pkg/usecases"
	"github.com/stretchr/testify/assert"
)

func TestReadKeywordsByUserID(t *testing.T) {
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

	readKeywordsByUserID := usecases.ReadKeywordsByUserID{
		UserStore: uStore,
	}

	keywords, err := readKeywordsByUserID.Exec(user.ID)
	if err != nil {
		t.Error(err)
	}

	assert.ElementsMatch(t, []*domain.Keyword{
		magnetKeyword,
		rocksKeyword,
		batteryKeyword,
		triangleKeyword,
	}, keywords)
}
