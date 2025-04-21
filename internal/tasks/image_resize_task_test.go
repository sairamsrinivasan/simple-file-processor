package tasks_test

import (
	"context"
	"simple-file-processor/internal/db"
	"simple-file-processor/internal/lib"
	"simple-file-processor/internal/mocks/mockdb"
	"simple-file-processor/internal/mocks/mocktasks"
	"simple-file-processor/internal/models"
	"simple-file-processor/internal/tasks"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestNewImageResizeTask tests the NewImageResizeTask function
func TestNewImageResizeTask(t *testing.T) {
	// test cases
	tests := []struct {
		name      string
		client    tasks.Client
		payload   *tasks.ImageResizePayload
		logger    *zerolog.Logger
		expectErr bool
	}{
		{
			name:      "valid task",
			client:    new(mocktasks.Client),
			payload:   &tasks.ImageResizePayload{Width: 100, Height: 100, FileID: "123", StoragePath: "/path/to/file", Filename: "test.jpg"},
			logger:    &zerolog.Logger{},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := tasks.NewImageResizeTask(tt.client, tt.payload, tt.logger)
			if (tt.expectErr && err == nil) || (!tt.expectErr && err != nil) {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if tt.expectErr {
				assert.Nil(t, task)
				assert.Error(t, err)
			} else {
				assert.NotNil(t, task)
			}
		})
	}
}

// TestNewImageResizeHandler tests the NewImageResizeHandler function
func TestNewImageResizeHandler(t *testing.T) {
	// test cases
	tests := []struct {
		name    string
		db      db.Database
		resizer lib.Resizer
		log     *zerolog.Logger
	}{
		{
			name:    "valid handler",
			db:      new(mockdb.Database),
			resizer: new(mocktasks.Resizer),
			log:     &zerolog.Logger{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := tasks.NewImageResizeHandler(tt.db, tt.resizer, tt.log)
			assert.NotNil(t, handler)
		})
	}
}

// Tests for the ProcessTask function
func TestProcessTask(t *testing.T) {
	log := zerolog.Nop()
	task := asynq.NewTask(tasks.ImageResizeTaskType, []byte(`{"Width":100,"Height":100,"FileID":"123","StoragePath":"/path/to/file","Filename":"test.jpg"}`))
	// test cases
	tests := []struct {
		name        string
		task        *asynq.Task
		mockDB      func(m *mockdb.Database)   // Function to set up mock behavior for the database
		mockResizer func(m *mocktasks.Resizer) // Function to set up mock behavior for the image resizer
		expectErr   bool                       // Whether we expect an error
	}{
		{
			name: "valid task",
			task: task,
			mockDB: func(m *mockdb.Database) {
				m.On("AddProcessedOutput", mock.Anything, mock.Anything).Return(nil)
			},
			mockResizer: func(m *mocktasks.Resizer) {
				m.On("ResizeImage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(models.ProcessedOutput{}, nil)
			},
			expectErr: false,
		},
		{
			name: "resize error",
			task: task,
			mockDB: func(m *mockdb.Database) {
				// No database interaction expected
			},
			mockResizer: func(m *mocktasks.Resizer) {
				m.On("ResizeImage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(models.ProcessedOutput{}, assert.AnError)
			},
			expectErr: true,
		},
		{
			name: "error adding processed output",
			task: task,
			mockDB: func(m *mockdb.Database) {
				m.On("AddProcessedOutput", mock.Anything, mock.Anything).Return(assert.AnError)
			},
			mockResizer: func(m *mocktasks.Resizer) {
				m.On("ResizeImage", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(models.ProcessedOutput{}, nil)
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockDB := new(mockdb.Database)
			mockResizer := new(mocktasks.Resizer)
			// Mock image resizer
			tt.mockDB(mockDB)
			tt.mockResizer(mockResizer)
			// Create the image resize handler with the mocked DB
			handler := tasks.NewImageResizeHandler(mockDB, mockResizer, &log)

			// Run the function
			err := handler.ProcessTask(context.Background(), tt.task)

			// Check if the error matches the expected outcome
			if (err != nil) != tt.expectErr {
				t.Errorf("ProcessTask() error = %v, expectedErr %v", err, tt.expectErr)
			}

			// Verify mock expectations
			mockDB.AssertExpectations(t)
			mockResizer.AssertExpectations(t)
			// Check if the task was processed correctly
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
