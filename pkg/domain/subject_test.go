package domain_test

import (
	"testing"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewSubject(t *testing.T) {
	expected := &domain.Subject{
		Id:   uuid.New(),
		Name: "Science",
	}

	s, err := domain.NewSubject("Science")
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, expected.Name, s.Name)
}

func TestNewSubjectEmptyName(t *testing.T) {
	s, err := domain.NewSubject("")
	if s != nil {
		t.Error()
	}
	assert.EqualError(t, err, domain.ErrSubjectNameEmpty)
}

func TestAddUnit(t *testing.T) {
	u, err := domain.NewUnit("Electromagnets")
	if err != nil {
		t.Error(err)
	}
	s, err := domain.NewSubject("Science")
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
	u1, err := domain.NewUnit("Electromagnets")
	if err != nil {
		t.Error(err)
	}
	u2, err := domain.NewUnit("Electromagnets")
	if err != nil {
		t.Error(err)
	}
	s, err := domain.NewSubject("Science")
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
