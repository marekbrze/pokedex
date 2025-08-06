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

type LocationsResult struct {
	Count    int
	Next     string
	Previous string
	Results  []Location
}

type SingleLocationResult struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type Location struct {
	Name string
	URL  string
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
			return nil, fmt.Errorf("Location not found")
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

func UnmarshalSingleLocation(body []byte) (SingleLocationResult, error) {
	var singleLocationResult SingleLocationResult
	if err := json.Unmarshal(body, &singleLocationResult); err != nil {
		return SingleLocationResult{}, nil
	}
	return singleLocationResult, nil
}
