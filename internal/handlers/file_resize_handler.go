package handlers

import (
	"net/http"
	"simple-file-processor/internal/models"

	"simple-file-processor/internal/tasks"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
)

type fileResizeRequest struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// FileResizeHandler handles the file resize request
func (h handler) FileResizeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fid := vars["id"]
	if vars["id"] == "" {
		h.log.Error().Msg("File ID is required")
		http.Error(w, `{"error": "File id is a required path parameter"}`, http.StatusUnprocessableEntity)
		return
	}

	var req fileResizeRequest
	if err := h.parseRequest(r, &req); err != nil {
		h.log.Error().Err(err).Msg("Failed to parse file resize request")
		http.Error(w, `{"error": "Failed to parse request"}`, http.StatusBadRequest)
		return
	}

	// Validate width and height
	if req.Width <= 0 || req.Height <= 0 {
		h.log.Error().Msg("Invalid width or height values")
		http.Error(w, `{"error": "Width and height must be greater than zero"}`, http.StatusBadRequest)
		return
	}

	// Get the file from the database
	f, err := h.db.FileByID(fid)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get file by ID")
		http.Error(w, `{"error": "File not found"}`, http.StatusNotFound)
		return
	}

	// If the file is not an image, return an error
	if !f.IsImage() {
		h.log.Error().Msg("File is not an image")
		http.Error(w, `{"error": "File is not an image"}`, http.StatusUnprocessableEntity)
		return
	}

	// Create the payload for the image resize task if the file is an image
	if err := h.ResizeImage(f, req, h.log); err != nil {
		h.log.Error().Err(err).Msg("Failed to enqueue image resize task")
		http.Error(w, `{"error": "Failed to enqueue resize task"}`, http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(`{"message": "Image resize task enqueued"}`))
}

// ResizeImage enqueues the image resize task to be processed by the async worker
func (h handler) ResizeImage(f *models.File, req fileResizeRequest, log zerolog.Logger) error {
	// Enqueue the image resize task
	// This will be handled by the async worker
	// and will be processed in the background
	payload := &tasks.ImageResizePayload{
		Width:        req.Width,
		Height:       req.Height,
		FileID:       f.ID,
		StoragePath:  f.StoragePath,
		OriginalName: f.OriginalName,
	}

	t, err := tasks.NewImageResizeTask(h.ac, log, payload)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create image resize task")
		return err
	}

	// Enqueue the image resize task
	if err := t.Enqueue(); err != nil {
		log.Error().Err(err).Msg("Failed to enqueue image resize task")
		return err
	}

	// Log the image resize task
	log.Info().Str("file_id", f.ID).Msg("Image resize task enqueued")
	return nil
}
