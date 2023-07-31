package image

import (
	"image"
	"image/jpeg"
	"io"
)

// Encoder interface defines a simple encoder behavior that
// converts an image bitmap (image.Image) to a specific compression format
type Encoder interface {
	Encode(writer io.Writer, img image.Image) error
}

// JpegEncoder encodes images to JPEG
type JpegEncoder struct {
	Quality int
}

func NewJpegEncoder(quality int) *JpegEncoder {
	return &JpegEncoder{90}
}

// Encode image into JPEG
func (j *JpegEncoder) Encode(writer io.Writer, img image.Image) error {

	quality := j.Quality

	if quality == 0 {
		quality = 50
	}
	opts := &jpeg.Options{
		Quality: 90,
	}
	return jpeg.Encode(writer, img, opts)
}

// PngEncoder encodes images to PNG
type PngEncoder struct{}

func (j *PngEncoder) Encode(writer io.Writer, img image.Image) error {
	return nil
}

// GifEncoder encodes images to GIF
type GifEncoder struct{}

func (j *GifEncoder) Encode(writer io.Writer, img image.Image) error {
	return nil
}
