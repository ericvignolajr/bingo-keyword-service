package domain

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Unit struct {
	gorm.Model
	ID           uuid.UUID `gorm:"primaryKey"`
	Name         string
	Keywords     []*Keyword
	Translations []*Translation `gorm:"polymorphic:Owner"`
	SubjectID    uuid.UUID
}

func NewUnit(name string) (*Unit, error) {
	return &Unit{
		ID:           uuid.New(),
		Name:         name,
		Keywords:     nil,
		Translations: nil,
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

func (u *Unit) Equal(other *Unit) bool {
	if u.ID != other.ID {
		return false
	}

	if u.Name != other.Name {
		return false
	}

	if u.SubjectID != other.SubjectID {
		return false
	}

	if len(u.Keywords) != len(other.Keywords) {
		return false
	}

	for i := range u.Keywords {
		if u.Keywords[i] != other.Keywords[i] {
			return false
		}
	}

	if len(u.Translations) != len(other.Translations) {
		return false
	}

	for i := range u.Translations {
		if u.Translations[i] != other.Translations[i] {
			return false
		}
	}

	return true
}
