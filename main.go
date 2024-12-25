package main

import (
	"os"

	_ "embed"

	"github.com/esiddiqui/gofoto/gui"
	"github.com/gopxl/pixel/v2/backends/opengl"
	log "github.com/sirupsen/logrus"
)

var rootPath string = "/"

func main() {
	log.Infof("gophoto is up")

	if len(os.Args) > 1 {
		rootPath = os.Args[1]
	} else {
		rootPath = "/"
	}

	log.Infof("root is %v", rootPath)

	window := &gui.Window{
		X:        1373,
		Y:        1063,
		ViewPort: nil,
		Path:     rootPath,
	}
	opengl.Run(window.Draw)

	// web
	// http.StartWebserverAtRoot(rootPath)

}
