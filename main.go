package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/almasx/pokedexcli/internal/pokecache"
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

func commandExit(config *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(config *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	fmt.Println("help - Displays a help message")
	fmt.Println("exit - Exit the Pokedex")
	return nil
}


type GetLocationAreas struct {
	Count int `json:"count"`
	Next string `json:"next"`
	Previous string `json:"previous"`
	Results []struct {
		Name string `json:"name"`
		Url string `json:"url"`
	} `json:"results"`
}

const BASE_URL = "https://pokeapi.co/api/v2/location-area/?offset=0&limit=20"

func fetchMapData(url string, config *config) (GetLocationAreas, error) {
	res := GetLocationAreas{}

	if data, ok := config.cache.Get(url); ok {
		fmt.Println("---- Getting from cache ----", url)
		err := json.Unmarshal(data, &res)
		if err != nil {
			return res, err
		}
		return res, nil
	}
	
	data, err := http.Get(url)
	if err != nil {
		return GetLocationAreas{}, err
	}
	defer data.Body.Close()
	body, err := io.ReadAll(data.Body)
	if err != nil {
		return GetLocationAreas{}, err
	}
	config.cache.Add(url, body)
	err = json.Unmarshal(body, &res)
	if err != nil {
		return GetLocationAreas{}, err
	}
	return res, nil
}

func commandMap(config *config, args []string) error {
	mapData := GetLocationAreas{}
	url := config.next
	if url == "" {
		url = BASE_URL
	}
	
	mapData, err := fetchMapData(url, config)
	if err != nil {
		return err
	}

	for _, result := range mapData.Results {
		fmt.Println(result.Name)
	}
	
	config.next = mapData.Next
	config.prev = mapData.Previous

	return nil
}

func commandMapb(config *config, args []string) error {
	mapData := GetLocationAreas{}
	url := config.prev
	if url == "" {
		fmt.Println("you're on the first page")
		return nil
	}
	
	mapData, err := fetchMapData(url, config)
	if err != nil {
		return err
	}

	for _, result := range mapData.Results {
		fmt.Println(result.Name)
	}
	
	config.prev = mapData.Previous
	config.next = mapData.Next

	return nil
}

type GetLocationAreaPokemons struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int           `json:"min_level"`
				MaxLevel        int           `json:"max_level"`
				ConditionValues []interface{} `json:"condition_values"`
				Chance          int           `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func fetchLocationAreaPokemons(url string, config *config) (GetLocationAreaPokemons, error) {
	res := GetLocationAreaPokemons{}

	if data, ok := config.cache.Get(url); ok {
		fmt.Println("---- Getting from cache ----", url)
		err := json.Unmarshal(data, &res)
		if err != nil {
			return res, err
		}
		return res, nil
	}
	
	data, err := http.Get(url)
	if err != nil {
		return GetLocationAreaPokemons{}, err
	}
	defer data.Body.Close()
	body, err := io.ReadAll(data.Body)
	if err != nil {
		return GetLocationAreaPokemons{}, err
	}
	config.cache.Add(url, body)
	err = json.Unmarshal(body, &res)
	if err != nil {
		return GetLocationAreaPokemons{}, err
	}
	return res, nil
}

func commandExplore(config *config, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("explore requires a location area")
	}
	location_area := args[0]
	if location_area == "" {
		return fmt.Errorf("location area is required")	
	}

	fmt.Println("Exploring", location_area, "...")

	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s", location_area)
	location_area_pokemons, err := fetchLocationAreaPokemons(url, config)
	if err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, pokemon := range location_area_pokemons.PokemonEncounters {
		fmt.Println(" - ", pokemon.Pokemon.Name)
	}

	return nil
}

type cliCommand struct {
	name        string
	description string
	callback     func(*config, []string) error
}

type config struct {
	next string
	prev string
	cache *pokecache.Cache
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
		callback:    commandMap,
	},

	"mapb": {
		name:        "mapb",
		description: "Show the previous map",
		callback:    commandMapb,
	},
	"explore": {
		name:        "explore",
		description: "Explore the map",
		callback:    commandExplore,
	},
}


func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cache := pokecache.NewCache(time.Second * 10)
	var user_input string

	config := config{
		next: "",
		prev: "",
		cache: cache,
	}

	for {
		fmt.Printf("Pokedex > ")
		scanner.Scan()
		user_input = scanner.Text()
		words := cleanInput(user_input)

		if len(words) == 0 {			continue
		}
		command, exists := commands[strings.ToLower(words[0])]
		if !exists {
			fmt.Println("Unknown command")
			continue
		}
		command.callback(&config, words[1:])
	}
}	
