package server

import (
	"fmt"
	"os"
	"os/signal"
	"simple-file-processor/internal/db"
	"simple-file-processor/internal/tasks"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

type workerServer struct {
	log   *zerolog.Logger
	rDB   int
	rAddr string
	db    db.Database
}

type WorkerServer interface {
	Start()
}

func NewWorkerServer(rAddr string, rDB int, db db.Database, log *zerolog.Logger) WorkerServer {
	return &workerServer{
		log:   log,
		rDB:   rDB,
		rAddr: rAddr,
		db:    db,
	}
}

// A background worker server that processes tasks from the task queue
// The worker server is responsible for consuming from the task queue
// and delegating the tasks to the appropriate handlers
func (ws *workerServer) Start() {
	// Initialize the worker server with the given redis address and database
	srv := asynq.NewServer(asynq.RedisClientOpt{Addr: ws.rAddr, DB: ws.rDB}, asynq.Config{
		Concurrency: 10, // Set the concurrency level
	})

	mux := asynq.NewServeMux()

	// Register the image resize handler with the task queue
	mux.Handle(tasks.ImageResizeTaskType, tasks.NewImageResizeHandler(ws.db, tasks.NewResizer(ws.log), ws.log))

	ws.log.Info().Msg("Starting worker server...")

	// Create a channel to listen for interrupt signals
	c := make(chan os.Signal, 1)

	// Listen for interrupt signals
	// This will allow the server to gracefully shutdown
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // Notify the channel when an interrupt signal is

	go func() {
		// Start the worker server with the given mux
		if err := srv.Run(mux); err != nil {
			ws.log.Fatal().Err(err).Msg("Failed to start worker")
		}
	}()

	// Wait for the signal
	<-c
	ws.log.Info().Msg("Received shutdown signal, shutting down worker server...")
	srv.Shutdown() // Shutdown the server
	fmt.Println("Server gracefully stopped")

	// Exit the process
	os.Exit(0)
}
