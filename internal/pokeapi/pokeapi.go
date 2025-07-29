// Package pokeapi helps to connect to the pokeapi
package pokeapi

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
	return LocationResult{}, nil
}
