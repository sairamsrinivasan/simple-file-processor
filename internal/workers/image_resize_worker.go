package workers

import (
	"context"
	"encoding/json"
	"simple-file-processor/internal/tasks"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

type ImageResizeWorker struct {
	// The handlers that will be used to handle the image resize task
	log zerolog.Logger
}

func NewImageResizeWorker(l zerolog.Logger) *ImageResizeWorker {
	return &ImageResizeWorker{
		log: l,
	}
}

// THe handler that will be used to handle the image resize task

func (i *ImageResizeWorker) Resize(ctx context.Context, t *asynq.Task) error {
	// Unmarshal the payload from the task
	var p *tasks.ImageResizeTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		i.log.Error().Err(err).Msg("Failed to unmarshal image resize task payload")
		return err
	}

	i.log.Info().Msg("Resizing image for file: " + p.FileID)
	return nil
}
