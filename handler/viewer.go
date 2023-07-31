package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/esiddiqui/gofoto/image"
	"github.com/esiddiqui/gofoto/templates"
	"github.com/sirupsen/logrus"
)

var extensions = []string{".jpg", ".jpeg", ".png"}

// GetViewingUIHandler builds & returns the album viewing ui
func GetViewingUIHandler(rootPath string) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path

		// remove prefixs to fetch path of the resource...
		path = strings.TrimPrefix(path, "/view")
		path = strings.TrimSuffix(path, "/")
		fileName := fmt.Sprintf("%v", getFile(r.URL.Query()))
		rotation := fmt.Sprintf("%v", getRotation(r.URL.Query()))
		scale := fmt.Sprintf("%v", getScale(r.URL.Query()))

		parent := path[:strings.LastIndex(path, "/")+1]

		prev, this, next := getPreviousThisAndNextFiles(fmt.Sprintf("%v/%v", rootPath, path), fileName)

		meta := &image.ImageMetadata{}
		if len(this) != 0 {
			logrus.Warn("Path", rootPath, path)
			absPath := fmt.Sprintf("%v/%v/%v", rootPath, path, this)
			meta, _ = image.Metadata(absPath)
		}

		dat := struct {
			Title      string
			Path       string
			Parent     string
			RootPath   string
			Metadata   image.ImageMetadata
			SrcAttr    string
			HrefBack   string
			HrefSelf   string
			HrefNext   string
			HrefParent string
			HrefRoot   string
			Scale      float32
			Rotation   int
		}{
			Title:      path,
			Path:       path,
			Parent:     parent,
			RootPath:   rootPath,
			Metadata:   *meta,
			SrcAttr:    fmt.Sprintf(`/show%v/%v?s=%v&r=%v`, path, this, scale, rotation),
			HrefBack:   fmt.Sprintf(`/view%v?f=%v&s=%v&r=0`, path, prev, scale),
			HrefNext:   fmt.Sprintf(`/view%v?f=%v&s=%v&r=0`, path, next, scale),
			HrefSelf:   fmt.Sprintf(`/view%v?f=%v`, path, this),
			HrefParent: fmt.Sprintf(`/browse%v`, parent),
			HrefRoot:   `/browse`,
			Scale:      getScale(r.URL.Query()),
			Rotation:   getRotation(r.URL.Query()),
		}
		templates.All.ExecuteTemplate(w, "viewer.tmpl", dat)
	}
}

// getPreviousThisAndNextFiles would return the previous, current & next files in alphabetical
// order for the supplied file.
//
// if name is supplied, it becomes the `this` file, else the first file in the
// returned listing for image files is returned as the `this` file
//
// if the path contains 0 files, then empty strings are returned
// if the path contains 1 file, then prev, this & next files are the same
//
// if the name matches the first file in alphabetical ordered list of chilidren,
// it becomes the `this` file, while the 2nd file is the next & last file is the previous
//
// if the name matches the last file in alphabetical ordered list of children,
// the 1st file is the next & the n-1-th file is the previous

func getPreviousThisAndNextFiles(path, name string) (string, string, string) {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return "", "", ""
	}

	var names []string
	for _, f := range files {
		if !f.IsDir() && !strings.HasPrefix(f.Name(), ".") && (strings.HasSuffix(f.Name(), ".JPG") || strings.HasSuffix(f.Name(), ".jpg")) || (strings.HasSuffix(f.Name(), ".CR2") || strings.HasSuffix(f.Name(), ".cr2")) {
			names = append(names, f.Name())
		}
	}

	if len(names) == 0 {
		return "", "", ""
	}

	sort.Strings(names)

	// if name is not supplied, use the first name as the name to use
	if len(name) == 0 {
		name = names[0]
	}

	// if there's only 1 file, then it's the prev, this & next
	if len(names) == 1 {
		return name, name, name
	}

	// for over 2 files in the list..

	length := len(names)
	for n, fn := range names {
		if fn == name {
			this := n
			prev := this - 1 //(this - 1) % length //int(math.Abs(float64(this-1))) % length
			if prev < 0 {
				prev = length + prev
			}
			next := (this + 1) % length
			return names[prev], names[this], names[next]
		}
	}
	return "", "", ""
}
