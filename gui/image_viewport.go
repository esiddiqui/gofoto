package gui

import (
	"fmt"
	"image"
	"image/color"

	// "image/color"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/text"

	log "github.com/sirupsen/logrus"
	"golang.org/x/image/font/basicfont"

	gofoto_image "github.com/esiddiqui/gofoto/image"
)

type imageCache struct {
	// size int           // TODO move to viewPort /// cache size, min 1 to hold the currently loaded picture
	// dat  []image.Image // TODO move to viewPort /// a cache of image.Image

	_todoFilename string
	_todoDat      image.Image
}

func (c *imageCache) put(key string, img image.Image) {
	// TODO impl this shit
	c._todoFilename = key
	c._todoDat = img
}

func (c imageCache) get(key string) (image.Image, bool) {
	if c._todoFilename == key {
		return c._todoDat, true
	}
	return nil, false
}
func (c imageCache) has(key string) bool {
	return c._todoFilename == key
}
func (c imageCache) cap() int { return 0 } //c.size }

var imgCache imageCache

type ImageViewPort struct {
	_state windowState
	// cache  imageCache

	// cacheSize: _cacheSize,           // load 1 picture at a time, 0 is invalid
	// cache:     []*image.Image{_img}, // loaded image cache

	// Path string
	// ForceRedraw     bool
	// ScaleToFit      bool
	// AngleOfRotation int

	// internal
	// lastDrawnPath string
}

// func NewImageViewPort(cacheSize int) ImageViewPort {
// 	cache := imageCache{
// 		size: cacheSize,
// 		dat:  make([]image.Image, cacheSize),
// 	}

// 	return ImageViewPort{
// 		cache: cache,
// 	}

// }

func (v *ImageViewPort) Draw(win *opengl.Window, state windowState) error {

	log.Infof("loading image to w.window")
	if v._state.equals(state) {
		log.Infof("nothing changed")
		return nil
	}

	// TODO
	// load image using gofoto or cache
	// load additional images if needed to fill the cache
	filename := state.files[state.current]

	var img image.Image
	var err error
	var ok bool

	if img, ok = imgCache.get(filename); !ok {
		log.Infof("cache miss; reading from file %v", filename)
		img, err = gofoto_image.Open(filename)
		if err != nil {
			log.Errorf("error opening image: %v\n", filename)
		}
		imgCache.put(filename, img)
	}

	// rotation
	if state.rotationAngle != 0 {
		var err error
		log.Infof("rotation %v by %v deg", filename, state.rotationAngle)
		img, err = gofoto_image.Rotate(img, state.rotationAngle)
		if err != nil {
			log.Errorf("error rotating image %v by %v deg", filename, state.rotationAngle)
		}
	}

	// scaling needed?
	scale_factor := state.scaleFactor
	if state.scaleToFit {
		scale_factor = state.effectiveScalefactor(win, img)
	}

	if scale_factor != 1.0 {
		img, err = gofoto_image.ResizeScale(img, float32(scale_factor))
		if err != nil {
			log.Errorf("error scaling image %v by %v percent", filename, scale_factor)
		}
	}

	// TODO need to ensure only viewport area is cleared
	// clear out window before drawing
	win.Clear(color.Black)

	// load a pic & make sprite from the image
	pic := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(pic, pic.Bounds())

	// relative position
	// TODO allow for scrolling if not scaled
	vec1 := win.Bounds().Center() // center of the w.win
	if !state.scaleToFit {
		vec1.X += float64(state.scrollOffsetX)
		vec1.Y += float64(state.scrollOffsetY)
		log.Infof("scroll offset %v %v", state.scrollOffsetX, state.scrollOffsetY)
	}
	sprite.Draw(win, pixel.IM.Moved(vec1))

	// TODO very expensive operation; so need a review here...
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(10, 10), basicAtlas)
	fmt.Fprintln(basicTxt, filename)
	basicTxt.Draw(win, pixel.IM)

	return nil
}
