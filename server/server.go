package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path"
	"strings"
	"time"

	"server/model"
)

type Server struct {
	GitHubToken string
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}
}

func (s *Server) getContributionData(username string) []model.ContributionEntry {
	var requestBody bytes.Buffer

	requestBodyObj := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}{
		Query: `
			query userInfo($LOGIN: String!) {
				user(login: $LOGIN) {
					name
					contributionsCollection {
					contributionCalendar {
						totalContributions
							weeks {
								contributionDays {
									contributionCount
									date
								}
							}
						}
					}
				}
			},
		`,
		Variables: map[string]interface{}{
			"LOGIN": username,
		},
	}

	err := json.NewEncoder(&requestBody).Encode(requestBodyObj)
	handleError(err)

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", &requestBody)
	handleError(err)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.GitHubToken))

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	handleError(err)
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	handleError(err)

	var result model.GitHubData
	json.Unmarshal(data, &result)

	t := result.GetContributionOfLastSevenDays()
	return t
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
		ContributionData: s.getContributionData("wst24365888"),
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
