package server

import (
	"fmt"
	"net/http"
	"os"
	"simple-file-processor/internal/config"
	"simple-file-processor/internal/db"
	"strconv"

	"github.com/rs/zerolog"
)

type server struct {
	conf   config.Config
	router Router
	log    zerolog.Logger
	db     db.DB
}

type Server interface {
	Start() error
}

func NewServer() Server {
	c := config.NewConfig()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	l := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	r := NewRouter(&c, l)
	db := db.NewDB(c.GetConnectionString(), l)

	// Initialize the server with the given configuration
	return &server{
		conf:   &c,
		router: r,
		log:    l,
		db:     db,
	}
}

func (s *server) Start() error {
	s.log.Info().Msg("Starting server on port " + strconv.Itoa(s.conf.GetPort()))
	s.router.InitRoutes()
	return http.ListenAndServe(fmt.Sprintf(":%d", s.conf.GetPort()), s.router.GetRouter())
}
