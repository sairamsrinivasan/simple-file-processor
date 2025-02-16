package handlers

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var (
	log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
)

func TestNewHandlers(t *testing.T) {
	h := NewHandlers(log)
	assert.NotNil(t, h)
}

// TestGetHandler tests the GetHandler function
func TestGetHandler(t *testing.T) {
	h := NewHandlers(log)
	assert.NotNil(t, h)
	handler := h.GetHandler("HealthCheckHandler")
	assert.NotNil(t, handler)
}

func TestGetHandlerNotFound(t *testing.T) {
	h := NewHandlers(log)
	assert.NotNil(t, h)
	handler := h.GetHandler("NotFoundHandler")
	assert.Nil(t, handler)
}
