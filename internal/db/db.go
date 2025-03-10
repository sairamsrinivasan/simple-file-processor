package db

import (
	"fmt"
	"simple-file-processor/internal/models"
	"time"

	"github.com/google/uuid"

	"github.com/rs/zerolog"
)

type DB struct {
	Gdb GormDB
	Log *zerolog.Logger
}

type Database interface {
	Migrate() error
	InsertFileMetadata(*models.File) error
	AddProcessedOutput(string, models.ProcessedOutput) error
	FileByID(string) (*models.File, error)
}

// NewDB creates a new database instance with the given configuration and gorm instance
// the gorm instance is passed as an interface to allow for mocking in tests
func NewDB(gdb GormDB, l *zerolog.Logger) Database {
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

// Insert file metadata into the database
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

// Adds a processed output to the file
func (db DB) AddProcessedOutput(fid string, po models.ProcessedOutput) error {
	// Adding a processed output to a file does not create a new record in the database
	// Instead, it updates the existing file record with the new processed output
	// Set the ID and timestamps for the processed output
	if po.ID == uuid.Nil {
		db.Log.Info().Msg("Setting ID for processed output")
		po.ID = uuid.New()
	}

	po.CreatedAt = time.Now()
	po.UpdatedAt = time.Now()

	// Add the processed output to the file
	db.Log.Info().Msg(fmt.Sprintf("Adding processed output to file: %s", fid))
	f, err := db.FileByID(fid)
	if err != nil {
		return err
	}

	f.ProcessedOutputs = append(f.ProcessedOutputs, po)
	if err := db.Gdb.Model(f).Update("processed_outputs", f.ProcessedOutputs).Error; err != nil {
		db.Log.Error().Err(err).Msg("Failed to add processed output to file")
		return err
	}

	db.Log.Info().Msg(fmt.Sprintf("Processed output added to file: %s", f.OriginalName))
	return nil
}

// GetFileByID returns the file with the given ID
func (db DB) FileByID(id string) (*models.File, error) {
	db.Log.Info().Msg(fmt.Sprintf("Getting file with ID: %s", id))
	f := &models.File{}
	if err := db.Gdb.Model(f).First(f, "id = ?", id).Error; err != nil {
		db.Log.Error().Err(err).Msg("Failed to get file by ID")
		return nil, err
	}

	db.Log.Info().Msg(fmt.Sprintf("File with ID: %s found", id))
	return f, nil
}
