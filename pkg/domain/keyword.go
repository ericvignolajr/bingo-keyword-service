package domain

import (
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/google/uuid"
)

type Keyword struct {
	Id               uuid.UUID
	Name, Definition string
	Picture          image.Image
	TranslateTo      map[string]Translation
}

type Translation struct {
	Name, Definition string
}

func NewKeyword(name, defintion string) (*Keyword, error) {
	return &Keyword{
		Id:          uuid.New(),
		Name:        name,
		Definition:  defintion,
		TranslateTo: map[string]Translation{},
	}, nil
}
