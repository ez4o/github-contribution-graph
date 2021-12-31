package util

import (
	"fmt"
	"html/template"
	"server/model"
	"strings"
)

func HydrateContributionData(data []model.ContributionEntry) template.JS {
	entryToHydrate := make([]string, len(data))
	for i, entry := range data {
		entryToHydrate[i] = fmt.Sprintf(
			"{dateString: '%s', amount: %d}",
			entry.DateString, entry.Amount,
		)
	}

	return template.JS(strings.Join(entryToHydrate, ",\n"))
}
