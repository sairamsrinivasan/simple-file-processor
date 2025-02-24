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
	t.Run("DatabasePassword", func(t *testing.T) {

		t.Run("Default Password", func(t *testing.T) {
			os.Unsetenv("FILE_DATABASE_PASSWORD")
			assert.Equal(t, c.DatabasePassword(), "password")
		})

		t.Run("Set Password", func(t *testing.T) {
			os.Setenv("FILE_DATABASE_PASSWORD", "test")
			assert.Equal(t, c.DatabasePassword(), "test")
			os.Unsetenv("FILE_DATABASE_PASSWORD")
		})
	})

	// Test Getting the Database Username
	t.Run("DatabaseUsername", func(t *testing.T) {
		t.Run("Default Username", func(t *testing.T) {
			os.Unsetenv("FILE_DATABASE_USERNAME")
			assert.Equal(t, c.DatabaseUsername(), "username")
		})

		t.Run("Set Username", func(t *testing.T) {
			os.Setenv("FILE_DATABASE_USERNAME", "test")
			assert.Equal(t, c.DatabaseUsername(), "test")
			os.Unsetenv("FILE_DATABASE_USERNAME")
		})
	})

	// Test Getting the Database Port
	t.Run("DatabasePort", func(t *testing.T) {
		t.Run("Default Port", func(t *testing.T) {
			assert.Equal(t, c.DatabasePort(), 5432)
		})

		t.Run("Set Port", func(t *testing.T) {
			os.Setenv("DB_PORT", "1234")
			assert.Equal(t, c.DatabasePort(), 1234)
			os.Unsetenv("DB_PORT")
		})
	})

	// Test Getting the Service Port
	t.Run("Port", func(t *testing.T) {
		t.Run("Default Port", func(t *testing.T) {
			assert.Equal(t, c.Port(), 8080)
		})

		t.Run("Set Port", func(t *testing.T) {
			os.Setenv("APP_PORT", "1234")
			assert.Equal(t, c.Port(), 1234)
			os.Unsetenv("APP_PORT")
		})
	})

	// Test Getting the database name
	t.Run("DatabaseName", func(t *testing.T) {
		t.Run("Default Name", func(t *testing.T) {
			assert.Equal(t, c.DatabaseName(), "file_processor")
		})

		t.Run("Set Name", func(t *testing.T) {
			os.Setenv("DB_NAME", "test")
			assert.Equal(t, c.DatabaseName(), "test")
			os.Unsetenv("DB_NAME")
		})
	})

	// Test Getting the database host
	t.Run("DatabaseHost", func(t *testing.T) {
		t.Run("Default Host", func(t *testing.T) {
			assert.Equal(t, c.DatabaseHost(), "localhost")
		})

		t.Run("Set Host", func(t *testing.T) {
			os.Setenv("DB_HOST", "test")
			assert.Equal(t, c.DatabaseHost(), "test")
			os.Unsetenv("DB_HOST")
		})
	})

	// Test Getting the database type
	t.Run("DatabaseType", func(t *testing.T) {
		t.Run("Default Type", func(t *testing.T) {
			assert.Equal(t, c.DatabaseType(), "postgres")
		})

		t.Run("Set Type", func(t *testing.T) {
			os.Setenv("DB_TYPE", "test")
			assert.Equal(t, c.DatabaseType(), "test")
			os.Unsetenv("DB_TYPE")
		})
	})

	// Test Getting the connection string for the database
	t.Run("ConnectionString", func(t *testing.T) {
		assert.Equal(t, c.ConnectionString(), "postgres://username:password@localhost:5432/file_processor?sslmode=disable")
	})

	// Test Getting the Redis URL
	t.Run("RedisURL", func(t *testing.T) {
		t.Run("Default URL", func(t *testing.T) {
			assert.Equal(t, c.RedisURL(), "redis://localhost:6379/0")
		})
	})

	t.Run("RedisAddress", func(t *testing.T) {
		t.Run("Default Address", func(t *testing.T) {
			assert.Equal(t, c.RedisAddress(), "localhost:6379")
		})
	})

	t.Run("RedisDB", func(t *testing.T) {
		t.Run("Default DB", func(t *testing.T) {
			assert.Equal(t, c.RedisDB(), 0)
		})
	})
}
