package handlers

import (
	"os"
	"testing"

	"simple-file-processor/internal/db"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	log = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()
)

func DBMock() *db.DB {
	mockDb, _, _ := sqlmock.New()
	p := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})

	d, _ := gorm.Open(p, &gorm.Config{})
	return &db.DB{Gdb: d, Log: log}
}

func TestNewHandlers(t *testing.T) {
	db := DBMock()
	h := NewHandlers(log, db)
	assert.NotNil(t, h)
}

// TestGetHandler tests the GetHandler function
func TestGetHandler(t *testing.T) {
	db := DBMock()
	h := NewHandlers(log, db)
	assert.NotNil(t, h)
	handler := h.GetHandler("HealthCheckHandler")
	assert.NotNil(t, handler)
}

func TestGetHandlerNotFound(t *testing.T) {
	db := DBMock()
	h := NewHandlers(log, db)
	assert.NotNil(t, h)
	handler := h.GetHandler("NotFoundHandler")
	assert.Nil(t, handler)
}
