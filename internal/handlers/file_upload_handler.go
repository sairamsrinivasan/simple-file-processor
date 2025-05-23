package handlers

import (
	"encoding/json"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"simple-file-processor/internal/models"
	"simple-file-processor/internal/tasks"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

var uploadBase = "uploads"

// FileUploadHandler handles the file upload request
func (h handler) FileUploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		http.Error(w, "File is too large", http.StatusRequestEntityTooLarge)
		return
	}

	// Get the file from the form data
	f, inf, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	defer f.Close()
	// Create the upload directory if it doesn't exist
	if err := os.MkdirAll(uploadBase, os.ModePerm); err != nil {
		h.log.Error().Err(err).Msg("Failed to create upload directory")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Generate unique id for the file and name
	id := uuid.New().String() // construct unique id for the file to be stored in the database and on the file system
	ext := filepath.Ext(inf.Filename)
	var tExt string
	if len(ext) > 1 {
		tExt = ext[1:]
	} else {
		tExt = "unknown" // if no extension is provided
	}

	gn := id + "_" + inf.Filename                          // construct unique name for the file to be stored on the file system
	up := filepath.Join(uploadBase, filepath.Join(id, gn)) // construct unique path for the file to be stored on the file system
	sp := filepath.Join(uploadBase, id)
	mt := mime.TypeByExtension(ext)
	if mt == "" {
		mt = "application/octet-stream" // default mime type
	}

	// Create the upload directory for the file
	if err := os.MkdirAll(filepath.Join(uploadBase, id), os.ModePerm); err != nil {
		h.log.Error().Err(err).Msg("Failed to create upload directory")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Create the file on the "server" (file system)
	if err := CreateFile(up, f, h.log); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Track upload info for database
	file := &models.File{
		ID:                id,
		GeneratedName:     gn,
		MimeType:          mt,
		OriginalName:      inf.Filename,
		Size:              inf.Size,
		StoragePath:       sp,
		UploadedExtension: tExt,
	}

	// Insert the file metadata info into the database
	if err := h.db.InsertFileMetadata(file); err != nil {
		h.log.Error().Err(err).Msg("Failed to insert file content into the database")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return success response
	Success(w, file)

	// Generate metadata for the file
	generateVideoMetadata(h, file)

	// Log the file upload
	h.log.Info().Str("file_id", id).
		Str("file_name", inf.Filename).
		Str("stored_path", sp).
		Msg("File uploaded successfully")

}

func generateVideoMetadata(h handler, f *models.File) {
	if !f.IsVideo() {
		return
	}

	// Create a task to generate metadata for the file
	p := &tasks.VideoMetadataTaskPayload{
		FileID:      f.ID,
		StoragePath: f.StoragePath,
		Filename:    f.GeneratedName,
	}

	// Create a new video metadata task
	task, err := tasks.NewVideoMetadataTask(h.ac, p, h.log)
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to create video metadata task")
		return
	}

	// Enqueue the task
	err = task.Enqueue()
	if err != nil {
		h.log.Error().Err(err).Msg("Failed to enqueue video metadata task")
		return
	}

	h.log.Info().Msg("Video metadata task enqueued successfully for file: " + f.ID)
}

func Success(w http.ResponseWriter, f *models.File) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	resp, _ := json.Marshal(f)
	w.Write(resp)
}

func CreateFile(path string, f io.Reader, log *zerolog.Logger) error {
	// Create the file on the "server" (file system)
	dst, err := os.Create(path)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create file")
		return err
	}
	defer dst.Close() // close the file

	// Copy the uploaded file to the destination file
	if _, err := io.Copy(dst, f); err != nil {
		log.Error().Err(err).Msg("Failed to copy file")
		return err
	}

	return nil
}
