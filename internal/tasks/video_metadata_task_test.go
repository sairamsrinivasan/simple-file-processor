package tasks_test

import (
	"bytes"
	"context"
	"errors"
	"simple-file-processor/internal/lib"
	"simple-file-processor/internal/mocks/mockdb"
	"simple-file-processor/internal/mocks/mocklib"
	"simple-file-processor/internal/mocks/mocktasks"
	"simple-file-processor/internal/models"
	"simple-file-processor/internal/tasks"
	"testing"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var log = zerolog.Nop()

// MockFile is a mock implementation of the os.File interface
// to simulate file operations in tests
// This is a simple mock that embeds bytes.Buffer, todo: Move to a separate file if needed
type MockFile struct {
	bytes.Buffer
}

func (m *MockFile) Close() error {
	return nil
}

// Test_NewVideoMetadataTask tests the NewVideoMetadataTask function
func Test_NewVideoMetadataTask(t *testing.T) {
	tests := []struct {
		name      string
		client    tasks.Client
		payload   *tasks.VideoMetadataTaskPayload
		logger    *zerolog.Logger
		expectErr bool
	}{
		{
			name:      "valid task",
			client:    new(mocktasks.Client),
			payload:   &tasks.VideoMetadataTaskPayload{FileID: "123", StoragePath: "/path/to/file", Filename: "test.mp4"},
			logger:    &log,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := tasks.NewVideoMetadataTask(tt.client, tt.payload, tt.logger)
			if (tt.expectErr && err == nil) || (!tt.expectErr && err != nil) {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
				return
			}
			if !tt.expectErr {
				assert.NotNil(t, task)
			}
		})
	}
}

// TestNewVideoMetadataHandler tests the NewVideoMetadataHandler function
func TestNewVideoMetadataHandler(t *testing.T) {
	// test cases
	tests := []struct {
		name    string
		db      *mockdb.Database
		resizer *mocklib.MetadataExtractor
		fs      *mocklib.FileSystem
		logger  *zerolog.Logger
	}{
		{
			name:    "valid handler",
			db:      new(mockdb.Database),
			resizer: new(mocklib.MetadataExtractor),
			logger:  &log,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := tasks.NewVideoMetadataHandler(tt.resizer, tt.db, tt.fs, tt.logger)
			assert.NotNil(t, handler)
		})
	}
}

// TestProcessTask tests the ProcessTask function
func TestVideoMetadataProcessTask(t *testing.T) {
	task := asynq.NewTask(tasks.VideoMetadataTaskType, []byte(`{"FileID":"123","StoragePath":"/path/to/file","Filename":"test.mp4"}`))

	tests := []struct {
		name           string
		mockDB         func(m *mockdb.Database)
		mockExtractor  func(m *mocklib.MetadataExtractor)
		mockFileSystem func(m *mocklib.FileSystem)
		task           *asynq.Task
		expectErr      bool
	}{
		{
			name: "valid task",
			mockDB: func(m *mockdb.Database) {
				fid := "123"
				m.On("FileByID", fid).Return(&models.File{ID: fid, StoragePath: "/path/to/file", OriginalName: "test.mp4", UploadedExtension: "mp4"}, nil)
				m.On("AddProcessedOutput", fid, mock.Anything).Return(nil, nil)
			},
			mockExtractor: func(m *mocklib.MetadataExtractor) {
				m.On("ExtractVideoMetadata", "/path/to/file/test.mp4").Return(&lib.VideoMetadata{}, nil)
			},
			mockFileSystem: func(m *mocklib.FileSystem) {
				m.On("Create", "/path/to/file/123-metadata.json").Return(&MockFile{}, nil)
			},
			task:      task,
			expectErr: false,
		},
		{
			name: "file is not a video",
			mockDB: func(m *mockdb.Database) {
				m.On("FileByID", "123").Return(&models.File{StoragePath: "/path/to/file", OriginalName: "test.txt", UploadedExtension: "txt"}, nil)
			},
			mockExtractor:  func(_ *mocklib.MetadataExtractor) {},
			mockFileSystem: func(m *mocklib.FileSystem) {},
			task:           task,
			expectErr:      true,
		},
		{
			name: "failed to extract video metadata",
			mockDB: func(m *mockdb.Database) {
				m.On("FileByID", "123").Return(&models.File{StoragePath: "/path/to/file", OriginalName: "test.mp4", UploadedExtension: "mp4"}, nil)
			},
			mockExtractor: func(m *mocklib.MetadataExtractor) {
				m.On("ExtractVideoMetadata", "/path/to/file/test.mp4").Return(nil, errors.New("extract error"))
			},
			mockFileSystem: func(m *mocklib.FileSystem) {},
			task:           task,
			expectErr:      true,
		},
		{
			name: "failed to add processed output",
			mockDB: func(m *mockdb.Database) {
				m.On("FileByID", "123").Return(&models.File{StoragePath: "/path/to/file", OriginalName: "test.mp4", UploadedExtension: "mp4"}, nil)
				m.On("AddProcessedOutput", "123", mock.Anything).Return(errors.New("db error"))
			},
			mockExtractor: func(m *mocklib.MetadataExtractor) {
				m.On("ExtractVideoMetadata", "/path/to/file/test.mp4").Return(&lib.VideoMetadata{}, nil)
			},
			mockFileSystem: func(m *mocklib.FileSystem) {},
			task:           task,
			expectErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := new(mockdb.Database)
			ext := new(mocklib.MetadataExtractor)
			fs := new(mocklib.FileSystem)
			handler := tasks.NewVideoMetadataHandler(ext, db, fs, &log)

			tt.mockDB(db)
			tt.mockExtractor(ext)
			tt.mockFileSystem(fs)

			err := handler.ProcessTask(context.Background(), tt.task)
			if (tt.expectErr && err == nil) || (!tt.expectErr && err != nil) {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
				return
			}

			if !tt.expectErr {
				assert.Nil(t, err)
			}
		})
	}
}
