package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type Subject struct {
	Id    uuid.UUID
	Name  string
	Units []*Unit
}

func NewSubject(name string) (*Subject, error) {
	return &Subject{
		Id:    uuid.New(),
		Name:  name,
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
		if v.Name == u.Name {
			return true
		}
	}

	return false
}
