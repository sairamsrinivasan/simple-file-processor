package tasks

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

const (
	ImageResizeTaskType = "image:resize" // Name of the task
)

// Holds the payload for the image resize task
type ImageResizePayload struct {
	Width        int
	Height       int
	FileID       string
	StoragePath  string
	OriginalName string
}

type imageResizeHandler struct {
	log zerolog.Logger
}

// Constructs a client for the image resize task
func NewImageResizeTask(c Client, l zerolog.Logger, p *ImageResizePayload) (Task, error) {
	payload, err := json.Marshal(p)
	if err != nil {
		l.Error().Err(err).Msg("Failed to marshal image resize task payload for file: " + p.FileID)
		return nil, err
	}

	l.Info().Msg("Creating image resize task with payload: " + string(payload))
	return &task{
		client: c,
		log:    l,
		task:   asynq.NewTask(ImageResizeTaskType, payload),
	}, nil
}

// Constructs a new image resize handler for the async worker
// This will handle the image resize task and ensures that the
// handler has access to the logger
func NewImageResizeHandler(l zerolog.Logger) asynq.Handler {
	return &imageResizeHandler{
		log: l,
	}
}

// Handles the image ressize task and resizes the image
// This will be called by the async worker
func (i *imageResizeHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	// Unmarshal the payload from the task
	var p *ImageResizePayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		i.log.Error().Err(err).Msg("Failed to unmarshal image resize task payload")
		return err
	}

	i.log.Info().Msg("Resizing image for file with payload: " + string(t.Payload()))
	return nil
}
