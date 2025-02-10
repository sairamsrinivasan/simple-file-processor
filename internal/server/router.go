package server

import (
	"fmt"
	"simple-file-processor/internal/config"
	"simple-file-processor/internal/handlers"

	"github.com/gorilla/mux"
)

type router struct {
	conf     config.Config
	handlers handlers.Handlers
	router   *mux.Router
}

type Router interface {
	InitRoutes()
	GetRouter() *mux.Router
}

func NewRouter(c config.Config) Router {
	// Initialize the router with the given configuration
	// and return the router instance
	r := &router{
		conf: c,
	}

	r.router = mux.NewRouter()
	r.handlers = handlers.NewHandlers()
	return r
}

// Initializes the routes for the server using the configuration
func (r *router) InitRoutes() {
	// Initialize routes here
	fmt.Println("Initializing routes")
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
