package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/marekbrze/pokedexcli/internal/pokeapi"
	"github.com/marekbrze/pokedexcli/internal/pokecache"
)

// INFO: Main types

type config struct {
	next     string
	previous string
	cache    *pokecache.Cache
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config, []string) error
}

var (
	PokeConfig      config
	commandRegistry = make(map[string]cliCommand)
)

// INFO: Commands list
func init() {
	commandRegistry["map"] = cliCommand{
		name:        "map",
		description: "Displays pokemon locations",
		callback:    commandMap,
	}
	commandRegistry["mapb"] = cliCommand{
		name:        "mapb",
		description: "Displays previous locations page",
		callback:    commandMap2,
	}
	commandRegistry["explore"] = cliCommand{
		name:        "explore",
		description: "List all the pokemon from the region passed as the second argument",
		callback:    commandExplore,
	}
	commandRegistry["help"] = cliCommand{
		name:        "help",
		description: "Displays a help message",
		callback:    commandHelp,
	}
	commandRegistry["exit"] = cliCommand{
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	}
	PokeConfig.next = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"
	PokeConfig.previous = ""
}

// INFO: Main Loop
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	PokeConfig.cache = pokecache.NewCache(15 * time.Second)
	fmt.Print("Pokedex > ")
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			fmt.Print("Pokedex > ")
			continue
		}
		cleanedInput := cleanInput(line)
		if len(cleanedInput) == 0 {
			fmt.Print("Pokedex > ")
			continue
		}
		firstCommand := cleanedInput[0]
		command, exists := commandRegistry[firstCommand]
		if exists && len(cleanedInput) > 1 {
			err := command.callback(&PokeConfig, cleanedInput[1:])
			if err != nil {
				fmt.Print(fmt.Errorf("error: %w", err))
			}
		} else if exists {
			err := command.callback(&PokeConfig, []string{})
			if err != nil {
				fmt.Print(fmt.Errorf("error: %w", err))
			}
		} else {
			fmt.Println("Unknown command")
		}
		fmt.Print("Pokedex > ")
	}
}

// INFO: Commands

func commandExit(config *config, params []string) error {
	fmt.Printf("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *config, params []string) error {
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	getCommandsDescriptions()
	fmt.Println("")
	return nil
}

func commandMap(config *config, params []string) error {
	err := printLocations(config, PokeConfig.next)
	if err != nil {
		return err
	}
	return nil
}

func commandMap2(config *config, params []string) error {
	err := printLocations(config, PokeConfig.previous)
	if err != nil {
		return err
	}
	return nil
}

func commandExplore(config *config, params []string) error {
	err := printPokemonFromArea(params[0])
	if err != nil {
		return err
	}
	return nil
}

// INFO: Additional functions

// For MAP and MAPB commands
func printLocations(config *config, url string) error {
	if url == "" {
		fmt.Println("There are no result.")
	} else {
		resp, err := pokeapi.GetLocations(url, PokeConfig.cache)
		if err != nil {
			return err
		}
		locations, err := pokeapi.UnmarshalLocations(resp)
		if err != nil {
			return err
		}
		config.next = locations.Next
		config.previous = locations.Previous
		for _, v := range locations.Results {
			fmt.Println(v.Name)
		}
	}
	return nil
}

// For EXPLORE command
func printPokemonFromArea(name string) error {
	resp, err := pokeapi.ExploreLocation(name, PokeConfig.cache)
	if err != nil {
		return err
	}
	singleLocation, err := pokeapi.UnmarshalSingleLocation(resp)
	if err != nil {
		return err
	}
	fmt.Printf("Exploring %s...\n", singleLocation.Name)
	fmt.Println("Found Pokemon:")
	for _, pokemon := range singleLocation.PokemonEncounters {
		fmt.Println("-", pokemon.Pokemon.Name)
	}
	return nil
}

// For HELP commands
func getCommandsDescriptions() {
	for _, value := range commandRegistry {
		fmt.Printf("%v: %v\n", value.name, value.description)
	}
}

// For Main REPL loop
func cleanInput(text string) []string {
	if text == "" {
		return []string{}
	}
	reg := regexp.MustCompile("[^a-z0-9- ]+")
	wordsList := strings.Fields(reg.ReplaceAllString(strings.ToLower(text), ""))
	return wordsList
}
