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

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"

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

func (s *Server) getContributionData(w http.ResponseWriter, id string) ([]model.ContributionEntry, string, error) {
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
			"LOGIN": id,
		},
	}

	err := json.NewEncoder(&requestBody).Encode(requestBodyObj)
	if err != nil {
		return nil, "", err
	}

	req, err := http.NewRequest("POST", "https://api.github.com/graphql", &requestBody)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.GitHubToken))

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	var tmp map[string]interface{}
	err = json.Unmarshal(data, &tmp)
	if err != nil {
		return nil, "", err
	}

	var result model.GitHubData
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, "", err
	}

	t, err := result.GetContributionOfLastSevenDays()
	if err != nil {
		return nil, "", err
	}

	username := result.Data.User.Name
	if username == "" {
		username = id
	}

	return t, username, nil
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

func (s *Server) handleRender(w http.ResponseWriter, r *http.Request) {
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

	contributiondata, username, err := s.getContributionData(w, params.Username)
	if err != nil {
		fmt.Fprintln(w, "Error while parsing response from GitHub, please check all your request parameters.")

		return
	}

	w.Header().Set("Content-Type", "text/html")

	tmpl.Execute(w, &model.TemplateData{
		ContributionData: contributiondata,
		Username:         username,
		ImgUrl:           "https://images.unsplash.com/photo-1518791841217-8f162f1e1131?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=800&q=60",
	})
}

func (s *Server) handleSVG(w http.ResponseWriter, r *http.Request) {
	params := getRequestParams(r.URL.Query())
	if params.Username == "" {
		fmt.Fprintln(w, "Url param 'username' is missing.")

		return
	}

	l := launcher.New()
	l.Set("no-sandbox")
	headlessBrowserUrl := l.MustLaunch()

	browser := rod.New().ControlURL(headlessBrowserUrl).MustConnect()

	page := browser.MustPage(fmt.Sprintf("http://localhost:8687/?username=%s", params.Username))
	page.WaitLoad()

	el := page.MustElement("#svg-container")

	w.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, el.MustHTML())
}

func (s *Server) Start() {
	ssrServer := http.NewServeMux()
	ssrServer.HandleFunc("/", s.handleRender)
	go http.ListenAndServe(":8687", ssrServer)

	svgServer := http.NewServeMux()
	svgServer.HandleFunc("/", s.handleSVG)
	http.ListenAndServe(":8686", svgServer)
}
