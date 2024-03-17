package domain

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrDuplicateKeyword = errors.New("keywords must be unique within a unit")

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

func (u *Unit) AddKeyword(k *Keyword) error {
	isDuplicate := u.IsDuplicateKeyword(*k)
	if isDuplicate {
		return fmt.Errorf("unit %s already contains keyword %s, %w", u.Name, k.Name, ErrDuplicateKeyword)
	}

	u.Keywords = append(u.Keywords, k)
	return nil
}

func (u *Unit) FindKeyword(keywordID uuid.UUID) (*Keyword, error) {
	for _, v := range u.Keywords {
		if v.Id == keywordID {
			return v, nil
		}
	}

	return nil, fmt.Errorf("keyword with ID: %s could not be found in unit: %s", keywordID, u.Name)
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

	if strings.ToLower(u.Name) != strings.ToLower(other.Name) {
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
