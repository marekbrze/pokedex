// Package pokeapi helps to connect to the pokeapi
package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var baseURL = "https://pokeapi.co/api/v2/"

type LocationResult struct {
	Count    int
	Next     string
	Previous string
	Results  []Location
}

type Location struct {
	Name string
	URL  string
}

func GetLocations(link string) (LocationResult, error) {
	var locationResult LocationResult
	res, err := http.Get(link)
	if err != nil {
		return LocationResult{}, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationResult{}, err
	}

	if err := json.Unmarshal(body, &locationResult); err != nil {
		return LocationResult{}, nil
	}

	fmt.Println(locationResult.Count, locationResult.Next)
	return locationResult, nil
}
