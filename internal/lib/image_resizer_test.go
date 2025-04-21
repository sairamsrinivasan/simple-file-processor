package lib_test

import (
	"image"
	"image/jpeg"
	"os"
	"simple-file-processor/internal/lib"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

var (
	logger = zerolog.Nop()
)

// Verifies that the image resizer can resize an image
// and that it returns the correct output
func TestResizeImage(t *testing.T) {

	// test cases
	tests := []struct {
		name      string
		sp        string
		fn        string
		w         int
		h         int
		expectErr bool
	}{
		{
			name:      "valid resize",
			fn:        "test.jpg",
			w:         100,
			h:         100,
			expectErr: false,
		},
		{
			name:      "invalid width",
			fn:        "test.jpg",
			w:         -1,
			h:         100,
			expectErr: true,
		},
		{
			name:      "invalid height",
			fn:        "test.jpg",
			w:         100,
			h:         -1,
			expectErr: true,
		},
		{
			name:      "invalid file",
			fn:        "invalid.jpg",
			w:         100,
			h:         100,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, err := os.MkdirTemp(".", "test")
			if err != nil {
				t.Fatalf("failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(dir)

			// Create a test image file if the name is not "invalid.jpg"
			if tt.fn != "invalid.jpg" {
				createTestImage(dir, tt.fn, image.Rect(0, 0, 100, 100))
			}

			// Create a new image resizer
			resizer := lib.NewResizer(&logger)
			// Call the ResizeImage method
			output, err := resizer.ResizeImage(dir, tt.fn, tt.w, tt.h)
			if (tt.expectErr && err == nil) || (!tt.expectErr && err != nil) {
				t.Errorf("expected error: %v, got: %v", tt.expectErr, err)
			}
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NotNil(t, output)
			}
		})
	}
}

func createTestImage(dir, filename string, rect image.Rectangle) {
	img := image.NewRGBA(rect)
	f, err := os.Create(dir + "/" + filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if err := jpeg.Encode(f, img, nil); err != nil {
		panic(err)
	}
}
