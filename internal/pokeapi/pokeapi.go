// Package pokeapi helps to connect to the pokeapi
package pokeapi

import (
	"encoding/json"
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

func GetLocations(link string) ([]byte, error) {
	res, err := http.Get(link)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func UnmarshalLocations(body []byte) (LocationResult, error) {
	var locationResult LocationResult
	if err := json.Unmarshal(body, &locationResult); err != nil {
		return LocationResult{}, nil
	}
	return locationResult, nil
}
