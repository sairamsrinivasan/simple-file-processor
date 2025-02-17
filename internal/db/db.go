package db

import (
	"fmt"
	"simple-file-processor/internal/models"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	Gdb *gorm.DB
	Log zerolog.Logger
}

type Database interface {
	Migrate() error
	InsertFileMetadata(*models.File)
}

// NewDB creates a new database instance with the given configuration
func NewDB(connStr string, l zerolog.Logger) *DB {
	// Initialize the database with the given configuration
	// and return the database instance
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to connect to database")
	}

	// Migrate the schema
	return &DB{Gdb: db, Log: l}
}

func (db DB) Migrate() error {
	// Perform database migrations here
	db.Log.Info().Msg("Migrating database")
	err := db.Gdb.AutoMigrate(&models.File{})

	if err != nil {
		db.Log.Error().Err(err).Msg("Failed to migrate database")
		return err
	}

	db.Log.Info().Msg("Database migrated successfully")
	return nil
}

func (db DB) InsertFileMetadata(f *models.File) error {

	// Insert the file content into the database
	db.Log.Info().Msg(fmt.Sprintf("Inserting file metadata into the database: %s", f.OriginalName))
	if err := db.Gdb.Create(f).Error; err != nil {
		db.Log.Error().Err(err).Msg("Failed to insert file content into the database")
		return err
	}

	db.Log.Info().Msg(fmt.Sprintf("File metadata inserted into the database: %s", f.OriginalName))
	return nil
}
