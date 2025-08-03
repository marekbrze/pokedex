package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/marekbrze/pokedexcli/internal/pokeapi"
)

// INFO: Main types

type config struct {
	next     string
	previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
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

// INFO: Main Loop
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

// INFO: Commands

func commandExit(config *config) error {
	fmt.Printf("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *config) error {
	fmt.Println("\nWelcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	getCommandsDescriptions()
	fmt.Println("")
	return nil
}

func commandMap(config *config) error {
	err := printLocations(config, PokeConfig.next)
	if err != nil {
		return err
	}
	return nil
}

func commandMap2(config *config) error {
	err := printLocations(config, PokeConfig.previous)
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
		resp, err := pokeapi.GetLocations(url)
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
	reg := regexp.MustCompile("[^a-z ]+")
	wordsList := strings.Fields(reg.ReplaceAllString(strings.ToLower(text), ""))
	return wordsList
}
