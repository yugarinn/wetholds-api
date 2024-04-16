package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)


type CragsList struct {
	CragsList []Crag `json:"crags"`
}

type Crag struct {
	Name        string `json:"name"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Disciplines []string `json:"disciplines"`
	Weather     CragWeather `json:"weather"`
}

type CragWeather struct {
	Temperature 			 map[string][]int `json:"temperature"`
	PrecipitationProbability map[string][]int `json:"precipitationProbability"`
	WindSpeed   			 map[string][]int `json:"windSpeed"`
	CloudsCover              map[string][]int `json:"cloudCover"`
}

type CragFilters struct {
	Location    string
	Disciplines string
	Name        string
}

func InitServer() {
    http.HandleFunc("/crags", cragsHandler)

	port := ":9990"
	log.Printf("starting server on port %s", port)
    log.Fatal(http.ListenAndServe(port, nil))
}

func cragsHandler(writer http.ResponseWriter, request *http.Request) {
	go logRequest(request)

	crags, isCached, loadErr := loadCrags()

	if loadErr != nil {
		log.Println(loadErr.Error())
	}

	if !isCached {
		go CacheCragsResponse(crags)
	}

	filters := parseQueryFilters(request)
	crags = filterCrags(crags, filters)
	crags = sortCragsByDistance(crags, filters)

	writer.Header().Set("Content-Type", "application/json")
	writer.Header().Set("Access-Control-Allow-Origin", "*")

	json.NewEncoder(writer).Encode(crags)
}

func parseQueryFilters(request *http.Request) CragFilters {
	queryParams := request.URL.Query()

	return CragFilters{
		Location:    queryParams.Get("location"),
		Disciplines: queryParams.Get("disciplines"),
		Name:        queryParams.Get("name"),
	}
}

func loadCrags() (list CragsList, isCached bool, err error) {
	cachedCrags, isValidCache := GetCachedCrags()

	if false || isValidCache {
		return cachedCrags, true, nil
	}

	dataPath := os.Getenv("DATA_PATH")
	file, _ := os.ReadFile(fmt.Sprintf("%s/crags.json", dataPath))
	crags := CragsList{}

	unmarshalError := json.Unmarshal([]byte(file), &crags)

	var wg sync.WaitGroup
    for i := 0; i < len(crags.CragsList); i++ {
		wg.Add(1)
		crag := &crags.CragsList[i]

		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			hydrateCragWithWeather(crag)
		}(&wg)
    }

	wg.Wait()

	return crags, false, unmarshalError
}

func hydrateCragWithWeather(crag *Crag) {
	weatherResponse := FetchForecast(crag.Lat, crag.Lon)

	crag.Weather = CragWeather{
		Temperature: weatherResponse.Temperatures(),
		PrecipitationProbability: weatherResponse.PrecipitationProbabilties(),
		WindSpeed: weatherResponse.WindSpeeds(),
		CloudsCover: weatherResponse.CloudCovers(),
	}
}

func logRequest(r *http.Request) {
	log.Printf("%s %s", r.Method, r.RequestURI)
}
