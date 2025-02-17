package handlers

import (
	"net/http"
	"simple-file-processor/internal/db"

	"github.com/rs/zerolog"
)

type handler struct {
	handlers map[string]func(w http.ResponseWriter, r *http.Request)
	log      zerolog.Logger
	db       *db.DB
}

type Handlers interface {
	GetHandler(name string) func(w http.ResponseWriter, r *http.Request)
}

func NewHandlers(log zerolog.Logger, db *db.DB) Handlers {
	h := &handler{
		log: log,
		db:  db,
	}

	// Initialize the handlers map
	h.handlers = make(map[string]func(w http.ResponseWriter, r *http.Request))
	h.handlers["HealthCheckHandler"] = http.HandlerFunc(h.HealthCheckHandler)
	h.handlers["FileUploadHandler"] = http.HandlerFunc(h.FileUploadHandler)

	return h
}

func (h handler) GetHandler(name string) func(w http.ResponseWriter, r *http.Request) {
	h.log.Debug().Msg("Getting handler for name: " + name)
	if hand, ok := h.handlers[name]; ok {
		h.log.Debug().Msg("Handler found for name: " + name)
		return hand
	}

	return nil
}
