package data

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Experiment struct {
	gorm.Model
	Id          string `json:"id"`
	Competitors string `json:"competitors"`
	Judge       string `json:"judge"`
	Prompt      string `json:"prompt"`
	Status      string `json:"status"`
}

func connectToDB(databasePath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	return db, err
}

func InitDB(databasePath string, schema *Experiment) (*gorm.DB, error) {
	db, err := connectToDB(databasePath)
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(schema)
	return db, nil
}
