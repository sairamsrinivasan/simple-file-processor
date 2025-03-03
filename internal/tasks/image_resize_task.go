package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"simple-file-processor/internal/db"
	"simple-file-processor/internal/models"

	"github.com/hibiken/asynq"
	"github.com/nfnt/resize"
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
	db  db.Database
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
func NewImageResizeHandler(db db.Database, l zerolog.Logger) *imageResizeHandler {
	return &imageResizeHandler{
		db:  db,
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

	po, err := resizeImage(p.StoragePath, p.OriginalName, p.Width, p.Height, i.log)
	if err != nil {
		i.log.Error().Err(err).Msg("Failed to resize image for file with payload: " + string(t.Payload()))
		return err
	}

	if err != nil {
		i.log.Error().Err(err).Msg("Failed to get file by ID: " + p.FileID)
	}

	// Insert the processed output into the database
	if err := i.db.AddProcessedOutput(p.FileID, po); err != nil {
		i.log.Error().Err(err).Msg(fmt.Sprintf("Failed to add processed output %s to file: %s", po.OriginalName, p.FileID))
	}

	i.log.Info().Msg(fmt.Sprintf("Added processed output %s to file: %s", po.OriginalName, p.FileID))
	return nil
}

// Resizes the image with the given payload
func resizeImage(sp string, on string, w, h int, l zerolog.Logger) (models.ProcessedOutput, error) {
	in, err := buildImage(sp, on)
	if err != nil {
		l.Error().Err(err).Msg(fmt.Sprintf("Failed to open image %s at storage path %s", on, sp))
		return models.ProcessedOutput{}, err
	}

	// Resize the image
	out := resize.Resize(uint(w), uint(h), in, resize.Lanczos3)

	// Create the output file
	of := fmt.Sprintf("%s/%s", sp, "resized_"+on)
	f, err := os.Create(of)
	if err != nil {
		l.Error().Err(err).Msg(fmt.Sprintf("Failed to create resized image %s at storage path: %s", of, sp))
	}
	defer f.Close()

	// Encode the image to the output file
	if err := jpeg.Encode(f, out, nil); err != nil {
		l.Error().Err(err).Msg(fmt.Sprintf("Failed to encode resized image %s at storage path: %s", of, sp))
		return models.ProcessedOutput{}, err
	}

	// Get the size of the output file
	fi, err := os.Stat(of)

	// Create the processed output
	po := models.ProcessedOutput{
		StoragePath:  sp,
		OriginalName: fi.Name(),
		Width:        w,
		Height:       h,
		Type:         "image",
		Size:         fi.Size(),
	}

	l.Info().Msg(fmt.Sprintf("Resized image %s at storage path: %s", of, sp))
	return po, nil
}

// obtains an image from the storage path
// and returns the image object
func buildImage(storagePath, originalName string) (image.Image, error) {
	p := fmt.Sprintf("%s/%s", storagePath, originalName)
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}
