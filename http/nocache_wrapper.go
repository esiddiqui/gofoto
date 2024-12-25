package http

import (
	"net/http"
	"time"
)

// wrapNoCache wraps 2 handlers into one & add no-cache headerse
func GetNoCacheWrapper(h http.Handler) func(w http.ResponseWriter, r *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {

		// log.Infof("path requested %v", r.URL.Path)
		var epoch = time.Unix(0, 0).Format(time.RFC1123)
		var noCacheHeaders = map[string]string{
			"Expires":         epoch,
			"Cache-Control":   "no-cache, private, max-age=0",
			"Pragma":          "no-cache",
			"X-Accel-Expires": "0",
		}

		// Set our NoCache headers
		for k, v := range noCacheHeaders {
			w.Header().Set(k, v)
		}

		w.Header().Set("x-served-by", "h0")
		h.ServeHTTP(w, r)
	}
}
