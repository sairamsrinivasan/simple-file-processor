package handlers

import (
	"net/http/httptest"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var (
	l = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
)

func TestHealthHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/health", nil)
	h := NewHandlers(l)
	h.GetHandler("HealthCheckHandler")(rec, r)
	assert.Equal(t, rec.Code, 200)
	assert.Equal(t, rec.Body.String(), "OK")
}
