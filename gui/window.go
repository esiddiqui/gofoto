package gui

import (
	"image"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"golang.org/x/image/colornames"

	gofoto_image "github.com/esiddiqui/gofoto/image"

	log "github.com/sirupsen/logrus"
)

// NewWindow creates a new window at the supplied path
func NewWindow(width, height float64, path string) (*Window, error) {

	_cacheSize := 1

	// load files @ path
	_files := listJpegFiles(path)
	log.Debug("%v files founds in path %v", len(_files), path)

	// TODO
	// look at cache size & load accordingly,
	// current image is loaded inline, others are concurrently async after
	// load iages to cache
	var _img *image.Image
	if len(_files) > 0 {
		img, err := gofoto_image.Open(_files[0]) // load first
		if err != nil {
			log.Errorf("error opening image: %v\n", path)
		}
		_img = &img
	}

	state := &windowState{
		path:      path,                 // wording director
		files:     _files,               // if nil, then no files were loaded, sorry
		current:   0,                    // starting index
		cacheSize: _cacheSize,           // load 1 picture at a time, 0 is invalid
		cache:     []*image.Image{_img}, // loaded image cache

		scaleFactor:   1.0,   // no scaling
		scaleToFit:    false, // do not scale to fit
		scrollOffsetX: 0,
		scrollOffsetY: 0,
	}

	return &Window{
		width:    width,  // 1373,
		height:   height, // 1063,
		state:    state,  // state
		viewPort: new(ImageViewPort),
		_ref:     nil, // opengl windows refererence; will be created & saved later...
	}, nil

}

type windowState struct {
	path          string         // working directory
	files         []string       // eligible files to be loaded in the directory
	current       int            // currently pointed to file
	cacheSize     int            // cache size, min 1 to hold the currently loaded picture
	cache         []*image.Image // a cache of image.Image
	rotationAngle int            // angle of rotation
	scaleFactor   float64        // scale factor 1:100%
	scaleToFit    bool           // should ignore scaleFactor, rather scale to fit
	scrollOffsetX int
	scrollOffsetY int
}

// skip current index by `by`, where by can be +ve for forward skip
// or -ve for backward skip. If the new current is less than 0, or more
// than max (len(files)-1) then we do a circular move...
func (s *windowState) skip(by int) {
	len := len(s.files)
	max := len - 1

	s.current += by
	// cyclical move in case out of bounds...
	if s.current < 0 {
		s.current = len + by
	} else if s.current > max {
		s.current = s.current - len
	}
}

func (s *windowState) rotateOrignal() {
	s.rotationAngle = 0
}

func (s *windowState) rotateClockwise() {
	// rotate right/clockwise
	switch s.rotationAngle {
	case -90, 0, 90:
		s.rotationAngle += 90
	default: // -180
		s.rotationAngle = -90
	}
}

func (s *windowState) rotateCounterClockwise() {
	// rotate left/counter-clockwise
	switch s.rotationAngle {
	case 0, 90, 180:
		s.rotationAngle -= 90
	default: // 180
		s.rotationAngle = 180
	}
}

// effectiveScalefactor returns a scale_factor based on the window * image sizes when
// scale_to_fit flag is SET. if NOT SET, then the statically stored scale_factor
// is returned
func (s *windowState) effectiveScalefactor(win *opengl.Window, img image.Image) float64 {

	if !s.scaleToFit {
		return s.scaleFactor
	}

	scaleFactor := 1.0

	// can not calculate scale factor if
	// either of the references is nil
	if win == nil || img == nil {
		return scaleFactor
	}

	windowWidth := win.Bounds().Max.X
	windowHeight := win.Bounds().Max.Y

	max := img.Bounds().Max
	imgWidth := max.X
	imgHeight := max.Y

	if imgWidth > imgHeight {
		scaleFactor = windowWidth / float64(imgWidth)
	} else {
		scaleFactor = windowHeight / float64(imgHeight)
	}

	return scaleFactor

}

type Window struct {
	width    float64
	height   float64
	state    *windowState
	viewPort ViewPort       // the viewport abstraction
	_ref     *opengl.Window // internal pointer to opengl w.window
}

