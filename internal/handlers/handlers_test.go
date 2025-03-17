package handlers

import (
	"os"
	"testing"

	"simple-file-processor/internal/mocks/mockdb"
	"simple-file-processor/internal/mocks/mocktasks"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var (
	log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
)

func TestNewHandlers(t *testing.T) {
	db := new(mockdb.Database)
	ac := new(mocktasks.Client)
	h := NewHandlers(&log, db, ac)
	assert.NotNil(t, h)
}

// TestGetHandler tests the GetHandler function
func TestGetHandler(t *testing.T) {
	db := new(mockdb.Database)
	ac := new(mocktasks.Client)
	h := NewHandlers(&log, db, ac)
	assert.NotNil(t, h)
	handler := h.GetHandler("HealthCheckHandler")
	assert.NotNil(t, handler)
}

func TestGetHandlerNotFound(t *testing.T) {
	db := new(mockdb.Database)
	ac := new(mocktasks.Client)
	h := NewHandlers(&log, db, ac)
	assert.NotNil(t, h)
	handler := h.GetHandler("NotFoundHandler")
	assert.Nil(t, handler)
}
