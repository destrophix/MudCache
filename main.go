package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type LocationData struct {
	Id              string  `json:"id"`
	SeismicActivity float64 `json:"seismic_activity"`
	TemperatureC    float64 `json:"temperature_c"`
	RadiationLevel  float64 `json:"radiation_level"`
	LocationId      string  `json:"location_id"`
}

var locations = make(map[string]LocationData)

func locationHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		pathParts := strings.Split(req.URL.Path, "/")
		json.NewEncoder(w).Encode(locations[pathParts[1]])

	case http.MethodPut:
		pathParts := strings.Split(req.URL.Path, "/")
		var inputData LocationData
		decoder := json.NewDecoder(req.Body)
		if err := decoder.Decode(&inputData); err != nil {
			http.Error(w, "Failed to parse request body", http.StatusBadRequest)
			return
		}
		locationID := pathParts[1]
		inputData.LocationId = locationID
		key := locationID
		locations[key] = inputData
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(inputData)
		if err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

func healthHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {
	var asciiArt = `
 __  __           _  ____           _          
|  \/  |_   _  __| |/ ___|__ _  ___| |__   ___ 
| |\/| | | | |/ _  | |   / _  |/ __| '_ \ / _ \
| |  | | |_| | (_| | |__| (_| | (__| | | |  __/
|_|  |_|\__,_|\__,_|\____\__,_|\___|_| |_|\___|
	`
	fmt.Println(asciiArt)

	http.HandleFunc("/", locationHandler)

	http.HandleFunc("/health", healthHandler)

	http.ListenAndServe(":8090", nil)
}
