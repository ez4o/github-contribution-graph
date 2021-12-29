package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"server/model"
)

type Server struct {
	GitHubToken string
}

func getRequestParams(v url.Values) model.RequestParams {
	username := v.Get("username")

	params := model.RequestParams{
		Username: username,
	}

	return params
}

func (s *Server) getContributionData(w http.ResponseWriter, username string) ([]model.ContributionEntry, error) {
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
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", &requestBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.GitHubToken))

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tmp map[string]interface{}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return nil, err
	}

	var result model.GitHubData
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	t, err := result.GetContributionOfLastSevenDays()
	if err != nil {
		return nil, err
	}

	return t, nil
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
	params := getRequestParams(r.URL.Query())
	if params.Username == "" {
		fmt.Fprintln(w, "Url param 'username' is missing.")

		return
	}

	tmpl, err := template.New("index.html").Funcs(
		template.FuncMap{
			"hydrateContributionData": s.hydrateContributionData,
		},
	).ParseFiles(
		path.Join("dist", "index.html"),
	)
	if err != nil {
		fmt.Fprintln(w, "Error while parsing template.")

		return
	}

	contributiondata, err := s.getContributionData(w, params.Username)
	if err != nil {
		fmt.Fprintln(w, "Error while parsing response from GitHub, please check all your request parameters.")

		return
	}

	tmpl.Execute(w, &model.TemplateData{
		ContributionData: contributiondata,
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
