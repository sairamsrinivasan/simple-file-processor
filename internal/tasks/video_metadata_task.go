package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"simple-file-processor/internal/db"
	"simple-file-processor/internal/models"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

const (
	VideoMetadataTaskType = "video:extract-metadata" // Name of the task
)

// Holds the payload for the video metadata task
type VideoMetadataTaskPayload struct {
	FileID      string
	StoragePath string
	Filename    string
}

// Constructs a client for the video metadata task
type videoMetadataHandler struct {
	db  db.Database
	ext MetadataExtractor
	log *zerolog.Logger
}

// Constructs a client for the video metadata task
func NewVideoMetadataTask(c Client, p *VideoMetadataTaskPayload, l *zerolog.Logger) (Task, error) {
	payload, err := json.Marshal(p)
	if err != nil {
		l.Error().Err(err).Msg("Failed to marshal image resize task payload for file: " + p.FileID)
		return nil, err
	}

	l.Info().Msg("Creating image resize task with payload: " + string(payload))
	return &task{
		client: c,
		log:    l,
		task:   asynq.NewTask(VideoMetadataTaskType, payload),
	}, nil
}

func NewVideoMetadataHandler(ext MetadataExtractor, db db.Database, l *zerolog.Logger) *videoMetadataHandler {
	return &videoMetadataHandler{
		db:  db,
		ext: ext,
		log: l,
	}
}

// Handles the video metadata task and resizes the image
// This will be called by the async worker
func (h *videoMetadataHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	// Unmarshal the payload
	var p VideoMetadataTaskPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		h.log.Error().Err(err).Msg("Failed to unmarshal video metadata task payload")
		return err
	}

	h.log.Info().Msgf("Processing video metadata task for file %s", p.FileID)

	// Get the file from the database
	f, err := h.db.FileByID(p.FileID)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get file by ID")
		return err
	}

	// If the file is not a video, return an error
	if !f.IsVideo() {
		h.log.Error().Msg("File is not a video")
		return fmt.Errorf("file is not a video")
	}

	// Extract the video metadata
	m, err := h.ext.ExtractVideoMetadata(f)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to extract video metadata")
		return err
	}

	// Create and add the processed output to the database
	po := processedOutput(m)
	if err := h.db.AddProcessedOutput(p.FileID, po); err != nil {
		h.log.Error().Err(err).Msg("Failed to add processed output to database")
		return err
	}

	h.log.Info().Msgf("Processed video metadata for file %s", p.FileID)
	return nil
}

func processedOutput(vm *VideoMetadata) models.ProcessedOutput {
	return models.ProcessedOutput{
		BitRate:    vm.BitRate,
		Codec:      vm.Codec,
		Duration:   vm.Duration,
		Height:     vm.Height,
		Resolution: vm.Resolution,
		Size:       vm.Size,
		Type:       "video_metadata",
		Width:      vm.Width,
	}
}
