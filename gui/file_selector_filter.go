package gui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Filter struct {
	extensions map[string]struct{} // allowed extensions
	dot        bool                // allow dot files?
	dir        bool                // allow directories, dot or not
	suffixes   map[string]struct{} // allowed suffixes
	prefixes   map[string]struct{} // allowed prefixes
}

type FilterOption func(*Filter)

func NewFilter(opts ...FilterOption) *Filter {

	f := &Filter{
		extensions: make(map[string]struct{}),
		dot:        false,
		dir:        false,
		suffixes:   make(map[string]struct{}),
		prefixes:   make(map[string]struct{}),
	}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

func WithExtensions(extensions ...string) FilterOption {
	return func(f *Filter) {
		for _, ext := range extensions {
			f.extensions[ext] = struct{}{}
		}
	}
}

func WithSuffixes(suffixes ...string) FilterOption {
	return func(f *Filter) {
		for _, suffix := range suffixes {
			f.suffixes[suffix] = struct{}{}
		}
	}
}

func withPrefixes(prefixes ...string) FilterOption {
	return func(f *Filter) {
		for _, prefix := range prefixes {
			f.prefixes[prefix] = struct{}{}
		}
	}
}

func WithDotFiles() FilterOption {
	return func(f *Filter) {
		f.dot = true
	}
}

func WithDirectories() FilterOption {
	return func(f *Filter) {
		f.dir = true
	}
}

func (f *Filter) Allowed(path string) bool {

	var rdr, rse, rss, rsp bool

	// check if dir
	file, err := os.Open(path)
	if err != nil {
		log.Errorf("error opening path %v to check if it's directory: %v", path, err)
		return false
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		log.Errorf("error checking file stats for %v", path)
		return false
	}
	if f.dir && stat.IsDir() {
		rdr = true
	}

	// check if dot file
	if !f.dot {
		sections := strings.Split(path, "/")
		for _, sec := range sections {
			if strings.HasPrefix(sec, ".") {
				return false
			}
		}
	}

	// extensions must be from the white-listed map..
	extension := strings.ToLower(filepath.Ext(path))

	if _, ok := f.extensions[extension]; ok {
		rse = true
	}

	for k, _ := range f.prefixes {
		if strings.HasPrefix(path, k) {
			rsp = true
			break
		}
	}

	for k, _ := range f.suffixes {
		if strings.HasSuffix(path, k) {
			rss = true
			break
		}
	}

	log.Debugf("[PASS] isExtension: good to open %v\n", path)
	return rdr || rse || rss || rsp
}

// convenience method to load only JPEG files from the supplied directory path
func listJpegFiles(path string) []string {

	log.Debugf("loading files from %v\n", path)
	files, err := os.ReadDir(path)
	if err != nil {
		return nil
	}

	// filter only jpg or jpeg extension-ed files
	var names []string
	filter := NewFilter(WithExtensions(".jpg", ".jpeg"))
	for _, f := range files {
		fqn := fmt.Sprintf("%v/%v", path, f.Name())
		if filter.Allowed(fqn) {
			names = append(names, fqn)
		}
	}

	sort.Strings(names)
	for _, k := range names {
		log.Info(k)
	}
	return names
}
