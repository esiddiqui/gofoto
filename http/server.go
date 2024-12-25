package http

import (
	"embed"
	"log"
	"net/http"
)

var addr string = "0.0.0.0:8080"

//go:embed static/*
var staticFS embed.FS // this will load static/* under the embedFS/static/

func StartWebserverAtRoot(rootPath string) {
	server := http.NewServeMux()
	staticfileServer := http.FileServer(http.FS(staticFS))
	noCacheStaticHandler := http.HandlerFunc(GetNoCacheWrapper(staticfileServer))
	server.Handle("/static/", noCacheStaticHandler)

	noCacheWrapperHandler := GetNoCacheWrapper(http.HandlerFunc(GetImageHandler(rootPath)))
	server.Handle("/", http.RedirectHandler("/browse", http.StatusMovedPermanently))
	server.HandleFunc("/browse/", http.HandlerFunc(GetListingUIHandler(rootPath)))
	server.HandleFunc("/view/", http.HandlerFunc(GetViewingUIHandler(rootPath)))
	server.HandleFunc("/show/", noCacheWrapperHandler)

	err := http.ListenAndServe(addr, server)
	if err != nil {
		log.Fatal(err)
	}
}
