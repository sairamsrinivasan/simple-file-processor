package models

import (
	"database/sql/driver"
	"encoding/json"
)

type ProcessedOutput struct {
	ID           string `json:"id" gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Height       int    `json:"height"`                           // The height of the processed output
	OriginalName string `json:"original_name"`                    // The original name of the processed output
	Size         int64  `json:"size"`                             // The size of the processed output
	StoragePath  string `json:"storage_path"`                     // The storage path of the processed output
	Type         string `json:"type"`                             // The type of the processed output e.g. image, video, document, other, etc.
	Width        int    `json:"width"`                            // The width of the processed output
	CreatedAt    string `json:"created_at" gorm:"autoCreateTime"` // The created at timestamp of the processed output
	UpdatedAt    string `json:"updated_at" gorm:"autoUpdateTime"` // The updated at timestamp of the processed output
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
