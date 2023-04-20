package domain

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewUnit(t *testing.T) {
	expected := &Unit{
		Id:   uuid.New(),
		Name: "Electricity and Magnets",
	}

	u, err := NewUnit(expected.Name)
	if err != nil {
		t.Error(err)
	}

	if expected.Name != u.Name {
		t.Error("expected unit should match the one returned by factory function")
	}
}

func TestAddKeyword(t *testing.T) {
	u, err := NewUnit("electricity and magnets")
	if err != nil {
		t.Errorf("could not construct unit using factory function %s", err)
	}

	k1, _ := NewKeyword("magnets", "defintion of a magnet")
	_, err = u.AddKeyword(*k1)
	if err != nil {
		t.Errorf("unexpected input when adding keyword to unit %s", err)
	}

	assert.Contains(t, u.Keywords, k1)
}

func TestAddKeyword_DuplicateKeyword(t *testing.T) {
	u, err := NewUnit("electricity and magnets")
	if err != nil {
		t.Errorf("could not construct unit using factory function %s", err)
	}

	k1, _ := NewKeyword("magnets", "defintion of a magnet")
	k2, _ := NewKeyword("magnets", "defintion of a magnet")

	_, err = u.AddKeyword(*k1)
	if err != nil {
		t.Errorf("unexpected input when adding keyword to unit %s", err)
	}
	_, err = u.AddKeyword(*k2)
	assert.EqualError(t, err, fmt.Sprintf("unit %s already contains keyword %s", u.Name, k2.Name))

}

func TestIsDuplicateKeyword(t *testing.T) {
	u, _ := NewUnit("magnets and electricity")

	k1, err := NewKeyword("magnet", "definition of a magnet")
	if err != nil {
		t.Errorf("unexpected input when adding keyword to unit %s", err)
	}

	_, err = u.AddKeyword(*k1)
	if err != nil {
		t.Errorf("unexpected input when adding keyword to unit %s", err)
	}

	k2, err := NewKeyword("magnet", "definition of a magnet")
	if err != nil {
		t.Errorf("unexpected input when adding keyword to unit %s", err)
	}

	isDuplicate := u.IsDuplicateKeyword(*k2)
	assert.Equal(t, true, isDuplicate)
}
