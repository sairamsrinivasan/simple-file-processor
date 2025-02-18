package handlers

import (
	"net/http"
	"simple-file-processor/internal/db"

	"github.com/rs/zerolog"
)

type handler struct {
	Handlers map[string]func(w http.ResponseWriter, r *http.Request)
	log      zerolog.Logger
	db       db.Database
}

type Handlers interface {
	GetHandler(name string) func(w http.ResponseWriter, r *http.Request)
}

func NewHandlers(log zerolog.Logger, db db.Database) Handlers {
	h := &handler{
		log: log,
		db:  db,
	}

	// Initialize the handlers map
	h.Handlers = make(map[string]func(w http.ResponseWriter, r *http.Request))
	h.Handlers["HealthCheckHandler"] = http.HandlerFunc(h.HealthCheckHandler)
	h.Handlers["FileUploadHandler"] = http.HandlerFunc(h.FileUploadHandler)

	return h
}

func (h handler) GetHandler(name string) func(w http.ResponseWriter, r *http.Request) {
	h.log.Debug().Msg("Getting handler for name: " + name)
	if hand, ok := h.Handlers[name]; ok {
		h.log.Debug().Msg("Handler found for name: " + name)
		return hand
	}

	return nil
}
