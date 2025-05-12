package tasks

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

const (
	TranscodingTaskType = "video:transcode" // Name of the task
)

// Holds the payload for the video transcoding task
type VideoTranscodePayload struct {
	FileID      string
	Format      string
	Quality     string
	Resolution  string
	StoragePath string
	Filename    string
}

func NewTranscodingTask(c Client, p *VideoTranscodePayload, l *zerolog.Logger) (Task, error) {
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
