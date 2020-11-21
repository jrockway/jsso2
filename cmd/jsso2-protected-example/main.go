package main

import (
	"net/http"

	"github.com/jrockway/opinionated-server/server"
)

func main() {
	server.AppName = "jsso2-protected-example"
	server.Setup()
	server.SetHTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("content-type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}))
	server.ListenAndServe()
}
