package domain

import (
	"image"
	_ "image/jpeg"
	_ "image/png"

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
