package model

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

func (g *GitHubData) GetContributionOfLastSevenDays() []ContributionEntry {
	weeks := g.Data.User.ContributionsCollection.ContributionCalendar.Weeks
	weekLength := len(g.Data.User.ContributionsCollection.ContributionCalendar.Weeks)

	var contributionOfLastTwoWeeks []ContributionEntry
	lastTwoWeeks := weeks[weekLength-2:]

	for _, week := range lastTwoWeeks {
		for _, day := range week.ContributionDays {
			if day.Date == "" {
				continue
			}

			contributionOfLastTwoWeeks = append(contributionOfLastTwoWeeks, day.GetEntry())
		}
	}

	return contributionOfLastTwoWeeks[len(contributionOfLastTwoWeeks)-7:]
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
