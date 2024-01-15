package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Translation struct {
	gorm.Model
	Language   string
	Word       string
	Definition string
	OwnerID    uuid.UUID
	OwnerType  string
}
