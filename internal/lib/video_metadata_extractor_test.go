package lib_test

import (
	"errors"
	"simple-file-processor/internal/lib"
	"simple-file-processor/internal/mocks/mocklib"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var (
	log = zerolog.Nop()
)

func TestVideoMetadataExtractor_ExtractVideoMetadata(t *testing.T) {
	tests := []struct {
		name        string
		path        string
		mockCommand func(*mocklib.CommandExecutor)
		wantErr     bool
	}{
		{
			name: "Valid video file",
			path: "tmp/test.mp4",
			mockCommand: func(m *mocklib.CommandExecutor) {
				m.On("Command", "ffprobe", "-v", "error", "-print_format", "json", "-show_format", "-show_streams", "tmp/test.mp4").Return(
					[]byte(`
						{
							"format": {
								"bit_rate": "1000000", 
								"duration": "60.0", 
								"size": "60000000"
							}, 
							"streams": [
								{
									"codec_name": "h264", 
									"codec_type": "video", 
									"height": 1080, 
									"width": 1920
								}
							]
						}`), nil)
			},
			wantErr: false,
		},
		{
			name: "Invalid video file",
			path: "tmp/invalid.mp4",
			mockCommand: func(m *mocklib.CommandExecutor) {
				m.On("Command", "ffprobe", "-v", "error", "-print_format", "json", "-show_format", "-show_streams", "tmp/invalid.mp4").Return(
					[]byte(``), errors.New("ffprobe error"))
			},
			wantErr: true,
		},
		{
			name: "json parsing error",
			path: "tmp/test.mp4",
			mockCommand: func(m *mocklib.CommandExecutor) {
				m.On("Command", "ffprobe", "-v", "error", "-print_format", "json", "-show_format", "-show_streams", "tmp/test.mp4").Return(
					[]byte(``), nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mock the command executor
			ce := mocklib.NewCommandExecutor(t)

			// set up the mock command
			tt.mockCommand(ce)

			// create the video metadata extractor
			ext := lib.NewMetadataExtractor(ce, &log)

			// call the function
			got, err := ext.ExtractVideoMetadata(tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractVideoMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && !tt.wantErr {
				t.Errorf("ExtractVideoMetadata() got = %v, want %v", got, nil)
				return
			}

			if tt.wantErr {
				assert.Nil(t, got)
				assert.Error(t, err)
				return
			}

			assert.NotNil(t, got)
		})
	}
}
