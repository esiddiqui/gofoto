package templates

import (
	"embed"
	"html/template"

	"github.com/sirupsen/logrus"
)

//go:embed *
var templateFS embed.FS

var All *template.Template

func init() {
	var err error
	All, err = template.ParseFS(templateFS, "*")
	if err != nil {
		logrus.Fatal("error parsing html templates")
	}
	_ = All
}
