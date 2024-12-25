package http

import (
	"net/url"
	"strconv"
)

func getRotation(parms url.Values) int {
	v := parms.Get("r")
	vint, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return 0
	}
	return int(vint)
}

func getScale(parms url.Values) float32 {
	v := parms.Get("s")
	vint, err := strconv.ParseFloat(v, 32)
	if err != nil {
		return 0.2
	}
	return float32(vint)
}

func getFile(parms url.Values) string {
	return parms.Get("f")
}
