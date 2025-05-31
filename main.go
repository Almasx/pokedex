package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
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
		fmt.Println("explore requires a location area")
		return fmt.Errorf("explore requires a location area")
	}
	location_area := args[0]
	if location_area == "" {
		fmt.Println("location area is required")
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

type GetPokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
	Abilities      []struct {
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
		Ability  struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
	} `json:"abilities"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	HeldItems []struct {
		Item struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"item"`
		VersionDetails []struct {
			Rarity  int `json:"rarity"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"held_items"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt int `json:"level_learned_at"`
			VersionGroup   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			Order int `json:"order"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Species struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string      `json:"back_default"`
		BackFemale       interface{} `json:"back_female"`
		BackShiny        string      `json:"back_shiny"`
		BackShinyFemale  interface{} `json:"back_shiny_female"`
		FrontDefault     string      `json:"front_default"`
		FrontFemale      interface{} `json:"front_female"`
		FrontShiny       string      `json:"front_shiny"`
		FrontShinyFemale interface{} `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string      `json:"front_default"`
				FrontFemale  interface{} `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string      `json:"front_default"`
				FrontFemale      interface{} `json:"front_female"`
				FrontShiny       string      `json:"front_shiny"`
				FrontShinyFemale interface{} `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string      `json:"back_default"`
				BackFemale       interface{} `json:"back_female"`
				BackShiny        string      `json:"back_shiny"`
				BackShinyFemale  interface{} `json:"back_shiny_female"`
				FrontDefault     string      `json:"front_default"`
				FrontFemale      interface{} `json:"front_female"`
				FrontShiny       string      `json:"front_shiny"`
				FrontShinyFemale interface{} `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault  string `json:"back_default"`
					BackGray     string `json:"back_gray"`
					FrontDefault string `json:"front_default"`
					FrontGray    string `json:"front_gray"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault  string `json:"back_default"`
					BackGray     string `json:"back_gray"`
					FrontDefault string `json:"front_default"`
					FrontGray    string `json:"front_gray"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"crystal"`
				Gold struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"gold"`
				Silver struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string      `json:"back_default"`
					BackFemale       interface{} `json:"back_female"`
					BackShiny        string      `json:"back_shiny"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string      `json:"back_default"`
					BackFemale       interface{} `json:"back_female"`
					BackShiny        string      `json:"back_shiny"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string      `json:"back_default"`
					BackFemale       interface{} `json:"back_female"`
					BackShiny        string      `json:"back_shiny"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string      `json:"back_default"`
						BackFemale       interface{} `json:"back_female"`
						BackShiny        string      `json:"back_shiny"`
						BackShinyFemale  interface{} `json:"back_shiny_female"`
						FrontDefault     string      `json:"front_default"`
						FrontFemale      interface{} `json:"front_female"`
						FrontShiny       string      `json:"front_shiny"`
						FrontShinyFemale interface{} `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string      `json:"back_default"`
					BackFemale       interface{} `json:"back_female"`
					BackShiny        string      `json:"back_shiny"`
					BackShinyFemale  interface{} `json:"back_shiny_female"`
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string      `json:"front_default"`
					FrontFemale  interface{} `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string      `json:"front_default"`
					FrontFemale      interface{} `json:"front_female"`
					FrontShiny       string      `json:"front_shiny"`
					FrontShinyFemale interface{} `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string      `json:"front_default"`
					FrontFemale  interface{} `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Cries struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	PastTypes []struct {
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
		Types []struct {
			Slot int `json:"slot"`
			Type struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"type"`
		} `json:"types"`
	} `json:"past_types"`
	PastAbilities []struct {
		Generation struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"generation"`
		Abilities []struct {
			Ability  interface{} `json:"ability"`
			IsHidden bool        `json:"is_hidden"`
			Slot     int         `json:"slot"`
		} `json:"abilities"`
	} `json:"past_abilities"`
}

func fetchPokemon(url string, config *config) (GetPokemon, error) {
	res := GetPokemon{}

	if data, ok := config.cache.Get(url); ok {
		err := json.Unmarshal(data, &res)
		if err != nil {
			return res, err
		}
		return res, nil
	}
	
	data, err := http.Get(url)
	if err != nil {
		return GetPokemon{}, err
	}
	defer data.Body.Close()
	body, err := io.ReadAll(data.Body)
	if err != nil {
		return GetPokemon{}, err
	}
	config.cache.Add(url, body)
	err = json.Unmarshal(body, &res)
	if err != nil {
		return GetPokemon{}, err
	}
	return res, nil
}


func catchPokemon(pokemon_data GetPokemon) bool {
	catch_rate := pokemon_data.BaseExperience 
	random_number := rand.Intn(catch_rate * 2)
	return random_number <= catch_rate
}

func commandCatch(config *config, args []string) error {
	if len(args) != 1 {
		fmt.Println("catch requires a pokemon")
		return fmt.Errorf("catch requires a pokemon")
	}
	pokemon := args[0]
	if pokemon == "" {
		fmt.Println("pokemon is required")
		return fmt.Errorf("pokemon is required")
	}

	if _, ok := config.pokedex[pokemon]; ok {
		fmt.Println("pokemon already in pokedex")
		return fmt.Errorf("pokemon already in pokedex")
	}


	fmt.Printf("Throwing a Pokeball at %v...\n", pokemon)	
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s", pokemon)
	pokemon_data, err := fetchPokemon(url, config)

	if err != nil {
		return err
	}

	caught := catchPokemon(pokemon_data)
	if caught {
		fmt.Println( pokemon, "was caught!")
		config.pokedex[pokemon] = pokemon_data
	} else {
		fmt.Println( pokemon, "escaped!")
	}

	return nil
}

func commandInspect(config *config, args []string) error {
	if len(args) != 1 {
		fmt.Println("inspect requires a pokemon")
		return fmt.Errorf("inspect requires a pokemon")
	}
	pokemon := args[0]
	if pokemon == "" {
		fmt.Println("pokemon is required")
		return fmt.Errorf("pokemon is required")
	}

	if _, ok := config.pokedex[pokemon]; !ok {
		fmt.Println("you have not caught that pokemon")
		return fmt.Errorf("you have not caught that pokemon")
	}

	fmt.Printf("Name: %v\n", config.pokedex[pokemon].Name)
	fmt.Printf("Height: %v\n", config.pokedex[pokemon].Height)
	fmt.Printf("Weight: %v\n", config.pokedex[pokemon].Weight)

	fmt.Printf("Stats: \n")
	for _, stat := range config.pokedex[pokemon].Stats {
		fmt.Printf("  -%v: %v\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Printf("Types: \n")
	for _, type_ := range config.pokedex[pokemon].Types {
		fmt.Printf("  - %v\n", type_.Type.Name)
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
	pokedex map[string]GetPokemon
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
	"catch": {
		name:        "catch",
		description: "Catch a pokemon",
		callback:    commandCatch,
	},
	"inspect": {
		name:        "inspect",
		description: "Inspect a pokemon",
		callback:    commandInspect,
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
		pokedex: make(map[string]GetPokemon),
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
