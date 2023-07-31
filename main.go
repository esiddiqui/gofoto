package main

import (
	"embed"
	"net/http"
	"os"

	_ "embed"

	"github.com/esiddiqui/gofoto/handler"
	log "github.com/sirupsen/logrus"
)

var addr string = "0.0.0.0:8080"
var rootPath string = "/"

//go:embed static/*
var staticFS embed.FS // this will load static/* under the embedFS/static/

func main() {

	if len(os.Args) > 1 {
		rootPath = os.Args[1]
	} else {
		rootPath = os.Getenv("HOME")
	}

	log.Infof("root is %v", rootPath)

	server := http.NewServeMux()

	staticfileServer := http.FileServer(http.FS(staticFS))
	noCacheStaticHandler := http.HandlerFunc(handler.GetNoCacheWrapper(staticfileServer))
	server.Handle("/static/", noCacheStaticHandler)

	noCacheWrapperHandler := handler.GetNoCacheWrapper(http.HandlerFunc(handler.GetImageHandler(rootPath)))
	server.Handle("/", http.RedirectHandler("/browse", http.StatusMovedPermanently))
	server.HandleFunc("/browse/", http.HandlerFunc(handler.GetListingUIHandler(rootPath)))
	server.HandleFunc("/view/", http.HandlerFunc(handler.GetViewingUIHandler(rootPath)))
	server.HandleFunc("/show/", noCacheWrapperHandler)

	err := http.ListenAndServe(addr, server)
	if err != nil {
		log.Fatal(err)
	}
}
