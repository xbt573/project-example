package database

import (
	"github.com/xbt573/project-example/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

var (
	// for database singleton
	once     = sync.Once{}
	database *gorm.DB
)

func InitDatabase(databaseUrl string) error {
	var db *gorm.DB
	var err error
	once.Do(func() {
		db, err = gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
		if err != nil {
			return
		}

		err = db.AutoMigrate(&models.TODO{})
		if err != nil {
			return
		}

		database = db
	})

	return err
}

func GetInstance() *gorm.DB {
	return database
}
