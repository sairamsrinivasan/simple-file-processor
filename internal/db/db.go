package db

import (
	"fmt"
	"simple-file-processor/internal/models"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type DB struct {
	Gdb GormDB
	Log zerolog.Logger
}

// Wrapper interface for gorm.DB
// So we can use dependency injection
// for testing purposes
type GormDB interface {
	Create(interface{}) *gorm.DB
	AutoMigrate(...interface{}) error
}

type Database interface {
	Migrate() error
	InsertFileMetadata(*models.File) error
}

// NewDB creates a new database instance with the given configuration and gorm instance
// the gorm instance is passed as an interface to allow for mocking in tests
func NewDB(gdb GormDB, l zerolog.Logger) Database {
	return &DB{Gdb: gdb, Log: l} // No need to take address of an interface
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
