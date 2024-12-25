package http

import (
	"fmt"
	"net/http"
	"strings"

	goimage "image"

	"github.com/esiddiqui/gofoto/image"
)

// GetImageHandler builds & returns a new http handler that responds to requests
// with an image scaled & rotated with the supplied
func GetImageHandler(rootPath string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		// remove prefixs to fetch path of the resource...
		path := strings.TrimPrefix(r.URL.Path, "/show")
		path = strings.TrimSuffix(path, "/")

		rotation := getRotation(r.URL.Query())
		scale := getScale(r.URL.Query())
		_ = scale
		fqp := fmt.Sprintf("%v/%v", rootPath, path)

		ext := ""
		var img goimage.Image
		if strings.Contains(fqp, ".") {
			ext = fqp[strings.LastIndex(fqp, ".")+1:]
		}
		switch strings.ToLower(ext) {
		case "cr2":
			img, _ = image.OpenCr2(fqp)
		default:
			img, _ = image.Open(fqp)
		}

		img, _ = image.ResizeScale(img, scale)
		img, _ = image.Rotate(img, rotation)
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")
		image.WriteTo(image.NewJpegEncoder(90), w, img) //encode JPEG
	}
}
