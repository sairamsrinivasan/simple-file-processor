package models

import (
	"time"
)

type File struct {
	ID          string    `sql:"type:uuid;primary_key;default:uuid_generate_v4()"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Extension   string    `json:"extension"`
	StoragePath string    `json:"storage_path"`
	Size        int64     `json:"size"`
}
