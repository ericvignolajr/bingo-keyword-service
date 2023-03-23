package domain

import "github.com/google/uuid"

type Subject struct {
	Id    uuid.UUID
	Name  string
	Units []*Unit
}