func (w *Window) CreateAndDraw() {

	// var angleOfRotation int // 0 deg
	// scaleToFit := false
	forceRedraw := false
	lastDrawnPath := ""

	// load eligible filename in the current directory & set initial pointer
	// w.index := 0
	// var files []string
	// files = listJpegFiles(w.path)
	// log.Infof("%v files founds in path %v", len(files), w.path)

	// w.window configuration
	cfg := opengl.WindowConfig{
		Title:     w.state.path,
		Bounds:    pixel.R(0, 0, w.width, w.height),
		VSync:     true,
		Resizable: true,
	}

	// create w.window
	var err error
	w._ref, err = opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var skip int
	w._ref.Clear(colornames.Red)

	// main w.window loop here
	for !w._ref.Closed() {

		skip = 1

		if w._ref.Pressed(pixel.KeyLeftShift) || w._ref.Pressed(pixel.KeyRightShift) {
			skip = 10
		}

		// with mouse-right or keypad DOWN button, we go to the NEXT picture
		if w._ref.JustPressed(pixel.MouseButtonRight) || w._ref.JustPressed(pixel.KeyDown) {
			w.state.skip(skip)
			w.state.rotateOrignal()

			// with mouse-left or keypad UP button, we go to the PREV picture
		} else if w._ref.JustPressed(pixel.MouseButtonLeft) || w._ref.JustPressed(pixel.KeyUp) {
			w.state.skip(-1 * skip)
			w.state.rotateOrignal()
			// with Space key pressed
		} else if w._ref.JustPressed(pixel.KeySpace) {
			log.Infof("Spacebar pressed, toggling scale to fit %v", !w.state.scaleToFit)
			w.state.scaleToFit = !w.state.scaleToFit
			forceRedraw = true /* requires a redraw() */

			// keypad LEFT pressed
		} else if w._ref.JustPressed(pixel.KeyLeft) {

			w.state.rotateCounterClockwise()
			forceRedraw = true
			log.Infof("rotation anagle %v", w.state.rotationAngle)

			// key pad RIGHT pressed
		} else if w._ref.JustPressed(pixel.KeyRight) {

			w.state.rotateClockwise()
			forceRedraw = true
			log.Infof("rotation anagle %v", w.state.rotationAngle)

			// Q pressed
		} else if w._ref.JustPressed(pixel.KeyQ) {
			w.state.scrollOffsetY = 0
			w.state.scrollOffsetY = 0
			forceRedraw = true
			// W pressed; viewport moves up
		} else if w._ref.Pressed(pixel.KeyW) {
			w.state.scrollOffsetY -= 15
			forceRedraw = true
			// S pressed; viewport moves down
		} else if w._ref.Pressed(pixel.KeyS) {
			w.state.scrollOffsetY += 15
			forceRedraw = true
			// A pressed; viewport
		} else if w._ref.Pressed(pixel.KeyA) {
			w.state.scrollOffsetX += 15
			forceRedraw = true
			// D pressed
		} else if w._ref.Pressed(pixel.KeyD) {
			w.state.scrollOffsetX -= 15
			forceRedraw = true
			// Esc key closes the w.windows
		} else if w._ref.JustPressed(pixel.KeyEscape) {
			w._ref.SetClosed(true)
		}

		// var xx, yy *float64 // nil, no scaling...
		// if scaleToFit {
		// 	_xx := w.win.Bounds().Max.X
		// 	_yy := w.win.Bounds().Max.Y
		// 	xx = &_xx
		// 	yy = &_yy
		// }

		// draw image on w.window only if a new path is selected, or forceRedraw is `true`
		pathToDraw := w.state.files[w.state.current] // will be drawn if we call draw
		if pathToDraw != lastDrawnPath || forceRedraw {
			w.viewPort.Draw(w._ref, *w.state)
			lastDrawnPath = pathToDraw
			forceRedraw = false // set force redraw false after draing successfully
			// w.viewPort.loadPicInw.windowScaledAndRotated(w.win, files[w.index], xx, yy, &angleOfRotation)
			// lastDrawnPath = files[w.index] // this was drawn
			// reset after the draw
		}

		w._ref.Update()

	}

	log.Info("bye bye !! w.window was closed...")
}

// processUserInput will see if a mouse or keyboard interaction
// happened & return a

// func (w *w.window) processUserInput() (any, bool) {

// 	jumpOverCnt := 1
// 	if w.win.Pressed(pixel.KeyLeftShift) || w.win.Pressed(pixel.KeyRightShift) {
// 		jumpOverCnt = 10
// 	}

// 	// with mouse-right or keypad DOWN button, we go to the NEXT picture
// 	if w.win.JustPressed(pixel.MouseButtonRight) || w.win.JustPressed(pixel.KeyDown) {
// 		// w.win.Clear(colornames.Whitesmoke)
// 		w.index += jumpOverCnt
// 		if w.index >= len(files) {
// 			w.index = 0
// 		}
// 		angleOfRotation = 0 // angle of rotation resets for new pic

// 		// with mouse-left or keypad UP button, we go to the PREV picture
// 	} else if w.win.JustPressed(pixel.MouseButtonLeft) || w.win.JustPressed(pixel.KeyUp) {
// 		// w.win.Clear(colornames.Whitesmoke)
// 		w.index -= 1
// 		if w.index < 0 {
// 			w.index = len(files) - jumpOverCnt
// 		}
// 		angleOfRotation = 0 // angle of rotation resets for new pic

// 		// with Space key pressed
// 	} else if w.win.JustPressed(pixel.KeySpace) {
// 		log.Infof("Spacebar pressed, toggling scale to fit %v", !scaleToFit)
// 		scaleToFit = !scaleToFit
// 		forceRedraw = true

// 		// keypad LEFT pressed
// 	} else if w.win.JustPressed(pixel.KeyLeft) {
// 		// rotate left
// 		switch angleOfRotation {
// 		case 0, 90, 180:
// 			angleOfRotation -= 90
// 		default: // 180
// 			angleOfRotation = 180
// 		}
// 		forceRedraw = true
// 		log.Infof("rotation anagle %v", angleOfRotation)

// 		// key pad RIGHT pressed
// 	} else if w.win.JustPressed(pixel.KeyLeft) {
// 		// rotate left
// 		switch angleOfRotation {
// 		case 0, 90, 180:
// 			angleOfRotation -= 90
// 		default: // 180
// 			angleOfRotation = 180
// 		}
// 		forceRedraw = true
// 		log.Infof("rotation anagle %v", angleOfRotation)

// 		// key pad RIGHT pressed
// 	} else if w.win.JustPressed(pixel.KeyRight) {
// 		// rotate right
// 		switch angleOfRotation {
// 		case -90, 0, 90:
// 			angleOfRotation += 90
// 		default: // -180
// 			angleOfRotation = -90
// 		}
// 		forceRedraw = true
// 		log.Infof("rotation anagle %v", angleOfRotation)

// 	// Esc key closes the w.windows
// 	} else if w.win.JustPressed(pixel.KeyEscape) {
// 		w.win.SetClosed(true)
// 	}
// }
