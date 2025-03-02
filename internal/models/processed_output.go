package models

type ProcessedOutput struct {
	ID           string `json:"id" gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	StoragePath  string `json:"storage_path"`                     // The storage path of the processed output
	OriginalName string `json:"original_name"`                    // The original name of the processed output
	Width        int    `json:"width"`                            // The width of the processed output
	Height       int    `json:"height"`                           // The height of the processed output
	Size         int64  `json:"size"`                             // The size of the processed output
	CreatedAt    string `json:"created_at" gorm:"autoCreateTime"` // The created at timestamp of the processed output
	UpdatedAt    string `json:"updated_at" gorm:"autoUpdateTime"` // The updated at timestamp of the processed output
}
