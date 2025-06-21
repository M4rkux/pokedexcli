package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
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

var currentLocation string
var pokedex map[string]Pokemon

func main() {
	currentLocation = "canalave-city-area"
	pokedex = make(map[string]Pokemon)
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
			name:        "explore <area>",
			description: "Explore the area and list the Pokemons",
			callback: func(input []string) error {
				if len(input) < 2 {
					return errors.New("Missing required parameter (area)")
				}
				return commandExplore(input[1])
			},
			config: Config{},
		},
		"catch": {
			name:        "catch <pokemon>",
			description: "Throw a pokeball and tries to catch a pokemon",
			callback: func(input []string) error {
				if len(input) < 2 {
					return errors.New("Missing required parameter (pokemon)")
				}
				return commandCatch(input[1])
			},
			config: Config{},
		},
		"inspect": {
			name:        "inspect <pokemon>",
			description: "Shows the name, height, weight, stats and type(s) of the Pokemon",
			callback: func(input []string) error {
				if len(input) < 2 {
					return errors.New("Missing required parameter (pokemon)")
				}
				return commandInspect(input[1])
			},
			config: Config{},
		},
		"pokedex": {
			name:        "pokedex",
			description: "Shows the list of pokemons you already got",
			callback: func(input []string) error {
				return commandPokedex()
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
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
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

	fmt.Printf("Exploring %s...\n", area)

	pokemonEncounters, err := GetListPokemonsInArea(area)
	if err != nil {
		return err
	}

	currentLocation = area
	fmt.Println("Found Pokemon:")

	for _, pokemonEncounter := range pokemonEncounters {
		fmt.Println(" -", pokemonEncounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(pokemonName string) error {
	pokemonLocationArea, err := GetPokemonLocationArea(pokemonName)
	if err != nil {
		return err
	}

	isInArea := false
	for _, locationArea := range pokemonLocationArea {
		if locationArea.LocationArea.Name == currentLocation {
			isInArea = true
			break
		}
	}

	if !isInArea {
		return errors.New("pokemon not in area")
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	pokemon, err := GetPokemon(pokemonName)
	if err != nil {
		return err
	}

	captureProbability := rand.Intn(100)
	chance := GetChanceByBaseExperience(pokemon.BaseExperience)

	if captureProbability > chance {
		pokedex[pokemonName] = pokemon
		fmt.Printf("%s was caught!\n", pokemonName)
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}

func commandInspect(pokemonName string) error {
	pokemon, ok := pokedex[pokemonName]
	if !ok {
		return fmt.Errorf("%s is not caught yet!", pokemonName)
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %v\n", pokemon.Height)
	fmt.Printf("Weight: %v\n", pokemon.Weight)

	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("    -%s: %v\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf("    - %s\n", t.Type.Name)
	}

	return nil
}

func commandPokedex() error {
	if len(pokedex) == 0 {
		fmt.Println("Your pokedex is empty")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for pokemonName := range pokedex {
		fmt.Printf(" - %s\n", pokemonName)
	}

	return nil
}

func GetChanceByBaseExperience(baseExperience int) int {
	if baseExperience < 100 {
		return 80
	} else if baseExperience < 200 {
		return 60
	}
	return 40
}
