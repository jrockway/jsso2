package main

import (
	"fmt"
	"net/http"
	"sort"
	"text/template"

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
			w.Write([]byte(template.HTMLEscapeString(username)))
		} else {
			w.Write([]byte("Not logged in"))
		}
		w.Write([]byte("</p><ul>\n"))
		var keys []string
		for k := range req.Header {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			for _, v := range req.Header[k] {
				w.Write([]byte(fmt.Sprintf("    <li>%s: %s</li>\n", template.HTMLEscapeString(k), template.HTMLEscapeString(v))))
			}
		}
		w.Write([]byte("\n</ul></body></html>\n"))
	}))
	server.ListenAndServe()
}
