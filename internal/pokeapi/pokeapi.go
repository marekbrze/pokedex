// Package pokeapi helps to connect to the pokeapi
package pokeapi

import "net/http"
var baseURL = "https://pokeapi.co/api/v2/"

type LocationResult struct {
	count    int
	next     string
	previous string
	result   []Location
}

type Location struct {
	name string
	url  string
}

func GetLocations(link string) (LocationResult, error) {
	res, err := http.Get(link)
	if err != nil {
		return 
	}
	return LocationResult{}, nil
}
