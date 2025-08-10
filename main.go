package main

import (
	"bufio"
	"fmt"
	"math/rand"
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
	pokedex  map[string]pokeapi.Pokemon
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
	commandRegistry["pokedex"] = cliCommand{
		name:        "pokedex",
		description: "Prints your pokedex",
		callback:    printPokedex,
	}
	commandRegistry["catch"] = cliCommand{
		name:        "catch",
		description: "Tries to catch selected Pokemon",
		callback:    commandCatch,
	}
	commandRegistry["inspect"] = cliCommand{
		name:        "inspect",
		description: "inspect caught pokemon",
		callback:    commandInspect,
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
	PokeConfig.pokedex = map[string]pokeapi.Pokemon{}
}

// INFO: Main Loop
func main() {
	scanner := bufio.NewScanner(os.Stdin)
	PokeConfig.cache = pokecache.NewCache(15 * time.Second)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		cleanedInput := cleanInput(line)
		if len(cleanedInput) == 0 {
			continue
		}
		firstCommand := cleanedInput[0]
		command, exists := commandRegistry[firstCommand]
		if exists && len(cleanedInput) > 1 {
			err := command.callback(&PokeConfig, cleanedInput[1:])
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
		} else if exists {
			err := command.callback(&PokeConfig, []string{})
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
		} else {
			fmt.Println("Unknown command")
			continue
		}
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

func commandCatch(config *config, params []string) error {
	wantedPokemon := params[0]
	_, exist := PokeConfig.pokedex[wantedPokemon]
	if exist {
		fmt.Println("You already caught this pokemon!")
	} else {
		res, err := pokeapi.GetPokemon(wantedPokemon, PokeConfig.cache)
		if err != nil {
			return err
		}
		pokemon, err := pokeapi.UnmarshalPokemon(res)
		if err != nil {
			return err
		}
		userChance := rand.Intn(500)
		fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
		if userChance > pokemon.BaseExperience {
			PokeConfig.pokedex[pokemon.Name] = pokemon
			fmt.Printf("Congrats! You caught %s!\n", pokemon.Name)
		} else {
			fmt.Println(pokemon.Name, "got away!")
		}
	}
	return nil
}

func commandInspect(config *config, params []string) error {
	inspectedPokemon := params[0]
	pokemon, exists := config.pokedex[inspectedPokemon]
	if exists {
		fmt.Println("Name:", pokemon.Name)
		fmt.Println("Height:", pokemon.Height)
		fmt.Println("Weight:", pokemon.Weight)
		fmt.Println("Stats:")
		for _, v := range pokemon.Stats {
			fmt.Printf("    -%s: %d\n", v.Stat.Name, v.BaseStat)
		}
		fmt.Println("Types:")
		for _, v := range pokemon.Types {
			fmt.Printf("    -%s\n", v.Type.Name)
		}
	} else {
		fmt.Println("you have not caught that pokemon")
	}
	return nil
}

func printPokedex(config *config, params []string) error {
	if len(config.pokedex) == 0 {
		fmt.Println("You haven't caught any pokemon!")
	} else {
		fmt.Println("Your Pokedex:")
		for k := range config.pokedex {
			fmt.Printf("    - %s\n", k)
		}
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
