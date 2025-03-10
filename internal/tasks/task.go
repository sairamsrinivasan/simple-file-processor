package tasks

import (
	"time"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

// The client that will be used to enqueue the image resize task
type task struct {
	client Client          // Client to interact with the task queue
	log    *zerolog.Logger // Logger to log messages
	task   *asynq.Task     // Task to be enqueued
}

// ImageResizeTask interface defines the methods that the image resize task client should implement
type Task interface {
	// Name returns the name of the task
	Enqueue() error // Enqueues the task with the given payload
}

// Enqueues the image resize task with the given payload
func (i *task) Enqueue() error {
	// Enqueue the task with the given payload
	_, err := i.client.Enqueue(i.task, asynq.MaxRetry(3), asynq.Timeout(60*time.Second))
	if err != nil {
		i.log.Error().Err(err).Msg("Failed to enqueue task with payload: " + string(i.task.Payload()))
		return err
	}

	i.log.Info().Msg("Enqueued task with payload: " + string(i.task.Payload()))
	return nil
}
