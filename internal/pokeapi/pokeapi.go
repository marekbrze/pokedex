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
	Url  string
}

func GetLocations(link string) (LocationResult, error) {
	fmt.Println("Inner test")
	var locationResult LocationResult
	res, err := http.Get(link)
	if err != nil {
		fmt.Println("Error 1")
		return LocationResult{}, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error 2")
		return LocationResult{}, err
	}

	if err := json.Unmarshal(body, &locationResult); err != nil {
		fmt.Println("Error 3")
		return LocationResult{}, nil
	}

	fmt.Println(locationResult.Count, locationResult.Next)
	return locationResult, nil
}
