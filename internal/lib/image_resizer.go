package lib

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"simple-file-processor/internal/models"

	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"github.com/rs/zerolog"
)

type imageResizer struct {
	log *zerolog.Logger
}

type Resizer interface {
	ResizeImage(sp string, fn string, w, h int) (models.ProcessedOutput, error)
}

func NewResizer(l *zerolog.Logger) Resizer {
	return &imageResizer{
		log: l,
	}
}

// Resizes the image with the given payload
func (r *imageResizer) ResizeImage(sp string, fn string, w, h int) (models.ProcessedOutput, error) {
	// Validate the input parameters
	if w <= 0 || h <= 0 || fn == "" {
		r.log.Error().Msg(fmt.Sprintf("Invalid width or height for image %s at storage path %s", fn, sp))
		return models.ProcessedOutput{}, fmt.Errorf("invalid width or height")
	}

	// Open the image file
	f, err := os.Open(fmt.Sprintf("%s/%s", sp, fn))
	if err != nil {
		return models.ProcessedOutput{}, err
	}
	defer f.Close()

	// Decode the image
	img, err := decode(f)
	if err != nil {
		r.log.Error().Err(err).Msg(fmt.Sprintf("Failed to open image %s at storage path %s", fn, sp))
		return models.ProcessedOutput{}, err
	}

	// Resize the image
	out := resize.Resize(uint(w), uint(h), img, resize.Lanczos3)

	// Create the output file with a unique ID
	poid := uuid.New()
	ofp := fmt.Sprintf("%s/%s", sp, "resized_"+poid.String()+filepath.Ext(f.Name()))
	of, err := os.Create(ofp)
	if err != nil {
		r.log.Error().Err(err).Msg(fmt.Sprintf("Failed to create resized image %s at storage path: %s", ofp, sp))
		return models.ProcessedOutput{}, err
	}
	defer of.Close()

	// Encode the image to the output file
	if err := jpeg.Encode(of, out, nil); err != nil {
		r.log.Error().Err(err).Msg(fmt.Sprintf("Failed to encode resized image %s at storage path: %s", ofp, sp))
		return models.ProcessedOutput{}, err
	}

	// Create the processed output
	fi, _ := os.Stat(ofp)
	po := models.ProcessedOutput{
		ID:          poid,
		StoragePath: sp,
		Name:        fi.Name(),
		Width:       w,
		Height:      h,
		Type:        models.ResizedImageType,
		Extension:   filepath.Ext(fi.Name()),
		Size:        fi.Size(),
	}

	r.log.Info().Msg(fmt.Sprintf("Resized image %s at storage path: %s", ofp, sp))
	return po, nil
}

// obtains an image from the storage path
// and returns the image object
func decode(f *os.File) (image.Image, error) {
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}
