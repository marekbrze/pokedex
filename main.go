package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/marekbrze/pokedexcli/internal/pokeapi"
)

type Config struct {
	next     string
	previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

var (
	PokeConfig      Config
	commandRegistry = make(map[string]cliCommand)
)

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
	PokeConfig.next = "https://pokeapi.co/api/v2/location-area/"
	PokeConfig.previous = ""
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
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
		if exists {
			err := command.callback(&PokeConfig)
			if err != nil {
				fmt.Print(fmt.Errorf("error: %w", err))
			}
		} else {
			fmt.Println("Unknown command")
		}
		fmt.Print("Pokedex > ")
	}
}

func cleanInput(text string) []string {
	if text == "" {
		return []string{}
	}
	reg := regexp.MustCompile("[^a-z ]+")
	wordsList := strings.Fields(reg.ReplaceAllString(strings.ToLower(text), ""))
	return wordsList
}

// INFO: Commands

func commandExit(config *Config) error {
	fmt.Printf("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *Config) error {
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	getCommandsDescriptions()
	fmt.Println("")
	return nil
}

func commandMap(config *Config) error {
	err := printLocations(config, PokeConfig.next)
	if err != nil {
		return err
	}
	return nil
}

func commandMap2(config *Config) error {
	err := printLocations(config, PokeConfig.previous)
	if err != nil {
		return err
	}
	return nil
}

// INFO: Additional functions

// For MAP and MAPB commands
func printLocations(config *Config, url string) error {
	if url == "" {
		fmt.Println("There are no result.")
	} else {
		locations, err := pokeapi.GetLocations(url)
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

// For HELP commands
func getCommandsDescriptions() {
	for _, value := range commandRegistry {
		fmt.Printf("%v: %v\n", value.name, value.description)
	}
}
