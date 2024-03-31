package domain

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
)

var (
	ErrSubjectNameEmpty = errors.New("subject name cannot be empty")
	ErrDuplicateUnit    = errors.New("unit names must be unique within a subject")
)

type Subject struct {
	ID     uuid.UUID `gorm:"primaryKey"`
	Name   string
	Units  []*Unit
	UserID uuid.UUID
}

func NewSubject(name string, ownerID uuid.UUID) (*Subject, error) {
	if name == "" {
		return nil, ErrSubjectNameEmpty
	}
	capitalizedName := strings.ToUpper(string(name[0])) + name[1:]
	return &Subject{
		ID:     uuid.New(),
		Name:   capitalizedName,
		Units:  nil,
		UserID: ownerID,
	}, nil
}

func (s *Subject) AddUnit(u *Unit) error {
	isDuplicate := s.IsDuplicateUnit(*u)
	if isDuplicate {
		return fmt.Errorf("%w, subject %s already contains unit %s", ErrDuplicateUnit, s.Name, u.Name)
	}

	capitalizedName := strings.ToUpper(string(u.Name[0])) + u.Name[1:]
	u.Name = capitalizedName

	s.Units = append(s.Units, u)
	sort.SliceStable(s.Units, func(i, j int) bool {
		return strings.ToLower(s.Units[i].Name) < strings.ToLower(s.Units[j].Name)
	})
	return nil
}

func (s *Subject) IsDuplicateUnit(u Unit) bool {
	for _, v := range s.Units {
		if strings.EqualFold(v.Name, u.Name) {
			return true
		}
	}

	return false
}

func (s *Subject) Equal(other *Subject) bool {
	if s.ID != other.ID {
		return false
	}

	if s.UserID != other.UserID {
		return false
	}

	if s.Name != other.Name {
		return false
	}

	if len(s.Units) != len(other.Units) {
		return false
	}

	return true
}

func (s *Subject) FindUnitByID(unitID uuid.UUID) (*Unit, error) {
	for _, v := range s.Units {
		if v.ID == unitID {
			return v, nil
		}
	}

	return nil, fmt.Errorf("unit with ID %s could not be found", unitID)
}

func (s *Subject) UpdateSubjectName(newName string) error {
	if newName == "" {
		return ErrSubjectNameEmpty
	}

	capitalizedName := strings.ToUpper(string(newName[0])) + newName[1:]
	s.Name = capitalizedName
	return nil
}
