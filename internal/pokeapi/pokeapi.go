// Package pokeapi helps to connect to the pokeapi
package pokeapi

import "net/url"

var baseURL = "https://pokeapi.co/api/v2/"

type LocationResult struct {
	count    int
	next     url.URL
	previous url.URL
	restult  []Location
}

type Location struct {
	name string
	url  url.URL
}

func GetLocations(link url.URL) (string, error) {
	apiLink := baseURL + "location-area/"
	if link != (url.URL{}) {
		apiLink = link.String()
	}
	return apiLink, nil
}
