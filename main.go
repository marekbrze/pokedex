package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/marekbrze/pokedexcli/internal/pokeapi"
)

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

func init() {
	commandRegistry["map"] = cliCommand{
		name:        "map",
		description: "Displays pokemon locations",
		callback:    commandMap,
	}
	// commandRegistry["mapb"] = cliCommand{
	// 	name:        "mapb",
	// 	description: "Displays previous locations page",
	// 	callback:    commandMap2,
	// }
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

func getCommandsDescriptions() {
	for _, value := range commandRegistry {
		fmt.Printf("%v: %v\n", value.name, value.description)
	}
}

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
	if PokeConfig.next == "" {
		fmt.Println("There are no result.")
	} else {
		fmt.Println("test")
		locations, err := pokeapi.GetLocations(PokeConfig.next)
		fmt.Println(locations.Count)
		if err != nil {
			fmt.Println("error 4")
			return err
		}
		fmt.Println(len(locations.Results))
		for _, v := range locations.Results {
			fmt.Println(v.Name)
		}
		fmt.Println("End test")
	}
	return nil
}

// func commandMap2(config *config) error {
// 	var err error
// 	if config.previous == "" {
// 		fmt.Println("There are no result.")
// 	} else {
// 		locations, err := pokeapi.GetLocations(config.previous)
// 	}
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
