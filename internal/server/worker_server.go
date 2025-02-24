package server

import (
	"simple-file-processor/internal/tasks"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

// A background worker server that processes tasks from the task queue
// The worker server is responsible for consuming from the task queue
// and delegating the tasks to the appropriate handlers
func StartWorkerServer(rAddr string, rDB int, log zerolog.Logger) {
	// Initialize the worker server with the given redis address and database
	srv := asynq.NewServer(asynq.RedisClientOpt{Addr: rAddr, DB: rDB}, asynq.Config{
		Concurrency: 10, // Set the concurrency level
	})

	mux := asynq.NewServeMux()

	// Register the image resize handler with the task queue
	mux.Handle(tasks.ImageResizeTaskType, tasks.NewImageResizeHandler(log))

	log.Info().Msg("Starting worker server...")

	// Start the worker server with the given mux
	if err := srv.Run(mux); err != nil {
		log.Fatal().Err(err).Msg("Failed to start worker")
	}
}
