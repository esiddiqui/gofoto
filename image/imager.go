package image

import (
	"fmt"
	"image"
	_ "image/jpeg"
	"io"
	"os"
	"time"

	"github.com/nf/cr2"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	log "github.com/sirupsen/logrus"
	"golang.org/x/image/draw"
)

// Open a file & return an image (supported formats only), or an error
func Open(path string) (image.Image, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error reading image at %v due to %w", path, err)
	}
	defer file.Close()

	img, s, err := image.Decode(file)
	_ = s
	if err != nil {
		return nil, fmt.Errorf("error decoding image at %v due to %w", path, err)
	}
	return img, nil
}

// OpenCr2 decodds the contents of the file at path & returns an image, or a non-nil error
// This function uses the github.com/nf/cr2 for the underlying decode operation
func OpenCr2(path string) (image.Image, error) {

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error reading image at %v due to %w", path, err)
	}
	defer file.Close()

	img, err := cr2.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("error decoding image at %v due to %w", path, err)
	}
	return img, nil
}

// Save the image at the specified path, or returns error
func Save(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("error creating files at %v", path)
	}
	defer file.Close()

	encoder := NewJpegEncoder(90)
	return WriteTo(encoder, file, img)
}

// WriteTo encodes & write the image to a supplied io.Writer
func WriteTo(encoder Encoder, writer io.Writer, img image.Image) error {
	return encoder.Encode(writer, img)

}

// Resize the original image to the supplied size
func ResizeAbs(original image.Image, x, y int) (image.Image, error) {

	orect := original.Bounds()

	nrect := image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{x, y},
	}

	nimg := image.NewRGBA(nrect)
	draw.CatmullRom.Scale(nimg, nrect, original, orect, draw.Over, nil)
	return nimg, nil

}

// ResizeScale the original image using scaled to percent
func ResizeScale(original image.Image, percent float32) (image.Image, error) {

	orect := original.Bounds()

	nrect := image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{int(float32(orect.Dx()) * percent), int(float32(orect.Dy()) * percent)},
	}

	nimg := image.NewRGBA(nrect)
	draw.CatmullRom.Scale(nimg, nrect, original, orect, draw.Over, nil)
	return nimg, nil

}

// Rotate the image by the specified degree
func Rotate(original image.Image, deg int) (image.Image, error) {

	if deg == 0 {
		return original, nil
	}

	if deg != -90 && deg != 90 && deg != 180 {
		return nil, fmt.Errorf("invalid rotation %v", deg)
	}

	orect := original.Bounds()

	// resized rectangle for the new image...
	var nrect image.Rectangle

	// if 90 deg CW or CCW, then we swap height/width dimension
	if deg == 90 || deg == -90 {
		nrect = image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{orect.Dy(), orect.Dx()},
		}
	} else {
		nrect = image.Rectangle{
			Min: image.Point{0, 0},
			Max: image.Point{orect.Dx(), orect.Dy()},
		}
	}

	nimg := image.NewRGBA(nrect)

	var origX, origY int
	for y := 0; y < nrect.Dy(); y++ {
		for x := 0; x < nrect.Dx(); x++ {
			switch deg {
			case -90:
				origX = orect.Dx() - 1 - y
				origY = x
			case 90:
				origX = y
				origY = orect.Dy() - 1 - x
			case 180:
				origX = orect.Dx() - 1 - x
				origY = orect.Dy() - 1 - y
			}
			nimg.Set(x, y, original.At(origX, origY))
		}
	}
	return nimg, nil
}

type ImageMetadata struct {

	//Source metadata
	Location string

	// Dimensions metadata
	SizeBytes int64
	DimX      int64
	DimY      int64

	//Exif Metadata
	Camera            *string
	DateTaken         *time.Time
	Lat               *float64
	Long              *float64
	FocalLengthNumber *int64
	FocalLengthDenom  *int64
}

func Metadata(path string) (*ImageMetadata, error) {

	metadata := &ImageMetadata{
		Location: path,
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}

	metadata.SizeBytes = stat.Size()

	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return nil, fmt.Errorf("error decoding image at %v due to %w", path, err)
	}

	colorModel := cfg.ColorModel
	_ = colorModel

	metadata.DimX = int64(cfg.Width)
	metadata.DimY = int64(cfg.Height)

	f.Close() //close file

	// repon for exif data
	f, err = os.Open(path)
	if err != nil {
		return nil, err
	}

	// Optionally register camera makenote data parsing - currently Nikon and
	// Canon are supported.
	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		log.Errorf("error retrieving exif model %v", err)
		return metadata, nil
	}

	camModel, err := x.Get(exif.Model) // normally, don't ignore errors!
	if err != nil {
		log.Errorf("error retrieving exif camera informaiton %v", err)
	} else {
		camodel, _ := camModel.StringVal()
		metadata.Camera = &camodel
	}

	focal, _ := x.Get(exif.FocalLength)
	numer, denom, err := focal.Rat2(0) // retrieve first (only) rat. value
	if err != nil {
		log.Errorf("error retrieving exif focal length data %v", err)
	} else {
		metadata.FocalLengthNumber = &numer
		metadata.FocalLengthDenom = &denom
	}

	// Two convenience functions exist for date/time taken and GPS coords:
	tm, err := x.DateTime()
	if err != nil {
		log.Errorf("error retrieving exif timestamp data %v", err)
	} else {
		metadata.DateTaken = &tm
	}

	lat, long, _ := x.LatLong()
	if err != nil {
		log.Errorf("error retrieving exif gps data %v", err)
	} else {
		metadata.Lat = &lat
		metadata.Long = &long
	}

	return metadata, nil
}
