package main

import (
	"html/template"
	"net/http"
	"path"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(path.Join("dist", "index.html")))
	tmpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", handleIndex)

	// Create file server at "./dist/assets".
	// Cut off the prefix "/assets/" of request path and forward to it.
	fs := http.FileServer(http.Dir("./dist/assets"))
	http.HandleFunc("/assets/", func(w http.ResponseWriter, r *http.Request) {
		if path.Ext(r.URL.Path) == ".js" {
			// https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types/Common_types
			w.Header().Set("Content-Type", "text/javascript")
		}

		http.StripPrefix("/assets/", fs).ServeHTTP(w, r)
	})

	http.ListenAndServe(":80", nil)
}
