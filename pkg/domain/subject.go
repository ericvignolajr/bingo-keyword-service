package domain

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
)

var ErrSubjectNameEmpty = errors.New("subject name is empty, cannot create subject")
var ErrDuplicateUnit = errors.New("unit names must be unique within a subject")

type Subject struct {
	Id      uuid.UUID
	Name    string
	Units   []*Unit
	OwnerID uuid.UUID
}

func NewSubject(name string, ownerID uuid.UUID) (*Subject, error) {
	if name == "" {
		return nil, ErrSubjectNameEmpty
	}
	capitalizedName := strings.ToUpper(string(name[0])) + name[1:]
	return &Subject{
		Id:      uuid.New(),
		Name:    capitalizedName,
		Units:   []*Unit{},
		OwnerID: ownerID,
	}, nil
}

func (s *Subject) AddUnit(u Unit) (*Unit, error) {
	isDuplicate := s.IsDuplicateUnit(u)
	if isDuplicate {
		return nil, fmt.Errorf("%w, subject %s already contains unit %s", ErrDuplicateUnit, s.Name, u.Name)
	}

	capitalizedName := strings.ToUpper(string(u.Name[0])) + u.Name[1:]
	u.Name = capitalizedName

	s.Units = append(s.Units, &u)
	sort.SliceStable(s.Units, func(i, j int) bool {
		return strings.ToLower(s.Units[i].Name) < strings.ToLower(s.Units[j].Name)
	})
	return &u, nil
}

func (s *Subject) IsDuplicateUnit(u Unit) bool {
	for _, v := range s.Units {
		if strings.ToLower(v.Name) == strings.ToLower(u.Name) {
			return true
		}
	}

	return false
}
