package gui

import (
	"fmt"
	"image"
	"image/color"

	"github.com/gopxl/pixel/v2"
	"github.com/gopxl/pixel/v2/backends/opengl"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"

	"github.com/gopxl/pixel/v2/ext/text"

	gofoto_image "github.com/esiddiqui/gofoto/image"

	log "github.com/sirupsen/logrus"
)

type Window struct {
	X, Y     float64
	ViewPort *ViewPort
	Path     string
	// TODO Check
	ForceRedraw bool
}

func (w *Window) Draw() {

	var files []string

	var angleOfRotation int // 0 deg
	scaleToFit := false
	forceRedraw := false
	lastDrawnPath := ""

	// load image
	index := 0
	files = listJpegFiles(w.Path)
	log.Infof("%v files founds in path %v", len(files), w.Path)

	// window configuration
	cfg := opengl.WindowConfig{
		Title:     w.Path,
		Bounds:    pixel.R(0, 0, w.X, w.Y),
		VSync:     true,
		Resizable: true,
	}

	// create window
	win, err := opengl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var jumpOverCnt int
	win.Clear(colornames.Red)
	for !win.Closed() {

		jumpOverCnt = 1
		if win.Pressed(pixel.KeyLeftShift) || win.Pressed(pixel.KeyRightShift) {
			jumpOverCnt = 10
		}

		// with mouse-right or keypad DOWN button, we go to the NEXT picture
		if win.JustPressed(pixel.MouseButtonRight) || win.JustPressed(pixel.KeyDown) {
			// win.Clear(colornames.Whitesmoke)
			index += jumpOverCnt
			if index >= len(files) {
				index = 0
			}
			angleOfRotation = 0 // angle of rotation resets for new pic

			// with mouse-left or keypad UP button, we go to the PREV picture
		} else if win.JustPressed(pixel.MouseButtonLeft) || win.JustPressed(pixel.KeyUp) {
			// win.Clear(colornames.Whitesmoke)
			index -= 1
			if index < 0 {
				index = len(files) - jumpOverCnt
			}
			angleOfRotation = 0 // angle of rotation resets for new pic

			// with Space key pressed
		} else if win.JustPressed(pixel.KeySpace) {
			log.Infof("Spacebar pressed, toggling scale to fit %v", !scaleToFit)
			scaleToFit = !scaleToFit
			forceRedraw = true

			// keypad LEFT pressed
		} else if win.JustPressed(pixel.KeyLeft) {
			// rotate left
			switch angleOfRotation {
			case 0, 90, 180:
				angleOfRotation -= 90
			default: // 180
				angleOfRotation = 180
			}
			forceRedraw = true
			log.Infof("rotation anagle %v", angleOfRotation)
			// key pad RIGHT pressed
		} else if win.JustPressed(pixel.KeyLeft) {
			// rotate left
			switch angleOfRotation {
			case 0, 90, 180:
				angleOfRotation -= 90
			default: // 180
				angleOfRotation = 180
			}
			forceRedraw = true
			log.Infof("rotation anagle %v", angleOfRotation)
			// key pad RIGHT pressed
		} else if win.JustPressed(pixel.KeyRight) {
			// rotate right
			switch angleOfRotation {
			case -90, 0, 90:
				angleOfRotation += 90
			default: // -180
				angleOfRotation = -90
			}
			forceRedraw = true
			log.Infof("rotation anagle %v", angleOfRotation)
		} else if win.JustPressed(pixel.KeyEscape) {
			win.SetClosed(true)
		}

		var xx, yy *float64 // nil, no scaling...
		if scaleToFit {
			_xx := win.Bounds().Max.X
			_yy := win.Bounds().Max.Y
			xx = &_xx
			yy = &_yy
		}

		pathToDraw := files[index]
		// draw image on window only if a new path is selected, or forceRedraw is `true`
		if pathToDraw != lastDrawnPath || forceRedraw {
			loadPicInWindowScaledAndRotated(win, files[index], xx, yy, &angleOfRotation)
			lastDrawnPath = files[index] // this was drawn
			// reset after the draw
			forceRedraw = false // set force redraw false after draing successfully
		}

		win.Update()

	}

	log.Info("bye bye !! window was closed...")
}

type ViewPort struct {
	image.Rectangle

	// Path string
	ForceRedraw     bool
	ScaleToFit      bool
	AngleOfRotation int

	// internal
	lastDrawnPath string
}

func (v *ViewPort) Draw() {}

func loadPicInWindowScaledAndRotated(win *opengl.Window, path string, x, y *float64, rotationAngle *int) {
	log.Infof("loading image to window")

	// load image using gofoto
	img, err := gofoto_image.Open(path)
	if err != nil {
		log.Errorf("error opening image: %v\n", path)
	}

	if rotationAngle != nil {
		var err error
		log.Infof("rotation %v by %v deg", path, *rotationAngle)
		img, err = gofoto_image.Rotate(img, *rotationAngle)
		if err != nil {
			log.Errorf("error rotating image %v by %v deg", path, *rotationAngle)
		}
	}

	// scale to fit
	if x != nil || y != nil {
		log.Infof("scaling %v to fit window size %v, %v", path, *x, *y)

		max := img.Bounds().Max
		imgX := max.X
		imgY := max.Y

		scaleFactor := 1.0
		if imgX > imgY {
			scaleFactor = *x / float64(imgX)
		} else {
			scaleFactor = *y / float64(imgY)
		}
		img, err = gofoto_image.ResizeScale(img, float32(scaleFactor))
		if err != nil {
			log.Errorf("error scaling image %v by %v percent", path, scaleFactor)
		}
	}

	win.Clear(color.Black) // black background ...

	// load a pic & make sprite from the image
	pic := pixel.PictureDataFromImage(img)
	sprite := pixel.NewSprite(pic, pic.Bounds())
	sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))

	// TODO very expensive operation; so need a review here...
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(10, 10), basicAtlas)
	fmt.Fprintln(basicTxt, path)
	basicTxt.Draw(win, pixel.IM)
}
