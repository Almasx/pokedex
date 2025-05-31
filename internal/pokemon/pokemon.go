package pokemon

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"

	"github.com/almasx/pokedexcli/internal/api"
	"github.com/almasx/pokedexcli/internal/cli"
)



func fetchPokemon(url string, config *cli.Config) (api.GetPokemon, error) {
	res := api.GetPokemon{}

	if data, ok := config.Cache.Get(url); ok {
		err := json.Unmarshal(data, &res)
		if err != nil {
			return res, err
		}
		return res, nil
	}
	
	data, err := http.Get(url)
	if err != nil {
		return api.GetPokemon{}, err
	}
	defer data.Body.Close()
	body, err := io.ReadAll(data.Body)
	if err != nil {
		return api.GetPokemon{}, err
	}
	config.Cache.Add(url, body)
	err = json.Unmarshal(body, &res)
	if err != nil {
		return api.GetPokemon{}, err
	}
	return res, nil
}

func catchPokemon(pokemon_data api.GetPokemon) bool {
	catch_rate := pokemon_data.BaseExperience 
	random_number := rand.Intn(catch_rate * 2)
	return random_number <= catch_rate
}

func CommandCatch(config *cli.Config, args []string) error {
	if len(args) != 1 {
		fmt.Println("catch requires a pokemon")
		return fmt.Errorf("catch requires a pokemon")
	}
	pokemon := args[0]
	if pokemon == "" {
		fmt.Println("pokemon is required")
		return fmt.Errorf("pokemon is required")
	}

	if _, ok := config.Pokedex[pokemon]; ok {
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
		fmt.Println("You may now inspect it with the inspect command.")
		config.Pokedex[pokemon] = pokemon_data
	} else {
		fmt.Println( pokemon, "escaped!")
	}

	return nil
}

func CommandInspect(config *cli.Config, args []string) error {
	if len(args) != 1 {
		fmt.Println("inspect requires a pokemon")
		return fmt.Errorf("inspect requires a pokemon")
	}
	pokemon := args[0]
	if pokemon == "" {
		fmt.Println("pokemon is required")
		return fmt.Errorf("pokemon is required")
	}

	if _, ok := config.Pokedex[pokemon]; !ok {
		fmt.Println("you have not caught that pokemon")
		return fmt.Errorf("you have not caught that pokemon")
	}

	fmt.Printf("Name: %v\n", config.Pokedex[pokemon].Name)
	fmt.Printf("Height: %v\n", config.Pokedex[pokemon].Height)
	fmt.Printf("Weight: %v\n", config.Pokedex[pokemon].Weight)

	fmt.Printf("Stats: \n")
	for _, stat := range config.Pokedex[pokemon].Stats {
		fmt.Printf("  -%v: %v\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Printf("Types: \n")
	for _, type_ := range config.Pokedex[pokemon].Types {
		fmt.Printf("  - %v\n", type_.Type.Name)
	}
	
	return nil
}

func CommandPokedex(config *cli.Config, args []string) error {	
	fmt.Println("Your Pokedex:")
	for _, pokemon := range config.Pokedex {
		fmt.Printf("  - %v\n", pokemon.Name)
	}
	return nil
}
