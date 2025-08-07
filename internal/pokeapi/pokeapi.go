// Package pokeapi helps to connect to the pokeapi
package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/marekbrze/pokedexcli/internal/pokecache"
)

var baseURL = "https://pokeapi.co/api/v2/"

type Locations struct {
	Name string
	URL  string
}

type LocationsResult struct {
	Count    int
	Next     string
	Previous string
	Results  []Locations
}

func GetLocations(link string, cache *pokecache.Cache) ([]byte, error) {
	entry, exist := cache.Get(link)
	if exist {
		fmt.Println("Entry from Cache")
		return entry, nil
	} else {
		res, err := http.Get(link)
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			return nil, err
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		go cache.Add(link, body)
		fmt.Println("Entry from API")
		return body, nil
	}
}

func UnmarshalLocations(body []byte) (LocationsResult, error) {
	var locationResult LocationsResult
	if err := json.Unmarshal(body, &locationResult); err != nil {
		return LocationsResult{}, nil
	}
	return locationResult, nil
}

func ExploreLocation(name string, cache *pokecache.Cache) ([]byte, error) {
	link := "https://pokeapi.co/api/v2/location-area/" + name
	entry, exist := cache.Get(link)
	if exist {
		fmt.Println("Entry from Cache")
		return entry, nil
	} else {
		res, err := http.Get(link)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("location not found")
		} else if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("StatusCode: %d", res.StatusCode)
		}
		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		go cache.Add(link, body)
		fmt.Println("Entry from API")
		return body, nil
	}
}

func UnmarshalSingleLocation(body []byte) (Location, error) {
	var singleLocationResult Location
	if err := json.Unmarshal(body, &singleLocationResult); err != nil {
		return Location{}, nil
	}
	return singleLocationResult, nil
}
