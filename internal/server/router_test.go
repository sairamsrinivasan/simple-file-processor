package server

import (
	"os"
	"simple-file-processor/internal/config"
	"simple-file-processor/internal/db"
	"testing"

	"simple-file-processor/internal/mocks"

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
	gdb := new(mocks.GormDB)
	c := config.NewConfig()
	l := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	db := db.NewDB(gdb, l)
	r := NewRouter(&c, l, db)
	assert.NotNil(t, r)
}
