package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	VideoMetadataType = "video_metadata" // The type of the processed output
	ResizedImageType  = "resized_image"  // The type of the resized image
)

type ProcessedOutput struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"` // The unique identifier of the processed output
	BitRate     string    `json:"bit_rate"`                            // The bit rate of the processed output
	Codec       string    `json:"codec"`                               // The codec of the processed output
	Duration    string    `json:"duration"`                            // The duration of the processed output
	Extension   string    `json:"extension"`                           // The file extension of the processed output
	Format      string    `json:"format"`                              // The format of the processed output
	Height      int       `json:"height"`                              // The height of the processed output
	Name        string    `json:"name"`                                // The name of the processed output
	Resolution  string    `json:"resolution"`                          // The resolution of the processed output
	Size        int64     `json:"size"`                                // The size of the processed output in bytes
	StoragePath string    `json:"storage_path"`                        // The storage path of the processed output
	Type        string    `json:"type"`                                // The type of the processed output e.g. image, video, document, other, etc.
	Width       int       `json:"width"`                               // The width of the processed output
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`    // The created at timestamp of the processed output
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`    // The updated at timestamp of the processed output
}

// Value implements the driver.Valuer interface for JSONB storage
func (po ProcessedOutput) Value() (driver.Value, error) {
	return json.Marshal(po)
}

// Scan implements the sql.Scanner interface for JSONB retrieval
func (po *ProcessedOutput) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		panic("Failed to scan processed output")
	}
	return json.Unmarshal(b, po)
}
