package handlers

import (
	"net/http"

	"github.com/rs/zerolog"
)

type handler struct {
	handlers map[string]func(w http.ResponseWriter, r *http.Request)
	log      zerolog.Logger
}

type Handlers interface {
	GetHandler(name string) func(w http.ResponseWriter, r *http.Request)
}

func NewHandlers(log zerolog.Logger) *handler {
	h := &handler{
		log: log,
	}

	// Initialize the handlers map
	h.handlers = make(map[string]func(w http.ResponseWriter, r *http.Request))
	h.handlers["HealthCheckHandler"] = http.HandlerFunc(h.HealthCheckHandler)
	h.handlers["FileUploadHandler"] = http.HandlerFunc(h.FileUploadHandler)

	return h
}

func (h handler) GetHandler(name string) func(w http.ResponseWriter, r *http.Request) {
	h.log.Debug().Msg("Getting handler for name: " + name)
	if handler, ok := h.handlers[name]; ok {
		h.log.Debug().Msg("Handler found for name: " + name)
		return handler
	}

	return nil
}
