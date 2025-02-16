package db

import (
	"simple-file-processor/internal/models"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	log zerolog.Logger
	db  *gorm.DB
}

type Database interface {
	Migrate() error
}

// NewDB creates a new database instance with the given configuration
func NewDB(connStr string, l zerolog.Logger) DB {
	// Initialize the database with the given configuration
	// and return the database instance
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		l.Error().Err(err).Msg("Failed to connect to database")
		panic(err)
	}

	// Migrate the schema

	return DB{db: db, log: l}
}

func (db DB) Migrate() error {
	// Perform database migrations here
	db.log.Info().Msg("Migrating database")
	err := db.db.AutoMigrate(&models.File{})

	if err != nil {
		db.log.Error().Err(err).Msg("Failed to migrate database")
		return err
	}

	db.log.Info().Msg("Database migrated successfully")
	return nil
}
