package handlers

import (
	"encoding/json"
	"net/http"
	"simple-file-processor/internal/db"
	"simple-file-processor/internal/tasks"

	"github.com/rs/zerolog"
)

type handler struct {
	Handlers map[string]func(w http.ResponseWriter, r *http.Request)
	log      *zerolog.Logger
	db       db.Database
	ac       tasks.Client
}

type Handlers interface {
	GetHandler(name string) func(w http.ResponseWriter, r *http.Request)
}

// Configures handlers for the server
func NewHandlers(log *zerolog.Logger, db db.Database, ac tasks.Client) Handlers {
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
	h.Handlers["FileResizeHandler"] = http.HandlerFunc(h.FileResizeHandler)
	h.Handlers["FileTranscodeHandler"] = http.HandlerFunc(h.FileTranscodeHandler)
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

// parseRequest parses the request body into the given struct
func (h handler) ParseRequest(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	return decoder.Decode(v)
}
