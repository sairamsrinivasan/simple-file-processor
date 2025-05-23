package models

import (
	"slices"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	imageTypes = []string{"image/jpeg", "image/png", "image/gif", "image/jpg"} // Supported image types
	videoTypes = []string{"video/mp4", "video/avi", "video/mkv", "video/mov"}  // Supported video types
)

type File struct {
	ID                string            `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	GeneratedName     string            `json:"generated_name"`                      // e.g. file name without extension
	MimeType          string            `json:"mime_type"`                           // e.g. file mime type
	ProcessedOutputs  []ProcessedOutput `json:"processed_outputs" gorm:"type:jsonb"` // e.g. processed outputs of the file, storing as jsonb
	OriginalName      string            `json:"original_name"`                       // e.g. file name with extension
	Size              int64             `json:"size"`                                // e.g. file size in bytes
	Status            string            `json:"status" gorm:"default:'pending'"`     // e.g. pending, processing, completed, failed
	StoragePath       string            `json:"storage_path"`                        // e.g. path where the file is stored
	Type              string            `json:"type"`                                // e.g. image, video, document, other, etc.
	UploadedExtension string            `json:"uploaded_extension"`                  // e.g. file extension
	CreatedAt         time.Time         `json:"created_at" gorm:"autoCreateTime"`    // e.g. file created at
	UpdatedAt         time.Time         `json:"updated_at" gorm:"autoUpdateTime"`    // e.g. file updated at
}

// A callback that is executed before a file is created
func (f *File) BeforeCreate(tx *gorm.DB) (err error) {
	// Set the type based on the mime type
	if f.Type == "" {
		if slices.Contains(imageTypes, f.MimeType) {
			f.Type = "image"
		} else if slices.Contains(videoTypes, f.MimeType) {
			f.Type = "video"
		} else {
			f.Type = "other"
		}
	}

	return nil
}

func (f *File) IsImage() bool {
	// Check if the file extension is valid
	ext := strings.ToLower(f.UploadedExtension)
	return ext == "jpg" || ext == "jpeg" || ext == "png" || ext == "gif"
}

func (f *File) IsVideo() bool {
	// Check if the file extension is valid
	ext := strings.ToLower(f.UploadedExtension)
	return ext == "mp4" || ext == "avi" || ext == "mkv" || ext == "mov"
}
