package handlers

import (
	"net/http"
	"simple-file-processor/internal/db"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

type handler struct {
	Handlers map[string]func(w http.ResponseWriter, r *http.Request)
	log      zerolog.Logger
	db       db.Database
	ac       *asynq.Client
}

type Handlers interface {
	GetHandler(name string) func(w http.ResponseWriter, r *http.Request)
}

func NewHandlers(log zerolog.Logger, db db.Database, ac *asynq.Client) Handlers {
	h := &handler{
		log: log,
		db:  db,
		ac:  ac,
	}

	// Initialize the handlers map
	// Each handler services a specific route
	// and is registered within the router
	// and defined in the config file
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
