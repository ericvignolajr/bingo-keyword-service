package sql

import (
	"log"

	"github.com/ericvignolajr/bingo-keyword-service/pkg/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func init() {
	dbConn, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		FullSaveAssociations: true,
		Logger:               logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("could not open connection to sql database")
	}
	err = dbConn.Migrator().AutoMigrate(
		domain.User{},
		domain.Subject{},
		domain.Unit{},
		domain.Translation{},
		domain.Keyword{},
	)
	if err != nil {
		log.Fatal("could not migrate database models")
	}
	db = dbConn
}
