package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/health", nil)
	h := NewHandlers()
	h.GetHandler("HealthCheckHandler")(rec, r)
	assert.Equal(t, rec.Code, 200)
	assert.Equal(t, rec.Body.String(), "OK")
}
