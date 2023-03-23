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
