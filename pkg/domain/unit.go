package domain

import (
	"fmt"

	"github.com/google/uuid"
)

type Unit struct {
	ID          uuid.UUID `gorm:"primaryKey"`
	Name        string
	Keywords    []*Keyword
	TranslateTo map[string]string `gorm:"type:string"`
	SubjectID   uuid.UUID
}

func NewUnit(name string) (*Unit, error) {
	return &Unit{
		ID:          uuid.New(),
		Name:        name,
		Keywords:    []*Keyword{},
		TranslateTo: map[string]string{},
	}, nil
}

func (u *Unit) AddKeyword(k Keyword) (*Keyword, error) {
	isDuplicate := u.IsDuplicateKeyword(k)
	if isDuplicate {
		return nil, fmt.Errorf("unit %s already contains keyword %s", u.Name, k.Name)
	}

	u.Keywords = append(u.Keywords, &k)
	return &k, nil
}

func (u *Unit) IsDuplicateKeyword(k Keyword) bool {
	for _, v := range u.Keywords {
		if v.Name == k.Name {
			return true
		}
	}

	return false
}
