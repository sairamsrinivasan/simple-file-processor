package lib

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

type videoMetadataExtractor struct {
	exec CommandExecutor
	log  *zerolog.Logger
}

// Extractor interface defines the methods that the metadata extractor should implement
type MetadataExtractor interface {
	ExtractVideoMetadata(path string) (*VideoMetadata, error)
}

// NewMetadataExtractor constructs a new metadata extractor
func NewMetadataExtractor(exec CommandExecutor, l *zerolog.Logger) MetadataExtractor {
	return &videoMetadataExtractor{
		log:  l,
		exec: exec,
	}
}

// Struct to hold the metadata
type VideoMetadata struct {
	BitRate    string `json:"bit_rate"`
	Codec      string `json:"codec"`
	Duration   string `json:"duration"`
	Height     int    `json:"height"`
	Resolution string `json:"resolution"`
	Width      int    `json:"width"`
	Size       int64  `json:"size"` // Size in bytes
	// Add more fields as needed
}

// Struct to hold the ffprobe output format
type ffprobeFormat struct {
	BitRate  string `json:"bit_rate"`
	Duration string `json:"duration"`
	Size     string `json:"size"`
}

// Struct to hold a single stream in the ffprobe output
type ffprobeStream struct {
	CodecName string `json:"codec_name"`
	CodecType string `json:"codec_type"`
	Height    int    `json:"height"`
	Width     int    `json:"width"`
}

type ffprobeOutput struct {
	Format  ffprobeFormat   `json:"format"`
	Streams []ffprobeStream `json:"streams"`
}

// ExtractMetadata extracts metadata from the file
func (e *videoMetadataExtractor) ExtractVideoMetadata(path string) (*VideoMetadata, error) {
	// Shell out to ffprobe to get the metadata
	// ffprobe -v error -print_format json -show_format -show_streams <file>
	out, err := e.exec.Command(
		"ffprobe",
		"-v", "error",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		path,
	)
	if err != nil {
		e.log.Error().Err(err).Msg("Failed to execute ffprobe for file " + path)
		return nil, err
	}

	var po ffprobeOutput
	if err := json.Unmarshal(out, &po); err != nil {
		e.log.Error().Err(err).Msg("Failed to unmarshal ffprobe output for file " + path)
		return nil, err
	}

	// Extract the metadata from the ffprobe output
	meta := videoMetadata(&po)

	// Update the file metadata in the database
	return meta, nil
}

// Extracts video metadata from the ffprobe output
func videoMetadata(po *ffprobeOutput) *VideoMetadata {
	var width, height int
	var codec, duration, bitRate string

	// Extract the metadata from the ffprobe output
	// by iterating over the streams and getting the first video stream
	// and the format
	if len(po.Streams) > 0 {
		for _, stream := range po.Streams {
			if strings.EqualFold(stream.CodecType, "video") && stream.Width > 0 && stream.Height > 0 {
				width = stream.Width
				height = stream.Height
				codec = stream.CodecName
				break
			}
		}
	}

	duration = po.Format.Duration
	bitRate = po.Format.BitRate
	return &VideoMetadata{
		BitRate:    bitRate,
		Codec:      codec,
		Duration:   duration,
		Height:     height,
		Resolution: fmt.Sprintf("%dx%d", width, height),
		Width:      width,
		Size:       toInt64(po.Format.Size),
	}
}

func toInt64(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}

	return i
}
