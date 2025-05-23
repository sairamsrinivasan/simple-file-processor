package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"simple-file-processor/internal/db"
	"simple-file-processor/internal/lib"
	"simple-file-processor/internal/models"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

const (
	VideoMetadataTaskType = "video:extract-metadata" // Name of the task
	metadataExt           = "json"                   // The file extension of the metadata file
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
	ext lib.MetadataExtractor
	fs  lib.FileSystem
	log *zerolog.Logger
}

// Constructs a client for the video metadata task
func NewVideoMetadataTask(c Client, p *VideoMetadataTaskPayload, l *zerolog.Logger) (Task, error) {
	payload, err := json.Marshal(p)
	if err != nil {
		l.Error().Err(err).Msg("Failed to marshal video metadata task payload for file: " + p.FileID)
		return nil, err
	}

	l.Info().Msg("Creating video metadata task with payload: " + string(payload))
	return &task{
		client: c,
		log:    l,
		task:   asynq.NewTask(VideoMetadataTaskType, payload),
	}, nil
}

func NewVideoMetadataHandler(ext lib.MetadataExtractor, db db.Database, fs lib.FileSystem, l *zerolog.Logger) *videoMetadataHandler {
	return &videoMetadataHandler{
		db:  db,
		ext: ext,
		fs:  fs,
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
	h.log.Info().Msgf("Extracting video metadata for file %s at path %s", p.FileID, f.StoragePath)
	path := fmt.Sprintf("%s/%s", p.StoragePath, p.Filename)
	m, err := h.ext.ExtractVideoMetadata(path)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to extract video metadata")
		return err
	}

	// Create and add the processed output to the database
	po := processedOutput(f, m)
	if err := h.db.AddProcessedOutput(p.FileID, po); err != nil {
		h.log.Error().Err(err).Msg("Failed to add processed output to database")
		return err
	}

	// Generate the metadata file
	_, err = generateMetadataFile(h.fs, f, m)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to generate metadata file")
		return err
	}

	h.log.Info().Msgf("Processed video metadata for file %s and saved to %s", p.FileID, po.StoragePath)
	return nil
}

func generateMetadataFile(fs lib.FileSystem, f *models.File, vm *lib.VideoMetadata) (string, error) {
	// Create the metadata file
	metadataFile := fmt.Sprintf("%s/%s-metadata.%s", f.StoragePath, f.ID, metadataExt)
	file, err := fs.Create(metadataFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Write the metadata to the file
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(vm); err != nil {
		return "", err
	}

	return metadataFile, nil
}

func processedOutput(f *models.File, vm *lib.VideoMetadata) models.ProcessedOutput {
	return models.ProcessedOutput{
		BitRate:     vm.BitRate,
		Codec:       vm.Codec,
		Duration:    vm.Duration,
		Height:      vm.Height,
		Resolution:  vm.Resolution,
		Size:        vm.Size,
		Type:        models.VideoMetadataType,
		Width:       vm.Width,
		Name:        fmt.Sprintf("%s-%s", f.ID, "metadata"),
		Extension:   metadataExt,
		StoragePath: f.StoragePath,
	}
}
