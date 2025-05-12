package handlers

import (
	"net/http"

	"simple-file-processor/internal/models"
	"simple-file-processor/internal/tasks"

	"github.com/gorilla/mux"
)

type transcodeRequest struct {
	Format     string `json:"format"`     // The format to transcode the video to
	Quality    string `json:"quality"`    // The quality of the transcoded video
	Resolution string `json:"resolution"` // The resolution of the transcoded video
}

// A handler that triggers video transcoding tasks
func (h handler) FileTranscodeHandler(w http.ResponseWriter, r *http.Request) {
	tr := transcodeRequest{}
	if err := h.ParseRequest(r, &tr); err != nil {
		h.log.Error().Err(err).Msg("Failed to parse request body")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if tr.Format == "" || tr.Quality == "" || tr.Resolution == "" {
		h.log.Error().Msg("Missing required fields in request body")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	fid := vars["id"]
	if fid == "" {
		h.log.Error().Msg("File ID is required")
		http.Error(w, `{"error": "File id is a required path parameter"}`, http.StatusUnprocessableEntity)
		return
	}

	// Get the file from the database
	f, err := h.db.FileByID(fid)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to get file by ID")
		http.Error(w, `{"error": "File not found"}`, http.StatusNotFound)
		return
	}

	if !f.IsVideo() {
		h.log.Error().Msg("Cannot transcode file, it is not a video")
		http.Error(w, `{"error": "File is not a video"}`, http.StatusUnprocessableEntity)
		return
	}

	transcode(f, &tr, h)

	h.log.Info().Msgf("Transcoding video file %s to %s format", f.ID, tr.format)
	w.WriteHeader(http.StatusAccepted)
}

// Enqueues a video transcoding task for the given file
func transcode(f *models.File, tr *transcodeRequest, h handler) {
	// Create the payload for the video transcode task
	payload := tasks.VideoTranscodePayload{
		FileID:     f.ID,
		Format:     tr.Format,
		Quality:    tr.Quality,
		Resolution: tr.Resolution,
	}

	task, err := tasks.NewTranscodingTask(h.ac, &payload, h.log)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to create video transcode task")
		return
	}

	// Enqueue the task
	if err := task.Enqueue(); err != nil {
		h.log.Error().Err(err).Msg("Failed to enqueue video transcode task")
	}

	h.log.Info().Msgf("Video transcode task enqueued for file %s", f.ID)
}
