package data

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func connectToDB(databasePath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	return db, err
}

func InitDB(databasePath string, schemas ...any) (*gorm.DB, error) {
	db, err := connectToDB(databasePath)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(schemas...)
	return db, nil
}
