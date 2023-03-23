package domain

import "github.com/google/uuid"

type Unit struct {
	Id          uuid.UUID
	Name        string
	Keywords    []*Keyword
	TranslateTo map[string]string
}
