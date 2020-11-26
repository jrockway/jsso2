package main

import (
	"net/http"

	"github.com/jrockway/opinionated-server/server"
)

func main() {
	server.AppName = "jsso2-protected-example"
	server.Setup()
	server.SetHTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("content-type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><head><title>Example Protected App</title></head><body><p>"))
		if username := req.Header.Get("x-jsso2-username"); username != "" {
			w.Write([]byte("Logged in as "))
			w.Write([]byte(username))
		} else {
			w.Write([]byte("Not logged in"))
		}
		w.Write([]byte("</p></body></html>"))
	}))
	server.ListenAndServe()
}
