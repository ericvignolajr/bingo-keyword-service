package domain

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

const (
	ErrSubjectNameEmpty = "subject name is empty, cannot create subject"
)

type Subject struct {
	Id      uuid.UUID
	Name    string
	Units   []*Unit
	OwnerID uuid.UUID
}

func NewSubject(name string) (*Subject, error) {
	if name == "" {
		return nil, errors.New(ErrSubjectNameEmpty)
	}
	capitalizedName := strings.ToUpper(string(name[0])) + name[1:]
	return &Subject{
		Id:    uuid.New(),
		Name:  capitalizedName,
		Units: []*Unit{},
	}, nil
}

func (s *Subject) AddUnit(u Unit) (*Unit, error) {
	isDuplicate := s.IsDuplicateUnit(u)
	if isDuplicate {
		return nil, fmt.Errorf("subject %s already contains unit %s", s.Name, u.Name)
	}

	s.Units = append(s.Units, &u)
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
