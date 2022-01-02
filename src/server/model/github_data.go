package model

import (
	"fmt"
	"math"
)

type GitHubData struct {
	Data struct {
		User struct {
			Name                    string `json:"name"`
			ContributionsCollection struct {
				ContributionCalendar struct {
					TotalContributions int `json:"totalContributions"`
					Weeks              []struct {
						ContributionDays []ContributionDays `json:"contributionDays"`
					} `json:"weeks"`
				} `json:"contributionCalendar"`
			} `json:"contributionsCollection"`
		} `json:"user"`
	} `json:"data"`
}

func (g *GitHubData) GetContributionOfLastNDays(lastNDays int) ([]ContributionEntry, error) {
	weeks := g.Data.User.ContributionsCollection.ContributionCalendar.Weeks
	weekLength := len(g.Data.User.ContributionsCollection.ContributionCalendar.Weeks)
	nWeek := math.Ceil(float64(lastNDays)/float64(7)) + 1

	if float64(weekLength) < nWeek {
		return nil, fmt.Errorf("Data is not enough to get last %d days", lastNDays)
	}

	var contributionOfLastTwoWeeks []ContributionEntry
	lastNWeeks := weeks[weekLength-int(nWeek):]

	for _, week := range lastNWeeks {
		for _, day := range week.ContributionDays {
			if day.Date == "" {
				continue
			}

			contributionOfLastTwoWeeks = append(contributionOfLastTwoWeeks, day.GetEntry())
		}
	}

	return contributionOfLastTwoWeeks[len(contributionOfLastTwoWeeks)-lastNDays:], nil
}

type ContributionDays struct {
	ContributionCount int    `json:"contributionCount"`
	Date              string `json:"date"`
}

func (c ContributionDays) GetEntry() ContributionEntry {
	return ContributionEntry{
		DateString: c.Date,
		Amount:     c.ContributionCount,
	}
}
