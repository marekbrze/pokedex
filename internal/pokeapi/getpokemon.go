package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/marekbrze/pokedexcli/internal/pokecache"
)

func GetPokemon(name string, cache *pokecache.Cache) ([]byte, error) {
	link := "https://pokeapi.co/api/v2/pokemon/" + name
	entry, exist := cache.Get(link)
	if exist {
		fmt.Println("entry from cache")
		return entry, nil
	} else {

		res, err := http.Get(link)
		if err != nil {
			return nil, fmt.Errorf("error when making call to the pokomen endpoint. Info: %w", err)
		}
		defer res.Body.Close()

		if res.StatusCode == http.StatusNotFound {
			return nil, fmt.Errorf("pokemon not found")
		}

		if res.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("error when making call to the pokomen api. StatusCode: %d", res.StatusCode)
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

func UnmarshalPokemon(body []byte) (Pokemon, error) {
	var pokemon Pokemon
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return Pokemon{}, err
	}
	return pokemon, nil
}
