package server

import (
	"os"
	"simple-file-processor/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	c *config.Config
)

func TestNewRouter(t *testing.T) {
	os.Chdir("../..")
	c := config.NewConfig()
	r := NewRouter(&c)
	assert.NotNil(t, r)
}
