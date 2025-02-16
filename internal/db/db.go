package db

import (
	"simple-file-processor/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type database struct {
	db *gorm.DB
}

// NewDB creates a new database instance with the given configuration
func NewDB(conf config.Config) database {
	// Initialize the database with the given configuration
	// and return the database instance
	db, err := gorm.Open(postgres.Open(ConnectionString(conf)), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return database{db: db}
}

// ConnectionString returns the connection string for the database
func ConnectionString(conf config.Config) string {
	// Construct the database connection string here
	return conf.GetDatabaseType() +
		"://" + conf.GetDatabaseUsername() + ":" + conf.GetDatabasePassword() +
		"@" + conf.GetDatabaseHost() + ":" + string(conf.GetDatabasePort()) + "/" + conf.GetDatabaseName() + "?sslmode=disable"
}
