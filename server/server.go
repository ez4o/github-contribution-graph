package server

import (
	"fmt"
	"html/template"
	"net/http"
	"path"
	"strings"

	"server/model"
)

type Server struct{}

func (s *Server) getContributionData() []model.ContributionEntry {
	return []model.ContributionEntry{
		{
			DateString: "2018-01-01",
			Amount:     18,
		},
		{
			DateString: "2018-01-02",
			Amount:     12,
		},
		{
			DateString: "2018-01-03",
			Amount:     16,
		},
	}
}

func (s *Server) hydrateContributionData(data []model.ContributionEntry) template.JS {
	entryToHydrate := make([]string, len(data))
	for i, entry := range data {
		entryToHydrate[i] = fmt.Sprintf(
			"{dateString: '%s', amount: %d}",
			entry.DateString, entry.Amount,
		)
	}

	return template.JS(strings.Join(entryToHydrate, ",\n"))
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("index.html").Funcs(
		template.FuncMap{
			"hydrateContributionData": s.hydrateContributionData,
		},
	).ParseFiles(
		path.Join("dist", "index.html"),
	)
	if err != nil {
		panic(err)
	}

	tmpl.Execute(w, &model.TemplateData{
		ContributionData: s.getContributionData(),
	})
}

func (s *Server) Start() {
	http.HandleFunc("/", s.handleIndex)

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
