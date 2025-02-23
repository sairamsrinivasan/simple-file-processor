package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

const (
	ImageResizeTaskType = "image:resize" // Name of the task
)

// Holds the payload for the image resize task
type ImageResizeTaskPayload struct {
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileID       string `json:"file_id"`
	StoragePath  string `json:"storage_path"`
	OriginalName string `json:"original_name"`
}

// The client that will be used to enqueue the image resize task
type imageResizeClient struct {
	client *asynq.Client  // Client to interact with the task queue
	log    zerolog.Logger // Logger to log messages
}

// ImageResizeTask interface defines the methods that the image resize task client should implement
type ImageResizeTask interface {
	// Name returns the name of the tas
	Enqueue(payload *ImageResizeTaskPayload) error
}

// Constructs a client for the image resize task
func NewImageResizeTask(c *asynq.Client, l zerolog.Logger) ImageResizeTask {
	return &imageResizeClient{
		client: c,
		log:    l,
	}
}

// Enqueues the image resize task with the given payload
func (i *imageResizeClient) Enqueue(p *ImageResizeTaskPayload) error {
	// Enqueue the task with the given payload
	payload, err := json.Marshal(p)
	if err != nil {
		i.log.Error().Err(err).Msg("Failed to marshal image resize task payload for file: " + p.FileID)
		return err
	}

	// Create a new task with the given payload
	task := asynq.NewTask(ImageResizeTaskType, payload)
	// Enqueue the task with the client
	_, err = i.client.Enqueue(task)
	if err != nil {
		i.log.Error().Err(err).Msg("Failed to enqueue image resize task for file: " + p.FileID)
		return err
	}

	return nil
}
