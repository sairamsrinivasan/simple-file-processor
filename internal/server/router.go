package server

import (
	"fmt"
	"simple-file-processor/internal/config"
	"simple-file-processor/internal/db"
	"simple-file-processor/internal/handlers"
	"simple-file-processor/internal/tasks"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type router struct {
	conf     config.Config
	handlers handlers.Handlers
	router   *mux.Router
	log      *zerolog.Logger
}

type Router interface {
	InitRoutes()
	Router() *mux.Router
}

// NewRouter initializes the router with the given configuration
func NewRouter(c config.Config, log *zerolog.Logger, db db.Database) Router {
	// Set up an async client to be used for async tasks
	// Initialize the router with the given configuration
	// and return the router instance
	return &router{
		conf:     c,
		log:      log,
		router:   mux.NewRouter(),
		handlers: handlers.NewHandlers(log, db, AsyncClient(c)),
	}
}

// AsyncClient initializes the async client
// This is used to send tasks to the async worker
func AsyncClient(c config.Config) tasks.Client {
	return tasks.NewAsyncClient(c.RedisAddress(), c.RedisDB())
}

// Initializes the routes for the server using the configuration
func (r *router) InitRoutes() {
	// Initialize routes here
	r.log.Debug().Msg("Initializing routes")
	rts := r.conf.GetRoutes()
	for _, rt := range rts {
		fmt.Println("Route: ", rt.Path, " Method: ", rt.Method)
		r.router.HandleFunc(rt.Path, r.handlers.GetHandler(rt.Handler)).Methods(rt.Method)
	}
}

// returns the router to be used in the main function
// so that it can be used to start the server
func (r *router) Router() *mux.Router {
	return r.router
}
