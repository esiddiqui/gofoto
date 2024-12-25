package http

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/esiddiqui/gofoto/templates"
	"github.com/esiddiqui/gofoto/util"
)

// GetListingUIHandler builds & returns the album listing ui
func GetListingUIHandler(rootPath string) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		path := r.URL.Path
		// remove prefixs to fetch path of the resource...
		path = strings.TrimPrefix(path, "/browse")
		path = strings.TrimSuffix(path, "/")

		// fully-qualified path on file system
		absPath := fmt.Sprintf("%v/%v", rootPath, path)
		parent := path[:strings.LastIndex(path, "/")+1]

		//dirs := listDirectories(absPath)

		dat := struct {
			Title      string
			Path       string
			AbsPath    string
			Parent     string
			RootPath   string
			SrcAttr    string
			HrefBack   string
			HrefNext   string
			HrefParent string
			Items      []string
		}{
			Title:    path,
			Path:     path,
			Parent:   parent,
			RootPath: rootPath,
			AbsPath:  absPath,
			Items:    listDirectories(absPath),
		}
		templates.All.ExecuteTemplate(w, "lister.tmpl", dat)
	}
}

// ListDirectories returns a list of all children for the supplied path,
// directory names
func listDirectories(path string) []string {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil
	}

	var names []string
	filter := util.NewFilter(util.WithDirectories())
	for _, f := range files {
		fqn := fmt.Sprintf("%v/%v", path, f.Name())
		if filter.Allowed(fqn) {
			if f.IsDir() {
				names = append(names, f.Name())
			} else {
				names = append(names, f.Name())
			}
		}
	}

	sort.Strings(names)
	return names
}
