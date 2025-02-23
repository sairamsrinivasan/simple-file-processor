package server

import (
	"fmt"
	"net/http"
	"os"
	"simple-file-processor/internal/config"
	"simple-file-processor/internal/db"
	"strconv"

	"github.com/rs/zerolog"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type server struct {
	conf   config.Config
	router Router
	log    zerolog.Logger
}

type Server interface {
	Start() error
}

func NewServer() Server {
	c := config.NewConfig()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	l := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
	gdb, err := gorm.Open(postgres.Open(c.GetConnectionString()), &gorm.Config{})
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to connect to database")
		panic(err)
	}

	db := db.NewDB(gdb, l)
	r := NewRouter(c, l, db)
	db.Migrate() // Migrate the database schema
	// Initialize the server with the given configuration
	return &server{
		conf:   c,
		router: r,
		log:    l,
	}
}

func (s *server) Start() error {
	s.log.Info().Msg("Starting server on port " + strconv.Itoa(s.conf.GetPort()))
	s.router.InitRoutes()
	return http.ListenAndServe(fmt.Sprintf(":%d", s.conf.GetPort()), s.router.GetRouter())
}
