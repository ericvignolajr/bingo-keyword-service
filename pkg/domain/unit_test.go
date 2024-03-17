package domain_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/google/go-cmp/cmp"
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
	err = u.AddKeyword(k1)
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

	err = u.AddKeyword(k1)
	if err != nil {
		t.Errorf("unexpected input when adding keyword to unit %s", err)
	}
	err = u.AddKeyword(k2)
	assert.Equal(t, true, errors.Is(err, domain.ErrDuplicateKeyword))

}

func TestIsDuplicateKeyword(t *testing.T) {
	u, _ := domain.NewUnit("magnets and electricity")

	k1, err := domain.NewKeyword("magnet", "definition of a magnet")
	if err != nil {
		t.Errorf("unexpected input when adding keyword to unit %s", err)
	}

	err = u.AddKeyword(k1)
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

func TestEqual(t *testing.T) {
	type testCase struct {
		name     string
		unit1    *domain.Unit
		unit2    *domain.Unit
		expected bool
	}

	u1, _ := domain.NewUnit("magnets and electricity")
	u2, _ := domain.NewUnit("weather and biology")

	testCase1 := testCase{
		name:     "two different units",
		unit1:    u1,
		unit2:    u2,
		expected: false,
	}

	u3, _ := domain.NewUnit(u1.Name)
	testCase2 := testCase{
		name:     "ids don't match",
		unit1:    u1,
		unit2:    u3,
		expected: false,
	}

	testCases := []testCase{
		testCase1,
		testCase2,
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := testCase.unit1.Equal(testCase.unit2)
			if result != testCase.expected {
				fmt.Println(cmp.Diff(testCase.unit1, testCase.unit2))
				t.Errorf("on test %s, got %t, expected %t", testCase.name, result, testCase.expected)
			}
		})
	}

	isEqual := cmp.Equal(u1, u2)
	assert.Equal(t, false, isEqual)

}
