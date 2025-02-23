package models

import (
	"slices"
	"strings"
	"time"

	"gorm.io/gorm"
)

var imageTypes = []string{"image/jpeg", "image/png", "image/gif", "image/jpg"} // Supported image types

type File struct {
	ID                string    `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	GeneratedName     string    `json:"generated_name"`                  // e.g. file name without extension
	MimeType          string    `json:"mime_type"`                       // e.g. file mime type
	OriginalName      string    `json:"original_name"`                   // e.g. file name with extension
	Size              int64     `json:"size"`                            // e.g. file size in bytes
	Status            string    `json:"status" gorm:"default:'pending'"` // e.g. pending, processing, completed, failed
	StoragePath       string    `json:"storage_path"`                    // e.g. path where the file is stored
	Type              string    `json:"type"`                            // e.g. image, video, document, other, etc.
	UploadedExtension string    `json:"uploaded_extension"`              // e.g. file extension
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
}

// A callback that is executed before a file is created
func (f *File) BeforeCreate(tx *gorm.DB) (err error) {
	// Set the type based on the mime type
	if f.Type == "" {
		if slices.Contains(imageTypes, f.MimeType) {
			f.Type = "image"
		} else {
			f.Type = "other"
		}
	}

	return nil
}

func (f *File) IsImage() bool {
	// Check if the file is an image
	if f.Type != "image" {
		return false
	}

	// Check if the file extension is valid
	ext := strings.ToLower(f.UploadedExtension)
	return ext == "jpg" || ext == "jpeg" || ext == "png" || ext == "gif"
}
