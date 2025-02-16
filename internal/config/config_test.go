package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConfig(t *testing.T) {
	// Change the working directory to the current directory
	// to load the test configuration file
	err := os.Chdir("../..")

	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	c := NewConfig()

	// check if the configuration is not nil
	assert.NotNil(t, c)

	// check if the configuration file is loaded correctly
	assert.Equal(t, c.GetPort(), 8080)

	// check if the routes are loaded correctly
	assert.NotEmpty(t, c.GetRoutes())

	// check database configuration
	assert.Equal(t, c.GetDatabase().Type, "postgres")
	assert.Equal(t, c.GetDatabase().ConnectionString, "localhost:5432")
}
