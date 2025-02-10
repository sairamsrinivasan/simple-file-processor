package server

import (
	"fmt"
	"net/http"
	"simple-file-processor/internal/config"
)

type server struct {
	conf   config.Config
	router Router
}

type Server interface {
	Start() error
}

func NewServer() Server {
	c := config.NewConfig()
	r := NewRouter(&c)

	return &server{
		conf:   &c,
		router: r,
	}
}

func (s *server) Start() error {
	fmt.Println("Starting server at port : ", s.conf.GetPort())
	s.router.InitRoutes()
	return http.ListenAndServe(fmt.Sprintf(":%d", s.conf.GetPort()), s.router.GetRouter())
}
