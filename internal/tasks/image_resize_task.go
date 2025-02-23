package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

const (
	taskName = "image:resize" // Name of the task
)

type ImageResizeTaskPayload struct {
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileID       string `json:"file_id"`
	StoragePath  string `json:"storage_path"`
	OriginalName string `json:"original_name"`
}

type imageResizeClient struct {
	client *asynq.Client  // Client to interact with the task queue
	log    zerolog.Logger // Logger to log messages
}

type ImageResizeTask interface {
	// Name returns the name of the tas
	Enqueue(payload *ImageResizeTaskPayload) error
}

func NewImageResizeTask(c *asynq.Client, l zerolog.Logger) ImageResizeTask {
	return &imageResizeClient{
		client: c,
		log:    l,
	}
}

func (i *imageResizeClient) Enqueue(p *ImageResizeTaskPayload) error {
	// Enqueue the task with the given payload
	payload, err := json.Marshal(p)
	if err != nil {
		i.log.Error().Err(err).Msg("Failed to marshal image resize task payload for file: " + p.FileID)
		return err
	}

	// Create a new task with the given payload
	task := asynq.NewTask(taskName, payload)
	// Enqueue the task with the client
	_, err = i.client.Enqueue(task)
	if err != nil {
		i.log.Error().Err(err).Msg("Failed to enqueue image resize task for file: " + p.FileID)
		return err
	}

	return nil
}
