package domain

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewSubject(t *testing.T) {
	expected := &Subject{
		Id:   uuid.New(),
		Name: "Science",
	}

	s, err := NewSubject("Science")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, expected.Name, s.Name)
}

func TestAddUnit(t *testing.T) {
	u, err := NewUnit("Electromagnets")
	if err != nil {
		t.Error(err)
	}
	s, err := NewSubject("Science")
	if err != nil {
		t.Error(err)
	}

	newU, err := s.AddUnit(*u)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, u, newU)

}

func TestIsDuplicateUnit(t *testing.T) {
	u1, err := NewUnit("Electromagnets")
	if err != nil {
		t.Error(err)
	}
	u2, err := NewUnit("Electromagnets")
	if err != nil {
		t.Error(err)
	}
	s, err := NewSubject("Science")
	if err != nil {
		t.Error(err)
	}

	u2, err = s.AddUnit(*u1)
	if err != nil {
		t.Error(err)
	}
	isDuplicate := s.IsDuplicateUnit(*u2)
	assert.Equal(t, true, isDuplicate)
}
