package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	// Change the working directory to the current directory
	// to load the test configuration file
	err := os.Chdir("../..")

	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Load the configuration
	c := NewConfig()

	// Test Getting the Database Password
	t.Run("GetDatabasePassword", func(t *testing.T) {

		t.Run("Default Password", func(t *testing.T) {
			os.Unsetenv("FILE_DATABASE_PASSWORD")
			assert.Equal(t, c.GetDatabasePassword(), "password")
		})

		t.Run("Set Password", func(t *testing.T) {
			os.Setenv("FILE_DATABASE_PASSWORD", "test")
			assert.Equal(t, c.GetDatabasePassword(), "test")
			os.Unsetenv("FILE_DATABASE_PASSWORD")
		})
	})

	// Test Getting the Database Username
	t.Run("GetDatabaseUsername", func(t *testing.T) {
		t.Run("Default Username", func(t *testing.T) {
			os.Unsetenv("FILE_DATABASE_USERNAME")
			assert.Equal(t, c.GetDatabaseUsername(), "username")
		})

		t.Run("Set Username", func(t *testing.T) {
			os.Setenv("FILE_DATABASE_USERNAME", "test")
			assert.Equal(t, c.GetDatabaseUsername(), "test")
			os.Unsetenv("FILE_DATABASE_USERNAME")
		})
	})

	// Test Getting the Database Port
	t.Run("GetDatabasePort", func(t *testing.T) {
		t.Run("Default Port", func(t *testing.T) {
			assert.Equal(t, c.GetDatabasePort(), 5432)
		})

		t.Run("Set Port", func(t *testing.T) {
			os.Setenv("DB_PORT", "1234")
			assert.Equal(t, c.GetDatabasePort(), 1234)
			os.Unsetenv("DB_PORT")
		})
	})

	// Test Getting the Service Port
	t.Run("GetPort", func(t *testing.T) {
		t.Run("Default Port", func(t *testing.T) {
			assert.Equal(t, c.GetPort(), 8080)
		})

		t.Run("Set Port", func(t *testing.T) {
			os.Setenv("APP_PORT", "1234")
			assert.Equal(t, c.GetPort(), 1234)
			os.Unsetenv("APP_PORT")
		})
	})

	// Test Getting the database name
	t.Run("GetDatabaseName", func(t *testing.T) {
		t.Run("Default Name", func(t *testing.T) {
			assert.Equal(t, c.GetDatabaseName(), "file_processor")
		})

		t.Run("Set Name", func(t *testing.T) {
			os.Setenv("DB_NAME", "test")
			assert.Equal(t, c.GetDatabaseName(), "test")
			os.Unsetenv("DB_NAME")
		})
	})

	// Test Getting the database host
	t.Run("GetDatabaseHost", func(t *testing.T) {
		t.Run("Default Host", func(t *testing.T) {
			assert.Equal(t, c.GetDatabaseHost(), "localhost")
		})

		t.Run("Set Host", func(t *testing.T) {
			os.Setenv("DB_HOST", "test")
			assert.Equal(t, c.GetDatabaseHost(), "test")
			os.Unsetenv("DB_HOST")
		})
	})

	// Test Getting the database type
	t.Run("GetDatabaseType", func(t *testing.T) {
		t.Run("Default Type", func(t *testing.T) {
			assert.Equal(t, c.GetDatabaseType(), "postgres")
		})

		t.Run("Set Type", func(t *testing.T) {
			os.Setenv("DB_TYPE", "test")
			assert.Equal(t, c.GetDatabaseType(), "test")
			os.Unsetenv("DB_TYPE")
		})
	})

	t.Run("GetConnectionString", func(t *testing.T) {
		assert.Equal(t, c.GetConnectionString(), "postgres://username:password@localhost:5432/file_processor?sslmode=disable")
	})
}
