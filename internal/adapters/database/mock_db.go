package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupMockDB() (*DB, error) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	mockedDb := DB{Db: db}
	mockedDb.AutoMigrate()

	return &mockedDb, nil
}
