package internal

import (
	"strings"
)


func filterCrags(crags CragsList, filters CragFilters) CragsList {
	var filteredCrags []Crag

	for _, crag := range crags.CragsList {
		if matchFilter(crag, filters) {
			filteredCrags = append(filteredCrags, crag)
		}
	}

	return CragsList{CragsList: filteredCrags}
}

func matchFilter(crag Crag, filters CragFilters) bool {
	matchName := true
	matchDiscipline := true

	if filters.Name != "" && !strings.Contains(strings.ToLower(crag.Name), strings.ToLower(filters.Name)) {
		matchName = false
	}

	if filters.Disciplines != "" {
		filterDisciplines := strings.Split(strings.ToLower(filters.Disciplines), ",")
		matchDiscipline = false

		for _, disc := range crag.Disciplines {
			for _, filterDisc := range filterDisciplines {
				if strings.ToLower(disc) == filterDisc {
					matchDiscipline = true
					break
				}
			}

			if matchDiscipline {
				break
			}
		}
	}

	return matchName && matchDiscipline
}
