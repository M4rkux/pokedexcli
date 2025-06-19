package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	next     string
	previous string
}

type cliCommand struct {
	name        string
	description string
	callback    func([]string) error
	config      Config
}

func main() {
	var commands map[string]cliCommand
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback: func(input []string) error {
				return commandExit()
			},
			config: Config{},
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback: func(input []string) error {
				return commandHelp(&commands)
			},
			config: Config{},
		},
		"map": {
			name:        "map",
			description: "Displays the names of the location areas",
			callback: func(input []string) error {
				return commandMap(&commands)
			},
			config: Config{},
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the names of the location areas (previous)",
			callback: func(input []string) error {
				return commandMapBack(&commands)
			},
			config: Config{},
		},
		"explore": {
			name:        "explore",
			description: "Explore the area and list the Pokemons",
			callback: func(input []string) error {
				if len(input) < 2 {
					return errors.New("Missing required parameter (area)")
				}
				return commandExplore(input[1])
			},
			config: Config{},
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())

		if len(input) == 0 {
			continue
		}

		if command, ok := commands[input[0]]; ok {
			err := command.callback(input)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Println("Unknown command")
		}
	}
}

func cleanInput(text string) []string {
	trimmed := strings.TrimSpace(text)
	lowered := strings.ToLower(trimmed)
	words := strings.Fields(lowered)
	return words
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(commands *map[string]cliCommand) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()
	for _, cmd := range *commands {
		fmt.Println(fmt.Sprintf("%s: %s", cmd.name, cmd.description))
	}
	return nil
}

func commandMap(commands *map[string]cliCommand) error {
	locationAreas, err := GetListLocationAreas((*commands)["map"].config.next)
	if err != nil {
		return err
	}

	mapCmd := (*commands)["map"]
	mapCmd.config = Config{
		next:     locationAreas.Next,
		previous: locationAreas.Previous,
	}
	(*commands)["map"] = mapCmd

	for _, locationArea := range locationAreas.Results {
		fmt.Println(locationArea.Name)
	}

	return nil
}

func commandMapBack(commands *map[string]cliCommand) error {
	if (*commands)["map"].config.previous == "" {
		fmt.Println("you're on the first page")
		return nil
	}

	locationAreas, err := GetListLocationAreas((*commands)["map"].config.previous)
	if err != nil {
		return err
	}

	mapCmd := (*commands)["map"]
	mapCmd.config = Config{
		next:     locationAreas.Next,
		previous: locationAreas.Previous,
	}
	(*commands)["map"] = mapCmd

	for _, locationArea := range locationAreas.Results {
		fmt.Println(locationArea.Name)
	}

	return nil
}

func commandExplore(area string) error {

	fmt.Println(fmt.Sprintf("Exploring %s...", area))

	pokemonEncounters, err := GetListPokemonsInArea(area)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")

	for _, pokemonEncounter := range pokemonEncounters {
		fmt.Println(" -", pokemonEncounter.Pokemon.Name)
	}

	return nil
}
