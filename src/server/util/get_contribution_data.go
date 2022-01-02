package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"server/model"
	"time"
)

func GetContributionData(w http.ResponseWriter, id string, githubToken string, lastNDays int) ([]model.ContributionEntry, string, error) {
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
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", githubToken))

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

	t, err := result.GetContributionOfLastNDays(lastNDays)
	if err != nil {
		return nil, "", err
	}

	username := result.Data.User.Name
	if username == "" {
		username = id
	}

	return t, username, nil
}
