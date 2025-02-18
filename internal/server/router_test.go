package server

import (
	"os"
	"simple-file-processor/internal/config"
	"simple-file-processor/internal/db"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var (
	c *config.Config
	l zerolog.Logger
	d db.Database
)

func TestNewRouter(t *testing.T) {
	os.Chdir("../..")
	c := config.NewConfig()
	l := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	d = db.NewDB(c.GetConnectionString(), l)
	r := NewRouter(&c, l, d)
	assert.NotNil(t, r)
}
