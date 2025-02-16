package db

import (
	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	log zerolog.Logger
	db  *gorm.DB
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

	return DB{db: db, log: l}
}
