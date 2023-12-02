package database

import (
	"github.com/xbt573/project-example/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func migrate(db *gorm.DB) error {
	return db.AutoMigrate([]any{
		&models.TODO{},
		&models.User{},
	}...)
}

func NewDB(databaseUrl string, dialector gorm.Dialector) (*gorm.DB, error) {
	if dialector == nil {
		dialector = postgres.Open(databaseUrl)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	return db, migrate(db)
}
