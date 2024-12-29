package gui

import "github.com/gopxl/pixel/v2/backends/opengl"

type ViewPort interface {
	Draw(win *opengl.Window, state windowState) error
}
