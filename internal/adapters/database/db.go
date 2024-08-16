package database

import (
	"fmt"
	"os"
	"strconv"

	"github.com/stivo-m/vise-resume/internal/core/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	Db *gorm.DB
}

func NewDatabase() (*DB, error) {
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, err
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), port,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &DB{Db: db}, nil
}

func (db *DB) AutoMigrate() {
	db.Db.AutoMigrate(&domain.User{})
	db.Db.AutoMigrate(&domain.Password{})
	db.Db.AutoMigrate(&domain.Verifications{})
	db.Db.AutoMigrate(&domain.Token{})
}
