package domain_test

import (
	"fmt"
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUnit(t *testing.T) {
	expected := &domain.Unit{
		ID:   uuid.New(),
		Name: "Electricity and Magnets",
	}

	u, err := domain.NewUnit(expected.Name)
	if err != nil {
		t.Error(err)
	}

	if expected.Name != u.Name {
		t.Error("expected unit should match the one returned by factory function")
	}
}

func TestAddKeyword(t *testing.T) {
	u, err := domain.NewUnit("electricity and magnets")
	if err != nil {
		t.Errorf("could not construct unit using factory function %s", err)
	}

	k1, _ := domain.NewKeyword("magnets", "defintion of a magnet")
	_, err = u.AddKeyword(*k1)
	if err != nil {
		t.Errorf("unexpected input when adding keyword to unit %s", err)
	}

	assert.Contains(t, u.Keywords, k1)
}

func TestAddKeyword_DuplicateKeyword(t *testing.T) {
	u, err := domain.NewUnit("electricity and magnets")
	if err != nil {
		t.Errorf("could not construct unit using factory function %s", err)
	}

	k1, _ := domain.NewKeyword("magnets", "defintion of a magnet")
	k2, _ := domain.NewKeyword("magnets", "defintion of a magnet")

	_, err = u.AddKeyword(*k1)
	if err != nil {
		t.Errorf("unexpected input when adding keyword to unit %s", err)
	}
	_, err = u.AddKeyword(*k2)
	assert.EqualError(t, err, fmt.Sprintf("unit %s already contains keyword %s", u.Name, k2.Name))

}

func TestIsDuplicateKeyword(t *testing.T) {
	u, _ := domain.NewUnit("magnets and electricity")

	k1, err := domain.NewKeyword("magnet", "definition of a magnet")
	if err != nil {
		t.Errorf("unexpected input when adding keyword to unit %s", err)
	}

	_, err = u.AddKeyword(*k1)
	if err != nil {
		t.Errorf("unexpected input when adding keyword to unit %s", err)
	}

	k2, err := domain.NewKeyword("magnet", "definition of a magnet")
	if err != nil {
		t.Errorf("unexpected input when adding keyword to unit %s", err)
	}

	isDuplicate := u.IsDuplicateKeyword(*k2)
	assert.Equal(t, true, isDuplicate)
}
