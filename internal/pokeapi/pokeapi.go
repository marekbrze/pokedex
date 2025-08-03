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

		body, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		go cache.Add(link, body)
		fmt.Println("Entry from API")
		return body, nil
	}
}

func UnmarshalLocations(body []byte) (LocationResult, error) {
	var locationResult LocationResult
	if err := json.Unmarshal(body, &locationResult); err != nil {
		return LocationResult{}, nil
	}
	return locationResult, nil
}
