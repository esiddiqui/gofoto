package gui

import (
	"fmt"
	"image/color"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"github.com/gopxl/pixel/v2/ext/text"

	log "github.com/sirupsen/logrus"
	"golang.org/x/image/font/basicfont"

	gofoto_image "github.com/esiddiqui/gofoto/image"
)

type ImageViewPort struct {

	// Path string
	ForceRedraw     bool
	ScaleToFit      bool
	AngleOfRotation int

	// internal
	// lastDrawnPath string
}

func (v *ImageViewPort) Draw(win *opengl.Window, state windowState) error {
	// 	return v.loadPicInWindowScaledAndRotated(win,state)
	// }

	// func (v *ImageViewPort) loadPicInWindowScaledAndRotated(win *opengl.Window, state windowState) {
	// TODO remove this
	// path string, x, y *float64, rotationAngle *int

	log.Infof("loading image to w.window")

	// TODO
	// load image using gofoto or cache
	// load additional images if needed to fill the cache
	filename := state.files[state.current]
	img, err := gofoto_image.Open(filename)
	if err != nil {
		log.Errorf("error opening image: %v\n", filename)
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

	// scale to fit
	// if x != nil || y != nil {
	// 	log.Infof("scaling %v to fit w.window size %v, %v", path, *x, *y)

	// 	max := img.Bounds().Max
	// 	imgX := max.X
	// 	imgY := max.Y

	// 	scaleFactor := 1.0
	// 	if imgX > imgY {
	// 		scaleFactor = *x / float64(imgX)
	// 	} else {
	// 		scaleFactor = *y / float64(imgY)
	// 	}
	// 	img, err = gofoto_image.ResizeScale(img, float32(scaleFactor))
	// 	if err != nil {
	// 		log.Errorf("error scaling image %v by %v percent", path, scaleFactor)
	// 	}
	// }

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
