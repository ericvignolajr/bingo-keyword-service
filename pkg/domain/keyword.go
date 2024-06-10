package domain

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"strings"

	"github.com/google/uuid"
)

type Keyword struct {
	ID               uuid.UUID `gorm:"primaryKey"`
	Name, Definition string
	Picture          image.Image    `gorm:"type:bytes"`
	Translations     []*Translation `gorm:"polymorphic:Owner"`
	UnitID           uuid.UUID
}

func NewKeyword(name, defintion string) (*Keyword, error) {
	return &Keyword{
		ID:           uuid.New(),
		Name:         name,
		Definition:   defintion,
		Translations: nil,
	}, nil
}

func (k *Keyword) Equal(other *Keyword) bool {
	if k == nil || other == nil {
		if k == nil && other == nil {
			return true
		}
		return false
	}

	if k.ID != other.ID {
		return false
	}

	if k.UnitID != other.UnitID {
		return false
	}

	if !strings.EqualFold(k.Name, other.Name) {
		return false
	}

	if !strings.EqualFold(k.Definition, other.Definition) {
		return false
	}

	if len(k.Translations) != len(other.Translations) {
		return false
	}

	for i := range k.Translations {
		if k.Translations[i] != other.Translations[i] {
			return false
		}
	}

	return true
}
