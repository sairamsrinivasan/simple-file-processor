package server

import (
	"os"
	"simple-file-processor/internal/config"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var (
	c *config.Config
	l zerolog.Logger
)

func TestNewRouter(t *testing.T) {
	os.Chdir("../..")
	c := config.NewConfig()
	l := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	r := NewRouter(&c, l)
	assert.NotNil(t, r)
}
