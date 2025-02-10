package handlers

import (
	"fmt"
	"net/http"
)

type handler struct {
	handlers map[string]func(w http.ResponseWriter, r *http.Request)
}

type Handlers interface {
	GetHandler(name string) func(w http.ResponseWriter, r *http.Request)
}

func NewHandlers() *handler {
	h := &handler{}
	h.handlers = make(map[string]func(w http.ResponseWriter, r *http.Request))
	h.handlers["HealthCheckHandler"] = http.HandlerFunc(h.HealthCheckHandler)
	return h
}

func (h handler) GetHandler(name string) func(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Getting handler for name: ", name)
	if handler, ok := h.handlers[name]; ok {
		fmt.Println("Handler found for name: ", name)
		return handler
	}

	return nil
}
