package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewHandlers(t *testing.T) {
	h := NewHandlers()
	assert.NotNil(t, h)
}

// TestGetHandler tests the GetHandler function
func TestGetHandler(t *testing.T) {
	h := NewHandlers()
	assert.NotNil(t, h)
	handler := h.GetHandler("HealthCheckHandler")
	assert.NotNil(t, handler)
}

func TestGetHandlerNotFound(t *testing.T) {
	h := NewHandlers()
	assert.NotNil(t, h)
	handler := h.GetHandler("NotFoundHandler")
	assert.Nil(t, handler)
}
