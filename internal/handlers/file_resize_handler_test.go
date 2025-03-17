package handlers_test

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"simple-file-processor/internal/handlers"
	"simple-file-processor/internal/mocks/mockdb"
	"simple-file-processor/internal/mocks/mocktasks"
	"simple-file-processor/internal/models"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFileResizeHandler(t *testing.T) {
	log := zerolog.Nop()
	var tests = []struct {
		name           string
		fileID         string
		mockDB         func(db *mockdb.Database)
		mockClient     func(client *mocktasks.Client)
		width          int
		height         int
		expectedStatus int
	}{
		{
			name:   "valid request",
			fileID: "valid-file-id",
			mockDB: func(db *mockdb.Database) {
				db.On("FileByID", "valid-file-id").Return(&models.File{
					ID:                "valid-file-id",
					Type:              "image/jpeg",
					UploadedExtension: "jpg",
				}, nil)
			},
			mockClient: func(client *mocktasks.Client) {
				client.On("Enqueue", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
			},
			width:          100,
			height:         100,
			expectedStatus: http.StatusAccepted,
		},
		{
			name:   "invalid file ID",
			fileID: "",
			mockDB: func(db *mockdb.Database) {
				// No database call needed
			},
			mockClient: func(client *mocktasks.Client) {
				// No client call needed
			},
			width:          100,
			height:         100,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:   "invalid width and height",
			fileID: "valid-file-id",
			mockDB: func(db *mockdb.Database) {
				// no database call needed
			},
			mockClient: func(client *mocktasks.Client) {
				// no client call needed
			},
			width:          -1,
			height:         -1,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:   "file not found",
			fileID: "not-found-file-id",
			mockDB: func(db *mockdb.Database) {
				db.On("FileByID", "not-found-file-id").Return(nil, fmt.Errorf("file not found"))
			},
			mockClient: func(client *mocktasks.Client) {
				// no client call needed
			},
			width:          100,
			height:         100,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:   "file is not an image",
			fileID: "not-an-image-file-id",
			mockDB: func(db *mockdb.Database) {
				db.On("FileByID", "not-an-image-file-id").Return(&models.File{
					ID:                "not-an-image-file-id",
					Type:              "text/plain",
					UploadedExtension: "txt",
				}, nil)
			},
			mockClient: func(client *mocktasks.Client) {
				// no client call needed
			},
			width:          100,
			height:         100,
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:   "failed to enqueue task",
			fileID: "valid-file-id",
			mockDB: func(db *mockdb.Database) {
				db.On("FileByID", "valid-file-id").Return(&models.File{
					ID:                "valid-file-id",
					Type:              "image/jpeg",
					UploadedExtension: "jpg",
				}, nil)
			},
			mockClient: func(client *mocktasks.Client) {
				client.On("Enqueue", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("failed to enqueue task"))
			},
			width:          100,
			height:         100,
			expectedStatus: http.StatusUnprocessableEntity,
		},
	}

	// Iterate over the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new mock database
			db := new(mockdb.Database)
			client := new(mocktasks.Client)

			// Set up the mocks
			tt.mockDB(db)
			tt.mockClient(client)

			// create a new recorder
			rec := httptest.NewRecorder()

			// create a new request with body
			body := bytes.NewBuffer([]byte(fmt.Sprintf(`{"width": %d, "height": %d}`, tt.width, tt.height)))
			req := httptest.NewRequest("PUT", "/file/"+tt.fileID+"/resize", body)
			req.Header.Set("Content-Type", "application/json")
			req = mux.SetURLVars(req, map[string]string{"id": tt.fileID})

			// Create a new handler
			handler := handlers.NewHandlers(&log, db, client).GetHandler("FileResizeHandler")

			// Call the handler
			handler(rec, req)

			// Check the status code
			assert.Equal(t, tt.expectedStatus, rec.Code)
		})
	}
}
