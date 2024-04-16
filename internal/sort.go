package internal

import (
	"fmt"
	"math"
	"sort"
)


func sortCragsByDistance(crags CragsList, filters CragFilters) CragsList {
	if filters.Location == "" || ! isValidLocation(filters.Location) {
		return crags
	}

    sorted := CragsList{
        CragsList: make([]Crag, len(crags.CragsList)),
    }
    copy(sorted.CragsList, crags.CragsList)

	var centerLat, centerLon float64
	fmt.Sscanf(filters.Location, "%f,%f", &centerLat, &centerLon)

	sort.Slice(sorted.CragsList, func(i, j int) bool {
		distI := haversine(centerLat, centerLon, sorted.CragsList[i].Lat, sorted.CragsList[i].Lon)
		distJ := haversine(centerLat, centerLon, sorted.CragsList[j].Lat, sorted.CragsList[j].Lon)

		return distI < distJ
	})

	return sorted
}

func isValidLocation(location string) bool {
	var lat, lon float64

	n, err := fmt.Sscanf(location, "%f,%f", &lat, &lon)

	if err != nil || n != 2 {
		return false
	}

	return lat >= -90 && lat <= 90 && lon >= -180 && lon <= 180
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	var rad = math.Pi / 180
	var h = hav(rad * (lat2 - lat1)) + math.Cos(rad * lat1) * math.Cos(rad * lat2) * hav(rad * (lon2 - lon1))

	return 2 * 6371 * math.Asin(math.Sqrt(h))
}

func hav(theta float64) float64 {
	return 0.5 - math.Cos(theta) / 2
}
