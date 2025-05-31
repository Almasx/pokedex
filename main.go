package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/almasx/pokedexcli/internal/api"
	"github.com/almasx/pokedexcli/internal/cli"
	explorepkg "github.com/almasx/pokedexcli/internal/explore"
	mappkg "github.com/almasx/pokedexcli/internal/map"
	"github.com/almasx/pokedexcli/internal/pokecache"
	"github.com/almasx/pokedexcli/internal/pokemon"
)

func cleanInput(text string) []string {
	var res []string
	for _, value := range strings.Split(text, " ") {
		if value != "" {
			res = append(res, value)
		}
	}
	return res
}

func commandExit(config *cli.Config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *cli.Config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("help - Displays a help message")
	fmt.Println("exit - Exit the Pokedex")
	fmt.Println("map - Show the map of the Pokemon world")
	fmt.Println("mapb - Show the previous page of the map")
	fmt.Println("explore <location_area> - Explore a location area")
	fmt.Println("catch <pokemon> - Catch a pokemon")
	fmt.Println("inspect <pokemon> - Inspect a pokemon")
	fmt.Println("pokedex - Show the pokedex")
	return nil
}

type cliCommand struct {
	name        string
	description string
	callback    func(*cli.Config, []string) error
}

var commands = map[string]cliCommand{
	"exit": {
		name:        "exit",
		description: "Exit the Pokedex",
		callback:    commandExit,
	},
	"help": {
		name:        "help",
		description: "Show all commands",
		callback:    commandHelp,
	},
	"map": {
		name:        "map",
		description: "Show the map",
		callback:    mappkg.CommandMap,
	},
	"mapb": {
		name:        "mapb",
		description: "Show the previous map",
		callback:    mappkg.CommandMapb,
	},
	"explore": {
		name:        "explore",
		description: "Explore the map",
		callback:    explorepkg.CommandExplore,
	},
	"catch": {
		name:        "catch",
		description: "Catch a pokemon",
		callback:    pokemon.CommandCatch,
	},
	"inspect": {
		name:        "inspect",
		description: "Inspect a pokemon",
		callback:    pokemon.CommandInspect,
	},
	"pokedex": {
		name:        "pokedex",
		description: "Show the pokedex",
		callback:    pokemon.CommandPokedex,
	},
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cache := pokecache.NewCache(time.Second * 10)
	var user_input string

	config := cli.Config{
		Next:    "",
		Prev:    "",
		Cache:   cache,
		Pokedex: make(map[string]api.GetPokemon),
	}

	for {
		fmt.Printf("Pokedex > ")
		scanner.Scan()
		user_input = scanner.Text()
		words := cleanInput(user_input)

		if len(words) == 0 {
			continue
		}
		command, exists := commands[strings.ToLower(words[0])]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}
		command.callback(&config, words[1:])
	}
}	
