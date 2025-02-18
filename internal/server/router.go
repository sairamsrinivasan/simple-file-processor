package server

import (
	"fmt"
	"simple-file-processor/internal/config"
	"simple-file-processor/internal/db"
	"simple-file-processor/internal/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type router struct {
	conf     config.Config
	handlers handlers.Handlers
	router   *mux.Router
	log      zerolog.Logger
}

type Router interface {
	InitRoutes()
	GetRouter() *mux.Router
}

func NewRouter(c config.Config, log zerolog.Logger, db db.Database) Router {
	// Initialize the router with the given configuration
	// and return the router instance
	return &router{
		conf:     c,
		log:      log,
		router:   mux.NewRouter(),
		handlers: handlers.NewHandlers(log, db),
	}
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
func (r *router) GetRouter() *mux.Router {
	return r.router
}
